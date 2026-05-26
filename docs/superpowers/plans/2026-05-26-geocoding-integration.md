# OpenCage Geocoding Integration Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Geocode addresses to real lat/lng coordinates before creating orders and hubs, replacing the current `coords: { lat: 0, lng: 0 }` placeholders.

**Architecture:** A `geocodeAddress()` utility in `frontend/src/lib/geocode.ts` calls the OpenCage Geocoding API using the `opencage-api-client` library (already installed). `OrderForm.vue` and `HubFormModal.vue` call this utility on submit, block on failure, and inject real coords into the payload before sending to the backend.

**Tech Stack:** Vue 3 (Composition API), TypeScript, `opencage-api-client` v2.1.2, Vite (`import.meta.env`), OpenCage Geocoding API

---

### Task 1: Set up environment variable

**Files:**
- Create: `frontend/.env`
- Modify: `frontend/.env.example`

- [ ] **Step 1: Create `frontend/.env` with the API key placeholder**

```ini
VITE_OPENCAGE_API_KEY=your-opencage-api-key-here
```

- [ ] **Step 2: Add the key to `frontend/.env.example`**

Read the existing `.env.example` first to check its format, then append the new variable:

```
VITE_OPENCAGE_API_KEY=your-opencage-api-key-here
```

- [ ] **Step 3: Commit**

```bash
git add frontend/.env frontend/.env.example
git commit -m "feat: add VITE_OPENCAGE_API_KEY env variable for geocoding"
```

---

### Task 2: Create geocodeAddress utility

**Files:**
- Create: `frontend/src/lib/geocode.ts`

- [ ] **Step 1: Write the geocodeAddress function**

```typescript
import { geocode } from "opencage-api-client";

export async function geocodeAddress(
  subDistrict: string,
  district: string,
  province: string,
): Promise<{ lat: number; lng: number }> {
  const q = `${subDistrict}, ${district}, ${province}`.trim();
  const key = import.meta.env.VITE_OPENCAGE_API_KEY as string;

  if (!key) {
    throw new Error("OpenCage API key is not configured.");
  }

  let data;
  try {
    data = await geocode({ q, key });
  } catch {
    throw new Error(
      "Location lookup failed. Please try again later.",
    );
  }

  if (!data.results || data.results.length === 0) {
    throw new Error("Could not resolve this address. Check the fields and try again.");
  }

  const { lat, lng } = data.results[0].geometry;
  return { lat, lng };
}
```

- [ ] **Step 2: Verify the build compiles**

```bash
cd frontend && npm run build
```

Expected: Build succeeds (vue-tsc + vite). If `opencage-api-client` has type issues, verify it exports `geocode` and has compatible types.

- [ ] **Step 3: Commit**

```bash
git add frontend/src/lib/geocode.ts
git commit -m "feat: add geocodeAddress utility for OpenCage geocoding"
```

---

### Task 3: Modify OrderForm.vue to geocode on submit

**Files:**
- Modify: `frontend/src/components/OrderForm.vue`

- [ ] **Step 1: Add the geocode import and error state**

Add to the `<script setup>` block (after existing imports):

```typescript
import { geocodeAddress } from "@/lib/geocode";
```

Add the geocode errors ref after the `errors` ref:

```typescript
const geocodeErrors = ref<Record<string, string>>({});
```

Add the geocoding-in-progress ref after `geocodeErrors`:

```typescript
const geocoding = ref(false);
```

- [ ] **Step 2: Update handleSubmit to geocode before emitting**

Replace the existing `handleSubmit` function:

```typescript
async function handleSubmit() {
  if (!validate()) return;
  geocodeErrors.value = {};
  geocoding.value = true;

  try {
    const [senderCoords, receiverCoords] = await Promise.all([
      geocodeAddress(
        sender.value.subDistrict,
        sender.value.district,
        sender.value.province,
      ),
      geocodeAddress(
        receiver.value.subDistrict,
        receiver.value.district,
        receiver.value.province,
      ),
    ]);

    emit("submit", {
      customer: {
        name: sender.value.name,
        zipcode: sender.value.zipcode,
        subDistrict: sender.value.subDistrict,
        district: sender.value.district,
        province: sender.value.province,
        coords: senderCoords,
      },
      receiver: {
        name: receiver.value.name,
        zipcode: receiver.value.zipcode,
        subDistrict: receiver.value.subDistrict,
        district: receiver.value.district,
        province: receiver.value.province,
        coords: receiverCoords,
      },
      carrier: carrier.value,
      weight: weight.value,
      items: items.value,
      estimatedDelivery: estimatedDeliveryRaw.value || "",
    });
  } catch (e) {
    const msg = e instanceof Error ? e.message : "Could not resolve address.";
    geocodeErrors.value = { sender: msg, receiver: msg };
  } finally {
    geocoding.value = false;
  }
}
```

Note: `Promise.all` (not `allSettled`) is used intentionally — if either geocode fails, we want to reject immediately and show an error, since we can't submit partial coords.

- [ ] **Step 3: Add inline error display in the template**

After the `<ThaiAddressGroup>` for sender, add the geocode error:

