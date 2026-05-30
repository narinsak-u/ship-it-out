<script setup lang="ts">
import { ref } from "vue";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import LoginForm from "@/components/LoginForm.vue";
import SignupForm from "@/components/SignupForm.vue";

defineProps<{ open?: boolean }>();
const emit = defineEmits<{ close: []; authenticated: []; guest: [] }>();

type Tab = "login" | "signup";
const activeTab = ref<Tab>("login");

function switchTab(tab: Tab) {
  activeTab.value = tab;
}

function onAuthenticated() {
  emit("authenticated");
}

function onGuest() {
  emit("guest");
}
</script>

<template>
  <Dialog :open="open" @update:open="(v) => !v && emit('close')">
    <DialogContent
      class="w-full max-w-md rounded-xl border border-border bg-card p-6 shadow-elegant sm:max-w-md"
    >
      <DialogHeader>
        <DialogTitle class="font-mono text-lg font-semibold">
          {{ activeTab === "login" ? "Sign in" : "Create account" }}
        </DialogTitle>
      </DialogHeader>

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

      <LoginForm
        v-if="activeTab === 'login'"
        @success="onAuthenticated"
        @guest="onGuest"
        @switch-to-signup="switchTab('signup')"
      />
      <SignupForm
        v-else
        @success="onAuthenticated"
        @guest="onGuest"
        @switch-to-login="switchTab('login')"
      />
    </DialogContent>
  </Dialog>
</template>
