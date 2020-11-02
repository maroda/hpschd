/*

	Mesostic RESTish API

	/app - API endpoint
	/ping - Readiness check
	/metrics - Prometheus metrics
	/homepage - Frontend displays the NASA APOD Mesostic

*/

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
)

// Submit ::: Data object model for JSON submissions
type Submit struct {
	Text        string
	SpineString string
}

// homepage ::: Home
/*
The idea with the homepage is that the Mesostic has already been built, and loading home will show one.

There could be chance operations to pick which date after 2000-01-01 to use.

But as the various dates are chosen over time, the cache of mesostics will increase, which means the longer an instance stays running, the more mesostic variation it gets to display.

In other words, every iteration pull a new APOD and create a mesostic for the library.

When index is loaded, pull a mesostic at random and display it.

Decoupling the cron fetching the text from the display.
*/

// MesoPrint ::: Elements for HTML rendering
type MesoPrint struct {
	Title    string // Page Title
	Mesostic string // The New Mesostic
}

func homepage(w http.ResponseWriter, r *http.Request) {
	hTimer := prometheus.NewTimer(hpschdHomeTimer)
	defer hTimer.ObserveDuration()
	_, _, fu := Envelope()

	w.WriteHeader(http.StatusOK)

	// struct for importing into the HTML template
	var formatMeso MesoPrint

	// this function reads the first item off the top of the channel
	// this channel is populated with the filename of the newest created mesostic
	// which is the result of a a 15m cronjob to fetch the NASA APOD
	// currently this channel is non-buffered but meant to be left open
	//
	// when the channel is empty, this function returns a special signal HPSCHD
	// (does the channel stay empty if the read function isn't run? does it block the cronjob? will the cronjob pile up jobs if this happens?)
	// this instructs the loading homepage to randomly select a previously created mesostic
	var mesoFile string = nasaNewREAD()

	switch mesoFile {
	case "HPSCHD":
		// The channel reader has returned the signal for "no more data".
		// TODO: currently this isn't random, it needs to be
		mesoDir := "store"               // An ephemeral 'datastore' of previously created mesostics.
		iMesoFile := ichingMeso(mesoDir) // The i-ching-like engine for choosing a random mesostic.
		formatMeso.Title = iMesoFile
		formatMeso.Mesostic = readMesoFile(&iMesoFile)

		log.Info().
			Str("fu", fu).
			Str("filename", mesoFile).
			Msg("Chance Operations Indicated")
	default:
		// A filename exists on the channel and has been returned.
		formatMeso.Title = mesoFile
		formatMeso.Mesostic = readMesoFile(&mesoFile)

		log.Info().
			Str("fu", fu).
			Str("filename", mesoFile).
			Msg("Mesostic formatted")
	}

	// display the new mesostic on the homepage
	hometmpl := template.Must(template.ParseFiles("public/index.html"))
	err := hometmpl.Execute(w, formatMeso)
	if err != nil {
		log.Fatal().Str("fu", fu).Msg("Cannot render HTML")
	}

	log.Info().
		Str("host", r.Host).
		Str("ref", r.RemoteAddr).
		Str("xref", r.Header.Get("X-Forwarded-For")).
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Str("proto", r.Proto).
		Str("agent", r.Header.Get("User-Agent")).
		Str("response", "200").
		Msg("")
}

// FSubmit ::: POST Method form submission.
func FSubmit(w http.ResponseWriter, r *http.Request) {
	hTimer := prometheus.NewTimer(hpschdFsubTimer)
	defer hTimer.ObserveDuration()

	w.WriteHeader(http.StatusOK)

	// Take the given path as the Spine String.
	args := mux.Vars(r)
	spine := args["arg"]
	fmt.Printf("spine = %s\n", spine)

	r.ParseForm()
	for k, v := range r.Form {
		fmt.Printf("%s = %s\n", k, v)
	}

	log.Info().
		Str("host", r.Host).
		Str("ref", r.RemoteAddr).
		Str("xref", r.Header.Get("X-Forwarded-For")).
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Str("proto", r.Proto).
		Str("agent", r.Header.Get("User-Agent")).
		Str("response", "200").
		Msg("New Form Submission")
}

// JSubmit ::: POST Method JSON submission.
func JSubmit(w http.ResponseWriter, r *http.Request) {
	hTimer := prometheus.NewTimer(hpschdJsubTimer)
	defer hTimer.ObserveDuration()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var subd Submit

	// decode body into struct
	if err := json.NewDecoder(r.Body).Decode(&subd); err != nil {
		log.Fatal().Err(err).Msg("failed to decode body")
	}
	source := subd.Text       // the multi-line source for the Mesostic
	spine := subd.SpineString // the SpineString for the Mesostic

	fileName := fileTmp(&spine, &source)
	mcMeso := make(chan string)
	go mesoMain(fileName, spine, mcMeso)

	// receive the channel data and display result
	showR := <-mcMeso
	fmt.Println(showR)
	fmt.Fprintf(w, "%s\n", showR)

	log.Info().
		Str("host", r.Host).
		Str("ref", r.RemoteAddr).
		Str("xref", r.Header.Get("X-Forwarded-For")).
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Str("proto", r.Proto).
		Str("agent", r.Header.Get("User-Agent")).
		Str("response", "200").
		Str("tmp", fileName).
		Msg("New JSON")
}

// readiness checks are Counted but not logged
func ping(w http.ResponseWriter, r *http.Request) {
	hpschdPingCount.Add(1)
	w.Write([]byte("pong\n"))
}
