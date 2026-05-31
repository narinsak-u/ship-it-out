# Pre-Deployment Hardening Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Fix all deployment blockers identified in readiness audit — Redis removal, JWT secret hardening, configurable CORS, security headers, health endpoint, and doc fixes.

**Architecture:** Backend-only changes except Task 7 (frontend). Config-driven approach: add env vars, validate at startup, clean stale references. Each task is independently testable.

**Tech Stack:** Go 1.24 / Fiber v2, Vue 3 / Vite 6

---

## File Structure

| Task | Files |
|------|-------|
| 1. Redis removal | `backend/.env`, `backend/.env.example`, `README.md` |
| 2. JWT_SECRET validation | `backend/internal/config/config.go` |
| 3. Configurable CORS | `backend/internal/config/config.go`, `backend/internal/middleware/cors.go`, `backend/.env`, `backend/.env.example` |
| 4. Security headers | `backend/internal/middleware/security.go` (new), `backend/cmd/server/main.go` |
| 5. Health endpoint | `backend/internal/health/handler.go` (new), `backend/internal/health/health_test.go` (new), `backend/cmd/server/main.go` |
| 6. Doc fixes | `README.md`, `backend/docs/WORKFLOW.md` |
| 7. Bundle split | `frontend/src/views/OrderFormView.vue` |

---

### Task 1: Remove all Redis references

**Files:**
- Modify: `backend/.env`
- Modify: `backend/.env.example`
- Modify: `README.md`

**Context:** The project no longer uses Redis. The in-memory rate limiter replaces the need, and no other Redis-dependent code exists (`internal/database/` has only `postgres.go`). Remove all stale references so new deployers don't waste time configuring a service they don't need.

- [ ] **Remove REDIS_URL from `backend/.env`**

Edit `backend/.env` — delete the `REDIS_URL=redis://redis:6379` line.

- [ ] **Remove REDIS_URL from `backend/.env.example`**

Edit `backend/.env.example` — delete the `REDIS_URL` row.

File becomes:
```
PORT=8080
DATABASE_URL=postgres://user:pass@localhost:5432/shipments
JWT_SECRET=change-me
```

- [ ] **Update `README.md` to remove Redis references**

3 changes needed:

1. Tech stack table — change "Cache" row from `go-redis/v9` to `N/A (in-memory rate limiter)`
2. Docker Compose line — change "Postgres 16 + Redis 7" to "Postgres 16"
3. Backend env vars table — delete the `REDIS_URL` row
4. Backend getting-started — change "Or manually (Postgres + Redis required)" to "Or manually (Postgres required)"

---

### Task 2: JWT_SECRET startup validation

**Files:**
- Modify: `backend/internal/config/config.go`

**Context:** A JWT secret of `change-me` or an empty string is a security vulnerability. Add validation in `config.Load()` that panics with a clear message if the secret is weak, preventing the server from starting with a known/default secret.

- [ ] **Add secret validation in `config.Load()`**

```go
func Load() {
	if err := godotenv.Load(); err != nil {
		l := zerolog.New(os.Stderr)
		l.Warn().Msg(".env file not found, using system env")
	}

	App = Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://user:pass@localhost:5432/shipments"),
		JWTSecret:   getEnv("JWT_SECRET", ""),
		JWTTTL:      24 * time.Hour,
	}

	if App.JWTSecret == "" || App.JWTSecret == "change-me" {
		log.Fatal().Msg("JWT_SECRET must be set to a strong, non-default value in .env or environment")
	}
}
```

- [ ] **Run `go vet ./...` and `go build ./...` to verify**

Run: `go vet ./... && go build ./...`
Expected: clean exit (no errors)

- [ ] **Commit**

```bash
git add backend/internal/config/config.go
git commit -m "fix: validate JWT_SECRET at startup, reject default/empty"
```

---

### Task 3: Make CORS origin configurable via env var

**Files:**
- Modify: `backend/internal/config/config.go`
- Modify: `backend/internal/middleware/cors.go`
- Modify: `backend/.env`
- Modify: `backend/.env.example`

**Context:** CORS origin is hardcoded to `http://localhost:5173`. In production the frontend lives at a different domain. Read the allowed origin from `CORS_ORIGIN` env var, defaulting to `http://localhost:5173` for local dev.

- [ ] **Add `CORSOrigin` to config struct and `Load()`**

