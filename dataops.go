/*

	Mesostic Data Operations

	- Filesystem / database access
	- Specialized random / hash values
	- Configurations

*/

package main

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"math/rand/v2"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
)

// rndDate ::: Produce a random date in the format YYYY-MM-DD.
func rndDate() string {
	// rand ranges are [0,r)
	rMi := 20 // Millinium
	rYr := 20 // Years
	rMo := 12 // Months
	rDy := 31 // Days

	// No random for millinium
	Mi := fmt.Sprint(rMi)

	// Yr can be zero
	Yr := fmt.Sprintf("%02d", rand.IntN(rYr))

	// Don't actually use the last number but then add it back.
	Mo := fmt.Sprintf("%02d", rand.IntN(rMo-1)+1)

	// Good thing for a test:
	// 	In rare cases this may be > 31,
	// 	but the API should return a 404
	// 	and that will trigger another random selection anyway.
	Dy := fmt.Sprintf("%02d", rand.IntN(rDy)+1)

	// Formatted YYYY-MM-DD date
	newdate := Mi + Yr + "-" + Mo + "-" + Dy

	return newdate
}

// envVar ::: Grab a single ENV VAR with a provided default
// This will not set an ENV VAR that exists but set to an empty string.
func envVar(env, alt string) string {
	value, ext := os.LookupEnv(env)
	if !ext {
		return alt
	}
	return value
}

// ichingMeso ::: Uses chance operations to select an existing NASA APOD Mesostic.
func ichingMeso(dir string) string {
	var fileList []string
	for _, entry := range dirents(dir) {
		fullPath := filepath.Join(dir, entry.Name())
		fileList = append(fileList, fullPath)
	}
	if fileList == nil {
		log.Error().Msg("ENOENT ::: Is the datastore available?")
		return "ENOENT"
	}

	randix := rand.IntN(len(fileList))
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
// TODO: Refactor to more modern file access
func readMesoFile(ctx context.Context, f *string) string {
	ctx, span := otel.Tracer("mesostic/web").Start(ctx, "readMesoFile")
	defer span.End()

	if len(*f) == 0 {
		err := fmt.Errorf("no file specified")
		span.RecordError(err)
		slog.Error("no file specified", slog.Any("error", err))
		return "ENOENT"
	}

	var mesoBuf []byte
	mesoBuf, err := os.ReadFile(*f)
	if err != nil {
		span.RecordError(err)
		slog.Error("error reading file", slog.Any("error", err))
		return "ENOENT"
	}

	return string(mesoBuf)
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
