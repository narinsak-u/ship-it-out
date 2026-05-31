# Backend Workflows

Describes the request/response flow for every major feature in the API.

---

## 1. Authentication Workflow

Handles user registration, login, session verification, and logout using JWT tokens stored in HTTP-only cookies.

### 1.1 Register (`POST /api/auth/register`)

**Purpose:** Create a new user account and immediately log them in.

**Flow:**
```
Client → POST /api/auth/register { name, email, password, role? }
          │
          ├─ RateLimitAuth middleware runs first
          │    ├─ In-memory sliding-window per IP: 5 req / 60s
          │    └─ If exceeded → 429 "too many requests"
          │
          ├─ 1. Parse JSON body into RegisterRequest
          ├─ 2. Validate required fields (name, email, password)
          ├─ 3. Hash password with bcrypt
          ├─ 4. Default role to "customer" if omitted
          ├─ 5. Insert User into Postgres
         │     └─ If email exists → 409 "email already registered"
         ├─ 6. Build JWT with claims: { user_id, role, exp (24h) }
         ├─ 7. Sign JWT with HS256 using JWTSecret
         ├─ 8. Set HTTP-only cookie "jwt" (SameSite=Lax, MaxAge=86400)
         └─ 9. Return { user } → 200
```

**Cookie config:** HTTP-only (not readable by JS), SameSite=Lax (sent on same-site navigations), 24h expiry.

### 1.2 Login (`POST /api/auth/login`)

**Purpose:** Authenticate an existing user and return a JWT.

**Flow:**
```
Client → POST /api/auth/login { email, password }
          │
          ├─ RateLimitAuth middleware runs first
          │    ├─ 5 req / 60s per IP sliding window
          │    └─ If exceeded → 429 "too many requests"
          │
          ├─ 1. Parse JSON body into LoginRequest
          ├─ 2. Look up user by email in Postgres
         │     └─ If not found → 401 "invalid email or password"
         ├─ 3. Compare password against stored bcrypt hash
         │     └─ If mismatch → 401 "invalid email or password"
         ├─ 4. Build + sign JWT with { user_id, role, exp }
         ├─ 5. Set "jwt" cookie
         └─ 6. Return { token, user } → 200
```

**Note:** The token string is returned in the JSON body too, for mobile/CLI clients that can't use cookies.

### 1.3 Me (`GET /api/auth/me`)

**Purpose:** Return the profile of the currently authenticated user.

**Flow:**
```
Client → GET /api/auth/me (with "jwt" cookie or Authorization header)
         │
         ├─ AuthRequired middleware runs first
         │    ├─ Extract token from "Authorization: Bearer <token>" or "jwt" cookie
         │    ├─ Parse + verify JWT signature
         │    │   └─ If invalid/expired → 401
         │    └─ Store user_id + role in c.Locals
         │
         ├─ 1. Read user_id from c.Locals ("user_id")
         ├─ 2. Query User by ID from Postgres
         │     └─ If not found → 404
         └─ 3. Return { id, name, email, role, created_at } → 200
```

### 1.4 Logout (`POST /api/auth/logout`)

**Purpose:** Clear the auth cookie so the browser no longer sends credentials.

**Flow:**
```
Client → POST /api/auth/logout
         │
         ├─ 1. Overwrite "jwt" cookie with empty value
         └─ 2. Set MaxAge=0 (browser deletes it immediately)
```

### 1.5 AuthRequired Middleware

Applied to protected route groups. Checks every incoming request for a valid JWT.

**Flow:**
```
Request arrives at protected route
         │
         ├─ 1. Read "Authorization" header
         │     ├─ If present and starts with "Bearer " → extract token
         │     └─ If absent → fall back to "jwt" cookie
         ├─ 2. If no token found → 401 "missing or invalid token"
         ├─ 3. Parse token with jwt.Parse() using JWTSecret as HMAC key
         │     └─ If invalid signature or expired → 401 "invalid or expired token"
         ├─ 4. Extract user_id (as uint) and role from token claims
         ├─ 5. Store in c.Locals("user_id") and c.Locals("role")
         └─ 6. Call c.Next() → route handler executes
```

