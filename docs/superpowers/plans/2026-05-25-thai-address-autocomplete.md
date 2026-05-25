# Thai Address Autocomplete Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) for syntax tracking.

**Goal:** Add zipcode-driven Thai address autocomplete to the order form using the `thai-data` library.

**Architecture:** Extract a reusable `ThaiAddressGroup.vue` component encapsulating zipcode→sub-district dropdown→auto-fill district/province. Use it twice (sender + receiver) in `OrderForm.vue`. All `thai-data` calls are synchronous against bundled JSON — no network requests.

**Tech Stack:** Vue 3 (Composition API, `<script setup lang="ts">`), shadcn-vue `<Input>`, native `<select>`, tailwind-data

---

### File Structure

```
frontend/src/components/
├── ThaiAddressGroup.vue    # CREATE — reusable Thai address autocomplete block
├── OrderForm.vue            # MODIFY — replace sender/receiver <fieldset> blocks
└── ui/
    └── Input.vue            # UNCHANGED — existing shadcn-vue input
```

---

### Task 1: Create ThaiAddressGroup.vue

**Files:**
- Create: `frontend/src/components/ThaiAddressGroup.vue`

This component encapsulates one address block (sender or receiver) with:
- Name (free-text `<Input>`)
- Zipcode (`<Input>` with maxlength=5)
- Sub-district (native `<select>` that appears on 5-digit zipcode)
- District (auto-filled `<Input disabled>`)
- Province (auto-filled `<Input disabled>`)

- [ ] **Step 1: Create the component scaffold**

```vue
<script setup lang="ts">
import { ref, watch, computed } from "vue";
import {
  getSubDistrictNames,
  getDistrictNames,
  getProvinceName,
} from "thai-data";
import Input from "@/components/ui/Input.vue";

export interface ThaiAddress {
  name: string;
  zipcode: string;
  subDistrict: string;
  district: string;
  province: string;
}

const props = defineProps<{
  label: string;
  modelValue: ThaiAddress;
  errors?: Record<string, string>;
}>();

const emit = defineEmits<{
  "update:modelValue": [value: ThaiAddress];
}>();

const availableSubDistricts = ref<string[]>([]);

function patch(field: keyof ThaiAddress, value: string) {
  emit("update:modelValue", { ...props.modelValue, [field]: value });
}

const districtDisplay = computed(() => {
  if (props.modelValue.district) return props.modelValue.district;
  const names = getDistrictNames(props.modelValue.zipcode);
  return names.length > 0 ? names[0] : "";
});

const provinceDisplay = computed(() => {
  if (props.modelValue.province) return props.modelValue.province;
  return getProvinceName(props.modelValue.zipcode) ?? "";
});

watch(
  () => props.modelValue.zipcode,
  (zip, oldZip) => {
    if (zip.length === 5) {
      availableSubDistricts.value = getSubDistrictNames(zip);
      // If zipcode changed after a prior lookup, clear previous selection
      if (oldZip && zip !== oldZip) {
        patch("subDistrict", "");
        patch("district", "");
        patch("province", "");
      }
    } else {
      availableSubDistricts.value = [];
    }
  },
  { immediate: true },
);

// Auto-select when exactly 1 sub-district matches
watch(availableSubDistricts, (list) => {
  if (list.length === 1 && !props.modelValue.subDistrict) {
    const zip = props.modelValue.zipcode;
    patch("subDistrict", list[0]);
    const districts = getDistrictNames(zip);
    const province = getProvinceName(zip);
    if (districts.length > 0) patch("district", districts[0]);
    if (province) patch("province", province);
  }
});

function onSubDistrictSelected(ev: Event) {
  const sub = (ev.target as HTMLSelectElement).value;
  if (!sub) return;
  const zip = props.modelValue.zipcode;
  patch("subDistrict", sub);
  const districts = getDistrictNames(zip);
  const province = getProvinceName(zip);
  if (districts.length > 0) patch("district", districts[0]);
  if (province) patch("province", province);
}
</script>

<template>
  <fieldset class="rounded-xl border border-border p-5">
    <legend class="font-mono text-xs uppercase tracking-widest text-primary px-2">
      {{ label }}
    </legend>
    <div class="grid gap-5 md:grid-cols-2">
      <!-- Name -->
      <div>
        <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground"
          >Name</label
        >
        <Input
          :value="modelValue.name"
          class="mt-1.5 font-mono text-sm"
          placeholder="e.g. ประวิทย์ ใจดี"
          @input="patch('name', ($event.target as HTMLInputElement).value)"
        />
        <p v-if="errors?.name" class="mt-1 font-mono text-xs text-destructive">
          {{ errors.name }}
        </p>
      </div>

      <!-- Zipcode -->
      <div>
        <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground"
          >Zipcode</label
        >
        <Input
          :value="modelValue.zipcode"
          class="mt-1.5 font-mono text-sm"
          maxlength="5"
          placeholder="e.g. 10200"
          @input="patch('zipcode', ($event.target as HTMLInputElement).value)"
        />
        <p v-if="errors?.zipcode" class="mt-1 font-mono text-xs text-destructive">
          {{ errors.zipcode }}
        </p>
      </div>

      <!-- Sub-district -->
      <div>
        <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground"
          >Sub-district</label
        >
        <select
          v-if="availableSubDistricts.length"
          :value="modelValue.subDistrict"
          class="mt-1.5 flex h-10 w-full rounded-lg border border-border bg-background px-3 font-mono text-sm"
          @change="onSubDistrictSelected"
        >
          <option disabled value="">Select sub-district...</option>
          <option
            v-for="sd in availableSubDistricts"
            :key="sd"
            :value="sd"
          >
            {{ sd }}
          </option>
        </select>
        <Input
          v-else
          disabled
          class="mt-1.5 font-mono text-sm"
          placeholder="Enter zipcode first"
        />
        <p v-if="errors?.subDistrict" class="mt-1 font-mono text-xs text-destructive">
          {{ errors.subDistrict }}
        </p>
      </div>

      <!-- District -->
      <div>
        <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground"
          >District</label
        >
        <Input :value="districtDisplay" disabled class="mt-1.5 font-mono text-sm" />
        <p v-if="errors?.district" class="mt-1 font-mono text-xs text-destructive">
          {{ errors.district }}
        </p>
      </div>

      <!-- Province -->
      <div>
        <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground"
          >Province</label
        >
        <Input :value="provinceDisplay" disabled class="mt-1.5 font-mono text-sm" />
        <p v-if="errors?.province" class="mt-1 font-mono text-xs text-destructive">
          {{ errors.province }}
        </p>
      </div>
    </div>
  </fieldset>
</template>
```

