import { describe, it, expect, vi, beforeEach } from "vitest";
import { mount } from "@vue/test-utils";
import { createRouter, createWebHistory } from "vue-router";
import { VueQueryPlugin, QueryClient } from "@tanstack/vue-query";
import { createPinia, setActivePinia } from "pinia";

vi.mock("@/lib/api/orders", () => ({
  fetchOrdersPaginated: vi.fn().mockResolvedValue({
    data: [
      {
        id: "ORD-001", trackingNumber: "TH202600001", status: "in_transit", progress: 62,
        customer: { name: "John", coords: { lat: 0, lng: 0 } },
        origin: "A", destination: "B", carrier: "Test", weight: 10, items: 2,
        estimatedDelivery: "Jun 1", createdAt: "May 28",
        currentCoords: { lat: 0, lng: 0 }, events: [],
      },
    ],
    pagination: { page: 1, limit: 10, total: 1, totalPages: 1 },
  }),
  deleteOrder: vi.fn().mockResolvedValue(undefined),
}));

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: "/", name: "home", component: { template: "<div>Home</div>" } },
    { path: "/orders", name: "orders", component: { template: "<div>Orders</div>" } },
    { path: "/orders/create", name: "order-create", component: { template: "<div>Create</div>" } },
    { path: "/orders/:orderId", name: "order-detail", component: { template: "<div>Detail</div>" } },
    { path: "/orders/:orderId/edit", name: "order-edit", component: { template: "<div>Edit</div>" } },
  ],
});

describe("OrdersView", () => {
  beforeEach(() => {
    setActivePinia(createPinia());
  });

  it("renders shipment manifest title", async () => {
    const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
    const { default: OrdersView } = await import("./OrdersView.vue");
    const wrapper = mount(OrdersView, {
      global: {
        plugins: [router, [VueQueryPlugin, { queryClient }], createPinia()],
        stubs: {
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
        },
      },
    });
    await new Promise(r => setTimeout(r, 300));
    expect(wrapper.text()).toContain("Shipment manifest");
  });
});
