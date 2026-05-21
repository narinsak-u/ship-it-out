<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { type ShipmentStatus, statusLabels } from '@/lib/orders'
import { cn } from '@/lib/utils'

const props = defineProps<{
  status: ShipmentStatus
  class?: HTMLAttributes['class']
}>()

const styles: Record<ShipmentStatus, string> = {
  pending: 'bg-muted text-muted-foreground border-border',
  in_transit: 'bg-info/15 text-info border-info/30',
  out_for_delivery: 'bg-primary/15 text-primary border-primary/30',
  delivered: 'bg-success/15 text-success border-success/30',
  delayed: 'bg-destructive/15 text-destructive border-destructive/30',
}
</script>

<template>
  <span
    :class="cn(
      'inline-flex items-center gap-1.5 rounded-full border px-2.5 py-0.5 text-xs font-mono uppercase tracking-wider',
      styles[status],
      props.class,
    )"
  >
    <span class="h-1.5 w-1.5 rounded-full bg-current" />
    {{ statusLabels[status] }}
  </span>
</template>
