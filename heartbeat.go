package main

import "github.com/rs/zerolog/log"

// nasaNewREAD ::: Consume the current filename for the current NASA APOD Mesostic.
func nasaNewREAD() string {
	_, _, fu := Envelope()
	m := <-nasaNewMESO
	log.Info().Str("fu", fu).Msg("Filename from nasaNewMESO consumed")
	return m

	// this select was in use in homepage()
	// might be pertinent here too
	/*
		select {
		case mesoFile := <-nasaNewMESO:
			formatMeso.Title = mesoFile
			formatMeso.Mesostic = readMesoFile(&mesoFile)
			log.Info().Str("fu", fu).Msg("new mesostic received")
		default:
			formatMeso.Title = "HPSCHD"
			log.Info().Str("fu", fu).Msg("NO ACTION")
		}
	*/
}
