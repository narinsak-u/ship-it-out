<script setup lang="ts">
import { ref, computed } from "vue";
import { X } from "lucide-vue-next";
import { useHubs, useCreateHub, useUpdateHub } from "@/hooks/useHubs";
import { hubStatusLabels, type HubStatus } from "@/lib/carriers";
import Input from "@/components/ui/Input.vue";
import Button from "@/components/ui/Button.vue";

const props = defineProps<{ hubId?: string | null }>();
const emit = defineEmits<{ close: [] }>();

const { data: hubsData } = useHubs();
const createHub = useCreateHub();
const updateHub = useUpdateHub();

const existing = computed(() => {
  if (!props.hubId || !hubsData.value) return null;
  return hubsData.value.find((h) => h.id === props.hubId) ?? null;
});

const name = ref(existing.value?.name ?? "");
const carrierId = "THUN";
const address = ref(existing.value?.address ?? "");
const capacity = ref(existing.value?.capacity ?? 1000);
const status = ref<HubStatus>(existing.value?.status ?? "active");

const isEditing = computed(() => !!props.hubId);

const submitPending = computed(() =>
  isEditing.value ? updateHub.isPending.value : createHub.isPending.value,
);

const submitError = computed(() =>
  isEditing.value ? updateHub.isError.value : createHub.isError.value,
);

async function handleSubmit() {
  const data = {
    name: name.value,
    carrierId,
    address: address.value,
    coords: { lat: 0, lng: 0 },
    capacity: capacity.value,
    currentUtilization: existing.value?.currentUtilization ?? 0,
    status: status.value,
  };

  try {
    if (isEditing.value && props.hubId) {
      await updateHub.mutateAsync({ id: props.hubId, data });
    } else {
      await createHub.mutateAsync(data);
    }
    emit("close");
  } catch {
    // Mutation error is surfaced via submitError; modal stays open
  }
}
</script>

<template>
  <div
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm"
    @click.self="emit('close')"
  >
    <div class="w-full max-w-md rounded-xl border border-border bg-card p-6 shadow-elegant">
      <div class="flex items-center justify-between">
        <h2 class="font-mono text-lg font-semibold">{{ isEditing ? "Edit Hub" : "Add Hub" }}</h2>
        <button @click="emit('close')" class="text-muted-foreground hover:text-foreground">
          <X class="h-5 w-5" />
        </button>
      </div>

      <div class="mt-5 space-y-4">
        <div>
          <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground"
            >Name</label
          >
          <Input v-model="name" class="mt-1.5 font-mono text-sm" placeholder="Hub name" />
        </div>
        <div>
          <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground"
            >Carrier</label
          >
          <div
            class="mt-1.5 flex h-10 w-full items-center rounded-lg border border-border bg-background px-3 font-mono text-sm text-muted-foreground"
          >
            Thun-u-der Express
          </div>
        </div>
        <div>
          <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground"
            >Address</label
          >
          <Input v-model="address" class="mt-1.5 font-mono text-sm" placeholder="Full address" />
        </div>
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground"
              >Capacity (units)</label
            >
            <Input v-model.number="capacity" type="number" class="mt-1.5 font-mono text-sm" />
          </div>
          <div>
            <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground"
              >Status</label
            >
            <select
              v-model="status"
              class="mt-1.5 flex h-10 w-full rounded-lg border border-border bg-background px-3 font-mono text-sm"
            >
              <option v-for="(label, key) in hubStatusLabels" :key="key" :value="key">
                {{ label }}
              </option>
            </select>
          </div>
        </div>
      </div>

      <div
        v-if="submitError"
        class="mt-3 font-mono text-xs text-destructive"
      >
        Failed to save hub. Please try again.
      </div>

      <div class="mt-6 flex justify-end gap-3">
        <Button variant="outline" @click="emit('close')">Cancel</Button>
        <Button
          :disabled="!name || submitPending"
          @click="handleSubmit"
        >
          {{ submitPending ? "Saving…" : isEditing ? "Update Hub" : "Create Hub" }}
        </Button>
      </div>
    </div>
  </div>
</template>
