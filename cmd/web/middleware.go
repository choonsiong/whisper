package main

import "net/http"

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