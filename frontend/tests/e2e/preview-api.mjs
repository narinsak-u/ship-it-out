import { createServer } from "node:http";

let shipment = {
  id: "ORD-TEST1",
  trackingNumber: "TH-TEST-1",
  status: "in_transit",
  carrier: "Flash Express",
  customer: { name: "Sender", zipcode: "10900", subDistrict: "ลาดยาว", district: "จตุจักร", province: "กรุงเทพมหานคร", coords: { lat: 13.7, lng: 100.5 } },
  receiver: { name: "Receiver", zipcode: "10800", subDistrict: "วงศ์สว่าง", district: "บางซื่อ", province: "กรุงเทพมหานคร", coords: { lat: 13.3, lng: 100.9 } },
  origin: "Bangkok", destination: "Chonburi", currentCoords: { lat: 13.5, lng: 100.7 }, weight: 2, items: 2,
  estimatedDelivery: "2026-12-01T00:00:00Z", createdAt: "2026-11-01T00:00:00Z", progress: 55,
};
const originalShipment = structuredClone(shipment);
const activeDeliveries = [shipment, ...Array.from({ length: 10 }, (_, index) => ({ ...shipment, id: `ORD-TEST${index + 2}`, trackingNumber: `TH-TEST-${index + 2}`, customer: { ...shipment.customer, name: `Sender ${index + 2}` } }))];
let deliveryReads = 0;
const manifestStatuses = ["in_transit", "delivered", "pending", "in_transit", "delivered", "pending"];
const createManifest = () => Array.from({ length: 12 }, (_, index) => {
  const number = String(index + 1).padStart(2, "0");
  return { ...structuredClone(shipment), id: `ORD-MANIFEST-${number}`, trackingNumber: `TH-MANIFEST-${number}`, status: manifestStatuses[index % manifestStatuses.length], customer: { ...shipment.customer, name: `Sender ${number}` }, receiver: { ...shipment.receiver, name: `Receiver ${number}` } };
});
let manifestRecords = createManifest();
let analyticsMode = "normal";
const authUsers = new Map([
  ["demo@shipitout.example", { id: 2, name: "Ship It Out Demo", email: "demo@shipitout.example", role: "customer" }],
  ["staff@shipitout.example", { id: 1, name: "Preview Staff", email: "staff@shipitout.example", role: "staff" }],
]);
const authPasswords = new Map([["demo@shipitout.example", "demo-read-only-2026"], ["staff@shipitout.example", "preview-staff-password"]]);
const authSessions = new Map();
let shipmentUpdateFailure = false;
let nextUserId = 3;
let manifestDeleteFailure = false;
let nextShipmentNumber = 1;

