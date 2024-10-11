package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	if err := InitDB(conf[dbPath]); err != nil {
		log.Fatalf("Failed to initialize LevelDB: %v", err)
	}

	defer CloseDB()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)

	if err := serve(sig); err != nil {
		log.Printf("cannot serve: %v", err)
	}

	log.Println("application exit")
}
