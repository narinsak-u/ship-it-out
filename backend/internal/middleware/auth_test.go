package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/narinsak-u/backend/internal/config"
	"github.com/stretchr/testify/assert"
)

func setupConfig() {
	config.App = config.Config{
		JWTSecret: "test-secret",
		JWTTTL:    24 * time.Hour,
	}
}

func validToken() string {
	claims := jwt.MapClaims{
		"user_id": float64(1),
		"role":    "customer",
		"exp":     time.Now().Add(1 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, _ := token.SignedString([]byte(config.App.JWTSecret))
	return s
}

func expiredToken() string {
	claims := jwt.MapClaims{
		"user_id": float64(1),
		"role":    "customer",
		"exp":     time.Now().Add(-1 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, _ := token.SignedString([]byte(config.App.JWTSecret))
	return s
}

func testApp() *fiber.App {
	app := fiber.New()
	app.Use(AuthRequired())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})
	return app
}

func TestAuthRequired_ValidToken_Header(t *testing.T) {
	setupConfig()
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+validToken())
	resp, _ := testApp().Test(req, 1000)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestAuthRequired_NoToken(t *testing.T) {
	setupConfig()
	resp, _ := testApp().Test(httptest.NewRequest("GET", "/test", nil), 1000)
	assert.Equal(t, 401, resp.StatusCode)
}

func TestAuthRequired_ExpiredToken(t *testing.T) {
	setupConfig()
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+expiredToken())
	resp, _ := testApp().Test(req, 1000)
	assert.Equal(t, 401, resp.StatusCode)
}

func TestAuthRequired_CookieFallback(t *testing.T) {
	setupConfig()
	req := httptest.NewRequest("GET", "/test", nil)
	req.AddCookie(&http.Cookie{Name: "jwt", Value: validToken()})
	resp, _ := testApp().Test(req, 1000)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestAuthRequired_InvalidSignature(t *testing.T) {
	setupConfig()
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalidtoken")
	resp, _ := testApp().Test(req, 1000)
	assert.Equal(t, 401, resp.StatusCode)
}
