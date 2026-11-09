package bootstrap

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	"github.com/narinsak-u/backend/internal/analytics"
	"github.com/narinsak-u/backend/internal/auth"
	"github.com/narinsak-u/backend/internal/platform/config"
	"github.com/narinsak-u/backend/internal/platform/database"
	"github.com/narinsak-u/backend/internal/hub"
	"github.com/narinsak-u/backend/internal/middleware"
	"github.com/narinsak-u/backend/internal/models"
	"github.com/narinsak-u/backend/internal/seed"
	"github.com/narinsak-u/backend/internal/server"
	"github.com/narinsak-u/backend/internal/shipment"
	"github.com/narinsak-u/backend/internal/tracking"
)

// NewApp composes shared middleware and existing feature routes around db.
func NewApp(db *gorm.DB) *fiber.App {
	app := fiber.New()
	app.Use(middleware.CORS())
	app.Use(middleware.SecurityHeaders())
	app.Use(middleware.Logger())

	authRepo := auth.NewGormRepository(db)
	authService := auth.NewService(authRepo, config.App.JWTSecret, config.App.JWTTTL)
	authHandler := auth.NewHandler(authService)
	hubRepo := hub.NewGormRepository(db)
	hubHandler := hub.NewHandler(hub.NewService(hubRepo))
	shipmentRepo := shipment.NewGormRepository(db)
	shipmentHandler := shipment.NewHandler(shipment.NewService(shipmentRepo, hubRepo))
	trackingHandler := tracking.NewHandler(tracking.NewService(shipmentRepo))
	analyticsRepo := analytics.NewGormRepository(db)
	analyticsHandler := analytics.NewHandler(analytics.NewService(analyticsRepo))
	server.RegisterRoutes(app, authRepo, authHandler, hubHandler, shipmentHandler, trackingHandler, analyticsHandler)
	return app
}

// Run loads configuration, initializes persistent state, and serves the API.
func Run() error {
	config.Load()
	config.Validate()
	database.ConnectPostgres(config.App.DatabaseURL)
	db := database.DB
	if db == nil {
		return fmt.Errorf("database connection is unavailable")
	}
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("get database connection: %w", err)
	}
	defer sqlDB.Close()

	if err := db.AutoMigrate(&models.User{}, &models.Shipment{}, &models.ShipmentEvent{}, &models.Hub{}); err != nil {
		return fmt.Errorf("database migration failed: %w", err)
	}
	if err := seed.BootstrapAdmin(db, config.App.AdminName, config.App.AdminEmail, config.App.AdminPassword, config.App.AdminRoleCutover); err != nil {
		return fmt.Errorf("admin bootstrap failed: %w", err)
	}
	if err := seed.BootstrapDemoUser(db); err != nil {
		return fmt.Errorf("demo account bootstrap failed: %w", err)
	}
	if err := seed.SeedHubs(db); err != nil {
		return fmt.Errorf("hub seed failed: %w", err)
	}
	if err := seed.SeedShipments(db); err != nil {
		return fmt.Errorf("shipment seed failed: %w", err)
	}

	log.Info().Str("port", config.App.Port).Msg("server starting")
	return serve(NewApp(db), config.App.Port)
}

func serve(app *fiber.App, port string) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	listenErr := make(chan error, 1)
	go func() { listenErr <- app.Listen(":" + port) }()

	select {
	case err := <-listenErr:
		if err != nil {
			return fmt.Errorf("server failed to start: %w", err)
		}
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := app.ShutdownWithContext(shutdownCtx); err != nil {
			return fmt.Errorf("server shutdown failed: %w", err)
		}
	}
	return nil
}
