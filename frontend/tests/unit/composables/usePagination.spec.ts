import { describe, it, expect } from "vitest";
import { ref, nextTick } from "vue";
import { usePagination } from "@/composables/usePagination";

describe("usePagination", () => {
  it("returns correct page count for 25 items with 10 per page", () => {
    const items = ref([
      1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25,
    ]);
    const { currentPage, totalPages, pageItems } = usePagination(items);
    expect(currentPage.value).toBe(1);
    expect(totalPages.value).toBe(3);
    expect(pageItems.value).toHaveLength(10);
  });

  it("handles empty array", () => {
    const items = ref<number[]>([]);
    const { totalPages, pageItems } = usePagination(items);
    expect(totalPages.value).toBe(0);
    expect(pageItems.value).toHaveLength(0);
  });

  it("handles single page (3 items)", () => {
    const items = ref([1, 2, 3]);
    const { totalPages, pageItems } = usePagination(items);
    expect(totalPages.value).toBe(1);
    expect(pageItems.value).toHaveLength(3);
  });

  it("setPage clamps to valid range", () => {
    const items = ref([1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11]);
    const { currentPage, setPage } = usePagination(items);
    setPage(0);
    expect(currentPage.value).toBe(1);
    setPage(5);
    expect(currentPage.value).toBe(2);
    setPage(2);
    expect(currentPage.value).toBe(2);
  });

  it("nextPage advances and prevPage goes back", () => {
    const items = ref([1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11]);
    const { currentPage, nextPage, prevPage } = usePagination(items);
    nextPage();
    expect(currentPage.value).toBe(2);
    nextPage();
    expect(currentPage.value).toBe(2);
    prevPage();
    expect(currentPage.value).toBe(1);
    prevPage();
    expect(currentPage.value).toBe(1);
  });

  it("resets to page 1 when items change", async () => {
    const items = ref([1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11]);
    const { currentPage, setPage } = usePagination(items);
    setPage(2);
    expect(currentPage.value).toBe(2);
    items.value = [1, 2, 3];
    await nextTick();
    expect(currentPage.value).toBe(1);
  });
});
