package main

import (
	"errors"
	"fmt"
	"github.com/choonsiong/whisper/pkg/models"
	"html/template"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		//http.NotFound(w, r)
		app.notFound(w)
		return
	}

	s, err := app.whispers.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	//for _, whisper := range s {
	//	fmt.Fprintf(w, "%v\n", whisper)
	//}

	data := &templateData{
		Whispers: s,
	}

	files := []string{
		"./ui/html/home.page.tmpl", // must be the first
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		//log.Println(err.Error())
		//app.errorLog.Println(err.Error())
		//http.Error(w, "Internal Server Error", 500)
		app.serverError(w, err)
		return
	}

	// Write the template content as the response body.
	// Pass in the templateData struct when executing the template.
	err = ts.Execute(w, data)
	if err != nil {
		//log.Println(err.Error())
		//app.errorLog.Println(err.Error())
		//http.Error(w, "Internal Server Error", 500)
		app.serverError(w, err)
	}

	//w.Write([]byte("Hello"))
}

func (app *application) showWhisper(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		//http.NotFound(w, r)
		app.notFound(w)
		return
	}

	s, err := app.whispers.Get(id)

	if err != nil {
		// If no matching record is found, return a 404 Not Found response.
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}

		return
	}

	// Create an instance of a templateData struct holding the whisper data.
	data := &templateData{Whisper: s}

	files := []string{
		"./ui/html/show.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)

	if err != nil {
		app.serverError(w, err)
		return
	}

	//err = ts.Execute(w, s)

	// Pass in the templateData struct when executing the template.
	err = ts.Execute(w, data)

	if err != nil {
		app.serverError(w, err)
	}

	if app.debug {
		app.infoLog.Printf("showWhisper: id = %d", id)
	}

	//w.Write([]byte("Show snippet"))
	//fmt.Fprintf(w, "Show whisper with ID %d...", id)
	//fmt.Fprintf(w, "%v", s)
}

func (app *application) createWhisper(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost) // Must call this before below methods, else no effect

		//w.WriteHeader(405) // WriteHeader() can only call once
		//w.Write([]byte("Method Not Allowed"))

		// Using http.Error() is more common than call the WriteHeader() and Write() above
		//http.Error(w, "Method Not Allowed", 405)
		app.clientError(w, http.StatusMethodNotAllowed)

		return
	}

	// Create dummy data for testing.
	title := "0 snail"
	content := "0 snail\nClimb Mount Fiji,\nBut slowly, slowly!\n\n- Kobayashi Issa"
	expires := "7"

	id, err := app.whispers.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	//w.Write([]byte("Create a new whisper..."))
	http.Redirect(w, r, fmt.Sprintf("/whisper?id=%d", id), http.StatusSeeOther)
}