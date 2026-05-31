import { describe, it, expect, vi, beforeEach } from "vitest";
import { mount, flushPromises } from "@vue/test-utils";
import { createRouter, createMemoryHistory } from "vue-router";
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

describe("OrderFormView", () => {
  beforeEach(() => {
    setActivePinia(createPinia());
  });

  async function createView(route: string) {
    const router = createRouter({
      history: createMemoryHistory(),
      routes: [
        { path: "/orders/create", name: "order-create", component: { template: "<div>Create</div>" } },
        { path: "/orders/:orderId/edit", name: "order-edit", component: { template: "<div>Edit</div>" } },
        { path: "/orders/:orderId", name: "order-detail", component: { template: "<div>Detail</div>" } },
        { path: "/orders", name: "orders", component: { template: "<div>Orders</div>" } },
      ],
    });
    const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
    await router.push(route);
    await router.isReady();
    const { default: OrderFormView } = await import("@/views/OrderFormView.vue");
    return mount(OrderFormView, {
      global: {
        plugins: [router, [VueQueryPlugin, { queryClient }], createPinia()],
        stubs: { OrderForm: true, Skeleton: true },
      },
    });
  }

  it("renders create mode", async () => {
    const wrapper = await createView("/orders/create");
    await flushPromises();
    expect(wrapper.text()).toContain("Create Order");
  });
});
