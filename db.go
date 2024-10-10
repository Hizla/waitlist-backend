package main

import (
	"github.com/syndtr/goleveldb/leveldb"
	"log"
)

var db *leveldb.DB

func InitDB(path string) error {
	var err error
	db, err = leveldb.OpenFile(path, nil)
	if err != nil {
		return err
	}
	return nil
}

func CloseDB() {
	if err := db.Close(); err != nil {
		log.Printf("Warning: Failed to close LevelDB: %v", err)
	}
}

func SaveEmail(email string) error {
	if err := db.Put([]byte(email), []byte("registered"), nil); err != nil {
		log.Printf("Error saving email to LevelDB: %v", err)
		return err
	}
	return nil
}

func IsEmailRegistered(email string) (bool, error) {
	_, err := db.Get([]byte(email), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return false, nil
		}
		log.Printf("Error checking email registration: %v", err)
		return false, err
	}
	return true, nil
}
