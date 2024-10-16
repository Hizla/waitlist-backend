package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/syndtr/goleveldb/leveldb"
)

func main() {
	var db *leveldb.DB

	if d, err := leveldb.OpenFile(conf[dbPath], nil); err != nil {
		log.Fatalf("cannot open database %q: %v", conf[dbPath], err)
	} else {
		db = d
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("cannot close database %q: %v", conf[dbPath], err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	if err := serve(sig, db); err != nil {
		log.Printf("cannot serve: %v", err)
	}

	log.Println("application exit")
}
