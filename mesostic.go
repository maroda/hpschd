/*

	Mesostic Engine

*/

package main

import (
	"fmt"
	"io/ioutil"
	"os"
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
	var wstack []string // slice for rebuilding the west fragment
	var estack []string // slice for rebuilding the east fragment
	var found bool      // the character was found in this line
	mode := 0           // the Mesostic algorithm mode, always starts with 0

CharLoop:
	// step through the current string and process mesostic rules
	for i := 0; i < len(s); i++ {
		char := string(s[i])

		/*
			WestSide ::: Everything to the LEFT AND INCLUDING the SpineString, this is mode 0
			EastSide ::: Everything to the RIGHT AND NOT the SpineString, includes:

			mode 1 50% Mesostic ::: No occurance of the current Spine String Character in back of it.
			mode 2 100% Mesostic ::: No occurance of the current Spine String Character in back OR in front of it.
			mode 3 Meso-Acrostic ::: No pre/post rules, any character can appear before or after the Spine String Char.
		*/
		switch mode {
		case 0:
			switch {
			case char != z[*ict]:
				// not found
				log.Debug().Str("z", z[*ict]).Msg("notfound")
				if i == len(s)-1 {
					// last character of the line
					if char != z[*ict] {
						// the final character is not the SpineString
						// in this version, the line is thrown out
						// in future versions, the line may only be thrown out
						// if a certain tolerance is met
						// e.g.: the number of false successes
						log.Debug().Str("SSCHAR", z[*ict]).Str("LNCHAR", char).Msg("EOL")
						found = false
						wstack = nil
						break CharLoop // We're done.
					}
				}
				wstack = append(wstack, char)
			case char == z[*ict]:
				// SpineString hit!
				log.Debug().Str("z", z[*ict]).Msg("spinestr")
				found = true
				char = strings.ToUpper(char)  // Spine Character is capitalized
				wstack = append(wstack, char) // Appended to the string
				mode = 1
				break // re-evaluate the switch with mode set
			}
		case 1:
			/*
				Allow the current SSchar in the remainder of the line,
				but do not print anything on this line at or beyond the next SSchar
				because that will appear on the next line and cannot have itself preceeding it.

				This method preserves the line returns found in the source.
			*/
			if char != z[*nex] {
				estack = append(estack, char)
			} else {
				break CharLoop // We're done.
			}
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

	/*
		Break down the file into lines and manipulate them into a new unordered mesostic.
		mesoLine() populates a global map, returning a boolean success status.

		If the SpineString character was found,
			Ictus() rotates the SpineString position forward one spot,
			and the next line processed will have the next character to match.
		If the SpineString character wasn't found,
			rewind the SpineString with Preus(),
			making the SpineString character "stay in place", so it is matched in the next line.

			Currently, if the SSchar is never found, the Mesostic will be effectively blank.
			There might need to be a tolerance setting here:
			If not found X number of times (say, a fraction of the length of source), then rotate fwd.
	*/
	for _, sline := range strings.Split(string(source), "\n") {
		lnc++

		success := mesoLine(strings.ToLower(sline), spineChars, lnc, &ictus, &nexus, &spaces)
		if !success {
			Preus(len(spineString), &ictus, &nexus)
		}
		Ictus(len(spineString), &ictus, &nexus)

		log.Debug().
			Int("lnc", lnc).
			Int("ictus", ictus).
			Int("spaces", spaces).
			Bool("success", success).
			Msg("")
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

	// tmp scratch is no longer needed
	var ferr = os.Remove(f)
	if ferr != nil {
		log.Error()
	}
}
