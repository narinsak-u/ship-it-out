package database

import (
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is the global GORM database handle. Once ConnectPostgres runs, every
// other package can access the database through this variable (e.g. database.DB.Find(...)).
// GORM is an ORM — it maps Go structs (models) to database tables so we can
// query using Go methods instead of writing raw SQL.
var DB *gorm.DB

// ConnectPostgres opens a connection to PostgreSQL using the provided DSN
// with retries. Postgres often takes a few seconds to be ready in Docker,
// so we retry every 2 seconds up to 30 seconds before giving up.
func ConnectPostgres(dsn string) {
	var err error
	for i := 0; i < 15; i++ {
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Warn),
		})
		if err == nil {
			log.Info().Msg("postgres connected")
			return
		}
		log.Warn().Err(err).Msgf("postgres not ready (attempt %d/15), retrying in 2s...", i+1)
		time.Sleep(2 * time.Second)
	}
	log.Fatal().Err(err).Msg("failed to connect to postgres after 15 attempts")
}
