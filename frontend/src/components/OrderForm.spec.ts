import { describe, it, expect, vi } from "vitest";
import { mount } from "@vue/test-utils";
import OrderForm from "./OrderForm.vue";

vi.mock("@/lib/geocode", () => ({
  geocodeAddress: vi.fn().mockResolvedValue({ lat: 13.75, lng: 100.5 }),
}));

const stubs = {
  ThaiAddressGroup: { template: "<div><legend>{{ label }}</legend></div>", props: ["label", "modelValue", "errors"] },
  Input: {
    template:
      '<input :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)" />',
    props: ["modelValue"],
  },
  Button: { template: "<button><slot /></button>" },
};

describe("OrderForm", () => {
  it("renders sender and receiver sections", () => {
    const wrapper = mount(OrderForm, {
      props: { initial: undefined, isEditing: false, pending: false },
      global: { stubs },
    });
    expect(wrapper.text()).toContain("Sender Info");
    expect(wrapper.text()).toContain("Receiver Info");
  });

  it("renders create mode button text", () => {
    const wrapper = mount(OrderForm, {
      props: { initial: undefined, isEditing: false, pending: false },
      global: { stubs },
    });
    expect(wrapper.text()).toContain("Create Order");
  });

  it("renders cancel button", () => {
    const wrapper = mount(OrderForm, {
      props: { initial: undefined, isEditing: false, pending: false },
      global: { stubs },
    });
    expect(wrapper.text()).toContain("Cancel");
  });

  it("cancels on cancel button click", async () => {
    const wrapper = mount(OrderForm, {
      props: { initial: undefined, isEditing: false, pending: false },
      global: { stubs },
    });
    const cancelBtns = wrapper.findAll("button").filter((b) => b.text().includes("Cancel"));
    if (cancelBtns.length > 0) {
      await cancelBtns[0].trigger("click");
      expect(wrapper.emitted("cancel")).toBeTruthy();
    }
  });
});
