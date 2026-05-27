<script setup lang="ts">
import { ref, computed, watch } from "vue";
import { toast } from "vue-sonner";
import { geocodeAddress } from "@/lib/geocode";
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

const PROVINCES = [
  "Amnat Charoen",
  "Bueng Kan",
  "Chaiyaphum",
  "Kalasin",
  "Loei",
  "Maha Sarakham",
  "Mukdahan",
  "Nakhon Phanom",
  "Nong Bua Lamphu",
  "Nong Khai",
  "Roi Et",
  "Sakon Nakhon",
  "Sisaket",
  "Surin",
  "Yasothon",
] as const;

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
const location = ref("");
const capacity = ref(1000);
const status = ref<HubStatus>("active");
const geocodeError = ref("");
const geocoding = ref(false);

function resetForm() {
  name.value = "";
  location.value = "";
  capacity.value = 1000;
  status.value = "active";
  geocodeError.value = "";
}

watch(existing, (hub) => {
  if (hub) {
    name.value = hub.name;
    location.value = hub.address;
    capacity.value = hub.capacity;
    status.value = hub.status;
  }
});

watch(
  () => props.open,
  (isOpen) => {
    if (isOpen && !props.hubId) {
      resetForm();
    }
  },
);

const isEditing = computed(() => !!props.hubId);

const submitPending = computed(() =>
  isEditing.value ? updateHub.isPending.value : createHub.isPending.value,
);

const submitError = computed(() =>
  isEditing.value ? updateHub.isError.value : createHub.isError.value,
);

async function handleSubmit() {
  geocodeError.value = "";
  geocoding.value = true;

  let coords: { lat: number; lng: number };
  if (isEditing.value && existing.value) {
    coords = existing.value.coords;
    geocoding.value = false;
  } else {
    try {
      coords = await geocodeAddress("", "", location.value);
    } catch (e) {
      geocodeError.value = e instanceof Error ? e.message : "Could not resolve location.";
      geocoding.value = false;
      return;
    }
  }

  const data = {
    name: name.value,
    carrierId,
    address: location.value,
    coords,
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
  } finally {
    geocoding.value = false;
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
          <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground">
            Location
          </label>
          <div v-if="isEditing">
            <div
              class="mt-1.5 flex h-10 w-full items-center rounded-lg border border-border bg-background px-3 font-mono text-sm text-muted-foreground"
            >
              {{ location || "Not set" }}
            </div>
          </div>
          <div v-else>
            <Select v-model="location">
              <SelectTrigger
                class="mt-1.5 flex h-10 w-full rounded-lg border border-border bg-background px-3 font-mono text-sm"
              >
                <SelectValue placeholder="Select province..." />
              </SelectTrigger>
              <SelectContent>
                <SelectGroup>
                  <SelectItem v-for="p in PROVINCES" :key="p" :value="p">
                    {{ p }}
                  </SelectItem>
                </SelectGroup>
              </SelectContent>
            </Select>
          </div>
          <p
            v-if="geocodeError"
            class="mt-1.5 rounded-md bg-destructive/15 px-3 py-2 font-mono text-xs text-destructive"
          >
            {{ geocodeError }}
          </p>
        </div>
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground"
              >Capacity (units)</label
            >
            <Input v-model.number="capacity" type="number" class="mt-1.5 font-mono text-sm" />
          </div>
          <div>
            <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground">
              Status
            </label>
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
        <Button :disabled="!name || !location || submitPending || geocoding" @click="handleSubmit">
          {{
            geocoding
              ? "Resolving location\u2026"
              : submitPending
                ? "Saving\u2026"
                : isEditing
                  ? "Update Hub"
                  : "Create Hub"
          }}
        </Button>
      </div>
    </DialogContent>
  </Dialog>
</template>
