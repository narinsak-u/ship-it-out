<script setup lang="ts">
import { ref } from 'vue'
import { statusLabels, type ShipmentStatus } from '@/lib/orders'
import type { OrderFormData } from '@/lib/api/orders'
import Input from '@/components/ui/Input.vue'
import Button from '@/components/ui/Button.vue'

const props = defineProps<{
  initial?: Partial<OrderFormData & { status?: ShipmentStatus }>
  isEditing?: boolean
  pending?: boolean
}>()

const emit = defineEmits<{
  submit: [data: OrderFormData & { status?: ShipmentStatus }]
  cancel: []
}>()

// Sender
const senderName = ref(props.initial?.customer?.name ?? '')
const senderZipcode = ref(props.initial?.customer?.zipcode ?? '')
const senderSubDistrict = ref(props.initial?.customer?.subDistrict ?? '')
const senderDistrict = ref(props.initial?.customer?.district ?? '')
const senderProvince = ref(props.initial?.customer?.province ?? '')

// Receiver
const receiverName = ref(props.initial?.receiver?.name ?? '')
const receiverZipcode = ref(props.initial?.receiver?.zipcode ?? '')
const receiverSubDistrict = ref(props.initial?.receiver?.subDistrict ?? '')
const receiverDistrict = ref(props.initial?.receiver?.district ?? '')
const receiverProvince = ref(props.initial?.receiver?.province ?? '')

// Parcel
const carrier = ref(props.initial?.carrier ?? 'Thun-u-der Express')
const weight = ref(props.initial?.weight ?? '')
const items = ref(props.initial?.items ?? 1)
const estimatedDelivery = ref(props.initial?.estimatedDelivery ?? '')
const status = ref<ShipmentStatus>(props.initial?.status ?? 'pending')

const errors = ref<Record<string, string>>({})

function validate(): boolean {
  const e: Record<string, string> = {}
  if (!senderName.value.trim()) e.senderName = 'Required'
  if (!senderZipcode.value.trim()) e.senderZipcode = 'Required'
  if (!senderSubDistrict.value.trim()) e.senderSubDistrict = 'Required'
  if (!senderDistrict.value.trim()) e.senderDistrict = 'Required'
  if (!senderProvince.value.trim()) e.senderProvince = 'Required'
  if (!receiverName.value.trim()) e.receiverName = 'Required'
  if (!receiverZipcode.value.trim()) e.receiverZipcode = 'Required'
  if (!receiverSubDistrict.value.trim()) e.receiverSubDistrict = 'Required'
  if (!receiverDistrict.value.trim()) e.receiverDistrict = 'Required'
  if (!receiverProvince.value.trim()) e.receiverProvince = 'Required'
  // carrier is fixed
  if (!weight.value.trim()) e.weight = 'Required'
  if (!items.value || items.value < 1) e.items = 'Must be at least 1'
  if (!estimatedDelivery.value.trim()) e.estimatedDelivery = 'Required'
  errors.value = e
  return Object.keys(e).length === 0
}

function handleSubmit() {
  if (!validate()) return
  emit('submit', {
    customer: {
      name: senderName.value,
      zipcode: senderZipcode.value,
      subDistrict: senderSubDistrict.value,
      district: senderDistrict.value,
      province: senderProvince.value,
      coords: { lat: 0, lng: 0 },
    },
    receiver: {
      name: receiverName.value,
      zipcode: receiverZipcode.value,
      subDistrict: receiverSubDistrict.value,
      district: receiverDistrict.value,
      province: receiverProvince.value,
      coords: { lat: 0, lng: 0 },
    },
    carrier: carrier.value,
    weight: weight.value,
    items: items.value,
    estimatedDelivery: estimatedDelivery.value,
    ...(props.isEditing ? { status: status.value } : {}),
  })
}
</script>

