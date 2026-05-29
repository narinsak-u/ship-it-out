# Frontend Test Coverage Analysis

**Date:** May 29, 2026
**Status:** Completed
**Scope:** Vitest (Unit/Integration) & Playwright (E2E)

---

## 1. Executive Summary

The Harbor Ops frontend features a modern, multi-layered testing suite that provides high confidence in the application's core functionality. All major layers—composables, stores, API clients, hooks, components, and views—have dedicated test files. The use of **Mock Service Worker (MSW)** for API interception is a highlight, allowing for realistic integration tests without a running backend.

While the "breadth" of coverage is excellent (27 test files, 107 tests passing), the "depth" in certain areas—specifically complex reactive logic and interactive E2E flows—needs enhancement to reach "production-grade" reliability.

---

## 2. Coverage Analysis by Layer

### Layer 1: Composables & Utilities
- **Status:** Excellent
- **Coverage:** `usePagination` and `useSearchFilter` are thoroughly tested for edge cases (empty arrays, boundary conditions, case-insensitivity).
- **Finding:** These are the most resilient parts of the codebase.

### Layer 2: Auth Store (Pinia)
- **Status:** Strong
- **Coverage:** Covers initialization, guest mode, login/signup success/failure, and logout.
- **Finding:** Correctly verifies that state updates match API responses.

### Layer 3: API Layer (Client & Mappers)
- **Status:** Strong
- **Coverage:** All 18+ endpoints are tested via MSW. Mappers are verified for date/timestamp formatting and field transformations.
- **Finding:** Reliable verification that the frontend can handle the contract defined by the Go backend.

### Layer 4: Hooks (TanStack Query)
- **Status:** Good
- **Coverage:** Basic verification that queries fetch data and mutations trigger API calls.
- **Finding:** These tests ensure the "plumbing" between the API and the UI components is working.

### Layer 5: Components
- **Status:** Mixed
- **Presentation Components:** (`StatusBadge`, `Pagination`, `ConfirmDialog`) have high-quality tests covering props and emits.
- **Complex Logic Components:** (`OrderForm`, `ThaiAddressGroup`, `AuthModal`) have **shallow coverage**. Tests verify they render sections, but often skip internal validation logic, async side-effects (zipcode lookup), and complex state transitions.

### Layer 6: Views (Pages)
- **Status:** Shallow
- **Coverage:** Verifies that pages render expected headings and sections.
- **Finding:** These serve as "smoke tests" for the router and top-level data fetching but don't exercise page-level interactivity.

### Layer 7: E2E Tests (Playwright)
- **Status:** Minimal
- **Coverage:** Basic navigation and visibility checks.
- **Finding:** No tests currently exercise a full business process (e.g., "Create Order -> Find in List -> View Detail").

---

## 3. Identified Gaps & Risks

1.  **Form Validation Logic:** `OrderForm.vue` and `AuthModal.vue` contain significant validation logic (required fields, password matching, weight constraints). These are currently untested at the unit level.
2.  **Reactive Watchers:** The `ThaiAddressGroup.vue` zipcode lookup logic (fetching sub-districts from `thai-data` on zipcode change) is complex but has no corresponding test cases.
3.  **UI Interactivity:**
    *   `OrdersView.vue` search debouncing (300ms) is not verified.
    *   Tab switching in `CarriersView.vue` is tested for rendering but not for the lazy-loading performance impact.
4.  **Error States & Toasts:** While API errors are tested at the client level, the UI's response to these errors (error banners, toast notifications) is largely unverified in component specs.
5.  **Analytics Transformations:** The data processing in `AnalyticsPanel.vue` (grouping by region, status distribution percentages) is complex enough to merit its own logic tests.

---

## 4. Suggestions for Improvement

### Immediate Actions (High Impact)
1.  **Strengthen `ThaiAddressGroup.spec.ts`**:
    *   Simulate entering a valid zipcode (e.g., "10100").
    *   Assert that the sub-district `<Select>` is populated with expected values from the `thai-data` package.
2.  **Deepen `OrderForm.spec.ts`**:
    *   Add "Validation" tests: verify that the "Create Order" button is disabled until all required fields are filled.
    *   Add "Geocoding Integration" tests: mock a geocode failure and verify the inline error message appears.
3.  **Implement a "Critical Path" E2E Flow**:
    *   Create a new test `tests/e2e/fulfillment.spec.ts` that:
        1. Logs in as a user.
        2. Fills and submits the "New Order" form.
        3. Verifies the new order appears in the `OrdersView` table.
        4. Navigates to the `OrderDetailView` for that order.

### Architectural Suggestions (Maintenance)
1.  **Extract Analytics Logic**: Move computed property logic from `AnalyticsPanel.vue` into a separate `src/lib/analytics-utils.ts` and add pure unit tests for the transformation functions.
2.  **Mock Toasts**: Add a global stub or mock for `vue-sonner` in `tests/setup.ts` to easily verify that success/error toasts are triggered during mutations.
3.  **Test Debounce**: Use `vi.useFakeTimers()` in `OrdersView.spec.ts` to verify that the API is only called once after a burst of typing.

---

## 5. Conclusion

The current test suite is a solid foundation that ensures the application "runs" and "talks to the API." By focusing on the "depth" of the interactive components and adding a few comprehensive E2E flows, the project can move from "well-tested" to "production-ready."