---

## 2. Shipment Workflows

### 2.1 List (`GET /api/shipments`)

**Purpose:** Return shipments with pagination, search, and status filtering. Public (no auth required).

**Query params:**

| Param | Type | Default | Description |
|-------|------|---------|-------------|
| `page` | int | 1 | Page number |
| `limit` | int | 10 | Items per page (use `-1` for all) |
| `search` | string | "" | ILIKE search on order_id, tracking_number, customer_name, destination |
| `status` | string | "" | Filter by status (e.g. "in_transit") |
| `exclude_status` | string | "" | Exclude by status (e.g. "delivered") |

**Flow:**
```
Client → GET /api/shipments?page=1&limit=10&search=bangkok&status=in_transit
         │
         ├─ 1. Parse query params (page defaults 1, limit defaults 10)
         ├─ 2. Build GORM query with optional WHERE filters
         │     ├─ status filter: WHERE status = ?
         │     ├─ exclude_status: WHERE status != ?
         │     └─ search: WHERE (order_id ILIKE ? OR tracking_number ILIKE ? OR ...)
         ├─ 3. Count total matching rows
         ├─ 4. Apply OFFSET = (page-1) * limit, LIMIT = limit
         ├─ 5. ORDER BY created_at DESC
         ├─ 6. AfterFind fires → reconstructs Coords for each shipment
         └─ 7. Return { data: [...], pagination: { page, limit, total, totalPages } } → 200
```

### 2.2 GetByID (`GET /api/shipments/:orderId`)

**Purpose:** Fetch a single shipment by its OrderID (e.g. `ORD-10245`). Public.

**Flow:**
```
Client → GET /api/shipments/ORD-10245
         │
         ├─ 1. Extract orderId from URL param
         ├─ 2. Query Shipment WHERE order_id = ?
         │     └─ If not found → 404 "shipment not found"
         ├─ 3. AfterFind fires → reconstructs Coords
         └─ 4. Return Shipment → 200
```

### 2.3 Create (`POST /api/shipments`)

**Purpose:** Register a new shipment and generate a tracking number.

**Request body:**
```json
{
  "customer": { "name": "...", "zipcode": "...", "subDistrict": "...", "district": "...", "province": "...", "coords": { "lat": 0, "lng": 0 } },
  "receiver": { "name": "...", "zipcode": "...", "subDistrict": "...", "district": "...", "province": "...", "coords": { "lat": 0, "lng": 0 } },
  "carrier": "Thun-u-der Express",
  "weight": "12.4 kg",
  "items": 3
}
```

**Flow:**
```
Client → POST /api/shipments (JWT required)
         │
         ├─ 1. Parse JSON into CreateRequest
         ├─ 2. Build Shipment model:
         │     ├─ OrderID = "ORD-" + max suffix + 1 (e.g. "ORD-10261")
         │     ├─ TrackingNumber = "TH" + year + 5-digit ms hash
         │     ├─ Customer/Receiver from request body
         │     ├─ Origin = "subDistrict, district, province" (customer)
         │     ├─ Destination = "subDistrict, district, province" (receiver)
         │     ├─ CurrentCoords = customer's coords
         │     ├─ Status = "pending"
         │     ├─ EstimatedDelivery = now + 72h
         │     └─ Progress = 0
         ├─ 3. GORM BeforeSave hook → copies Coords to flat columns
         ├─ 4. Insert Shipment into Postgres
         │     └─ If error → 500 "failed to create shipment"
         ├─ 5. Create initial tracking event:
         │     ├─ Status = "Label Created"
         │     ├─ Description = "Awaiting pickup."
         │     └─ Location = customer address + coords
         └─ 6. Return Shipment → 200
```

