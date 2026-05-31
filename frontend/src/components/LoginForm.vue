<script setup lang="ts">
import { ref } from "vue";
import { LogIn, Loader2 } from "lucide-vue-next";
import Input from "@/components/ui/Input.vue";
import Button from "@/components/ui/Button.vue";
import Separator from "@/components/ui/Separator.vue";
import { useAuthForm } from "@/composables/useAuthForm";

const emit = defineEmits<{ success: []; guest: []; "switch-to-signup": [] }>();

const { submitting, error, handleLogin, handleGuest } = useAuthForm("login");

const email = ref("admin@gmail.com");
const password = ref("123456789");

async function onSubmit() {
  const ok = await handleLogin(email.value, password.value);
  if (ok) emit("success");
}

function onGuest() {
  handleGuest();
  emit("guest");
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
        placeholder="admin@harborops.io"
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

    <div class="relative my-4">
      <div class="absolute inset-0 flex items-center"><Separator /></div>
      <div class="relative flex justify-center text-xs uppercase">
        <span class="bg-card px-2 font-mono text-muted-foreground">or</span>
      </div>
    </div>

    <Button type="button" variant="ghost" class="w-full" @click="onGuest">Continue as Guest</Button>

    <p class="text-center font-mono text-xs text-muted-foreground">
      Don't have an account?
      <button type="button" @click="emit('switch-to-signup')" class="text-primary hover:underline">
        Sign up
      </button>
    </p>
  </form>
</template>
