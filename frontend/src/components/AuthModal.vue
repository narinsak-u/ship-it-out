<script setup lang="ts">
import { ref } from "vue";
import { X, LogIn, UserPlus } from "lucide-vue-next";
import { useAuthStore } from "@/stores/auth";
import Input from "@/components/ui/Input.vue";
import Button from "@/components/ui/Button.vue";

const emit = defineEmits<{ close: []; authenticated: []; guest: [] }>();

const store = useAuthStore();

type Tab = "login" | "signup";
const activeTab = ref<Tab>("login");

// Login fields
const loginEmail = ref("");
const loginPassword = ref("");

// Signup fields
const signupName = ref("");
const signupEmail = ref("");
const signupPassword = ref("");
const signupConfirm = ref("");

const error = ref("");

function switchTab(tab: Tab) {
  activeTab.value = tab;
  error.value = "";
}

function handleLogin() {
  const ok = store.login(loginEmail.value, loginPassword.value);
  if (!ok) {
    error.value = "Please fill in all fields";
    return;
  }
  emit("authenticated");
}

function handleSignup() {
  const err = store.signup(
    signupName.value,
    signupEmail.value,
    signupPassword.value,
    signupConfirm.value,
  );
  if (err) {
    error.value = err;
    return;
  }
  emit("authenticated");
}

function handleGuest() {
  emit("guest");
}
</script>

<template>
  <div
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm"
    @click.self="emit('close')"
  >
    <div
      class="w-full max-w-md rounded-xl border border-border bg-card p-6 shadow-elegant"
    >
      <div class="flex items-center justify-between">
        <h2 class="font-mono text-lg font-semibold">
          {{ activeTab === "login" ? "Sign in" : "Create account" }}
        </h2>
        <button
          @click="emit('close')"
          class="text-muted-foreground hover:text-foreground"
        >
          <X class="h-5 w-5" />
        </button>
      </div>

      <!-- Tab bar -->
      <div class="mt-5 flex gap-0 border-b border-border">
        <button
          @click="switchTab('login')"
          class="flex-1 pb-3 font-mono text-sm transition-colors"
          :class="
            activeTab === 'login'
              ? 'text-foreground border-b-2 border-primary'
              : 'text-muted-foreground'
          "
        >
          Sign In
        </button>
        <button
          @click="switchTab('signup')"
          class="flex-1 pb-3 font-mono text-sm transition-colors"
          :class="
            activeTab === 'signup'
              ? 'text-foreground border-b-2 border-primary'
              : 'text-muted-foreground'
          "
        >
          Sign Up
        </button>
      </div>

      <!-- Login form -->
      <form
        v-if="activeTab === 'login'"
        @submit.prevent="handleLogin"
        class="mt-6 space-y-4"
      >
        <div>
          <label
            class="font-mono text-xs uppercase tracking-widest text-muted-foreground"
            >Email</label
          >
          <Input
            v-model="loginEmail"
            type="email"
            class="mt-1.5 font-mono text-sm"
            placeholder="admin@harborops.io"
          />
        </div>
        <div>
          <label
            class="font-mono text-xs uppercase tracking-widest text-muted-foreground"
            >Password</label
          >
          <Input
            v-model="loginPassword"
            type="password"
            class="mt-1.5 font-mono text-sm"
            placeholder="&#xb7;&#xb7;&#xb7;&#xb7;&#xb7;&#xb7;&#xb7;&#xb7;"
          />
        </div>

        <p v-if="error" class="font-mono text-xs text-destructive">
          {{ error }}
        </p>

        <Button type="submit" class="w-full gap-2">
          <LogIn class="h-4 w-4" /> Sign In
        </Button>

        <div class="relative my-4">
          <div class="absolute inset-0 flex items-center">
            <span class="w-full border-t border-border" />
          </div>
          <div class="relative flex justify-center text-xs uppercase">
            <span class="bg-card px-2 font-mono text-muted-foreground"
              >or</span
            >
          </div>
        </div>

        <Button type="button" variant="ghost" class="w-full" @click="handleGuest">
          Continue as Guest
        </Button>

        <p
          class="text-center font-mono text-xs text-muted-foreground"
        >
          Don't have an account?
          <button
            type="button"
            @click="switchTab('signup')"
            class="text-primary hover:underline"
          >
            Sign up
          </button>
        </p>
      </form>

      <!-- Signup form -->
      <form v-else @submit.prevent="handleSignup" class="mt-6 space-y-4">
        <div>
          <label
            class="font-mono text-xs uppercase tracking-widest text-muted-foreground"
            >Name</label
          >
          <Input
            v-model="signupName"
            class="mt-1.5 font-mono text-sm"
            placeholder="Your name"
          />
        </div>
        <div>
          <label
            class="font-mono text-xs uppercase tracking-widest text-muted-foreground"
            >Email</label
          >
          <Input
            v-model="signupEmail"
            type="email"
            class="mt-1.5 font-mono text-sm"
            placeholder="email@example.com"
          />
        </div>
        <div>
          <label
            class="font-mono text-xs uppercase tracking-widest text-muted-foreground"
            >Password</label
          >
          <Input
            v-model="signupPassword"
            type="password"
            class="mt-1.5 font-mono text-sm"
          />
        </div>
        <div>
          <label
            class="font-mono text-xs uppercase tracking-widest text-muted-foreground"
            >Confirm Password</label
          >
          <Input
            v-model="signupConfirm"
            type="password"
            class="mt-1.5 font-mono text-sm"
          />
        </div>

        <p v-if="error" class="font-mono text-xs text-destructive">
          {{ error }}
        </p>

        <Button type="submit" class="w-full gap-2">
          <UserPlus class="h-4 w-4" /> Create Account
        </Button>

        <div class="relative my-4">
          <div class="absolute inset-0 flex items-center">
            <span class="w-full border-t border-border" />
          </div>
          <div class="relative flex justify-center text-xs uppercase">
            <span class="bg-card px-2 font-mono text-muted-foreground"
              >or</span
            >
          </div>
        </div>

        <Button type="button" variant="ghost" class="w-full" @click="handleGuest">
          Continue as Guest
        </Button>

        <p
          class="text-center font-mono text-xs text-muted-foreground"
        >
          Already have an account?
          <button
            type="button"
            @click="switchTab('login')"
            class="text-primary hover:underline"
          >
            Sign in
          </button>
        </p>
      </form>
    </div>
  </div>
</template>
