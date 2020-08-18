package main

import (
	"bytes"
	"fmt"
	"strings"
)

// this is an example of recursion, writing a comma between every 3 digits via concatenation
func comma(s string) string {
	n := len(s)
	if n <= 3 {
		return s
	}
	return comma(s[:n-3]) + "," + s[n-3:]
}

// intsToString is like fmt.Sprint(values) but adds commas.
// takes an array of ints as input
func intsToString(values []int) string {
	var buf bytes.Buffer
	buf.WriteByte('[')
	for i, v := range values {
		if i > 0 {
			buf.WriteString(", ")
		}
		fmt.Fprintf(&buf, "%d", v)
	}
	buf.WriteByte(']')
	return buf.String()
}

func ln2cap(s string) string {
	spineC := "q"
	spineCU := strings.ToUpper(spineC)
	re := strings.NewReplacer(spineC, spineCU)
	res := re.Replace(s)
	return res
}

func main() {
	// this passes a slice of ints
	// fmt.Println(intsToString([]int{4, 2, 1})) // array literal

	// this capitalizes the letter in the string,
	// which would work on any size string
	cmLN := "craquemattic"
	cmUL := ln2cap(cmLN)
	fmt.Println(cmUL)

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

	// used for the recursive comma function
	//for i := 1; i < len(os.Args); i++ {
	//	fmt.Printf("  %s\n", comma(os.Args[i]))
	//}
}
