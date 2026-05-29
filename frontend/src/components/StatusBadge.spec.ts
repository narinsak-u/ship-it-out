import { describe, it, expect } from "vitest";
import { mount } from "@vue/test-utils";
import StatusBadge from "./StatusBadge.vue";

const statuses = ["pending", "picked_up", "departed", "in_transit", "out_for_delivery", "delivered", "delayed"] as const;

describe("StatusBadge", () => {
  for (const status of statuses) {
    it(`renders correctly for ${status}`, () => {
      const wrapper = mount(StatusBadge, { props: { status } });
      expect(wrapper.text()).toBeTruthy();
    });
  }
});