createServer(async (req, res) => {
  res.setHeader("content-type", "application/json");
  const cookies = Object.fromEntries((req.headers.cookie || "").split("; ").filter(Boolean).map((item) => item.split("=")));
  const user = authSessions.get(cookies.jwt);
  if (req.url === "/__test/delivery-reads") return res.end(JSON.stringify({ reads: deliveryReads }));
  if (req.method === "POST" && req.url === "/__test/deliveries?reset=1") { shipment = structuredClone(originalShipment); activeDeliveries[0] = shipment; deliveryReads = 0; res.statusCode = 204; return res.end(); }
  if (req.method === "POST" && req.url === "/__test/shipment?reset=1") { shipment = structuredClone(originalShipment); activeDeliveries[0] = shipment; manifestRecords = createManifest(); shipmentUpdateFailure = false; res.statusCode = 204; return res.end(); }
  if (req.method === "POST" && req.url === "/__test/shipment?update=failure") { shipmentUpdateFailure = true; res.statusCode = 204; return res.end(); }
  if (req.method === "POST" && req.url === "/__test/manifest?reset=1") { manifestRecords = createManifest(); manifestDeleteFailure = false; res.statusCode = 204; return res.end(); }
  if (req.method === "POST" && req.url === "/__test/manifest?delete=failure") { manifestDeleteFailure = true; res.statusCode = 204; return res.end(); }
  if (req.method === "DELETE" && new URL(req.url ?? "/", "http://localhost").pathname.startsWith("/api/shipments/")) {
    if (!user) { res.statusCode = 401; return res.end(JSON.stringify({ success: false, error: "unauthorized" })); }
    if (user.role !== "staff") { res.statusCode = 403; return res.end(JSON.stringify({ success: false, error: "forbidden" })); }
    if (manifestDeleteFailure) { res.statusCode = 503; return res.end(JSON.stringify({ success: false, error: "shipment delete unavailable" })); }
    const id = decodeURIComponent(new URL(req.url, "http://localhost").pathname.slice("/api/shipments/".length));
    manifestRecords = manifestRecords.filter((record) => record.id !== id);
    return res.end(JSON.stringify({ success: true, data: null }));
  }
  if (req.method === "POST" && req.url?.startsWith("/__test/analytics?mode=")) { analyticsMode = new URL(req.url, "http://localhost").searchParams.get("mode") ?? "normal"; res.statusCode = 204; return res.end(); }
  if (analyticsMode === "failure" && req.url?.startsWith("/api/analytics/")) { res.statusCode = 503; return res.end(JSON.stringify({ success: false, error: "unavailable" })); }
  if (analyticsMode === "empty" && req.url === "/api/analytics/overview") return res.end(JSON.stringify({ success: true, data: { total: 0, active: 0, delivered: 0, by_status: null, by_region: ["Central", "East", "North", "West", "North-east", "South"].map((name) => ({ name, total: 0 })) } }));
  if (analyticsMode === "empty" && req.url === "/api/analytics/timeseries") return res.end(JSON.stringify({ success: true, data: { by_month: null, by_day_of_week: null } }));
  if (req.url === "/api/analytics/overview") return res.end(JSON.stringify({ success: true, data: { total: 126, active: 44, delivered: 82, by_status: [{ status: "delivered", count: 82 }, { status: "in_transit", count: 28 }, { status: "pending", count: 16 }], by_region: [{ name: "Bangkok", total: 64 }, { name: "Chonburi", total: 38 }, { name: "Chiang Mai", total: 24 }, { name: "Phuket", total: 0 }] } }));
  if (req.url === "/api/analytics/timeseries") return res.end(JSON.stringify({ success: true, data: { by_month: [{ month: "2026-06", count: 18 }, { month: "2026-07", count: 24 }, { month: "2026-08", count: 31 }, { month: "2026-09", count: 26 }, { month: "2026-10", count: 27 }], by_day_of_week: ["Sunday", "Thursday", "Monday", "Saturday", "Wednesday", "Tuesday", "Friday"].map((day, index) => ({ day, count: [8, 16, 21, 13, 26, 18, 24][index] })) } }));
  if (req.url === "/api/auth/me") { if (!user) { res.statusCode = 401; return res.end(JSON.stringify({ success: false, error: "not authenticated" })); } return res.end(JSON.stringify({ success: true, data: user })); }
  if (req.url === "/api/auth/login" || req.url === "/api/auth/register") {
    const chunks = []; for await (const chunk of req) chunks.push(chunk); const body = JSON.parse(Buffer.concat(chunks).toString() || "{}"); let account = authUsers.get(body.email);
    if (req.url === "/api/auth/register") { if (account) { res.statusCode = 409; return res.end(JSON.stringify({ success: false, error: "email already registered" })); } account = { id: nextUserId++, name: body.name, email: body.email, role: "customer" }; authUsers.set(body.email, account); authPasswords.set(body.email, body.password); }
    else if (!account || authPasswords.get(body.email) !== body.password) { res.statusCode = 401; return res.end(JSON.stringify({ success: false, error: "invalid email or password" })); }
    const token = `preview-session-${account.id}-${Date.now()}`; authSessions.set(token, account); res.setHeader("set-cookie", `jwt=${token}; Path=/; HttpOnly; SameSite=None; Secure`); return res.end(JSON.stringify({ success: true, data: { user: account, token } }));
  }
  if (req.url === "/api/auth/logout" && req.method === "POST") { authSessions.delete(cookies.jwt); res.setHeader("set-cookie", "jwt=; Path=/; HttpOnly; SameSite=None; Secure; Max-Age=0"); return res.end(JSON.stringify({ success: true, data: { message: "logged out" } })); }
  if (req.method === "POST" && req.url === "/api/shipments") {
    if (!user) { res.statusCode = 401; return res.end(JSON.stringify({ success: false, error: "not authenticated" })); }
    if (user.role !== "admin" && user.role !== "staff") { res.statusCode = 403; return res.end(JSON.stringify({ success: false, error: "staff role required" })); }
    let body = ""; for await (const chunk of req) body += chunk; const input = JSON.parse(body); const number = nextShipmentNumber++;
    const record = { ...structuredClone(shipment), id: `ORD-CREATED-${number}`, trackingNumber: `TH-CREATED-${number}`, status: "pending", customer: input.customer, receiver: input.receiver, carrier: input.carrier, origin: `${input.customer.subDistrict}, ${input.customer.district}, ${input.customer.province}`, destination: `${input.receiver.subDistrict}, ${input.receiver.district}, ${input.receiver.province}`, currentCoords: input.customer.coords, weight: input.weight, items: input.items, createdAt: new Date().toISOString(), progress: 0 };
    manifestRecords.unshift(record); res.statusCode = 201; return res.end(JSON.stringify({ success: true, data: record }));
  }
  if (req.url?.startsWith("/api/shipments?")) {
    if (req.url.includes("limit=-1&exclude_status=delivered")) { deliveryReads++; return res.end(JSON.stringify({ success: true, data: activeDeliveries })); }
    const query = new URL(req.url, "http://localhost").searchParams;
    if (query.has("exclude_status")) return res.end(JSON.stringify({ success: true, data: [shipment], pagination: { page: 1, limit: 3, total: 1, totalPages: 1 } }));
    const search = (query.get("search") ?? "").toLowerCase(); const status = query.get("status");
    const filtered = manifestRecords.filter((record) => `${record.id} ${record.trackingNumber} ${record.customer.name} ${record.receiver.name} ${record.origin} ${record.destination}`.toLowerCase().includes(search) && (!status || record.status === status));
    const page = Number(query.get("page") || 1); const limit = Number(query.get("limit") || 10); const start = (page - 1) * limit;
    return res.end(JSON.stringify({ success: true, data: filtered.slice(start, start + limit), pagination: { page, limit, total: filtered.length, totalPages: Math.max(1, Math.ceil(filtered.length / limit)) } }));
  }
  if (req.url === "/api/hubs") return res.end(JSON.stringify({ success: true, data: Array.from({ length: 12 }, (_, index) => ({ id: `HUB-${String(index + 1).padStart(3, "0")}`, name: ["Bangkok", "Chonburi", "Chiang Mai", "Phuket", "Khon Kaen", "Nakhon Ratchasima", "Hat Yai", "Udon Thani", "Rayong", "Ayutthaya", "Surat Thani", "Lampang"][index] + " Hub", carrierId: "THUN", address: "Thailand", coords: { lat: 13.7563, lng: 100.5018 }, capacity: 1200, currentUtilization: 840, status: index === 1 ? "maintenance" : "active", createdAt: "2026-01-01T00:00:00Z" })) }));
  if (req.method === "PATCH" && req.url === "/api/shipments/ORD-TEST1/status") {
    const session = req.headers.cookie?.match(/(?:^|; )jwt=([^;]+)/)?.[1]; const statusUser = authSessions.get(session);
    if (!statusUser) { res.statusCode = 401; return res.end(JSON.stringify({ error: "unauthorized" })); }
    if (statusUser.role !== "staff" && statusUser.role !== "admin") { res.statusCode = 403; return res.end(JSON.stringify({ error: "forbidden" })); }
    let body = ""; for await (const chunk of req) body += chunk; const update = JSON.parse(body); shipment = { ...shipment, status: update.status, ...(update.hubId === "HUB-001" ? { hubId: "HUB-001", currentCoords: { lat: 13.7563, lng: 100.5018 } } : {}) }; activeDeliveries[0] = shipment; return res.end(JSON.stringify({ success: true, data: shipment }));
  }
  if (req.method === "PUT" && req.url === "/api/shipments/ORD-TEST1") {
    if (!user) { res.statusCode = 401; return res.end(JSON.stringify({ success: false, error: "not authenticated" })); }
    if (user.role !== "staff") { res.statusCode = 403; return res.end(JSON.stringify({ success: false, error: "staff role required" })); }
    if (shipmentUpdateFailure) { res.statusCode = 503; return res.end(JSON.stringify({ success: false, error: "shipment update unavailable" })); }
    let body = ""; for await (const chunk of req) body += chunk; const update = JSON.parse(body); shipment = { ...shipment, ...update, origin: `${update.customer.subDistrict}, ${update.customer.district}, ${update.customer.province}`, destination: `${update.receiver.subDistrict}, ${update.receiver.district}, ${update.receiver.province}`, currentCoords: update.customer.coords }; manifestRecords = manifestRecords.map((record) => record.id === shipment.id ? shipment : record); activeDeliveries[0] = shipment; return res.end(JSON.stringify({ success: true, data: shipment }));
  }
  const createdShipment = manifestRecords.find((record) => req.url === `/api/shipments/${encodeURIComponent(record.id)}`);
  if (createdShipment) return res.end(JSON.stringify({ success: true, data: createdShipment }));
  if (req.url === "/api/shipments/ORD-TEST1") return res.end(JSON.stringify({ success: true, data: shipment }));
  if (req.url === "/api/shipments/missing") { res.statusCode = 404; return res.end(JSON.stringify({ success: false, error: "shipment not found" })); }
  if (req.url === "/api/track/TH-TEST-1") return res.end(JSON.stringify({ success: true, data: { shipment, events: [{ status: "In Transit", location: { name: "Bangkok Hub", lat: 13.5, lng: 100.7 }, description: "On the way", timestamp: "2026-11-07T12:00:00Z" }] } }));
  res.statusCode = 404; res.end(JSON.stringify({ success: false, error: "not found" }));
}).listen(8081, "127.0.0.1");
