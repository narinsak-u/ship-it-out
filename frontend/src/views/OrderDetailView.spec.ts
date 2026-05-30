import { describe, it, expect, vi, beforeEach } from "vitest";
import { mount, flushPromises } from "@vue/test-utils";
import { createRouter, createMemoryHistory } from "vue-router";
import { VueQueryPlugin, QueryClient } from "@tanstack/vue-query";
import { createPinia, setActivePinia } from "pinia";
import type { ShipmentStatus } from "@/lib/orders";

vi.mock("@/lib/api/orders", () => ({
  fetchOrder: vi.fn(),
  fetchOrderEvents: vi.fn(),
}));
import { fetchOrder, fetchOrderEvents } from "@/lib/api/orders";

async function createView(orderId = "ORD-001") {
  const router = createRouter({
    history: createMemoryHistory(),
    routes: [
      { path: "/", name: "home", component: { template: "<div>Home</div>" } },
      { path: "/orders", name: "orders", component: { template: "<div>Orders</div>" } },
      { path: "/orders/:orderId", name: "order-detail", component: { template: "<div>Detail</div>" } },
    ],
  });
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  const { default: OrderDetailView } = await import("./OrderDetailView.vue");
  await router.push(`/orders/${orderId}`);
  await router.isReady();
  return mount(OrderDetailView, {
    global: {
      plugins: [router, [VueQueryPlugin, { queryClient }], createPinia()],
      stubs: { StatusBadge: true, Skeleton: true, ShipmentMap: true, Card: true, CardHeader: true, CardTitle: true, CardContent: true },
    },
  });
}

const mockOrder = {
  id: "ORD-001",
  trackingNumber: "TH202600001",
  status: "in_transit" as ShipmentStatus,
  progress: 62,
  customer: { name: "John Doe", zipcode: "10100", subDistrict: "A", district: "B", province: "C", coords: { lat: 13.75, lng: 100.5 } },
  receiver: { name: "Jane Doe", zipcode: "50100", subDistrict: "D", district: "E", province: "F", coords: { lat: 18.78, lng: 98.98 } },
  origin: "Bangkok",
  destination: "Chiang Mai",
  carrier: "Pacific Freight",
  weight: 12.4,
  items: 3,
  estimatedDelivery: "Jun 1, 2026",
  createdAt: "May 28, 2026",
  currentCoords: { lat: 16.0, lng: 99.5 },
  events: [],
};

const mockEvents = [
  {
    timestamp: "May 28, 10:00",
    location: { name: "Bangkok", lat: 13.75, lng: 100.5 },
    status: "Picked Up",
    description: "",
  },
];

describe("OrderDetailView", () => {
  beforeEach(() => {
    setActivePinia(createPinia());
    vi.clearAllMocks();
  });

  it("renders tracking number", async () => {
    vi.mocked(fetchOrder).mockResolvedValue(mockOrder);
    vi.mocked(fetchOrderEvents).mockResolvedValue(mockEvents);
    const wrapper = await createView();
    await flushPromises();
    expect(wrapper.text()).toContain("TH202600001");
  });

  it("shows 404 for non-existent order", async () => {
    vi.mocked(fetchOrder).mockRejectedValue(new Error("Not found"));
    const wrapper = await createView("FAKE");
    await flushPromises();
    await flushPromises();
    expect(wrapper.text()).toContain("404");
  });
});
