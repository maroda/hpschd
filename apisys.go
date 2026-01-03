package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type DataAPI struct {
	Text        string `json:"text"`
	SpineString string `json:"spinestring"`
}

// ServePoems is called by main() and contains the mux
type ServePoems struct {
	Server *http.Server
	Mux    *mux.Router
	Ticker *time.Ticker
}

// SetupMux provides a new Mux with its internal routing configured
// These are the control points for Toadlester
func (sp *ServePoems) SetupMux() *mux.Router {
	r := mux.NewRouter()

	r.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/", sp.HomeHandler)
	r.HandleFunc("/healthz", sp.HealthzHandler)

	api := r.PathPrefix("/app").Subrouter()
	api.HandleFunc("", sp.GetJSON).Methods(http.MethodPost)

	return r
}

func (sp *ServePoems) HealthzHandler(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok")) }

type HomeMesostic struct {
	MU    sync.Mutex
	Title string `json:"title"`
	Poem  string `json:"poem"`
}

// HomeHandler displays the new mesostic on the homepage
func (sp *ServePoems) HomeHandler(w http.ResponseWriter, r *http.Request) {
	hometmpl := template.Must(template.ParseFiles("public/poem.html"))

	// This is created every time the homepage is requested
	hm := HomeMesostic{}
	hm.MU.Lock()

	cache := "store"                 // Datastore of created poems
	rndFile := ichingMeso(cache)     // Random filename from existing poems
	hm.Title = rndFile               // Read title of mesostic file
	hm.Poem = readMesoFile(&rndFile) // Load poem from mesostic file
	hm.MU.Unlock()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	err := hometmpl.Execute(w, hm)
	if err != nil {
		slog.Error("cannot render html")
		http.Error(w, "cannot render html", http.StatusInternalServerError)
	}
}

// GetJSON is the v1 API that will be obviated by JSONHandler
func (sp *ServePoems) GetJSON(w http.ResponseWriter, r *http.Request) {
	di := &DataAPI{}

	// Rate limit first, then read the body for processing
	maxbytes := int64(1024 * 1024)
	r.Body = http.MaxBytesReader(w, r.Body, maxbytes)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("body unreadable or exceeded limit")
		http.Error(w, "body unreadable or exceeded limit", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	fmt.Println(string(body))
	err = json.Unmarshal(body, di)
	if err != nil {
		slog.Error("cannot unmarshal body")
		http.Error(w, "cannot unmarshal body", http.StatusInternalServerError)
		return
	}

	// Validate: If either these two fields are not filled, it's a bad request
	if di.SpineString == "" || di.Text == "" {
		slog.Error("empty text or spinestring")
		http.Error(w, "empty text or spinestring", http.StatusBadRequest)
		return
	}
	title := di.SpineString
	os.Unsetenv("HPSCHD_SPINESTRING") // The API overrides this setting
	m := NewMesostic(title, string(body), di)

	m.MU.Lock()
	m.SourceTxt = di.Text
	m.Date = time.Now().Format("2006-01-02")
	m.MU.Unlock()

	mesostic := m.BuildMeso()

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(mesostic)
	if err != nil {
		slog.Error("Encode Error", slog.Any("error", err))
		http.Error(w, "encode error", http.StatusInternalServerError)
		return
	}
}

func (sp *ServePoems) JSONHandler(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}
