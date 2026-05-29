import { describe, it, expect, beforeEach } from "vitest";
import { nextTick } from "vue";
import { mount } from "@vue/test-utils";
import { VueQueryPlugin, QueryClient } from "@tanstack/vue-query";
import { createPinia, setActivePinia } from "pinia";

describe("CarriersView", () => {
  beforeEach(() => {
    setActivePinia(createPinia());
  });

  it("renders tab headers", async () => {
    const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
    const { default: CarriersView } = await import("./CarriersView.vue");
    const wrapper = mount(CarriersView, {
      global: {
        plugins: [[VueQueryPlugin, { queryClient }], createPinia()],
        stubs: { Card: true, CardContent: true },
      },
    });
    await new Promise((r) => setTimeout(r, 200));
    expect(wrapper.text()).toContain("Active Deliveries");
    expect(wrapper.text()).toContain("Hubs");
    expect(wrapper.text()).toContain("Analytics");
  });

  it("switches tabs on click", async () => {
    const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
    const { default: CarriersView } = await import("./CarriersView.vue");
    const wrapper = mount(CarriersView, {
      global: {
        plugins: [[VueQueryPlugin, { queryClient }], createPinia()],
        stubs: { Card: true, CardContent: true },
      },
    });
    await new Promise((r) => setTimeout(r, 200));

    const buttons = wrapper.findAll("button");
    const hubsBtn = buttons.find((b) => b.text().includes("Hubs"));
    expect(hubsBtn).toBeDefined();

    await hubsBtn!.trigger("click");
    await nextTick();

    const updatedButtons = wrapper.findAll("button");
    const activeBtn = updatedButtons.find((b) => b.classes().includes("bg-card"));
    expect(activeBtn).toBeDefined();
    expect(activeBtn!.text()).toContain("Hubs");
  });
});
