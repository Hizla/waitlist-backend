package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/contrib/hcaptcha"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/syndtr/goleveldb/leveldb"
)

func serve(sig chan os.Signal, db *leveldb.DB) error {
	app := fiber.New()

	// cors
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{conf[allowedURL]},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept"},
	}))

	// rate limiting
	app.Use(limiter.New(limiter.Config{
		Max:        5,
		Expiration: 7 * 24 * time.Hour, // 1 week expiration
		LimitReached: func(c fiber.Ctx) error {
			log.Printf("Rate limit exceeded for IP: %s", c.IP())
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"message": "Rate limit exceeded. Max 5 registrations per week.",
			})
		},
	}))

	var captcha fiber.Handler
	hCaptchaEnable := conf[hCaptchaSiteKey] != "unset" && conf[hCaptchaSecretKey] != "unset"

	if hCaptchaEnable {
		// create hCaptcha middleware if enabled
		captcha = hcaptcha.New(hcaptcha.Config{
			SecretKey: conf[hCaptchaSecretKey],
		})

		log.Printf("hCaptcha enabled with site key %q", conf[hCaptchaSiteKey])
	} else {
		// empty middleware if disabled
		captcha = func(c fiber.Ctx) error {
			return c.Next()
		}

		log.Printf("hCaptcha disabled because one or both of %q and %q are unset",
			confEnv[hCaptchaSiteKey][0], confEnv[hCaptchaSecretKey][0])
	}

	routeHCaptchaSiteKey(app, "/api", !hCaptchaEnable, conf[hCaptchaSiteKey])
	routeRegister(app, "/api/register", db, captcha)

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
