package main

import (
	"testing"
)

func TestMesostic_ParseSpine(t *testing.T) {
	mesostic := Mesostic{
		Title: "music has the rights to children",
	}

	mesostic.ParseSpine()

	if mesostic.Spine[0] != "m" {
		t.Errorf("Spine does not start with Title character")
	}
	if mesostic.Spine[5] != "h" {
		t.Errorf("Spine does not skip the space after the first word")
	}
}

func TestMesostic_ParseLine(t *testing.T) {
	mesostic := Mesostic{
		Title:     "wander",
		SpineIdx:  0,
		MLines:    make([]string, 0),
		MLinesIdx: 0,
	}
	mesostic.ParseSpine()

	if mesostic.FormatLine("the quick brown fox") {
		want := "the quick broWn fox"
		if mesostic.MLines[mesostic.MLinesIdx] != want {
			t.Errorf("MLines[%d] != %s", mesostic.MLinesIdx, want)
		}
		if mesostic.SpineIdx == 0 {
			t.Errorf("Expected SpineIdx to be non-zero, got: %d", mesostic.SpineIdx)
		}
	}
}
