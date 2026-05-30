<script setup lang="ts">
defineProps<{
  counts: {
    total: number;
    active: number;
    maintenance: number;
    closed: number;
    full: number;
  };
}>();
</script>

<template>
  <div class="grid grid-cols-2 gap-4 md:grid-cols-5">
    <div class="rounded-lg border border-border bg-secondary/50 p-4">
      <div class="font-mono text-[11px] uppercase tracking-widest text-muted-foreground">
        Total Hubs
      </div>
      <div class="mt-1 font-mono text-3xl font-semibold">{{ counts.total }}</div>
    </div>
    <div class="rounded-lg border border-border bg-secondary/50 p-4">
      <div class="font-mono text-[11px] uppercase tracking-widest text-muted-foreground">
        Active
      </div>
      <div class="mt-1 font-mono text-3xl font-semibold text-success">{{ counts.active }}</div>
    </div>
    <div class="rounded-lg border border-border bg-secondary/50 p-4">
      <div class="font-mono text-[11px] uppercase tracking-widest text-muted-foreground">
        Maintenance
      </div>
      <div class="mt-1 font-mono text-3xl font-semibold text-warning">{{ counts.maintenance }}</div>
    </div>
    <div class="rounded-lg border border-border bg-secondary/50 p-4">
      <div class="font-mono text-[11px] uppercase tracking-widest text-muted-foreground">Full</div>
      <div class="mt-1 font-mono text-3xl font-semibold text-destructive">{{ counts.full }}</div>
    </div>
    <div class="rounded-lg border border-border bg-secondary/50 p-4">
      <div class="font-mono text-[11px] uppercase tracking-widest text-muted-foreground">
        Closed
      </div>
      <div class="mt-1 font-mono text-3xl font-semibold text-destructive">{{ counts.closed }}</div>
    </div>
  </div>
</template>
