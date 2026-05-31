import { describe, it, expect, vi } from "vitest";
import { defineComponent } from "vue";
import { mount, flushPromises } from "@vue/test-utils";
import { VueQueryPlugin, QueryClient } from "@tanstack/vue-query";
import { createPinia, setActivePinia } from "pinia";

vi.mock("@/lib/api/hubs", () => ({
  fetchHubs: vi.fn().mockResolvedValue([
    {
      id: "h1",
      name: "Bangkok Hub",
      carrierId: "c1",
      address: "123 St",
      coords: { lat: 13.75, lng: 100.5 },
      capacity: 1000,
      currentUtilization: 500,
      status: "active",
    },
  ]),
  createHub: vi.fn().mockResolvedValue({
    id: "h-new",
    name: "New Hub",
    carrierId: "c1",
    address: "456 St",
    coords: { lat: 0, lng: 0 },
    capacity: 500,
    currentUtilization: 0,
    status: "active",
  }),
  deleteHub: vi.fn().mockResolvedValue(undefined),
}));

import { useHubs, useCreateHub } from "@/hooks/useHubs";

function mountComposable<T>(setup: () => T) {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  const TestComponent = defineComponent({ setup, template: "<div></div>" });
  return mount(TestComponent, {
    global: { plugins: [[VueQueryPlugin, { queryClient }], createPinia()] },
  });
}

describe("useHubs", () => {
  it("provides reactive data", async () => {
    let data: any;
    mountComposable(() => {
      const query = useHubs();
      data = query.data;
      return {};
    });
    await flushPromises();
    expect(data.value).toHaveLength(1);
  });
});

describe("useCreateHub", () => {
  it("creates a hub via mutation", async () => {
    let mutation!: ReturnType<typeof useCreateHub>;
    mountComposable(() => {
      mutation = useCreateHub();
      return {};
    });
    const result = await mutation.mutateAsync({
      name: "New Hub",
      carrierId: "c1",
      address: "456 St",
      coords: { lat: 0, lng: 0 },
      capacity: 500,
      currentUtilization: 0,
      status: "active",
    });
    expect(result.id).toBe("h-new");
  });
});
