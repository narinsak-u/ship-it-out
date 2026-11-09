package server

import (
	"github.com/gofiber/fiber/v2"

	"github.com/narinsak-u/backend/internal/analytics/presentation"
	authpresentation "github.com/narinsak-u/backend/internal/auth/presentation"
	"github.com/narinsak-u/backend/internal/health"
	hubpresentation "github.com/narinsak-u/backend/internal/hub/presentation"
	"github.com/narinsak-u/backend/internal/middleware"
	"github.com/narinsak-u/backend/internal/shipment"
	"github.com/narinsak-u/backend/internal/tracking"
)

// RegisterRoutes attaches feature handlers and access middleware to the API.
func RegisterRoutes(app *fiber.App, authRepo middleware.UserRepository, authHandler *authpresentation.Handler, hubHandler *hubpresentation.Handler, shipmentHandler *shipment.Handler, trackingHandler *tracking.Handler, analyticsHandler *presentation.Handler) {
	api := app.Group("/api")
	api.Get("/health", health.Check)

	authGroup := api.Group("/auth")
	authGroup.Post("/register", middleware.RateLimitAuth(), authHandler.Register)
	authGroup.Post("/login", middleware.RateLimitAuth(), authHandler.Login)
	authGroup.Get("/me", middleware.AuthRequired(), authHandler.Me)
	authGroup.Post("/logout", authHandler.Logout)

	api.Get("/hubs", hubHandler.List)
	api.Get("/hubs/:id", hubHandler.GetByID)
	hubGroup := api.Group("/hubs", middleware.AuthRequired(), middleware.StaffOnly(authRepo))
	hubGroup.Post("/", hubHandler.Create)
	hubGroup.Put("/:id", hubHandler.Update)
	hubGroup.Delete("/:id", hubHandler.Delete)

	api.Get("/shipments", shipmentHandler.List)
	api.Get("/shipments/:orderId", shipmentHandler.GetByID)
	shipmentGroup := api.Group("/shipments", middleware.AuthRequired(), middleware.StaffOnly(authRepo))
	shipmentGroup.Post("/", shipmentHandler.Create)
	shipmentGroup.Patch("/:orderId/status", shipmentHandler.UpdateStatus)
	shipmentGroup.Put("/:orderId", shipmentHandler.Update)
	shipmentGroup.Delete("/:orderId", shipmentHandler.Delete)

	api.Get("/track/:trackingNumber", trackingHandler.Track)
	api.Get("/analytics/overview", analyticsHandler.Overview)
	api.Get("/analytics/timeseries", analyticsHandler.TimeSeries)
}
