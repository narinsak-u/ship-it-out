# AGENTS.md — ship-it-out

## Project structure

- `frontend/`: sole supported Next.js App Router application, deployed on Vercel.
- `backend/`: separately deployed Go/Fiber API and PostgreSQL configuration.
- `docs/adr/`: architecture and authorization decisions.

## Frontend commands (run from `frontend/`)

| Command | Purpose |
| --- | --- |
| `bun run dev` | Start Next.js development server |
| `bun run build` | Build production Next.js app |
| `bun run start` | Start production Next.js app |
| `bun run lint` | Run ESLint |
| `bun run test:e2e` | Run supported Next.js Playwright suite |

Set server-only `API_BASE_URL` to the Go API origin without `/api`. Browser API calls use same-origin `/api/*` and Next.js proxies to Go. Never expose API or geocoding secrets as public variables. Demo credentials stay visible; demo account is read-only. No guest mode.

## Backend commands (run from `backend/`)

| Command | Purpose |
| --- | --- |
| `go run ./cmd/server` | Start API server (needs Postgres) |
| `docker compose up -d` | Start Postgres 16 |
| `go build ./...` | Compile packages |
| `go vet ./...` | Static analysis |
| `go test ./...` | Run tests |
| `go mod tidy` | Clean dependencies |

Backend listens on `http://localhost:8080`; frontend listens on `http://localhost:3000`.

## Backend conventions

- Fiber handlers return `pkg/utils` response helpers.
- All API routes use `/api/` and response format `{"success":true,"data":...}` or `{"success":false,"error":"..."}`.
- JWT sessions use HTTP-only cookies; Go is authoritative for roles and protected writes.
- Explicitly handle errors; use pointer request fields for partial updates.

## Reference docs

- `frontend/README.md`: setup, environment, security, and browser coverage.
- `frontend/docs/OVERVIEW.md`: supported frontend architecture.
- `backend/docs/OVERVIEW.md`: API architecture.
- `backend/docs/WORKFLOW.md`: request and response flows.
- `docs/adr/0001-nextjs-frontend-api-origin.md`: deployment and browser API origin.
- `docs/adr/0002-staff-only-writes.md`: authorization and production cutover.

## Agent skills

### Issue tracker

Issues live in GitHub Issues. Use the `gh` CLI. See `docs/agents/issue-tracker.md`.

### Triage labels

Use the five canonical labels: needs-triage, needs-info, ready-for-agent, ready-for-human, wontfix. See `docs/agents/triage-labels.md`.

### Domain docs

Single-context: root GLOSSARY.md and docs/adr/. See `docs/agents/domain.md`.
