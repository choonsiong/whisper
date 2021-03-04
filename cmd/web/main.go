package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

// Define an application struct to hold the application-wide dependencies.
type application struct {
	errorLog *log.Logger
	infoLog *log.Logger
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	// Create a logger for writing information messages.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Initialize a new instance of application containing the dependencies.
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	mux := http.NewServeMux()
	//mux.HandleFunc("/", home)
	//mux.HandleFunc("/whisper", showWhisper)
	//mux.HandleFunc("/whisper/create", createWhisper)
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

	// Initialize a new http server struct so that we can use our custom errorLog.
	srv := 	&http.Server{
		Addr: *addr,
		ErrorLog: errorLog,
		Handler: mux,
	}

	infoLog.Printf("Starting server on %s", *addr)
	//err := http.ListenAndServe(*addr, mux)
	err := srv.ListenAndServe()

	// As a rule of thumb, you should avoid using the Panic() and Fatal() variations outside
	// of main(), it's good practice to return errors instead, and only panic or exit directly
	// from main().
	errorLog.Fatal(err)
}
