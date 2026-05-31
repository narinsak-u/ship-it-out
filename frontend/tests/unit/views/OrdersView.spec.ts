import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { mount, flushPromises } from "@vue/test-utils";
import { createRouter, createMemoryHistory } from "vue-router";
import { VueQueryPlugin, QueryClient } from "@tanstack/vue-query";
import { createPinia, setActivePinia } from "pinia";

vi.mock("@/lib/api/orders", () => ({
  fetchOrdersPaginated: vi.fn().mockResolvedValue({
    data: [
      {
        id: "ORD-001",
        trackingNumber: "TH202600001",
        status: "in_transit",
        progress: 62,
        customer: { name: "John", coords: { lat: 0, lng: 0 } },
        origin: "A",
        destination: "B",
        carrier: "Test",
        weight: 10,
        items: 2,
        estimatedDelivery: "Jun 1",
        createdAt: "May 28",
        currentCoords: { lat: 0, lng: 0 },
        events: [],
      },
    ],
    pagination: { page: 1, limit: 10, total: 1, totalPages: 1 },
  }),
  deleteOrder: vi.fn().mockResolvedValue(undefined),
}));

import { fetchOrdersPaginated } from "@/lib/api/orders";

async function buildRouter() {
  const { createRouter, createMemoryHistory } = await import("vue-router");
  return createRouter({
    history: createMemoryHistory(),
    routes: [
      { path: "/", name: "home", component: { template: "<div>Home</div>" } },
      { path: "/orders", name: "orders", component: { template: "<div>Orders</div>" } },
      { path: "/orders/create", name: "order-create", component: { template: "<div>Create</div>" } },
      { path: "/orders/:orderId", name: "order-detail", component: { template: "<div>Detail</div>" } },
      { path: "/orders/:orderId/edit", name: "order-edit", component: { template: "<div>Edit</div>" } },
    ],
  });
}

const defaultStubs = {
  StatusBadge: true,
  Pagination: true,
  AuthModal: true,
  ConfirmDialog: true,
  Skeleton: true,
  Input: true,
  Button: true,
  Table: true,
  TableHeader: true,
  TableBody: true,
  TableRow: true,
  TableHead: true,
  TableCell: true,
  ShipmentFilters: true,
  OrderTableRow: true,
};

describe("OrdersView", () => {
  beforeEach(() => {
    setActivePinia(createPinia());
  });

  it("renders shipment manifest title", async () => {
    const router = await buildRouter();
    const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
    const { default: OrdersView } = await import("@/views/OrdersView.vue");
    const wrapper = mount(OrdersView, {
      global: {
        plugins: [router, [VueQueryPlugin, { queryClient }], createPinia()],
        stubs: defaultStubs,
      },
    });
    await flushPromises();
    expect(wrapper.text()).toContain("Shipment manifest");
  });

  it("debounces search input before fetching", async () => {
    vi.useFakeTimers({ shouldAdvanceTime: true });
    const router = await buildRouter();
    const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
    const { default: OrdersView } = await import("@/views/OrdersView.vue");
    const wrapper = mount(OrdersView, {
      global: {
        plugins: [router, [VueQueryPlugin, { queryClient }], createPinia()],
        stubs: {
          ...defaultStubs,
          Input: {
            template:
              '<input :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)" />',
            props: ["modelValue"],
          },
          ShipmentFilters: false,
        },
      },
    });
    await flushPromises();
    vi.mocked(fetchOrdersPaginated).mockClear();

    const searchInput = wrapper.find("input");
    await searchInput.setValue("test");

    vi.advanceTimersByTime(200);
    expect(fetchOrdersPaginated).not.toHaveBeenCalled();

    vi.advanceTimersByTime(300);
    await flushPromises();
    expect(fetchOrdersPaginated).toHaveBeenCalled();

    vi.useRealTimers();
  });
});
