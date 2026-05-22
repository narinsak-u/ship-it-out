# Real Auth Integration Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace mock Pinia auth with real Go backend API using httpOnly JWT cookies.

**Architecture:** Backend sets JWT as httpOnly cookie on login/register. Frontend calls `GET /api/auth/me` to validate sessions. Guest mode is a sessionStorage flag for read-only mock data access.

**Tech Stack:** Go (Fiber, JWT), Vue 3 (Pinia), Tailwind CSS v4

**Plan split:** This plan covers only auth integration. Data API integration is deferred.

---

## File Structure

### Backend (create/modify):
| File | Action | Responsibility |
|------|--------|----------------|
| `internal/middleware/cors.go` | Modify | Restrict origin to `localhost:5173`, allow credentials |
| `internal/middleware/auth.go` | Modify | Fall back to jwt cookie when Authorization header missing |
| `internal/auth/handler.go` | Modify | Add Me/Logout handlers, cookie-setting on Login/Register |
| `cmd/server/main.go` | Modify | Register /auth/me and /auth/logout routes |

### Frontend (create/modify):
| File | Action | Responsibility |
|------|--------|----------------|
| `src/lib/api/client.ts` | Create | Shared `fetch` wrapper with `credentials: 'include'` |
| `src/stores/auth.ts` | Rewrite | Real API calls, `/me` session check, guest mode |
| `src/components/AuthModal.vue` | Modify | Async handlers, loading state on submit |
| `src/components/SiteHeader.vue` | Modify | Loading skeleton, guest label |
| `src/views/OrdersView.vue` | Modify | Guest-aware New Order button |

---

### Task 1: Update CORS for cookie credentials

**Files:**
- Modify: `backend/internal/middleware/cors.go`

- [ ] **Step 1: Replace CORS with origin-specific + credentials**

```go
package middleware

import "github.com/gofiber/fiber/v2"

func CORS() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Set("Access-Control-Allow-Credentials", "true")
		c.Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		if c.Method() == "OPTIONS" {
			return c.SendStatus(204)
		}
		return c.Next()
	}
}
```

- [ ] **Step 2: Verify it compiles**

Run: `cd backend && go build ./...`
Expected: no errors

- [ ] **Step 3: Commit**

```bash
git add backend/internal/middleware/cors.go
git commit -m "fix: restrict CORS to frontend origin with credentials for cookie auth"
```

---

### Task 2: Update AuthRequired middleware to read jwt cookie as fallback

**Files:**
- Modify: `backend/internal/middleware/auth.go`

- [ ] **Step 1: Modify middleware to check jwt cookie when Authorization header is missing**

```go
package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/narinsak-u/backend/internal/config"
	"github.com/narinsak-u/backend/pkg/utils"
)

func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		tokenStr := ""

		if auth != "" && strings.HasPrefix(auth, "Bearer ") {
			tokenStr = strings.TrimPrefix(auth, "Bearer ")
		} else {
			// Fallback to jwt cookie
			tokenStr = c.Cookies("jwt")
		}

		if tokenStr == "" {
			return utils.Error(c, 401, "missing or invalid token")
		}

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(config.App.JWTSecret), nil
		})
		if err != nil || !token.Valid {
			return utils.Error(c, 401, "invalid or expired token")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return utils.Error(c, 401, "invalid token claims")
		}

		c.Locals("user_id", claims["user_id"])
		c.Locals("role", claims["role"])
		return c.Next()
	}
}
```

- [ ] **Step 2: Verify it compiles**

Run: `cd backend && go build ./...`
Expected: no errors

- [ ] **Step 3: Commit**

```bash
git add backend/internal/middleware/auth.go
git commit -m "feat: fall back to jwt cookie when Authorization header is missing"
```

---

### Task 3: Add cookie-setting, Me handler, and Logout handler to backend auth

**Files:**
- Modify: `backend/internal/auth/handler.go`

- [ ] **Step 1: Add cookie helper constant and import `net/http`**

Add import `"net/http"` to the import block, then add at package level:

```go
const cookieName = "jwt"
```

- [ ] **Step 2: Add `setAuthCookie` helper**

Add the function:

```go
func setAuthCookie(c *fiber.Ctx, token string) {
	c.Cookie(&fiber.Cookie{
		Name:     cookieName,
		Value:    token,
		Path:     "/",
		HTTPOnly: true,
		SameSite: "Lax",
		MaxAge:   86400,
	})
}
```

- [ ] **Step 3: Modify `Login` to set cookie after JWT generation**

After `tokenStr` is generated (line 81), add the cookie call:

```go
	setAuthCookie(c, tokenStr)
```

The full modified function should look like:

