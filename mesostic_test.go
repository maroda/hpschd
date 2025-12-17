/*

	Tests

*/

package main

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

// TestTmesoMain ::: Create a mesostic from a static source file.
// For now this test is display-only.
// TODO: Checksum or something.
func TestTmesoMain(t *testing.T) {
	fmt.Printf("\n\t::: Test Target mesoMain() :::\n")

	// This mimics the calls made by various API endpoints.
	// There is no automatic filename creation, a known file is used to match.
	// this test is "failing" because the mesoMain function now *removes* the file it is passed.

	// Set up store
	TTstore, err := os.MkdirTemp(".", "store")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(TTstore)

	// Set up tmp
	TTdir, err := os.MkdirTemp(".", "txrx")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(TTdir)

	spine := "craque"
	sourceFile := "sources/lorenipsum-plaintext.txt"

	// simulate reading an input source
	bRead, berr := os.ReadFile(sourceFile)
	if berr != nil {
		t.Error(berr)
	}

	// write a scratch tmp file in the test temp directory, which mesoMain removes
	TTfile, err := os.CreateTemp(TTdir, "lorenipsum_*.txt")
	if err != nil {
		t.Error(err)
	}
	testTmp := TTfile.Name()
	TTfile.Close()

	werr := os.WriteFile(testTmp, bRead, 0644)
	if werr != nil {
		t.Error(werr)
	}

	// mesoMain receives ::: tmp filename, the SpineString, data channel
	mcMeso := make(chan string)
	go mesoMain(testTmp, spine, mcMeso)

	// receive the channel data and display result
	mesostic := <-mcMeso
	fmt.Println(mesostic)

	// check if mesoMain removed the scratch tmp correctly
	if _, lserr := os.Stat(testTmp); lserr == nil {
		t.Errorf("mesoMain() did not remove temp file '%s' as expected", testTmp)
	}

	/*
	                       lorem ipsum dolor sit amet, Consectetu
	                         elit, sed do eiusmod tempoR incididunt ut l
	                                           dolore mAgna ali
	   nostrud exercitation ullamco laboris nisi ut aliQ
	                               ex ea commodo conseqUat. duis aut
	                                                  rEprehenderit in voluptate velit esse
	                       eu fugiat nulla pariatur. exCepteu
	                                    cupidatat non pRoident, sunt in culp
	                                   deserunt mollit Anim id est laborum.


	*/
}

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
