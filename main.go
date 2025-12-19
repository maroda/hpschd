/*

	HPSCHD Main

*/

package main

import (
	"flag"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Zerolog
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Runtime Flags
	debug := flag.Bool("debug", false, "Log Level: DEBUG")
	nofetch := flag.Bool("nofetch", false, "Do not start NASA APOD cronjob")

	// Parse Flags
	flag.Parse()

	// Flag Options
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Info().Msg("Log level set to DEBUG")
	}

	/*
		Confirm / initiate data locations

		store ::: ephemeral mesostic cache
		txrx ::: tmp scratch files
	*/
	datadirs := []string{"store", "txrx"}
	localDirs(datadirs)

	// Fetching the NASA APOD for the homepage display is default behavior.
	// The 'nofetch' flag turns this off.
	if *nofetch {
		log.Info().Msg("Running with integrated NASA APOD fetch disabled.")
	} else {
		// Fetch initial APOD mesostic to populate store before starting web server
		// This prevents ENOENT errors when users visit homepage before cronjob runs
		apiKey := envVar("NASA_API_KEY", "DEMO_KEY")
		apodURL := "https://api.nasa.gov/planetary/apod?api_key=" + apiKey

		log.Info().Msg("Fetching initial NASA APOD mesostic...")
		NASAetl(apodURL)
		log.Info().Msg("Initial mesostic created, starting cronjob.")

		timer := envVar("HPSCHD_TIMER", "77")
		timerI, err := strconv.Atoi(timer)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to parse HPSCHD_TIMER")
		}

		// Start up ticker for fetching source text to display on the homepage as a Mesostic.
		// The NASA APOD API has a query limit of 1k/hr, every 15s is 240/hr.
		// TODO: There is a possible retry bug here...
		//		increasingly EXISTENT (existing) mesostics trip up a fast fetch (e.g. 15s)
		//		try ~11m for something that keeps things fresh enough
		//		but will hopefully avoid the EXISTENT pileup
		go fetchTicker(uint64(timerI))
	}

	// Prometheus
	prometheus.MustRegister(hpschdPingCount)
	prometheus.MustRegister(hpschdHomeTimer)
	prometheus.MustRegister(hpschdJsubTimer)
	prometheus.MustRegister(hpschdFsubTimer)
	prometheus.MustRegister(hpschdMesolineTimer)
	prometheus.MustRegister(hpschdNASAetlTimer)

	// Deploy the web server
	rt := mux.NewRouter()

	// Basic Pages
	rt.Handle("/metrics", promhttp.Handler())
	rt.HandleFunc("/", homepage)
	rt.HandleFunc("/ping", ping)

	// API Features
	api := rt.PathPrefix("/app").Subrouter()
	api.HandleFunc("", JSubmit).Methods(http.MethodPost)       // JSON submission POST
	api.HandleFunc("/{arg}", FSubmit).Methods(http.MethodPost) // Form submission POST

	if err := http.ListenAndServe(":9999", rt); err != nil {
		log.Fatal().Err(err).Msg("startup failed!")
	}
}
