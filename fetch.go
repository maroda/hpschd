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
	_, _, fu := Envelope()

	url := u

	log.Info().
		Str("fu", fu).
		Str("url", url).
		Msg("URL Received")

	// new HTTP client, request timeout of 10s
	apodClient := http.Client{
		Timeout: time.Second * 10,
	}

	// new request object
	req, reqErr := http.NewRequest(http.MethodGet, url, nil)
	if reqErr != nil {
		log.Error().
			Str("fu", fu).
			Err(reqErr).
			Msg("")
	}

	// make the request
	req.Header.Set("User-Agent", "Go Mesostic hpschd.xyz")
	result, resErr := apodClient.Do(req)
	if resErr != nil {
		log.Error().
			Str("fu", fu).
			Err(resErr).
			Msg("")
	}

	if result.Body != nil {
		defer result.Body.Close()
	}

	body, readErr := ioutil.ReadAll(result.Body)
	if readErr != nil {
		log.Error().
			Str("fu", fu).
			Err(readErr).
			Msg("")
	}

	ae := apodE{}

	jsonErr := json.Unmarshal(body, &ae)
	if jsonErr != nil {
		log.Error().
			Str("fu", fu).
			Err(jsonErr).
			Msg("unable to parse value")
	}

	if ae.Code == 404 {
		log.Warn().
			Str("fu", fu).
			Int("code", ae.Code).
			Str("msg", ae.Msg).
			Msg("no data")
		return "err", "404", ae.Msg
	}

	log.Info().
		Str("fu", fu).
		Str("date", ae.Date).
		Str("title", ae.Title).
		Msg("Source Extracted")

	log.Debug().
		Str("fu", fu).
		Str("date", ae.Date).
		Str("title", ae.Title).
		Str("source", ae.Explain).
		Msg("Source Extracted")

	return ae.Date, ae.Title, ae.Explain
}

// fetchRandURL ::: Returns a constructed string using a random date for the NASA APOD API query.
func fetchRandURL() string {
	salt := time.Now().Unix()
	date := rndDate(salt)
	url := "https://api.nasa.gov/planetary/apod?date=" + date + "&api_key=Ijb0zLeEt71HMQdy8YjqB583FK3bdh1yThVJYzpu"
	return url
}
