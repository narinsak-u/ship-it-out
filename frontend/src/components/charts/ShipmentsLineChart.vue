<script setup lang="ts">
import type { HTMLAttributes } from "vue";
import { VisLine, VisXYContainer, VisAxis } from "@unovis/vue";
import type { ChartConfig } from "@/components/ui/chart";
import {
  ChartContainer,
  ChartTooltip,
  ChartCrosshair,
  ChartTooltipContent,
  componentToString,
} from "@/components/ui/chart";
import type { MonthlyCount } from "@/lib/api/analytics";

const props = defineProps<{
  data: MonthlyCount[];
  class?: HTMLAttributes["class"];
}>();

type Data = MonthlyCount;

const chartConfig = {
  count: {
    label: "Shipments",
    color: "hsl(var(--primary))",
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
      <VisLine
        :x="(d: Data) => new Date(d.month + '-01').getTime()"
        :y="(d: Data) => d.count"
        :color="[chartConfig.count.color]"
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
        :color="[chartConfig.count.color]"
      />
    </VisXYContainer>
  </ChartContainer>
</template>
