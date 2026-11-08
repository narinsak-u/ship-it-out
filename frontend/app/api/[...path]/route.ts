import { NextRequest } from "next/server";

async function proxy(request: NextRequest, context: { params: Promise<{ path: string[] }> }) {
  const { path } = await context.params;
  const base = process.env.API_BASE_URL;
  if (!base) return Response.json({ error: "API_BASE_URL is required" }, { status: 500 });
  const target = new URL(
    `/api/${path.map(encodeURIComponent).join("/")}${request.nextUrl.search}`,
    base,
  );
  const headers = new Headers(request.headers);
  headers.delete("host");
  const response = await fetch(target, {
    method: request.method,
    headers,
    body: ["GET", "HEAD"].includes(request.method) ? undefined : await request.arrayBuffer(),
    cache: "no-store",
    redirect: "manual",
  });
  return new Response(response.body, { status: response.status, headers: response.headers });
}

export const GET = proxy;
export const POST = proxy;
export const PUT = proxy;
export const PATCH = proxy;
export const DELETE = proxy;
