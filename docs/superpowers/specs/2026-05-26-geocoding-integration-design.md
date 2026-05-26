# OpenCage Geocoding Integration

**Date:** 2026-05-26
**Status:** Approved

## Overview

Integrate the OpenCage Geocoding API to convert addresses into real latitude/longitude coordinates before creating orders and hubs. Currently the app sends `coords: { lat: 0, lng: 0 }` as placeholders — this replaces those with real geocoded coordinates so the Leaflet map in OrderDetailView displays accurate routes.

The `opencage-api-client` npm package (v2.1.2) is already installed in the frontend.

## Decisions

| Decision | Choice |
|----------|--------|
| Trigger | Auto on form submit (not manual button, not auto-while-typing) |
| API key location | Client-side (`VITE_OPENCAGE_API_KEY`) |
| Geocode failure | Block submit + show inline error on the address group |
| Geocoding order | Parallel for order form (sender + receiver), sequential for hub form |
| Query string | Composed from `subDistrict, district, province` |

## Data Flow

### Order form

```
User clicks "Create Order"
  → validate form fields (existing)
  → disable button, show "Resolving addresses…"
  → Promise.allSettled([
      geocodeAddress(sender.subDistrict, sender.district, sender.province),
      geocodeAddress(receiver.subDistrict, receiver.district, receiver.province)
    ])
  → if both succeed:
      → replace coords with real lat/lng
      → emit OrderFormData (unchanged type)
      → existing OrderFormView.handleSubmit() flow
  → if any fail:
      → re-enable form
      → show inline error on failed address group
      → block submission
```

### Hub form

```
User clicks "Create Hub"
  → validate (name + address required)
  → disable button, show "Resolving address…"
  → geocodeAddress(hub address fields)
  → if success: replace coords, submit to backend
  → if failure: inline error, block submission
```

## Architecture

### New file: `frontend/src/lib/geocode.ts`

A single exported async function:

```typescript
export async function geocodeAddress(
  subDistrict: string,
  district: string,
  province: string,
): Promise<{ lat: number; lng: number }> {
  const q = `${subDistrict}, ${district}, ${province}`.trim();
  const key = import.meta.env.VITE_OPENCAGE_API_KEY;
  if (!key) throw new Error("OpenCage API key is not configured.");
  // Uses opencage-api-client library (v2.1.2)
  // If browser compatibility issues arise, fall back to raw fetch:
  // POST https://api.opencagedata.com/geocode/v1/json?q=...&key=...
  const { geocode } = await import("opencage-api-client");
  const data = await geocode({ q, key });
  if (!data.results || data.results.length === 0) {
    throw new Error("Could not resolve this address.");
  }
  const { lat, lng } = data.results[0].geometry;
  return { lat, lng };
}
```

- The key is read from `import.meta.env.VITE_OPENCAGE_API_KEY`
- Errors are caught and re-thrown with a user-friendly message
- The function is pure — no side effects, no state
- If `opencage-api-client` doesn't work in the Vite browser context, fall back to a raw `fetch()` call to the OpenCage REST API (same endpoint, same params)

### Modified file: `frontend/src/components/OrderForm.vue`

`handleSubmit()` changes:
1. After validation passes, call `geocodeAddress` for both sender and receiver in parallel (`Promise.allSettled`)
2. On success: merge real coords into the emitted `OrderFormData`
3. On failure: set a `geocodeErrors` reactive record, display inline error below the address group title, do NOT emit

A `geocodeErrors` ref is added: `ref<Record<string, string>>({})` with keys like `"sender"` and `"receiver"`.

The submit button shows `"Resolving addresses…"` while geocoding is in progress.

### Modified file: `frontend/src/components/HubFormModal.vue`

`handleSubmit()` changes:
1. After building data object, geocode the address before sending
2. On success: replace `coords: { lat: 0, lng: 0 }` with real coords, proceed with mutation
3. On failure: set a `geocodeError` ref, display inline error message, block submission

The submit button shows `"Resolving address…"` while geocoding.

### Environment variable

Add `VITE_OPENCAGE_API_KEY=<your-key>` to `frontend/.env` (documented in `.env.example`).

OpenCage free tier allows 2,500 requests/day — sufficient for development. For production, restrict the key to the app's domain in the OpenCage dashboard.

## Error handling

| Scenario | Behavior |
|----------|----------|
| Network error (no internet, API down) | Catch, show inline "Could not resolve this address. Check your connection and try again." |
| OpenCage returns no results | Catch, show inline "Could not resolve this address. Check the fields and try again." |
| OpenCage returns error (401, 429, etc.) | Catch, show inline "Location lookup failed. Please try again later." |
| One address fails, other succeeds | Block submit. Show error on failed group only. |

## Testing

- No test framework exists yet. Manual testing:
  - Create order with valid Thai address → coords in DB should be real (check map)
  - Create order with garbage address → error shown, form not submitted
  - Create hub with valid address → coords in DB should be real
  - Create hub with garbage address → error shown, modal stays open

## Files changed

| File | Change |
|------|--------|
| `frontend/src/lib/geocode.ts` | **New** — `geocodeAddress()` utility |
| `frontend/.env` | **New** — `VITE_OPENCAGE_API_KEY` |
| `frontend/.env.example` | **Updated** — add `VITE_OPENCAGE_API_KEY` |
| `frontend/src/components/OrderForm.vue` | **Modified** — geocode before submit |
| `frontend/src/components/HubFormModal.vue` | **Modified** — geocode before submit |