In `backend/internal/config/config.go`, add to the `Config` struct:
```go
type Config struct {
	Port        string
	DatabaseURL string
	CORSOrigin  string
	JWTSecret   string
	JWTTTL      time.Duration
}
```

In `config.Load()`, add to the `App` initialization:
```go
App = Config{
	Port:        getEnv("PORT", "8080"),
	DatabaseURL: getEnv("DATABASE_URL", "postgres://user:pass@localhost:5432/shipments"),
	CORSOrigin:  getEnv("CORS_ORIGIN", "http://localhost:5173"),
	JWTSecret:   getEnv("JWT_SECRET", ""),
	JWTTTL:      24 * time.Hour,
}
```

- [ ] **Update `backend/internal/middleware/cors.go` to use config**

Replace the hardcoded origin:

```go
package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/narinsak-u/backend/internal/config"
)

func CORS() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", config.App.CORSOrigin)
		c.Set("Access-Control-Allow-Credentials", "true")
		c.Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		if c.Method() == "OPTIONS" {
			return c.SendStatus(204)
		}
		return c.Next()
	}
}
```

- [ ] **Add `CORS_ORIGIN` to `backend/.env` and `backend/.env.example`**

Append to both files:
```
CORS_ORIGIN=http://localhost:5173
```

- [ ] **Run backend tests and build**

Run: `go vet ./... && go build ./... && go test ./...`
Expected: all pass (no test relies on CORS — no test regressions expected).

- [ ] **Commit**

```bash
git add backend/internal/config/config.go backend/internal/middleware/cors.go backend/.env backend/.env.example
git commit -m "feat: make CORS origin configurable via CORS_ORIGIN env var"
```

---

### Task 4: Add security headers middleware

**Files:**
- Create: `backend/internal/middleware/security.go`
- Modify: `backend/cmd/server/main.go`

**Context:** The server emits no security headers, leaving it vulnerable to clickjacking, MIME-sniffing, and lacking CSP/XSS protection. Add a middleware that sets recommended headers on every response.

- [ ] **Create `backend/internal/middleware/security.go`**

```go
package middleware

import "github.com/gofiber/fiber/v2"

// SecurityHeaders returns middleware that sets HTTP security headers on every
// response. These headers help protect against common web vulnerabilities:
//
//   - X-Content-Type-Options: prevents MIME-type sniffing (CVE-2012-5874)
//   - X-Frame-Options: prevents clickjacking by blocking <frame>/<iframe> embedding
//   - Referrer-Policy: controls what referrer info is sent with cross-origin requests
//   - Permissions-Policy: restricts which browser APIs pages can use
//   - Content-Security-Policy: mitigates XSS and data injection attacks
func SecurityHeaders() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Frame-Options", "DENY")
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		c.Set("Content-Security-Policy",
			"default-src 'self'; "+
				"script-src 'self'; "+
				"style-src 'self' 'unsafe-inline'; "+
				"img-src 'self' data: https:; "+
				"connect-src 'self' https://api.opencagedata.com; "+
				"frame-ancestors 'none'; "+
				"base-uri 'self'; "+
				"form-action 'self'",
		)
		return c.Next()
	}
}
```

**CSP notes:** `'unsafe-inline'` is needed for style because Tailwind v4 injects inline styles. `connect-src` allows OpenCage geocoding. When you add a TLS reverse proxy, also set `Strict-Transport-Security` at the proxy level.

- [ ] **Register the middleware in `main.go`**

In `backend/cmd/server/main.go`, add the import and middleware line. After the CORS and Logger middleware, add:

```go
// Security headers protect against common web vulnerabilities
app.Use(middleware.SecurityHeaders())
```

Place it between `app.Use(middleware.CORS())` and `app.Use(middleware.Logger())` (or any order — security headers are idempotent).

- [ ] **Build to verify**

Run: `go build ./...`
Expected: clean build.

- [ ] **Commit**

```bash
git add backend/internal/middleware/security.go backend/cmd/server/main.go
git commit -m "feat: add security headers middleware (CSP, XFO, HSTS-ready)"
```

---

### Task 5: Add health check endpoint

**Files:**
- Create: `backend/internal/health/handler.go`
- Create: `backend/internal/health/health_test.go`
- Modify: `backend/cmd/server/main.go`

**Context:** Container orchestrators (Docker, k8s) rely on health check endpoints for liveness probes. Add `GET /api/health` that returns `{"success":true,"data":{"status":"ok"}}`.

