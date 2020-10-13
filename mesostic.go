/*

	Mesostic Engine

*/

package main

import (
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

// Hash table of line fragments
var fragMents = make(map[string]LineFrag)

// this will need to be passed as a pointer like the ictus/nexus stuff <<< ??? why did i write this
// for some reason when i pull this into the mesoLine() function,
// the Ictus rotation breaks, or the display of it breaks. not sure.
// but it only works when this is a global.
var fragCount int // global to count total fragment combinations (i.e. lines)

// Spine ::: Process the SpineString
//	Construct a slice of lowercase SpineString characters that can be rotated by Ictus().
func Spine(z string) []string {
	var zch []string

	for h := 0; h < len(z); h++ {
		zch = append(zch, strings.ToLower(string(z[h])))
	}
	return zch
}

// mesoLine ::: finds the current SpineString character in the current line
//
// The West Fragment is everything to the left of the SpineString character.
// The East Fragment is everything to the right of the SpineString character.
//
//		s == the current line to process
//		z == slice of SpineString characters
//		c == line number
//		ict == ictus of the SpineString characters
//		nex == next ictus (not always ict + 1)
//		spaces == Pointer ::: current left-aligned whitespace
//
func mesoLine(s string, z []string, c int, ict *int, nex *int, spaces *int) bool {
	var wstack []string   // slice for rebuilding the west fragment
	var estack []string   // slice for rebuilding the east fragment
	var found bool = true // assert the character will be found
	var mode int          // Choose the Mesostic algorithm mode

CharLoop:
	// step through the current string and process mesostic rules
	for i := 0; i < len(s); i++ {
		char := string(s[i])

		switch mode {
		// WestSide
		case 0:
			// as long as the character isn't the current Spine Character fill in the WestSide fragment.
			if char != z[*ict] {
				wstack = append(wstack, char)
				if char == string(s[len(s)-1]) {
					// if we've reached this point,
					// the SpineString character didn't appear.
					// currently, the line is thown out.
					log.Debug().Str("final", char).Msg("EOL")
					found = false
					wstack = nil
					wstack = append(wstack, " ")
				}
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
			if char != z[*nex] {
				estack = append(estack, char)
			} else {
				break CharLoop // We're done.
			}

			/*
				More Cases:

				A bug in the algorithm:
				If the SpineString char doesn't exist in the line at all don't print the line and retain the ictus.
				Currently if this happens the line is printed without a SpineString character and moves on.

				100% Mesostic ::: No occurance of the current Spine String Character in back OR in front of it.
				Meso-Acrostic ::: No pre/post rules, any character can appear before or after the Spine String Char.
			*/
		}
	}

	// Post processing
	fragmentW := strings.Join(wstack, "")                // WestSide fragment
	fragmentE := strings.Join(estack, "")                // EastSide fragment
	fragkey := shakey(fragmentW + fmt.Sprint(fragCount)) // unique identifier and consistent key sizes
	fragCount++

	// Add results to a new map entry
	// Some or all of this might be better off as pointers...
	fragMents[fragkey] = LineFrag{Index: c, LineNum: fragCount, WChars: len(fragmentW), Data: fragmentW + fragmentE}

	// record the longest WestSide fragment length, but calculate it from the passed value
	if len(fragmentW) > *spaces {
		*spaces = len(fragmentW)
	}

	return found
}

// Ictus ::: Enables the rotation of SpineString characters by operating on the index.
//
//			lss == length of SpineString
//			isp == pointer to the ictus
//			nsp == pointer to the next ictus
//
func Ictus(lss int, isp *int, nsp *int) {
	// a mesostic line has been finished,
	// increase ictus, i.e. the current character position
	if *isp < lss-1 {
		*isp++
	} else if *isp == lss-1 { // last element, rotate
		*isp = 0
	}

	// this is used for lookahead
	//  update nexus, i.e. the next character position
	//	 if we've hit the end, we wrap to the first character
	if *isp == lss-1 { // now last element, next is 0
		*nsp = 0
	} else {
		*nsp = *isp + 1
	}
}

// Preus ::: Ictus, but rewind.
func Preus(lss int, isp *int, nsp *int) {
	if *isp > 0 { // not first element, simple subtraction
		*isp--
	} else { // now first element, rotate backwards
		*isp = lss - 1
	}

	// nexus is still calcualted the same way
	if *isp == lss-1 {
		*nsp = 0
	} else {
		*nsp = *isp + 1
	}
}

// Sort Interface ::: linefragments by LineNum (lnc)
func (ls LineFrags) Len() int {
	return len(ls)
}
func (ls LineFrags) Swap(i, j int) {
	ls[i], ls[j] = ls[j], ls[i]
}
func (ls LineFrags) Less(i, j int) bool {
	return ls[i].LineNum < ls[j].LineNum
}

// mesoMain ::: Takes a filename as input for processing.
//	Alternate main()
//	TODO: only launch the api if a "server" flag is given.
//			Otherwise, it's a standalone CLI tool.
//
// f == filename for processing
// z == Spine String
// o == channel for return
//
func mesoMain(f string, z string, o chan<- string) {
	var lnc int               // line counts for the Index
	var ictus int             // SpineString character address
	var nexus int = ictus + 1 // Next SpineString character address
	var spaces int = 0        // Left-aligned whitespace for all lines

	// split the SpineString into a slice of characters
	spineChars := Spine(z)

	spineString := strings.Join(spineChars, "")
	// DEBUG ::: fmt.Sprint(spineString)

	source, err := ioutil.ReadFile(f)
	if err != nil {
		log.Error()
	}

	// Break down the file into lines and manipulate them into a new unordered mesostic
	for _, sline := range strings.Split(string(source), "\n") {
		lnc++

		// mesoLine populates a global map, returning a boolean success status
		success := mesoLine(strings.ToLower(sline), spineChars, lnc, &ictus, &nexus, &spaces)

		log.Debug().
			Int("lnc", lnc).
			Int("ictus", ictus).
			Int("spaces", spaces).
			Bool("success", success).
			Msg("")

		// Once the mesostic line has been created and added to the data map,
		// operate on the Ictus values to construct the next line
		if !success {
			// The character wasn't found, rewind the SpineString.
			// The following Ictus() call makes this effectively "staying still"

			/// UGH this is not working
			// Preus(len(spineString), &ictus, &nexus)
		}
		Ictus(len(spineString), &ictus, &nexus)
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

	for i := 0; i < len(linefragments); i++ {
		// define 'West Side' whitespace as
		//  (length of the longest fragment) - (length of the current fragment)
		padMe := spaces - linefragments[i].WChars
		printspace := strings.Repeat(" ", padMe)

		// format the new line with leading whitespace and trailing line return
		fragstack = append(fragstack, printspace)
		fragstack = append(fragstack, linefragments[i].Data)
		fragstack = append(fragstack, "\n")
	}
	mesostic := strings.Join(fragstack, "")
	o <- fmt.Sprint(mesostic)
	close(o)
}
