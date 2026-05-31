package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Port        string
	DatabaseURL string
	CORSOrigin  string
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
		CORSOrigin:  getEnv("CORS_ORIGIN", "http://localhost:5173"),
		JWTSecret:   getEnv("JWT_SECRET", ""),
		JWTTTL:      24 * time.Hour,
	}
}

func Validate() {
	if App.JWTSecret == "" || App.JWTSecret == "change-me" {
		log.Panic().Msg("JWT_SECRET must be set to a strong, non-default value in .env or environment")
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
