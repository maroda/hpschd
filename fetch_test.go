/*

	Tests

*/

package main

// TestTFetchSource ::: Match a fetched static URL (i.e. not the latest APOD) with known values.
/* this has a bug, if DEMO_KEY hits a rate limit (code 429) this test fails
func TestTFetchSource(t *testing.T) {
	fmt.Printf("\n\t::: Test Target fetchSource() :::\n")

	// Remove this if normal logging from the called fuction is needed.
	zerolog.SetGlobalLevel(zerolog.Disabled)

	url := "https://api.nasa.gov/planetary/apod?date=2000-01-01&api_key=DEMO_KEY"
	matchDate := "2000-01-01"
	matchTitle := "The Millennium that Defines Universe"

	date, title, _ := fetchSource(url)

	if date != matchDate {
		t.Errorf("%s does not match %s\n", date, matchDate)
	}

	if title != matchTitle {
		t.Errorf("%s does not match %s\n", title, matchTitle)
	}
}
*/
