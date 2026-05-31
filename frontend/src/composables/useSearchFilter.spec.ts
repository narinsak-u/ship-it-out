import { describe, it, expect } from "vitest";
import { ref } from "vue";
import { useSearchFilter } from "./useSearchFilter";

interface TestItem {
  id: number;
  name: string;
  email: string;
}

const testData: TestItem[] = [
  { id: 1, name: "Alice Wonderland", email: "alice@example.com" },
  { id: 2, name: "Bob Builder", email: "bob@test.com" },
  { id: 3, name: "Charlie Brown", email: "charlie@example.com" },
];

describe("useSearchFilter", () => {
  it("returns all items when query is empty", () => {
    const items = ref(testData);
    const query = ref("");
    const result = useSearchFilter(items, query, ["name"]);
    expect(result.value).toHaveLength(3);
  });

  it("filters by name case-insensitively", () => {
    const items = ref(testData);
    const query = ref("alice");
    const result = useSearchFilter(items, query, ["name"]);
    expect(result.value).toHaveLength(1);
    expect(result.value[0].name).toBe("Alice Wonderland");
  });

  it("returns empty array when no match", () => {
    const items = ref(testData);
    const query = ref("zzzz");
    const result = useSearchFilter(items, query, ["name"]);
    expect(result.value).toHaveLength(0);
  });

  it("searches across multiple fields", () => {
    const items = ref(testData);
    const query = ref("bob");
    const result = useSearchFilter(items, query, ["name", "email"]);
    expect(result.value).toHaveLength(1);
  });

  it("skips null or undefined field values", () => {
    const items = ref<TestItem[]>([{ id: 4, name: "Dave", email: null as unknown as string }]);
    const query = ref("dave");
    const result = useSearchFilter(items, query, ["name", "email"]);
    expect(result.value).toHaveLength(1);
  });

  it("handles undefined items gracefully", () => {
    const items = ref<TestItem[] | undefined>(undefined);
    const query = ref("test");
    const result = useSearchFilter(items, query, ["name"]);
    expect(result.value).toEqual([]);
  });
});
