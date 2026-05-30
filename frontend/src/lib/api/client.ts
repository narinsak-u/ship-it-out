export const BASE_URL = import.meta.env.VITE_API_URL || "http://localhost:8080/api";

interface ApiSuccess<T> {
  data: T;
  error?: never;
}

interface ApiError {
  error: string;
  data?: never;
}

type ApiResult<T> = ApiSuccess<T> | ApiError;

async function request<T = unknown>(
  path: string,
  options?: RequestInit,
  raw = false,
): Promise<ApiResult<T>> {
  try {
    const res = await fetch(`${BASE_URL}${path}`, {
      credentials: "include",
      headers: { "Content-Type": "application/json", ...options?.headers },
      ...options,
    });
    const json = await res.json();
    if (!res.ok) return { error: json.error || `Request failed (${res.status})` };

    return { data: raw ? (json as T) : (json.data as T) };
  } catch {
    return { error: "Network error -- is the backend running?" };
  }
}

export const api = {
  get: <T = unknown>(path: string) => request<T>(path),
  getRaw: <T = unknown>(path: string) => request<T>(path, undefined, true),
  post: <T = unknown>(path: string, body?: unknown) =>
    request<T>(path, { method: "POST", body: body ? JSON.stringify(body) : undefined }),
  del: <T = unknown>(path: string) => request<T>(path, { method: "DELETE" }),
  put: <T = unknown>(path: string, body?: unknown) =>
    request<T>(path, { method: "PUT", body: body ? JSON.stringify(body) : undefined }),
  patch: <T = unknown>(path: string, body?: unknown) =>
    request<T>(path, { method: "PATCH", body: body ? JSON.stringify(body) : undefined }),
};
