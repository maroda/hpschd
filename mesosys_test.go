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

func TestMesostic_BuildMeso(t *testing.T) {
	t.Run("Correct mesostic text returned for APOD source", func(t *testing.T) {
		// This is a real NASA APOD entry
		sourceJSON := `
{
    "date": "2000-01-01",
    "explanation": "Welcome to the millennial year at the threshold of millennium three.  During millennium two, humanity continually redefined its concept of \"Universe\": first as spheres centered on the Earth, in mid-millennium as the Solar System, a few centuries ago as the Galaxy, and within the last century as the matter emanating from the Big Bang.  During millennium three humanity may hope to discover alien life, to understand the geometry and composition of our present concept of Universe, and even to travel through this Universe.  Whatever our accomplishments, humanity will surely find adventure and discovery in the space above and beyond, and possibly define the surrounding Universe in ways and colors we cannot yet imagine by the threshold of millennium four.",
    "hdurl": "https://apod.nasa.gov/apod/image/0001/flammarion_halfcolor.gif",
    "media_type": "image",
    "service_version": "v1",
    "title": "The Millennium that Defines Universe",
    "url": "https://apod.nasa.gov/apod/image/0001/flammarion_halfcolor_big.gif"
}
`
		want := `
                 welcome To

             first as sphEr
                      in Mid-
             a few centurIes ago as the galaxy
          and within the Last century as the matter emanating from the big bang
                during miL
                   to undErstand th
                        aNd eve
whatever our accomplishmeNts
                    humanIty w
and possibly define the sUrro`

		title := "The Millennium that Defines Universe"
		ae := &DataAPOD{}
		meso := NewMesostic(title, sourceJSON, ae)

		gotae, ok := meso.SourceData.(*DataAPOD)
		if !ok {
			t.Errorf("Source data does not implement DataAPOD")
		}
		meso.MU.Lock()
		meso.SourceTxt = gotae.Explaination
		meso.MU.Unlock()

		got := meso.BuildMeso()
		if got != want {
			t.Errorf("Expected:\n%s\n\nGot:\n%s", want, got)
		}
	})

	t.Run("Correct mesostic text returned for DataAPI source", func(t *testing.T) {
		apiJSON := `{"text": "the quick brown; fox jumps over; the lazy dog", "spinestring": "craque"}`
		want := `
      the quiCk brown
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
