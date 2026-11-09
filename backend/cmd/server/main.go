package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/narinsak-u/backend/internal/bootstrap"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	if err := bootstrap.Run(); err != nil {
		log.Fatal().Err(err).Msg("server startup failed")
	}
}
