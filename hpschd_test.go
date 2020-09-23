/*

	Tests

*/

package main

import (
	"fmt"
	"strings"
	"testing"
)

// TestFetchSource
// This URL: "https://api.nasa.gov/planetary/apod?date=2000-01-01&api_key=Ijb0zLeEt71HMQdy8YjqB583FK3bdh1yThVJYzpu"
// Returns this data: `{"date":"2000-01-01","explanation":"Welcome to the millennial year at the threshold of millennium three.  During millennium two, humanity continually redefined its concept of \"Universe\": first as spheres centered on the Earth, in mid-millennium as the Solar System, a few centuries ago as the Galaxy, and within the last century as the matter emanating from the Big Bang.  During millennium three humanity may hope to discover alien life, to understand the geometry and composition of our present concept of Universe, and even to travel through this Universe.  Whatever our accomplishments, humanity will surely find adventure and discovery in the space above and beyond, and possibly define the surrounding Universe in ways and colors we cannot yet imagine by the threshold of millennium four.","hdurl":"https://apod.nasa.gov/apod/image/0001/flammarion_halfcolor.gif","media_type":"image","service_version":"v1","title":"The Millennium that Defines Universe","url":"https://apod.nasa.gov/apod/image/0001/flammarion_halfcolor_big.gif"}`
// fetchSource() is called with this URL and matched with these known values.
func TestFetchSource(t *testing.T) {
	fmt.Printf("\n\t::: TestFetchSource :::\n")

	url := "https://api.nasa.gov/planetary/apod?date=2000-01-01&api_key=Ijb0zLeEt71HMQdy8YjqB583FK3bdh1yThVJYzpu"
	matchDate := "2000-01-01"
	matchTitle := "The Millennium that Defines Universe"

	// this test only checks date and title
	date, title, _ := fetchSource(url)

	if date != matchDate {
		t.Errorf("%s does not match %s\n", date, matchDate)
	}

	if title != matchTitle {
		t.Errorf("%s does not match %s\n", title, matchTitle)
	}

	fmt.Printf("Date '%s' matches '%s'\nTitle '%s' matches '%s'\n", date, matchDate, title, matchTitle)
}

// TestSpineStringRotation
// Using a specific string, test this function's ability to rotate through each character.
func TestSpineStringRotation(t *testing.T) {
	fmt.Printf("\n\t::: TestSpineStringRotation :::\n")

	spineString := "cra"
	ictus := 0
	nexus := 1

	var spineChars []string
	for h := 0; h < len(spineString); h++ {
		spineChars = append(spineChars, strings.ToLower(string(spineString[h])))
	}

	if spineChars[ictus] != "c" {
		t.Errorf("%q\n", spineString)
	}
	fmt.Println(ictus, spineChars[ictus])

	Ictus(len(spineString), &ictus, &nexus)

	if spineChars[ictus] != "r" {
		t.Errorf("%q\n", spineString)
	}
	fmt.Println(ictus, spineChars[ictus])

	Ictus(len(spineString), &ictus, &nexus)

	if spineChars[ictus] != "a" {
		t.Errorf("%q\n", spineString)
	}
	fmt.Println(ictus, spineChars[ictus])

	Ictus(len(spineString), &ictus, &nexus)

	if spineChars[ictus] != "c" {
		t.Errorf("Rotation failed! %q\n", spineString)
	}
	fmt.Println(ictus, spineChars[ictus])
}
