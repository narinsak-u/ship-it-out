package http_test

import (
	"encoding/json"
	"testing"

	"github.com/gofiber/fiber/v2"
	platformhttp "github.com/narinsak-u/backend/internal/platform/http"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)


func TestSuccess(t *testing.T) {
	app := fiber.New()
	c := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(c)

	err := platformhttp.Success(c, fiber.Map{"key": "value"})
	assert.NoError(t, err)
	assert.Equal(t, 200, c.Response().StatusCode())

	var resp map[string]interface{}
	json.Unmarshal(c.Response().Body(), &resp)
	assert.True(t, resp["success"].(bool))
	assert.Equal(t, "value", resp["data"].(map[string]interface{})["key"])
}

func TestError(t *testing.T) {
	app := fiber.New()
	c := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(c)

	err := platformhttp.Error(c, 400, "bad request")
	assert.NoError(t, err)
	assert.Equal(t, 400, c.Response().StatusCode())

	var resp map[string]interface{}
	json.Unmarshal(c.Response().Body(), &resp)
	assert.False(t, resp["success"].(bool))
	assert.Equal(t, "bad request", resp["error"].(string))
}

func TestSuccessWithPagination(t *testing.T) {
	app := fiber.New()
	c := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(c)

	data := []string{"a", "b"}
	err := platformhttp.SuccessWithPagination(c, data, 1, 10, 25)
	assert.NoError(t, err)
	assert.Equal(t, 200, c.Response().StatusCode())

	var resp map[string]interface{}
	json.Unmarshal(c.Response().Body(), &resp)
	pagination := resp["pagination"].(map[string]interface{})
	assert.Equal(t, float64(1), pagination["page"])
	assert.Equal(t, float64(10), pagination["limit"])
	assert.Equal(t, float64(25), pagination["total"])
	assert.Equal(t, float64(3), pagination["totalPages"])
}

func TestSuccessWithPagination_ZeroTotal(t *testing.T) {
	app := fiber.New()
	c := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(c)

	err := platformhttp.SuccessWithPagination(c, []interface{}{}, 1, 10, 0)
	assert.NoError(t, err)
	var resp map[string]interface{}
	json.Unmarshal(c.Response().Body(), &resp)
	pagination := resp["pagination"].(map[string]interface{})
	assert.Equal(t, float64(1), pagination["totalPages"])
}
