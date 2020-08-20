package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
)

// LineFrag ::: Data model describing a processed LineFragment.
type LineFrag struct {
	Index   int    // Cardinality, starting with 0
	SSChar  string // The single byte of SpineString that appears in this fragment.
	BigEnd  bool   // Where in the fragment the SSChar appears: 'BigEnd=false' means the SSChar appears on the far left end of the LineFragment.
	LineNum int    // The assigned line number for this fragment. Each will appear exactly twice: one with BigEnd=true, one with BigEnd=false.
	Data    string // Equivalent to the Key
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

// pstack takes a string and finds a 'spineC' character.
func pstack(s string, c int) string {
	// DEBUG ::: line := []string{s}
	// DEBUG ::: fmt.Printf("Finding in: %q: ", line)

	var stack []string
	SC := "m"

	for i := 0; i < len(s); i++ {
		appC := string(s[i])
		if appC == SC {
			// DEBUG ::: fmt.Println("found")
			appC = strings.ToUpper(appC)
			stack = append(stack, appC)
			// copy(stack, ???) <<< this should be copied to a map of slices, which are then concatenated together
			break
		}
		stack = append(stack, appC)
	}
	fragment := strings.Join(stack, "")
	fragCount++
	// this doesn't currently remove the SSChar from the fragment!!!
	fragMents[fragment] = LineFrag{Index: c, SSChar: SC, BigEnd: false, LineNum: fragCount, Data: fragment}
	// DEBUG ::: fmt.Println(fragMents[fragment])
	return fragment
}

func main() {
	var lnc int

	// process the given file, silent if no file given
	for _, origtxt := range os.Args[1:] {
		data, err := ioutil.ReadFile(origtxt)
		if err != nil {
			log.Fatal(err)
			break
		}
		for _, line := range strings.Split(string(data), "\n") {
			// DEBUG ::: fmt.Println("^^^ ", line)
			// this might not be the right place to pass the SpineChar
			// instead, pass the line number
			lnc++
			pstack(strings.ToLower(line), lnc)
			// PS := pstack(strings.ToLower(line), lnc)
			// fmt.Println(PS)
		}
	}

	var linefragments LineFrags

	fragKeys := make([]string, 0, len(fragMents))

	for k, _ := range fragMents {
		// fmt.Println(fragMents[k].LineNum)
		linefragments = append(linefragments, fragMents[k])
		// fmt.Println(k)
		fragKeys = append(fragKeys, k)
	}

	// now print the mesostic with lines in sorted order
	// DEBUG ::: fmt.Println(linefragments)
	// DEBUG ::: fmt.Println(sort.IsSorted(linefragments))
	sort.Sort(linefragments)
	for i := 0; i < len(linefragments); i++ {
		fmt.Println(linefragments[i].LineNum, linefragments[i].Data)
	}
}

// Configure linefragments Sort
func (ls LineFrags) Len() int {
	return len(ls)
}
func (ls LineFrags) Swap(i, j int) {
	ls[i], ls[j] = ls[j], ls[i]
}
func (ls LineFrags) Less(i, j int) bool {
	return ls[i].LineNum < ls[j].LineNum
}
