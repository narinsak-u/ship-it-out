package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
	JWTTTL      time.Duration
}

var App Config

func Load() {
	if err := godotenv.Load(); err != nil {
		l := zerolog.New(os.Stderr)
		l.Warn().Msg(".env file not found, using system env")
	}

	App = Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://user:pass@localhost:5432/shipments"),
		JWTSecret:   getEnv("JWT_SECRET", "change-me"),
		JWTTTL:      24 * time.Hour,
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
