<script setup lang="ts">
import { ref, computed } from "vue";
import { Search, X } from "lucide-vue-next";
import { useDrivers, useAssignDriver } from "@/hooks/useDrivers";
import { orders } from "@/lib/orders";
import { getCarrier } from "@/lib/carriers";
import Input from "@/components/ui/Input.vue";
import Button from "@/components/ui/Button.vue";

const emit = defineEmits<{ close: [] }>();

const { data: drivers } = useDrivers();
const assignMutation = useAssignDriver();

const driverQuery = ref("");
const selectedDriverId = ref("");
const selectedOrderId = ref("");

const availableDrivers = computed(() => {
  if (!drivers.value) return [];
  const q = driverQuery.value.trim().toLowerCase();
  return drivers.value.filter((d) => {
    if (d.status !== "available") return false;
    if (!q) return true;
    return d.name.toLowerCase().includes(q) || d.vehicleInfo.toLowerCase().includes(q);
  });
});

const unfinishedOrders = computed(() =>
  orders.filter((o) => o.status !== "delivered" && !o.driverId),
);

const canSubmit = computed(() => selectedDriverId.value && selectedOrderId.value);

async function handleAssign() {
  if (!canSubmit.value) return;
  await assignMutation.mutateAsync({
    driverId: selectedDriverId.value,
    orderId: selectedOrderId.value,
  });
  emit("close");
}
</script>

<template>
  <div class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm" @click.self="emit('close')">
    <div class="w-full max-w-lg rounded-xl border border-border bg-card p-6 shadow-elegant">
      <div class="flex items-center justify-between">
        <h2 class="font-mono text-lg font-semibold">Assign Driver to Shipment</h2>
        <button @click="emit('close')" class="text-muted-foreground hover:text-foreground">
          <X class="h-5 w-5" />
        </button>
      </div>

      <!-- Driver selection -->
      <div class="mt-5">
        <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground">Select Driver</label>
        <div class="mt-2 flex items-center gap-2 rounded-lg border border-border bg-background px-3">
          <Search class="h-4 w-4 text-muted-foreground" />
          <Input
            v-model="driverQuery"
            placeholder="Search available drivers..."
            class="h-10 border-0 bg-transparent font-mono text-sm shadow-none focus-visible:ring-0"
          />
        </div>
        <div class="mt-2 max-h-40 overflow-y-auto space-y-1">
          <button
            v-for="d in availableDrivers"
            :key="d.id"
            @click="selectedDriverId = d.id"
            class="w-full rounded-lg px-3 py-2 text-left font-mono text-sm transition-colors"
            :class="selectedDriverId === d.id ? 'bg-primary/15 text-primary' : 'hover:bg-secondary'"
          >
            <div>{{ d.name }}</div>
            <div class="text-xs text-muted-foreground">{{ d.vehicleInfo }} · {{ getCarrier(d.carrierId)?.name }}</div>
          </button>
          <div v-if="availableDrivers.length === 0" class="py-4 text-center font-mono text-xs text-muted-foreground">
            No available drivers found.
          </div>
        </div>
      </div>

      <!-- Order selection -->
      <div class="mt-5">
        <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground">Target Shipment</label>
        <div class="mt-2 space-y-1 max-h-40 overflow-y-auto">
          <button
            v-for="o in unfinishedOrders"
            :key="o.id"
            @click="selectedOrderId = o.id"
            class="w-full rounded-lg px-3 py-2 text-left font-mono text-sm transition-colors"
            :class="selectedOrderId === o.id ? 'bg-primary/15 text-primary' : 'hover:bg-secondary'"
          >
            <div>{{ o.id }} — {{ o.customer }}</div>
            <div class="text-xs text-muted-foreground">{{ o.origin }} → {{ o.destination }}</div>
          </button>
        </div>
      </div>

      <div v-if="assignMutation.isError" class="mt-3 font-mono text-xs text-destructive">
        Failed to assign driver. Please try again.
      </div>

      <div class="mt-6 flex justify-end gap-3">
        <Button variant="outline" @click="emit('close')">Cancel</Button>
        <Button :disabled="!canSubmit || assignMutation.isPending" @click="handleAssign">
          {{ assignMutation.isPending ? "Assigning…" : "Assign Driver" }}
        </Button>
      </div>
    </div>
  </div>
</template>
