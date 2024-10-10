package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"log"
	"os"
	"regexp"
	"time"
)

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

func main() {

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: os.Getenv("ALLOWED_ORIGINS"),
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	if err := InitDB("waitlistdb"); err != nil {
		log.Fatalf("Failed to initialize LevelDB: %v", err)
	}

	defer CloseDB()

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

	// Waitlist registration route
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


	// Graceful shutdown
	app.Use(func(c *fiber.Ctx) error {
		return c.Next()
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Fatal(app.Listen(":" + port))
}