```go
func Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, 400, "invalid request body")
	}

	var user models.User
	if result := database.DB.Where("email = ?", req.Email).First(&user); result.Error != nil {
		return utils.Error(c, 401, "invalid email or password")
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		return utils.Error(c, 401, "invalid email or password")
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(config.App.JWTTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(config.App.JWTSecret))
	if err != nil {
		return utils.Error(c, 500, "failed to generate token")
	}

	setAuthCookie(c, tokenStr)

	return utils.Success(c, fiber.Map{
		"token": tokenStr,
		"user":  user,
	})
}
```

- [ ] **Step 4: Modify `Register` to set cookie after user creation**

After `database.DB.Create(&user)` succeeds, generate a JWT and set the cookie:

```go
func Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, 400, "invalid request body")
	}
	if req.Name == "" || req.Email == "" || req.Password == "" {
		return utils.Error(c, 400, "name, email, and password are required")
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return utils.Error(c, 500, "failed to hash password")
	}

	role := req.Role
	if role == "" {
		role = "customer"
	}

	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hash,
		Role:     role,
	}

	if result := database.DB.Create(&user); result.Error != nil {
		return utils.Error(c, 409, "email already registered")
	}

	// Auto-login: generate JWT and set cookie
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(config.App.JWTTTL).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(config.App.JWTSecret))
	if err != nil {
		return utils.Error(c, 500, "failed to generate token")
	}

	setAuthCookie(c, tokenStr)

	return utils.Success(c, fiber.Map{"user": user})
}
```

- [ ] **Step 5: Add `Me` handler**

```go
func Me(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return utils.Error(c, 401, "not authenticated")
	}

	var user models.User
	if result := database.DB.First(&user, userID); result.Error != nil {
		return utils.Error(c, 404, "user not found")
	}

	return utils.Success(c, fiber.Map{
		"id":         user.ID,
		"name":       user.Name,
		"email":      user.Email,
		"role":       user.Role,
		"created_at": user.CreatedAt,
	})
}
```

- [ ] **Step 6: Add `Logout` handler**

```go
func Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/",
		HTTPOnly: true,
		SameSite: "Lax",
		MaxAge:   0,
	})
	return utils.Success(c, fiber.Map{"message": "logged out"})
}
```

- [ ] **Step 7: Verify it compiles**

Run: `cd backend && go build ./...`
Expected: no errors

- [ ] **Step 8: Commit**

```bash
git add backend/internal/auth/handler.go
git commit -m "feat: set httpOnly JWT cookie on login/register, add /me and /logout handlers"
```

---

### Task 4: Register new auth routes

**Files:**
- Modify: `backend/cmd/server/main.go`

- [ ] **Step 1: Add /me and /logout routes**

In `main.go`, after the existing auth routes:

```go
authGroup.Get("/me", middleware.AuthRequired(), auth.Me)
authGroup.Post("/logout", auth.Logout)
```

The full `/auth` group should look like:

```go
authGroup := api.Group("/auth")
authGroup.Post("/register", auth.Register)
authGroup.Post("/login", auth.Login)
authGroup.Get("/me", middleware.AuthRequired(), auth.Me)
authGroup.Post("/logout", auth.Logout)
```

- [ ] **Step 2: Verify it compiles**

Run: `cd backend && go build ./...`
Expected: no errors

- [ ] **Step 3: Commit**

```bash
git add backend/cmd/server/main.go
git commit -m "feat: register /auth/me and /auth/logout routes"
```

---

### Task 5: Create frontend API client

**Files:**
- Create: `frontend/src/lib/api/client.ts`

- [ ] **Step 1: Create the shared API client**

```ts
const BASE = 'http://localhost:8080/api';

interface ApiSuccess<T> {
  data: T;
  error?: never;
}

interface ApiError {
  error: string;
  data?: never;
}

type ApiResult<T> = ApiSuccess<T> | ApiError;

async function request<T = unknown>(path: string, options?: RequestInit): Promise<ApiResult<T>> {
  try {
    const res = await fetch(`${BASE}${path}`, {
      credentials: 'include',
      headers: { 'Content-Type': 'application/json', ...options?.headers },
      ...options,
    });
    const json = await res.json();
    if (!res.ok) return { error: json.error || `Request failed (${res.status})` };
    return { data: json.data as T };
  } catch {
    return { error: 'Network error — is the backend running?' };
  }
}

export const api = {
  get: <T = unknown>(path: string) => request<T>(path),
  post: <T = unknown>(path: string, body?: unknown) =>
    request<T>(path, { method: 'POST', body: body ? JSON.stringify(body) : undefined }),
  del: <T = unknown>(path: string) => request<T>(path, { method: 'DELETE' }),
};
```

- [ ] **Step 2: Verify TypeScript**

