package health

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestCheck(t *testing.T) {
	app := fiber.New()
	app.Get("/api/health", Check)

	resp, _ := app.Test(httptest.NewRequest("GET", "/api/health", nil), 1000)
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode)

	var body map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&body)
	assert.True(t, body["success"].(bool))
	data := body["data"].(map[string]interface{})
	assert.Equal(t, "ok", data["status"])
}
