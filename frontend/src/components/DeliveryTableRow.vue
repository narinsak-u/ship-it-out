<script setup lang="ts">
import { computed } from "vue";
import type { AcceptableValue } from "reka-ui";
import { Eye, Check } from "lucide-vue-next";
import { RouterLink } from "vue-router";
import { TableRow, TableCell } from "@/components/ui/table";
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { statusLabels, type ShipmentStatus } from "@/lib/orders";
import type { Hub } from "@/lib/hubs";

interface OrderRow {
  id: string;
  trackingNumber: string;
  customer: { name: string };
  carrier: string;
  status: ShipmentStatus;
  hubId?: string;
  estimatedDelivery: string;
}

const props = defineProps<{
  order: OrderRow;
  hubs: Hub[];
  isAuthenticated: boolean;
  statusDraft?: ShipmentStatus;
  hubDraft?: string;
}>();

const emit = defineEmits<{
  "update:statusDraft": [value: ShipmentStatus];
  "update:hubDraft": [value: string];
  update: [orderId: string];
}>();

const selectedStatus = computed(() => props.statusDraft ?? props.order.status);

const selectedHubId = computed(() => props.hubDraft ?? props.order.hubId ?? "");

function usesHubSelector(status: ShipmentStatus) {
  return (
    status === "departed" ||
    status === "in_transit" ||
    status === "out_for_delivery" ||
    status === "delayed"
  );
}

const canUpdate = computed(() => {
  const ds = selectedStatus.value;
  const dh = selectedHubId.value;
  const changed =
    ds !== props.order.status || (usesHubSelector(ds) && dh !== (props.order.hubId ?? ""));
  if (!changed) return false;
  if (usesHubSelector(ds) && !dh) return false;
  return true;
});
</script>

<template>
  <TableRow class="cursor-pointer border-b border-border transition-colors hover:bg-secondary/40">
    <TableCell class="font-mono text-sm text-primary">{{ order.id }}</TableCell>
    <TableCell class="font-mono text-xs text-muted-foreground">{{
      order.trackingNumber
    }}</TableCell>
    <TableCell class="font-mono text-sm">{{ order.customer.name }}</TableCell>
    <TableCell class="font-mono text-sm text-muted-foreground">{{ order.carrier }}</TableCell>
    <TableCell>
      <Select
        :model-value="selectedStatus"
        @update:model-value="(v: AcceptableValue) => emit('update:statusDraft', v as ShipmentStatus)"
        :disabled="!isAuthenticated"
      >
        <SelectTrigger
          class="h-7 rounded-lg border border-border bg-background px-2 font-mono text-xs disabled:opacity-40"
        >
          <SelectValue />
        </SelectTrigger>
        <SelectContent>
          <SelectGroup>
            <SelectItem v-for="(label, key) in statusLabels" :key="key" :value="key">
              {{ label }}
            </SelectItem>
          </SelectGroup>
        </SelectContent>
      </Select>
    </TableCell>
    <TableCell>
      <Select
        v-if="usesHubSelector(selectedStatus)"
        :model-value="selectedHubId"
        @update:model-value="(v: AcceptableValue) => emit('update:hubDraft', (v ?? '') as string)"
        :disabled="!isAuthenticated"
      >
        <SelectTrigger
          class="h-7 w-full rounded-lg border border-border bg-background px-2 font-mono text-xs disabled:opacity-40"
        >
          <SelectValue placeholder="Select hub..." />
        </SelectTrigger>
        <SelectContent>
          <SelectGroup>
            <SelectItem v-for="h in hubs" :key="h.id" :value="h.id">
              {{ h.name }}
            </SelectItem>
          </SelectGroup>
        </SelectContent>
      </Select>
      <span v-else class="font-mono text-xs text-muted-foreground">&mdash;</span>
    </TableCell>
    <TableCell class="font-mono text-xs text-muted-foreground">{{
      order.estimatedDelivery
    }}</TableCell>
    <TableCell class="flex gap-1">
      <RouterLink
        :to="{ name: 'order-detail', params: { orderId: order.id } }"
        class="rounded p-1.5 text-muted-foreground transition-colors hover:text-foreground"
        title="View details"
      >
        <Eye class="h-4 w-4" />
      </RouterLink>
      <button
        @click="emit('update', order.id)"
        :disabled="!isAuthenticated || !canUpdate"
        class="rounded cursor-pointer p-1.5 text-muted-foreground transition-colors hover:text-primary disabled:opacity-30 disabled:pointer-events-none"
        title="Update"
      >
        <Check class="h-4 w-4" />
      </button>
    </TableCell>
  </TableRow>
</template>