```html
<ThaiAddressGroup label="Sender Info" v-model="sender" :errors="senderErrorMap" />
<p
  v-if="geocodeErrors.sender"
  class="mt-2 rounded-md bg-destructive/15 px-3 py-2 font-mono text-xs text-destructive"
>
  {{ geocodeErrors.sender }}
</p>
```

After the `<ThaiAddressGroup>` for receiver, add:

```html
<ThaiAddressGroup label="Receiver Info" v-model="receiver" :errors="receiverErrorMap" />
<p
  v-if="geocodeErrors.receiver"
  class="mt-2 rounded-md bg-destructive/15 px-3 py-2 font-mono text-xs text-destructive"
>
  {{ geocodeErrors.receiver }}
</p>
```

- [ ] **Step 4: Update submit button text during geocoding**

Replace the submit button text:

```html
<Button type="submit" :disabled="pending || !canSubmit || geocoding">
  {{
    geocoding
      ? "Resolving addresses\u2026"
      : pending
        ? "Saving\u2026"
        : isEditing
          ? "Save Changes"
          : "Create Order"
  }}
</Button>
```

Note: The `geocoding` check is added to `:disabled` so the button can't be double-clicked while geocoding.

- [ ] **Step 5: Verify the build compiles**

```bash
cd frontend && npm run build
```

Expected: Build succeeds.

- [ ] **Step 6: Commit**

```bash
git add frontend/src/components/OrderForm.vue
git commit -m "feat: geocode sender and receiver addresses before order submit"
```

---

### Task 4: Modify HubFormModal.vue to geocode on submit

**Files:**
- Modify: `frontend/src/components/HubFormModal.vue`

- [ ] **Step 1: Add the geocode import and error state**

Add to the `<script setup>` block (after existing imports):

```typescript
import { geocodeAddress } from "@/lib/geocode";
```

Add after the `submitError` computed:

```typescript
const geocodeError = ref("");
const geocoding = ref(false);
```

- [ ] **Step 2: Update handleSubmit to geocode before sending**

Replace the existing `handleSubmit` function:

```typescript
async function handleSubmit() {
  geocodeError.value = "";
  geocoding.value = true;

  let coords: { lat: number; lng: number };
  try {
    coords = await geocodeAddress(
      address.value,
      "",
      "",
    );
    // For hubs, address is a free-text field (e.g. "Haven 12, Rotterdam, NL"),
    // so we pass it as subDistrict with empty district/province.
  } catch (e) {
    geocodeError.value = e instanceof Error ? e.message : "Could not resolve address.";
    geocoding.value = false;
    return;
  }

  const data = {
    name: name.value,
    carrierId,
    address: address.value,
    coords,
    capacity: capacity.value,
    currentUtilization: existing.value?.currentUtilization ?? 0,
    status: status.value,
  };

  try {
    if (isEditing.value && props.hubId) {
      await updateHub.mutateAsync({ id: props.hubId, data });
      toast.success("Hub updated");
    } else {
      await createHub.mutateAsync(data);
      toast.success("Hub created");
    }
    emit("close");
  } catch {
    // Mutation errors are surfaced via submitError; modal stays open
  } finally {
    geocoding.value = false;
  }
}
```

- [ ] **Step 3: Add geocode error display in the template**

After the address input and before the capacity/status grid:

```html
<div>
  <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground"
    >Address</label
  >
  <Input v-model="address" class="mt-1.5 font-mono text-sm" placeholder="Full address" />
  <p
    v-if="geocodeError"
    class="mt-1.5 rounded-md bg-destructive/15 px-3 py-2 font-mono text-xs text-destructive"
  >
    {{ geocodeError }}
  </p>
</div>
```

- [ ] **Step 4: Update submit button text during geocoding**

Replace the existing button:

```html
<Button :disabled="!name || submitPending || geocoding" @click="handleSubmit">
  {{
    geocoding
      ? "Resolving address\u2026"
      : submitPending
        ? "Saving\u2026"
        : isEditing
          ? "Update Hub"
          : "Create Hub"
  }}
</Button>
```

- [ ] **Step 5: Verify the build compiles**

```bash
cd frontend && npm run build
```

Expected: Build succeeds.

- [ ] **Step 6: Commit**

```bash
git add frontend/src/components/HubFormModal.vue
git commit -m "feat: geocode hub address before hub create/update"
```

---

## Verification

After all tasks are committed:

1. **Start the backend:** `cd backend && go run .`
2. **Start the frontend:** `cd frontend && npm run dev`
3. **Set the API key** in `frontend/.env` with a valid OpenCage key
4. **Create an order** — fill in a real Thai address (e.g. Bang Rak, Bang Rak, Bangkok) and submit. Check the backend DB to confirm `customer_lat` and `customer_lng` are real coordinates (not 0,0).
5. **Create a hub** — fill in a real address and submit. Check the backend DB for real `lat`/`lng`.
6. **Check the map** — navigate to OrderDetailView for the created order and verify the route markers appear at correct locations.
7. **Test error handling** — enter a garbage address and submit. Verify the inline error appears and the form is NOT submitted.
