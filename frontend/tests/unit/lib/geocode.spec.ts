import { describe, it, expect, vi, beforeEach } from "vitest";

const mockGeocode = vi.hoisted(() => vi.fn());

vi.mock("opencage-api-client", () => ({
  geocode: mockGeocode,
}));

import { geocodeAddress } from "@/lib/geocode";

describe("geocodeAddress", () => {
  beforeEach(() => {
    vi.stubEnv("VITE_OPENCAGE_API_KEY", "test-key");
    mockGeocode.mockReset();
  });

  it("returns coordinates on successful lookup", async () => {
    mockGeocode.mockResolvedValue({
      results: [{ geometry: { lat: 13.75, lng: 100.5 } }],
    });
    const result = await geocodeAddress("Phra Nakhon", "Bangkok", "Bangkok");
    expect(result).toEqual({ lat: 13.75, lng: 100.5 });
  });

  it("throws when API key is not configured", async () => {
    vi.stubEnv("VITE_OPENCAGE_API_KEY", "");
    await expect(geocodeAddress("A", "B", "C")).rejects.toThrow(
      "OpenCage API key is not configured.",
    );
    expect(mockGeocode).not.toHaveBeenCalled();
  });

  it("throws when geocode request fails", async () => {
    mockGeocode.mockRejectedValue(new Error("Network error"));
    await expect(geocodeAddress("A", "B", "C")).rejects.toThrow(
      "Location lookup failed. Please try again later.",
    );
  });

  it("throws when no results returned", async () => {
    mockGeocode.mockResolvedValue({ results: [] });
    await expect(geocodeAddress("A", "B", "C")).rejects.toThrow(
      "Could not resolve this address. Check the fields and try again.",
    );
  });

  it("throws when results lack geometry", async () => {
    mockGeocode.mockResolvedValue({ results: [{}] });
    await expect(geocodeAddress("A", "B", "C")).rejects.toThrow(
      "Could not resolve this address. Check the fields and try again.",
    );
  });
});
