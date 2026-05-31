import { describe, it, expect, beforeEach, vi } from "vitest";
import { setActivePinia, createPinia } from "pinia";
import { useAuthStore } from "@/stores/auth";

const mockUser = {
  id: 1,
  name: "Test User",
  email: "test@test.com",
  role: "admin",
  created_at: "2026-01-01T00:00:00Z",
};

vi.mock("@/lib/api/client", () => {
  const mockApi = {
    get: vi.fn(),
    post: vi.fn(),
  };
  return { api: mockApi };
});

import { api } from "@/lib/api/client";

describe("auth store", () => {
  beforeEach(() => {
    setActivePinia(createPinia());
    vi.clearAllMocks();
    sessionStorage.clear();
  });

  it("init() fetches user and sets state on success", async () => {
    vi.mocked(api.get).mockResolvedValue({ data: mockUser });
    const store = useAuthStore();
    await store.init();
    expect(store.user).toEqual(mockUser);
    expect(store.loading).toBe(false);
  });

  it("init() handles API failure gracefully", async () => {
    vi.mocked(api.get).mockRejectedValue(new Error("Network error"));
    const store = useAuthStore();
    await store.init();
    expect(store.user).toBeNull();
    expect(store.loading).toBe(false);
  });

  it("init() with guest mode skips API call", async () => {
    sessionStorage.setItem("harborops_guest", "true");
    const store = useAuthStore();
    store.enterGuestMode();
    expect(api.get).not.toHaveBeenCalled();
    expect(store.loading).toBe(false);
    expect(store.isGuest).toBe(true);
  });

  it("login() sets user on success", async () => {
    vi.mocked(api.post).mockResolvedValue({ data: { user: mockUser } });
    vi.mocked(api.get).mockResolvedValue({ data: mockUser });
    const store = useAuthStore();
    const err = await store.login("test@test.com", "password");
    expect(err).toBeNull();
    expect(store.user).toEqual(mockUser);
  });

  it("login() returns error message on failure", async () => {
    vi.mocked(api.post).mockResolvedValue({ error: "Invalid credentials" });
    const store = useAuthStore();
    const err = await store.login("bad@test.com", "wrong");
    expect(err).toBe("Invalid credentials");
    expect(store.user).toBeNull();
  });

  it("signup() sets user on success", async () => {
    vi.mocked(api.post).mockResolvedValue({ data: { user: mockUser } });
    vi.mocked(api.get).mockResolvedValue({ data: mockUser });
    const store = useAuthStore();
    const err = await store.signup("Test User", "test@test.com", "password");
    expect(err).toBeNull();
    expect(store.user).toEqual(mockUser);
  });

  it("signup() returns error on failure", async () => {
    vi.mocked(api.post).mockResolvedValue({ error: "Email taken" });
    const store = useAuthStore();
    const err = await store.signup("Test", "exists@test.com", "password");
    expect(err).toBe("Email taken");
  });

  it("logout() clears user and sessionStorage", async () => {
    vi.mocked(api.post).mockResolvedValue({ data: { success: true } });
    sessionStorage.setItem("harborops_guest", "true");
    const store = useAuthStore();
    store.user = mockUser;
    store.isGuest = false;
    await store.logout();
    expect(store.user).toBeNull();
    expect(store.isGuest).toBe(false);
    expect(sessionStorage.getItem("harborops_guest")).toBeNull();
  });

  it("enterGuestMode() sets sessionStorage flag", () => {
    const store = useAuthStore();
    store.enterGuestMode();
    expect(store.isGuest).toBe(true);
    expect(sessionStorage.getItem("harborops_guest")).toBe("true");
  });

  it("isAuthenticated is true when user is set", () => {
    const store = useAuthStore();
    expect(store.isAuthenticated).toBe(false);
    store.user = mockUser;
    expect(store.isAuthenticated).toBe(true);
  });
});
