/*

	Mesostic Fetch Source

*/

package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type apodE struct {
	Copyright string `json:"copyright"`
	Date      string `json:"date"`
	Explain   string `json:"explanation"`
	HDURL     string `json:"hdurl"`
	Media     string `json:"media_type"`
	Version   string `json:"service_version"`
	Title     string `json:"title"`
	URL       string `json:"url"`
	Code      int    `json:"code"` // this is typically 404
	Msg       string `json:"msg"`  // typically the 404 reason
}

// fetchSource ::: Accepts a URL and returns three elements from the struct used to unmarshal the JSON response.
func fetchSource(u string) (string, string, string) {
	url := u
	log.Info().Str("url", url).Msg("URL Received")

	// new HTTP client
	apodClient := http.Client{
		Timeout: time.Second * 2,
	}

	// new request object
	req, reqErr := http.NewRequest(http.MethodGet, url, nil)
	if reqErr != nil {
		log.Error().Err(reqErr).Msg("")
	}

	// make the request
	req.Header.Set("User-Agent", "Go Mesostic hpschd.xyz")
	result, resErr := apodClient.Do(req)
	if resErr != nil {
		log.Error().Err(resErr).Msg("")
	}
	// the timeout above caused the following error:
	// i was testing the 404 condition, so it's possible that NASA is updating the API?
	//	{"level":"fatal","error":"Get \"https://api.nasa.gov/planetary/apod?api_key=Ijb0zLeEt71HMQdy8YjqB583FK3bdh1yThVJYzpu\": context deadline exceeded (Client.Timeout exceeded while awaiting headers)","time":"2020-09-23T18:05:38-07:00"}

	if result.Body != nil {
		defer result.Body.Close()
	}

	body, readErr := ioutil.ReadAll(result.Body)
	if readErr != nil {
		log.Error().Err(readErr).Msg("")
	}

	ae := apodE{}

	jsonErr := json.Unmarshal(body, &ae)
	if jsonErr != nil {
		log.Error().Err(jsonErr).Msg("unable to parse value")
	}

	if ae.Code == 404 {
		log.Warn().Int("code", ae.Code).Str("msg", ae.Msg).Msg("no data")
		return "err", "404", ae.Msg
	}

	log.Info().
		Str("date", ae.Date).
		Str("title", ae.Title).
		Str("source", ae.Explain).
		Msg("Source Extracted")
	return ae.Date, ae.Title, ae.Explain
}
