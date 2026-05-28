<script setup lang="ts">
import { computed } from "vue";
import type { HTMLAttributes } from "vue";
import { VisDonut, VisDonutSelectors, VisSingleContainer } from "@unovis/vue";
import type { ChartConfig } from "@/components/ui/chart";
import {
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
  componentToString,
} from "@/components/ui/chart";

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

const tooltipContentFn = componentToString(chartConfig.value, ChartTooltipContent, {
  labelKey: "count",
  nameKey: "status",
});
</script>

<template>
  <div v-if="data.length === 0" class="flex h-48 items-center justify-center text-sm text-muted-foreground">
    No data
  </div>
  <ChartContainer
    v-else
    :config="chartConfig"
    :class="props.class"
  >
    <VisSingleContainer :data="data">
      <VisDonut
        :value="(d: Data) => d.count"
        :color="(d: Data) => `var(--color-${d.status})`"
      />
      <ChartTooltip
        v-if="tooltipContentFn"
        :triggers="{
          [VisDonutSelectors.segment]: tooltipContentFn,
        }"
      />
    </VisSingleContainer>
  </ChartContainer>
</template>
