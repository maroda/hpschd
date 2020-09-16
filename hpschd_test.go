/*

	This requires more, but global vars get in the way a lot.

*/

package main

import (
	"fmt"
	"strings"
	"testing"
)

// TestIctusIncr
func TestIctusOnce(t *testing.T) {
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
}
