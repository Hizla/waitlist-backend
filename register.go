package main

import (
	"log"
	"regexp"

	"github.com/gofiber/fiber/v2"
)

var emailRegexp = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

type registration struct {
	Email string `json:"email"`
}

// Waitlist registration route
func routeRegister(app *fiber.App) {
	app.Post("/register", func(c *fiber.Ctx) error {
		req := new(registration)

		if err := c.BodyParser(req); err != nil {
			log.Printf("Invalid request body: %v", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid request",
			})
		}

		// Validate email format
		if !emailRegexp.MatchString(req.Email) {
			log.Printf("Invalid email format: %s", req.Email)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid email format",
			})
		}

		// Check if email is already registered
		registered, err := IsEmailRegistered(req.Email)
		if err != nil {
			log.Printf("Error checking email registration: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Error checking registration status",
			})
		}
		if registered {
			log.Printf("Email already registered: %s", req.Email)
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"message": "Email already registered",
			})
		}

		// Save the email to the waitlist
		if err := SaveEmail(req.Email); err != nil {
			log.Printf("Error registering email: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Error registering email",
			})
		}

		log.Printf("Email registered successfully: %s", req.Email)
		return c.JSON(fiber.Map{
			"message": "Email registered successfully",
		})
	})
}
