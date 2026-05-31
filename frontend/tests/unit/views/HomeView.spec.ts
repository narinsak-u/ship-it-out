import { describe, it, expect, vi, beforeEach } from "vitest";
import { mount, flushPromises } from "@vue/test-utils";
import { createRouter, createMemoryHistory } from "vue-router";
import { VueQueryPlugin, QueryClient } from "@tanstack/vue-query";
import { createPinia, setActivePinia } from "pinia";

vi.mock("@/lib/api/orders", () => ({
  fetchActiveDeliveries: vi.fn().mockResolvedValue([
    {
      id: "ORD-001",
      trackingNumber: "TH202600001",
      status: "in_transit",
      progress: 62,
      customer: { name: "John", coords: { lat: 0, lng: 0 } },
      origin: "Bangkok",
      destination: "Chiang Mai",
      carrier: "Test",
      weight: 10,
      items: 2,
      estimatedDelivery: "Jun 1",
      createdAt: "May 28",
      currentCoords: { lat: 0, lng: 0 },
      events: [],
    },
  ]),
}));

vi.mock("@/lib/api/analytics", () => ({
  fetchAnalytics: vi
    .fn()
    .mockResolvedValue({ total: 100, active: 45, delivered: 55, by_status: [], by_region: [] }),
}));

vi.mock("@/lib/api/tracking", () => ({
  trackShipment: vi.fn().mockResolvedValue({ shipment: { id: "ORD-001" } }),
}));

async function createView() {
  const router = createRouter({
    history: createMemoryHistory(),
    routes: [
      { path: "/", name: "home", component: { template: "<div>Home</div>" } },
      {
        path: "/orders/:orderId",
        name: "order-detail",
        component: { template: "<div>Detail</div>" },
      },
    ],
  });
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  const { default: HomeView } = await import("@/views/HomeView.vue");
  await router.push("/");
  await router.isReady();
  return mount(HomeView, {
    global: {
      plugins: [router, [VueQueryPlugin, { queryClient }], createPinia()],
      stubs: { StatusBadge: true, Input: true, Button: true },
    },
  });
}

describe("HomeView", () => {
  beforeEach(() => {
    setActivePinia(createPinia());
    vi.clearAllMocks();
  });

  it("renders hero section", async () => {
    const wrapper = await createView();
    await flushPromises();
    expect(wrapper.text()).toContain("Move fast");
  });

  it("shows tracking search form", async () => {
    const wrapper = await createView();
    await flushPromises();
    expect(wrapper.find("form").exists()).toBe(true);
  });
});
