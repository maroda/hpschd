/*

	Data Operations Tests

*/

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestTlocalDirs ::: Create a local dirset, remove when done.
func TestTlocalDirs(t *testing.T) {
	fmt.Printf("\n\t::: Test Target localDirs() :::\n")

	// testing directory
	testDir := "./TTlocal_" + fmt.Sprint(time.Now().Unix())

	// configured directories
	locals := []string{"TTstore, maTTest"}

	err := os.Mkdir(testDir, 0755)
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(testDir)

	cderr := os.Chdir(testDir)
	if cderr != nil {
		t.Error(err)
	}

	// Create and then check the existence of each configured directory.
	localDirs(locals)
	for _, dir := range locals {
		if _, lderr := os.Stat(dir); lderr != nil {
			t.Error()
		}
	}

	cberr := os.Chdir("../")
	if cberr != nil {
		t.Error(cberr)
	}
}

// TestTextent ::: Stat a known file, stat an unknown file.
func TestTextent(t *testing.T) {
	fmt.Printf("\n\t::: Test Target extent() :::\n")

	knownFile := "/etc/passwd"
	unknownFile := "/tmp/nofile" + fmt.Sprint(time.Now().Unix())

	if !extent(knownFile) {
		t.Error()
	}

	if extent(unknownFile) {
		t.Error()
	}
}

// TestTdirents ::: Create a tmp dir and file in it, get the dir list, match the file in the list.
func TestTdirents(t *testing.T) {
	fmt.Printf("\n\t::: Test Target dirents() :::\n")

	// Set up tmp
	TTdir, err := os.MkdirTemp(".", "TT")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(TTdir)

	TTfile, err := os.CreateTemp(TTdir, "")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(TTfile.Name())

	// Call dirents()
	for _, entry := range dirents(TTdir) {
		pathENT := filepath.Clean(filepath.Join(TTdir, entry.Name()))
		pathTT := filepath.Clean(TTfile.Name())
		if pathENT != pathTT {
			t.Errorf("Local file '%s' does NOT match created file '%s'.\n", pathENT, pathTT)
		}
	}
}

// TestTichingMeso ::: When working, this will test the presence of a random filename. ???
func TestTichingMeso(t *testing.T) {
	fmt.Printf("\n\t::: Test Target ichingMeso() :::\n")

	// Set up tmp
	TTdir, err := os.MkdirTemp(".", "TT")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(TTdir)

	TTfile, err := os.CreateTemp(TTdir, "")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(TTfile.Name())

	// Call ichingMeso()
}

// TestTfileTmp ::: Special Mesostic TMP file creation.
// Given a set of strings, match the created filename and verify its presence on disk,
// perhaps match the content itself, then delete.
func TestTfileTmp(t *testing.T) {
	fmt.Printf("\n\t::: Test Target fileTmp() :::\n")

	// Set up tmp
	TTdir, err := os.MkdirTemp(".", "txrx")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(TTdir)

	spine := "cra"
	source := "que"
	fileName := fileTmp(&spine, &source)
	if !strings.Contains(fileName, spine) {
		t.Errorf("Local filename '%s' does not contain '%s'.\n", fileName, spine)
	}
}

// TestTenvVar ::: Process environment variables correctly with a given fallback option.
func TestTenvVar(t *testing.T) {
	fmt.Printf("\n\t::: Test Target envVar() :::\n")

	var getvar string
	var testvar string
	var testval string
	var fallval string

	// testvar does not exist, fallback provided (good config)
	// expected return: the fallback value
	fallval = "fallback_NoVAR"
	getvar = envVar(testvar, fallval)
	if getvar != fallval {
		t.Errorf("%s, %s, %s", testvar, testval, fallval)
	}
	t.Logf("fallback received: %s", getvar)

	// testvar exists, but is unset, no fallback (error condition)
	// correct return: empty
	testvar = "TTVAR"
	fallval = ""
	getvar = envVar(testvar, fallval)
	if getvar != "" {
		t.Errorf("%s, %s, %s", testvar, testval, fallval)
	}
	t.Logf("empty received: %s", getvar)

	// testvar exists, but is unset, fallback provided (good config)
	// correct return: value for fallval
	fallval = "fallback_NoValue"
	getvar = envVar(testvar, fallval)
	if getvar != fallval {
		t.Errorf("%s, %s, %s", testvar, testval, fallval)
	}
	t.Logf("fallback received: %s", getvar)

	// Finally testvar is set, fallback provided (good config)
	// correct return: value for testval
	testval = "TestTenvVar"
	fallval = "fallback_NoValue"
	os.Setenv(testvar, testval) // testvar := testval
	getvar = envVar(testvar, fallval)
	if getvar != testval {
		t.Errorf("%s, %s, %s", testvar, testval, fallval)
	}
	t.Logf("set value received: %s", getvar)
}

// TestTrndDate ::: Test the creation of a random date
func TestTrndDate(t *testing.T) {
	fmt.Printf("\n\t::: Test Target rndDate() :::\n")

	salt := time.Now().Unix()
	randomdate := rndDate(salt)

	// split the string by "-"
	// now test by checking its values mathematically.
	// see if the random values fall within the ranges given?

	fmt.Println(randomdate)
}
