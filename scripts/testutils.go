package main

import (
	"database/sql"
	"flag"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func init() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	flag.Parse()

	exit := 1
	if flag.Arg(0) == "isrunning" {
		exit = isrunning()
	} else if flag.Arg(0) == "wait" {
		exit = wait()
	} else {
		log.Print("No command specified")
	}
	os.Exit(exit)
}

func isrunning() int {
	db, err := sql.Open("postgres", os.Getenv("POSTGRES"))
	if err != nil {
		return 1
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		return 1
	}
	return 0
}

func wait() int {
	for i := 0; i < 60; i += 1 {
		timer := time.NewTimer(1 * time.Second)
		if isrunning() == 0 {
			return 0
		}
		<-timer.C
	}
	log.Print("Timed out trying to connect to postgres")
	return 1
}
