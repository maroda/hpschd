/*

	Mesostic Scheduler and Tasks

*/

package main

import (
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
)

// Channel for NASAapod Mesostic publishing
// Buffered with capacity 1 to prevent blocking on initial startup fetch
var nasaNewMESO = make(chan string, 1)

// fetchTicker takes fetch frequency in seconds (ffs) and runs the ETL job
func fetchTicker(ffs uint64) {
	// NASA official Astronomy Picture of the Day endpoint URL using NASA's demo API key
	apiKey := envVar("NASA_API_KEY", "DEMO_KEY")
	apodnow := "https://api.nasa.gov/planetary/apod?api_key=" + apiKey
	apodenv := "HPSCHD_NASA_APOD_URL" // Optional ENV VAR for full URL override
	url := envVar(apodenv, apodnow)   // NASA APOD URL to query, default if no ENV VAR

	// The ticker runs at the rate of the fetch frequency
	ticker := time.NewTicker(time.Second * time.Duration(ffs))
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			NASAetl(url)
		}
	}
}

// NASAetl ::: Retrieve Astronomy Picture of the Day (APOD) metadata,
// process it through the Mesostic engine, save it in a library of ephemeral copies,
// pass the new data point (filename path) to a channel for use with displays.
func NASAetl(url string) {
	hTimer := prometheus.NewTimer(hpschdNASAetlTimer)
	defer hTimer.ObserveDuration()

	_, _, fu := Envelope()

	log.Info().
		Str("fu", fu).
		Msg("NASA APOD Mesostic Begin")

	// the title as the spine, for now :)
	date, spine, source := fetchSource(url)

	// There is typically a long stretch of time from ~0000UTC to
	// sometime the next morning while the APOD for the next day is being updated.
	// NASA APOD API will return: 'no data available for date: YYYY-MM-DD'
	if spine == "404" {
		log.Warn().Str("fu", fu).Str("code", "404").
			Msg("Remote data not available, waiting until next timed request.")
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
	nasaNewMESO <- mesoFile

	log.Info().
		Str("spinestring", spine).
		Str("filename", mesoFile).
		Msg("NASA APOD Mesostic End")

	log.Debug().
		Str("fu", fu).
		Str("fetchdate", date).
		Str("spinestring", spine).
		Str("filename", mesoFile).
		Str("mesostic", showR).
		Msg("NASA APOD Mesostic End")
}
