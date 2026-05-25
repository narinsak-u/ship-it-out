<script setup lang="ts">
import { ref, computed } from "vue";
import { LogIn, UserPlus, Loader2 } from "lucide-vue-next";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import Separator from "@/components/ui/Separator.vue";
import { useAuthStore } from "@/stores/auth";
import Input from "@/components/ui/Input.vue";
import Button from "@/components/ui/Button.vue";

const props = defineProps<{ open?: boolean }>();
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
const signupConfirm = ref("");

const error = computed(() => store.error);

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
  store.error = "";
  if (
    !signupName.value.trim() ||
    !signupEmail.value.trim() ||
    !signupPassword.value.trim() ||
    !signupConfirm.value.trim()
  ) {
    store.error = "Please fill in all fields";
    return;
  }
  if (signupPassword.value !== signupConfirm.value) {
    store.error = "Passwords do not match";
    return;
  }
  submitting.value = true;
  const err = await store.signup(signupName.value, signupEmail.value, signupPassword.value);
  submitting.value = false;
  if (!err) emit("authenticated");
}

function handleGuest() {
  store.enterGuestMode();
  emit("guest");
}
</script>

<template>
  <Dialog :open="open" @update:open="(v) => !v && emit('close')">
    <DialogContent class="w-full max-w-md rounded-xl border border-border bg-card p-6 shadow-elegant sm:max-w-md">
      <DialogHeader>
        <DialogTitle class="font-mono text-lg font-semibold">
          {{ activeTab === "login" ? "Sign in" : "Create account" }}
        </DialogTitle>
      </DialogHeader>

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
      <form v-if="activeTab === 'login'" @submit.prevent="handleLogin" class="mt-6 space-y-4">
        <div>
          <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground"
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
          <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground"
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

        <Button type="submit" class="w-full gap-2" :disabled="submitting">
          <LogIn v-if="!submitting" class="h-4 w-4" />
          <Loader2 v-else class="h-4 w-4 animate-spin" />
          {{ submitting ? "Signing in..." : "Sign In" }}
        </Button>

        <div class="relative my-4">
          <div class="absolute inset-0 flex items-center">
            <Separator />
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
          <button type="button" @click="switchTab('signup')" class="text-primary hover:underline">
            Sign up
          </button>
        </p>
      </form>

      <!-- Signup form -->
      <form v-else @submit.prevent="handleSignup" class="mt-6 space-y-4">
        <div>
          <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground"
            >Name</label
          >
          <Input v-model="signupName" class="mt-1.5 font-mono text-sm" placeholder="Your name" />
        </div>
        <div>
          <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground"
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
          <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground"
            >Password</label
          >
          <Input v-model="signupPassword" type="password" class="mt-1.5 font-mono text-sm" />
        </div>
        <div>
          <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground"
            >Confirm Password</label
          >
          <Input v-model="signupConfirm" type="password" class="mt-1.5 font-mono text-sm" />
        </div>

        <p v-if="error" class="font-mono text-xs text-destructive">
          {{ error }}
        </p>

        <Button type="submit" class="w-full gap-2" :disabled="submitting">
          <UserPlus v-if="!submitting" class="h-4 w-4" />
          <Loader2 v-else class="h-4 w-4 animate-spin" />
          {{ submitting ? "Creating account..." : "Create Account" }}
        </Button>

        <div class="relative my-4">
          <div class="absolute inset-0 flex items-center">
            <Separator />
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
          <button type="button" @click="switchTab('login')" class="text-primary hover:underline">
            Sign in
          </button>
        </p>
      </form>
    </DialogContent>
  </Dialog>
</template>
