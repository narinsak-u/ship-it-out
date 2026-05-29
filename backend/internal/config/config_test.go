package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/narinsak-u/backend/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestLoad_Defaults(t *testing.T) {
	os.Unsetenv("PORT")
	os.Unsetenv("DATABASE_URL")
	os.Unsetenv("JWT_SECRET")
	config.Load()
	assert.Equal(t, "8080", config.App.Port)
	assert.Equal(t, "postgres://user:pass@localhost:5432/shipments", config.App.DatabaseURL)
	assert.Equal(t, "change-me", config.App.JWTSecret)
	assert.Equal(t, 24*time.Hour, config.App.JWTTTL)
}

func TestLoad_FromEnv(t *testing.T) {
	os.Setenv("PORT", "9090")
	os.Setenv("DATABASE_URL", "postgres://test:test@localhost:9999/testdb")
	os.Setenv("JWT_SECRET", "my-secret-key")
	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("JWT_SECRET")
	}()
	config.Load()
	assert.Equal(t, "9090", config.App.Port)
	assert.Equal(t, "postgres://test:test@localhost:9999/testdb", config.App.DatabaseURL)
	assert.Equal(t, "my-secret-key", config.App.JWTSecret)
}
