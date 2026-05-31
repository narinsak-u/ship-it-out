<script setup lang="ts">
import { ArrowRight, Eye, Pencil, Trash2 } from "lucide-vue-next";
import { RouterLink } from "vue-router";
import StatusBadge from "@/components/StatusBadge.vue";
import { TableRow, TableCell } from "@/components/ui/table";
import type { Order } from "@/lib/orders";

const props = defineProps<{
  order: Order;
  isAuthenticated: boolean;
}>();

const emit = defineEmits<{
  edit: [id: string];
  delete: [id: string];
}>();
</script>

<template>
  <TableRow class="border-b border-border transition-colors hover:bg-secondary/40">
    <TableCell>
      <RouterLink
        :to="{ name: 'order-detail', params: { orderId: order.id } }"
        class="font-mono text-sm text-primary"
      >
        {{ order.id }}
      </RouterLink>
    </TableCell>
    <TableCell class="font-mono text-sm text-muted-foreground">
      {{ order.trackingNumber }}
    </TableCell>
    <TableCell class="text-sm">{{ order.customer.name }}</TableCell>
    <TableCell>
      <span class="flex items-center gap-2 font-mono text-xs text-muted-foreground">
        <span>{{ order.origin }}</span>
        <ArrowRight class="h-3 w-3 text-primary" />
        <span>{{ order.destination }}</span>
      </span>
    </TableCell>
    <TableCell><StatusBadge :status="order.status" /></TableCell>
    <TableCell class="font-mono text-xs text-muted-foreground">
      {{ order.estimatedDelivery }}
    </TableCell>
    <TableCell class="flex">
      <RouterLink
        :to="{ name: 'order-detail', params: { orderId: order.id } }"
        class="rounded p-1.5 text-muted-foreground transition-colors hover:text-foreground"
        title="View details"
      >
        <Eye class="h-4 w-4" />
      </RouterLink>
      <div v-if="isAuthenticated">
        <button
          @click.stop="emit('edit', order.id)"
          class="rounded cursor-pointer p-1.5 text-muted-foreground hover:text-primary"
        >
          <Pencil class="h-4 w-4" />
        </button>
        <button
          @click.stop="emit('delete', order.id)"
          class="rounded cursor-pointer p-1.5 text-muted-foreground hover:text-destructive"
        >
          <Trash2 class="h-4 w-4" />
        </button>
      </div>
    </TableCell>
  </TableRow>
</template>
