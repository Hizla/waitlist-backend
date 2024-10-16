package main

import (
	"log"
	"regexp"

	"github.com/gofiber/fiber/v3"
	"github.com/syndtr/goleveldb/leveldb"
)

var emailRegexp = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

type registration struct {
	Email string `json:"email"`
}

// Waitlist registration route
func routeRegister(app *fiber.App, db *leveldb.DB, captcha fiber.Handler) {
	app.Post("/register", func(c fiber.Ctx) error {
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

		// Save the email to the waitlist
		if err := db.Put([]byte(req.Email), []byte{'x'}, nil); err != nil {
			log.Printf("cannot register email %q: %v", req.Email, err)
			return c.Status(fiber.StatusInternalServerError).JSON(newMessage(false, "Cannot register email"))
		}

		log.Printf("registered email %q", req.Email)
		return c.JSON(newMessage(true, "Email registered successfully"))
	}, captcha)
}
