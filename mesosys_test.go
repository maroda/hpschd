package main

import (
	"testing"
)

func TestMesostic_ParseSpine(t *testing.T) {
	mesostic := Mesostic{
		Title: "music has the rights to children",
	}

	newspine := envVar("HPSCHD_SPINESTRING", "")
	mesostic.ParseSpine(newspine)

	if mesostic.Spine[0] != "m" {
		t.Errorf("Spine does not start with Title character")
	}
	if mesostic.Spine[5] != "h" {
		t.Errorf("Spine does not skip the space after the first word")
	}
}

func TestMesostic_BuildMeso(t *testing.T) {
	t.Run("Correct mesostic text returned for APOD source", func(t *testing.T) {
		title := "The Millennium that Defines Universe"
		ae := &DataAPOD{}
		meso := NewMesostic(title, testApodJSON, ae)

		gotae, ok := meso.SourceData.(*DataAPOD)
		if !ok {
			t.Errorf("Source data does not implement DataAPOD")
		}
		meso.MU.Lock()
		meso.SourceTxt = gotae.Explaination
		meso.MU.Unlock()

		got := meso.BuildMeso()
		assertStringContains(t, got, mesosticApod)
		t.Log(got)
	})

	t.Run("Correct mesostic text returned for DataAPI source", func(t *testing.T) {
		apiJSON = `{"text": "the quick brown; fox jumps over; the lazy dog", "spinestring": "craque"}`
		want := `
      the quiCk b
fox jumps oveR
        the lAzy dog`

		title := "craque"
		ae := &DataAPI{}
		meso := NewMesostic(title, apiJSON, ae)

		gotae, ok := meso.SourceData.(*DataAPI)
		if !ok {
			t.Errorf("Source data does not implement DataAPOD")
		}
		meso.MU.Lock()
		meso.SourceTxt = gotae.Text
		meso.MU.Unlock()

		got := meso.BuildMeso()
		if got != want {
			t.Errorf("Expected:\n%s\n\nGot:\n%s", want, got)
		}
	})
}
