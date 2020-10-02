/*

	Mesostic Scheduler and Tasks

*/

package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
)

// Channel for NASAapod Mesostic publishing
var nasaNewMESO = make(chan string)

// fetchCron ::: Cron emulator, currently just a single job.
func fetchCron() {
	// NASA official Astronomy Picture of the Day endpoint URL using a freely available API key
	apodnow := "https://api.nasa.gov/planetary/apod?api_key=Ijb0zLeEt71HMQdy8YjqB583FK3bdh1yThVJYzpu"
	apodenv := "HPSCHD_NASA_APOD_URL" // Optional ENV VAR
	url := envVar(apodenv, apodnow)   // NASA APOD URL to query, default if no ENV VAR
	var afreq uint64 = 900            // Frequency (s) to check

	// Start a new fetch job immediately, followed every afreq seconds.
	fcron := gocron.NewScheduler(time.UTC)
	_, ferr := fcron.Every(afreq).Seconds().StartImmediately().Do(NASAetl, url)
	if ferr != nil {
		log.Error()
	}
	defer fcron.StartBlocking()

	// TODO: add another job to randomize the source with dates since 2000-01-01
	// this will help populate the datastore to give homepage() more random choices
}

// NASAetl ::: Retrieve Astronomy Picture of the Day (APOD) metadata,
// process it through the Mesostic engine, save it in a library of ephemeral copies,
// pass the new data point (filename path) to a channel for use with displays.
func NASAetl(url string) {
	_, _, fu := Envelope()

	fmt.Println("NASAetl Running")

	// the title as the spine, for now :)
	date, spine, source := fetchSource(url)

	// The most up-to-date image can throw a 404 towards the end of the day.
	// For now this triggers the fetch of 2004-04-04
	if spine == "404" {
		// TODO: Have this trigger a random fetch (not yet implemented) to populate the url argument
		url := "https://api.nasa.gov/planetary/apod?date=2004-04-04&api_key=Ijb0zLeEt71HMQdy8YjqB583FK3bdh1yThVJYzpu"

		go NASAetl(url)
		log.Warn().
			Str("fu", fu).
			Str("code", "404").
			Msg("Remote data not available, alternate ETL triggered.")

		// there's no need for this ETL to continue
		return
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
		// this would be a good place to issue a random date display
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
		Msg("NASA APOD Mesostic complete")

	log.Debug().
		Str("fu", fu).
		Str("date", date).
		Str("spinestring", spine).
		Str("filename", mesoFile).
		Str("mesostic", showR).
		Msg("NASA APOD Mesostic complete")
}
