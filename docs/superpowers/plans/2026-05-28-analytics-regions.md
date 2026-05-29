# Analytics: Shipments by Region — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace client-side hardcoded analytics with backend API data, adding per-region shipment breakdown.

**Architecture:** Backend groups shipments by `receiver_province` via SQL, maps provinces to regions in Go, and returns a `by_region` array. Frontend creates an analytics API client + hook, then consumes all API data (KPIs, status distribution, regions) in AnalyticsPanel, removing the `orders` array dependency.

**Tech Stack:** Go (Fiber + GORM), Vue 3 (Composition API), TanStack Vue Query

---

### Task 1: Backend — Add province-to-region mapping and by_region query

**Files:**
- Modify: `backend/internal/analytics/handler.go`

- [ ] **Step 1: Add province→region map, query struct, and region aggregation logic**

Replace the file content with:

```go
package analytics

import (
	"github.com/gofiber/fiber/v2"
	"github.com/narinsak-u/backend/internal/database"
	"github.com/narinsak-u/backend/internal/models"
	"github.com/narinsak-u/backend/pkg/utils"
)

var provinceRegion = map[string]string{
	"Bangkok": "Central", "Nonthaburi": "Central", "Pathum Thani": "Central",
	"Samut Prakan": "Central", "Phra Nakhon Si Ayutthaya": "Central",
	"Nakhon Pathom": "Central", "Samut Sakhon": "Central", "Samut Songkhram": "Central",
	"Sing Buri": "Central", "Ang Thong": "Central", "Lopburi": "Central",
	"Saraburi": "Central", "Chai Nat": "Central", "Suphan Buri": "Central",
	"Chonburi": "East", "Rayong": "East", "Chachoengsao": "East",
	"Chanthaburi": "East", "Trat": "East", "Prachinburi": "East",
	"Sa Kaeo": "East", "Nakhon Nayok": "East",
	"Chiang Mai": "North", "Chiang Rai": "North", "Lampang": "North",
	"Lamphun": "North", "Phrae": "North", "Nan": "North", "Phayao": "North",
	"Mae Hong Son": "North", "Uttaradit": "North",
	"Kanchanaburi": "West", "Phetchaburi": "West", "Prachuap Khiri Khan": "West",
	"Ratchaburi": "West", "Tak": "West",
	"Khon Kaen": "North-east", "Nakhon Ratchasima": "North-east",
	"Udon Thani": "North-east", "Ubon Ratchathani": "North-east",
	"Buriram": "North-east", "Surin": "North-east", "Sisaket": "North-east",
	"Chaiyaphum": "North-east", "Loei": "North-east", "Nong Bua Lamphu": "North-east",
	"Nong Khai": "North-east", "Maha Sarakham": "North-east", "Roi Et": "North-east",
	"Kalasin": "North-east", "Sakon Nakhon": "North-east", "Nakhon Phanom": "North-east",
	"Mukdahan": "North-east", "Amnat Charoen": "North-east", "Bueng Kan": "North-east",
	"Yasothon": "North-east",
	"Phuket": "South", "Songkhla": "South", "Surat Thani": "South",
	"Nakhon Si Thammarat": "South", "Krabi": "South", "Trang": "South",
	"Phatthalung": "South", "Satun": "South", "Chumphon": "South",
	"Ranong": "South", "Phangnga": "South", "Pattani": "South",
	"Yala": "South", "Narathiwat": "South",
}

type provinceCount struct {
	Province  string
	Total     int64
	Delivered int64
}

type regionCount struct {
	Name      string `json:"name"`
	Total     int64  `json:"total"`
	Delivered int64  `json:"delivered"`
}

func Overview(c *fiber.Ctx) error {
	var total int64
	database.DB.Model(&models.Shipment{}).Count(&total)

	var active int64
	database.DB.Model(&models.Shipment{}).Where("status NOT IN ?", []string{"DELIVERED", "RETURNED"}).Count(&active)

	var delivered int64
	database.DB.Model(&models.Shipment{}).Where("status = ?", "DELIVERED").Count(&delivered)

	type StatusCount struct {
		Status string
		Count  int64
	}
	var byStatus []StatusCount
	database.DB.Model(&models.Shipment{}).Select("status, count(*) as count").Group("status").Scan(&byStatus)

	var byProvince []provinceCount
	database.DB.Model(&models.Shipment{}).
		Select("receiver_province as province, count(*) as total, sum(case when status = 'DELIVERED' then 1 else 0 end) as delivered").
		Group("receiver_province").
		Scan(&byProvince)

	regionMap := make(map[string]*regionCount)
	for _, p := range byProvince {
		region := provinceRegion[p.Province]
		if region == "" {
			region = "Other"
		}
		r, ok := regionMap[region]
		if !ok {
			r = &regionCount{Name: region}
			regionMap[region] = r
		}
		r.Total += p.Total
		r.Delivered += p.Delivered
	}

	byRegion := make([]regionCount, 0, len(regionMap))
	for _, r := range regionMap {
		byRegion = append(byRegion, *r)
	}

	return utils.Success(c, fiber.Map{
		"total":     total,
		"active":    active,
		"delivered": delivered,
		"by_status": byStatus,
		"by_region": byRegion,
	})
}
```

