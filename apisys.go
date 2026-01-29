package main

import (
	"encoding/json"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
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

	r.HandleFunc("/", sp.HomeHandler)
	r.HandleFunc("/healthz", sp.HealthzHandler)

	api := r.PathPrefix("/app").Subrouter()
	api.HandleFunc("", sp.GetJSON).Methods(http.MethodPost)

	return r
}

func (sp *ServePoems) HealthzHandler(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok")) }

// HomeMesostic is the final poem that is served on the homepage or returned by the API.
type HomeMesostic struct {
	mu    sync.Mutex
	Title string `json:"title"`
	ADate string `json:"date"`
	Poem  string `json:"poem"`
}

// HomeHandler displays the new mesostic on the homepage
func (sp *ServePoems) HomeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := otel.Tracer("mesostic/web").Start(ctx, "Homepage")
	defer span.End()

	hometmpl := template.Must(template.ParseFiles("public/index.html"))

	// This is created every time the homepage is requested
	hm := HomeMesostic{}
	cache := "store"             // Datastore of created poems
	rndFile := ichingMeso(cache) // Random filename from existing poems

	// Create human-friendly date and title
	nameParts := strings.Split(rndFile, "_")
	fullDate, ok := strings.CutPrefix(nameParts[0], cache+"/")
	if !ok {
		slog.Warn("Could not parse date")
	}
	fullTitle := nameParts[1:]

	// Update the struct, locking during file access
	hm.mu.Lock()
	hm.Title = strings.TrimSpace(strings.Join(fullTitle, " "))
	hm.ADate = fullDate
	hm.Poem = readMesoFile(ctx, &rndFile) // Load poem from mesostic file
	hm.mu.Unlock()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	// Write the HTML via template
	err := hometmpl.Execute(w, hm)
	if err != nil {
		span.RecordError(err)
		slog.Error("cannot render html")
		http.Error(w, "cannot render html", http.StatusInternalServerError)
	}
}

// GetJSON returns a plaintext mesostic from the submitted JSON
func (sp *ServePoems) GetJSON(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := otel.Tracer("mesostic/api").Start(ctx, "GetJSON")
	defer span.End()

	di := &DataAPI{}

	// Rate limit first, then read the body for processing
	maxbytes := int64(1024 * 1024)
	r.Body = http.MaxBytesReader(w, r.Body, maxbytes)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		span.RecordError(err)
		slog.Error("body unreadable or exceeded limit")
		http.Error(w, "body unreadable or exceeded limit", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, di)
	if err != nil {
		span.RecordError(err)
		slog.Error("cannot unmarshal body")
		http.Error(w, "cannot unmarshal body", http.StatusInternalServerError)
		return
	}

	// Validate: If either these two fields are not filled, it's a bad request
	if di.SpineString == "" || di.Text == "" {
		span.AddEvent("empty text or spinestring", trace.WithAttributes(attribute.String("spinestring", di.SpineString)))
		slog.Error("empty text or spinestring")
		http.Error(w, "empty text or spinestring", http.StatusBadRequest)
		return
	}

	title := di.SpineString
	os.Unsetenv("HPSCHD_SPINESTRING") // The API overrides this setting
	m := NewMesostic(ctx, title, string(body), di)

	m.MU.Lock()
	m.SourceTxt = di.Text
	m.Date = time.Now().Format("2006-01-02")
	m.MU.Unlock()

	// HomeMesostic is used, identical to how it's served on the homepage
	mesojson := &HomeMesostic{
		Title: title,
		Poem:  m.BuildMeso(ctx),
	}

	// Write the mesostic back
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(mesojson)
	if err != nil {
		span.RecordError(err)
		slog.Error("Encode Error", slog.Any("error", err))
		http.Error(w, "encode error", http.StatusInternalServerError)
		return
	}
}
