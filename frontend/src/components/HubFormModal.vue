<script setup lang="ts">
import { ref, computed, watch } from "vue";
import { toast } from "vue-sonner";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { useHubs, useCreateHub, useUpdateHub } from "@/hooks/useHubs";
import { hubStatusLabels, type HubStatus } from "@/lib/carriers";
import Input from "@/components/ui/Input.vue";
import Button from "@/components/ui/Button.vue";
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

const props = defineProps<{ hubId?: string | null; open?: boolean }>();
const emit = defineEmits<{ close: [] }>();

const { data: hubsData } = useHubs();
const createHub = useCreateHub();
const updateHub = useUpdateHub();

const existing = computed(() => {
  if (!props.hubId || !hubsData.value) return null;
  return hubsData.value.find((h) => h.id === props.hubId) ?? null;
});

const name = ref("");
const carrierId = "THUN";
const address = ref("");
const capacity = ref(1000);
const status = ref<HubStatus>("active");

watch(existing, (hub) => {
  if (hub) {
    name.value = hub.name;
    address.value = hub.address;
    capacity.value = hub.capacity;
    status.value = hub.status;
  }
});

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
      toast.success("Hub updated");
    } else {
      await createHub.mutateAsync(data);
      toast.success("Hub created");
    }
    emit("close");
  } catch {
    // Mutation error is surfaced via submitError; modal stays open
  }
}
</script>

<template>
  <Dialog :open="open" @update:open="(v) => !v && emit('close')">
    <DialogContent
      class="w-full max-w-md rounded-xl border border-border bg-card p-6 shadow-elegant sm:max-w-md"
    >
      <DialogHeader>
        <DialogTitle class="font-mono text-lg font-semibold">
          {{ isEditing ? "Edit Hub" : "Add Hub" }}
        </DialogTitle>
      </DialogHeader>

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
            <Select v-model="status">
              <SelectTrigger
                class="mt-1.5 flex h-10 w-full rounded-lg border border-border bg-background px-3 font-mono text-sm"
              >
                <SelectValue placeholder="Select status..." />
              </SelectTrigger>
              <SelectContent>
                <SelectGroup>
                  <SelectItem v-for="(label, key) in hubStatusLabels" :key="key" :value="key">
                    {{ label }}
                  </SelectItem>
                </SelectGroup>
              </SelectContent>
            </Select>
          </div>
        </div>
      </div>

      <div v-if="submitError" class="mt-3 font-mono text-xs text-destructive">
        Failed to save hub. Please try again.
      </div>

      <div class="mt-6 flex justify-end gap-3">
        <Button variant="outline" @click="emit('close')">Cancel</Button>
        <Button :disabled="!name || submitPending" @click="handleSubmit">
          {{ submitPending ? "Saving…" : isEditing ? "Update Hub" : "Create Hub" }}
        </Button>
      </div>
    </DialogContent>
  </Dialog>
</template>
