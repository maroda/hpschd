package main

import (
	"fmt"
	"strings"
)

// LineFrag ::: Data model describing a processed LineFragment.
type LineFrag struct {
	Index   int  // Cardinality, starting with 0
	SSChar  byte // The single byte of SpineString that appears in this fragment.
	BigEnd  bool // Where in the fragment the SSChar appears: 'BigEnd=false' means the SSChar appears on the far left end of the LineFragment.
	LineNum int  // The assigned line number for this fragment. Each will appear exactly twice: one with BigEnd=true, one with BigEnd=false.
}

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
var lineFrag = make(map[string]LineFrag)
var lineNew = make(map[string]int)

// if this ends up being used, ln2cap() should take both the string and the character to replace with its capital
// this is a lot more processing though, it might be better to do it inline with the ToUpper(SSChar)
func ln2cap(s string) string {
	spineC := "q"
	spineCU := strings.ToUpper(spineC)
	re := strings.NewReplacer(spineC, spineCU)
	res := re.Replace(s)
	return res
}

// pstack takes a string and finds a 'spineC' character.
func pstack(s string) string {
	line := []string{s}
	fmt.Printf("Finding in: %q: ", line)

	var stack []string
	spineC := "q"

	for i := 0; i < len(s); i++ {
		appC := string(s[i])
		if appC == spineC {
			fmt.Println("found")
			appC = strings.ToUpper(appC)
			stack = append(stack, appC)
			// copy(stack, ???) <<< this should be copied to a map of slices, which are then concatenated together
			break
		}
		stack = append(stack, appC)
	}
	result := strings.Join(stack, "")
	return result
}

func main() {
	cmLN := "craquemattic"
	cmPS := pstack(cmLN)
	fmt.Println(cmPS)

	// this also capitalizes the letter in the string,
	// which would work on any size string
	// cmUL := ln2cap(cmLN)
	// fmt.Println(cmUL)

	/*

		what i might need here is a recursive function
		that steps through letters, comparing them against each other

		or, a function that gradually increases the slice ahead of it until it gets to the next anchor letter.

			String to array
			Begin slice1 at BOL, increasing len until the contents equal Anchor1 (case agnostic)
		Anchor1 = c
			ToUpper Anchor1
		Anchor1 = C
			Assign slice1 to new array LineAnchor1L
			Begin slice2 after Anchor1, increasing len until the contents equal Anchor2
		Anchor2 = r
			ToUpper Anchor2
		Anchor2 = R
			Assign slice2 to new array LineAnchor2L

			and so on

		Anchor3 = a
		Anchor4 = q
		Anchor5 = u
		Anchor6 = e

	*/
}
