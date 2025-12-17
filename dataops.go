/*

	Mesostic Data Operations

	- Filesystem / database access
	- Specialized random / hash values
	- Configurations

*/

package main

import (
	"crypto/sha1"
	"fmt"
	"io/fs"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

// rndDate ::: Produce a random date in the format YYYY-MM-DD.
func rndDate(salt int64) string {

	// rand ranges are [0,r)
	rMi := 20 // Millinium
	rYr := 20 // Years
	rMo := 12 // Months
	rDy := 31 // Days

	rand.Seed(salt)

	// No random for millinium
	Mi := fmt.Sprint(rMi)

	// Yr can be zero
	Yr := fmt.Sprintf("%02d", rand.Intn(rYr))

	// Don't actually use the last number but then add it back.
	Mo := fmt.Sprintf("%02d", rand.Intn(rMo-1)+1)

	// Good thing for a test:
	// 	In rare cases this may be > 31,
	// 	but the API should return a 404
	// 	and that will trigger another random selection anyway.
	Dy := fmt.Sprintf("%02d", rand.Intn(rDy)+1)

	// Formatted YYYY-MM-DD date
	newdate := Mi + Yr + "-" + Mo + "-" + Dy

	return newdate
}

// envVar ::: Grab a single ENV VAR and provide a fallback configuration.
func envVar(env, alt string) string {
	url, ext := os.LookupEnv(env)
	if !ext {
		url = alt
	}
	return url
}

// ichingMeso ::: Uses chance operations to select an existing NASA APOD Mesostic.
func ichingMeso(dir string) string {
	var fileList []string
	for _, entry := range dirents(dir) {
		fullPath := filepath.Join(dir, entry.Name())
		fmt.Println(fullPath)
		fileList = append(fileList, fullPath)
	}
	if fileList == nil {
		log.Error().Msg("ENOENT ::: Is the datastore available?")
		return "ENOENT"
	}

	rand.Seed(time.Now().Unix())
	randix := rand.Intn(len(fileList))
	return fileList[randix]
}

// dirents ::: read a directory and return its contents
func dirents(d string) []fs.DirEntry {
	ents, err := os.ReadDir(d)
	if err != nil {
		log.Error()
		return nil
	}
	return ents
}

// extent ::: file system entry exists
func extent(fs string) bool {
	if _, err := os.Stat(fs); err != nil {
		return false
	}
	return true
}

// localDirs ::: set up permanent data directories
func localDirs(ld []string) {
	for _, dir := range ld {
		if !extent(dir) {
			log.Info().Str("directory", dir).Msg("Dir not found, creating.")
			err := os.Mkdir(dir, 0700)
			if err != nil {
				log.Error()
			}
		}
	}
}

// readMesoFile ::: Open and read the Mesostic
func readMesoFile(f *string) string {
	if len(*f) == 0 {
		log.Error().Msg("no path given")
		return "error"
	}

	var mesoBuf []byte
	mesoBuf, err := os.ReadFile(*f)
	if err != nil {
		log.Error()
	}

	return string(mesoBuf)
}

// apodNEW ::: Check if a disk file exists in the Mesostic store or create a new one.
// The return values are the filename and whether the function wrote a new file.
func apodNew(sp *string, da *string, me *string) (string, bool) {
	_, _, fu := Envelope()

	mDir := "store"
	tr := strings.NewReplacer(" ", "_")
	spn := tr.Replace(*sp)
	fP := fmt.Sprintf("%s/%s__%s", mDir, *da, spn)

	// NASA API returned a 404
	if spn == "404" {
		log.Warn().Str("fu", fu).Msg("404 NOFILE")
		return fP, false
	}

	// Mesostic file exists
	if _, err := os.Stat(fP); err == nil {
		log.Warn().Str("fu", fu).Msg("EXISTENT")
		return fP, false
	}

	// Write data to a new file
	sB := []byte(*me)
	err := os.WriteFile(fP, sB, 0644)
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
	fN := fmt.Sprintf("txrx/%s__%d", *sp, fS)
	sB := []byte(*so)
	err := os.WriteFile(fN, sB, 0644)
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

// SHA1 for consistent size keys
func shakey(k string) string {
	s := sha1.New()
	s.Write([]byte(k))
	bash := s.Sum(nil)
	hash := fmt.Sprintf("%x", bash)
	return hash
}
