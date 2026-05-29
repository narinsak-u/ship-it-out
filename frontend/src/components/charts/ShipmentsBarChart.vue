<script setup lang="ts">
import { computed, type HTMLAttributes } from "vue";
import { VisGroupedBar, VisXYContainer, VisAxis } from "@unovis/vue";
import type { ChartConfig } from "@/components/ui/chart";
import {
  ChartContainer,
  ChartTooltip,
  ChartCrosshair,
  ChartTooltipContent,
  componentToString,
} from "@/components/ui/chart";
import type { DayOfWeekCount } from "@/lib/api/analytics";

const DAYS = [
  "Monday",
  "Tuesday",
  "Wednesday",
  "Thursday",
  "Friday",
  "Saturday",
  "Sunday",
] as const;

const props = defineProps<{
  data: DayOfWeekCount[];
  class?: HTMLAttributes["class"];
}>();

interface DataPoint {
  day: string;
  count: number;
  dayIndex: number;
}

const mappedData = computed<DataPoint[]>(() => {
  return props.data.map((d) => ({
    dayIndex: DAYS.indexOf(d.day as (typeof DAYS)[number]),
    ...d,
  }));
});

const yMax = computed(() => Math.max(...mappedData.value.map((d) => d.count), 0));

const yTicks = computed(() => {
  const ticks: number[] = [];
  const max = yMax.value % 2 === 0 ? yMax.value : yMax.value + 1;
  for (let i = 0; i <= max; i += 2) {
    ticks.push(i);
  }
  return ticks;
});

const chartConfig = {
  count: {
    label: "Shipments",
    color: "var(--color-primary)",
  },
} satisfies ChartConfig;
</script>

<template>
  <div
    v-if="mappedData.length === 0"
    class="flex h-48 items-center justify-center text-sm text-muted-foreground"
  >
    No data
  </div>
  <ChartContainer v-else :config="chartConfig" :class="props.class">
    <VisXYContainer :data="mappedData">
      <VisGroupedBar
        :x="(d: DataPoint) => d.dayIndex"
        :y="(d: DataPoint) => d.count"
        :color="['var(--color-count)']"
        :rounded-corners="4"
        bar-padding="0.1"
        group-padding="0"
      />
      <VisAxis
        type="x"
        :x="(d: DataPoint) => d.dayIndex"
        :tick-line="false"
        :domain-line="false"
        :grid-line="false"
        :tick-format="(v: number) => DAYS[v].slice(0, 3)"
        :tick-values="DAYS.map((_, i) => i)"
      />
      <VisAxis
        type="y"
        :tick-line="false"
        :domain-line="false"
        :grid-line="true"
        :tick-values="yTicks"
      />
      <ChartTooltip />
      <ChartCrosshair
        :template="componentToString(chartConfig, ChartTooltipContent)"
        :color="['var(--color-count)']"
      />
    </VisXYContainer>
  </ChartContainer>
</template>