### 2.4 Update (`PUT /api/shipments/:orderId`)

**Purpose:** Partially update an existing shipment's fields.

**Request body (all fields optional):**
```json
{
  "customer": { ... },
  "receiver": { ... },
  "carrier": "Pacific Freight",
  "weight": "15.0 kg",
  "items": 5,
  "estimatedDelivery": "2026-06-01T00:00:00Z"
}
```

**Flow:**
```
Client → PUT /api/shipments/ORD-10245 (JWT required)
         │
         ├─ 1. Extract orderId from URL param
         ├─ 2. Parse JSON body into UpdateRequest (all pointer fields)
         ├─ 3. Find existing Shipment by OrderID
         │     └─ If not found → 404 "shipment not found"
         ├─ 4. Update only provided fields:
         │     ├─ Customer → updates Customer + Origin + CurrentCoords
         │     ├─ Receiver → updates Receiver + Destination
         │     ├─ Carrier, Weight, Items, EstimatedDelivery → direct update
         ├─ 5. BeforeSave fires → syncs coords
         ├─ 6. Save shipment
         └─ 7. Return updated Shipment → 200
```

### 2.5 UpdateStatus (`PATCH /api/shipments/:orderId/status`)

**Purpose:** Change a shipment's status and log a tracking event with context-aware location. If a hubId is provided and the hub is found, sets the shipment's currentCoords to the hub's location.

**Request body:**
```json
{
  "status": "in_transit",
  "hubId": "HUB-003"
}
```

**Flow:**
```
Client → PATCH /api/shipments/ORD-10245/status (JWT required)
         │
         ├─ 1. Extract orderId from URL param
         ├─ 2. Parse JSON body { status, hubId? }
         ├─ 3. Find Shipment by OrderID → 404 if missing
         ├─ 4. Update shipment.Status
         ├─ 5. If hubId provided:
         │     ├─ Look up Hub by ID in Postgres
         │     └─ If found:
         │          ├─ Set shipment.HubID = body.hubId
         │          └─ Set shipment.CurrentCoords = hub.Lat/Lng
         ├─ 6. Save shipment (BeforeSave fires → syncs coords)
         ├─ 7. Build ShipmentEvent via statusToEvent():
         │     ├─ Status changed to a hub-based status (departed, in_transit,
         │       out_for_delivery, delayed):
         │       Location = hub's name, address, lat, lng (if hub provided)
         │       Falls back to shipment's current coords / origin / destination
         │     ├─ Status = "pending" → location = customer address
         │     ├─ Status = "picked_up" → location = customer address
         │     ├─ Status = "delivered" → location = receiver address
         │     └─ Each status maps to a human-readable label + description
         ├─ 8. Insert event (ShipmentID = shipment.ID)
         └─ 9. Return updated Shipment → 200
```

**Status-to-event mapping:**

| Status | Event label | Description | Location source |
|--------|-------------|-------------|-----------------|
| `pending` | "Label Created" | "Awaiting pickup." | Customer address |
| `picked_up` | "Picked Up" | "Parcel collected from sender." | Customer address |
| `departed` | "Departed" | "In transit to hub." | Hub or origin |
| `in_transit` | "In Transit" | "Transit to next hub." | Hub or destination |
| `out_for_delivery` | "Out for Delivery" | "Out for delivery." | Hub or destination |
| `delivered` | "Delivered" | "Delivered to recipient." | Receiver address |
| `delayed` | "Delayed" | "Unexpected issue encountered." | Hub or destination |

### 2.6 Delete (`DELETE /api/shipments/:orderId`)

**Purpose:** Remove a shipment and all its tracking events from the database.

**Flow:**
```
Client → DELETE /api/shipments/ORD-10245 (JWT required)
         │
         ├─ 1. Extract orderId from URL param
         ├─ 2. Find Shipment by OrderID → 404 if missing
         ├─ 3. Delete all ShipmentEvents WHERE shipment_id = shipment.ID
         ├─ 4. Delete the Shipment itself
         └─ 5. Return { message: "shipment deleted" } → 200
```

