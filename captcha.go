package main

import (
	"github.com/gofiber/fiber/v3"
)

type respHSiteKey struct {
	Success bool   `json:"success"`
	SiteKey string `json:"hcaptcha_site_key"`
}

// Route to expose hCaptcha site key.
// Returns a constant pre-generated response
// to avoid unnecessary allocations or serialisations
func routeHCaptchaSiteKey(app *fiber.App, stub bool, siteKey string) {
	var resp string
	if stub {
		resp = mustConstResp(newMessage(false, "hCaptcha is not enabled on this instance."))
	} else {
		resp = mustConstResp(respHSiteKey{true, siteKey})
	}

	app.Get("/captcha", func(c fiber.Ctx) error {
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.SendString(resp)
	})
}
