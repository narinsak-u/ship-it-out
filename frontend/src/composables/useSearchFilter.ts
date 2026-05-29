import { type Ref, computed } from "vue";

type SearchableItem = Record<string, unknown>;

export function useSearchFilter<T extends SearchableItem>(
  items: Ref<T[] | undefined>,
  query: Ref<string>,
  fields: (keyof T)[],
) {
  return computed(() => {
    if (!items.value) return [];

    const q = query.value.trim().toLowerCase();
    if (!q) return items.value;

    return items.value.filter((item) =>
      fields.some((field) => {
        const value = item[field];
        if (value == null) return false;
        return String(value).toLowerCase().includes(q);
      }),
    );
  });
}
