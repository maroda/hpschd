/*

	This requires more, but global vars get in the way a lot.

*/

package main

import (
	"fmt"
	"testing"
)

// TestIctusIncr
func TestIctusOnce(t *testing.T) {
	ss := "mat"
	for h := 0; h < len(ss); h++ {
		sca = append(sca, string(ss[h]))
	}

	if sca[ictus] != "m" {
		t.Errorf("%q\n", ss)
	}
	fmt.Println(ictus, sca[ictus])
	Ictus(1)
	if sca[ictus] != "a" {
		t.Errorf("%q\n", ss)
	}
	fmt.Println(ictus, sca[ictus])
	Ictus(1)
	if sca[ictus] != "t" {
		t.Errorf("%q\n", ss)
	}
	fmt.Println(ictus, sca[ictus])
}
