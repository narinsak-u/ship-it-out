import { describe, it, expect, vi, beforeEach } from "vitest";
import { mount } from "@vue/test-utils";
import { createRouter, createWebHistory } from "vue-router";
import { VueQueryPlugin, QueryClient } from "@tanstack/vue-query";
import { createPinia, setActivePinia } from "pinia";

vi.mock("@/lib/api/orders", () => ({
  fetchOrder: vi.fn().mockResolvedValue({
    id: "ORD-001",
    trackingNumber: "TH202600001",
    status: "in_transit",
    progress: 62,
    customer: { name: "John Doe", coords: { lat: 13.75, lng: 100.5 } },
    receiver: { name: "Jane Doe", coords: { lat: 18.78, lng: 98.98 } },
    origin: "Bangkok",
    destination: "Chiang Mai",
    carrier: "Pacific Freight",
    weight: 12.4,
    items: 3,
    estimatedDelivery: "Jun 1, 2026",
    createdAt: "May 28, 2026",
    currentCoords: { lat: 16.0, lng: 99.5 },
    events: [],
  }),
  fetchOrderEvents: vi.fn().mockResolvedValue([
    {
      timestamp: "May 28, 10:00",
      location: { name: "Bangkok", lat: 13.75, lng: 100.5 },
      status: "Picked Up",
      description: "",
    },
  ]),
}));

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: "/", name: "home", component: { template: "<div>Home</div>" } },
    { path: "/orders", name: "orders", component: { template: "<div>Orders</div>" } },
    {
      path: "/orders/:orderId",
      name: "order-detail",
      component: { template: "<div>Detail</div>" },
    },
  ],
});

describe("OrderDetailView", () => {
  beforeEach(() => {
    setActivePinia(createPinia());
  });

  it("renders tracking number", async () => {
    const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
    const { default: OrderDetailView } = await import("./OrderDetailView.vue");
    await router.push("/orders/ORD-001");
    await router.isReady();
    const wrapper = mount(OrderDetailView, {
      global: {
        plugins: [router, [VueQueryPlugin, { queryClient }], createPinia()],
        stubs: {
          StatusBadge: true,
          Skeleton: true,
          ShipmentMap: true,
          Card: true,
          CardHeader: true,
          CardTitle: true,
          CardContent: true,
        },
      },
    });
    await new Promise((r) => setTimeout(r, 500));
    expect(wrapper.text()).toContain("TH202600001");
  });

  it("shows 404 for non-existent order", async () => {
    vi.mocked((await import("@/lib/api/orders")).fetchOrder).mockRejectedValue(
      new Error("Not found"),
    );
    const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
    const { default: OrderDetailView } = await import("./OrderDetailView.vue");
    await router.push("/orders/FAKE");
    await router.isReady();
    const wrapper = mount(OrderDetailView, {
      global: {
        plugins: [router, [VueQueryPlugin, { queryClient }], createPinia()],
        stubs: {
          StatusBadge: true,
          Skeleton: true,
          ShipmentMap: true,
          Card: true,
          CardHeader: true,
          CardTitle: true,
          CardContent: true,
        },
      },
    });
    await new Promise((r) => setTimeout(r, 500));
    expect(wrapper.text()).toContain("404");
  });
});
