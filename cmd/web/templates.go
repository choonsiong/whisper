package main

import "github.com/choonsiong/whisper/pkg/models"

// Define a templateData type to act as the holding structure for any dynamic
// data that we want to pass to our HTML templates.
type templateData struct {
	Whisper *models.Whisper
}
