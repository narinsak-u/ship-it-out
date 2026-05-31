import { ref, computed } from "vue";
import { toast } from "vue-sonner";
import { useAuthStore } from "@/stores/auth";

export function useAuthForm(mode: "login" | "signup") {
  const store = useAuthStore();
  const submitting = ref(false);

  const error = computed(() => store.error);

  async function handleLogin(email: string, password: string): Promise<boolean> {
    submitting.value = true;
    store.error = "";
    const err = await store.login(email, password);
    submitting.value = false;
    if (!err) {
      toast.success("Signed in successfully");
      return true;
    }
    return false;
  }

  async function handleSignup(
    name: string,
    email: string,
    password: string,
    confirm: string,
  ): Promise<boolean> {
    store.error = "";
    if (!name.trim() || !email.trim() || !password.trim() || !confirm.trim()) {
      store.error = "Please fill in all fields";
      return false;
    }
    if (password !== confirm) {
      store.error = "Passwords do not match";
      return false;
    }
    submitting.value = true;
    const err = await store.signup(name, email, password);
    submitting.value = false;
    if (!err) {
      toast.success("Account created");
      return true;
    }
    return false;
  }

  function handleGuest() {
    store.enterGuestMode();
  }

  return { submitting, error, handleLogin, handleSignup, handleGuest };
}
