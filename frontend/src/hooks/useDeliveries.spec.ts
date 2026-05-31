import { describe, it, expect, vi } from "vitest";
import { defineComponent } from "vue";
import { mount, flushPromises } from "@vue/test-utils";
import { VueQueryPlugin, QueryClient } from "@tanstack/vue-query";
import { createPinia, setActivePinia } from "pinia";

vi.mock("@/lib/api/orders", () => ({
  fetchActiveDeliveries: vi.fn().mockResolvedValue([
    {
      id: "ORD-001",
      trackingNumber: "TH202600001",
      status: "in_transit",
      progress: 50,
      customer: { name: "John" },
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
  ]),
  updateShipmentStatus: vi.fn().mockResolvedValue({ id: "ORD-001", status: "delivered" }),
}));

import { useActiveDeliveries, useUpdateShipmentStatus } from "./useDeliveries";

function mountComposable<T>(setup: () => T) {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  const TestComponent = defineComponent({ setup, template: "<div></div>" });
  return mount(TestComponent, {
    global: { plugins: [[VueQueryPlugin, { queryClient }], createPinia()] },
  });
}

describe("useActiveDeliveries", () => {
  it("returns reactive query data", async () => {
    let data: any;
    mountComposable(() => {
      const query = useActiveDeliveries();
      data = query.data;
      return {};
    });
    await flushPromises();
    expect(data.value).toHaveLength(1);
  });
});

describe("useUpdateShipmentStatus", () => {
  it("mutates shipment status", async () => {
    let mutation!: ReturnType<typeof useUpdateShipmentStatus>;
    mountComposable(() => {
      mutation = useUpdateShipmentStatus();
      return {};
    });
    const result = await mutation.mutateAsync({
      orderId: "ORD-001",
      status: "delivered",
      hubId: "hub-1",
    });
    expect(result.status).toBe("delivered");
  });
});
