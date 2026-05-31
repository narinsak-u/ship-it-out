package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/narinsak-u/backend/internal/config"
)

func CORS() fiber.Handler {
	return func(c *fiber.Ctx) error {
		origin := c.Get("Origin")
		allowedOrigin := resolveOrigin(origin, config.App.CORSOrigin)
		c.Set("Access-Control-Allow-Origin", allowedOrigin)
		c.Set("Access-Control-Allow-Credentials", "true")
		c.Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		if c.Method() == "OPTIONS" {
			return c.SendStatus(204)
		}
		return c.Next()
	}
}

func resolveOrigin(requestOrigin, allowed string) string {
	if allowed == "*" {
		return "*"
	}
	for _, o := range strings.Split(allowed, ",") {
		o = strings.TrimSpace(o)
		if o == requestOrigin {
			return o
		}
	}
	return allowed
}
