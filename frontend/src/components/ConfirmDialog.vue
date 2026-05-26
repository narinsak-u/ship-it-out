<script setup lang="ts">
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import Button from "@/components/ui/Button.vue";

defineProps<{
  open: boolean;
  title: string;
  description: string;
  pending?: boolean;
}>();

const emit = defineEmits<{
  confirm: [];
  cancel: [];
}>();
</script>

<template>
  <Dialog :open="open" @update:open="(v) => !v && emit('cancel')">
    <DialogContent
      class="w-full max-w-md rounded-xl border border-border bg-card p-6 shadow-elegant sm:max-w-md"
    >
      <DialogHeader>
        <DialogTitle class="font-mono text-lg font-semibold">{{ title }}</DialogTitle>
      </DialogHeader>
      <p class="font-mono text-sm text-muted-foreground">{{ description }}</p>
      <div class="mt-4 flex justify-end gap-3">
        <Button variant="outline" @click="emit('cancel')">Cancel</Button>
        <Button variant="destructive" :disabled="pending" @click="emit('confirm')">
          {{ pending ? "Deleting…" : "Delete" }}
        </Button>
      </div>
    </DialogContent>
  </Dialog>
</template>
