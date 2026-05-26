<script setup lang="ts">
import { ref, computed } from "vue";

import type { OrderFormData } from "@/lib/api/orders";
import Input from "@/components/ui/Input.vue";
import Button from "@/components/ui/Button.vue";
import ThaiAddressGroup from "./ThaiAddressGroup.vue";
import { geocodeAddress } from "@/lib/geocode";

const props = defineProps<{
  initial?: Partial<OrderFormData> & { estimatedDeliveryRaw?: string };
  isEditing?: boolean;
  pending?: boolean;
}>();

const emit = defineEmits<{
  submit: [data: OrderFormData];
  cancel: [];
}>();

const sender = ref({
  name: props.initial?.customer?.name ?? "",
  zipcode: props.initial?.customer?.zipcode ?? "",
  subDistrict: props.initial?.customer?.subDistrict ?? "",
  district: props.initial?.customer?.district ?? "",
  province: props.initial?.customer?.province ?? "",
});

const receiver = ref({
  name: props.initial?.receiver?.name ?? "",
  zipcode: props.initial?.receiver?.zipcode ?? "",
  subDistrict: props.initial?.receiver?.subDistrict ?? "",
  district: props.initial?.receiver?.district ?? "",
  province: props.initial?.receiver?.province ?? "",
});

// Parcel
const carrier = ref(props.initial?.carrier ?? "Thun-u-der Express");
const weight = ref(props.initial?.weight ?? "");
const items = ref(props.initial?.items ?? 1);
const estimatedDelivery = ref(props.initial?.estimatedDelivery ?? "");
const estimatedDeliveryRaw = ref(props.initial?.estimatedDeliveryRaw ?? "");

const errors = ref<Record<string, string>>({});
const geocodeErrors = ref<Record<string, string>>({});
const geocoding = ref(false);

const filled = computed(() => {
  return (
    sender.value.name.trim() &&
    sender.value.zipcode.trim() &&
    sender.value.subDistrict.trim() &&
    sender.value.district.trim() &&
    sender.value.province.trim() &&
    receiver.value.name.trim() &&
    receiver.value.zipcode.trim() &&
    receiver.value.subDistrict.trim() &&
    receiver.value.district.trim() &&
    receiver.value.province.trim() &&
    weight.value.trim() &&
    (items.value ?? 0) >= 1 &&
    (!props.isEditing || estimatedDelivery.value.trim())
  );
});

const formChanged = computed(() => {
  if (!props.isEditing || !props.initial) return true;
  return (
    sender.value.name !== (props.initial.customer?.name ?? "") ||
    sender.value.zipcode !== (props.initial.customer?.zipcode ?? "") ||
    sender.value.subDistrict !== (props.initial.customer?.subDistrict ?? "") ||
    sender.value.district !== (props.initial.customer?.district ?? "") ||
    sender.value.province !== (props.initial.customer?.province ?? "") ||
    receiver.value.name !== (props.initial.receiver?.name ?? "") ||
    receiver.value.zipcode !== (props.initial.receiver?.zipcode ?? "") ||
    receiver.value.subDistrict !== (props.initial.receiver?.subDistrict ?? "") ||
    receiver.value.district !== (props.initial.receiver?.district ?? "") ||
    receiver.value.province !== (props.initial.receiver?.province ?? "") ||
    weight.value !== (props.initial.weight ?? "") ||
    (items.value ?? 0) !== (props.initial.items ?? 1) ||
    estimatedDelivery.value !== (props.initial.estimatedDelivery ?? "")
  );
});

const canSubmit = computed(() => filled.value && formChanged.value);

function validate(): boolean {
  const e: Record<string, string> = {};
  if (!sender.value.name.trim()) e["sender.name"] = "Required";
  if (!sender.value.zipcode.trim()) e["sender.zipcode"] = "Required";
  if (!sender.value.subDistrict.trim()) e["sender.subDistrict"] = "Required";
  if (!sender.value.district.trim()) e["sender.district"] = "Required";
  if (!sender.value.province.trim()) e["sender.province"] = "Required";
  if (!receiver.value.name.trim()) e["receiver.name"] = "Required";
  if (!receiver.value.zipcode.trim()) e["receiver.zipcode"] = "Required";
  if (!receiver.value.subDistrict.trim()) e["receiver.subDistrict"] = "Required";
  if (!receiver.value.district.trim()) e["receiver.district"] = "Required";
  if (!receiver.value.province.trim()) e["receiver.province"] = "Required";
  if (!weight.value.trim()) e.weight = "Required";
  if (!items.value || items.value < 1) e.items = "Must be at least 1";
  if (props.isEditing && !estimatedDelivery.value.trim()) e.estimatedDelivery = "Required";
  errors.value = e;
  return Object.keys(e).length === 0;
}

