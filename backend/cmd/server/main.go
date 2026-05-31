package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/narinsak-u/backend/internal/analytics"
	"github.com/narinsak-u/backend/internal/auth"
	"github.com/narinsak-u/backend/internal/config"
	"github.com/narinsak-u/backend/internal/database"
	"github.com/narinsak-u/backend/internal/health"
	"github.com/narinsak-u/backend/internal/hub"
	"github.com/narinsak-u/backend/internal/middleware"
	"github.com/narinsak-u/backend/internal/models"
	"github.com/narinsak-u/backend/internal/seed"
	"github.com/narinsak-u/backend/internal/shipment"
	"github.com/narinsak-u/backend/internal/tracking"
)

// main starts the backend server. It:
//  1. Loads config (env vars) and connects to Postgres
//  2. Runs auto-migration so database tables match our Go models
//  3. Creates a Fiber HTTP server with CORS + logging middleware
//  4. Registers all routes:
//     - /auth/* — public auth endpoints (register, login, me, logout)
//     - /shipments/* — CRUD for shipments (requires auth)
//     - /track/:trackingNumber — public tracking lookup
//     - /analytics/overview — dashboard stats
//  5. Starts the HTTP server on the configured port
func main() {
	// Use Unix timestamps (e.g. 1700000000) in log output instead of RFC3339
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// --- Bootstrap: load config, connect database, migrate schemas ---
	config.Load()
	config.Validate()
	database.ConnectPostgres(config.App.DatabaseURL)

	// Auto-create/update tables so they match our model structs
	database.DB.AutoMigrate(&models.User{}, &models.Shipment{}, &models.ShipmentEvent{}, &models.Hub{})

	// Seed demo data (skips if tables already have rows)
	seed.SeedHubs(database.DB)
	seed.SeedShipments(database.DB)

	// --- Create the Fiber app and attach global middleware ---
	app := fiber.New()

	// CORS allows the frontend (different origin) to call these APIs
	app.Use(middleware.CORS())
	// Security headers protect against common web vulnerabilities
	app.Use(middleware.SecurityHeaders())
	// Logger prints every incoming request (method, path, status, duration)
	app.Use(middleware.Logger())

	// --- Route registration ---

	// Everything under /api will have the /api prefix
	api := app.Group("/api")

	// Health check for container orchestrators (Docker, K8s liveness probe)
	api.Get("/health", health.Check)

	// --- Auth routes (public — no auth required) ---
	authHandler := auth.NewHandler(auth.NewGormRepository(database.DB))
	authGroup := api.Group("/auth")
	authGroup.Post("/register", middleware.RateLimitAuth(), authHandler.Register) // POST /api/auth/register (rate limited)
	authGroup.Post("/login", middleware.RateLimitAuth(), authHandler.Login)       // POST /api/auth/login (rate limited)
	authGroup.Get("/me", middleware.AuthRequired(), authHandler.Me)               // GET  /api/auth/me (needs valid JWT cookie)
	authGroup.Post("/logout", authHandler.Logout)                                 // POST /api/auth/logout

	// --- Hub routes (public read, auth required for write) ---
	hubRepo := hub.NewGormRepository(database.DB)
	hubHandler := hub.NewHandler(hubRepo)

	api.Get("/hubs", hubHandler.List)        // GET /api/hubs (public)
	api.Get("/hubs/:id", hubHandler.GetByID) // GET /api/hubs/:id (public)
	hubGroup := api.Group("/hubs", middleware.AuthRequired())
	hubGroup.Post("/", hubHandler.Create)
	hubGroup.Put("/:id", hubHandler.Update)
	hubGroup.Delete("/:id", hubHandler.Delete)

	// --- Shipment routes (public read, auth required for write) ---
	shipmentRepo := shipment.NewGormRepository(database.DB)
	shipmentHandler := shipment.NewHandler(shipmentRepo, hubRepo)

	api.Get("/shipments", shipmentHandler.List)             // GET    /api/shipments (public)
	api.Get("/shipments/:orderId", shipmentHandler.GetByID) // GET    /api/shipments/:orderId (public)
	shipmentGroup := api.Group("/shipments", middleware.AuthRequired())
	shipmentGroup.Post("/", shipmentHandler.Create)                       // POST   /api/shipments
	shipmentGroup.Patch("/:orderId/status", shipmentHandler.UpdateStatus) // PATCH  /api/shipments/:orderId/status
	shipmentGroup.Put("/:orderId", shipmentHandler.Update)                // PUT    /api/shipments/:orderId
	shipmentGroup.Delete("/:orderId", shipmentHandler.Delete)             // DELETE /api/shipments/:orderId

	// --- Public tracking (anyone can look up a shipment by tracking number) ---
	trackingHandler := tracking.NewHandler(shipmentRepo)
	api.Get("/track/:trackingNumber", trackingHandler.Track)

	// --- Analytics (no auth required) ---
	analyticsHandler := analytics.NewHandler(shipmentRepo)
	api.Get("/analytics/overview", analyticsHandler.Overview)
	api.Get("/analytics/timeseries", analyticsHandler.TimeSeries)

	// --- Start the server ---
	log.Info().Str("port", config.App.Port).Msg("server starting")
	if err := app.Listen(":" + config.App.Port); err != nil {
		log.Fatal().Err(err).Msg("server failed to start")
	}
}
