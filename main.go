package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
  "github.com/nats-io/go-nats"
	_ "github.com/lib/pq"
)

var listen string
var postgres string

func init() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	listen = os.Getenv("LISTEN")

	// Database
	db := connect()
  // NATs
  nc := connectNats()
	// Handler
	h := NewHandler(db, nc)
	// Routes
	router := NewRouter(h)

	log.Printf("starting server on %s", listen)
	log.Fatal(http.ListenAndServe(listen, router))
}

func connect() *sql.DB {
	postgres = os.Getenv("POSTGRES")

	// Open postgres
	log.Printf("connecting to postgres %s", postgres)
	db, err := sql.Open("postgres", postgres)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func connectNats() *nats.Conn {
  natsHost := os.Getenv("NATS")
  var nc *nats.Conn
  var err error
  if natsHost != "" {
    log.Printf("connecting to nats %s", natsHost)
    nc, err = nats.Connect(natsHost)
    if err != nil {
      log.Fatal(err)
    }
  }
  return nc
}
