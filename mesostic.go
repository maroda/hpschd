package main

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
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

var padCount int                          // global to track right-align padding
var fragCount int                         // global to count total fragment combinations (i.e. lines)
var fragMents = make(map[string]LineFrag) // Hash table of line fragments
var ss string = "craque"                  // SpineString
var sca []string                          // Characters in SpineString
var ictus int                             // SpineString character address

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
			// The WestSide fragment is complete, fill in the remainder for the EastSide fragment.
			/*

				TODO: This isn't actually checking for the next char yet,
				it's checking for the current one. It should look ahead,
				and that might be another option to the Ictus() function.

				For example, Ictus() could also easily increment a "next value"
				so that other resources don't need to do math to navigate
				the mesostic ruleset.

			*/
			if char != sca[ictus] {
				estack = append(estack, char)
			} else {
				break CharLoop // We're done.
			}
		}
	}

	// Post processing
	Ictus(1)                                         // Increases the rotation of the SpineString
	fragmentW := strings.Join(wstack, "")            // WestSide fragment
	fragmentE := strings.Join(estack, "")            // EastSide fragment
	fragCount++                                      // new line number of the current two fragments
	fragkey := shakey(fragmentW + string(fragCount)) // unique identifier and consistent key sizes

	/*

		TODO: Lines without the current StringChar are being printed without that StringChar and the StringChar advances anyway

	*/

	// Some or all of this might be better off as pointers...
	fragMents[fragkey] = LineFrag{Index: c, LineNum: fragCount, WChars: len(fragmentW), Data: fragmentW + fragmentE}

	// record the longest WestSide fragment length
	if len(fragmentW) > padCount {
		padCount = len(fragmentW)
	}
}

// Ictus :: Rotates SpineString characters
func Ictus(i int) {
	if ictus < len(ss)-1 {
		ictus += i
	} else if ictus == len(ss)-1 {
		ictus = 0
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

func main() {
	// line counts for the Index
	var lnc int

	for h := 0; h < len(ss); h++ {
		sca = append(sca, string(ss[h]))
	}

	// process the given file, silent if no file given
	for _, origtxt := range os.Args[1:] {
		data, err := ioutil.ReadFile(origtxt)
		if err != nil {
			log.Fatal(err)
			break
		}

		// run mesoLine - which needs a new name - to process each line
		for _, line := range strings.Split(string(data), "\n") {
			lnc++
			mesoLine(strings.ToLower(line), lnc)
		}
	}

	// Sort & Print //
	//
	// Lines to be sorted are pushed to a slice.
	// Sort is configured on LineNum.
	// Uneven padding is accomplished by subtracting
	//   the length of the current WestSide fragment
	//   from the length of the longest WestSide fragment (padCount)

	var linefragments LineFrags
	for k := range fragMents {
		linefragments = append(linefragments, fragMents[k])
	}

	sort.Sort(linefragments)

	for i := 0; i < len(linefragments); i++ {
		padMe := padCount - linefragments[i].WChars
		spaces := strings.Repeat(" ", padMe)
		fmt.Printf("%s%s\n", spaces, linefragments[i].Data)
	}
}
