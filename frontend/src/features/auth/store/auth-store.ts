"use client";

import { create } from "zustand";
import { getSession, login, logout, register } from "../api/auth-api";
import type { User } from "../types";

type AuthState = {
  user: User | null;
  isStaff: boolean;
  loading: boolean;
  initialize: () => Promise<void>;
  login: (email: string, password: string) => Promise<void>;
  signup: (name: string, email: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
};

let sessionInitialization: Promise<void> | undefined;

export const useAuthStore = create<AuthState>((set) => {
  const setUser = (user: User | null) =>
    set({ user, isStaff: user?.role === "admin" || user?.role === "staff" });

  return {
    user: null,
    isStaff: false,
    loading: true,
    initialize() {
      sessionInitialization ??= getSession()
        .then(setUser)
        .catch(() => setUser(null))
        .finally(() => set({ loading: false }));
      return sessionInitialization;
    },

    async login(email, password) {
      const result = await login(email, password);
      setUser(result.user);
    },
    async signup(name, email, password) {
      const result = await register(name, email, password);
      setUser(result.user);
    },
    async logout() {
      await logout();
      setUser(null);
    },
  };
});
