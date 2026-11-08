package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/narinsak-u/backend/internal/analytics"
	"github.com/narinsak-u/backend/internal/auth"
	"github.com/narinsak-u/backend/internal/config"
	"github.com/narinsak-u/backend/internal/database"
	"github.com/narinsak-u/backend/internal/hub"
	"github.com/narinsak-u/backend/internal/middleware"
	"github.com/narinsak-u/backend/internal/models"
	"github.com/narinsak-u/backend/internal/seed"
	"github.com/narinsak-u/backend/internal/server"
	"github.com/narinsak-u/backend/internal/shipment"
	"github.com/narinsak-u/backend/internal/tracking"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	config.Load()
	config.Validate()
	database.ConnectPostgres(config.App.DatabaseURL)

	if err := database.DB.AutoMigrate(&models.User{}, &models.Shipment{}, &models.ShipmentEvent{}, &models.Hub{}); err != nil {
		log.Fatal().Err(err).Msg("database migration failed")
	}
	if err := seed.BootstrapAdmin(database.DB, config.App.AdminName, config.App.AdminEmail, config.App.AdminPassword, config.App.AdminRoleCutover); err != nil {
		log.Fatal().Err(err).Msg("admin bootstrap failed")
	}
	if err := seed.BootstrapDemoUser(database.DB); err != nil {
		log.Fatal().Err(err).Msg("demo account bootstrap failed")
	}
	if err := seed.SeedHubs(database.DB); err != nil {
		log.Fatal().Err(err).Msg("hub seed failed")
	}
	if err := seed.SeedShipments(database.DB); err != nil {
		log.Fatal().Err(err).Msg("shipment seed failed")
	}

	app := fiber.New()
	app.Use(middleware.CORS())
	app.Use(middleware.SecurityHeaders())
	app.Use(middleware.Logger())
	authRepo := auth.NewGormRepository(database.DB)
	authService := auth.NewService(authRepo, config.App.JWTSecret, config.App.JWTTTL)
	authHandler := auth.NewHandler(authService)
	hubRepo := hub.NewGormRepository(database.DB)
	hubHandler := hub.NewHandler(hub.NewService(hubRepo))
	shipmentRepo := shipment.NewGormRepository(database.DB)
	shipmentHandler := shipment.NewHandler(shipmentRepo, shipment.NewService(shipmentRepo, hubRepo))
	trackingHandler := tracking.NewHandler(tracking.NewService(shipmentRepo))
	analyticsRepo := analytics.NewGormRepository(database.DB)
	analyticsHandler := analytics.NewHandler(analytics.NewService(analyticsRepo))
	server.RegisterRoutes(app, authRepo, authHandler, hubHandler, shipmentHandler, trackingHandler, analyticsHandler)

	log.Info().Str("port", config.App.Port).Msg("server starting")
	if err := app.Listen(":" + config.App.Port); err != nil {
		log.Fatal().Err(err).Msg("server failed to start")
	}
}
