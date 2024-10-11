package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/syndtr/goleveldb/leveldb"
)

func serve(sig chan os.Signal, db *leveldb.DB) error {
	app := fiber.New()

	// cors
	app.Use(cors.New(cors.Config{
		AllowOrigins: conf[allowedOrigins],
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// rate limiting
	app.Use(limiter.New(limiter.Config{
		Max:        5,
		Expiration: 7 * 24 * time.Hour, // 1 week expiration
		LimitReached: func(c *fiber.Ctx) error {
			log.Printf("Rate limit exceeded for IP: %s", c.IP())
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"message": "Rate limit exceeded. Max 5 registrations per week.",
			})
		},
	}))

	// /register
	routeRegister(app, db)

	// Graceful shutdown
	app.Use(func(c *fiber.Ctx) error {
		return c.Next()
	})

	// graceful shutdown
	go func() {
		<-sig
		log.Println("shutting down")
		if err := app.Shutdown(); err != nil {
			fmt.Printf("cannot shutdown: %v", err)
		}
	}()

	return app.Listen(conf[listenAddr])
}
