import { describe, it, expect, vi, beforeEach } from "vitest";
import { mount, flushPromises } from "@vue/test-utils";
import ThaiAddressGroup from "./ThaiAddressGroup.vue";

const mocks = vi.hoisted(() => ({
  getSubDistrictNames: vi.fn(),
  getDistrictNames: vi.fn(),
  getProvinceName: vi.fn(),
}));

vi.mock("thai-data", () => mocks);

const stubs = {
  Select: { template: "<div><slot /></div>", props: ["modelValue"] },
  SelectContent: { template: "<div><slot /></div>" },
  SelectGroup: { template: "<div><slot /></div>" },
  SelectItem: { template: "<div><slot /></div>", props: ["value"] },
  SelectTrigger: { template: "<button><slot /></button>" },
  SelectValue: { template: "<span><slot /></span>" },
  Input: {
    template:
      '<div><input :value="modelValue" :placeholder="placeholder" @input="$emit(\'update:modelValue\', $event.target.value)" /><span>{{ placeholder }}</span></div>',
    props: ["modelValue", "placeholder"],
  },
};

const emptyModel = { name: "", zipcode: "", subDistrict: "", district: "", province: "" };

describe("ThaiAddressGroup", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mocks.getDistrictNames.mockReturnValue([]);
    mocks.getProvinceName.mockReturnValue("");
  });

  it("renders fields with label", () => {
    const wrapper = mount(ThaiAddressGroup, {
      props: { label: "Sender Info", modelValue: emptyModel, errors: {} },
      global: { stubs },
    });
    expect(wrapper.text()).toContain("Sender Info");
  });

  it("shows error messages when provided", () => {
    const wrapper = mount(ThaiAddressGroup, {
      props: {
        label: "Receiver",
        modelValue: emptyModel,
        errors: { name: "Required", zipcode: "Required" },
      },
      global: { stubs },
    });
    expect(wrapper.text()).toContain("Required");
  });

  describe("ThaiAddressGroup zipcode watcher", () => {
    it("triggers zipcode lookup when zipcode reaches 5 characters", async () => {
      mocks.getSubDistrictNames.mockReturnValue(["Klong Tan", "Klong Toei"]);

      const wrapper = mount(ThaiAddressGroup, {
        props: { label: "Test", modelValue: emptyModel },
        global: { stubs },
      });

      expect(mocks.getSubDistrictNames).not.toHaveBeenCalled();

      await wrapper.setProps({
        modelValue: { ...emptyModel, zipcode: "10110" },
      });

      expect(mocks.getSubDistrictNames).toHaveBeenCalledWith("10110");
      expect(wrapper.text()).toContain("Klong Tan");
      expect(wrapper.text()).toContain("Klong Toei");
    });

    it("clears subDistrict/district/province when zipcode changes to a different zip", async () => {
      mocks.getSubDistrictNames.mockReturnValue(["Klong Tan", "Klong Toei"]);

      const wrapper = mount(ThaiAddressGroup, {
        props: {
          label: "Test",
          modelValue: {
            name: "",
            zipcode: "10110",
            subDistrict: "Klong Tan",
            district: "Wattana",
            province: "Bangkok",
          },
        },
        global: { stubs },
      });

      mocks.getSubDistrictNames.mockReturnValue(["Bang Na", "Bang Na Tai"]);

      await wrapper.setProps({
        modelValue: {
          name: "",
          zipcode: "10260",
          subDistrict: "Klong Tan",
          district: "Wattana",
          province: "Bangkok",
        },
      });

      const allEmits = wrapper.emitted("update:modelValue") ?? [];
      const lastEmit = allEmits.at(-1)?.[0];
      expect(lastEmit).toMatchObject({
        zipcode: "10260",
        subDistrict: "",
        district: "",
        province: "",
      });
    });

    it("clears all address fields when zipcode shortens from 5 characters", async () => {
      mocks.getSubDistrictNames.mockReturnValue(["Klong Tan", "Klong Toei"]);

      const wrapper = mount(ThaiAddressGroup, {
        props: {
          label: "Test",
          modelValue: {
            name: "Name",
            zipcode: "10110",
            subDistrict: "Klong Tan",
            district: "Wattana",
            province: "Bangkok",
          },
        },
        global: { stubs },
      });

      await wrapper.setProps({
        modelValue: {
          name: "Name",
          zipcode: "10",
          subDistrict: "Klong Tan",
          district: "Wattana",
          province: "Bangkok",
        },
      });

      const allEmits = wrapper.emitted("update:modelValue") ?? [];
      const lastEmit = allEmits.at(-1)?.[0];
      expect(lastEmit).toMatchObject({
        name: "Name",
        zipcode: "10",
        subDistrict: "",
        district: "",
        province: "",
      });
    });

    it("shows no-results when zipcode has no matching sub-districts", async () => {
      mocks.getSubDistrictNames.mockReturnValue([]);

      const wrapper = mount(ThaiAddressGroup, {
        props: { label: "Test", modelValue: emptyModel },
        global: { stubs },
      });

      expect(wrapper.text()).toContain("Enter zipcode first");

      await wrapper.setProps({
        modelValue: { ...emptyModel, zipcode: "99999" },
      });

      expect(mocks.getSubDistrictNames).toHaveBeenCalledWith("99999");
      expect(wrapper.text()).toContain("No sub-districts found");
    });

    it("auto-selects sub-district when only 1 result", async () => {
      mocks.getSubDistrictNames.mockReturnValue(["Klong Tan"]);
      mocks.getDistrictNames.mockReturnValue(["Wattana"]);
      mocks.getProvinceName.mockReturnValue("Bangkok");

      const wrapper = mount(ThaiAddressGroup, {
        props: { label: "Test", modelValue: emptyModel },
        global: { stubs },
      });

      await wrapper.setProps({
        modelValue: { ...emptyModel, zipcode: "10110" },
      });

      expect(mocks.getDistrictNames).toHaveBeenCalledWith("10110");
      expect(mocks.getProvinceName).toHaveBeenCalledWith("10110");

      const lastEmit = (wrapper.emitted("update:modelValue") ?? []).at(-1)?.[0];
      expect(lastEmit).toMatchObject({
        subDistrict: "Klong Tan",
        district: "Wattana",
        province: "Bangkok",
      });
    });

    it("does not auto-select when subDistrict is already set", async () => {
      mocks.getSubDistrictNames.mockReturnValue(["Klong Tan"]);

      const wrapper = mount(ThaiAddressGroup, {
        props: {
          label: "Test",
          modelValue: {
            name: "",
            zipcode: "10110",
            subDistrict: "Klong Tan",
            district: "Wattana",
            province: "Bangkok",
          },
        },
        global: { stubs },
      });

      await flushPromises();
      // Should not emit update:modelValue because:
      // 1. subDistrict "Klong Tan" IS in results (no stale clear)
      // 2. subDistrict is already set (auto-select guard prevents it)
      expect(wrapper.emitted("update:modelValue")).toBeFalsy();
    });

    it("clears stale sub-district when saved value not in results", async () => {
      mocks.getSubDistrictNames.mockReturnValue(["Klong Tan", "Klong Toei"]);

      const wrapper = mount(ThaiAddressGroup, {
        props: {
          label: "Test",
          modelValue: {
            name: "",
            zipcode: "10110",
            subDistrict: "Old Town",
            district: "Old",
            province: "Old Province",
          },
        },
        global: { stubs },
      });

      const lastEmit = (wrapper.emitted("update:modelValue") ?? []).at(-1)?.[0];
      expect(lastEmit).toMatchObject({
        subDistrict: "",
        district: "",
        province: "",
      });
    });
  });
});
