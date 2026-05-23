package database

import (
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
// (connection string) and stores the handle in the global DB variable.
//
// How it works:
//  1. gorm.Open() tells GORM to connect using the Postgres driver
//  2. The logger is set to "Warn" level — only slow queries and errors are printed
//  3. If the connection fails (wrong URL, Postgres not running, etc.), the
//     app shuts down immediately with log.Fatal()
//  4. On success, a confirmation message is logged
func ConnectPostgres(dsn string) {
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to postgres")
	}
	log.Info().Msg("postgres connected")
}
