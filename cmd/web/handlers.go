package main

import (
	"errors"
	"fmt"
	"github.com/choonsiong/whisper/pkg/models"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// No longer needed, because pat matches "/" path exactly.
	//if r.URL.Path != "/" {
	//	//http.NotFound(w, r)
	//	app.notFound(w)
	//	return
	//}

	// Introduce panic for testing recoverPanic
	//panic("Ooops!")

	s, err := app.whispers.Latest()

	if err != nil {
		app.serverError(w, err)
		return
	}

	//for _, whisper := range s {
	//	fmt.Fprintf(w, "%v\n", whisper)
	//}

	// Use the new render helper instead of below...
	app.render(w, r, "home.page.tmpl", &templateData{
		Whispers: s,
	})

	if app.debug {
		app.debugLog.Printf("home: %v", s)
	}

	/*
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
	}*/

	//w.Write([]byte("Hello"))
}

func (app *application) showWhisper(w http.ResponseWriter, r *http.Request) {
	//id, err := strconv.Atoi(r.URL.Query().Get("id"))

	// Pat doesn't strip the colon from the named capture key, so we need to
	// get the value of ":id" from the query string instead of "id".
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))

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

	app.render(w, r, "show.page.tmpl", &templateData{
		Whisper: s,
	})

	// Create an instance of a templateData struct holding the whisper data.
	/*data := &templateData{Whisper: s}

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
	}*/

	if app.debug {
		app.debugLog.Printf("showWhisper: id = %d", id)
		app.debugLog.Printf("showWhisper: %v", s)
	}

	//w.Write([]byte("Show snippet"))
	//fmt.Fprintf(w, "Show whisper with ID %d...", id)
	//fmt.Fprintf(w, "%v", s)
}

func (app *application) createWhisperForm(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("Create a new whisper form"))
	app.render(w, r, "create.page.tmpl", nil)
}

func (app *application) createWhisper(w http.ResponseWriter, r *http.Request) {

	// Checking if the request method is a POST is now superfluous and can be
	// removed.
	//if r.Method != http.MethodPost {
	//	w.Header().Set("Allow", http.MethodPost) // Must call this before below methods, else no effect
	//
	//	//w.WriteHeader(405) // WriteHeader() can only call once
	//	//w.Write([]byte("Method Not Allowed"))
	//
	//	// Using http.Error() is more common than call the WriteHeader() and Write() above
	//	//http.Error(w, "Method Not Allowed", 405)
	//	app.clientError(w, http.StatusMethodNotAllowed)
	//
	//	return
	//}

	//title := "1 snail"
	//content := "1 snail\nOne two three\n\n- Foo Bar"
	//expires := "1000"

	// Call r.ParseForm() whichs add any data in POST request bodies to the
	// r.PostForm map. This also works in the same way for PUT and PATCH requests.
	err := r.ParseForm()

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Use the r.PostForm.Get() to retrieve the data fields from the r.PostForm map.
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires := r.PostForm.Get("expires")

	// Create a new whisper record in the database using the form data.
	id, err := app.whispers.Insert(title, content, expires)

	if err != nil {
		app.serverError(w, err)
		return
	}

	if app.debug {
		app.debugLog.Printf("createWhisper: {%s, %s, %s}\n", title, content, expires)
	}

	// Create dummy data for testing.
	//title := "0 snail"
	//content := "0 snail\nClimb Mount Fiji,\nBut slowly, slowly!\n\n- Kobayashi Issa"
	//expires := "7"

	//id, err := app.whispers.Insert(title, content, expires)
	//if err != nil {
	//	app.serverError(w, err)
	//	return
	//}

	//w.Write([]byte("Create a new whisper..."))
	//http.Redirect(w, r, fmt.Sprintf("/whisper?id=%d", id), http.StatusSeeOther)

	// Change the redirect to use the new semantic URL style of /whisper/:id
	http.Redirect(w, r, fmt.Sprintf("/whisper/%d", id), http.StatusSeeOther)
}