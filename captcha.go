package main

import (
	"github.com/gofiber/fiber/v3"
)

// Route to expose hCaptcha site key
func routeHCaptchaSiteKey(app *fiber.App, stub bool, siteKey string) {
	if stub {
		app.Get("/captcha", func(c fiber.Ctx) error {
			return c.JSON(fiber.Map{
				"success": false,
				"message": "hCaptcha is not enabled on this instance",
			})
		})
	} else {
		app.Get("/captcha", func(c fiber.Ctx) error {
			return c.JSON(fiber.Map{
				"success":           true,
				"hcaptcha_site_key": siteKey,
			})
		})
	}
}
