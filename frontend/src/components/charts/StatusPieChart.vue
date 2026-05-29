<script setup lang="ts">
import { computed } from "vue";
import type { HTMLAttributes } from "vue";
import { cn } from "@/lib/utils";
import { VisDonut, VisSingleContainer } from "@unovis/vue";
import { Donut } from "@unovis/ts";
import type { ChartConfig } from "@/components/ui/chart";
import { ChartContainer, ChartTooltip } from "@/components/ui/chart";

export interface StatusPieEntry {
  status: string;
  label: string;
  count: number;
  pct: number;
  fill: string;
}

const props = defineProps<{
  data: StatusPieEntry[];
  class?: HTMLAttributes["class"];
}>();

type Data = StatusPieEntry;

const chartConfig = computed<ChartConfig>(() => {
  const config: ChartConfig = {};
  for (const entry of props.data) {
    config[entry.status] = { label: entry.label, color: entry.fill };
  }
  return config;
});

const tooltipFn = (d: Record<string, any>) => {
  const item = d?.data ?? d;
  if (!item?.label) return "";
  return `<div class="flex items-center gap-2 rounded-lg border border-border/50 bg-background px-2.5 py-1.5 text-xs shadow-xl">
    <div class="h-2.5 w-2.5 shrink-0 rounded-[2px]" style="background:${item.fill}"></div>
    <span class="text-muted-foreground">${item.label}</span>
    <span class="font-mono font-medium tabular-nums text-foreground">${item.count}</span>
  </div>`;
};

const labelFormatter = (d: Data) => `${d.count} ${d.label}`;
</script>

<template>
  <div
    v-if="data.length === 0"
    class="flex h-48 items-center justify-center text-sm text-muted-foreground"
  >
    No data
  </div>
  <ChartContainer v-else :config="chartConfig" :class="cn('aspect-square', props.class)">
    <VisSingleContainer :data="data">
      <VisDonut
        :value="(d: Data) => d.count"
        :color="(d: Data) => `var(--color-${d.status})`"
        :arc-width="0"
        :show-labels="true"
        :label="labelFormatter"
        :label-color="'var(--background)'"
        stroke="var(--background)"
        :stroke-width="1"
      />
      <ChartTooltip
        :triggers="{
          [Donut.selectors.segment]: tooltipFn,
        }"
      />
    </VisSingleContainer>
  </ChartContainer>
</template>
