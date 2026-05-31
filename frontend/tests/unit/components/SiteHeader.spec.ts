import { describe, it, expect, vi, beforeEach } from "vitest";
import { mount, flushPromises } from "@vue/test-utils";
import { createPinia, setActivePinia } from "pinia";
import { createRouter, createMemoryHistory } from "vue-router";
import { useAuthStore } from "@/stores/auth";
import SiteHeader from "@/components/SiteHeader.vue";

function createTestRouter() {
  return createRouter({
    history: createMemoryHistory(),
    routes: [
      { path: "/", name: "home", component: { template: "<div>Home</div>" } },
      { path: "/orders", name: "orders", component: { template: "<div>Orders</div>" } },
      { path: "/carriers", name: "carriers", component: { template: "<div>Carriers</div>" } },
    ],
  });
}

const stubs = {
  AuthModal: { template: "<div />" },
};

describe("SiteHeader", () => {
  let pinia: ReturnType<typeof createPinia>;

  beforeEach(() => {
    pinia = createPinia();
    setActivePinia(pinia);
  });

  it("shows sign in button for unauthenticated users", async () => {
    const router = createTestRouter();
    const wrapper = mount(SiteHeader, {
      global: { plugins: [router, pinia], stubs },
    });
    const store = useAuthStore();
    store.loading = false;
    await flushPromises();
    expect(wrapper.text()).toContain("Sign in");
  });

  it("shows user name when authenticated", async () => {
    const router = createTestRouter();
    const wrapper = mount(SiteHeader, {
      global: { plugins: [router, pinia], stubs },
    });
    const store = useAuthStore();
    store.user = { id: 1, name: "Admin", email: "admin@test.com", role: "admin", created_at: "" };
    store.loading = false;
    await flushPromises();
    expect(wrapper.text()).toContain("Admin");
    expect(wrapper.text()).toContain("Sign out");
  });

  it("shows guest label when in guest mode", async () => {
    const router = createTestRouter();
    const wrapper = mount(SiteHeader, {
      global: { plugins: [router, pinia], stubs },
    });
    const store = useAuthStore();
    store.isGuest = true;
    store.loading = false;
    await flushPromises();
    expect(wrapper.text()).toContain("Guest");
  });

  it("renders navigation links", async () => {
    const router = createTestRouter();
    const wrapper = mount(SiteHeader, {
      global: { plugins: [router, pinia], stubs },
    });
    await router.isReady();
    expect(wrapper.text()).toContain("Home");
    expect(wrapper.text()).toContain("Orders");
    expect(wrapper.text()).toContain("Carriers");
  });
});
