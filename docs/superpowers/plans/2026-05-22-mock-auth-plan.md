# Mock Authentication Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add mock authentication to gate order management actions behind a login/signup modal.

**Architecture:** Pinia store with `isAuthenticated` boolean persisted to localStorage. `AuthModal.vue` with login/signup tabs. OrdersView conditionally renders actions based on auth state.

**Tech Stack:** Vue 3, Pinia, shadcn-vue Input/Button, lucide-vue-next

---

### Task 1: Create Pinia auth store

**Files:**
- Create: `frontend/src/stores/auth.ts`

- [ ] **Step 1: Create auth store**

Create `frontend/src/stores/auth.ts`:

```typescript
import { defineStore } from 'pinia';
import { ref, watch } from 'vue';

const STORAGE_KEY = 'harborops_auth';

export const useAuthStore = defineStore('auth', () => {
  const isAuthenticated = ref(localStorage.getItem(STORAGE_KEY) === 'true');

  watch(isAuthenticated, (val) => {
    localStorage.setItem(STORAGE_KEY, val ? 'true' : '');
  });

  function login(email: string, password: string): boolean {
    if (!email.trim() || !password.trim()) return false;
    isAuthenticated.value = true;
    return true;
  }

  function signup(name: string, email: string, password: string, confirmPassword: string): string | null {
    if (!name.trim() || !email.trim() || !password.trim() || !confirmPassword.trim()) {
      return 'Please fill in all fields';
    }
    if (password !== confirmPassword) {
      return 'Passwords do not match';
    }
    isAuthenticated.value = true;
    return null;
  }

  function logout() {
    isAuthenticated.value = false;
  }

  return { isAuthenticated, login, signup, logout };
});
```

- [ ] **Step 2: Commit**

```bash
git add frontend/src/stores/auth.ts
git commit -m "feat: add Pinia auth store with localStorage persistence"
```

---

### Task 2: Create AuthModal component

**Files:**
- Create: `frontend/src/components/AuthModal.vue`

- [ ] **Step 1: Create AuthModal.vue**

Create `frontend/src/components/AuthModal.vue`:

