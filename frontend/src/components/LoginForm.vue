<script setup lang="ts">
import { ref } from "vue";
import { LogIn, Loader2 } from "lucide-vue-next";
import Input from "@/components/ui/Input.vue";
import Button from "@/components/ui/Button.vue";
import { useAuthForm } from "@/composables/useAuthForm";

const emit = defineEmits<{ success: []; "switch-to-signup": [] }>();

const { submitting, error, handleLogin } = useAuthForm("login");

const email = ref("demo@shipitout.example");
const password = ref("demo-read-only-2026");

async function onSubmit() {
  const ok = await handleLogin(email.value, password.value);
  if (ok) emit("success");
}
</script>

<template>
  <form @submit.prevent="onSubmit" class="mt-6 space-y-4">
    <div>
      <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground">Email</label>
      <Input
        v-model="email"
        type="email"
        class="mt-1.5 font-mono text-sm"
        placeholder="demo@shipitout.example"
      />
    </div>
    <div>
      <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground"
        >Password</label
      >
      <Input
        v-model="password"
        type="password"
        class="mt-1.5 font-mono text-sm"
        placeholder="&#xb7;&#xb7;&#xb7;&#xb7;&#xb7;&#xb7;&#xb7;&#xb7;"
      />
    </div>

    <p v-if="error" class="font-mono text-xs text-destructive">{{ error }}</p>

    <Button type="submit" class="w-full gap-2" :disabled="submitting">
      <LogIn v-if="!submitting" class="h-4 w-4" />
      <Loader2 v-else class="h-4 w-4 animate-spin" />
      {{ submitting ? "Signing in..." : "Sign In" }}
    </Button>
    <p
      class="rounded-md border border-border bg-secondary/40 p-3 font-mono text-xs text-muted-foreground"
    >
      Read-only demo account. Email: demo@shipitout.example · Password: demo-read-only-2026
    </p>

    <p class="text-center font-mono text-xs text-muted-foreground">
      Don't have an account?
      <button type="button" @click="emit('switch-to-signup')" class="text-primary hover:underline">
        Sign up
      </button>
    </p>
  </form>
</template>
