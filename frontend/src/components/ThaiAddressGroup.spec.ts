import { describe, it, expect } from "vitest";
import { mount } from "@vue/test-utils";
import ThaiAddressGroup from "./ThaiAddressGroup.vue";

const stubs = {
  Select: { template: "<div><slot /></div>", props: ["modelValue"] },
  SelectContent: { template: "<div><slot /></div>" },
  SelectGroup: { template: "<div><slot /></div>" },
  SelectItem: { template: "<div><slot /></div>", props: ["value"] },
  SelectTrigger: { template: "<button><slot /></button>" },
  SelectValue: { template: "<span><slot /></span>" },
  Input: {
    template:
      '<input :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)" />',
    props: ["modelValue"],
  },
};

describe("ThaiAddressGroup", () => {
  it("renders fields with label", () => {
    const wrapper = mount(ThaiAddressGroup, {
      props: {
        label: "Sender Info",
        modelValue: { name: "", zipcode: "", subDistrict: "", district: "", province: "" },
        errors: {},
      },
      global: { stubs },
    });
    expect(wrapper.text()).toContain("Sender Info");
  });

  it("shows error messages when provided", () => {
    const wrapper = mount(ThaiAddressGroup, {
      props: {
        label: "Receiver",
        modelValue: { name: "", zipcode: "", subDistrict: "", district: "", province: "" },
        errors: { name: "Required", zipcode: "Required" },
      },
      global: { stubs },
    });
    expect(wrapper.text()).toContain("Required");
  });
});
