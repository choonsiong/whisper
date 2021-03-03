package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Write([]byte("Hello"))
}

func showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	//w.Write([]byte("Show snippet"))
	fmt.Fprintf(w, "Show snippet with ID %d...", id)
}

func createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost) // Must call this before below methods, else no effect

		//w.WriteHeader(405) // WriteHeader() can only call once
		//w.Write([]byte("Method Not Allowed"))

		// Using http.Error() is more common than call the WriteHeader() and Write() above
		http.Error(w, "Method Not Allowed", 405)

		return
	}

	w.Write([]byte("Create snippet"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}