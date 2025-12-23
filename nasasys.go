package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	webTimeout = 10 * time.Second
)

type DataAPOD struct {
	Copyright      string `json:"copyright"`
	Date           string `json:"date"`
	Explaination   string `json:"explanation"`
	HDURL          string `json:"media"`
	MediaType      string `json:"media_type"`
	ServiceVersion string `json:"service_version"`
	Title          string `json:"title"`
	URL            string `json:"url"`
	Code           int    `json:"code"`
	Msg            string `json:"msg"`
}

// HTTPClient is implemented to enable testing and
// supports the global client with SingleFetchWithClient
type HTTPClient interface {
	Get(string) (*http.Response, error)
}

// sharedHTTPClient is a global client to ensure we're reusing connections
// (putting this directly in SingleFetch would
// create a new client every request every second)
var sharedHTTPClient = &http.Client{
	Timeout: webTimeout,
	Transport: &http.Transport{
		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 2,
		IdleConnTimeout:     30 * time.Second,
	},
}

// TickerAPOD takes a frequency in seconds (freq) and runs GetAPOD
// This is a method of the ServePoems struct that does the data handling
func (sp *ServePoems) TickerAPOD() {
	// NASA official Astronomy Picture of the Day endpoint URL using NASA's demo API key
	apiKey := envVar("NASA_API_KEY", "DEMO_KEY")
	apodnow := "https://api.nasa.gov/planetary/apod?api_key=" + apiKey

	apodenv := "HPSCHD_NASA_APOD_URL" // Optional ENV VAR for full URL override
	url := envVar(apodenv, apodnow)   // NASA APOD URL to query, default if no ENV VAR

	fs := RealFS{}

	// The first time this runs, it fetches the current default date
	_, err := GetAPOD(url, fs)
	if err != nil {
		log.Error().Err(err).Msg("could not get APOD")
		slog.Error("Error fetching APOD from " + url)
	}

	for {
		select {
		case <-sp.Ticker.C:
			// Randomized dates are used for all subsequent fetches.
			date := rndDate(time.Now().UnixNano())
			url = "https://api.nasa.gov/planetary/apod?date=" + date + "&api_key=" + apiKey
			_, err = GetAPOD(url, fs)
			if err != nil {
				log.Error().Err(err).Msg("could not get APOD")
				slog.Error("Error fetching APOD from " + url)
			}
		}
	}
}

// SingleFetchWithClient handles the messy business of the HTTP connection
// and is testable with dependency injection, called by SingleFetch
func SingleFetchWithClient(url string, c HTTPClient) (int, []byte, error) {
	resp, err := c.Get(url)
	if err != nil {
		slog.Error("Fetch Error", slog.Any("Error", err))
		return 0, nil, err
	}

	// This io.ReadAll block does not have test coverage
	// Accepting this because of how difficult it is to mock
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Could not read body", slog.Any("Error", err))
		return 0, nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Error("Close Error", slog.Any("Error", err))
			return
		}
	}()

	return resp.StatusCode, body, err
}

// SingleFetch returns the Response Code, raw byte stream body, and error
// This uses a Shared HTTP Client:
// - to reuse existing endpoint connections
// - to avoid stale connections that eat up OS FDs
func SingleFetch(url string) (int, []byte, error) {
	return SingleFetchWithClient(url, sharedHTTPClient)
}

// FileSystem is for operating with local configs and/or data.
type FileSystem interface {
	Open(name string) (*os.File, error)
	Stat(name string) (os.FileInfo, error)
	WriteFile(name string, data []byte, perm os.FileMode) error
}

type RealFS struct{}

func (fs RealFS) WriteFile(name string, data []byte, perm os.FileMode) error {
	return os.WriteFile(name, data, perm)
}

func (fs RealFS) Open(name string) (*os.File, error) {
	return os.Open(name)
}

func (fs RealFS) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

// GetAPOD writes the resulting APOD output to the store
// There is no trigger here yet for fetching a random APOD
// which was formerly done by NASAetl when no file gets created
// i.e. it's already been created. (see: /if !created/ and the recursive func )
// That type of functionality should be added in here, so we
// get the continuously updating random URL dates.
func GetAPOD(url string, fs FileSystem) (string, error) {
	_, body, err := SingleFetch(url)
	if err != nil {
		slog.Error("Fetch Error", slog.Any("Error", err))
		return "", fmt.Errorf("fetch url error: %s", url)
	}

	// Each "Get____" function will have its own defined structure
	dd := &DataAPOD{}
	err = json.Unmarshal(body, dd)
	if err != nil {
		slog.Error("Unmarshal Error", slog.Any("Error", err))
		return "", fmt.Errorf("unmarshal url error: %s", url)
	}
	title := dd.Title
	m := NewMesostic(title, string(body), dd)

	// When it needs to write directly to the struct, a lock is required
	m.MU.Lock()
	m.SourceTxt = dd.Explaination
	m.Date = dd.Date
	m.MU.Unlock()

	// When using methods from the struct, no lock is used
	mesostic := m.BuildMeso()
	filename := fmt.Sprintf("store/%s__%s", m.Date, strings.Join(m.Spine, ""))
	if err = fs.WriteFile(filename, []byte(mesostic), 0644); err != nil {
		log.Error().Err(err).Msgf("Failed to write mesostic")
		slog.Error("write file error")
		return "", fmt.Errorf("write file error: %s", filename)
	}

	log.Info().Str("filename", filename).Msg("apod mesostic stored")
	return m.Poem, nil
}
