package main

import (
	"database/sql"
	"flag"
	_ "github.com/lib/pq"
	"github.com/tbarisic/letsgo-snippetbox/internal/models"
	"html/template"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func main() {

	address := flag.String("address", ":1337", "HTTP network address")

	dsn := flag.String("dsn", "postgres://snippetbox:test@localhost:5433/snippetbox?sslmode=disable", "Postgresql data source name")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)

	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	cache, err := newTemplateCache()

	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &models.SnippetModel{DB: db, LOG: infoLog},
		templateCache: cache,
	}

	server := &http.Server{
		Addr:     *address,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", *address)
	err = server.ListenAndServe()
	errorLog.Fatal(err)
}
