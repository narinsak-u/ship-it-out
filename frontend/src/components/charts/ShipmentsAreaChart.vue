<script setup lang="ts">
import type { HTMLAttributes } from "vue";
import { VisArea, VisXYContainer, VisAxis } from "@unovis/vue";
import type { ChartConfig } from "@/components/ui/chart";
import {
  ChartContainer,
  ChartTooltip,
  ChartCrosshair,
  ChartTooltipContent,
  componentToString,
} from "@/components/ui/chart";

export interface CumulativeEntry {
  month: string;
  count: number;
}

const props = defineProps<{
  data: CumulativeEntry[];
  class?: HTMLAttributes["class"];
}>();

type Data = CumulativeEntry;

const chartConfig = {
  count: {
    label: "Total Shipments",
    color: "var(--color-primary)",
  },
} satisfies ChartConfig;
</script>

<template>
  <div v-if="data.length === 0" class="flex h-64 items-center justify-center text-sm text-muted-foreground">
    No data
  </div>
  <ChartContainer
    v-else
    :config="chartConfig"
    :class="props.class"
  >
    <VisXYContainer :data="data">
      <VisArea
        :x="(d: Data) => new Date(d.month + '-01').getTime()"
        :y="(d: Data) => d.count"
        :color="['var(--color-count)']"
      />
      <VisAxis
        type="x"
        :x="(d: Data) => new Date(d.month + '-01').getTime()"
        :tick-line="false"
        :domain-line="false"
        :grid-line="false"
        :tick-format="(v: number) => {
          const date = new Date(v);
          return date.toLocaleDateString('en-US', { month: 'short', year: '2-digit' });
        }"
        :tick-values="data.map(d => new Date(d.month + '-01').getTime())"
      />
      <VisAxis
        type="y"
        :tick-line="false"
        :domain-line="false"
        :grid-line="true"
      />
      <ChartTooltip />
      <ChartCrosshair
        :template="componentToString(chartConfig, ChartTooltipContent)"
        :color="['var(--color-count)']"
      />
    </VisXYContainer>
  </ChartContainer>
</template>
