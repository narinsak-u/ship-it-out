<script setup lang="ts">
import { ref, watch, computed } from "vue";
import {
  getSubDistrictNames,
  getDistrictNames,
  getProvinceName,
} from "thai-data";
import Input from "@/components/ui/Input.vue";

export interface ThaiAddress {
  name: string;
  zipcode: string;
  subDistrict: string;
  district: string;
  province: string;
}

const props = defineProps<{
  label: string;
  modelValue: ThaiAddress;
  errors?: Record<string, string>;
}>();

const emit = defineEmits<{
  "update:modelValue": [value: ThaiAddress];
}>();

const availableSubDistricts = ref<string[]>([]);

function patch(field: keyof ThaiAddress, value: string) {
  emit("update:modelValue", { ...props.modelValue, [field]: value });
}

function selectSubDistrict(sub: string) {
  const zip = props.modelValue.zipcode;
  patch("subDistrict", sub);
  const districts = getDistrictNames(zip);
  const province = getProvinceName(zip);
  if (districts.length > 0) patch("district", districts[0]);
  if (province) patch("province", province);
}

const districtDisplay = computed(() => {
  if (props.modelValue.district) return props.modelValue.district;
  const names = getDistrictNames(props.modelValue.zipcode);
  return names.length > 0 ? names[0] : "";
});

const provinceDisplay = computed(() => {
  if (props.modelValue.province) return props.modelValue.province;
  return getProvinceName(props.modelValue.zipcode) ?? "";
});

watch(
  () => props.modelValue.zipcode,
  (zip, oldZip) => {
    if (zip.length === 5) {
      availableSubDistricts.value = getSubDistrictNames(zip);
      if (oldZip && zip !== oldZip) {
        patch("subDistrict", "");
        patch("district", "");
        patch("province", "");
      }
    } else {
      availableSubDistricts.value = [];
      if (oldZip && oldZip.length === 5) {
        patch("subDistrict", "");
        patch("district", "");
        patch("province", "");
      }
    }
  },
  { immediate: true },
);

watch(availableSubDistricts, (list) => {
  if (list.length === 1 && !props.modelValue.subDistrict) {
    selectSubDistrict(list[0]);
  }
});

function onSubDistrictSelected(ev: Event) {
  const sub = (ev.target as HTMLSelectElement).value;
  if (!sub) return;
  selectSubDistrict(sub);
}
</script>

<template>
  <fieldset class="rounded-xl border border-border p-5">
    <legend class="font-mono text-xs uppercase tracking-widest text-primary px-2">
      {{ label }}
    </legend>
    <div class="grid gap-5 md:grid-cols-2">
      <div>
        <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground"
          >Name</label
        >
        <Input
          :value="modelValue.name"
          class="mt-1.5 font-mono text-sm"
          placeholder="e.g. ประวิทย์ ใจดี"
          @input="patch('name', ($event.target as HTMLInputElement).value)"
        />
        <p v-if="errors?.name" class="mt-1 font-mono text-xs text-destructive">
          {{ errors.name }}
        </p>
      </div>
      <div>
        <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground"
          >Zipcode</label
        >
        <Input
          :value="modelValue.zipcode"
          class="mt-1.5 font-mono text-sm"
          maxlength="5"
          inputmode="numeric"
          placeholder="e.g. 10200"
          @input="patch('zipcode', ($event.target as HTMLInputElement).value)"
        />
        <p v-if="errors?.zipcode" class="mt-1 font-mono text-xs text-destructive">
          {{ errors.zipcode }}
        </p>
      </div>
      <div>
        <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground"
          >Sub-district</label
        >
        <select
          v-if="availableSubDistricts.length"
          :value="modelValue.subDistrict"
          class="mt-1.5 flex h-10 w-full rounded-lg border border-border bg-background px-3 font-mono text-sm"
          @change="onSubDistrictSelected"
        >
          <option disabled value="">Select sub-district...</option>
          <option
            v-for="sd in availableSubDistricts"
            :key="sd"
            :value="sd"
          >
            {{ sd }}
          </option>
        </select>
        <Input
          v-else
          disabled
          class="mt-1.5 font-mono text-sm"
          placeholder="Enter zipcode first"
        />
        <p v-if="errors?.subDistrict" class="mt-1 font-mono text-xs text-destructive">
          {{ errors.subDistrict }}
        </p>
      </div>
      <div>
        <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground"
          >District</label
        >
        <Input :value="districtDisplay" disabled class="mt-1.5 font-mono text-sm" />
        <p v-if="errors?.district" class="mt-1 font-mono text-xs text-destructive">
          {{ errors.district }}
        </p>
      </div>
      <div>
        <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground"
          >Province</label
        >
        <Input :value="provinceDisplay" disabled class="mt-1.5 font-mono text-sm" />
        <p v-if="errors?.province" class="mt-1 font-mono text-xs text-destructive">
          {{ errors.province }}
        </p>
      </div>
    </div>
  </fieldset>
</template>
