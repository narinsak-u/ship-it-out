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

## 2. Shipment CRUD Workflow

### 2.1 List (`GET /api/shipments`)

**Purpose:** Return all shipments from the database.

**Flow:**
```
Client → GET /api/shipments (JWT required)
         │
         ├─ 1. Query all Shipment rows via GORM Find()
         └─ 2. Return [] → 200
```

**Note:** No pagination, filtering, or sorting yet. Returns all rows.

### 2.2 Create (`POST /api/shipments`)

**Purpose:** Register a new shipment and generate a tracking number.

**Flow:**
```
Client → POST /api/shipments (JWT required)
         { customer: { name, zipcode, subDistrict, district, province, coords },
           receiver: { name, zipcode, subDistrict, district, province, coords },
           carrier, weight, items }
         │
         ├─ 1. Parse JSON into CreateRequest
         ├─ 2. Build Shipment model:
         │     ├─ TrackingNumber = "TH" + year + 5-digit ms hash
         │     ├─ Customer/Receiver from request
         │     ├─ Origin = "subDistrict, district, province" (customer)
         │     ├─ Destination = "subDistrict, district, province" (receiver)
         │     ├─ CurrentCoords = customer's coords
         │     ├─ Status = "pending"
         │     ├─ EstimatedDelivery = now + 72h
         │     └─ Progress = 0
         ├─ 3. GORM BeforeSave hook fires → copies Coords → flat _lat/_lng columns
         ├─ 4. Insert into Postgres
         │     └─ If error → 500
         └─ 5. Return Shipment → 200
```

### 2.3 GetByID (`GET /api/shipments/:id`)

**Purpose:** Fetch a single shipment by its database ID.

**Flow:**
```
Client → GET /api/shipments/42 (JWT required)
         │
         ├─ 1. Parse :id param as uint
         │     └─ If not a number → 400
         ├─ 2. Query Shipment by primary key
         │     └─ If not found → 404
         ├─ 3. GORM AfterFind hook fires → reconstructs Coords from flat columns
         └─ 4. Return Shipment → 200
```

### 2.4 UpdateStatus (`PATCH /api/shipments/:id/status`)

**Purpose:** Change a shipment's status and log a tracking event.

**Flow:**
```
Client → PATCH /api/shipments/42/status (JWT required)
         { status: "in_transit" }
         │
         ├─ 1. Parse :id param + JSON body { status }
         ├─ 2. Find Shipment by ID → 404 if missing
         ├─ 3. Update shipment.Status
         ├─ 4. Save shipment (BeforeSave fires → syncs coords)
         │
         ├─ 5. Create ShipmentEvent:
         │     ├─ ShipmentID = shipment.ID
         │     ├─ Status = body.status
         │     ├─ Location = { name: shipment.Destination,
         │     │               lat: shipment.CurrentCoords.Lat,
         │     │               lng: shipment.CurrentCoords.Lng }
         │     └─ Description = "Status updated to <status>"
         ├─ 6. Insert event into Postgres
         └─ 7. Return updated Shipment → 200
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
         ├─ 3. GORM AfterFind fires → reconstructs Coords
         ├─ 4. Query ShipmentEvents where shipment_id = shipment.ID
         │     ordered by created_at asc
         └─ 5. Return { shipment, events } → 200
```

---

## 4. Analytics Workflow

### 4.1 Overview (`GET /api/analytics/overview`)

**Purpose:** Return aggregate statistics for the dashboard.

**Flow:**
```
Client → GET /api/analytics/overview (JWT required)
         │
         ├─ 1. Count total shipments → total
         ├─ 2. Count shipments NOT in (DELIVERED, RETURNED) → active
         ├─ 3. Count shipments with status = "DELIVERED" → delivered
         ├─ 4. Group all shipments by status, count each → by_status
         └─ 5. Return { total, active, delivered, by_status } → 200
```

**Note:** The `by_status` array is a raw SQL `GROUP BY` via GORM's `Select().Group().Scan()`.

---

## 5. WebSocket Workflow

**Purpose:** Provide real-time tracking updates to connected clients.

### 5.1 Connection

```
Client → GET /ws/tracking/TH202612345 (WebSocket upgrade)
         │
         ├─ 1. Fiber checks for WebSocket upgrade headers
         │     └─ If not a WS request → 400
         ├─ 2. Determine room:
         │     ├─ /ws/tracking/:trackingNumber → room = the tracking number
         │     ├─ /ws/admin → room = "global"
         │     └─ /ws/driver → room = "global"
         ├─ 3. Create Client struct: { Room, Conn, Send (buffered chan, 256) }
         ├─ 4. Register client in DefaultHub
         │     (adds to map[Client]bool under a mutex)
         ├─ 5. Start goroutine: reads from client.Send chan → writes to WS connection
         ├─ 6. Enter read loop: blocks on conn.ReadMessage()
         │     (messages from client are currently discarded)
         └─ 7. On read error → break → Unregister client → close conn
```

