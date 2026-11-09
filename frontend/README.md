# ship-it-out Frontend

Next.js App Router shipment dashboard, deployed on Vercel. Go API remains a separate deployment; the browser uses same-origin `/api` routes proxied by Next.js.

## Quick start

```bash
cd frontend
bun install
cp .env.example .env.local
# Set API_BASE_URL to the Go API origin, without /api. Set
# NEXT_PUBLIC_CARTO_API_KEY for dark raster map tiles.
bun run dev
```

Next.js runs at `http://localhost:3000`; local Go API defaults to `http://localhost:8080`.

## Scripts

| Command | Purpose |
| --- | --- |
| `bun run dev` | Start Next.js development server |
| `bun run build` | Build production Next.js app |
| `bun run start` | Start production Next.js app |
| `bun run lint` | Run ESLint |
| `bun run test:e2e` | Run supported Next.js Playwright coverage |

## Supported app and security

Next.js is the sole frontend. Public shipment and analytics pages render on the server; interactive navigation and management use client components. Browser calls, including login and writes, go to same-origin `/api/*`. `app/api/[...path]/route.ts` forwards requests and cookies to the separately deployed Go API using server-only `API_BASE_URL`. Do not use a browser API origin variable or expose the Go endpoint as a public environment variable. Go remains authoritative for authentication, roles, and protected writes.

The login screen visibly displays the read-only demo credentials: `demo@shipitout.example` / `demo-read-only-2026`. No guest mode or guest-stored authentication exists. Demo account is customer-role and cannot write.

For production cutover, configure Vercel `API_BASE_URL` to the Go API origin and keep Go API credentials and `OPENCAGE_API_KEY` server-side. Follow `docs/adr/0002-staff-only-writes.md`: retain existing accounts and shipment data, demote pre-existing admins during the one-time cutover, and permit writes only for the newly seeded admin and staff. Do not expose secrets through `NEXT_PUBLIC_*`.

## Environment

| Variable | Required | Description |
| --- | --- | --- |
| `API_BASE_URL` | Yes | Go API origin, no `/api` suffix; server-only |
| `OPENCAGE_API_KEY` | For geocoding | Server-side OpenCage key |
| `NEXT_PUBLIC_CARTO_API_KEY` | Yes, for map tiles | Public CARTO raster tile key; required by CARTO's raster basemap service |

## Routes and coverage

Supported routes: `/`, `/orders`, `/orders/create`, `/orders/[orderId]`, `/orders/[orderId]/edit`, `/carriers`. Playwright tests in `tests/e2e/next-*.spec.ts` cover navigation, public SSR, auth cookie restore, staff writes, demo write denial, shipment management, hubs, and analytics. Tests start a local Next server and Go API-compatible fixture; live deployment smoke coverage must additionally exercise the deployed Next and Go origins.

See `docs/OVERVIEW.md`, root `docs/adr/0001-nextjs-frontend-api-origin.md`, and root `docs/adr/0003-tailwind-shadcn-zustand.md` for frontend architecture and UI decisions.
