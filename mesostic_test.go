/*

	Tests

*/

package main

import (
	"fmt"
	"strings"
	"testing"
)

// TestTIctus ::: Using a specific string, test this function's ability to rotate through each character.
func TestTIctus(t *testing.T) {
	fmt.Printf("\n\t::: Test Target Ictus() :::\n")

	spineString := "cra"
	ictus := 0
	nexus := 1

	var spineChars []string
	for h := 0; h < len(spineString); h++ {
		spineChars = append(spineChars, strings.ToLower(string(spineString[h])))
	}

	// Test each character and then that it rotates to the first successfully.
	if spineChars[ictus] != "c" {
		t.Errorf("%q\n", spineString)
	}

	Ictus(len(spineString), &ictus, &nexus)

	if spineChars[ictus] != "r" {
		t.Errorf("%q\n", spineString)
	}

	Ictus(len(spineString), &ictus, &nexus)

	if spineChars[ictus] != "a" {
		t.Errorf("%q\n", spineString)
	}

	Ictus(len(spineString), &ictus, &nexus)

	if spineChars[ictus] != "c" {
		t.Errorf("Rotation failed! %q\n", spineString)
	}
}

// TestTPreus ::: Using a specific string, test this function's ability to rewind the SpineString rotation done by Ictus().
func TestTPreus(t *testing.T) {
	fmt.Printf("\n\t::: Test Target Preus() :::\n")

	spineString := "cra"
	ictus := 0
	nexus := 1

	var spineChars []string
	for h := 0; h < len(spineString); h++ {
		spineChars = append(spineChars, strings.ToLower(string(spineString[h])))
	}

	// Test each character and then that it rewinds to the one before it.
	if spineChars[ictus] != "c" {
		t.Errorf("%q\n", spineString)
	}

	Preus(len(spineString), &ictus, &nexus)

	if spineChars[ictus] != "a" {
		t.Errorf("%q\n", spineString)
	}

	Preus(len(spineString), &ictus, &nexus)

	if spineChars[ictus] != "r" {
		t.Errorf("%q\n", spineString)
	}

	Preus(len(spineString), &ictus, &nexus)

	if spineChars[ictus] != "c" {
		t.Errorf("Rewind failed! %q\n", spineString)
	}
}

/*

Ed's pillar test file:
https://github.com/Everbridge/generate-secure-pillar/blob/main/main_test.go

Interesting structure to try out:
https://github.com/rs/zerolog/blob/master/log_test.go

func TestMe() {
	t.Run("empty", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out)
		log.Log().Msg("")
		if got, want := decodeIfBinaryToString(out.Bytes()), "{}\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})

	t.Run ...

  }

*/
