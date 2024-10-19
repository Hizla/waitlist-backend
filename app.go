package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync/atomic"

	"github.com/gofiber/contrib/hcaptcha"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/syndtr/goleveldb/leveldb"
)

func serve(sig chan os.Signal, db *leveldb.DB) error {
	app := fiber.New()

	// cors
	if conf[allowedURL] != "unset" {
		app.Use(cors.New(cors.Config{
			AllowOrigins: []string{conf[allowedURL]},
			AllowHeaders: []string{"Origin", "Content-Type", "Accept"},
		}))
	} else {
		log.Println("CORS disabled")
	}

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

	count := new(atomic.Uint64)

	routeHCaptchaSiteKey(app, "/api", !hCaptchaEnable, conf[hCaptchaSiteKey])
	routeRegister(app, "/api/register", db, count, captcha)

	if err := routeCount(app, "/api/count", db, count); err != nil {
		return err
	}

	// graceful shutdown
	go func() {
		<-sig
		log.Println("shutting down")
		if err := app.Shutdown(); err != nil {
			fmt.Printf("cannot shutdown: %v", err)
		}
	}()

	if conf[listen] == "unset" {
		return app.Listen(conf[listenAddr])
	} else {
		if l, err := net.Listen("unix", conf[listen]); err != nil {
			return err
		} else {
			if err = os.Chmod(conf[listen], 0777); err != nil {
				log.Printf("cannot change ownership of socket %q: %v", conf[listen], err)
			}
			return app.Listener(l)
		}
	}
}
