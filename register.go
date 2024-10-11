package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

// Waitlist registration route
func routeRegister(app *fiber.App) {
	app.Post("/register", func(c *fiber.Ctx) error {
		type Request struct {
			Email string `json:"email"`
		}
		req := new(Request)

		if err := c.BodyParser(req); err != nil {
			log.Printf("Invalid request body: %v", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid request",
			})
		}

		// Validate email format
		if !emailRegex.MatchString(req.Email) {
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