```typescript
<script setup lang="ts">
import { ref } from 'vue'
import { X, LogIn, UserPlus } from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'
import Input from '@/components/ui/Input.vue'
import Button from '@/components/ui/Button.vue'

const emit = defineEmits<{ close: []; authenticated: []; guest: [] }>()

const store = useAuthStore()

type Tab = 'login' | 'signup'
const activeTab = ref<Tab>('login')

// Login fields
const loginEmail = ref('')
const loginPassword = ref('')

// Signup fields
const signupName = ref('')
const signupEmail = ref('')
const signupPassword = ref('')
const signupConfirm = ref('')

const error = ref('')

function switchTab(tab: Tab) {
  activeTab.value = tab
  error.value = ''
}

function handleLogin() {
  const ok = store.login(loginEmail.value, loginPassword.value)
  if (!ok) {
    error.value = 'Please fill in all fields'
    return
  }
  emit('authenticated')
}

function handleSignup() {
  const err = store.signup(signupName.value, signupEmail.value, signupPassword.value, signupConfirm.value)
  if (err) {
    error.value = err
    return
  }
  emit('authenticated')
}

function handleGuest() {
  emit('guest')
}
</script>

<template>
  <div class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm" @click.self="emit('close')">
    <div class="w-full max-w-md rounded-xl border border-border bg-card p-6 shadow-elegant">
      <div class="flex items-center justify-between">
        <h2 class="font-mono text-lg font-semibold">{{ activeTab === 'login' ? 'Sign in' : 'Create account' }}</h2>
        <button @click="emit('close')" class="text-muted-foreground hover:text-foreground">
          <X class="h-5 w-5" />
        </button>
      </div>

      <!-- Tab bar -->
      <div class="mt-5 flex gap-0 border-b border-border">
        <button
          @click="switchTab('login')"
          class="flex-1 pb-3 font-mono text-sm transition-colors"
          :class="activeTab === 'login' ? 'text-foreground border-b-2 border-primary' : 'text-muted-foreground'"
        >
          Sign In
        </button>
        <button
          @click="switchTab('signup')"
          class="flex-1 pb-3 font-mono text-sm transition-colors"
          :class="activeTab === 'signup' ? 'text-foreground border-b-2 border-primary' : 'text-muted-foreground'"
        >
          Sign Up
        </button>
      </div>

      <!-- Login form -->
      <form v-if="activeTab === 'login'" @submit.prevent="handleLogin" class="mt-6 space-y-4">
        <div>
          <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground">Email</label>
          <Input v-model="loginEmail" type="email" class="mt-1.5 font-mono text-sm" placeholder="admin@harborops.io" />
        </div>
        <div>
          <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground">Password</label>
          <Input v-model="loginPassword" type="password" class="mt-1.5 font-mono text-sm" placeholder="&#xb7;&#xb7;&#xb7;&#xb7;&#xb7;&#xb7;&#xb7;&#xb7;" />
        </div>

        <p v-if="error" class="font-mono text-xs text-destructive">{{ error }}</p>

        <Button type="submit" class="w-full gap-2">
          <LogIn class="h-4 w-4" /> Sign In
        </Button>

        <div class="relative my-4">
          <div class="absolute inset-0 flex items-center">
            <span class="w-full border-t border-border" />
          </div>
          <div class="relative flex justify-center text-xs uppercase">
            <span class="bg-card px-2 font-mono text-muted-foreground">or</span>
          </div>
        </div>

        <Button type="button" variant="ghost" class="w-full" @click="handleGuest">
          Continue as Guest
        </Button>

        <p class="text-center font-mono text-xs text-muted-foreground">
          Don't have an account?
          <button type="button" @click="switchTab('signup')" class="text-primary hover:underline">Sign up</button>
        </p>
      </form>

      <!-- Signup form -->
      <form v-else @submit.prevent="handleSignup" class="mt-6 space-y-4">
        <div>
          <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground">Name</label>
          <Input v-model="signupName" class="mt-1.5 font-mono text-sm" placeholder="Your name" />
        </div>
        <div>
          <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground">Email</label>
          <Input v-model="signupEmail" type="email" class="mt-1.5 font-mono text-sm" placeholder="email@example.com" />
        </div>
        <div>
          <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground">Password</label>
          <Input v-model="signupPassword" type="password" class="mt-1.5 font-mono text-sm" />
        </div>
        <div>
          <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground">Confirm Password</label>
          <Input v-model="signupConfirm" type="password" class="mt-1.5 font-mono text-sm" />
        </div>

        <p v-if="error" class="font-mono text-xs text-destructive">{{ error }}</p>

        <Button type="submit" class="w-full gap-2">
          <UserPlus class="h-4 w-4" /> Create Account
        </Button>

        <div class="relative my-4">
          <div class="absolute inset-0 flex items-center">
            <span class="w-full border-t border-border" />
          </div>
          <div class="relative flex justify-center text-xs uppercase">
            <span class="bg-card px-2 font-mono text-muted-foreground">or</span>
          </div>
        </div>

        <Button type="button" variant="ghost" class="w-full" @click="handleGuest">
          Continue as Guest
        </Button>

        <p class="text-center font-mono text-xs text-muted-foreground">
          Already have an account?
          <button type="button" @click="switchTab('login')" class="text-primary hover:underline">Sign in</button>
        </p>
      </form>
    </div>
  </div>
</template>
```

- [ ] **Step 2: Commit**

```bash
git add frontend/src/components/AuthModal.vue
git commit -m "feat: add AuthModal with login and signup tabs"
```

---

### Task 3: Update OrdersView with auth gating

**Files:**
- Modify: `frontend/src/views/OrdersView.vue`

- [ ] **Step 1: Add auth gating to OrdersView**

Read `frontend/src/views/OrdersView.vue`. Make these changes:

**Script section** — add `useAuthStore` import and `showAuthModal` ref:

After `import { cn } from '@/lib/utils'`, add:
```typescript
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const showAuthModal = ref(false)
```

**"New Order" button** — change the `RouterLink` to a conditional handler. Replace the existing button section:

