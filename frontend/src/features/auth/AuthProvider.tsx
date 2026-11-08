"use client";

import { createContext, useContext, useEffect, useMemo, useState } from "react";

type User = {
  id: number;
  name: string;
  email: string;
  role: string;
  created_at?: string;
};

type AuthContextValue = {
  user: User | null;
  isStaff: boolean;
  loading: boolean;
  login: (email: string, password: string) => Promise<void>;
  signup: (name: string, email: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
};

const AuthContext = createContext<AuthContextValue | null>(null);

async function authRequest(path: string, method: string, body?: object) {
  const response = await fetch(`/api/auth/${path}`, {
    method,
    credentials: "same-origin",
    headers: body ? { "content-type": "application/json" } : undefined,
    body: body ? JSON.stringify(body) : undefined,
    cache: "no-store",
  });
  const result = await response.json();
  if (!response.ok) throw new Error(result.error || `Request failed (${response.status})`);
  return result.data;
}

export function AuthProvider({ children }: Readonly<{ children: React.ReactNode }>) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    authRequest("me", "GET")
      .then(setUser)
      .catch(() => setUser(null))
      .finally(() => setLoading(false));
  }, []);

  const value = useMemo<AuthContextValue>(
    () => ({
      user,
      isStaff: user?.role === "admin" || user?.role === "staff",
      loading,
      async login(email, password) {
        const result = await authRequest("login", "POST", { email, password });
        setUser(result.user);
      },
      async signup(name, email, password) {
        const result = await authRequest("register", "POST", { name, email, password });
        setUser(result.user);
      },
      async logout() {
        await authRequest("logout", "POST");
        setUser(null);
      },
    }),
    [user, loading],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) throw new Error("useAuth must be used inside AuthProvider");
  return context;
}
