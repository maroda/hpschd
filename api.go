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
	"io/ioutil"
	"net/http"
	"os"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	_, _, fu := Envelope()

	w.WriteHeader(http.StatusOK)

	// struct for importing into the HTML template
	var formatMeso MesoPrint

	// read the cron channel for a new mesostic
	// this is ALMOST working
	// what's happening is it clears the channel of the new mesostic
	// and if there are more mesostics waiting in the channel, it will display them
	// but when the channel is empty, it will default here, and nothing is displayed
	// so maybe the fix here is to "cache" the current mesostic filename somehow
	// and if the channel is empty (no new ones), display a chance derived one from the library.

	// this function reads the first item off the top of the channel
	var mesoFile string = nasaNewREAD()

	switch mesoFile {
	case "HPSCHD":
		mesoDir := "store"
		iMesoFile := ichingMeso(mesoDir)
		formatMeso.Title = iMesoFile
		formatMeso.Mesostic = readMesoFile(&iMesoFile)

		log.Info().
			Str("fu", fu).
			Str("filename", mesoFile).
			Msg("Chance Operations Indicated")
	default:
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

// readMesoFile ::: Open and read the Mesostic
func readMesoFile(f *string) string {
	if len(*f) == 0 {
		log.Error().Msg("no path given")
		return "error"
	}

	var mesoBuf []byte
	mesoBuf, err := ioutil.ReadFile(*f)
	if err != nil {
		log.Error()
	}

	return string(mesoBuf)
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
	mcMeso := make(chan string)
	spine := "maroda" // figuring out where this should be set, for now it's manual
	go mesoMain(handle.Filename, spine, mcMeso)
	// go mesoMain(handle.Filename, mcMeso)

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

	// decode body into struct
	if err := json.NewDecoder(r.Body).Decode(&subd); err != nil {
		log.Fatal().Err(err).Msg("failed to decode body")
	}
	source := subd.Text       // the multi-line source for the Mesostic
	spine := subd.SpineString // the SpineString for the Mesostic

	// DEBUG ::: fmt.Fprintf(w, "Source:\n%s\nSpine:\n%s\n", source, spine)

	// dump the data to a tmp file
	// 	this mimics the multi-part upload version
	// 	placing data in a tmp file is extensible to
	// 	placing it in a database or other fast storage
	/*
		fT := time.Now()
		fS := fT.Unix()
		fN := fmt.Sprintf("%s__%d", spine, fS)
		sB := []byte(source)
		err := ioutil.WriteFile(fN, sB, 0644)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	*/

	fileName := fileTmp(&spine, &source)

	// mesoMain receives ::: tmp filename, the SpineString, data channel
	mcMeso := make(chan string)
	go mesoMain(fileName, spine, mcMeso)

	// receive the channel data and display result
	showR := <-mcMeso
	fmt.Println(showR)
	fmt.Fprintf(w, "%s\n", showR)

	// tmp file deletion should be non-blocking,
	// but we should know about it, and log it below.
	var ferr = os.Remove(fileName)
	if ferr != nil {
		log.Error()
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
		Str("tmp", fileName).
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

	// Start up scheduler for fetching source text to display on the homepage as a Mesostic
	// the right place for this might be a / page handler that issues the goroutine
	// fetchCron() may have to send data back on a channel?
	go fetchCron()

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
