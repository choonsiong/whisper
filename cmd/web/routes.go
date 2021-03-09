package main

import (
	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
	"net/http"
)

// Note: Update the routes() to use bmizerany/pat package to help implement
// RESTful routes.
func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	mux := pat.New()

	// Pat doesn't allow us to register handler functions directly, so we need to
	// convert them using the http.HandlerFunc() adapter.
	// Pat matches patterns in the order that they are registered.
	mux.Get("/", http.HandlerFunc(app.home))
	mux.Get("/whisper/create", http.HandlerFunc(app.createWhisperForm))
	mux.Post("/whisper/create", http.HandlerFunc(app.createWhisper))

	// The named capture (:id) acts like a wildcard, whereas the rest of the pattern
	// matches literally. Pat will add the contents of the named capture to the URL
	// query string at runtime behind the scenes.
	mux.Get("/whisper/:id", http.HandlerFunc(app.showWhisper))

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}

// Note: We update the routes to use the justinas/alice package to help us manage
// our middleware/handler chains.

//func (app *application) routes() *http.ServeMux {
//func (app *application) routes() http.Handler {
//	// Create a middleware chain containing our 'standard' middleware which will
//	// be used for every request our application receives.
//	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
//
//	mux := http.NewServeMux()
//
//	mux.HandleFunc("/", app.home)
//	mux.HandleFunc("/whisper", app.showWhisper)
//	mux.HandleFunc("/whisper/create", app.createWhisper)
//
//	// Note that the path given to the http.Dir() is relative to the
//	// project directory root.
//	fileServer := http.FileServer(http.Dir("./ui/static"))
//
//	// Use the mux.Handle() function to register the file server as the handler for
//	// all URL paths that start with "/static/". For matching paths, we strip the
//	// "/static" prefix before the request reaches the file server.
//	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
//
//	//return mux
//
//	// Pass the servemux as the 'next' parameter to the secureHeaders middleware.
//	// Because secureHeaders is just a function, and the function returns a
//	// http.Handler we don't need to do anything else.
//	//return secureHeaders(mux)
//
//	// Wrap the existing chain with the logRequest middleware.
//	//return app.logRequest(secureHeaders(mux))
//
//	// Wrap the existing chain with the recoverPanic middleware.
//	//return app.recoverPanic(app.logRequest(secureHeaders(mux)))
//
//	// Return the 'standard' middleware chain followed by the servemux.
//	return standardMiddleware.Then(mux)
//}
