package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	files := []string{
		"./ui/html/home.page.tmpl", // must be the first
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	// Write the template content as the response body.
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}

	//w.Write([]byte("Hello"))
}

func showWhisper(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	//w.Write([]byte("Show snippet"))
	fmt.Fprintf(w, "Show whisper with ID %d...", id)
}

func createWhisper(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost) // Must call this before below methods, else no effect

		//w.WriteHeader(405) // WriteHeader() can only call once
		//w.Write([]byte("Method Not Allowed"))

		// Using http.Error() is more common than call the WriteHeader() and Write() above
		http.Error(w, "Method Not Allowed", 405)

		return
	}

	w.Write([]byte("Create a new whisper..."))
}