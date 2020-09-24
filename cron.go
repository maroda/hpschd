/*

	Mesostic Fetch Automator

	A static URL is configured for fetching a source text via API to be transmogrified into a Mesostic for display.

*/

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
)

// generic crontab emulator
func fetchCron() {
	fcron := gocron.NewScheduler(time.UTC)

	_, err := fcron.Every(10).Seconds().StartImmediately().Do(NASAetl)
	if err != nil {
		log.Error()
	}

	fcron.StartBlocking()
}

// NASAetl ::: Retrieve the Astronomy Picture of the Day (APOD) metadata and process it through the Mesostic engine.
// This is where the URL is configured.
//
func NASAetl() {
	url := "https://api.nasa.gov/planetary/apod?api_key=Ijb0zLeEt71HMQdy8YjqB583FK3bdh1yThVJYzpu"

	fmt.Println("NASAetl Running")

	// the title as the spine should be interesting.
	date, spine, source := fetchSource(url)
	if spine == "404" {
		// TODO: The most up-to-date image can throw a 404 towards the end of the day. Provide a fall-back that displays yesterday's.
		log.Error()
		return
	}

	// fmt.Println(date, spine, source)

	// the remainder of this function mimics the JSON API calls

	fileName := fileTmp(&spine, &source)

	mcMeso := make(chan string)
	go mesoMain(fileName, spine, mcMeso)
	showR := <-mcMeso
	fmt.Println(showR)
	// fmt.Fprintf(w, "%s\n", showR)

	var ferr = os.Remove(fileName)
	if ferr != nil {
		log.Error()
	}

	log.Info().Str("date", date).Str("spine", spine).Msg("")
}

// fileTmp ::: Take a source string and place it in a file name after the spinestring.
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
