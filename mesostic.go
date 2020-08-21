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

// string slice for the collection of LineFrag entries to be sorted
type LineFrags []LineFrag

var padCount int                          // global to track right-align padding
var fragCount int                         // global to count total fragment combinations (i.e. lines)
var fragMents = make(map[string]LineFrag) // Hash table of line fragments

// mesoLine ::: finds the SpineString (SS) characters (sc)
func mesoLine(s string, ss string, c int) {
	var wstack []string // slice for rebuilding the west fragment
	var estack []string // slice for rebuilding the east fragment
	var sca []string    // Spine Characters in SpineString
	var mode int

	// step through each character of the string... is what we want to do here. but that is the tough part.
	// get formatting fixed first with just a single character of the passed SpineString (ss)
	for h := 0; h < len(ss); h++ {
		sca = append(sca, string(ss[h]))
	}

	// step through the current string and process mesostic rules
	for i := 0; i < len(s); i++ {
		appC := string(s[i])

		switch mode {
		// WestSide
		case 0:
			// as long as the character isn't the current Spine Character fill in the WestSide fragment.
			if appC != sca[0] {
				wstack = append(wstack, appC)
			} else {
				// SpineString hit!
				appC = strings.ToUpper(appC)  // Spine Character is capitalized
				wstack = append(wstack, appC) // Appended to the string
				mode = 1
				break
			}
		// EastSide
		case 1:
			// The WestSide fragment is complete, fill in the remainder for the EastSide fragment.
			//
			// Make this configurable?
			//
			estack = append(estack, appC) // Full rest of line
			// x := 0 + 1
			//if appC != sca[x] {
			//	estack = append(estack, appC)
			//} else {
			//	break
			//}
		}
	}

	fragmentW := strings.Join(wstack, "")
	fragmentE := strings.Join(estack, "")
	fragCount++                                      // new line number of the current two fragments
	fragkey := shakey(fragmentW + string(fragCount)) // unique identifier and consistent key sizes

	// Some or all of this might be better off as pointers...
	fragMents[fragkey] = LineFrag{Index: c, LineNum: fragCount, WChars: len(fragmentW), Data: fragmentW + fragmentE}

	// the padding algo needs to be different. it will have to be the longest - present = padding
	// record the longest WestSide fragment length
	if len(fragmentW) > padCount {
		padCount = len(fragmentW)
	}
}

func main() {
	// line counts for the Index
	var lnc int

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
			mesoLine(strings.ToLower(line), "mat", lnc)
		}
	}

	// Sort & Print //
	//
	// Lines to be sorted are pushed to a slice.
	// Sort is configured on LineNum.
	// Uneven padding is accomplished by subtracting the length of the current WestSide fragment
	// from the length of the longest WestSide fragment (padCount)
	var linefragments LineFrags
	for k, _ := range fragMents {
		linefragments = append(linefragments, fragMents[k])
	}

	sort.Sort(linefragments)

	for i := 0; i < len(linefragments); i++ {
		padMe := padCount - linefragments[i].WChars
		spaces := strings.Repeat(" ", padMe)
		fmt.Printf("%s%s\n", spaces, linefragments[i].Data)
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
