/*

	Mesostic Fetch Automator

	A static URL is configured for fetching a source text via API to be transmogrified into a Mesostic for display.

*/

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
)

// generic crontab emulator
func fetchCron() {
	fcron := gocron.NewScheduler(time.UTC)

	_, err := fcron.Every(30).Seconds().StartImmediately().Do(NASAetl)
	if err != nil {
		log.Error()
	}

	fcron.StartBlocking()
}

// NASAetl ::: Retrieve the Astronomy Picture of the Day (APOD) metadata and process it through the Mesostic engine.
// This is where the URL is configured.
//
func NASAetl() {
	fmt.Println("NASAetl Running")

	// NASA official Astronomy Picture of the Day endpoint URL using a freely available API key
	// url := "https://api.nasa.gov/planetary/apod?api_key=Ijb0zLeEt71HMQdy8YjqB583FK3bdh1yThVJYzpu"
	url := "https://api.nasa.gov/planetary/apod?date=2000-01-01&api_key=Ijb0zLeEt71HMQdy8YjqB583FK3bdh1yThVJYzpu"

	// the title as the spine, for now :)
	date, spine, source := fetchSource(url)

	// The most up-to-date image can throw a 404 towards the end of the day.
	// This is meant to provide a fall-back that displays a random one in the past.
	// Which might expand this conditional, or become a switch statement.
	if spine == "404" {
		log.Error().Str("code", "404").Msg("Remote data not available, deploying fallback [coming soon]")
	}

	// remove spaces from spine (camel case it)
	trcc := strings.NewReplacer(" ", "")
	spn := trcc.Replace(spine)

	// convert each phrase into a line by replacing commas and periods with newlines.
	trnl := strings.NewReplacer(". ", "\n", ", ", "\n")
	source = trnl.Replace(source)

	// the remainder of this function mimics the JSON API calls
	fileName := fileTmp(&spn, &source)
	mcMeso := make(chan string)
	go mesoMain(fileName, spn, mcMeso)
	showR := <-mcMeso

	mesoFile := apodNew(&spine, &date, &showR)
	fmt.Printf("%s populated with the following Mesostic: \n%s", mesoFile, showR)

	// remove the temp source file
	var ferr = os.Remove(fileName)
	if ferr != nil {
		log.Error()
	}

	log.Info().Str("date", date).Str("spine", spine).Msg("")
}

func apodNew(sp *string, da *string, me *string) string {
	tr := strings.NewReplacer(" ", "_")
	spn := tr.Replace(*sp)
	fP := fmt.Sprintf("public/%s__%s", *da, spn)

	// check here if fP exists and how old it is.
	// if it does and is under X days old, do nothing

	sB := []byte(*me)
	err := ioutil.WriteFile(fP, sB, 0644)
	if err != nil {
		log.Error()
	}
	return fP
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
