/*

	Mesostic Fetch Automator

	A static URL is configured for fetching a source text via API to be transmogrified into a Mesostic for display.

*/

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
)

// Channel for NASAapod Mesostic publishing
var nasaNewMESO = make(chan string)

// generic crontab emulator
func fetchCron() {
	// Job for fetching a new source
	fcron := gocron.NewScheduler(time.UTC)
	_, ferr := fcron.Every(30).Seconds().StartImmediately().Do(NASAetl)
	if ferr != nil {
		log.Error()
	}
	defer fcron.StartBlocking()

	// Job for displaying the new mesostic
	// This is intended as a 'heartbeat' with the side effect of unblocking nasaNewMESO
	/*
		dcron := gocron.NewScheduler(time.UTC)
		_, derr := dcron.Every(13).Seconds().StartImmediately().Do(NASAetl)
		if derr != nil {
			log.Error()
		}

		defer dcron.StartBlocking()
	*/
}

// NASAetl ::: Retrieve the Astronomy Picture of the Day (APOD) metadata and process it through the Mesostic engine.
// This is where the URL is configured.
//
func NASAetl() {
	_, _, fu := Envelope()

	fmt.Println("NASAetl Running")

	// NASA official Astronomy Picture of the Day endpoint URL using a freely available API key
	url := "https://api.nasa.gov/planetary/apod?api_key=Ijb0zLeEt71HMQdy8YjqB583FK3bdh1yThVJYzpu"
	// url := "https://api.nasa.gov/planetary/apod?date=2000-01-01&api_key=Ijb0zLeEt71HMQdy8YjqB583FK3bdh1yThVJYzpu"

	// the title as the spine, for now :)
	date, spine, source := fetchSource(url)

	// The most up-to-date image can throw a 404 towards the end of the day.
	// This is meant to provide a fall-back that displays a random one in the past.
	// Which might expand this conditional, or become a switch statement.
	if spine == "404" {
		log.Error().Str("code", "404").Msg("Remote data not available, deploying fallback [coming soon]")
	}

	// we don't want spaces in the spine string
	trcc := strings.NewReplacer(" ", "")
	spn := trcc.Replace(spine)

	// convert each phrase into a line by replacing commas and periods with newlines.
	trnl := strings.NewReplacer(". ", "\n", ", ", "\n")
	source = trnl.Replace(source)

	// get a mesostic
	// this mimics the JSON API calls
	// which will probably need to be revisited once this section is done
	tmpFileName := fileTmp(&spn, &source)
	mcMeso := make(chan string)
	go mesoMain(tmpFileName, spn, mcMeso)
	showR := <-mcMeso

	// create new Mesostic file
	// this *might* be better *before* the mesostic is created?
	// but then, how do we know to send it data?
	mesoFile, created := apodNew(&spine, &date, &showR)
	if !created {
		fmt.Printf("Entry exists at '%s', skipping file creation for '%s'\n", mesoFile, spine)
	}

	// remove the tmp source file
	var ferr = os.Remove(tmpFileName)
	if ferr != nil {
		log.Error()
	}

	/*

		I have tested and seen that what is being produced can be displayed.
		The problems are:
		1. Providing a filename to the homepage loader
		2. Formatting correctly in HTML... in plaintext, e.g. curl, the string formatting remains.
		3. Clearing the channel of new filenames

		So, similar to the mesostic files themselves,
		the filenames will be kept "as a database" that the homepage function can check,
		and select with chance operations.

		But the channel method could still be valuable and faster.
		So it will be kept... maybe a heartbeat cronjob makes sense?

		Making it a buffered channel would also help, but increase memory a bit.

	*/

	// this is BLOCKING, but it works if the homepage is accessed
	// when NOT accessed, this will continue to block, but run, and pile up mesostics in the channel
	// additionally, the tmp file is never deleted
	nasaNewMESO <- mesoFile // push filename of new Mesostic

	// The blocking happening here is maybe important to deal with.
	// The channel will "fill up", most probably, eating up memory?
	// Or maybe if the file exists, homepage() doesn't need to know?
	// this line is only reached when the channel is unblocked by accessing homepage()

	log.Info().
		Str("fu", fu).
		Str("date", date).
		Str("spinestring", spine).
		Str("filename", mesoFile).
		// Str("mesostic", showR). // leave for DEBUG mode
		Msg("NASA APOD Mesostic complete")
}

// apodNEW ::: Check if a disk file exists in the Mesostic store or create a new one.
// The return values are the filename and whether the function wrote a new file.
func apodNew(sp *string, da *string, me *string) (string, bool) {
	mDir := "store"
	tr := strings.NewReplacer(" ", "_")
	spn := tr.Replace(*sp)
	fP := fmt.Sprintf("%s/%s__%s", mDir, *da, spn)

	// TODO: check how old it is and update if a retention period is met
	if _, err := os.Stat(fP); err == nil {
		return fP, false
	}

	// write data to a new file if it doesn't exist
	sB := []byte(*me)
	err := ioutil.WriteFile(fP, sB, 0644)
	if err != nil {
		log.Error()
	}
	return fP, true
}

// fileTmp ::: Take a source string and place it in a file name after the spinestring.
// This only creates the file by a straight byte copy.
// Calling functions are responsible for file deletion when finished.
func fileTmp(sp *string, so *string) string {
	fT := time.Now()
	fS := fT.Unix()
	fN := fmt.Sprintf("%s__%d", *sp, fS)
	sB := []byte(*so)
	err := ioutil.WriteFile(fN, sB, 0644)
	if err != nil {
		log.Error()
	}
	return fN
}

// Envelope ::: Returns details about the current code execution point.
// This enables tracing in log events, for instance:
//		_, _, fu := Envelope()
//		fmt.Printf("current function: %s", fu)
func Envelope() (string, int, string) {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return frame.File, frame.Line, frame.Function
}
