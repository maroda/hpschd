/*

	Mesostic Observability

	Global definitions and functions related to metrics and events.

*/

package main

import (
	"runtime"

	"github.com/prometheus/client_golang/prometheus"
)

/*
	Global Prometheus Definitions
	LinearBuckets are defined as: 50 Buckets, 10ms each, starting at 1ms
*/

// Access Counts
var hpschdPingCount = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "hpschdPingCount",
	Help: "Total number of Readiness pings.",
})

// Function Timers
var hpschdHomeTimer = prometheus.NewHistogram(prometheus.HistogramOpts{
	Name:    "hpschdHomeTimer",
	Help:    "Historgram for the runtime of homepage.",
	Buckets: prometheus.LinearBuckets(0.001, 0.01, 50),
})

var hpschdJsubTimer = prometheus.NewHistogram(prometheus.HistogramOpts{
	Name:    "hpschdJsubTimer",
	Help:    "Historgram for the runtime of jsubmit (JSON).",
	Buckets: prometheus.LinearBuckets(0.001, 0.01, 50),
})

var hpschdFsubTimer = prometheus.NewHistogram(prometheus.HistogramOpts{
	Name:    "hpschdFsubTimer",
	Help:    "Historgram for the runtime of fsubmit (Form).",
	Buckets: prometheus.LinearBuckets(0.001, 0.01, 50),
})

var hpschdMesolineTimer = prometheus.NewHistogram(prometheus.HistogramOpts{
	Name:    "hpschdMesolineTimer",
	Help:    "Historgram for the runtime of mesoLine.",
	Buckets: prometheus.LinearBuckets(0.001, 0.01, 50),
})

var hpschdNASAetlTimer = prometheus.NewHistogram(prometheus.HistogramOpts{
	Name:    "hpschdNASAetlTimer",
	Help:    "Historgram for the runtime of NASAetl.",
	Buckets: prometheus.LinearBuckets(0.001, 0.01, 50),
})

// Envelope ::: Returns details about the current code execution point.
// This enables tracing in log events, for instance from within a function:
//		_, _, fu := Envelope()
//		fmt.Printf("current function: %s", fu)
func Envelope() (string, int, string) {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return frame.File, frame.Line, frame.Function
}
