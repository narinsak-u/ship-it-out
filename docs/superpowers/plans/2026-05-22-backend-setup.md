# Backend Scaffold Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Scaffold the Go backend — folder structure, all dependencies, Docker setup, and a working HTTP server with model layer, route stubs, and env config.

**Architecture:** Fiber-based HTTP server with GORM for PostgreSQL, go-redis for Redis, gorilla/websocket for realtime. Layered as cmd entrypoint → internal packages → pkg utilities.

**Tech Stack:** Go 1.24, Fiber v2, GORM, PostgreSQL 16, Redis 7, gorilla/websocket, golang-jwt v5, zerolog

---

## File Structure (Creation Order)

```
backend/
├── cmd/server/main.go
├── internal/
│   ├── config/config.go
│   ├── database/postgres.go
│   ├── database/redis.go
│   ├── models/user.go
│   ├── models/shipment.go
│   ├── models/shipment_event.go
│   ├── middleware/auth.go
│   ├── middleware/cors.go
│   ├── middleware/logger.go
│   ├── auth/handler.go
│   ├── shipment/handler.go
│   ├── tracking/handler.go
│   ├── analytics/handler.go
│   ├── websocket/hub.go
│   └── websocket/client.go
├── pkg/utils/hash.go
├── pkg/utils/response.go
├── Dockerfile
├── docker-compose.yml
├── .env.example
├── go.mod
└── go.sum
```

---

### Task 1: Create folder structure and go.mod

**Files:**

- Create: `backend/internal/config/.gitkeep`
- Create: `backend/internal/database/.gitkeep`
- Create: `backend/internal/models/.gitkeep`
- Create: `backend/internal/middleware/.gitkeep`
- Create: `backend/internal/auth/.gitkeep`
- Create: `backend/internal/shipment/.gitkeep`
- Create: `backend/internal/tracking/.gitkeep`
- Create: `backend/internal/analytics/.gitkeep`
- Create: `backend/internal/websocket/.gitkeep`
- Create: `backend/pkg/utils/.gitkeep`
- Create: `backend/cmd/server/.gitkeep`
- Modify: `backend/go.mod`

- [ ] **Step 1: Create directory structure**

```bash
mkdir -p backend/internal/{config,database,models,middleware,auth,shipment,tracking,analytics,websocket}
mkdir -p backend/pkg/utils
mkdir -p backend/cmd/server
touch backend/internal/config/.gitkeep
touch backend/internal/database/.gitkeep
touch backend/internal/models/.gitkeep
touch backend/internal/middleware/.gitkeep
touch backend/internal/auth/.gitkeep
touch backend/internal/shipment/.gitkeep
touch backend/internal/tracking/.gitkeep
touch backend/internal/analytics/.gitkeep
touch backend/internal/websocket/.gitkeep
touch backend/pkg/utils/.gitkeep
```

- [ ] **Step 2: Read existing go.mod**

Run: `cat backend/go.mod`

The current content should be:

```
module github.com/narinsak-u/backend

go 1.24.2
```

- [ ] **Step 3: Add backend .gitkeep files to git**

```bash
git add backend/internal/**/.gitkeep backend/pkg/utils/.gitkeep backend/cmd/server/.gitkeep
git commit -m "chore: create backend folder structure"
```

---

### Task 2: Install Go dependencies

**Files:**

- Modify: `backend/go.mod`
- Create: `backend/go.sum`

- [ ] **Step 1: Install all backend dependencies**

```bash
cd backend
go get github.com/gofiber/fiber/v2@latest
go get github.com/gofiber/contrib/websocket@latest
go get github.com/gorilla/websocket@latest
go get gorm.io/gorm@latest
go get gorm.io/driver/postgres@latest
go get github.com/redis/go-redis/v9@latest
go get github.com/golang-jwt/jwt/v5@latest
go get github.com/joho/godotenv@latest
go get github.com/rs/zerolog@latest
```

Expected: all packages download and resolve. `go.sum` created.

- [ ] **Step 2: Verify go.mod has all dependencies**

Run: `cat backend/go.mod`

Expected: module declaration plus all `require` blocks with version pins.

- [ ] **Step 3: Commit**

```bash
git add backend/go.mod backend/go.sum
git commit -m "feat: add backend Go dependencies"
```

---

### Task 3: Config package

**Files:**

