<script setup lang="ts">
import { ref } from "vue";
import { UserPlus, Loader2 } from "lucide-vue-next";
import Input from "@/components/ui/Input.vue";
import Button from "@/components/ui/Button.vue";
import Separator from "@/components/ui/Separator.vue";
import { useAuthForm } from "@/composables/useAuthForm";

const emit = defineEmits<{ success: []; guest: []; "switch-to-login": [] }>();

const { submitting, error, handleSignup, handleGuest } = useAuthForm("signup");

const name = ref("");
const email = ref("");
const password = ref("");
const confirm = ref("");

async function onSubmit() {
  const ok = await handleSignup(name.value, email.value, password.value, confirm.value);
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
      <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground">Name</label>
      <Input v-model="name" class="mt-1.5 font-mono text-sm" placeholder="Your name" />
    </div>
    <div>
      <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground">Email</label>
      <Input
        v-model="email"
        type="email"
        class="mt-1.5 font-mono text-sm"
        placeholder="email@example.com"
      />
    </div>
    <div>
      <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground"
        >Password</label
      >
      <Input v-model="password" type="password" class="mt-1.5 font-mono text-sm" />
    </div>
    <div>
      <label class="font-mono text-xs uppercase tracking-widest text-muted-foreground"
        >Confirm Password</label
      >
      <Input v-model="confirm" type="password" class="mt-1.5 font-mono text-sm" />
    </div>

    <p v-if="error" class="font-mono text-xs text-destructive">{{ error }}</p>

    <Button type="submit" class="w-full gap-2" :disabled="submitting">
      <UserPlus v-if="!submitting" class="h-4 w-4" />
      <Loader2 v-else class="h-4 w-4 animate-spin" />
      {{ submitting ? "Creating account..." : "Create Account" }}
    </Button>

    <div class="relative my-4">
      <div class="absolute inset-0 flex items-center"><Separator /></div>
      <div class="relative flex justify-center text-xs uppercase">
        <span class="bg-card px-2 font-mono text-muted-foreground">or</span>
      </div>
    </div>

    <Button type="button" variant="ghost" class="w-full" @click="onGuest">Continue as Guest</Button>

    <p class="text-center font-mono text-xs text-muted-foreground">
      Already have an account?
      <button type="button" @click="emit('switch-to-login')" class="text-primary hover:underline">
        Sign in
      </button>
    </p>
  </form>
</template>
