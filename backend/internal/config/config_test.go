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
	assert.Equal(t, "change-me", config.App.JWTSecret)
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
