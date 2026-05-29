<script setup lang="ts">
import { computed } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useQuery } from "@tanstack/vue-query";
import { toast } from "vue-sonner";
import { useCreateOrder, useUpdateOrder } from "@/hooks/useOrders";
import { fetchOrder } from "@/lib/api/orders";
import type { OrderFormData } from "@/lib/api/orders";
import { orderKeys } from "@/lib/api/queryKeys";
import OrderForm from "@/components/OrderForm.vue";
import Skeleton from "@/components/ui/Skeleton.vue";

const route = useRoute();
const router = useRouter();
const createOrder = useCreateOrder();
const updateOrder = useUpdateOrder();
const hasSaveError = computed(() => createOrder.isError.value || updateOrder.isError.value);

const isEditing = computed(() => !!route.params.orderId);
const orderId = computed(() => route.params.orderId as string | undefined);

const { data: order } = useQuery({
  queryKey: orderKeys.detail(orderId.value!),
  queryFn: () => fetchOrder(orderId.value!),
  enabled: isEditing,
});

const isPending = computed(() => createOrder.isPending.value || updateOrder.isPending.value);

async function handleSubmit(data: OrderFormData) {
  if (isEditing.value && orderId.value) {
    await updateOrder.mutateAsync({ id: orderId.value, data });
    toast.success("Order updated");
    router.push({ name: "order-detail", params: { orderId: orderId.value } });
  } else {
    const created = await createOrder.mutateAsync(data);
    toast.success("Order created");
    router.push({ name: "orders" });
  }
}

function handleCancel() {
  if (isEditing.value && orderId.value) {
    router.push({ name: "order-detail", params: { orderId: orderId.value } });
  } else {
    router.push({ name: "orders" });
  }
}
</script>

<template>
  <div>
    <section class="border-b border-border bg-gradient-hero">
      <div class="mx-auto max-w-7xl px-6 py-14">
        <span class="font-mono text-xs uppercase tracking-widest text-primary">/ orders</span>
        <h1 class="mt-3 text-4xl font-semibold tracking-tight md:text-5xl">
          {{ isEditing ? "Edit Order" : "Create Order" }}
        </h1>
        <p class="mt-3 max-w-2xl text-muted-foreground">
          {{
            isEditing
              ? "Update shipment details and status."
              : "Register a new shipment in the system."
          }}
        </p>
      </div>
    </section>

    <section class="mx-auto max-w-3xl px-6 py-10">
      <div v-if="isEditing && !order" class="py-16 text-center">
        <p class="font-mono text-lg">Loading order...</p>
        <Skeleton class="mt-4 h-96 rounded-xl" />
      </div>

      <div v-else-if="isEditing && !order">
        <Skeleton class="h-96 rounded-xl" />
      </div>

      <div v-else class="rounded-xl border border-border bg-card p-6 shadow-elegant">
        <OrderForm
          :initial="order ?? undefined"
          :is-editing="isEditing"
          :pending="isPending"
          @submit="handleSubmit"
          @cancel="handleCancel"
        />
      </div>

      <div
        v-if="hasSaveError"
        class="mt-4 rounded-lg bg-destructive/15 px-4 py-3 font-mono text-sm text-destructive"
      >
        Failed to save order. Please try again.
      </div>
    </section>
  </div>
</template>
