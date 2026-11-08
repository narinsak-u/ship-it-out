package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Port             string
	DatabaseURL      string
	CORSOrigin       string
	JWTSecret        string
	AdminName        string
	AdminEmail       string
	AdminPassword    string
	AdminRoleCutover bool
	JWTTTL           time.Duration
}

var App Config

func Load() {
	if err := godotenv.Load(); err != nil {
		l := zerolog.New(os.Stderr)
		l.Warn().Msg(".env file not found, using system env")
	}

	cutover, _ := strconv.ParseBool(getEnv("ADMIN_ROLE_CUTOVER", "false"))
	App = Config{
		Port:             getEnv("PORT", "8080"),
		DatabaseURL:      getEnv("DATABASE_URL", "postgres://user:pass@localhost:5432/shipments"),
		CORSOrigin:       getEnv("CORS_ORIGIN", "http://localhost:3000"),
		JWTSecret:        getEnv("JWT_SECRET", ""),
		AdminName:        getEnv("ADMIN_NAME", "Administrator"),
		AdminEmail:       getEnv("ADMIN_EMAIL", ""),
		AdminPassword:    getEnv("ADMIN_PASSWORD", ""),
		AdminRoleCutover: cutover,
		JWTTTL:           24 * time.Hour,
	}
}

func Validate() {
	if App.JWTSecret == "" || App.JWTSecret == "change-me" {
		log.Panic().Msg("JWT_SECRET must be set to a strong, non-default value in .env or environment")
	}
	if App.AdminEmail == "" || App.AdminPassword == "" {
		log.Panic().Msg("ADMIN_EMAIL and ADMIN_PASSWORD must be set")
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
