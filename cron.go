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

	// There is typically a long stretch of time from ~0000UTC to
	// sometime the next morning while the APOD for the next day is being updated.
	// NASA APOD API will return: 'no data available for date: YYYY-MM-DD'
	//
	// Trigger a new fetch for a new mesostic added to the store and quit.
	if spine == "404" {
		go NASAetl(fetchRandURL())

		log.Warn().
			Str("fu", fu).
			Str("code", "404").
			Msg("Remote data not available, randomized ETL triggered.")

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
	mesoFile, created := apodNew(&spine, &date, &showR)

	// If the mesostic file already exists, no more action is needed.
	// Trigger a new fetch for a new mesostic added to the store and quit.
	// TODO: This check should go *before* creating the mesostic at all.
	// 			e.g. construct the filename and check against dirents()
	if !created {
		go NASAetl(fetchRandURL())

		log.Warn().
			Str("fu", fu).
			Str("code", "204").
			Msg("Local mesostic exists, randomized ETL triggered.")

		return
	}

	// remove the tmp source file
	var ferr = os.Remove(tmpFileName)
	if ferr != nil {
		log.Error()
	}

	// push filename of new Mesostic
	// for now this is a non-buffered blocking operation
	nasaNewMESO <- mesoFile

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
