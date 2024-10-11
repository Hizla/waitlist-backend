package main

import (
	"log"
	"regexp"

	"github.com/gofiber/fiber/v2"
	"github.com/syndtr/goleveldb/leveldb"
)

var emailRegexp = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

type registration struct {
	Email string `json:"email"`
}

// Waitlist registration route
func routeRegister(app *fiber.App, db *leveldb.DB) {
	app.Post("/register", func(c *fiber.Ctx) error {
		req := new(registration)

		if err := c.BodyParser(req); err != nil {
			if verbose {
				log.Printf("invalid request from %q: %v", c.IP(), err)
			}
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid request",
			})
		}

		// Validate email format
		if !emailRegexp.MatchString(req.Email) {
			if verbose {
				log.Printf("invalid email from %q: %s", c.IP(), req.Email)
			}
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid email format",
			})
		}

		// Check if email is already registered
		if ok, err := db.Has([]byte(req.Email), nil); err != nil {
			log.Printf("cannot check for existence of email %q: %v", req.Email, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Error checking registration status",
			})
		} else if ok {
			if verbose {
				log.Printf("duplicate email from %q: %s", c.IP(), req.Email)
			}
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"message": "Email already registered",
			})
		}

		// Save the email to the waitlist
		if err := db.Put([]byte(req.Email), []byte{'x'}, nil); err != nil {
			log.Printf("cannot register email %q: %v", req.Email, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Error registering email",
			})
		}

		log.Printf("registered email %q", req.Email)
		return c.JSON(fiber.Map{
			"message": "Email registered successfully",
		})
	})
}
