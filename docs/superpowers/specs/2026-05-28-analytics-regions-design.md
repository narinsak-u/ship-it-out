# Analytics: Shipments by Region

## Goal

Replace the hardcoded client-side "Shipments by Region" section in the AnalyticsPanel with data from the backend `/analytics/overview` API endpoint. The entire AnalyticsPanel should switch to consuming API data instead of the local `orders` array.

## Backend Changes

### Province-to-Region Mapping

A Go map in `backend/internal/analytics/handler.go` mapping Thai provinces to regions:

| Region | Provinces |
|--------|-----------|
| Central | Bangkok, Nonthaburi, Pathum Thani, Samut Prakan, Phra Nakhon Si Ayutthaya, Nakhon Pathom, Samut Sakhon, Samut Songkhram, Sing Buri, Ang Thong, Lopburi, Saraburi, Chai Nat, Suphan Buri |
| East | Chonburi, Rayong, Chachoengsao, Chanthaburi, Trat, Prachinburi, Sa Kaeo, Nakhon Nayok |
| North | Chiang Mai, Chiang Rai, Lampang, Lamphun, Phrae, Nan, Phayao, Mae Hong Son, Uttaradit |
| West | Kanchanaburi, Phetchaburi, Prachuap Khiri Khan, Ratchaburi, Tak |
| North-east | Khon Kaen, Nakhon Ratchasima, Udon Thani, Ubon Ratchathani, Buriram, Surin, Sisaket, Chaiyaphum, Loei, Nong Bua Lamphu, Nong Khai, Maha Sarakham, Roi Et, Kalasin, Sakon Nakhon, Nakhon Phanom, Mukdahan, Amnat Charoen, Bueng Kan, Yasothon |
| South | Phuket, Songkhla, Surat Thani, Nakhon Si Thammarat, Krabi, Trang, Phatthalung, Satun, Chumphon, Ranong, Phangnga, Pattani, Yala, Narathiwat |

### Query

Add a query to `Overview` handler that groups shipments by `receiver_province`:

```sql
SELECT receiver_province,
       COUNT(*) AS total,
       SUM(CASE WHEN status = 'DELIVERED' THEN 1 ELSE 0 END) AS delivered
FROM shipments
GROUP BY receiver_province
```

In Go, iterate results, map each province to its region via the static map, and aggregate total/delivered counts per region. Provinces not found in the map are grouped as "Other".

### API Response Shape

Add a `by_region` array to the existing response:

```json
{
  "total": 100,
  "active": 60,
  "delivered": 40,
  "by_status": [{ "status": "PENDING", "count": 10 }, ...],
  "by_region": [{ "name": "Central", "total": 15, "delivered": 8 }, ...]
}
```

The `by_region` array is sorted by total descending.

## Frontend Changes

### New Files

1. **`frontend/src/lib/api/analytics.ts`** — exports `fetchAnalytics()` using `api.get("/analytics/overview")`. Types match the API response shape.

2. **`frontend/src/hooks/useAnalytics.ts`** — exports `useAnalytics()` wrapping `useQuery({ queryKey: ["analytics"], queryFn: fetchAnalytics })`. Follows the same pattern as `useHubs`, `useDeliveries`, etc.

### Modified File: `AnalyticsPanel.vue`

- Remove `import { orders, statusLabels } from "@/lib/orders"` (no longer needed)
- Remove `import { Carriers, useCarriers }` if unused — actually `useCarriers` is still needed for the Active Carriers KPI
- Import `useAnalytics` hook
- Replace client-side computeds (`kpis`, `carrierPerformance`, `regionPerformance`, `statusDistribution`, `maxRegionOrders`, `maxStatusCount`) with API data
- KPI cards: `total` from API, `onTime` computed from `delivered/total`, `activeCarriers` from `useCarriers`, `avgDeliveryTime` stays hardcoded
- Status Distribution: maps `by_status` from API using a local `statusLabels` record
- Shipments by Region: renders `by_region` from API directly

### Data Flow

```
AnalyticsPanel.vue
  ├── useAnalytics() → GET /api/analytics/overview → { total, active, delivered, by_status, by_region }
  └── useCarriers()  → GET /api/carriers          → [{ id, name, status, ... }]
```

No more dependency on the hardcoded `orders` array for analytics.

## Files Touched

| File | Change |
|------|--------|
| `backend/internal/analytics/handler.go` | Add province→region map, query, by_region aggregation |
| `frontend/src/lib/api/analytics.ts` | New — fetchAnalytics + types |
| `frontend/src/hooks/useAnalytics.ts` | New — useAnalytics hook |
| `frontend/src/components/AnalyticsPanel.vue` | Replace client-side computeds with API hook |
