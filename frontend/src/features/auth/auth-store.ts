"use client";

import { create } from "zustand";

export type User = {
  id: number;
  name: string;
  email: string;
  role: string;
  created_at?: string;
};

type AuthState = {
  user: User | null;
  isStaff: boolean;
  loading: boolean;
  initialize: () => Promise<void>;
  login: (email: string, password: string) => Promise<void>;
  signup: (name: string, email: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
};

async function authRequest<T>(path: string, method: string, body?: object): Promise<T> {
  const response = await fetch(`/api/auth/${path}`, {
    method,
    credentials: "same-origin",
    headers: body ? { "content-type": "application/json" } : undefined,
    body: body ? JSON.stringify(body) : undefined,
    cache: "no-store",
  });
  const result = await response.json();
  if (!response.ok) throw new Error(result.error || `Request failed (${response.status})`);
  return result.data as T;
}

let sessionInitialization: Promise<void> | undefined;

export const useAuthStore = create<AuthState>((set) => {
  const setUser = (user: User | null) =>
    set({ user, isStaff: user?.role === "admin" || user?.role === "staff" });

  return {
    user: null,
    isStaff: false,
    loading: true,
    initialize() {
      sessionInitialization ??= authRequest<User>("me", "GET")
        .then(setUser)
        .catch(() => setUser(null))
        .finally(() => set({ loading: false }));
      return sessionInitialization;
    },
    async login(email, password) {
      const result = await authRequest<{ user: User }>("login", "POST", { email, password });
      setUser(result.user);
    },
    async signup(name, email, password) {
      const result = await authRequest<{ user: User }>("register", "POST", {
        name,
        email,
        password,
      });
      setUser(result.user);
    },
    async logout() {
      await authRequest<unknown>("logout", "POST");
      setUser(null);
    },
  };
});