Run: `cd frontend && npx vue-tsc --noEmit`
Expected: no errors

- [ ] **Step 3: Commit**

```bash
git add frontend/src/lib/api/client.ts
git commit -m "feat: add shared API client with credentials-include"
```

---

### Task 6: Rewrite auth store with real API calls

**Files:**
- Modify: `frontend/src/stores/auth.ts`

- [ ] **Step 1: Replace the entire store content**

```ts
import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import { api } from '@/lib/api/client';

export interface AuthUser {
  id: number;
  name: string;
  email: string;
  role: string;
  created_at: string;
}

export const useAuthStore = defineStore('auth', () => {
  const user = ref<AuthUser | null>(null);
  const loading = ref(true);
  const error = ref('');
  const isGuest = ref(sessionStorage.getItem('harborops_guest') === 'true');

  const isAuthenticated = computed(() => user.value !== null);

  function init() {
    loading.value = true;
    if (isGuest.value) {
      loading.value = false;
      return;
    }
    api.get<AuthUser>('/auth/me').then((res) => {
      if (res.data) {
        user.value = res.data;
      }
      loading.value = false;
    });
  }

  async function login(email: string, password: string): Promise<string | null> {
    error.value = '';
    const res = await api.post<{ user: AuthUser }>('/auth/login', { email, password });
    if (res.error) {
      error.value = res.error;
      return res.error;
    }
    // Verify session by calling /me
    const me = await api.get<AuthUser>('/auth/me');
    if (me.data) {
      user.value = me.data;
      return null;
    }
    error.value = 'Login failed — session not established';
    return error.value;
  }

  async function signup(name: string, email: string, password: string): Promise<string | null> {
    error.value = '';
    const res = await api.post<{ user: AuthUser }>('/auth/register', { name, email, password });
    if (res.error) {
      error.value = res.error;
      return res.error;
    }
    const me = await api.get<AuthUser>('/auth/me');
    if (me.data) {
      user.value = me.data;
      return null;
    }
    error.value = 'Signup failed — session not established';
    return error.value;
  }

  async function logout() {
    await api.post('/auth/logout');
    user.value = null;
    isGuest.value = false;
    sessionStorage.removeItem('harborops_guest');
  }

  function enterGuestMode() {
    isGuest.value = true;
    sessionStorage.setItem('harborops_guest', 'true');
    loading.value = false;
  }

  return { user, loading, error, isGuest, isAuthenticated, init, login, signup, logout, enterGuestMode };
});
```

- [ ] **Step 2: Verify TypeScript**

Run: `cd frontend && npx vue-tsc --noEmit`
Expected: no errors

- [ ] **Step 3: Commit**

```bash
git add frontend/src/stores/auth.ts
git commit -m "feat: rewrite auth store with real backend API calls"
```

---

### Task 7: Update AuthModal for async API calls

**Files:**
- Modify: `frontend/src/components/AuthModal.vue`

- [ ] **Step 1: Rewrite the script section**

Replace the entire `<script setup lang="ts">` block:

```vue
<script setup lang="ts">
import { ref } from "vue";
import { X, LogIn, UserPlus, Loader2 } from "lucide-vue-next";
import { useAuthStore } from "@/stores/auth";
import Input from "@/components/ui/Input.vue";
import Button from "@/components/ui/Button.vue";

const emit = defineEmits<{ close: []; authenticated: []; guest: [] }>();

const store = useAuthStore();

type Tab = "login" | "signup";
const activeTab = ref<Tab>("login");
const submitting = ref(false);

// Login fields
const loginEmail = ref("");
const loginPassword = ref("");

// Signup fields
const signupName = ref("");
const signupEmail = ref("");
const signupPassword = ref("");

function switchTab(tab: Tab) {
  activeTab.value = tab;
  store.error = "";
}

async function handleLogin() {
  submitting.value = true;
  store.error = "";
  const err = await store.login(loginEmail.value, loginPassword.value);
  submitting.value = false;
  if (!err) emit("authenticated");
}

async function handleSignup() {
  submitting.value = true;
  store.error = "";
  const err = await store.signup(
    signupName.value,
    signupEmail.value,
    signupPassword.value,
  );
  submitting.value = false;
  if (!err) emit("authenticated");
}

function handleGuest() {
  store.enterGuestMode();
  emit("guest");
}
</script>
```

- [ ] **Step 2: Update the submit buttons for loading state**

In the login form's submit button (`<Button type="submit" class="w-full gap-2">`), replace:

```html
<Button type="submit" class="w-full gap-2" :disabled="submitting">
  <LogIn v-if="!submitting" class="h-4 w-4" />
  <Loader2 v-else class="h-4 w-4 animate-spin" />
  {{ submitting ? 'Signing in...' : 'Sign In' }}
</Button>
```

