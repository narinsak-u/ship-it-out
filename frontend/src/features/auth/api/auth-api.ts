import type { User } from "../types";

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

export const getSession = () => authRequest<User>("me", "GET");
export const login = (email: string, password: string) =>
  authRequest<{ user: User }>("login", "POST", { email, password });
export const register = (name: string, email: string, password: string) =>
  authRequest<{ user: User }>("register", "POST", { name, email, password });
export const logout = () => authRequest<unknown>("logout", "POST");
