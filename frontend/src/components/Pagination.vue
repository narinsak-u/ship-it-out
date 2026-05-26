<script setup lang="ts">
import { computed } from "vue";
import { ChevronLeft, ChevronRight } from "lucide-vue-next";
import { cn } from "@/lib/utils";

interface Props {
  currentPage: number;
  totalPages: number;
  totalItems: number;
  pageSize: number;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  "update:currentPage": [page: number];
}>();

const pages = computed(() => {
  const total = props.totalPages;
  const current = props.currentPage;
  if (total <= 7) {
    return Array.from({ length: total }, (_, i) => i + 1);
  }
  const result: (number | "...")[] = [1];
  if (current > 3) result.push("...");
  const start = Math.max(2, current - 1);
  const end = Math.min(total - 1, current + 1);
  for (let i = start; i <= end; i++) result.push(i);
  if (current < total - 2) result.push("...");
  result.push(total);
  return result;
});

function goTo(page: number) {
  if (page < 1 || page > props.totalPages || page === props.currentPage) return;
  emit("update:currentPage", page);
}
</script>

<template>
  <div
    v-if="totalPages > 1"
    class="flex items-center justify-between border-t border-border px-4 py-3"
  >
    <span class="font-mono text-xs text-muted-foreground">
      Showing {{ (currentPage - 1) * pageSize + 1 }}–{{
        Math.min(currentPage * pageSize, totalItems)
      }}
      of {{ totalItems }}
    </span>
    <div class="flex items-center gap-1">
      <button
        :disabled="currentPage <= 1"
        class="rounded p-1.5 text-muted-foreground transition-colors hover:text-foreground disabled:opacity-30 disabled:pointer-events-none"
        @click="goTo(currentPage - 1)"
      >
        <ChevronLeft class="h-4 w-4" />
      </button>
      <template v-for="p in pages" :key="typeof p === 'string' ? p : p">
        <span v-if="p === '...'" class="px-1 font-mono text-xs text-muted-foreground">…</span>
        <button
          v-else
          @click="goTo(p)"
          :class="
            cn(
              'flex h-7 min-w-7 items-center justify-center rounded px-1.5 font-mono text-xs transition-colors',
              p === currentPage
                ? 'border border-primary bg-primary/15 text-primary'
                : 'text-muted-foreground hover:text-foreground',
            )
          "
        >
          {{ p }}
        </button>
      </template>
      <button
        :disabled="currentPage >= totalPages"
        class="rounded p-1.5 text-muted-foreground transition-colors hover:text-foreground disabled:opacity-30 disabled:pointer-events-none"
        @click="goTo(currentPage + 1)"
      >
        <ChevronRight class="h-4 w-4" />
      </button>
    </div>
  </div>
</template>
