/*

	Mesostic Observability

	Global definitions and functions related to metrics and events.

*/

package main

import (
	"runtime"

	"github.com/prometheus/client_golang/prometheus"
)

// Global Prometheus Definitions
var msgPostCnt = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "mesostic_post_app_total",
	Help: "Total number of POST me-api requests.",
})

var pingCnt = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "mesostic_ping_total",
	Help: "Total number of Readiness pings.",
})

var msgPostDur = prometheus.NewHistogram(prometheus.HistogramOpts{
	Name: "mesostic_post_app_timer_seconds",
	Help: "Historgram for the runtime of POST to /app",
	// 50 Buckets, 10ms each, starting at 1ms
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