- [ ] **Create `backend/internal/health/handler.go`**

```go
package health

import (
	"github.com/gofiber/fiber/v2"
	"github.com/narinsak-u/backend/pkg/utils"
)

// Check returns a simple health check response. Used by container
// orchestrators (Docker healthcheck, Kubernetes liveness probe) to
// verify the server is alive and serving requests.
func Check(c *fiber.Ctx) error {
	return utils.Success(c, fiber.Map{"status": "ok"})
}
```

- [ ] **Create `backend/internal/health/health_test.go`**

```go
package health

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestCheck(t *testing.T) {
	app := fiber.New()
	app.Get("/api/health", Check)

	resp, _ := app.Test(httptest.NewRequest("GET", "/api/health", nil), 1000)
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode)

	var body map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&body)
	assert.True(t, body["success"].(bool))
	data := body["data"].(map[string]interface{})
	assert.Equal(t, "ok", data["status"])
}
```

- [ ] **Register route in `main.go`**

Add the import and route. At the top of `main.go` imports:
```go
"github.com/narinsak-u/backend/internal/health"
```

Register the route right after `api := app.Group("/api")`:
```go
api.Get("/health", health.Check)
```

- [ ] **Run tests and build**

Run: `go test ./internal/health/ -v && go build ./...`
Expected: `TestCheck` passes, build succeeds.

- [ ] **Commit**

```bash
git add backend/internal/health/ backend/cmd/server/main.go
git commit -m "feat: add /api/health endpoint for container liveness probes"
```

---

### Task 6: Fix documentation mismatches

**Files:**
- Modify: `README.md`
- Modify: `backend/docs/WORKFLOW.md`

- [ ] **Fix README.md line 169 — analytics overview auth requirement**

Change:
```
| GET | `/api/analytics/overview` | JWT | Dashboard aggregate stats |
```
To:
```
| GET | `/api/analytics/overview` | No | Dashboard aggregate stats |
```

- [ ] **Add the analytics timeseries row to the README API table**

Insert after the overview row:
```
| GET | `/api/analytics/timeseries` | No | Shipment trends by month and day-of-week |
```

- [ ] **Verify WORKFLOW.md is accurate**

Confirm the WORKFLOW.md analytics sections already match the implementation (they do — overview and timeseries are both documented correctly from the May 31 session).

- [ ] **Commit**

```bash
git add README.md
git commit -m "docs: fix analytics endpoint auth status in README, add timeseries row"
```

---

### Task 7: Split large OrderFormView bundle

**Files:**
- Modify: `frontend/src/views/OrderFormView.vue`

**Context:** Vite build reports `OrderFormView-DgS33PmN.js` at 1,039 kB (uncompressed) / 116 kB gzipped. This is because Vite hoists shared dependencies (reka-ui, cmdk-vue, lucide-vue-next) into this chunk. The lightweight fix is to add a `build.rollupOptions.output.manualChunks` config in `vite.config.ts` to separate vendor code from app code.

- [ ] **Add manualChunks to `vite.config.ts`**

```typescript
build: {
  rollupOptions: {
    output: {
      manualChunks: {
        vendor: [
          "vue",
          "vue-router",
          "pinia",
          "@tanstack/vue-query",
          "lucide-vue-next",
          "vue-sonner",
        ],
        ui: [
          "radix-vue",
          "reka-ui",
          "cmdk-vue",
          "vaul-vue",
          "class-variance-authority",
        ],
        maps: ["leaflet"],
      },
    },
  },
},
```

Add this inside the existing `defineConfig({})` call alongside `plugins` and `resolve`.

- [ ] **Build to verify**

Run: `npm run build`
Expected: Build passes. Main chunks are now separated into `vendor.*.js`, `ui.*.js`, `maps.*.js`. The `OrderFormView` chunk should drop below 500 kB uncompressed.

- [ ] **Commit**

```bash
git add frontend/vite.config.ts
git commit -m "perf: split vendor/ui/maps chunks in Vite build config"
```

---

## Execution

**Plan complete and saved to `docs/superpowers/plans/2026-05-31-pre-deployment-hardening.md`. Two execution options:**

**1. Subagent-Driven (recommended)** — I dispatch a fresh subagent per task, review between tasks, fast iteration

**2. Inline Execution** — Execute tasks in this session using executing-plans, batch execution with checkpoints

**Which approach?**