---

## 3. Public Tracking Workflow

### 3.1 Track (`GET /api/track/:trackingNumber`)

**Purpose:** Allow anyone (no auth) to look up a shipment and its event history by tracking number.

**Flow:**
```
Client → GET /api/track/TH202612345
         │
         ├─ 1. Extract trackingNumber from URL param
         ├─ 2. Query Shipment by tracking_number column
         │     └─ If not found → 404
         ├─ 3. AfterFind fires → reconstructs Coords
         ├─ 4. Query ShipmentEvents WHERE shipment_id = shipment.ID
         │     ordered by created_at asc
         └─ 5. Return { shipment, events } → 200
```

---

## 4. Hub Workflows

### 4.1 List (`GET /api/hubs`)

**Purpose:** Return all hubs. Public (no auth required).

**Flow:**
```
Client → GET /api/hubs
         ├─ Query all Hub rows
         ├─ AfterFind fires → reconstruct Coords for each
         └─ Return [...] → 200
```

### 4.2 GetByID (`GET /api/hubs/:id`)

**Purpose:** Fetch a single hub by its string ID (e.g. HUB-001). Public.

**Flow:**
```
Client → GET /api/hubs/HUB-001
         ├─ 1. Query Hub by ID (string PK, not uint)
         │     └─ If not found → 404
         ├─ 2. AfterFind fires → reconstruct Coords
         └─ 3. Return Hub → 200
```

### 4.3 Create (`POST /api/hubs`)

**Purpose:** Add a new logistics hub. Auto-generates "HUB-NNN" ID if not provided.

**Request body:**
```json
{
  "name": "Rayong Hub",
  "carrierId": "CAR-001",
  "address": "Industrial Estate, Rayong",
  "coords": { "lat": 12.6814, "lng": 101.2817 },
  "capacity": 5000,
  "currentUtilization": 1200,
  "status": "active"
}
```

**Flow:**
```
Client → POST /api/hubs (JWT required)
         │
         ├─ 1. Parse JSON body into Hub
         ├─ 2. If ID not provided, generate "HUB-NNN" (max suffix + 1)
         ├─ 3. BeforeSave fires → Coords → Lat/Lng
         ├─ 4. Insert into Postgres
         │     └─ If error → 500
         └─ 5. Return Hub → 200
```

### 4.4 Update (`PUT /api/hubs/:id`)

**Purpose:** Modify an existing hub's fields.

**Flow:**
```
Client → PUT /api/hubs/HUB-001 (JWT required)
         │
         ├─ 1. Find existing Hub by ID → 404 if missing
         ├─ 2. Parse JSON body into existing Hub (BodyParser overwrites fields)
         ├─ 3. Re-set Hub.ID (BodyParser clears the PK field)
         ├─ 4. Save (BeforeSave fires → syncs coords)
         └─ 5. Return updated Hub → 200
```

### 4.5 Delete (`DELETE /api/hubs/:id`)

**Purpose:** Remove a hub from the database.

**Flow:**
```
Client → DELETE /api/hubs/HUB-001 (JWT required)
         ├─ 1. Delete Hub by ID
         │     └─ If error → 500
         └─ 2. Return { message: "hub deleted" } → 200
```

---

## 5. Analytics Workflow

### 5.1 Overview (`GET /api/analytics/overview`)

**Purpose:** Return aggregate statistics for the dashboard.

**Flow:**
```
Client → GET /api/analytics/overview
          │
          ├─ 1. Count total shipments → total
          ├─ 2. Count shipments NOT in (DELIVERED, RETURNED) → active
          ├─ 3. Count shipments with status = "DELIVERED" → delivered
          ├─ 4. Group all shipments by status, count each → by_status
          ├─ 5. Group shipments by province, map to 6 Thai regions → by_region
          └─ 6. Return { total, active, delivered, by_status, by_region } → 200
```

