package main

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/narinsak-u/backend/internal/analytics"
	"github.com/narinsak-u/backend/internal/auth"
	"github.com/narinsak-u/backend/internal/health"
	"github.com/narinsak-u/backend/internal/hub"
	"github.com/narinsak-u/backend/internal/middleware"
	"github.com/narinsak-u/backend/internal/shipment"
	"github.com/narinsak-u/backend/internal/tracking"
)

func registerRoutes(app *fiber.App, db *gorm.DB) {
	api := app.Group("/api")
	api.Get("/health", health.Check)

	authRepo := auth.NewGormRepository(db)
	authHandler := auth.NewHandler(authRepo)
	authGroup := api.Group("/auth")
	authGroup.Post("/register", middleware.RateLimitAuth(), authHandler.Register)
	authGroup.Post("/login", middleware.RateLimitAuth(), authHandler.Login)
	authGroup.Get("/me", middleware.AuthRequired(), authHandler.Me)
	authGroup.Post("/logout", authHandler.Logout)

	hubRepo := hub.NewGormRepository(db)
	hubHandler := hub.NewHandler(hubRepo)
	api.Get("/hubs", hubHandler.List)
	api.Get("/hubs/:id", hubHandler.GetByID)
	hubGroup := api.Group("/hubs", middleware.AuthRequired(), middleware.StaffOnly(authRepo))
	hubGroup.Post("/", hubHandler.Create)
	hubGroup.Put("/:id", hubHandler.Update)
	hubGroup.Delete("/:id", hubHandler.Delete)

	shipmentRepo := shipment.NewGormRepository(db)
	shipmentHandler := shipment.NewHandler(shipmentRepo, hubRepo)
	api.Get("/shipments", shipmentHandler.List)
	api.Get("/shipments/:orderId", shipmentHandler.GetByID)
	shipmentGroup := api.Group("/shipments", middleware.AuthRequired(), middleware.StaffOnly(authRepo))
	shipmentGroup.Post("/", shipmentHandler.Create)
	shipmentGroup.Patch("/:orderId/status", shipmentHandler.UpdateStatus)
	shipmentGroup.Put("/:orderId", shipmentHandler.Update)
	shipmentGroup.Delete("/:orderId", shipmentHandler.Delete)

	trackingHandler := tracking.NewHandler(shipmentRepo)
	api.Get("/track/:trackingNumber", trackingHandler.Track)

	analyticsHandler := analytics.NewHandler(shipmentRepo)
	api.Get("/analytics/overview", analyticsHandler.Overview)
	api.Get("/analytics/timeseries", analyticsHandler.TimeSeries)
}
