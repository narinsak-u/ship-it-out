# Use Tailwind, shadcn/ui, and Zustand for the frontend UI

Style all Next.js routes with Tailwind CSS v4 and use shadcn/ui primitives for shared interactive controls. Use Recharts through shadcn chart components for analytics. Keep cross-route client auth state in an in-memory Zustand store; the HTTP-only session cookie and Go API remain authoritative, with no token persistence. Keep Leaflet vendor CSS and the map-specific marker rules as the only global style exception. Preserve existing routes, API contracts, shipment workflows, and authorization behavior.