<template>
  <form @submit.prevent="handleSubmit" class="space-y-8">
    <!-- Section 1: Sender Info -->
    <fieldset class="rounded-xl border border-border p-5">
      <legend class="font-mono text-xs uppercase tracking-widest text-primary px-2">
        &#x1f9d1;&#x200d;&#x1f4ed; Sender Info
      </legend>
      <div class="grid gap-5 md:grid-cols-2">
        <div>
          <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground">Name</label>
          <Input v-model="senderName" class="mt-1.5 font-mono text-sm" placeholder="e.g. Aria Nakamura" />
          <p v-if="errors.senderName" class="mt-1 font-mono text-xs text-destructive">{{ errors.senderName }}</p>
        </div>
        <div>
          <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground">Zipcode</label>
          <Input v-model="senderZipcode" class="mt-1.5 font-mono text-sm" placeholder="e.g. 3011" />
          <p v-if="errors.senderZipcode" class="mt-1 font-mono text-xs text-destructive">{{ errors.senderZipcode }}</p>
        </div>
        <div>
          <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground">Sub-district</label>
          <Input v-model="senderSubDistrict" class="mt-1.5 font-mono text-sm" placeholder="e.g. Stadsdriehoek" />
          <p v-if="errors.senderSubDistrict" class="mt-1 font-mono text-xs text-destructive">{{ errors.senderSubDistrict }}</p>
        </div>
        <div>
          <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground">District</label>
          <Input v-model="senderDistrict" class="mt-1.5 font-mono text-sm" placeholder="e.g. Centrum" />
          <p v-if="errors.senderDistrict" class="mt-1 font-mono text-xs text-destructive">{{ errors.senderDistrict }}</p>
        </div>
        <div>
          <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground">Province</label>
          <Input v-model="senderProvince" class="mt-1.5 font-mono text-sm" placeholder="e.g. Zuid-Holland" />
          <p v-if="errors.senderProvince" class="mt-1 font-mono text-xs text-destructive">{{ errors.senderProvince }}</p>
        </div>
      </div>
    </fieldset>

    <!-- Section 2: Receiver Info -->
    <fieldset class="rounded-xl border border-border p-5">
      <legend class="font-mono text-xs uppercase tracking-widest text-primary px-2">
        &#x1f9d1;&#x200d;&#x1f381; Receiver Info
      </legend>
      <div class="grid gap-5 md:grid-cols-2">
        <div>
          <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground">Name</label>
          <Input v-model="receiverName" class="mt-1.5 font-mono text-sm" placeholder="e.g. James Mitchell" />
          <p v-if="errors.receiverName" class="mt-1 font-mono text-xs text-destructive">{{ errors.receiverName }}</p>
        </div>
        <div>
          <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground">Zipcode</label>
          <Input v-model="receiverZipcode" class="mt-1.5 font-mono text-sm" placeholder="e.g. 11201" />
          <p v-if="errors.receiverZipcode" class="mt-1 font-mono text-xs text-destructive">{{ errors.receiverZipcode }}</p>
        </div>
        <div>
          <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground">Sub-district</label>
          <Input v-model="receiverSubDistrict" class="mt-1.5 font-mono text-sm" placeholder="e.g. DUMBO" />
          <p v-if="errors.receiverSubDistrict" class="mt-1 font-mono text-xs text-destructive">{{ errors.receiverSubDistrict }}</p>
        </div>
        <div>
          <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground">District</label>
          <Input v-model="receiverDistrict" class="mt-1.5 font-mono text-sm" placeholder="e.g. Brooklyn" />
          <p v-if="errors.receiverDistrict" class="mt-1 font-mono text-xs text-destructive">{{ errors.receiverDistrict }}</p>
        </div>
        <div>
          <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground">Province</label>
          <Input v-model="receiverProvince" class="mt-1.5 font-mono text-sm" placeholder="e.g. New York" />
          <p v-if="errors.receiverProvince" class="mt-1 font-mono text-xs text-destructive">{{ errors.receiverProvince }}</p>
        </div>
      </div>
    </fieldset>

    <!-- Section 3: Parcel Info -->
    <fieldset class="rounded-xl border border-border p-5">
      <legend class="font-mono text-xs uppercase tracking-widest text-primary px-2">
        &#x1f4e6; Parcel Info
      </legend>
      <div class="grid gap-5 md:grid-cols-2">
        <div>
          <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground">Carrier</label>
          <div class="mt-1.5 flex h-10 w-full items-center rounded-lg border border-border bg-background px-3 font-mono text-sm text-muted-foreground">
            Thun-u-der Express
          </div>
        </div>
        <div>
          <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground">Weight</label>
          <Input v-model="weight" class="mt-1.5 font-mono text-sm" placeholder="e.g. 12.4 kg" />
          <p v-if="errors.weight" class="mt-1 font-mono text-xs text-destructive">{{ errors.weight }}</p>
        </div>
        <div>
          <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground">Items</label>
          <Input v-model.number="items" type="number" min="1" class="mt-1.5 font-mono text-sm" />
          <p v-if="errors.items" class="mt-1 font-mono text-xs text-destructive">{{ errors.items }}</p>
        </div>
        <div>
          <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground">Estimated Delivery</label>
          <Input v-model="estimatedDelivery" class="mt-1.5 font-mono text-sm" placeholder="e.g. May 25, 2026" />
          <p v-if="errors.estimatedDelivery" class="mt-1 font-mono text-xs text-destructive">{{ errors.estimatedDelivery }}</p>
        </div>
        <div v-if="isEditing">
          <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground">Status</label>
          <select
            v-model="status"
            class="mt-1.5 flex h-10 w-full rounded-lg border border-border bg-background px-3 font-mono text-sm"
          >
            <option v-for="(label, key) in statusLabels" :key="key" :value="key">{{ label }}</option>
          </select>
        </div>
      </div>
    </fieldset>

    <div class="flex justify-end gap-3 pt-4 border-t border-border">
      <Button variant="outline" type="button" @click="emit('cancel')">Cancel</Button>
      <Button type="submit" :disabled="pending">
        {{ pending ? 'Saving\u2026' : isEditing ? 'Save Changes' : 'Create Order' }}
      </Button>
    </div>
  </form>
</template>
