package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var listen string
var postgres string

func main() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	listen = os.Getenv("LISTEN")
	postgres = os.Getenv("POSTGRES")

	// Open postgres
	log.Printf("connecting to postgres %s", postgres)
	db, err := sql.Open("postgres", postgres)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Handler
	h := NewHandler(db)
	// Routes
	router := NewRouter(h)

	log.Printf("starting server on %s", listen)
	log.Fatal(http.ListenAndServe(listen, router))
}
