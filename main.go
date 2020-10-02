/*

	HPSCHD Main

*/

package main

import (
	"flag"
	"net/http"

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

	// Parse Flags
	flag.Parse()

	// Flag Options
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// Prometheus
	prometheus.MustRegister(msgPostCnt)
	prometheus.MustRegister(msgPostDur)
	prometheus.MustRegister(pingCnt)

	// Start up scheduler for fetching source text to display on the homepage as a Mesostic
	go fetchCron()

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