- Create: `backend/internal/config/config.go`

- [ ] **Step 1: Write config/config.go**

```go
package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

type Config struct {
	Port        string
	DatabaseURL string
	RedisURL    string
	JWTSecret   string
	JWTTTL     time.Duration
}

var App Config

func Load() {
	if err := godotenv.Load(); err != nil {
		l := zerolog.New(os.Stderr)
		l.Warn().Msg(".env file not found, using system env")
	}

	App = Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://user:pass@localhost:5432/shipments"),
		RedisURL:    getEnv("REDIS_URL", "redis://localhost:6379"),
		JWTSecret:   getEnv("JWT_SECRET", "change-me"),
		JWTTTL:     24 * time.Hour,
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
```

- [ ] **Step 2: Commit**

```bash
git add backend/internal/config/config.go
git commit -m "feat: add config loader with env support"
```

---

### Task 4: Database package

**Files:**

- Create: `backend/internal/database/postgres.go`
- Create: `backend/internal/database/redis.go`

- [ ] **Step 1: Write database/postgres.go**

```go
package database

import (
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectPostgres(dsn string) {
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to postgres")
	}
	log.Info().Msg("postgres connected")
}
```

- [ ] **Step 2: Write database/redis.go**

```go
package database

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

var Redis *redis.Client

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
```

- [ ] **Step 3: Commit**

```bash
git add backend/internal/database/
git commit -m "feat: add postgres and redis connection helpers"
```

---

### Task 5: Models package

**Files:**

- Create: `backend/internal/models/user.go`
- Create: `backend/internal/models/shipment.go`
- Create: `backend/internal/models/shipment_event.go`

- [ ] **Step 1: Write models/user.go**

```go
package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Email     string    `gorm:"unique;not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"`
	Role      string    `gorm:"not null;default:customer" json:"role"`
	CreatedAt time.Time `json:"created_at"`
}
```

- [ ] **Step 2: Write models/shipment.go**

```go
package models

import "time"

type Shipment struct {
	ID                uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	TrackingNumber    string    `gorm:"unique;index;not null" json:"tracking_number"`
	SenderName        string    `gorm:"not null" json:"sender_name"`
	ReceiverName      string    `gorm:"not null" json:"receiver_name"`
	OriginAddress     string    `gorm:"not null" json:"origin_address"`
	DestinationAddress string   `gorm:"not null" json:"destination_address"`
	Weight            float64   `gorm:"not null" json:"weight"`
	Status            string    `gorm:"not null;default:CREATED" json:"status"`
	EstimatedDelivery time.Time `json:"estimated_delivery"`
	CreatedAt         time.Time `json:"created_at"`
}
```

- [ ] **Step 3: Write models/shipment_event.go**

```go
package models

import "time"

type ShipmentEvent struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ShipmentID  uint      `gorm:"not null;index" json:"shipment_id"`
	Status      string    `gorm:"not null" json:"status"`
	Location    string    `json:"location,omitempty"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}
```

- [ ] **Step 4: Commit**

```bash
git add backend/internal/models/
git commit -m "feat: add User, Shipment, ShipmentEvent models"
```

---

### Task 6: Utils package

**Files:**

- Create: `backend/pkg/utils/hash.go`
- Create: `backend/pkg/utils/response.go`

- [ ] **Step 1: Write utils/hash.go**

```go
package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
```

Note: `golang.org/x/crypto` is a transitive dependency from golang-jwt. It will already be in go.sum.

- [ ] **Step 2: Write utils/response.go**

```go
package utils

import "github.com/gofiber/fiber/v2"

func Success(c *fiber.Ctx, data interface{}) error {
	return c.JSON(fiber.Map{"success": true, "data": data})
}

func Error(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{"success": false, "error": message})
}
```

- [ ] **Step 3: Commit**

```bash
git add backend/pkg/utils/
git commit -m "feat: add hash and response utilities"
```

---

### Task 7: Middleware package

**Files:**

- Create: `backend/internal/middleware/auth.go`
- Create: `backend/internal/middleware/cors.go`
- Create: `backend/internal/middleware/logger.go`

- [ ] **Step 1: Write middleware/cors.go**

```go
package middleware

import "github.com/gofiber/fiber/v2"

