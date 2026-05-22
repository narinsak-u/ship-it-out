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
