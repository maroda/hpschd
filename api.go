/*

	Mesostic API Front

	/app - API endpoint
	/ping - a readiness check
	/metrics - prometheus metrics

*/

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

// define prometheus metrics
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

// Submit ::: Data object model for JSON submissions
type Submit struct {
	Text        string
	SpineString string
}

// HTML template
var templates = template.Must(template.ParseFiles("public/upload.html"))

// Render HTML for FPage + FUpload functionality.
func display(w http.ResponseWriter, page string, data interface{}) {
	templates.ExecuteTemplate(w, page+".html", data)
}

// FPage ::: GET Method file upload.
func FPage(w http.ResponseWriter, r *http.Request) {
	//msgPostCnt.Add(1)
	//msgTimer := prometheus.NewTimer(msgPostDur)
	//defer msgTimer.ObserveDuration()

	w.WriteHeader(http.StatusOK)

	display(w, "upload", nil)

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

// FUpload ::: POST Method file upload.
func FUpload(w http.ResponseWriter, r *http.Request) {
	msgPostCnt.Add(1)
	msgTimer := prometheus.NewTimer(msgPostDur)
	defer msgTimer.ObserveDuration()

	w.WriteHeader(http.StatusOK)

	// 5MB memory limit
	r.ParseMultipartForm(5 << 20)

	// Form File handler
	ffile, handle, err := r.FormFile("source")
	if err != nil {
		log.Error()
		return
	}

	defer ffile.Close()
	fmt.Printf("Uploaded: %+v\n", handle.Filename)
	fmt.Printf("Size: %+v\n", handle.Size)
	fmt.Printf("MIME Type: %+v\n", handle.Header)

	// for now we'll just test this and create a local copy of the file,
	// which could then be passed to mesostic as an alternate method if the
	// subroutine is unavailable. both should be tested.

	// create the local file
	f, err := os.Create(handle.Filename)
	defer f.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// copy the form upload_file to the disk_file
	if _, err := io.Copy(f, ffile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// pass the name of the disk_file to the function
	// mesoMain(handle.Filename)
	mcMeso := make(chan string)
	go mesoMain(handle.Filename, mcMeso)
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
		Msg("New Upload")
}

// FSubmit ::: POST Method form submission.
func FSubmit(w http.ResponseWriter, r *http.Request) {
	msgPostCnt.Add(1)
	msgTimer := prometheus.NewTimer(msgPostDur)
	defer msgTimer.ObserveDuration()

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
	msgPostCnt.Add(1)
	msgTimer := prometheus.NewTimer(msgPostDur)
	defer msgTimer.ObserveDuration()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var subd Submit

	// decode body into struct and pull the value to be encoded
	if err := json.NewDecoder(r.Body).Decode(&subd); err != nil {
		log.Fatal().Err(err).Msg("failed to decode body")
	}
	source := subd.Text
	spine := subd.SpineString

	fmt.Fprintf(w, "Source: %s ::: Spine: %s\n", source, spine)

	// this needs to be modified so that the JSON version
	// can use the same function as the upload / form version
	// mcMeso := make(chan string)
	// go Mesostic(source, spine, mcMeso)
	// showR := <-mcMeso
	// fmt.Println(showR)

	log.Info().
		Str("host", r.Host).
		Str("ref", r.RemoteAddr).
		Str("xref", r.Header.Get("X-Forwarded-For")).
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Str("proto", r.Proto).
		Str("agent", r.Header.Get("User-Agent")).
		Str("response", "200").
		Msg("New JSON")
}

// readiness checks are counted but not logged
func ping(w http.ResponseWriter, r *http.Request) {
	pingCnt.Add(1)
	w.Write([]byte("pong\n"))
}

// HTTP frontend for Mesostic API
// This controls flow.
func main() {
	// Prometheus
	prometheus.MustRegister(msgPostCnt)
	prometheus.MustRegister(msgPostDur)
	prometheus.MustRegister(pingCnt)

	rt := mux.NewRouter()
	rt.HandleFunc("/ping", ping)
	rt.Handle("/metrics", promhttp.Handler())

	// currently this upload method uses the hardcoded SpineString
	// TODO: Add SpineString entry in the form
	upload := rt.PathPrefix("/upload").Subrouter()
	upload.HandleFunc("", FPage).Methods(http.MethodGet)    // Upload page GET
	upload.HandleFunc("", FUpload).Methods(http.MethodPost) // File upload POST

	api := rt.PathPrefix("/app").Subrouter()
	api.HandleFunc("", JSubmit).Methods(http.MethodPost)       // JSON submission POST
	api.HandleFunc("/{arg}", FSubmit).Methods(http.MethodPost) // Form submission POST

	if err := http.ListenAndServe(":9999", rt); err != nil {
		log.Fatal().Err(err).Msg("startup failed!")
	}
}
