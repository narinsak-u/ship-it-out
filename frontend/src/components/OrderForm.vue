<script setup lang="ts">
import { ref } from 'vue'
import { carriers as carrierList } from '@/lib/carriers'
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

const customer = ref(props.initial?.customer ?? '')
const origin = ref(props.initial?.origin ?? '')
const destination = ref(props.initial?.destination ?? '')
const carrier = ref(props.initial?.carrier ?? carrierList[0]?.name ?? '')
const weight = ref(props.initial?.weight ?? '')
const items = ref(props.initial?.items ?? 1)
const estimatedDelivery = ref(props.initial?.estimatedDelivery ?? '')
const status = ref<ShipmentStatus>(props.initial?.status ?? 'pending')

const errors = ref<Record<string, string>>({})

function validate(): boolean {
  const e: Record<string, string> = {}
  if (!customer.value.trim()) e.customer = 'Required'
  if (!origin.value.trim()) e.origin = 'Required'
  if (!destination.value.trim()) e.destination = 'Required'
  if (!carrier.value.trim()) e.carrier = 'Required'
  if (!weight.value.trim()) e.weight = 'Required'
  if (!items.value || items.value < 1) e.items = 'Must be at least 1'
  if (!estimatedDelivery.value.trim()) e.estimatedDelivery = 'Required'
  errors.value = e
  return Object.keys(e).length === 0
}

function handleSubmit() {
  if (!validate()) return
  emit('submit', {
    customer: customer.value,
    origin: origin.value,
    destination: destination.value,
    carrier: carrier.value,
    weight: weight.value,
    items: items.value,
    estimatedDelivery: estimatedDelivery.value,
    ...(props.isEditing ? { status: status.value } : {}),
  })
}
</script>

<template>
  <form @submit.prevent="handleSubmit" class="space-y-5">
    <div class="grid gap-5 md:grid-cols-2">
      <!-- Customer -->
      <div>
        <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground">Customer</label>
        <Input v-model="customer" class="mt-1.5 font-mono text-sm" placeholder="e.g. Aria Nakamura" />
        <p v-if="errors.customer" class="mt-1 font-mono text-xs text-destructive">{{ errors.customer }}</p>
      </div>

      <!-- Carrier -->
      <div>
        <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground">Carrier</label>
        <select
          v-model="carrier"
          class="mt-1.5 flex h-10 w-full rounded-lg border border-border bg-background px-3 font-mono text-sm"
        >
          <option v-for="c in carrierList" :key="c.id" :value="c.name">{{ c.name }}</option>
        </select>
        <p v-if="errors.carrier" class="mt-1 font-mono text-xs text-destructive">{{ errors.carrier }}</p>
      </div>

      <!-- Origin -->
      <div>
        <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground">Origin</label>
        <Input v-model="origin" class="mt-1.5 font-mono text-sm" placeholder="e.g. Rotterdam, NL" />
        <p v-if="errors.origin" class="mt-1 font-mono text-xs text-destructive">{{ errors.origin }}</p>
      </div>

      <!-- Destination -->
      <div>
        <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground">Destination</label>
        <Input v-model="destination" class="mt-1.5 font-mono text-sm" placeholder="e.g. Brooklyn, NY" />
        <p v-if="errors.destination" class="mt-1 font-mono text-xs text-destructive">{{ errors.destination }}</p>
      </div>

      <!-- Weight -->
      <div>
        <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground">Weight</label>
        <Input v-model="weight" class="mt-1.5 font-mono text-sm" placeholder="e.g. 12.4 kg" />
        <p v-if="errors.weight" class="mt-1 font-mono text-xs text-destructive">{{ errors.weight }}</p>
      </div>

      <!-- Items -->
      <div>
        <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground">Items</label>
        <Input v-model.number="items" type="number" min="1" class="mt-1.5 font-mono text-sm" />
        <p v-if="errors.items" class="mt-1 font-mono text-xs text-destructive">{{ errors.items }}</p>
      </div>

      <!-- Estimated Delivery -->
      <div>
        <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground">Estimated Delivery</label>
        <Input v-model="estimatedDelivery" class="mt-1.5 font-mono text-sm" placeholder="e.g. May 25, 2026" />
        <p v-if="errors.estimatedDelivery" class="mt-1 font-mono text-xs text-destructive">{{ errors.estimatedDelivery }}</p>
      </div>

      <!-- Status (edit only) -->
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

    <div class="flex justify-end gap-3 pt-4 border-t border-border">
      <Button variant="outline" type="button" @click="emit('cancel')">Cancel</Button>
      <Button type="submit" :disabled="pending">
        {{ pending ? 'Saving\u2026' : isEditing ? 'Save Changes' : 'Create Order' }}
      </Button>
    </div>
  </form>
</template>
