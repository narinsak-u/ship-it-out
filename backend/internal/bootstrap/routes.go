package bootstrap

import (
	"github.com/gofiber/fiber/v2"
	"github.com/narinsak-u/backend/internal/analytics/presentation"
	authpresentation "github.com/narinsak-u/backend/internal/auth/presentation"
	hubpresentation "github.com/narinsak-u/backend/internal/hub/presentation"
	platformhttp "github.com/narinsak-u/backend/internal/platform/http"
	shipmentpresentation "github.com/narinsak-u/backend/internal/shipment/presentation"
	trackingpresentation "github.com/narinsak-u/backend/internal/tracking/presentation"
)

func registerRoutes(app *fiber.App, authRepo authpresentation.CurrentUserRepository, authHandler *authpresentation.Handler, hubHandler *hubpresentation.Handler, shipmentHandler *shipmentpresentation.Handler, trackingHandler *trackingpresentation.Handler, analyticsHandler *presentation.Handler) {
	api := app.Group("/api")
	api.Get("/health", func(c *fiber.Ctx) error { return platformhttp.Success(c, fiber.Map{"status": "ok"}) })

	authGroup := api.Group("/auth")
	authGroup.Post("/register", authpresentation.RateLimitAuth(), authHandler.Register)
	authGroup.Post("/login", authpresentation.RateLimitAuth(), authHandler.Login)
	authGroup.Get("/me", authpresentation.AuthRequired(), authHandler.Me)
	authGroup.Post("/logout", authHandler.Logout)

	api.Get("/hubs", hubHandler.List)
	api.Get("/hubs/:id", hubHandler.GetByID)
	hubGroup := api.Group("/hubs", authpresentation.AuthRequired(), authpresentation.StaffOnly(authRepo))
	hubGroup.Post("/", hubHandler.Create)
	hubGroup.Put("/:id", hubHandler.Update)
	hubGroup.Delete("/:id", hubHandler.Delete)

	api.Get("/shipments", shipmentHandler.List)
	api.Get("/shipments/:orderId", shipmentHandler.GetByID)
	shipmentGroup := api.Group("/shipments", authpresentation.AuthRequired(), authpresentation.StaffOnly(authRepo))
	shipmentGroup.Post("/", shipmentHandler.Create)
	shipmentGroup.Patch("/:orderId/status", shipmentHandler.UpdateStatus)
	shipmentGroup.Put("/:orderId", shipmentHandler.Update)
	shipmentGroup.Delete("/:orderId", shipmentHandler.Delete)

	api.Get("/track/:trackingNumber", trackingHandler.Track)
	api.Get("/analytics/overview", analyticsHandler.Overview)
	api.Get("/analytics/timeseries", analyticsHandler.TimeSeries)
}
