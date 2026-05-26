import { ref, computed, watch, type Ref } from "vue";

export function usePagination<T>(items: Ref<T[] | null | undefined>, pageSize = 10) {
  const currentPage = ref(1);

  const totalPages = computed(() => {
    if (!items.value) return 0;
    return Math.max(1, Math.ceil(items.value.length / pageSize));
  });

  const pageItems = computed(() => {
    if (!items.value) return [] as T[];
    const start = (currentPage.value - 1) * pageSize;
    return items.value.slice(start, start + pageSize);
  });

  watch(
    () => items.value?.length,
    () => {
      currentPage.value = 1;
    },
  );

  function setPage(n: number) {
    currentPage.value = Math.max(1, Math.min(n, totalPages.value));
  }

  function nextPage() {
    setPage(currentPage.value + 1);
  }

  function prevPage() {
    setPage(currentPage.value - 1);
  }

  return {
    currentPage,
    totalPages,
    pageItems,
    setPage,
    nextPage,
    prevPage,
  };
}
