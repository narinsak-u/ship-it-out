import { describe, it, expect, vi } from "vitest";
import { defineComponent } from "vue";
import { mount } from "@vue/test-utils";
import { VueQueryPlugin, QueryClient } from "@tanstack/vue-query";
import { createPinia, setActivePinia } from "pinia";

vi.mock("@/lib/api/orders", () => ({
  createOrder: vi.fn().mockResolvedValue({ id: "ORD-NEW", trackingNumber: "TH2026NEW" }),
  updateOrder: vi.fn().mockResolvedValue({ id: "ORD-001", trackingNumber: "TH202600001" }),
}));

import { useCreateOrder, useUpdateOrder } from "./useOrders";

function mountComposable<T>(setup: () => T) {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  setActivePinia(createPinia());
  const TestComponent = defineComponent({ setup, template: "<div></div>" });
  return mount(TestComponent, {
    global: { plugins: [[VueQueryPlugin, { queryClient }], createPinia()] },
  });
}

describe("useCreateOrder", () => {
  it("returns mutation with mutateAsync function", async () => {
    let mutation!: ReturnType<typeof useCreateOrder>;
    mountComposable(() => { mutation = useCreateOrder(); return {}; });
    const result = await mutation.mutateAsync({
      customer: { name: "T", zipcode: "1", subDistrict: "A", district: "B", province: "C", coords: { lat: 0, lng: 0 } },
      receiver: { name: "R", zipcode: "2", subDistrict: "D", district: "E", province: "F", coords: { lat: 0, lng: 0 } },
      carrier: "Test", weight: 1, items: 1, estimatedDelivery: "",
    });
    expect(result.id).toBe("ORD-NEW");
  });
});

describe("useUpdateOrder", () => {
  it("returns mutation for updating orders", async () => {
    let mutation!: ReturnType<typeof useUpdateOrder>;
    mountComposable(() => { mutation = useUpdateOrder(); return {}; });
    const result = await mutation.mutateAsync({ id: "ORD-001", data: { weight: 20 } });
    expect(result.id).toBe("ORD-001");
  });
});
