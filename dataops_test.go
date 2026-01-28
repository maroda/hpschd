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

// TODO: Most of dataops_test is from v1 and can be improved

// TestTlocalDirs ::: Create a local dirset, remove when done.
func TestTlocalDirs(t *testing.T) {
	t.Run("Test target localDirs()", func(t *testing.T) {
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
	})
}

// TestTextent ::: Stat a known file, stat an unknown file.
func TestTextent(t *testing.T) {
	t.Run("Test Target extent()", func(t *testing.T) {
		knownFile := "/etc/passwd"
		unknownFile := "/tmp/nofile" + fmt.Sprint(time.Now().Unix())

		if !extent(knownFile) {
			t.Error()
		}

		if extent(unknownFile) {
			t.Error()
		}
	})
}

// TestTdirents ::: Create a tmp dir and file in it, get the dir list, match the file in the list.
func TestTdirents(t *testing.T) {
	t.Run("Test Tdirents()", func(t *testing.T) {
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
	})
}

// TestTfileTmp ::: Special Mesostic TMP file creation.
// Given a set of strings, match the created filename and verify its presence on disk,
// perhaps match the content itself, then delete.
func TestTfileTmp(t *testing.T) {
	t.Run("Test Target fileTmp()", func(t *testing.T) {
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
	})
}

// TestTenvVar ::: Process environment variables correctly with a given fallback option.
func TestTenvVar(t *testing.T) {
	t.Run("Test Target envVar()", func(t *testing.T) {
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

		// testvar exists, but is unset, no fallback (error condition)
		// correct return: empty
		testvar = "TTVAR"
		fallval = ""
		getvar = envVar(testvar, fallval)
		if getvar != "" {
			t.Errorf("%s, %s, %s", testvar, testval, fallval)
		}

		// testvar exists, but is unset, fallback provided (good config)
		// correct return: value for fallval
		fallval = "fallback_NoValue"
		getvar = envVar(testvar, fallval)
		if getvar != fallval {
			t.Errorf("%s, %s, %s", testvar, testval, fallval)
		}

		// Finally testvar is set, fallback provided (good config)
		// correct return: value for testval
		testval = "TestTenvVar"
		fallval = "fallback_NoValue"
		os.Setenv(testvar, testval) // testvar := testval
		getvar = envVar(testvar, fallval)
		if getvar != testval {
			t.Errorf("%s, %s, %s", testvar, testval, fallval)
		}
	})
}

// TestTrndDate ::: Test the creation of a random date
func TestTrndDate(t *testing.T) {
	fmt.Printf("\n\t::: Test Target rndDate() :::\n")
	t.Run("RandomDate", func(t *testing.T) {
		randomdate := rndDate()
		t.Logf("random date: %s", randomdate)
	})
}
