# Frontend architecture

Next.js App Router is the sole supported frontend and runs on Vercel. The Go API is deployed separately. Next server components fetch public data from `API_BASE_URL`; browser components call the same-origin `/api/*` proxy in `app/api/[...path]/route.ts`. The proxy forwards methods, query strings, request bodies, and cookies without changing the Go API contract.

## Application structure

- `app/`: routes, layouts, and server-only route handlers.
- `src/features/auth/`: client auth provider and sign-in controls.
- `src/features/shipments/`: public homepage, server-rendered data helpers, manifest, shipment forms, and tracking search.
- `src/features/operations/`: carrier operations, hubs, deliveries, and maps.
- `src/features/analytics/`: API-backed analytics.
- `src/lib/`: shared shipment and hub types, API response mappers, and Next server API helpers.
- `tests/e2e/next-*.spec.ts`: supported Playwright browser coverage.

## Request and security model

`API_BASE_URL` is required and server-only. Configure it as the Go API origin without an `/api` suffix. Browser code must request `/api/...` on the Next.js origin; Next forwards to Go. This keeps the HTTP-only `jwt` cookie first-party in the browser and avoids cross-origin browser API requests. Go enforces customer/staff roles and remains authoritative for all writes. Public pages do not require login; no guest mode or client-persisted auth exists.

Login UI visibly shows demo credentials `demo@shipitout.example` / `demo-read-only-2026`. This account is read-only; protected writes require staff. Geocoding is server-side and requires `OPENCAGE_API_KEY`. Never expose API credentials through `NEXT_PUBLIC_*`.

For production cutover, retain existing accounts and shipments, demote existing admin accounts once, and grant write access only to the newly seeded admin and staff. See `docs/adr/0002-staff-only-writes.md`.

## Rendering and verification

Public shipment/analytics data uses server-side fetch with `cache: "no-store"`. Interactive filtering, pagination, forms, auth, and operations run in client components. Playwright tests cover public SSR, navigation, cookie login and restore, staff writes, and demo denial against a local Go-compatible fixture. `bun run build` builds the production Next.js application; `bun run test:e2e` runs the supported browser suite. Deployment smoke must use the live Next.js and Go deployments.
