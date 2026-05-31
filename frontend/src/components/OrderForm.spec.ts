import { describe, it, expect, vi } from "vitest";
import { mount, flushPromises } from "@vue/test-utils";
import OrderForm from "./OrderForm.vue";

vi.mock("@/lib/geocode", () => ({
  geocodeAddress: vi.fn().mockResolvedValue({ lat: 13.75, lng: 100.5 }),
}));

import { geocodeAddress } from "@/lib/geocode";

const stubs = {
  ThaiAddressGroup: {
    template: "<div><legend>{{ label }}</legend></div>",
    props: ["label", "modelValue", "errors"],
  },
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
    expect(cancelBtns.length).toBeGreaterThan(0);
    await cancelBtns[0].trigger("click");
    expect(wrapper.emitted("cancel")).toBeTruthy();
  });

  it("renders save changes in edit mode", () => {
    const wrapper = mount(OrderForm, {
      props: { initial: undefined, isEditing: true, pending: false },
      global: { stubs },
    });
    expect(wrapper.text()).toContain("Save Changes");
    expect(wrapper.text()).toContain("Estimated Delivery");
  });

  it("shows weight validation error on submit with empty form", async () => {
    const wrapper = mount(OrderForm, {
      props: { initial: undefined, isEditing: false, pending: false },
      global: { stubs },
    });
    await wrapper.find("form").trigger("submit");
    expect(wrapper.text()).toContain("Required");
  });

  it("shows geocode error inline on submission failure", async () => {
    vi.mocked(geocodeAddress).mockRejectedValue(
      new Error("Location lookup failed. Please try again later."),
    );
    const wrapper = mount(OrderForm, {
      props: {
        initial: {
          customer: {
            name: "John",
            zipcode: "10100",
            subDistrict: "A",
            district: "B",
            province: "C",
            coords: { lat: 13.75, lng: 100.5 },
          },
          receiver: {
            name: "Jane",
            zipcode: "10200",
            subDistrict: "D",
            district: "E",
            province: "F",
            coords: { lat: 13.75, lng: 100.5 },
          },
          weight: 12.4,
          items: 3,
        },
        isEditing: false,
        pending: false,
      },
      global: { stubs },
    });
    await wrapper.find("form").trigger("submit");
    await flushPromises();
    expect(wrapper.text()).toContain("Location lookup failed");
  });

  it("emits submit event on successful submission", async () => {
    vi.mocked(geocodeAddress).mockResolvedValue({ lat: 13.75, lng: 100.5 });
    const wrapper = mount(OrderForm, {
      props: {
        initial: {
          customer: {
            name: "John",
            zipcode: "10100",
            subDistrict: "A",
            district: "B",
            province: "C",
            coords: { lat: 13.75, lng: 100.5 },
          },
          receiver: {
            name: "Jane",
            zipcode: "10200",
            subDistrict: "D",
            district: "E",
            province: "F",
            coords: { lat: 13.75, lng: 100.5 },
          },
          weight: 12.4,
          items: 3,
        },
        isEditing: false,
        pending: false,
      },
      global: { stubs },
    });
    await wrapper.find("form").trigger("submit");
    await flushPromises();
    expect(wrapper.emitted("submit")).toBeTruthy();
  });
});
