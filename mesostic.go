package main

import (
	"fmt"
	"strings"
)

// define a global map
// whose keys are the new slices
// and values are the ordering/placements

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
	cmLN := "the quick brown fox jumped over the lazy dog"
	//cmLN := "craquemattic"

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
