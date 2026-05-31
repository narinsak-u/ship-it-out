package config_test

import (
	"testing"
	"time"

	"github.com/narinsak-u/backend/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestLoad_Defaults(t *testing.T) {
	oldCfg := config.App
	t.Cleanup(func() { config.App = oldCfg })

	config.Load()
	assert.Equal(t, "8080", config.App.Port)
	assert.Equal(t, "postgres://user:pass@localhost:5432/shipments", config.App.DatabaseURL)
	assert.Equal(t, "", config.App.JWTSecret)
	assert.Equal(t, 24*time.Hour, config.App.JWTTTL)
}

func TestLoad_FromEnv(t *testing.T) {
	t.Setenv("PORT", "9090")
	t.Setenv("DATABASE_URL", "postgres://test:test@localhost:9999/testdb")
	t.Setenv("JWT_SECRET", "my-secret-key")

	oldCfg := config.App
	t.Cleanup(func() { config.App = oldCfg })

	config.Load()
	assert.Equal(t, "9090", config.App.Port)
	assert.Equal(t, "postgres://test:test@localhost:9999/testdb", config.App.DatabaseURL)
	assert.Equal(t, "my-secret-key", config.App.JWTSecret)
}

func TestValidate_PanicsOnEmptySecret(t *testing.T) {
	oldCfg := config.App
	t.Cleanup(func() { config.App = oldCfg })

	config.Load()
	assert.Panics(t, func() { config.Validate() })
}

func TestValidate_PanicsOnDefaultSecret(t *testing.T) {
	oldCfg := config.App
	t.Cleanup(func() { config.App = oldCfg })

	t.Setenv("JWT_SECRET", "change-me")
	config.Load()
	assert.Panics(t, func() { config.Validate() })
}

func TestValidate_PassesWithValidSecret(t *testing.T) {
	oldCfg := config.App
	t.Cleanup(func() { config.App = oldCfg })

	t.Setenv("JWT_SECRET", "a-strong-random-secret-that-is-at-least-32-chars-long!!")
	config.Load()
	assert.NotPanics(t, func() { config.Validate() })
	assert.Equal(t, "a-strong-random-secret-that-is-at-least-32-chars-long!!", config.App.JWTSecret)
}
