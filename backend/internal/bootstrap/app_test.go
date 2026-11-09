package bootstrap

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewAppPreservesRouteComposition(t *testing.T) {
	app := NewApp(nil)

	health, err := app.Test(httptest.NewRequest("GET", "/api/health", nil))
	require.NoError(t, err)
	defer health.Body.Close()
	require.Equal(t, 200, health.StatusCode)
	var body struct {
		Success bool `json:"success"`
		Data    struct {
			Status string `json:"status"`
		} `json:"data"`
	}
	require.NoError(t, json.NewDecoder(health.Body).Decode(&body))
	require.True(t, body.Success)
	require.Equal(t, "ok", body.Data.Status)

	write, err := app.Test(httptest.NewRequest("POST", "/api/hubs", nil))
	require.NoError(t, err)
	defer write.Body.Close()
	require.Equal(t, 401, write.StatusCode)
}
