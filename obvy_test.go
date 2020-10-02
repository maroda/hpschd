/*

	Observability Tests

*/

package main

import (
	"fmt"
	"strings"
	"testing"
)

// TestTEnvelope ::: Test that tracing is working by checking myself.
func TestTEnvelope(t *testing.T) {
	fmt.Printf("\n\t::: Test Target Envelope() :::\n")

	// Every function calls Envelope() as needed for tracing.
	fi, li, fu := Envelope()     // REQUIRED: myline == THIS LINE NUMBER
	myline := 20                 // Test Value for Line Number
	me := "hpschd.TestTEnvelope" // Test Value for Function Name
	myfile := "obvy_test.go"     // Test Value for File Name

	fiS := strings.Split(fi, "/")
	if fiS[len(fiS)-1] != myfile {
		t.Errorf("The value '%s' is not the filename '%s' of the caller.\n", fi, myfile)
	}

	if li != myline {
		t.Errorf("The value '%d' is not the line number of the caller at '%d'.\n", li, myline)
	}

	fuS := strings.Split(fu, "/")
	if fuS[len(fuS)-1] != me {
		t.Errorf("The value '%s' is not the name of the caller '%s'.\n", fu, me)
	}
}
