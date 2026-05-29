import { describe, it, expect } from "vitest";
import { mount } from "@vue/test-utils";
import ConfirmDialog from "./ConfirmDialog.vue";

const stubs = {
  Dialog: { template: "<div><slot /></div>", props: ["open"] },
  DialogContent: { template: "<div><slot /></div>" },
  DialogHeader: { template: "<div><slot /></div>" },
  DialogTitle: { template: "<div><slot /></div>" },
};

describe("ConfirmDialog", () => {
  it("renders title and description", () => {
    const wrapper = mount(ConfirmDialog, { props: { open: true, title: "Delete?", description: "Are you sure?" }, global: { stubs } });
    expect(wrapper.text()).toContain("Delete?");
    expect(wrapper.text()).toContain("Are you sure?");
  });

  it("emits confirm on confirm button click", async () => {
    const wrapper = mount(ConfirmDialog, { props: { open: true, title: "Delete?", description: "Sure?" }, global: { stubs } });
    const buttons = wrapper.findAll("button");
    const deleteBtn = buttons.find((b) => b.text().includes("Delete"));
    await deleteBtn?.trigger("click");
    expect(wrapper.emitted("confirm")).toBeTruthy();
  });

  it("emits cancel on cancel button click", async () => {
    const wrapper = mount(ConfirmDialog, { props: { open: true, title: "Delete?", description: "Sure?" }, global: { stubs } });
    const buttons = wrapper.findAll("button");
    const cancelBtn = buttons.find((b) => b.text().includes("Cancel"));
    await cancelBtn?.trigger("click");
    expect(wrapper.emitted("cancel")).toBeTruthy();
  });

  it("shows pending state on confirm button", () => {
    const wrapper = mount(ConfirmDialog, { props: { open: true, title: "Delete?", description: "Sure?", pending: true }, global: { stubs } });
    expect(wrapper.text()).toContain("Deleting");
  });
});
