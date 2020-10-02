/*

	Mesostic Data Operations

	Currently a lot of filesystem stuff,
		intended to be expandable into a database.

	Also configurations and external variable handling.

*/

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

// envVar ::: Grab a single ENV VAR and provide a fallback configuration.
func envVar(env, alt string) string {
	url, ext := os.LookupEnv(env)
	if !ext {
		url = alt
	}
	return url
}

// function to set up local directories if not already present
func localDirs() {
}

// ichingMeso ::: Uses chance operations to select an existing NASA APOD Mesostic.
func ichingMeso(dir string) string {
	var fileList []string
	for _, entry := range dirents(dir) {
		fullPath := filepath.Join(dir, entry.Name())
		fmt.Println(fullPath)
		fileList = append(fileList, fullPath)
	}

	// now return only one. for now, the first. soon, randomized.
	return fileList[0]
}

// read a directory and return its contents
func dirents(d string) []os.FileInfo {
	ents, err := ioutil.ReadDir(d)
	if err != nil {
		log.Error()
		return nil
	}
	return ents
}

// apodNEW ::: Check if a disk file exists in the Mesostic store or create a new one.
// The return values are the filename and whether the function wrote a new file.
func apodNew(sp *string, da *string, me *string) (string, bool) {
	mDir := "store"
	tr := strings.NewReplacer(" ", "_")
	spn := tr.Replace(*sp)
	fP := fmt.Sprintf("%s/%s__%s", mDir, *da, spn)

	// TODO: check how old it is and update if a retention period is met
	if _, err := os.Stat(fP); err == nil {
		return fP, false
	}

	// write data to a new file if it doesn't exist
	sB := []byte(*me)
	err := ioutil.WriteFile(fP, sB, 0644)
	if err != nil {
		log.Error()
	}
	return fP, true
}

// fileTmp ::: Take a source string and place it in a file name after the spinestring.
// This only creates the file by a straight byte copy.
// Calling functions are responsible for file deletion when finished.
func fileTmp(sp *string, so *string) string {
	fT := time.Now()
	fS := fT.Unix()
	fN := fmt.Sprintf("%s__%d", *sp, fS)
	sB := []byte(*so)
	err := ioutil.WriteFile(fN, sB, 0644)
	if err != nil {
		log.Error()
	}
	return fN
}

// nasaNewREAD ::: Consume the current filename for the current NASA APOD Mesostic.
// No new data returns the string 'HPSCHD'
func nasaNewREAD() string {
	_, _, fu := Envelope()

	// The purpose here is to display the current APOD first,
	// and random ones subsequently, including the present one.
	select {
	case mesoFile := <-nasaNewMESO:
		log.Info().Str("fu", fu).Msg("Filename from nasaNewMESO consumed")
		return mesoFile
	default:
		log.Info().Str("fu", fu).Msg("No new filename presented, initiate chance operations.")
		return "HPSCHD"
	}
}
