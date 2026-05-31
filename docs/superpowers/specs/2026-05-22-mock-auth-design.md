# Mock Authentication Design

**Date:** 2026-05-22
**Project:** ship-it-out / Harbor Ops
**Status:** Approved

## Overview

Add a mock authentication system to gate order management actions. A Pinia store manages a boolean session flag persisted to localStorage. An `AuthModal` with login/signup tabs appears when an unauthenticated user tries to create or edit an order. Authenticated users see the actions column (edit/delete) in the orders table; unauthenticated users see a read-only table.

## Approach

Pinia store + modal component. Pinia is already installed and wired up in `main.ts`. The store provides global access to auth state across components.

## Auth Store

New file: `src/stores/auth.ts`

```typescript
useAuthStore {
  // state
  isAuthenticated: boolean  // default false, persisted to localStorage

  // actions
  login(email: string, password: string): boolean
  signup(name: string, email: string, password: string): boolean
  logout(): void
}
```

- **Persistence:** `isAuthenticated` is synced to `localStorage` key `"harborops_auth"` via a `watch` or direct `localStorage.setItem` calls in each action
- **login():** Accepts any non-empty email and password. Returns `false` if either is empty, otherwise sets `isAuthenticated = true` and returns `true`.
- **signup():** Accepts non-empty name, email, password. Returns `false` if any field empty, otherwise sets `isAuthenticated = true` and returns `true`.
- **logout():** Clears `isAuthenticated`. Clears `localStorage`.
- Init: On store creation, reads `localStorage` to restore session.

## AuthModal Component

New file: `src/components/AuthModal.vue`

**Layout:** Modal overlay matching the existing pattern (same as `AssignDriverModal.vue`, `HubFormModal.vue`).

**Two tabs:** Login and Signup. Tab headers at the top of the modal.

**Login tab fields:**
- Email (text input, shadcn Input)
- Password (password input, shadcn Input)
- "Sign In" submit button
- "Don't have an account? Sign up" link switches to Signup tab

**Signup tab fields:**
- Name (text input)
- Email (text input)
- Password (password input)
- Confirm Password (password input)
- "Create Account" submit button
- "Already have an account? Sign in" link switches to Login tab

**Validation:**
- All fields required (checked on submit)
- Signup: passwords must match
- On validation failure: inline error messages

**Flow:**
1. User clicks "New Order" while unauthenticated → `showAuthModal = true`
2. User fills form, submits → `login()`/`signup()` returns true
3. Modal closes → navigate to `/orders/create`
4. If login fails (empty fields) → inline error "Please fill in all fields"
5. Cancel/click outside → modal closes, stays on orders page

## Auth Gating

### OrdersView

- **"New Order" button:** If `!store.isAuthenticated`, set `showAuthModal = true` instead of navigating. If authenticated, navigates to `/orders/create` as before.
- **Actions column:** The entire Actions column header and cell content are conditionally rendered only when `store.isAuthenticated` is true. The table grid adjusts its column count accordingly.
  - Authenticated: 7 columns (Order ID, Tracking, Customer, Route, Status, ETA, Actions)
  - Unauthenticated: 6 columns (no Actions)
- Create and edit routes are NOT protected at the route level — the edit route from URL still works. The gating is purely for discoverability in the list view.

### SiteHeader

- Right side of the header nav, add:
  - Unauthenticated: "Sign in" button → navigates to `/orders` (user clicks New Order there to trigger modal)
  - Authenticated: small user indicator + "Sign out" button → calls `store.logout()`

## File Structure

```
src/
  stores/
    auth.ts                # Pinia auth store (new)
  components/
    AuthModal.vue          # Login/signup modal (new)
  views/
    OrdersView.vue         # Modified: auth gating, modal trigger
  components/
    SiteHeader.vue         # Modified: sign in/out button
```

## Error & Loading States

| Scenario | Behavior |
|---|---|
| Login empty fields | Inline error: "Please fill in all fields" |
| Signup empty fields | Inline error: "Please fill in all fields" |
| Signup password mismatch | Inline error: "Passwords do not match" |
| Login/signup success | Modal closes, navigates to `/orders/create` |
| Logout | Clears session, page stays, actions disappear |
| Page reload | Session restored from localStorage |

## Non-Goals

- No real password validation or hashing
- No protected routes (route guards)
- No backend integration
- No user roles or permissions
- No session expiry
- No registration limit
