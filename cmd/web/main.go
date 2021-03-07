package main

import (
	"database/sql"
	"flag"
	"github.com/choonsiong/whisper/pkg/models/mysql"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// Define an application struct to hold the application-wide dependencies.
type application struct {
	debug bool
	errorLog *log.Logger
	infoLog *log.Logger
	whispers *mysql.WhisperModel
	templateCache map[string]*template.Template
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	debug := flag.Bool("debug", false,"Turn on debug mode")
	dsn := flag.String("dsn", "whisperadmin:password@/whisper?parseTime=true", "MySQL data source name")

	flag.Parse()

	// Create a logger for writing information messages.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	// Initialize a new template cache...
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	// Initialize a new instance of application containing the dependencies.
	app := &application{
		debug: *debug,
		errorLog: errorLog,
		infoLog:  infoLog,
		whispers: &mysql.WhisperModel{DB: db},
		templateCache: templateCache,
	}

	// Initialize a new http server struct so that we can use our custom errorLog.
	srv := 	&http.Server{
		Addr: *addr,
		ErrorLog: errorLog,
		Handler: app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	//err := http.ListenAndServe(*addr, mux)
	err = srv.ListenAndServe()

	// As a rule of thumb, you should avoid using the Panic() and Fatal() variations outside
	// of main(), it's good practice to return errors instead, and only panic or exit directly
	// from main().
	errorLog.Fatal(err)
}

// The openDB() wraps sql.Open() and returns a sql.DB connection pool for a given DSN.
func openDB(dsn string) (*sql.DB, error) {
	// The sql.Open() doesn't actually create any connections, all it does is initialize the
	// pool for future use. Actual connections to the database are established lazily, as
	// and when needed for the first time. So to verify that everything is set up correctly we
	// need to use the db.Ping() to create a connection and check for any errors.
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}