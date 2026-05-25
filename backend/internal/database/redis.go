package database

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

// Redis is the global Redis client handle. After ConnectRedis runs, other
// packages can use it (e.g. database.Redis.Get(...)). Redis is an in-memory
// key-value store used here for real-time features (WebSocket presence,
// tracking updates, etc.) and fast caching.
var Redis *redis.Client

// ConnectRedis parses a Redis connection URL and establishes a connection.
//
// How it works:
//  1. redis.ParseURL() converts a string like "redis://user:pass@host:port/db"
//     into a redis.Options struct (host, port, password, database number)
//  2. redis.NewClient() creates a client with those options (but doesn't connect yet)
//  3. Redis.Ping() actually sends a "PING" command to the server to verify the
//     connection works. If the server is unreachable, the app shuts down with
//     log.Fatal()
//  4. On success, a confirmation message is logged
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
