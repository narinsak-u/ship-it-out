import { NextRequest } from "next/server";

export async function POST(request: NextRequest) {
  const base = process.env.API_BASE_URL;
  if (!base) return Response.json({ error: "API_BASE_URL is required" }, { status: 500 });
  let profile;
  try {
    profile = await fetch(`${base.replace(/\/$/, "")}/api/auth/me`, {
      headers: { cookie: request.headers.get("cookie") ?? "" },
      cache: "no-store",
    });
  } catch {
    return Response.json({ error: "Authentication service unavailable." }, { status: 502 });
  }
  if (profile.status === 401)
    return Response.json({ error: "Not authenticated." }, { status: 401 });
  if (!profile.ok)
    return Response.json({ error: "Authentication service unavailable." }, { status: 502 });
  const user = (await profile.json()).data;
  if (user?.role !== "admin" && user?.role !== "staff") {
    return Response.json({ error: "Staff role required." }, { status: 403 });
  }
  const key = process.env.OPENCAGE_API_KEY;
  if (!key) return Response.json({ error: "OpenCage API key is not configured." }, { status: 503 });

  let address: { subDistrict?: string; district?: string; province?: string };
  try {
    address = await request.json();
    if (!address || typeof address !== "object") throw new Error("Invalid address");
  } catch {
    return Response.json({ error: "Invalid address." }, { status: 400 });
  }
  const query = [address.subDistrict, address.district, address.province]
    .filter((part): part is string => typeof part === "string" && !!part.trim())
    .join(", ");
  if (!query) return Response.json({ error: "Address is required." }, { status: 400 });

  try {
    const url = new URL("https://api.opencagedata.com/geocode/v1/json");
    url.search = new URLSearchParams({ q: query, key }).toString();
    const response = await fetch(url, { cache: "no-store" });
    if (!response.ok) throw new Error("lookup failed");
    const result = await response.json();
    const geometry = result.results?.[0]?.geometry;
    if (typeof geometry?.lat !== "number" || typeof geometry?.lng !== "number") {
      return Response.json(
        { error: "Could not resolve this address. Check the fields and try again." },
        { status: 422 },
      );
    }
    return Response.json({ data: { lat: geometry.lat, lng: geometry.lng } });
  } catch {
    return Response.json(
      { error: "Location lookup failed. Please try again later." },
      { status: 502 },
    );
  }
}
