import { describe, it, expect } from "vitest";
import { http, HttpResponse } from "msw";
import { server } from "../../../tests/msw/server";
import { api, BASE_URL } from "./client";

describe("api client", () => {
  it("GET returns data on success", async () => {
    const result = await api.get("/auth/me");
    expect(result.data).toBeDefined();
    expect(result.error).toBeUndefined();
  });
  it("POST sends body and returns data", async () => {
    const result = await api.post("/auth/login", { email: "test@test.com", password: "pass" });
    expect(result.data).toBeDefined();
  });
  it("returns error object on 4xx response", async () => {
    server.use(
      http.get(`${BASE_URL}/auth/me`, () =>
        HttpResponse.json({ error: "Unauthorized" }, { status: 401 }),
      ),
    );
    const result = await api.get("/auth/me");
    expect(result.error).toBe("Unauthorized");
    expect(result.data).toBeUndefined();
  });
  it("returns network error when fetch fails", async () => {
    server.use(http.get(`${BASE_URL}/auth/me`, () => HttpResponse.error()));
    const result = await api.get("/auth/me");
    expect(result.error).toContain("Network error");
  });
  it("DELETE works", async () => {
    const result = await api.del("/shipments/ORD-001");
    expect(result.data).toBeDefined();
  });
  it("PUT works", async () => {
    const result = await api.put("/shipments/ORD-001", { weight: 15 });
    expect(result.data).toBeDefined();
  });
});
