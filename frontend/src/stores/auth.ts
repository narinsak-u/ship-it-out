import { defineStore } from "pinia";
import { ref, computed } from "vue";
import { api } from "@/lib/api/client";

export interface AuthUser {
  id: number;
  name: string;
  email: string;
  role: string;
  created_at: string;
}

export const useAuthStore = defineStore("auth", () => {
  const user = ref<AuthUser | null>(null);
  const loading = ref(true);
  const error = ref("");
  const isStaff = computed(() => user.value?.role === "admin" || user.value?.role === "staff");

  const isAuthenticated = computed(() => user.value !== null);

  async function init() {
    loading.value = true;
    try {
      const res = await api.get<AuthUser>("/auth/me");
      if (res.data) user.value = res.data;
    } catch {
      // API failure — user stays null
    }
    loading.value = false;
  }

  async function login(email: string, password: string): Promise<string | null> {
    error.value = "";
    const res = await api.post<{ user: AuthUser }>("/auth/login", { email, password });
    if (res.error) {
      error.value = res.error;
      return res.error;
    }
    const me = await api.get<AuthUser>("/auth/me");
    if (me.data) {
      user.value = me.data;
      return null;
    }
    error.value = "Login failed -- session not established";
    return error.value;
  }

  async function signup(name: string, email: string, password: string): Promise<string | null> {
    error.value = "";
    const res = await api.post<{ user: AuthUser }>("/auth/register", { name, email, password });
    if (res.error) {
      error.value = res.error;
      return res.error;
    }
    const me = await api.get<AuthUser>("/auth/me");
    if (me.data) {
      user.value = me.data;
      return null;
    }
    error.value = "Signup failed -- session not established";
    return error.value;
  }

  async function logout() {
    await api.post("/auth/logout");
    user.value = null;
    loading.value = false;
  }

  return {
    user,
    loading,
    error,
    isStaff,
    isAuthenticated,
    init,
    login,
    signup,
    logout,
  };
});