func CORS() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		if c.Method() == "OPTIONS" {
			return c.SendStatus(204)
		}
		return c.Next()
	}
}
```

- [ ] **Step 2: Write middleware/logger.go**

```go
package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start)
		log.Info().
			Str("method", c.Method()).
			Str("path", c.Path()).
			Int("status", c.Response().StatusCode()).
			Dur("duration", duration).
			Msg("request")
		return err
	}
}
```

- [ ] **Step 3: Write middleware/auth.go**

```go
package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/narinsak-u/backend/internal/config"
	"github.com/narinsak-u/backend/pkg/utils"
)

func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			return utils.Error(c, 401, "missing or invalid token")
		}
		tokenStr := strings.TrimPrefix(auth, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(config.App.JWTSecret), nil
		})
		if err != nil || !token.Valid {
			return utils.Error(c, 401, "invalid or expired token")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return utils.Error(c, 401, "invalid token claims")
		}

		c.Locals("user_id", claims["user_id"])
		c.Locals("role", claims["role"])
		return c.Next()
	}
}
```

- [ ] **Step 4: Commit**

```bash
git add backend/internal/middleware/
git commit -m "feat: add CORS, logger, and JWT auth middleware"
```

---

### Task 8: Auth handlers

**Files:**

- Create: `backend/internal/auth/handler.go`

- [ ] **Step 1: Write auth/handler.go**

```go
package auth

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/narinsak-u/backend/internal/config"
	"github.com/narinsak-u/backend/internal/database"
	"github.com/narinsak-u/backend/internal/models"
	"github.com/narinsak-u/backend/pkg/utils"
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, 400, "invalid request body")
	}
	if req.Name == "" || req.Email == "" || req.Password == "" {
		return utils.Error(c, 400, "name, email, and password are required")
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return utils.Error(c, 500, "failed to hash password")
	}

	role := req.Role
	if role == "" {
		role = "customer"
	}

	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hash,
		Role:     role,
	}

	if result := database.DB.Create(&user); result.Error != nil {
		return utils.Error(c, 409, "email already registered")
	}

	return utils.Success(c, fiber.Map{"user": user})
}

func Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, 400, "invalid request body")
	}

	var user models.User
	if result := database.DB.Where("email = ?", req.Email).First(&user); result.Error != nil {
		return utils.Error(c, 401, "invalid email or password")
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		return utils.Error(c, 401, "invalid email or password")
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(config.App.JWTTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(config.App.JWTSecret))
	if err != nil {
		return utils.Error(c, 500, "failed to generate token")
	}

	return utils.Success(c, fiber.Map{
		"token": tokenStr,
		"user":  user,
	})
}
```

- [ ] **Step 2: Commit**

```bash
git add backend/internal/auth/
git commit -m "feat: add register and login handlers"
```

---

### Task 9: Shipment handlers

**Files:**

- Create: `backend/internal/shipment/handler.go`

- [ ] **Step 1: Write shipment/handler.go**

```go
package shipment

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/narinsak-u/backend/internal/database"
	"github.com/narinsak-u/backend/internal/models"
	"github.com/narinsak-u/backend/pkg/utils"
)

type CreateRequest struct {
	SenderName        string  `json:"sender_name"`
	ReceiverName      string  `json:"receiver_name"`
	OriginAddress     string  `json:"origin_address"`
	DestinationAddress string `json:"destination_address"`
	Weight            float64 `json:"weight"`
}

func generateTrackingNumber() string {
	return fmt.Sprintf("TH%d%05d", time.Now().Year(), time.Now().UnixMilli()%100000)
}

func List(c *fiber.Ctx) error {
	var shipments []models.Shipment
	database.DB.Find(&shipments)
	return utils.Success(c, shipments)
}

func Create(c *fiber.Ctx) error {
	var req CreateRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, 400, "invalid request body")
	}

	shipment := models.Shipment{
		TrackingNumber:     generateTrackingNumber(),
		SenderName:         req.SenderName,
		ReceiverName:       req.ReceiverName,
		OriginAddress:      req.OriginAddress,
		DestinationAddress: req.DestinationAddress,
		Weight:             req.Weight,
		Status:             "CREATED",
	}

	if result := database.DB.Create(&shipment); result.Error != nil {
		return utils.Error(c, 500, "failed to create shipment")
	}

	return utils.Success(c, shipment)
}

func GetByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return utils.Error(c, 400, "invalid id")
	}

	var shipment models.Shipment
	if result := database.DB.First(&shipment, id); result.Error != nil {
		return utils.Error(c, 404, "shipment not found")
	}

	return utils.Success(c, shipment)
}

func UpdateStatus(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return utils.Error(c, 400, "invalid id")
	}

	var body struct {
		Status string `json:"status"`
	}
	if err := c.BodyParser(&body); err != nil {
		return utils.Error(c, 400, "invalid request body")
	}

	var shipment models.Shipment
	if result := database.DB.First(&shipment, id); result.Error != nil {
		return utils.Error(c, 404, "shipment not found")
	}

	shipment.Status = body.Status
	database.DB.Save(&shipment)

	event := models.ShipmentEvent{
		ShipmentID:  shipment.ID,
		Status:      body.Status,
		Description: fmt.Sprintf("Status updated to %s", body.Status),
	}
	database.DB.Create(&event)

	return utils.Success(c, shipment)
}
```

- [ ] **Step 2: Commit**

```bash
git add backend/internal/shipment/
git commit -m "feat: add shipment CRUD handlers"
```

---

### Task 10: Tracking and analytics handlers

**Files:**

- Create: `backend/internal/tracking/handler.go`
- Create: `backend/internal/analytics/handler.go`

- [ ] **Step 1: Write tracking/handler.go**

```go
package tracking

import (
	"github.com/gofiber/fiber/v2"
	"github.com/narinsak-u/backend/internal/database"
	"github.com/narinsak-u/backend/internal/models"
	"github.com/narinsak-u/backend/pkg/utils"
)

func Track(c *fiber.Ctx) error {
	trackingNumber := c.Params("trackingNumber")

	var shipment models.Shipment
	if result := database.DB.Where("tracking_number = ?", trackingNumber).First(&shipment); result.Error != nil {
		return utils.Error(c, 404, "shipment not found")
	}

	var events []models.ShipmentEvent
	database.DB.Where("shipment_id = ?", shipment.ID).Order("created_at asc").Find(&events)

	return utils.Success(c, fiber.Map{
		"shipment": shipment,
		"events":   events,
	})
}
```

- [ ] **Step 2: Write analytics/handler.go**

```go
package analytics

import (
	"github.com/gofiber/fiber/v2"
	"github.com/narinsak-u/backend/internal/database"
	"github.com/narinsak-u/backend/internal/models"
	"github.com/narinsak-u/backend/pkg/utils"
)

func Overview(c *fiber.Ctx) error {
	var total int64
	database.DB.Model(&models.Shipment{}).Count(&total)

	var active int64
	database.DB.Model(&models.Shipment{}).Where("status NOT IN ?", []string{"DELIVERED", "RETURNED"}).Count(&active)

	var delivered int64
	database.DB.Model(&models.Shipment{}).Where("status = ?", "DELIVERED").Count(&delivered)

	type StatusCount struct {
		Status string
		Count  int64
	}
	var byStatus []StatusCount
	database.DB.Model(&models.Shipment{}).Select("status, count(*) as count").Group("status").Scan(&byStatus)

	return utils.Success(c, fiber.Map{
		"total":     total,
		"active":    active,
		"delivered": delivered,
		"by_status": byStatus,
	})
}
```

- [ ] **Step 3: Commit**

```bash
git add backend/internal/tracking/ backend/internal/analytics/
git commit -m "feat: add tracking lookup and analytics overview"
```

---

### Task 11: WebSocket hub

**Files:**

- Create: `backend/internal/websocket/hub.go`
- Create: `backend/internal/websocket/client.go`

- [ ] **Step 1: Write websocket/hub.go**

```go
package websocket

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	mu      sync.RWMutex
	clients map[*Client]bool
}

var DefaultHub = &Hub{
	clients: make(map[*Client]bool),
}

func (h *Hub) Register(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[client] = true
}

func (h *Hub) Unregister(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		client.Conn.Close()
	}
}

func (h *Hub) BroadcastToRoom(room string, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for client := range h.clients {
		if client.Room == room {
			select {
			case client.Send <- message:
			default:
				close(client.Send)
				delete(h.clients, client)
			}
		}
	}
}
```

- [ ] **Step 2: Write websocket/client.go**

```go
package websocket

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/contrib/websocket"
)

type Client struct {
	Room string
	Conn *websocket.Conn
	Send chan []byte
}