- [ ] **Step 2: Build and verify backend compiles**

Run: `cd backend && go build ./...`
Expected: success (no errors)

- [ ] **Step 3: Commit**

```bash
git add backend/internal/analytics/handler.go
git commit -m "feat(analytics): add by_region aggregation to overview endpoint"
```

---

### Task 2: Frontend — Create analytics API client

**Files:**
- Create: `frontend/src/lib/api/analytics.ts`

- [ ] **Step 1: Create the API client file**

```typescript
import { api } from "@/lib/api/client";

export interface RegionCount {
  name: string;
  total: number;
  delivered: number;
}

export interface StatusCount {
  status: string;
  count: number;
}

export interface AnalyticsOverview {
  total: number;
  active: number;
  delivered: number;
  by_status: StatusCount[];
  by_region: RegionCount[];
}

export async function fetchAnalytics(): Promise<AnalyticsOverview> {
  const result = await api.get<AnalyticsOverview>("/analytics/overview");
  if (result.error) throw new Error(result.error);
  return result.data!;
}
```

- [ ] **Step 2: Commit**

```bash
git add frontend/src/lib/api/analytics.ts
git commit -m "feat(analytics): add frontend API client for analytics overview"
```

---

### Task 3: Frontend — Create useAnalytics hook

**Files:**
- Create: `frontend/src/hooks/useAnalytics.ts`

- [ ] **Step 1: Create the hook file**

```typescript
import { useQuery } from "@tanstack/vue-query";
import { fetchAnalytics } from "@/lib/api/analytics";

export function useAnalytics() {
  return useQuery({
    queryKey: ["analytics"],
    queryFn: fetchAnalytics,
  });
}
```

- [ ] **Step 2: Commit**

```bash
git add frontend/src/hooks/useAnalytics.ts
git commit -m "feat(analytics): add useAnalytics Vue Query hook"
```

---

### Task 4: Frontend — Update AnalyticsPanel to consume API data

**Files:**
- Modify: `frontend/src/components/AnalyticsPanel.vue`

- [ ] **Step 1: Replace the script section**

Replace the entire `<script setup>` block with:

```vue
<script setup lang="ts">
import { computed } from "vue";
import { useCarriers } from "@/hooks/useCarriers";
import { useAnalytics } from "@/hooks/useAnalytics";
import { statusLabels } from "@/lib/orders";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import Skeleton from "@/components/ui/Skeleton.vue";
import Button from "@/components/ui/Button.vue";

const { data: carriersData } = useCarriers();
const { data: analytics, isLoading, isError, refetch } = useAnalytics();

const kpis = computed(() => {
  if (!analytics.value) return null;
  const total = analytics.value.total;
  const delivered = analytics.value.delivered;
  const onTime = Math.round((delivered / Math.max(total, 1)) * 100);
  const activeCarriers = carriersData.value?.filter((c) => c.status === "active").length ?? 0;
  return { total, onTime, activeCarriers, avgDeliveryTime: "3.2 days" };
});

const regionPerformance = computed(() => {
  if (!analytics.value) return [];
  return analytics.value.by_region
    .map((r) => ({
      ...r,
      onTimeRate: r.total > 0 ? Math.round((r.delivered / r.total) * 100) : 0,
    }))
    .sort((a, b) => b.total - a.total);
});

const statusDistribution = computed(() => {
  if (!analytics.value) return [];
  const total = analytics.value.total;
  return analytics.value.by_status.map((s) => ({
    status: s.status.toLowerCase(),
    label: (statusLabels as Record<string, string>)[s.status.toLowerCase()] ?? s.status,
    count: s.count,
    pct: Math.round((s.count / Math.max(total, 1)) * 100),
  }));
});

const maxRegionOrders = computed(() =>
  Math.max(...regionPerformance.value.map((r) => r.total), 1),
);

const maxStatusCount = computed(() =>
  Math.max(...statusDistribution.value.map((s) => s.count), 1),
);
</script>
```

- [ ] **Step 2: Build and fix any type errors**

Run: `cd frontend && npm run build`
Expected: Build succeeds. If any type errors, fix them and re-run.

- [ ] **Step 3: Commit**

```bash
git add frontend/src/components/AnalyticsPanel.vue
git commit -m "feat(analytics): switch AnalyticsPanel to consume API data"
```
