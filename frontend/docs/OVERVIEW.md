# Frontend architecture

Next.js App Router is the sole supported frontend and runs on Vercel. The Go API is deployed separately. Next server components fetch public data from `API_BASE_URL`; browser components call the same-origin `/api/*` proxy in `src/app/api/[...path]/route.ts`. The proxy forwards methods, query strings, and request bodies, plus only `Accept`, `Authorization`, `Content-Type`, `Origin`, and the `jwt` cookie; unrelated browser headers and cookies stay at the Next boundary.

## Application structure

- `src/app/`: routes, layouts, and server-only route handlers; pages compose feature-owned UI and data.
- `src/features/auth/`: auth API operations, user types, Zustand session state, and sign-in controls.
- `src/features/shipments/`: canonical Shipment model and transport mappings, server/browser APIs, homepage data, manifest, and forms.
- `src/features/tracking/`: tracking lookup, tracking events, server API access, and route visualization.
- `src/features/hubs/`: canonical Hub model and management API, forms, and controls.
- `src/features/operations/`: active delivery monitoring, status workflow, and maps; consumes Shipments and Hubs through their feature APIs and models.
- `src/features/analytics/`: analytics API access, data mapping, and Recharts charts.
- `src/components/ui/`: business-independent shadcn/ui primitives; `src/lib/`: business-independent helpers, including generic API utilities and the server API origin helper.
- `src/app/globals.css`: Tailwind v4 theme and base layer; Leaflet vendor CSS and map marker rules are the only global style exceptions.
- `tests/e2e/next-*.spec.ts`: supported Playwright browser coverage.

## Request and security model

- `API_BASE_URL` is required and server-only. Configure it as the Go API origin without an `/api` suffix. Browser code must request `/api/...` on the Next.js origin; Next forwards to Go. This keeps the HTTP-only `jwt` cookie first-party in the browser and avoids cross-origin browser API requests. Go enforces customer/staff roles and remains authoritative for all writes. Zustand holds auth state in memory only; no guest mode or token persistence.

Login UI visibly shows demo credentials `demo@shipitout.example` / `demo-read-only-2026`. This account is read-only; protected writes require staff. Geocoding is server-side and requires `OPENCAGE_API_KEY`. Never expose API credentials through `NEXT_PUBLIC_*`.

For production cutover, retain existing accounts and shipments, demote existing admin accounts once, and grant write access only to the newly seeded admin and staff. See `docs/adr/0002-staff-only-writes.md`.

Frontend styling, component, and client auth state decisions are recorded in `docs/adr/0003-tailwind-shadcn-zustand.md`.

## Rendering and verification

Public shipment/analytics data uses server-side fetch with `cache: "no-store"`. Interactive filtering, pagination, forms, auth, and operations run in client components. Playwright tests cover public SSR, navigation, cookie login and restore, staff writes, and demo denial against a local Go-compatible fixture. `bun run build` builds the production Next.js application; `bun run test:e2e` runs the supported browser suite. Deployment smoke must use the live Next.js and Go deployments.
