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

var ss string = "craque"  // SpineString
var sca []string          // Characters in SpineString
var ictus int             // SpineString character address
var nexus int = ictus + 1 // Next SpineString character address

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
				because that will appear itself on the next line, and cannot have one preceeding it.
			*/
			if char != sca[nexus] {
				estack = append(estack, char)
			} else {
				break CharLoop // We're done.
			}

			// 100% Mesostic ::: No occurance of the current Spine String Character in back OR in front of it.
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

// api version that accepts text blocks
//		this will help break it down into smaller funcs
//		don't even care if the json sent over is any good
// t = SourceText
// s = SpineString
// func Mesostic(t string, s string) string {
func Mesostic(t string, s string, o chan<- string) {
	var lnc int

	ss = s
	for h := 0; h < len(ss); h++ {
		sca = append(sca, string(ss[h]))
	}

	// run mesoLine - which needs a new name - to process each line
	for _, line := range strings.Split(string(t), "\n") {
		lnc++
		mesoLine(strings.ToLower(line), lnc)
	}

	var linefragments LineFrags
	for k := range fragMents {
		linefragments = append(linefragments, fragMents[k])
	}

	sort.Sort(linefragments)

	for i := 0; i < len(linefragments); i++ {
		padMe := padCount - linefragments[i].WChars
		spaces := strings.Repeat(" ", padMe)
		// fmt.Printf("%s%s\n", spaces, linefragments[i].Data)
		//
		// ooooh this almost works! it is properly being sent back...
		// but i'm gonna have to figure out how to send back multiple lines
		// and send back a formatted string... Sprint won't work, but
		// the channel type seems wrong for Sprintf...
		// it correctly sends the constructed Mesostic line back,
		// but the intentional whitespace is being chopped
		o <- fmt.Sprint(spaces, linefragments[i].Data)
	}
	close(o)
}

// standalone version that reads text files
func mesomain() {
	var lnc int // line counts for the Index

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
