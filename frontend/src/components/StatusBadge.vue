<script setup lang="ts">
import type { HTMLAttributes } from "vue";
import { type ShipmentStatus, statusLabels } from "@/lib/orders";
import Badge from "@/components/ui/Badge.vue";

const props = defineProps<{
  status: ShipmentStatus;
  class?: HTMLAttributes["class"];
}>();
</script>

<template>
  <Badge
    :variant="status"
    class="gap-1.5 font-mono text-xs uppercase tracking-wider"
    :class="props.class"
  >
    <span class="h-1.5 w-1.5 rounded-full bg-current" />
    {{ statusLabels[status] }}
  </Badge>
</template>