### 5.2 Broadcasting

```
External trigger (not yet implemented)
         │
         └─ DefaultHub.BroadcastToRoom(room, message)
              │
              ├─ 1. Acquire read lock
              ├─ 2. Iterate all clients
              │     └─ For each client with matching Room:
              │          ├─ Try: message ← client.Send (non-blocking send to chan)
              │          └─ On full buffer: close chan + delete client (slow consumer drop)
              └─ 3. Release lock
```

**Current state:** The WebSocket infrastructure (connection upgrade, room registration, broadcasting) is fully wired. What's missing is the trigger that calls `BroadcastToRoom()` when a shipment status changes.

---

## 6. Hub Management Workflow

**Purpose:** Manage logistics hubs (warehouses/sorting centers).

### 6.1 List (`GET /api/hubs`)

```
Client → GET /api/hubs (JWT required)
         ├─ Query all Hub rows
         ├─ AfterFind fires → reconstruct Coords for each
         └─ Return [] → 200
```

### 6.2 Create (`POST /api/hubs`)

```
Client → POST /api/hubs (JWT required)
         { id, name, carrierId, address, coords, capacity, currentUtilization, status }
         │
         ├─ 1. Parse JSON body into Hub
         ├─ 2. BeforeSave fires → Coords → Lat/Lng
         ├─ 3. Insert into Postgres
         │     └─ If error → 500
         └─ 4. Return Hub → 200
```

### 6.3 GetByID (`GET /api/hubs/:id`)

```
Client → GET /api/hubs/HUB-001 (JWT required)
         ├─ 1. Query Hub by string ID (not uint)
         │     └─ If not found → 404
         ├─ 2. AfterFind fires → reconstruct Coords
         └─ 3. Return Hub → 200
```

### 6.4 Update (`PUT /api/hubs/:id`)

```
Client → PUT /api/hubs/HUB-001 (JWT required)
         { name, address, coords, ... }
         │
         ├─ 1. Find existing Hub by ID → 404 if missing
         ├─ 2. Parse JSON body into existing Hub (BodyParser overwrites fields)
         ├─ 3. Re-set ID (BodyParser clears the PK)
         ├─ 4. Save (BeforeSave fires → syncs coords)
         └─ 5. Return updated Hub → 200
```

### 6.5 Delete (`DELETE /api/hubs/:id`)

```
Client → DELETE /api/hubs/HUB-001 (JWT required)
         ├─ 1. Delete Hub by ID
         │     └─ If error → 500
         └─ 2. Return { message: "hub deleted" } → 200
```

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

```
SeedHubs(db)
  ├─ 1. Count Hub rows
  ├─ 2. If count > 0 → return (already seeded, skip)
  └─ 3. Insert 6 hubs one by one:
       ├─ Laem Chabang Port Hub (HUB-001, ชลบุรี, active)
       ├─ Pattaya Hub (HUB-002, ชลบุรี, active)
       ├─ Rayong Hub (HUB-003, ระยอง, active)
       ├─ Chanthaburi Hub (HUB-004, จันทบุรี, active)
       ├─ Chachoengsao Hub (HUB-005, ฉะเชิงเทรา, active)
       └─ Trat Hub (HUB-006, ตราด, maintenance)
```

### 7.2 SeedShipments

```
SeedShipments(db)
  ├─ 1. Count Shipment rows
  ├─ 2. If count > 0 → return (already seeded, skip)
  └─ 3. Insert 3 shipments + their events:
       │
       ├─ Shipment 1: TRK-9F2A-44B1
       │   Customer: สมชาย วงศ์เจริญ (แหลมฉบัง, ศรีราชา, ชลบุรี)
       │   Receiver: มาลี ทองดี (จันทนิมิต, เมือง, จันทบุรี)
       │   Status: in_transit | Carrier: Thun-u-der Express
       │   Events: Picked up (Laem Chabang Port Hub) → Departed (Laem Chabang) → In transit (Ban Bueng)
       │
       ├─ Shipment 2: TRK-FF02-1188
       │   Customer: วิมล ศรีสุวรรณ (หน้าเมือง, เมือง, ฉะเชิงเทรา)
       │   Receiver: กิตติพงศ์ แก้ววิเศษ (หนองปรือ, บางละมุง, ชลบุรี)
       │   Status: pending | Carrier: Thun-u-der Express
       │   Events: Label created (Chachoengsao Hub)
       │
       └─ Shipment 3: TRK-5E73-220B
           Customer: วิชัย สมบูรณ์ (ท่าประดู่, เมือง, ระยอง)
           Receiver: ประภาสิริ วัฒนา (บางพระ, เมือง, ตราด)
           Status: in_transit | Carrier: Thun-u-der Express
           Events: Picked up (Rayong Hub) → Departed (Rayong) → In transit (Klaeng)
```

Each shipment is created first via `db.Create()`, then its events are inserted with the auto-assigned `ShipmentID` foreign key.
