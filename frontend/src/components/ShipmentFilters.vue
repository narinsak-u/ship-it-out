<script setup lang="ts">
import { Search, Filter } from "lucide-vue-next";
import Input from "@/components/ui/Input.vue";
import { cn } from "@/lib/utils";
import type { ShipmentStatus } from "@/lib/orders";

interface FilterDef {
  key: ShipmentStatus | "all";
  label: string;
}

defineProps<{
  search: string;
  filter: string;
  filters: FilterDef[];
}>();

const emit = defineEmits<{
  "update:search": [value: string];
  "update:filter": [value: ShipmentStatus | "all"];
}>();
</script>

<template>
  <div class="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
    <div class="flex items-center gap-2 rounded-lg border border-border bg-card px-3 md:w-96">
      <Search class="h-4 w-4 text-muted-foreground" />
      <Input
        :model-value="search"
        placeholder="Search by ID, tracking, customer, destination"
        class="h-11 border-0 bg-transparent font-mono text-sm shadow-none focus-visible:ring-0"
        @update:model-value="emit('update:search', String($event))"
      />
    </div>
    <div class="flex items-center gap-2 overflow-x-auto">
      <Filter class="h-4 w-4 shrink-0 text-muted-foreground" />
      <button
        v-for="f in filters"
        :key="f.key"
        @click="emit('update:filter', f.key)"
        :class="
          cn(
            'rounded-full border cursor-pointer px-3 py-1.5 font-mono text-xs uppercase tracking-wider transition-colors',
            filter === f.key
              ? 'border-primary bg-primary/15 text-primary'
              : 'border-border text-muted-foreground hover:text-foreground',
          )
        "
      >
        {{ f.label }}
      </button>
    </div>
  </div>
</template>
