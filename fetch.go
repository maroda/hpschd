/*

	Mesostic Fetch Source

*/

package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
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
}

// fetchSource ::: Accepts a URL and returns three elements from the struct used to unmarshal the JSON response.
func fetchSource(u string) (string, string, string) {
	url := u

	// new HTTP client
	apodClient := http.Client{
		Timeout: time.Second * 2,
	}

	// new request object
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	// make the request
	req.Header.Set("User-Agent", "Go Mesostic hpschd.xyz")
	result, resErr := apodClient.Do(req)
	if resErr != nil {
		log.Fatal(resErr)
	}

	if result.Body != nil {
		defer result.Body.Close()
	}

	body, readErr := ioutil.ReadAll(result.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	ae := apodE{}

	jsonErr := json.Unmarshal(body, &ae)
	if jsonErr != nil {
		// log.Fatal(jsonErr)
		log.Fatalf("unable to parse value: %q, error: %s", string(body), jsonErr.Error())
	}

	return ae.Date, ae.Title, ae.Explain
}
