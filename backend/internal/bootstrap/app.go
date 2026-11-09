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

	"github.com/narinsak-u/backend/internal/analytics/application"
	"github.com/narinsak-u/backend/internal/analytics/infrastructure"
	"github.com/narinsak-u/backend/internal/analytics/presentation"
	authapplication "github.com/narinsak-u/backend/internal/auth/application"
	authinfra "github.com/narinsak-u/backend/internal/auth/infrastructure"
	authpresentation "github.com/narinsak-u/backend/internal/auth/presentation"
	"github.com/narinsak-u/backend/internal/platform/config"
	"github.com/narinsak-u/backend/internal/platform/database"
	hubapplication "github.com/narinsak-u/backend/internal/hub/application"
	hubinfra "github.com/narinsak-u/backend/internal/hub/infrastructure"
	hubpresentation "github.com/narinsak-u/backend/internal/hub/presentation"
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

	authRepo := authinfra.NewGormRepository(db)
	authService := authapplication.NewService(authRepo, authinfra.PasswordAdapter{}, authinfra.NewJWTAdapter(config.App.JWTSecret, config.App.JWTTTL))
	authHandler := authpresentation.NewHandler(authService)
	hubRepo := hubinfra.NewGormRepository(db)
	hubHandler := hubpresentation.NewHandler(hubapplication.NewService(hubRepo))
	shipmentRepo := shipment.NewGormRepository(db)
	shipmentHandler := shipment.NewHandler(shipment.NewService(shipmentRepo, hubRepo))
	trackingHandler := tracking.NewHandler(tracking.NewService(shipmentRepo))
	analyticsRepo := infrastructure.NewGormRepository(db)
	analyticsHandler := presentation.NewHandler(application.NewService(analyticsRepo, analyticsRepo))
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

	if err := db.AutoMigrate(&models.User{}, &models.Shipment{}, &models.ShipmentEvent{}, &hubinfra.Record{}); err != nil {
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
