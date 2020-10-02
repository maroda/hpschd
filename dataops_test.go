/*

	Data Operations Tests

*/

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestTdirents ::: Create a tmp dir and file in it, get the dir list, match the file in the list.
func TestTdirents(t *testing.T) {
	fmt.Printf("\n\t::: Test Target dirents() :::\n")

	// Set up tmp
	TTdir, err := ioutil.TempDir(".", "TT")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(TTdir)

	TTfile, err := ioutil.TempFile(TTdir, "")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(TTfile.Name())

	// Call dirents()
	for _, entry := range dirents(TTdir) {
		pathENT := filepath.Join(TTdir, entry.Name())
		pathTT := TTfile.Name()
		if pathENT != pathTT {
			t.Errorf("Local file '%s' does NOT match created file '%s'.\n", pathENT, pathTT)
		}
	}
}

// TestTichingMeso ::: When working, this will test the presence of a random filename. ???
func TestTichingMeso(t *testing.T) {
	fmt.Printf("\n\t::: Test Target ichingMeso() :::\n")

	// Set up tmp
	TTdir, err := ioutil.TempDir(".", "TT")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(TTdir)

	TTfile, err := ioutil.TempFile(TTdir, "")
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

	spine := "cra"
	source := "que"
	fileName := fileTmp(&spine, &source)
	if !strings.HasPrefix(fileName, spine) {
		t.Errorf("Local filename '%s' does not contain '%s'.\n", fileName, spine)
	}

	// fileTmp does not remove the file because it's used by multiple functions
	// so that a file removal of the passed filename works fine.
	var ferr = os.Remove(fileName)
	if ferr != nil {
		t.Error(ferr)
	}
}
