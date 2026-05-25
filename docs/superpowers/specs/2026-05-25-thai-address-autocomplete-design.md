# Thai Address Autocomplete for Order Form

**Date:** 2026-05-25
**Status:** Design

## Problem

The OrderForm currently has free-text `<Input>` fields for address components (sub-district, district, province) with non-Thai placeholder text. This is inaccurate for the real data domain (Thailand) and error-prone — users can type any value, leading to inconsistent addresses in the system.

## Solution

Replace the sender and receiver address sections in OrderForm with a reusable `ThaiAddressGroup.vue` component that uses the `thai-data` library (already installed at v3.0.2) to provide zipcode-driven address autocomplete.

## Flow

1. User types a 5-digit Thai zipcode in the zipcode field
2. On 5 characters entered, the sub-district `<Input>` is replaced with a native `<select>` dropdown populated by `thai-data`'s `getSubDistrictNames()`
3. User selects a sub-district from the dropdown
4. District and province fields auto-fill via `getDistrictNames()` / `getProvinceName()` and become read-only
5. If user modifies the zipcode after selection, the sub-district dropdown resets and district/province clear

## Component: `ThaiAddressGroup.vue`

**Note on API:** The installed `thai-data` package (v3.0.2) exports v2 API names despite the README showing v3 names. The actual importable functions are:
- `getSubDistrictNames(zipcode)` → `string[]`
- `getDistrictNames(zipcode)` → `string[]`
- `getProvinceName(zipcode)` → `string | null`

These are used directly — no network request, the data is bundled JSON.

**Path:** `frontend/src/components/ThaiAddressGroup.vue`

### Props

```ts
interface ThaiAddress {
  name: string;
  zipcode: string;
  subDistrict: string;
  district: string;
  province: string;
}

defineProps<{
  label: string;            // "Sender Info" or "Receiver Info"
  modelValue: ThaiAddress;
  errors?: Record<string, string>;
}>();
```

### Emits

```ts
defineEmits<{
  'update:modelValue': [value: ThaiAddress];
}>();
```

### Internal state

```ts
const availableSubDistricts = ref<string[]>([]); // populated on zipcode lookup
```

### Watchers

- `watch(() => modelValue.zipcode, ...)` — on 5 digits, calls `getSubDistrictNames()`. Clears district/province if zipcode changed after previous lookup. Runs `{ immediate: true }` for edit mode.
- `watch(availableSubDistricts, ...)` — if exactly 1 sub-district exists (rare), auto-select it; otherwise user picks from dropdown.

### Behavior details

- **Zipcode < 5 digits:** Sub-district is a disabled `<Input>` (can't free-type)
- **Zipcode = 5 digits + results found:** Sub-district becomes `<select>` populated with available sub-districts; district + province auto-fill
- **Zipcode = 5 digits + no results:** Sub-district stays as disabled `<Input>`, district/province stay empty — user sees no dropdown (postal code not in thai-data)
- **Edit mode:** On mount, if zipcode is filled, trigger lookup immediately. If saved sub-district matches a returned option, pre-select it. District and province from saved data take precedence (in case thai-data doesn't have an exact match for older addresses).
- **Read-only fields:** District and province are `<Input disabled>` once auto-filled
- **Name field:** Remains free-text `<Input>`, unchanged

### Template structure

```
<fieldset>
  <legend>{{ label }}</legend>
  <div class="grid gap-5 md:grid-cols-2">
    <div>  <!-- Name -->
      <label>Name</label>
      <Input v-model="localName" />
    </div>
    <div>  <!-- Zipcode -->
      <label>Zipcode</label>
      <Input v-model="localZipcode" />
    </div>
    <div>  <!-- Sub-district -->
      <label>Sub-district</label>
      <select v-if="availableSubDistricts.length" v-model="localSubDistrict">
        <option disabled value="">Select sub-district...</option>
        <option v-for="sd in availableSubDistricts" :key="sd" :value="sd">{{ sd }}</option>
      </select>
      <Input v-else disabled placeholder="Enter zipcode first" />
    </div>
    <div>  <!-- District -->
      <label>District</label>
      <Input :value="districtDisplay" disabled />
    </div>
    <div>  <!-- Province -->
      <label>Province</label>
      <Input :value="provinceDisplay" disabled />
    </div>
  </div>
</fieldset>
```

### Error display

Error messages for name, zipcode, and sub-district (the required fields) are displayed inline per field using the same pattern as the current OrderForm.

## Changes to OrderForm.vue

- Remove the sender and receiver `<fieldset>` blocks
- Import and use two `<ThaiAddressGroup>` instances instead
- Keep parcel info `<fieldset>` as-is
- The form data model stays exactly the same (still emits `OrderFormData`)

## Changes to OrderFormView.vue

No changes needed — OrderFormView only passes data through and handles submit/cancel.

## Validation

- `zipcode`: Required, must be 5 digits
- `subDistrict`: Required, must be selected from dropdown (not free-text)
- `district` and `province`: Required, auto-filled by lookup
- All other validations (name, weight, items, estimatedDelivery) unchanged

## Edge Cases

| Scenario | Behavior |
|----------|----------|
| Typing fewer than 5 digits | Sub-district field stays disabled, no lookup |
| No results for zipcode | Sub-district field stays disabled, no dropdown shown |
| User changes zipcode after selection | Sub-district dropdown resets, district/province clear |
| Exactly 1 sub-district result | Auto-selected, district/province auto-fill immediately |
| Edit mode + saved data | On mount, lookup runs; saved sub-district pre-selected if match found |
| Special chars in Thai names | Native `<select>` handles Unicode correctly; no extra work needed |

## What's NOT changing

- No new shadcn-vue components added (uses native `<select>` styled like existing Status dropdown)
- No API changes — thai-data is client-side only, no backend changes
- No dependency changes — thai-data is already in `package.json`
- No changes to `OrderFormView.vue`
- No changes to hooks, API client, mappers, or backend
