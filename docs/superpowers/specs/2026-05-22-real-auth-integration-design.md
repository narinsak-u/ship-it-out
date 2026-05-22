# Real Auth Integration Design

**Date:** 2026-05-22
**Status:** Approved

## Goal

Replace the mock Pinia auth store with real backend API calls using httpOnly JWT cookies. Keep guest mode for read-only access. Keep all order/carrier data as mock (data API swap deferred).

## Architecture

JWT tokens are set as httpOnly cookies by the Go backend on login/register. The frontend calls `GET /api/auth/me` to validate sessions and fetch the current user. Guest mode is a separate flag that bypasses auth entirely — no token, no backend calls, mock data only.

```
POST /api/auth/login  ──> 200 + Set-Cookie: jwt=<token>; HttpOnly; Path=/; SameSite=Lax; Max-Age=86400
POST /api/auth/register ──> same as login
GET  /api/auth/me   <── 200 { id, name, email, role, created_at } | 401
POST /api/auth/logout ──> Set-Cookie: jwt=; Max-Age=0
```

## Backend Changes

### CORS (`internal/middleware/cors.go`)

- Change `Access-Control-Allow-Origin` from `*` to `http://localhost:5173`
- Add `Access-Control-Allow-Credentials: true`
- Keep existing methods and headers

### Auth handler (`internal/auth/handler.go`)

**Login handler (modified):**
- After JWT generation, set `Set-Cookie` header:
  - Name: `jwt`
  - Value: token string
  - `HttpOnly`
  - `Path=/`
  - `SameSite=Lax`
  - `Max-Age=86400` (matches 24h JWT TTL)
  - `Secure` omitted in dev (no HTTPS on localhost)
- Keep existing JSON response body (token + user) for backward compat

**Register handler (modified):**
- Same cookie-setting after user creation

**New `Me` handler** (`GET /api/auth/me`):
- Uses `middleware.AuthRequired()` (already exists)
- Reads `user_id` from `c.Locals("user_id")`
- Looks up user in DB by ID (select id, name, email, role, created_at)
- Returns `utils.Success(c, user)`

**New `Logout` handler** (`POST /api/auth/logout`):
- Sets `Set-Cookie: jwt=; Path=/; HttpOnly; Max-Age=0`
- Returns `utils.Success(c, fiber.Map{"message": "logged out"})`

### Server routes (`cmd/server/main.go`)

Add to the `api.Group("/auth")`:
```go
authGroup.Get("/me", middleware.AuthRequired(), auth.Me)
authGroup.Post("/logout", auth.Logout)
```

## Frontend Changes

### Auth store (`src/stores/auth.ts`) — rewrite

```
State: user (User | null), loading (boolean), error (string), isGuest (boolean)
```

- `user` ref is the single source of truth for auth state. `isAuthenticated` is a computed getter: `user !== null`.
- On store init (module-level IIFE or explicit `init()`), call `GET /api/auth/me` with `credentials: 'include'`. Set `loading = true` during fetch.
- `login(email, password)` → `POST /api/auth/login` → sets cookie, then calls `/me` to hydrate user. Returns `null` on success, error string on failure.
- `signup(name, email, password)` → `POST /api/auth/register` → same flow.
- `logout()` → `POST /api/auth/logout` → clears user ref.
- `isGuest` persisted in `sessionStorage` (lost on tab close). Guest mode skips all auth checks — no `/me`, no cookie. Read-only mock data.
- On 401 from any auth call → clear user. On network error → set `error` ref with user-friendly message.

**User type** (defined inline or in store):
```ts
interface AuthUser {
  id: number;
  name: string;
  email: string;
  role: string;
  created_at: string;
}
```

### API client (`src/lib/api/client.ts`) — new

```ts
const BASE = 'http://localhost:8080/api';

async function request(path: string, options?: RequestInit) {
  try {
    const res = await fetch(`${BASE}${path}`, {
      credentials: 'include',
      headers: { 'Content-Type': 'application/json', ...options?.headers },
      ...options,
    });
    const json = await res.json();
    if (!res.ok) return { error: json.error || `Request failed (${res.status})` };
    return { data: json.data };
  } catch {
    return { error: 'Network error — is the backend running?' };
  }
}

export const api = {
  get: (path: string) => request(path),
  post: (path: string, body?: unknown) => request(path, { method: 'POST', body: body ? JSON.stringify(body) : undefined }),
  del: (path: string) => request(path, { method: 'DELETE' }),
};
```

### AuthModal (`src/components/AuthModal.vue`)

- `handleLogin` / `handleSignup` → call store methods (now async), set `error` from return value
- Loading state on submit button while awaiting API
- Guest button stays
- Close on backdrop click + close button stays

### SiteHeader (`src/components/SiteHeader.vue`)

- During `store.loading`, show nothing or skeleton for auth section
- When `store.user` is set → show "Admin ({name})" + "Sign out"
- When `store.isGuest` → show "Guest" label + "Sign in" button
- When neither → show "Sign in" button

### OrdersView (`src/views/OrdersView.vue`)

- New Order button: if `store.user` → `router.push(order-create)`. If `store.isGuest` or unauthenticated → show AuthModal
- Actions column: visible only when `store.user` is set
- No changes to mock data reading

## Error Handling

- **Token expiry**: `/me` returns 401 → user ref cleared → UI reverts to logged-out state
- **Network error**: API client returns structured error, surfaced in AuthModal and console-logged elsewhere
- **Loading state**: Auth store exposes `loading` ref, components conditionally render
- **Race on login**: Login sets cookie → calls `/me`. If `/me` fails, login considered failed

## Security

- JWT never accessible to JavaScript (httpOnly)
- SameSite=Lax prevents CSRF from external links
- CORS restricted to frontend dev origin with credentials
- Backend already validates JWT signature on every request via middleware

## Non-Goals

- No token refresh (24h TTL is sufficient)
- No order/carrier data API integration (deferred)
- No HTTPS in dev (Secure flag omitted)
- No role-based UI beyond showing "Admin" label
