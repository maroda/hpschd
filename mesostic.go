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
	Count   int    // Character count, not including SpineChar.
	SSChar  string // The single byte of SpineString that borders this fragment.
	BigEnd  bool   // Where against the fragment SSChar appears: 'BigEnd=false' means the SSChar appears on the far left end of the LineFragment.
	LineNum int    // The assigned line number for this fragment. Each will appear exactly twice: one with BigEnd=true, one with BigEnd=false.
	Data    string // Equivalent to the Key, which may be changed to something else (hash value)
}
type LineFrags []LineFrag

/* Hash table of LineFragments

Key = The LineFragment itself
Value = LineFrag struct

Example:
spineString = craque
lineFeed = craquemattic
lineFrag["cra"] == {0, "q", 1, false}
lineFrag["uemattic"] == {1, "q", 1, true}

Notice that this *removes* the SSChar from the line so the two strings can be concatenated on either side of a ToUpper(SSChar).

So to construct a single line:

lineFeed == "craquemattic"
lineCat == get line number; ( lineFrag key BigEnd=false ("cra") + lineFrag ToUpper(SSChar) ("Q") + lineFrag key BigEnd=true ("uemattic") )
lineEntry == lineCat <<< "craQuemattic"
lineNew[lineEntry] == mesoLineNum <<< "craQuemattic" == 1

*/

var fragCount int
var fileLines = make(map[string]int)      // The file broken into strings with line numbers
var fragMents = make(map[string]LineFrag) // Conglomerate line fragments, the line as the key
var newLines = make(map[string]int)       // Concatenated fragments that make up the new lines

// pstack takes a string and finds a 'spineC' (SC) character.
// nothing is returned, instead the fragMents hash table is updated
func pstack(s string, c int) {
	// DEBUG ::: line := []string{s}
	// DEBUG ::: fmt.Printf("Finding in: %q: ", line)
	var stack []string
	var charCount int
	var SC string

	SS := "mat"

	// Step through each character in SpineString (SS)
	for h := 0; h < len(SS); h++ {
		SC := string(SS[h])

		// !!! current problem with this is that it restarts the first string every time.
		for i := 0; i < len(s); i++ {
			appC := string(s[i])
			if appC == SC {
				appC = strings.ToUpper(appC)
				stack = append(stack, appC)
				// continue // <<< keeps processing the line, but at this point, we want it in a new fragment
				break // <<< stops processing the line
			} else {
				stack = append(stack, appC)
				charCount++
			}
		}
		// break exits here, so we move on to the next SC, but we lost our place in the string and have to start over!
	}

	fragment := strings.Join(stack, "")
	fragCount++
	fragkey := shakey(fragment)
	fragMents[fragkey] = LineFrag{Index: c, Count: charCount, SSChar: SC, BigEnd: false, LineNum: fragCount, Data: fragment}

	// record the longest fragment length for padding
	if charCount > padCount {
		padCount = charCount
	}
}

var padCount int

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
		// run pstack - which needs a new name - to process each line
		for _, line := range strings.Split(string(data), "\n") {
			lnc++
			pstack(strings.ToLower(line), lnc)
		}
	}

	// string slice (LineFrags) of the LineFrag data structure filled with all entries to be sorted
	var linefragments LineFrags
	for k, _ := range fragMents {
		linefragments = append(linefragments, fragMents[k])
	}

	// sorted and printed with indention and padding
	sort.Sort(linefragments)
	for i := 0; i < len(linefragments); i++ {
		fmt.Printf("\t%*s\n", padCount, linefragments[i].Data)
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
