# Serve the Next.js frontend and Go API through one browser origin

Deploy the Next.js App Router frontend on Vercel and the Go API separately. Proxy browser `/api` requests through the Next.js origin to preserve HTTP-only cookie sessions without cross-origin browser requests; server-render public data and keep interactive UI client-side. This retains the Go API contract and existing `/orders` URLs while requiring a Next.js server runtime instead of static hosting.