const senderErrorMap = computed(() => {
  const result: Record<string, string> = {};
  for (const key of Object.keys(errors.value)) {
    if (key.startsWith("sender.")) result[key.slice(7)] = errors.value[key];
  }
  return result;
});

const receiverErrorMap = computed(() => {
  const result: Record<string, string> = {};
  for (const key of Object.keys(errors.value)) {
    if (key.startsWith("receiver.")) result[key.slice(9)] = errors.value[key];
  }
  return result;
});

async function handleSubmit() {
  if (!validate()) return;
  geocodeErrors.value = {};
  geocoding.value = true;

  try {
    const [senderCoords, receiverCoords] = await Promise.all([
      geocodeAddress(
        sender.value.subDistrict,
        sender.value.district,
        sender.value.province,
      ),
      geocodeAddress(
        receiver.value.subDistrict,
        receiver.value.district,
        receiver.value.province,
      ),
    ]);

    emit("submit", {
      customer: {
        name: sender.value.name,
        zipcode: sender.value.zipcode,
        subDistrict: sender.value.subDistrict,
        district: sender.value.district,
        province: sender.value.province,
        coords: senderCoords,
      },
      receiver: {
        name: receiver.value.name,
        zipcode: receiver.value.zipcode,
        subDistrict: receiver.value.subDistrict,
        district: receiver.value.district,
        province: receiver.value.province,
        coords: receiverCoords,
      },
      carrier: carrier.value,
      weight: weight.value,
      items: items.value,
      estimatedDelivery: estimatedDeliveryRaw.value || "",
    });
  } catch (e) {
    const msg = e instanceof Error ? e.message : "Could not resolve address.";
    geocodeErrors.value = { sender: msg, receiver: msg };
  } finally {
    geocoding.value = false;
  }
}
</script>

<template>
  <form @submit.prevent="handleSubmit" class="space-y-8">
    <ThaiAddressGroup label="Sender Info" v-model="sender" :errors="senderErrorMap" />
    <p
      v-if="geocodeErrors.sender"
      class="mt-2 rounded-md bg-destructive/15 px-3 py-2 font-mono text-xs text-destructive"
    >
      {{ geocodeErrors.sender }}
    </p>

    <ThaiAddressGroup label="Receiver Info" v-model="receiver" :errors="receiverErrorMap" />
    <p
      v-if="geocodeErrors.receiver"
      class="mt-2 rounded-md bg-destructive/15 px-3 py-2 font-mono text-xs text-destructive"
    >
      {{ geocodeErrors.receiver }}
    </p>

    <!-- Section 3: Parcel Info -->
    <fieldset class="rounded-xl border border-border p-5">
      <legend class="font-mono text-xs uppercase tracking-widest text-primary px-2">
        Parcel Info
      </legend>
      <div class="grid gap-5 md:grid-cols-2">
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
            >Weight</label
          >
          <Input v-model="weight" class="mt-1.5 font-mono text-sm" placeholder="e.g. 12.4 kg" />
          <p v-if="errors.weight" class="mt-1 font-mono text-xs text-destructive">
            {{ errors.weight }}
          </p>
        </div>
        <div>
          <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground"
            >Items</label
          >
          <Input v-model.number="items" type="number" min="1" class="mt-1.5 font-mono text-sm" />
          <p v-if="errors.items" class="mt-1 font-mono text-xs text-destructive">
            {{ errors.items }}
          </p>
        </div>
        <div v-if="isEditing">
          <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground"
            >Estimated Delivery</label
          >
          <Input
            v-model="estimatedDelivery"
            class="mt-1.5 font-mono text-sm"
            placeholder="e.g. May 25, 2026"
          />
          <p v-if="errors.estimatedDelivery" class="mt-1 font-mono text-xs text-destructive">
            {{ errors.estimatedDelivery }}
          </p>
        </div>
      </div>
    </fieldset>

    <div class="flex justify-end gap-3 pt-4">
      <Button variant="outline" type="button" @click="emit('cancel')">Cancel</Button>
      <Button type="submit" :disabled="pending || !canSubmit || geocoding">
        {{
          geocoding
            ? "Resolving addresses\u2026"
            : pending
              ? "Saving\u2026"
              : isEditing
                ? "Save Changes"
                : "Create Order"
        }}
      </Button>
    </div>
  </form>
</template>