func HandleWebSocket(c *fiber.Ctx) error {
	room := c.Params("trackingNumber", "global")
	if websocket.IsWebSocketUpgrade(c) {
		return websocket.New(func(conn *websocket.Conn) {
			client := &Client{
				Room: room,
				Conn: conn,
				Send: make(chan []byte, 256),
			}
			DefaultHub.Register(client)
			defer DefaultHub.Unregister(client)

			go func() {
				for msg := range client.Send {
					if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
						return
					}
				}
			}()

			for {
				_, _, err := conn.ReadMessage()
				if err != nil {
					break
				}
			}
		})(c)
	}
	return c.SendStatus(400)
}
```

- [ ] **Step 3: Commit**

```bash
git add backend/internal/websocket/
git commit -m "feat: add websocket hub and client handler"
```

---

### Task 12: Main entry point

**Files:**

- Create: `backend/cmd/server/main.go`

- [ ] **Step 1: Write cmd/server/main.go**

```go
package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/narinsak-u/backend/internal/auth"
	"github.com/narinsak-u/backend/internal/config"
	"github.com/narinsak-u/backend/internal/database"
	"github.com/narinsak-u/backend/internal/middleware"
	"github.com/narinsak-u/backend/internal/models"
	"github.com/narinsak-u/backend/internal/shipment"
	"github.com/narinsak-u/backend/internal/tracking"
	"github.com/narinsak-u/backend/internal/websocket"
	"github.com/narinsak-u/backend/internal/analytics"
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
```

- [ ] **Step 2: Verify the project compiles**

```bash
cd backend
go build ./cmd/server/
```

Expected: binary `server.exe` (Windows) or `server` (Linux) created with no compilation errors.

- [ ] **Step 3: Clean up binary**

```bash
rm -f server server.exe
```

- [ ] **Step 4: Commit**

```bash
git add backend/cmd/server/main.go
git commit -m "feat: add server entry point with routes and auto-migrate"
```

---

### Task 13: Dockerfile

**Files:**

- Create: `backend/Dockerfile`

- [ ] **Step 1: Write Dockerfile**

```dockerfile
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server/

FROM alpine:3.19

RUN addgroup -S app && adduser -S app -G app
COPY --from=builder /server /server
USER app

EXPOSE 8080
CMD ["/server"]
```

- [ ] **Step 2: Commit**

```bash
git add backend/Dockerfile
git commit -m "feat: add multi-stage Dockerfile"
```

---

### Task 14: docker-compose.yml and .env.example

**Files:**

- Create: `backend/docker-compose.yml`
- Create: `backend/.env.example`

- [ ] **Step 1: Write docker-compose.yml**

```yaml
version: "3.9"

services:
  db:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: user
      POSTGRES_PASSWORD: pass
      POSTGRES_DB: shipments
    ports:
      - "5432:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

  backend:
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    environment:
      PORT: "8080"
      DATABASE_URL: "postgres://user:pass@db:5432/shipments"
      REDIS_URL: "redis://redis:6379"
      JWT_SECRET: "change-me"
    depends_on:
      - db
      - redis

volumes:
  pgdata:
```

- [ ] **Step 2: Write .env.example**

```
PORT=8080
DATABASE_URL=postgres://user:pass@db:5432/shipments
REDIS_URL=redis://redis:6379
JWT_SECRET=change-me
```

- [ ] **Step 3: Commit**

```bash
git add backend/docker-compose.yml backend/.env.example
git commit -m "feat: add docker-compose with postgres, redis, and backend"
```

---

### Task 15: Clean up .gitkeep files and final verification

**Files:**

- Modify: `backend/` (remove .gitkeep files)

- [ ] **Step 1: Remove .gitkeep files (all directories now have real files)**

```bash
find backend -name '.gitkeep' -delete
```

- [ ] **Step 2: Full build verification**

```bash
cd backend
go vet ./...
go build ./cmd/server/
```

Expected: `go vet` passes, binary compiles with no errors.

- [ ] **Step 3: Remove binary**

```bash
rm -f server server.exe
```

- [ ] **Step 4: Final git status check**

```bash
git status --short
```

Expected: all files tracked, no uncommitted changes.

- [ ] **Step 5: Commit cleanup**

```bash
git add -A
git commit -m "chore: clean up gitkeep files and finalize backend scaffold"
```
