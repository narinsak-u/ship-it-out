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

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	config.Load()
	database.ConnectPostgres(config.App.DatabaseURL)
	database.ConnectRedis(config.App.RedisURL)

	database.DB.AutoMigrate(&models.User{}, &models.Shipment{}, &models.ShipmentEvent{})

	app := fiber.New()

	app.Use(middleware.CORS())
	app.Use(middleware.Logger())

	api := app.Group("/api")

	authGroup := api.Group("/auth")
	authGroup.Post("/register", auth.Register)
	authGroup.Post("/login", auth.Login)

	shipmentGroup := api.Group("/shipments", middleware.AuthRequired())
	shipmentGroup.Get("/", shipment.List)
	shipmentGroup.Post("/", shipment.Create)
	shipmentGroup.Get("/:id", shipment.GetByID)
	shipmentGroup.Patch("/:id/status", shipment.UpdateStatus)

	api.Get("/track/:trackingNumber", tracking.Track)
	api.Get("/analytics/overview", middleware.AuthRequired(), analytics.Overview)

	app.Get("/ws/tracking/:trackingNumber", websocket.HandleWebSocket)
	app.Get("/ws/admin", websocket.HandleWebSocket)
	app.Get("/ws/driver", websocket.HandleWebSocket)

	log.Info().Str("port", config.App.Port).Msg("server starting")
	if err := app.Listen(":" + config.App.Port); err != nil {
		log.Fatal().Err(err).Msg("server failed to start")
	}
}