**Notes:**
- `by_status` is a GROUP BY via GORM's `Select().Group().Scan()`.
- `by_region` maps each shipment's destination province to one of 6 Thai regions (Central, East, North, West, North-east, South) via `internal/data/regions.go`.

**Response shape:**
```json
{
  "success": true,
  "data": {
    "total": 12,
    "active": 8,
    "delivered": 3,
    "by_status": [
      { "status": "in_transit", "count": 4 },
      { "status": "pending", "count": 2 },
      ...
    ],
    "by_region": [
      { "name": "Central", "total": 5 },
      { "name": "East", "total": 4 },
      ...
    ]
  }
}
```

### 5.2 TimeSeries (`GET /api/analytics/timeseries`)

**Purpose:** Return shipment creation trends grouped by month and day of week for dashboard charts.

**Flow:**
```
Client → GET /api/analytics/timeseries
         │
         ├─ 1. Group shipments by month (created_at truncated to YYYY-MM),
         │     count each → by_month
         ├─ 2. Group shipments by day-of-week (0=Sunday..6=Saturday),
         │     count each → by_day_of_week
         └─ 3. Return { by_month, by_day_of_week } → 200
```

**Response shape:**
```json
{
  "success": true,
  "data": {
    "by_month": [
      { "month": "2026-05", "count": 5 },
      { "month": "2026-06", "count": 7 }
    ],
    "by_day_of_week": [
      { "day_of_week": 1, "count": 3 },
      { "day_of_week": 3, "count": 5 }
    ]
  }
}
```

---

## 6. Rate Limiting Middleware

**Purpose:** Prevent brute-force attacks on authentication endpoints using an in-memory sliding-window rate limiter keyed by client IP.

**Configuration:**

| Setting | Value |
|---------|-------|
| Limit | 5 requests |
| Window | 60 seconds (sliding) |
| Applied to | `POST /api/auth/register`, `POST /api/auth/login` |
| Storage | In-memory map (single instance; replace with Redis for horizontal scaling) |

**Flow:**
```
Request arrives at login or register
         │
         ├─ 1. Extract client IP via c.IP()
         ├─ 2. Prune timestamps older than 60s for this IP
         ├─ 3. If remaining count >= 5 → 429 "too many requests"
         ├─ 4. Otherwise → record timestamp, call c.Next()
         └─ 5. Background cleanup goroutine evicts stale IPs every 60s
```

**Note:** Single-instance in-process limiter. For horizontally-scaled deployments, replace with a shared Redis-based implementation.

---

## 7. Seed Data Workflow

**Purpose:** Populate the database with demo data on first startup so developers can test immediately.

### Entry Point

`cmd/server/main.go` after `AutoMigrate`:

```go
seed.SeedHubs(database.DB)
seed.SeedShipments(database.DB)
```

### 7.1 SeedHubs

Inserts 6 hubs in Eastern Thailand (idempotent — skips if rows exist):

| ID | Hub | Lat | Lng | Status |
|----|-----|-----|-----|--------|
| HUB-001 | Laem Chabang Port Hub | 13.0833 | 100.8833 | active |
| HUB-002 | Pattaya Hub | 12.9236 | 100.8825 | active |
| HUB-003 | Rayong Hub | 12.6814 | 101.2817 | active |
| HUB-004 | Chanthaburi Hub | 12.6096 | 102.1041 | active |
| HUB-005 | Chachoengsao Hub | 13.6883 | 101.0719 | active |
| HUB-006 | Trat Hub | 12.2417 | 102.5125 | maintenance |

### 7.2 SeedShipments

Inserts 12 shipments (ORD-10245 through ORD-10260) with full contact info for customer/receiver (Thai names and addresses), tracking numbers in TH2026xxxxx format, and 1-6 tracking events each spanning a plausible status lifecycle. Idempotent — skips if rows exist.
