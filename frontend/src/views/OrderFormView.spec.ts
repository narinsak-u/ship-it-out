import { describe, it, expect, vi, beforeEach } from "vitest";
import { mount } from "@vue/test-utils";
import { createRouter, createWebHistory } from "vue-router";
import { VueQueryPlugin, QueryClient } from "@tanstack/vue-query";
import { createPinia, setActivePinia } from "pinia";

vi.mock("@/lib/api/orders", () => ({
  fetchOrder: vi.fn().mockResolvedValue({
    id: "ORD-001",
    trackingNumber: "TH202600001",
    status: "pending",
    progress: 0,
    customer: { name: "John", coords: { lat: 0, lng: 0 } },
    receiver: { name: "Jane", coords: { lat: 0, lng: 0 } },
    origin: "A",
    destination: "B",
    carrier: "Test",
    weight: 10,
    items: 2,
    estimatedDelivery: "Jun 1",
    createdAt: "May 28",
    currentCoords: { lat: 0, lng: 0 },
    events: [],
  }),
  createOrder: vi.fn().mockResolvedValue({ id: "ORD-NEW" }),
  updateOrder: vi.fn().mockResolvedValue({ id: "ORD-001" }),
}));

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: "/orders/create", name: "order-create", component: { template: "<div>Create</div>" } },
    {
      path: "/orders/:orderId/edit",
      name: "order-edit",
      component: { template: "<div>Edit</div>" },
    },
    {
      path: "/orders/:orderId",
      name: "order-detail",
      component: { template: "<div>Detail</div>" },
    },
    { path: "/orders", name: "orders", component: { template: "<div>Orders</div>" } },
  ],
});

describe("OrderFormView", () => {
  beforeEach(() => {
    setActivePinia(createPinia());
  });

  it("renders create mode", async () => {
    const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
    await router.push("/orders/create");
    await router.isReady();
    const { default: OrderFormView } = await import("./OrderFormView.vue");
    const wrapper = mount(OrderFormView, {
      global: {
        plugins: [router, [VueQueryPlugin, { queryClient }], createPinia()],
        stubs: { OrderForm: true, Skeleton: true },
      },
    });
    await new Promise((r) => setTimeout(r, 200));
    expect(wrapper.text()).toContain("Create Order");
  });
});
