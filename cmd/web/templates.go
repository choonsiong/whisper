package main

import (
	"github.com/choonsiong/whisper/pkg/models"
	"html/template"
	"path/filepath"
)

// Define a templateData type to act as the holding structure for any dynamic
// data that we want to pass to our HTML templates.
type templateData struct {
	Whisper *models.Whisper
	Whispers []*models.Whisper // Store the parsed templates in an in-memory cache.
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	// Initialize a new map to act as the cache.
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	// Loop through the pages one-by-one.
	for _, page := range pages {
		// Extract the file name.
		name := filepath.Base(page)

		// Parse the page template file in to a template set.
		ts, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Add layout templates to the template set.
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		// Add the template set to the cache, using the name of the page as the key.
		cache[name] = ts
	}

	// Return the map.
	return cache, nil
}
