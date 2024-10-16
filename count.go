package main

import (
	"log"
	"strconv"
	"sync/atomic"

	"github.com/gofiber/fiber/v3"
	"github.com/syndtr/goleveldb/leveldb"
)

func routeCount(app *fiber.App, p string, db *leveldb.DB, count *atomic.Uint64) error {
	var c uint64
	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		c++
	}
	iter.Release()
	if err := iter.Error(); err != nil {
		return err
	}
	count.Store(c)

	log.Printf("registration count is %d on startup", c)

	app.Get(p, func(c fiber.Ctx) error {
		c.Set("Content-Type", "text/plain")
		return c.SendString(strconv.FormatUint(count.Load(), 10))
	})

	return nil
}
