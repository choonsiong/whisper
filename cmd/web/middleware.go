package main

import (
	"fmt"
	"net/http"
)

// This 'middleware' will act on every request that is received, so it need to be
// executed before a request hits servemux.
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Implements additional security measures to help prevent XSS and Clickjacking attack.
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")

		// Code above will execute on the way down the chain
		next.ServeHTTP(w, r)
		// Code below will execute on the way back up the chain
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

// We don't want panic to bring down the whole application, so any panic is isolated to the
// goroutine serving the active HTTP request. Specifically, following a panic our server will
// log a stack trace to the server log, unwind the stack for the affected goroutine (calling any
// deferred functions along the way) and close the underlying HTTP connection. But it won't terminate
// the application, so importantly, any panic in our handlers won't bring down the server.
// Note: It's important to realize that this middleware will only recover panics that happen in the same
// goroutine that executed the recoverPanic() middleware. If, for example, you have a handler which
// spins up another goroutine (e.g. to do some background processing), then any panics that happen in
// the second goroutine will not be recovered - not by the recoverPanic() and not by the panic recovery
// built into Go HTTP server. They will cause your application to exit and bring down the server.
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a defer function (which will always be run in the event of a
		// panic as Go unwinds the stack).
		defer func() {
			// Use the builtin recover function to check if there has been a
			// paic or not. If there has...
			if err := recover(); err != nil {
				// Set a "Connection: close" header on the response so that Go's HTTP server
				// automatically close the current connection after a response has been
				// sent. It also informs the user that the connection will be closed.
				w.Header().Set("Connection", "close")
				// Return a 500 Internal Server Error response.
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}