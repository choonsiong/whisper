package main

import "net/http"

//func (app *application) routes() *http.ServeMux {
func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/whisper", app.showWhisper)
	mux.HandleFunc("/whisper/create", app.createWhisper)

	// Note that the path given to the http.Dir() is relative to the
	// project directory root.
	fileServer := http.FileServer(http.Dir("./ui/static"))

	// Use the mux.Handle() function to register the file server as the handler for
	// all URL paths that start with "/static/". For matching paths, we strip the
	// "/static" prefix before the request reaches the file server.
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	//return mux

	// Pass the servemux as the 'next' parameter to the secureHeaders middleware.
	// Because secureHeaders is just a function, and the function returns a
	// http.Handler we don't need to do anything else.
	//return secureHeaders(mux)

	// Wrap the existing chain with the logRequest middleware.
	return app.logRequest(secureHeaders(mux))
}
