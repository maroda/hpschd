/*

	Mesostic Engine

	BUG ::: if the number of lines in the source are less than the cardinality of the string,
		the code will wrap the spine string **and leave it that way in memory**

		something is needed here that resets the spine string for each goroutine that is run

		which may mean it cannot be treated globally... :(
		  the solution may ultimately be to re-do the global variable nature of
		  the mesostic processing variables and data map and use args and pointers.

	BUG ::: ??? Unsure what's going on here, but after left running for several days on ECS,
		the whitespace padding disappeared. When the tasks/containers were restarted,
		the whitespace padding returned. ¯\_(ツ)_/¯

*/

package main

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"

	"github.com/rs/zerolog/log"
)

// LineFrag ::: Data model describing a processed LineFragment.
type LineFrag struct {
	Index   int    // Line number from the original text.
	LineNum int    // The assigned line number for these fragments.
	WChars  int    // WestSide character count.
	Data    string // The new Mesostic line.
}

// LineFrags ::: string slice for the collection of LineFrag entries to be sorted
type LineFrags []LineFrag

/* These globals should be changed to... ???  */

var padCount int                          // global to track right-align padding
var fragCount int                         // global to count total fragment combinations (i.e. lines)
var fragMents = make(map[string]LineFrag) // Hash table of line fragments

var ss string = "craque"  // SpineString
var sca []string          // Characters in SpineString
var ictus int             // SpineString character address
var nexus int = ictus + 1 // Next SpineString character address

// Spine ::: Process the SpineString
func Spine() {
	for h := 0; h < len(ss); h++ {
		sca = append(sca, string(ss[h]))
	}
}

// mesoLine ::: finds the SpineString (SS) characters (sc)
func mesoLine(s string, c int) {
	var wstack []string // slice for rebuilding the west fragment
	var estack []string // slice for rebuilding the east fragment
	var mode int

CharLoop:
	// step through the current string and process mesostic rules
	for i := 0; i < len(s); i++ {
		char := string(s[i])

		switch mode {
		// WestSide
		case 0:
			// as long as the character isn't the current Spine Character fill in the WestSide fragment.
			if char != sca[ictus] {
				wstack = append(wstack, char)
			} else {
				// SpineString hit!
				char = strings.ToUpper(char)  // Spine Character is capitalized
				wstack = append(wstack, char) // Appended to the string
				mode = 1
				break // re-evaluate the switch with mode set
			}
		// EastSide
		case 1:
			/*
				50% Mesostic ::: No occurance of the current Spine String Character in back of it.

				Allow the current SpineChar in the remainder of the line,
				but do not print anything on this line at or beyond the next SpinChar
				because that will appear on the next line and cannot have itself preceeding it.
			*/
			if char != sca[nexus] {
				estack = append(estack, char)
			} else {
				break CharLoop // We're done.
			}

			// 100% Mesostic ::: No occurance of the current Spine String Character in back OR in front of it.
			// Meso-Acrostic ::: No pre/post rules, any character can appear before or after the Spine String Char.
		}
	}

	// Post processing
	fragmentW := strings.Join(wstack, "")            // WestSide fragment
	fragmentE := strings.Join(estack, "")            // EastSide fragment
	fragkey := shakey(fragmentW + string(fragCount)) // unique identifier and consistent key sizes
	if len(fragmentW) > padCount {                   // record the longest WestSide fragment length
		padCount = len(fragmentW)
	}
	fragCount++

	// Some or all of this might be better off as pointers...
	fragMents[fragkey] = LineFrag{Index: c, LineNum: fragCount, WChars: len(fragmentW), Data: fragmentW + fragmentE}

	Ictus(1) // Increases the rotation of the SpineString, i.e. the next address in sca[ictus].
}

// Ictus :: Rotates SpineString characters
func Ictus(i int) {
	// update global ictus
	if ictus < len(ss)-1 {
		ictus += i
	} else if ictus == len(ss)-1 { // last element, rotate
		ictus = 0
	}

	// update global nexus
	if ictus == len(ss)-1 { // now last element, next is 0
		nexus = 0
	} else {
		nexus = ictus + 1
	}
}

// Sort ::: linefragments by LineNum
func (ls LineFrags) Len() int {
	return len(ls)
}
func (ls LineFrags) Swap(i, j int) {
	ls[i], ls[j] = ls[j], ls[i]
}
func (ls LineFrags) Less(i, j int) bool {
	return ls[i].LineNum < ls[j].LineNum
}

// SHA1 for consistent size keys
func shakey(k string) string {
	s := sha1.New()
	s.Write([]byte(k))
	bash := s.Sum(nil)
	hash := fmt.Sprintf("%x", bash)
	return hash
}

// mesoMain ::: Takes a filename as input for processing.
//	Alternate main()
//	TODO: only launch the api if a "server" flag is given.
//			Otherwise, it's a standalone CLI tool.
//
func mesoMain(f string, o chan<- string) {
	var lnc int // line counts for the Index

	Spine()

	source, err := ioutil.ReadFile(f)
	if err != nil {
		log.Error()
	}

	for _, sline := range strings.Split(string(source), "\n") {
		lnc++
		mesoLine(strings.ToLower(sline), lnc)
	}

	// Sort & Print //

	// Lines are moved from the global map to a slice to be sorted.
	var fragstack []string
	var linefragments LineFrags
	for k := range fragMents {
		linefragments = append(linefragments, fragMents[k])
		delete(fragMents, k)
	}

	// Sort is configured on LineNum.
	sort.Sort(linefragments)

	// Uneven padding is accomplished by subtracting
	//   the length of the current WestSide fragment
	//   from the length of the longest WestSide fragment (padCount)
	for i := 0; i < len(linefragments); i++ {
		padMe := padCount - linefragments[i].WChars
		spaces := strings.Repeat(" ", padMe)
		fragstack = append(fragstack, spaces)
		fragstack = append(fragstack, linefragments[i].Data)
		fragstack = append(fragstack, "\n")
	}
	mesostic := strings.Join(fragstack, "")
	o <- fmt.Sprint(mesostic)
	close(o)
}
