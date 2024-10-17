package main

import (
	"log"
	"regexp"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/syndtr/goleveldb/leveldb"
)

var emailRegexp = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

type registration struct {
	Email string `json:"email"`
}

// Waitlist registration route
func routeRegister(app *fiber.App, p string, db *leveldb.DB, count *atomic.Uint64, captcha fiber.Handler) {
	app.Post(p, rateLimiter(), func(c fiber.Ctx) error {
		t := time.Now().UTC()
		req := new(registration)

		// Parse and validate the request
		if err := c.Bind().Body(req); err != nil {
			if verbose {
				log.Printf("invalid json from %q: %v", c.IP(), err)
			}
			return c.Status(fiber.StatusBadRequest).JSON(newMessage(false, "Invalid request"))
		}

		// Validate email format
		if !emailRegexp.MatchString(req.Email) {
			if verbose {
				log.Printf("invalid email from %q: %s", c.IP(), req.Email)
			}
			return c.Status(fiber.StatusBadRequest).JSON(newMessage(false, "Invalid email address"))
		}

		// Check if email is already registered
		if ok, err := db.Has([]byte(req.Email), nil); err != nil {
			log.Printf("cannot check for existence of email %q: %v", req.Email, err)
			return c.Status(fiber.StatusInternalServerError).JSON(newMessage(false, "Cannot check registration status"))
		} else if ok {
			if verbose {
				log.Printf("duplicate email from %q: %s", c.IP(), req.Email)
			}
			return c.Status(fiber.StatusConflict).JSON(newMessage(false, "Email already registered"))
		}

		// represent current time
		var bt []byte
		if b, err := t.MarshalBinary(); err != nil {
			log.Printf("cannot encode current time: %v", err)
			bt = []byte(strconv.Itoa(int(t.Unix())))
		} else {
			bt = b
		}

		// Save the email to the waitlist
		if err := db.Put([]byte(req.Email), bt, nil); err != nil {
			log.Printf("cannot register email %q: %v", req.Email, err)
			return c.Status(fiber.StatusInternalServerError).JSON(newMessage(false, "Cannot register email"))
		}

		log.Printf("registered email %q", req.Email)
		count.Add(1)
		return c.JSON(newMessage(true, "Email registered successfully"))
	}, captcha)
}