In the signup form's submit button (`<Button type="submit" class="w-full gap-2">`), replace:

```html
<Button type="submit" class="w-full gap-2" :disabled="submitting">
  <UserPlus v-if="!submitting" class="h-4 w-4" />
  <Loader2 v-else class="h-4 w-4 animate-spin" />
  {{ submitting ? 'Creating account...' : 'Create Account' }}
</Button>
```

- [ ] **Step 3: Add `Loader2` to the icon import**

Edit the import line to include `Loader2`:

```ts
import { X, LogIn, UserPlus, Loader2 } from "lucide-vue-next";
```

- [ ] **Step 4: Verify TypeScript**

Run: `cd frontend && npx vue-tsc --noEmit`
Expected: no errors

- [ ] **Step 5: Commit**

```bash
git add frontend/src/components/AuthModal.vue
git commit -m "feat: wire AuthModal to real auth store with loading state"
```

---

### Task 8: Update SiteHeader for loading state and guest label

**Files:**
- Modify: `frontend/src/components/SiteHeader.vue`

- [ ] **Step 1: Update the header auth section**

Replace the content from `<div class="ml-4 flex items-center gap-2 border-l border-border pl-4">` to the end of `</header>`:

```vue
        <div class="ml-4 flex items-center gap-2 border-l border-border pl-4">
          <template v-if="authStore.loading">
            <span class="font-mono text-xs text-muted-foreground">...</span>
          </template>
          <template v-else-if="authStore.user">
            <span class="font-mono text-xs text-muted-foreground">Admin ({{ authStore.user.name }})</span>
            <button
              @click="authStore.logout(); router.push({ name: 'home' })"
              class="flex items-center gap-1.5 rounded-md px-3 py-1.5 font-mono text-xs text-muted-foreground transition-colors hover:text-foreground"
            >
              <LogOut class="h-3.5 w-3.5" /> Sign out
            </button>
          </template>
          <template v-else-if="authStore.isGuest">
            <span class="font-mono text-xs text-muted-foreground">Guest</span>
            <button
              @click="router.push({ name: 'orders' })"
              class="flex items-center gap-1.5 rounded-md px-3 py-1.5 font-mono text-xs transition-colors text-primary hover:text-foreground"
            >
              <LogIn class="h-3.5 w-3.5" /> Sign in
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

- [ ] **Step 2: Verify TypeScript**

Run: `cd frontend && npx vue-tsc --noEmit`
Expected: no errors

- [ ] **Step 3: Commit**

```bash
git add frontend/src/components/SiteHeader.vue
git commit -m "feat: add loading skeleton and guest label to SiteHeader"
```

---

### Task 9: Update OrdersView for guest-aware New Order

**Files:**
- Modify: `frontend/src/views/OrdersView.vue`

- [ ] **Step 1: Update New Order button logic**

Replace the hero section New Order buttons (the `v-if`/`v-else` block starting with `New Order`):

```vue
          <div v-if="authStore.user" class="shrink-0">
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

Also update the actions column header visibility:

Change `v-if="authStore.isAuthenticated"` on the actions column header to `v-if="authStore.user"` (line 130).

Change `v-if="authStore.isAuthenticated"` on the actions column content to `v-if="authStore.user"` (line 154).

- [ ] **Step 2: Verify TypeScript**

Run: `cd frontend && npx vue-tsc --noEmit`
Expected: no errors

- [ ] **Step 3: Commit**

```bash
git add frontend/src/views/OrdersView.vue
git commit -m "feat: use authStore.user for OrdersView action gating"
```

---

### Task 10: Initialize auth store in App.vue

**Files:**
- Modify: `frontend/src/App.vue`

- [ ] **Step 1: Call `init()` on mount**

```vue
<script setup lang="ts">
import { onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import SiteHeader from '@/components/SiteHeader.vue'
import SiteFooter from '@/components/SiteFooter.vue'

const authStore = useAuthStore()

onMounted(() => {
  authStore.init()
})
</script>
```

- [ ] **Step 2: Verify TypeScript**

Run: `cd frontend && npx vue-tsc --noEmit`
Expected: no errors

- [ ] **Step 3: Commit**

```bash
git add frontend/src/App.vue
git commit -m "feat: initialize auth store on app mount"
```

---

### Task 11: Build verification

- [ ] **Step 1: Verify backend compiles**

Run: `cd backend && go build ./...`
Expected: no errors

- [ ] **Step 2: Verify frontend builds**

Run: `cd frontend && npm run build`
Expected: `vue-tsc` + Vite build succeeds with no errors

- [ ] **Step 3: Commit any final fixes**

```bash
git add -A
git commit -m "chore: fix build errors after auth integration"
```
