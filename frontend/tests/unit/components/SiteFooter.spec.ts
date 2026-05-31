import { describe, it, expect } from "vitest";
import { mount } from "@vue/test-utils";
import SiteFooter from "@/components/SiteFooter.vue";

describe("SiteFooter", () => {
  it("renders copyright text", () => {
    const wrapper = mount(SiteFooter);
    expect(wrapper.text()).toContain("2026");
  });
});
