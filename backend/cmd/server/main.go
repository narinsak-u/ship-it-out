package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/narinsak-u/backend/internal/analytics"
	"github.com/narinsak-u/backend/internal/auth"
	"github.com/narinsak-u/backend/internal/config"
	"github.com/narinsak-u/backend/internal/database"
	"github.com/narinsak-u/backend/internal/middleware"
	"github.com/narinsak-u/backend/internal/models"
	"github.com/narinsak-u/backend/internal/shipment"
	"github.com/narinsak-u/backend/internal/tracking"
	"github.com/narinsak-u/backend/internal/websocket"
)

// main is the entry point that starts the entire backend server. It:
//  1. Loads config (env vars) and connects to Postgres + Redis databases
//  2. Runs auto-migration so database tables match our Go models
//  3. Creates a Fiber HTTP server with CORS + logging middleware
//  4. Registers all routes:
//     - /auth/* — public auth endpoints (register, login, me, logout)
//     - /shipments/* — CRUD for shipments (requires auth)
//     - /track/:trackingNumber — public tracking lookup
//     - /analytics/overview — dashboard stats (requires auth)
//     - /ws/* — real-time WebSocket connections
//  5. Starts the HTTP server on the configured port
func main() {
	// Use Unix timestamps (e.g. 1700000000) in log output instead of RFC3339
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// --- Bootstrap: load config, connect databases, migrate schemas ---
	config.Load()
	database.ConnectPostgres(config.App.DatabaseURL)
	database.ConnectRedis(config.App.RedisURL)

	// Auto-create/update tables so they match our model structs
	database.DB.AutoMigrate(&models.User{}, &models.Shipment{}, &models.ShipmentEvent{})

	// --- Create the Fiber app and attach global middleware ---
	app := fiber.New()

	// CORS allows the frontend (different origin) to call these APIs
	app.Use(middleware.CORS())
	// Logger prints every incoming request (method, path, status, duration)
	app.Use(middleware.Logger())

	// --- Route registration ---

	// Everything under /api will have the /api prefix
	api := app.Group("/api")

	// --- Auth routes (public — no auth required) ---
	authGroup := api.Group("/auth")
	authGroup.Post("/register", auth.Register)               // POST /api/auth/register
	authGroup.Post("/login", auth.Login)                     // POST /api/auth/login
	authGroup.Get("/me", middleware.AuthRequired(), auth.Me) // GET  /api/auth/me (needs valid JWT cookie)
	authGroup.Post("/logout", auth.Logout)                   // POST /api/auth/logout

	// --- Shipment routes (auth required) ---
	shipmentGroup := api.Group("/shipments", middleware.AuthRequired())
	shipmentGroup.Get("/", shipment.List)                     // GET    /api/shipments
	shipmentGroup.Post("/", shipment.Create)                  // POST   /api/shipments
	shipmentGroup.Get("/:id", shipment.GetByID)               // GET    /api/shipments/:id
	shipmentGroup.Patch("/:id/status", shipment.UpdateStatus) // PATCH  /api/shipments/:id/status

	// --- Public tracking (anyone can look up a shipment by tracking number) ---
	api.Get("/track/:trackingNumber", tracking.Track)

	// --- Analytics (auth required) ---
	api.Get("/analytics/overview", middleware.AuthRequired(), analytics.Overview)

	// --- WebSocket endpoints for real-time tracking updates ---
	app.Get("/ws/tracking/:trackingNumber", websocket.HandleWebSocket)
	app.Get("/ws/admin", websocket.HandleWebSocket)
	app.Get("/ws/driver", websocket.HandleWebSocket)

	// --- Start the server ---
	log.Info().Str("port", config.App.Port).Msg("server starting")
	if err := app.Listen(":" + config.App.Port); err != nil {
		log.Fatal().Err(err).Msg("server failed to start")
	}
}
