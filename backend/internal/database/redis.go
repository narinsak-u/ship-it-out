package database

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

var Redis *redis.Client

func ConnectRedis(addr string) {
	opts, err := redis.ParseURL(addr)
	if err != nil {
		log.Fatal().Err(err).Msg("invalid redis url")
	}

	Redis = redis.NewClient(opts)

	if err := Redis.Ping(context.Background()).Err(); err != nil {
		log.Fatal().Err(err).Msg("failed to connect to redis")
	}
	log.Info().Msg("redis connected")
}
