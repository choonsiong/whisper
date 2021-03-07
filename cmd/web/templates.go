package main

import (
	"github.com/choonsiong/whisper/pkg/models"
	"html/template"
	"path/filepath"
	"time"
)

// Define a templateData type to act as the holding structure for any dynamic
// data that we want to pass to our HTML templates.
type templateData struct {
	CurrentYear int
	Whisper *models.Whisper
	Whispers []*models.Whisper // Store the parsed templates in an in-memory cache.
}

// Create a humanDate function which returns a nicely formatted string
// representation of a time.Time object.
// Note: Custom template function (like humanDate() below) can accept as many parameters
// as they need to, but they must return one value only. The only exception to this is if
// you want to return an error as the second value.
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// Initialize a template.FuncMap object and store it in a global variable.
// This is essentially a string-keyed map which acts a a lookup between the names
// of our custom template functions and the functions themselves.
var functions = template.FuncMap{
	"humanDate": humanDate,
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

		// The template.FuncMap must be registered with the template set before you
		// call the ParseFiles() method. This means we have to use template.New() to
		// create an empty template set, use the Funcs() method to register the
		// template.FuncMap, and then parse the file as normal.
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)

		// Parse the page template file in to a template set.
		//ts, err := template.ParseFiles(page)

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