package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
)

// Route to expose hCaptcha site key
func routeHCaptchaSiteKey(app *fiber.App) {
	app.Get("/hcaptcha-site-key", func(c fiber.Ctx) error {
		if conf[hCaptchaSiteKey] == "unset" {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "hCaptcha site key not configured",
			})
		}
		return c.JSON(fiber.Map{
			"hcaptcha_site_key": conf[hCaptchaSiteKey],
		})
	})
}

// Middleware to conditionally apply hCaptcha
func conditionalCaptcha(captcha fiber.Handler) fiber.Handler {
	return func(c fiber.Ctx) error {
		if conf[hCaptchaSecret] == "unset" {
			if verbose {
				log.Printf("Captcha bypassed for %q", c.IP())
			}
			return c.Next()
		}
		return captcha(c)
	}
}
