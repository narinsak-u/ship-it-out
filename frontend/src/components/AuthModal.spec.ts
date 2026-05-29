import { describe, it, expect, vi, beforeEach } from "vitest";
import { mount } from "@vue/test-utils";
import { createPinia, setActivePinia } from "pinia";
import { toast } from "vue-sonner";
import AuthModal from "./AuthModal.vue";

vi.mock("@/lib/api/client", () => ({
  api: { get: vi.fn(), post: vi.fn() },
}));

import { api } from "@/lib/api/client";

const mockUser = {
  id: 1,
  name: "Test Admin",
  email: "admin@gmail.com",
  role: "admin",
  created_at: "2026-01-01T00:00:00Z",
};

const stubs = {
  Dialog: { template: "<div><slot /></div>", props: ["open"] },
  DialogContent: { template: "<div><slot /></div>" },
  DialogHeader: { template: "<div><slot /></div>" },
  DialogTitle: { template: "<div><slot /></div>" },
  Input: {
    template:
      '<input :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)" />',
    props: ["modelValue"],
  },
  Button: { template: "<button><slot /></button>" },
  Separator: { template: "<hr />" },
};

describe("AuthModal", () => {
  let pinia: ReturnType<typeof createPinia>;

  beforeEach(() => {
    pinia = createPinia();
    setActivePinia(pinia);
  });

  it("shows sign in form by default", async () => {
    const wrapper = mount(AuthModal, {
      props: { open: true },
      global: { plugins: [pinia], stubs },
    });
    expect(wrapper.text()).toContain("Sign In");
  });

  it("switches to sign up tab", async () => {
    const wrapper = mount(AuthModal, {
      props: { open: true },
      global: { plugins: [pinia], stubs },
    });
    const signupBtns = wrapper.findAll("button").filter((b) => b.text().includes("Sign Up"));
    if (signupBtns.length > 0) {
      await signupBtns[0].trigger("click");
      await wrapper.vm.$nextTick();
    }
    expect(wrapper.text()).toContain("Create Account");
  });

  it("has guest mode button", async () => {
    const wrapper = mount(AuthModal, {
      props: { open: true },
      global: { plugins: [pinia], stubs },
    });
    expect(wrapper.text()).toContain("Continue as Guest");
  });

  it("shows error on signup with empty fields", async () => {
    const wrapper = mount(AuthModal, {
      props: { open: true },
      global: { plugins: [pinia], stubs },
    });
    const signupBtns = wrapper.findAll("button").filter((b) => b.text().includes("Sign Up"));
    if (signupBtns.length > 0) {
      await signupBtns[0].trigger("click");
      await wrapper.vm.$nextTick();
    }
    await wrapper.find("form").trigger("submit");
    await wrapper.vm.$nextTick();
    expect(wrapper.text()).toContain("Please fill in all fields");
  });

  it("shows error when signup passwords do not match", async () => {
    const wrapper = mount(AuthModal, {
      props: { open: true },
      global: { plugins: [pinia], stubs },
    });
    const signupBtns = wrapper.findAll("button").filter((b) => b.text().includes("Sign Up"));
    if (signupBtns.length > 0) {
      await signupBtns[0].trigger("click");
      await wrapper.vm.$nextTick();
    }
    const inputs = wrapper.findAll("input");
    await inputs[0].setValue("John");
    await inputs[1].setValue("john@test.com");
    await inputs[2].setValue("password123");
    await inputs[3].setValue("different");
    await wrapper.find("form").trigger("submit");
    await wrapper.vm.$nextTick();
    expect(wrapper.text()).toContain("Passwords do not match");
  });

  it("calls toast.success and emits authenticated on successful login", async () => {
    vi.mocked(api.post).mockResolvedValue({ data: { user: mockUser } });
    vi.mocked(api.get).mockResolvedValue({ data: mockUser });
    const wrapper = mount(AuthModal, {
      props: { open: true },
      global: { plugins: [pinia], stubs },
    });
    await wrapper.find("form").trigger("submit");
    await new Promise((r) => setTimeout(r, 50));
    expect(toast.success).toHaveBeenCalledWith("Signed in successfully");
    expect(wrapper.emitted("authenticated")).toBeTruthy();
  });

  it("shows error message on failed login", async () => {
    vi.mocked(api.post).mockResolvedValue({ error: "Invalid credentials" });
    const wrapper = mount(AuthModal, {
      props: { open: true },
      global: { plugins: [pinia], stubs },
    });
    await wrapper.find("form").trigger("submit");
    await new Promise((r) => setTimeout(r, 50));
    expect(wrapper.text()).toContain("Invalid credentials");
    expect(wrapper.emitted("authenticated")).toBeFalsy();
  });
});
