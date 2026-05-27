import { geocode } from "opencage-api-client";

export async function geocodeAddress(
  subDistrict: string,
  district: string,
  province: string,
): Promise<{ lat: number; lng: number }> {
  const q = `${subDistrict}, ${district}, ${province}`.trim();
  const key = import.meta.env.VITE_OPENCAGE_API_KEY as string;

  if (!key) {
    throw new Error("OpenCage API key is not configured.");
  }

  let data;
  try {
    data = await geocode({ q, key });
  } catch {
    throw new Error("Location lookup failed. Please try again later.");
  }

  if (!data.results || data.results.length === 0) {
    throw new Error("Could not resolve this address. Check the fields and try again.");
  }

  if (!data.results[0].geometry) {
    throw new Error("Could not resolve this address. Check the fields and try again.");
  }
  return data.results[0].geometry;
}
