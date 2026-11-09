package middleware

import "github.com/gofiber/fiber/v2"

// SecurityHeaders returns middleware that sets HTTP security headers on every
// response. These headers help protect against common web vulnerabilities:
//
//   - X-Content-Type-Options: prevents MIME-type sniffing
//   - X-Frame-Options: prevents clickjacking by blocking <frame>/<iframe> embedding
//   - Strict-Transport-Security: enforces HTTPS (6-month max-age, all subdomains)
//   - Referrer-Policy: limits referrer info on cross-origin requests
//   - Permissions-Policy: restricts powerful browser features (geolocation, camera, etc.)
//   - Content-Security-Policy: mitigates XSS and data injection attacks
//
// Strict-Transport-Security should ideally be set at the reverse proxy level (nginx/Caddy).
// It is included here as a fallback for direct-to-Fiber deployments.
func SecurityHeaders() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Frame-Options", "DENY")
		c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Set("Permissions-Policy", "geolocation=(), microphone=(), camera=(), payment=(), usb=()")
		c.Set("Content-Security-Policy",
			"default-src 'self'; "+
				"script-src 'self'; "+
				"style-src 'self' 'unsafe-inline'; "+
				"img-src 'self' data: https:; "+
				"connect-src 'self' https://api.opencagedata.com; "+
				"frame-ancestors 'none'; "+
				"base-uri 'self'; "+
				"form-action 'self'; "+
				"object-src 'none'",
		)
		return c.Next()
	}
}
