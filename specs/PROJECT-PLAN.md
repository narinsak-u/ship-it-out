
# Real-Time Shipment Tracking System — Full Project Plan

A portfolio-grade logistics platform inspired by:
- DHL Express
- FedEx
- Flash Express Thailand

Goal:
Build a modern realtime shipment platform using:
- Vue.js
- Go
- WebSockets
- Event-driven architecture
- Realtime updates
- Analytics dashboards

---

# 1. Core Concept

Users can:
- Create shipments
- Track parcels
- View realtime delivery updates
- Estimate delivery time
- Manage logistics operations

Admins/operators can:
- Assign drivers
- Update shipment status
- Manage hubs/warehouses
- Monitor analytics
- Track active deliveries live

---

# 2. System Architecture

```text
                ┌────────────────┐
                │ Vue Frontend   │
                └───────┬────────┘
                        │ REST/WebSocket
                ┌───────▼────────┐
                │ API Gateway    │
                └───────┬────────┘
                        │
        ┌───────────────┼────────────────┐
        ▼               ▼                ▼
┌────────────┐  ┌──────────────┐  ┌──────────────┐
│ Shipment   │  │ Tracking     │  │ Analytics    │
│ Service    │  │ Service      │  │ Service      │
└────────────┘  └──────────────┘  └──────────────┘
        │               │                │
        └───────┬───────┴────────────────┘
                ▼
         ┌──────────────┐
         │ PostgreSQL   │
         └──────────────┘
                │
                ▼
         ┌──────────────┐
         │ Redis PubSub │
         └──────────────┘
```

---

# 3. Tech Stack

# Frontend

## Core
- Vue.js
- Vite
- TypeScript
- Pinia
- Vue Router

## UI
- Tailwind CSS
- shadcn/vue
- Lucide icons

## Data
- Vue Query
- Axios

## Maps
- Leaflet
or
- Google Maps API

---

# Backend

## Core
- Go
- Fiber

## Database
- PostgreSQL
- Redis

## ORM
- GORM

## Realtime
- Gorilla WebSocket

## Auth
- JWT
- Refresh tokens

## Infrastructure
- Docker
- Docker Compose
- Nginx

---

# 4. Main Features

# Phase 1 — MVP

## Authentication

### Roles
- Admin
- Operator
- Driver
- Customer

### Features
- Login/register
- JWT auth
- Protected routes
- Role permissions

---

# Shipment Management

## Create Shipment

Fields:
- sender
- receiver
- addresses
- parcel weight
- dimensions
- delivery type
- COD option

---

## Shipment Status

```text
CREATED
PICKED_UP
AT_SORTING_CENTER
IN_TRANSIT
OUT_FOR_DELIVERY
DELIVERED
FAILED_DELIVERY
RETURNED
```

---

## Tracking Number System

Example:
```text
TH20260001234
```

Generate:
- unique
- searchable
- sortable

---

# Tracking Page

Public route:
```text
/track/:trackingNumber
```

Shows:
- shipment info
- realtime status
- timeline
- estimated arrival
- parcel history

---

# Realtime Updates

Using WebSocket:
- live status changes
- ETA updates
- driver location

---

# Shipment Timeline UI

Example:
```text
[✓] Parcel Created
[✓] Picked Up
[✓] Sorting Center Bangkok
[ ] Out For Delivery
[ ] Delivered
```

---

# Phase 2 — Advanced Features

# Driver System

## Driver Dashboard

Features:
- assigned deliveries
- route navigation
- delivery completion
- failed delivery reporting

---

# Live Map Tracking

Display:
- driver location
- delivery route
- destination

Using:
- GPS updates
- WebSockets

---

# ETA Prediction

Simple formula first:
```text
distance + traffic + delivery queue
```

Later:
- AI prediction
- historical delivery data

---

# Warehouse / Hub System

Entities:
- hubs
- warehouses
- sorting centers

Shipment moves:
```text
Bangkok Hub
→ Ayutthaya Hub
→ Chiang Mai Hub
```

---

# Analytics Dashboard

Charts:
- active shipments
- delivery success rate
- avg delivery time
- failed deliveries
- busiest routes

---

# Phase 3 — Portfolio “WOW” Features

# Event-Driven Architecture

Use Redis Pub/Sub.

Example:
```text
Shipment Updated
      ↓
Event Published
      ↓
Tracking Service
      ↓
Analytics Service
```

---

# Background Workers

Workers process:
- ETA calculations
- cleanup jobs
- report generation

---

# Optimistic UI

Frontend updates instantly:
```text
User changes shipment status
→ UI updates immediately
→ backend confirms later
```

---

# Offline Support

Drivers can:
- update deliveries offline
- sync later

Use:
- IndexedDB
- sync queue

---

# Barcode / QR System

Generate:
- QR labels
- barcode labels

Can scan:
- package updates
- delivery confirmations

---

# 5. Database Design

# users

```sql
id
name
email
password
role
created_at
```

---

# shipments

```sql
id
tracking_number
sender_name
receiver_name
origin_address
destination_address
weight
status
estimated_delivery
created_at
```

---

# shipment_events

```sql
id
shipment_id
status
location
description
created_at
```

This table powers:
- timelines
- history
- analytics

---

# drivers

```sql
id
user_id
vehicle_type
phone
current_lat
current_lng
```

---

# hubs

```sql
id
name
province
lat
lng
```

---

# 6. Frontend Pages

# Public

## Home
- tracking search
- shipment lookup

---

## Tracking Detail
- timeline
- shipment status
- map
- ETA

---

# Admin

## Dashboard
- analytics
- shipment overview

---

## Shipments
- create/edit/manage shipments

---

## Drivers
- driver assignment

---

## Hubs
- warehouse management

---

## Analytics
- charts/reports

---

# Driver

## Driver App Dashboard
- assigned packages
- route map
- delivery actions

---

# 7. API Design

# Auth

```http
POST /api/auth/login
POST /api/auth/register
POST /api/auth/refresh
```

---

# Shipments

```http
GET /api/shipments
POST /api/shipments
GET /api/shipments/:id
PATCH /api/shipments/:id/status
```

---

# Tracking

```http
GET /api/track/:trackingNumber
```

---

# WebSocket

```text
/ws/tracking/:trackingNumber
/ws/admin
/ws/driver
```

---

# 8. Realtime Architecture

# WebSocket Flow

```text
Driver updates status
        ↓
Backend updates DB
        ↓
Redis event published
        ↓
WebSocket server broadcasts
        ↓
Frontend updates instantly
```

---

# 9. Suggested Folder Structure

# Frontend

```text
src/
 ├── api/
 ├── components/
 ├── composables/
 ├── layouts/
 ├── pages/
 ├── stores/
 ├── types/
 └── websocket/
```

---

# Backend

```text
/cmd
/internal
   /auth
   /shipment
   /tracking
   /analytics
   /websocket
/pkg
```

---

# 10. UI/UX Ideas

# Modern Logistics UI

Style inspiration:
- Linear
- Stripe Dashboard
- DHL

Use:
- dark mode
- realtime animations
- live indicators
- status chips

---

# 11. Deployment Plan

# Frontend
- Cloudflare Pages
or
- Vercel

---

# Backend
- Railway
- Fly.io

---

# Database
- Supabase
- Neon

---
