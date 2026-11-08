import "server-only";

export function getApiOrigin() {
  const origin = process.env.API_BASE_URL;
  if (!origin) throw new Error("API_BASE_URL is required");
  return origin.replace(/\/$/, "");
}
