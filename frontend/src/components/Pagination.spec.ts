import { describe, it, expect } from "vitest";
import { mount } from "@vue/test-utils";
import Pagination from "./Pagination.vue";

describe("Pagination", () => {
  it("shows correct page range info", () => {
    const wrapper = mount(Pagination, { props: { currentPage: 1, totalPages: 5, totalItems: 50, pageSize: 10 } });
    expect(wrapper.text()).toContain("Showing");
    expect(wrapper.text()).toContain("50");
  });

  it("disables prev button on first page", () => {
    const wrapper = mount(Pagination, { props: { currentPage: 1, totalPages: 5, totalItems: 50, pageSize: 10 } });
    const buttons = wrapper.findAll("button");
    const prevBtn = buttons[0];
    expect(prevBtn.attributes("disabled")).toBeDefined();
  });

  it("emits update:currentPage when clicking a page number", async () => {
    const wrapper = mount(Pagination, { props: { currentPage: 1, totalPages: 5, totalItems: 50, pageSize: 10 } });
    const buttons = wrapper.findAll("button");
    for (const btn of buttons) {
      if (btn.text() === "2") {
        await btn.trigger("click");
        break;
      }
    }
    expect(wrapper.emitted("update:currentPage")).toBeTruthy();
  });

  it("shows ellipsis for large page counts", () => {
    const wrapper = mount(Pagination, { props: { currentPage: 5, totalPages: 20, totalItems: 200, pageSize: 10 } });
    expect(wrapper.text()).toContain("\u2026");
  });

  it("does not render when totalPages <= 1", () => {
    const wrapper = mount(Pagination, { props: { currentPage: 1, totalPages: 1, totalItems: 5, pageSize: 10 } });
    expect(wrapper.find(".flex").exists()).toBe(false);
  });
});
