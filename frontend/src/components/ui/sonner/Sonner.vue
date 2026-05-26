<script lang="ts" setup>
import type { ToasterProps } from "vue-sonner";
import {
  CircleCheckIcon,
  CircleXIcon,
  InfoIcon,
  Loader2Icon,
  TriangleAlertIcon,
  XIcon,
} from "@lucide/vue";
import { Toaster as Sonner } from "vue-sonner";
import { cn } from "@/lib/utils";

const props = defineProps<ToasterProps>();
</script>

<template>
  <Sonner
    :class="cn('toaster group', props.class)"
    :style="{
      '--normal-bg': 'var(--toast-bg)',
      '--normal-text': 'var(--card-foreground)',
      '--normal-border': 'var(--border)',
      '--success-bg': 'var(--toast-bg)',
      '--success-text': 'var(--card-foreground)',
      '--error-bg': 'var(--toast-bg)',
      '--error-text': 'var(--card-foreground)',
      '--warning-bg': 'var(--toast-bg)',
      '--warning-text': 'var(--card-foreground)',
      '--info-bg': 'var(--toast-bg)',
      '--info-text': 'var(--card-foreground)',
      '--border-radius': 'var(--radius)',
    }"
    v-bind="props"
  >
    <template #success-icon>
      <CircleCheckIcon class="size-4" />
    </template>
    <template #info-icon>
      <TriangleAlertIcon class="size-4" />
    </template>
    <template #warning-icon>
      <InfoIcon class="size-4" />
    </template>
    <template #error-icon>
      <CircleXIcon class="size-4" />
    </template>
    <template #loading-icon>
      <div>
        <Loader2Icon class="size-4 animate-spin" />
      </div>
    </template>
    <template #close-icon>
      <XIcon class="size-4" />
    </template>
  </Sonner>
</template>