Write this file to `frontend/src/components/ThaiAddressGroup.vue`.

- [ ] **Step 2: Verify it compiles**

Run: `cd frontend && npm run build`
Expected: Build succeeds with no errors. (Component isn't used yet so no change in behavior.)

- [ ] **Step 3: Commit**

```bash
git add frontend/src/components/ThaiAddressGroup.vue
git commit -m "feat: add ThaiAddressGroup component for zipcode-driven address autocomplete"
```

---

### Task 2: Update OrderForm.vue

**Files:**
- Modify: `frontend/src/components/OrderForm.vue`

Replace the two `<fieldset>` blocks (Sender Info, Receiver Info) with two `<ThaiAddressGroup>` instances. The parcel info section stays untouched.

- [ ] **Step 1: Update the `<script>` section**

Add the import and new reactive state. Remove the 10 individual refs for sender/receiver fields and replace with two refs that match the ThaiAddress interface:

```vue
<script setup lang="ts">
import { ref } from "vue";
import { statusLabels, type ShipmentStatus } from "@/lib/orders";
import type { OrderFormData } from "@/lib/api/orders";
import Input from "@/components/ui/Input.vue";
import Button from "@/components/ui/Button.vue";
import ThaiAddressGroup from "./ThaiAddressGroup.vue";

const props = defineProps<{
  initial?: Partial<OrderFormData & { status?: ShipmentStatus }>;
  isEditing?: boolean;
  pending?: boolean;
}>();

const emit = defineEmits<{
  submit: [data: OrderFormData & { status?: ShipmentStatus }];
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
const status = ref<ShipmentStatus>(props.initial?.status ?? "pending");

const errors = ref<Record<string, string>>({});

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
  if (!estimatedDelivery.value.trim()) e.estimatedDelivery = "Required";
  errors.value = e;
  return Object.keys(e).length === 0;
}

function senderErrors(errors: Record<string, string>): Record<string, string> {
  const result: Record<string, string> = {};
  for (const key of Object.keys(errors)) {
    if (key.startsWith("sender.")) result[key.slice(7)] = errors[key];
  }
  return result;
}

function receiverErrors(errors: Record<string, string>): Record<string, string> {
  const result: Record<string, string> = {};
  for (const key of Object.keys(errors)) {
    if (key.startsWith("receiver.")) result[key.slice(9)] = errors[key];
  }
  return result;
}

function handleSubmit() {
  if (!validate()) return;
  emit("submit", {
    customer: {
      name: sender.value.name,
      zipcode: sender.value.zipcode,
      subDistrict: sender.value.subDistrict,
      district: sender.value.district,
      province: sender.value.province,
      coords: { lat: 0, lng: 0 },
    },
    receiver: {
      name: receiver.value.name,
      zipcode: receiver.value.zipcode,
      subDistrict: receiver.value.subDistrict,
      district: receiver.value.district,
      province: receiver.value.province,
      coords: { lat: 0, lng: 0 },
    },
    carrier: carrier.value,
    weight: weight.value,
    items: items.value,
    estimatedDelivery: estimatedDelivery.value,
    ...(props.isEditing ? { status: status.value } : {}),
  });
}
</script>
```

- [ ] **Step 2: Replace the `<template>` sections for Sender and Receiver**

Delete the two existing `<fieldset>` blocks (lines 92–233 in the current file) and replace with:

```vue
<template>
  <form @submit.prevent="handleSubmit" class="space-y-8">
    <ThaiAddressGroup
      label="Sender Info"
      v-model="sender"
      :errors="senderErrors(errors)"
    />

    <ThaiAddressGroup
      label="Receiver Info"
      v-model="receiver"
      :errors="receiverErrors(errors)"
    />

    <!-- Section 3: Parcel Info -->
    <fieldset>...</fieldset>

    <div class="flex justify-end gap-3 pt-4 border-t border-border">
      <Button variant="outline" type="button" @click="emit('cancel')">Cancel</Button>
      <Button type="submit" :disabled="pending">
        {{ pending ? "Saving\u2026" : isEditing ? "Save Changes" : "Create Order" }}
      </Button>
    </div>
  </form>
</template>
```

The Parcel Info section stays exactly as it is in the current file.

- [ ] **Step 3: Verify the build**

Run: `cd frontend && npm run build`
Expected: Build succeeds with no errors.

- [ ] **Step 4: Commit**

```bash
git add frontend/src/components/OrderForm.vue
git commit -m "feat: integrate ThaiAddressGroup into OrderForm for sender/receiver addresses"
```
