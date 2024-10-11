package main

import (
	"log"
	"regexp"
)

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

func main() {

	if err := InitDB(conf[dbPath]); err != nil {
		log.Fatalf("Failed to initialize LevelDB: %v", err)
	}

	defer CloseDB()

	if err := serve(); err != nil {
		log.Printf("cannot serve: %s", err)
	}

	log.Println("application exit")
}
