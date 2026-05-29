import { describe, it, expect, vi, beforeEach } from "vitest";
import { mount } from "@vue/test-utils";
import { createPinia, setActivePinia } from "pinia";
import AuthModal from "./AuthModal.vue";

vi.mock("@/lib/api/client", () => ({
  api: { get: vi.fn(), post: vi.fn() },
}));

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
});