```html
          <div v-if="authStore.isAuthenticated" class="shrink-0">
            <RouterLink :to="{ name: 'order-create' }">
              <Button class="gap-2">
                <Plus class="h-4 w-4" /> New Order
              </Button>
            </RouterLink>
          </div>
          <div v-else class="shrink-0">
            <Button class="gap-2" @click="showAuthModal = true">
              <Plus class="h-4 w-4" /> New Order
            </Button>
          </div>
```

**Actions column** — conditionally render the Actions column header and cells. Replace the header line:

```html
          <span v-if="authStore.isAuthenticated" class="text-right">Actions</span>
```

Replace the actions cell div:

```html
            <div v-if="authStore.isAuthenticated" class="flex justify-end gap-1">
              <button
                @click.stop="router.push({ name: 'order-edit', params: { orderId: o.id } })"
                class="rounded p-1.5 text-muted-foreground hover:text-primary"
              >
                <Pencil class="h-4 w-4" />
              </button>
              <button
                @click.stop="deleteOrder(o.id)"
                class="rounded p-1.5 text-muted-foreground hover:text-destructive"
              >
                <Trash2 class="h-4 w-4" />
              </button>
            </div>
```

**Modal** — add the AuthModal component import at the top of the script:

```typescript
import AuthModal from '@/components/AuthModal.vue'
```

And add the modal at the end of the template (before the closing `</div>`):

```html
    <AuthModal
      v-if="showAuthModal"
      @close="showAuthModal = false"
      @authenticated="onAuthenticated"
      @guest="onGuest"
    />
  </div>
Add the handlers in the script:

```typescript
function onAuthenticated() {
  showAuthModal = false
  router.push({ name: 'order-create' })
}

function onGuest() {
  showAuthModal = false
  router.push({ name: 'order-create' })
}
```

Make sure the existing `RouterLink` import is present (it should already be from the existing code since `New Order` button uses `RouterLink`).

- [ ] **Step 2: Commit**

```bash
git add frontend/src/views/OrdersView.vue
git commit -m "feat: add auth gating to orders view"
```

---

### Task 4: Update SiteHeader with sign in/out

**Files:**
- Modify: `frontend/src/components/SiteHeader.vue`

- [ ] **Step 1: Add auth controls to header**

Read `frontend/src/components/SiteHeader.vue`. Make these changes:

**Script section** — add auth store import:

After `const route = useRoute();`, add:
```typescript
import { useAuthStore } from '@/stores/auth';
import { LogIn, LogOut } from 'lucide-vue-next';
import { useRouter } from 'vue-router';

const authStore = useAuthStore();
const router = useRouter();
```

Also keep the existing `Package` icon import and add `LogIn, LogOut`.

Update the Package import:
```typescript
import { Package, LogIn, LogOut } from 'lucide-vue-next';
```

**Template** — add auth controls after the nav links, before the closing `</nav>`:

```html
        <div class="ml-4 flex items-center gap-2 border-l border-border pl-4">
          <template v-if="authStore.isAuthenticated">
            <span class="font-mono text-xs text-muted-foreground">Admin</span>
            <button
              @click="authStore.logout()"
              class="flex items-center gap-1.5 rounded-md px-3 py-1.5 font-mono text-xs text-muted-foreground transition-colors hover:text-foreground"
            >
              <LogOut class="h-3.5 w-3.5" /> Sign out
            </button>
          </template>
          <button
            v-else
            @click="router.push({ name: 'orders' })"
            class="flex items-center gap-1.5 rounded-md px-3 py-1.5 font-mono text-xs transition-colors text-primary hover:text-foreground"
          >
            <LogIn class="h-3.5 w-3.5" /> Sign in
          </button>
        </div>
```

- [ ] **Step 2: Commit**

```bash
git add frontend/src/components/SiteHeader.vue
git commit -m "feat: add sign in/out controls to site header"
```

---

### Task 5: Build verification

- [ ] **Step 1: Run build**

```bash
cd frontend && npm run build
```

Expected: clean build with no type errors.

If there are issues, fix them and re-run.

- [ ] **Step 2: Push**

```bash
git push
```
