import { createServer } from "node:http";

let shipment = {
  id: "ORD-TEST1",
  trackingNumber: "TH-TEST-1",
  status: "in_transit",
  carrier: "Flash Express",
  customer: {
    name: "Sender",
    zipcode: "10900",
    subDistrict: "ลาดยาว",
    district: "จตุจักร",
    province: "กรุงเทพมหานคร",
    coords: { lat: 13.7, lng: 100.5 },
  },
  receiver: {
    name: "Receiver",
    zipcode: "10800",
    subDistrict: "วงศ์สว่าง",
    district: "บางซื่อ",
    province: "กรุงเทพมหานคร",
    coords: { lat: 13.3, lng: 100.9 },
  },
  origin: "Bangkok",
  destination: "Chonburi",
  currentCoords: { lat: 13.5, lng: 100.7 },
  weight: 2,
  items: 2,
  estimatedDelivery: "2026-12-01T00:00:00Z",
  createdAt: "2026-11-01T00:00:00Z",
  progress: 55,
};
const originalShipment = structuredClone(shipment);

const manifestStatuses = [
  "in_transit",
  "delivered",
  "pending",
  "in_transit",
  "delivered",
  "pending",
];
const createManifest = () =>
  Array.from({ length: 12 }, (_, index) => {
    const number = String(index + 1).padStart(2, "0");
    return {
      ...structuredClone(shipment),
      id: `ORD-MANIFEST-${number}`,
      trackingNumber: `TH-MANIFEST-${number}`,
      status: manifestStatuses[index % manifestStatuses.length],
      customer: { ...shipment.customer, name: `Sender ${number}` },
      receiver: { ...shipment.receiver, name: `Receiver ${number}` },
    };
  });
let manifestRecords = createManifest();

let analyticsMode = "normal";
const authUsers = new Map([
  [
    "demo@shipitout.example",
    { id: 2, name: "Ship It Out Demo", email: "demo@shipitout.example", role: "customer" },
  ],
  [
    "staff@shipitout.example",
    { id: 1, name: "Preview Staff", email: "staff@shipitout.example", role: "staff" },
  ],
]);
const authPasswords = new Map([
  ["demo@shipitout.example", "demo-read-only-2026"],
  ["staff@shipitout.example", "preview-staff-password"],
]);
const authSessions = new Map();
let shipmentUpdateFailure = false;
let manifestDeleteFailure = false;
let nextShipmentNumber = 1;

createServer(async (req, res) => {
  res.setHeader("content-type", "application/json");
  const cookies = Object.fromEntries(
    (req.headers.cookie || "")
      .split("; ")
      .filter(Boolean)
      .map((item) => item.split("=")),
  );
  const user = authSessions.get(cookies.jwt);
  if (req.method === "POST" && req.url === "/__test/shipment?reset=1") {
    shipment = structuredClone(originalShipment);
    manifestRecords = createManifest();
    shipmentUpdateFailure = false;
    res.statusCode = 204;
    return res.end();
  }
  if (req.method === "POST" && req.url === "/__test/shipment?update=failure") {
    shipmentUpdateFailure = true;
    res.statusCode = 204;
    return res.end();
  }
  if (req.method === "POST" && req.url === "/__test/manifest?reset=1") {
    manifestRecords = createManifest();
    manifestDeleteFailure = false;
    res.statusCode = 204;
    return res.end();
  }
  if (req.method === "POST" && req.url === "/__test/manifest?delete=failure") {
    manifestDeleteFailure = true;
    res.statusCode = 204;
    return res.end();
  }
  if (
    req.method === "DELETE" &&
    new URL(req.url ?? "/", "http://localhost").pathname.startsWith("/api/shipments/")
  ) {
    if (!user) {
      res.statusCode = 401;
      return res.end(JSON.stringify({ success: false, error: "unauthorized" }));
    }
    if (user.role !== "staff") {
      res.statusCode = 403;
      return res.end(JSON.stringify({ success: false, error: "forbidden" }));
    }
    if (manifestDeleteFailure) {
      res.statusCode = 503;
      return res.end(JSON.stringify({ success: false, error: "shipment delete unavailable" }));
    }
    const id = decodeURIComponent(
      new URL(req.url, "http://localhost").pathname.slice("/api/shipments/".length),
    );
    manifestRecords = manifestRecords.filter((record) => record.id !== id);
    return res.end(JSON.stringify({ success: true, data: null }));
  }
  if (req.method === "POST" && req.url?.startsWith("/__test/analytics?mode=")) {
    analyticsMode = new URL(req.url, "http://localhost").searchParams.get("mode") ?? "normal";
    res.statusCode = 204;
    return res.end();
  }
  if (analyticsMode === "failure" && req.url?.startsWith("/api/analytics/")) {
    res.statusCode = 503;
    return res.end(JSON.stringify({ success: false, error: "unavailable" }));
  }
  if (analyticsMode === "empty" && req.url === "/api/analytics/overview")
    return res.end(
      JSON.stringify({
        success: true,
        data: {
          total: 0,
          active: 0,
          delivered: 0,
          by_status: null,
          by_region: ["Central", "East", "North", "West", "North-east", "South"].map((name) => ({
            name,
            total: 0,
          })),
        },
      }),
    );
  if (analyticsMode === "empty" && req.url === "/api/analytics/timeseries")
    return res.end(
      JSON.stringify({ success: true, data: { by_month: null, by_day_of_week: null } }),
    );
  if (req.url === "/api/analytics/overview")
    return res.end(
      JSON.stringify({
        success: true,
        data: {
          total: 126,
          active: 44,
          delivered: 82,
          by_status: [
            { status: "delivered", count: 82 },
            { status: "in_transit", count: 28 },
            { status: "pending", count: 16 },
          ],
          by_region: [
            { name: "Bangkok", total: 64 },
            { name: "Chonburi", total: 38 },
            { name: "Chiang Mai", total: 24 },
          ],
        },
      }),
    );
  if (req.url === "/api/analytics/timeseries")
    return res.end(
      JSON.stringify({
        success: true,
        data: {
          by_month: [
            { month: "Jun", count: 18 },
            { month: "Jul", count: 24 },
            { month: "Aug", count: 31 },
            { month: "Sep", count: 26 },
            { month: "Oct", count: 27 },
          ],
          by_day_of_week: [
            { day: "Mon", count: 21 },
            { day: "Tue", count: 18 },
            { day: "Wed", count: 26 },
            { day: "Thu", count: 16 },
            { day: "Fri", count: 24 },
            { day: "Sat", count: 13 },
            { day: "Sun", count: 8 },
          ],
        },
      }),
    );
  if (req.url === "/api/auth/me") {
    if (!user) {
      res.statusCode = 401;
      return res.end(JSON.stringify({ success: false, error: "not authenticated" }));
    }
    return res.end(JSON.stringify({ success: true, data: user }));
  }
  if (req.url === "/api/auth/login" || req.url === "/api/auth/register") {
    const chunks = [];
    for await (const chunk of req) chunks.push(chunk);
    const body = JSON.parse(Buffer.concat(chunks).toString() || "{}");
    let account = authUsers.get(body.email);
    if (req.url === "/api/auth/register") {
      if (account) {
        res.statusCode = 409;
        return res.end(JSON.stringify({ success: false, error: "email already registered" }));
      }
      account = { id: nextUserId++, name: body.name, email: body.email, role: "customer" };
      authUsers.set(body.email, account);
      authPasswords.set(body.email, body.password);
    } else if (!account || authPasswords.get(body.email) !== body.password) {
      res.statusCode = 401;
      return res.end(JSON.stringify({ success: false, error: "invalid email or password" }));
    }
    const token = `preview-session-${account.id}-${Date.now()}`;
    authSessions.set(token, account);
    res.setHeader("set-cookie", `jwt=${token}; Path=/; HttpOnly; SameSite=None; Secure`);
    return res.end(JSON.stringify({ success: true, data: { user: account, token } }));
  }
  if (req.url === "/api/auth/logout" && req.method === "POST") {
    authSessions.delete(cookies.jwt);
    res.setHeader("set-cookie", "jwt=; Path=/; HttpOnly; SameSite=None; Secure; Max-Age=0");
    return res.end(JSON.stringify({ success: true, data: { message: "logged out" } }));
  }
  if (req.url === "/api/shipments" && req.method === "POST") {
    if (!user) {
      res.statusCode = 401;
      return res.end(JSON.stringify({ success: false, error: "not authenticated" }));
    }
    if (user.role !== "admin" && user.role !== "staff") {
      res.statusCode = 403;
      return res.end(JSON.stringify({ success: false, error: "staff role required" }));
    }
    let body = "";
    for await (const chunk of req) body += chunk;
    const input = JSON.parse(body);
    const number = nextShipmentNumber++;
    const record = {
      ...structuredClone(shipment),
      id: `ORD-CREATED-${number}`,
      trackingNumber: `TH-CREATED-${number}`,
      status: "pending",
      customer: input.customer,
      receiver: input.receiver,
      carrier: input.carrier,
      origin: `${input.customer.subDistrict}, ${input.customer.district}, ${input.customer.province}`,
      destination: `${input.receiver.subDistrict}, ${input.receiver.district}, ${input.receiver.province}`,
      currentCoords: input.customer.coords,
      weight: input.weight,
      items: input.items,
      createdAt: new Date().toISOString(),
      progress: 0,
    };
    manifestRecords.unshift(record);
    res.statusCode = 201;
    return res.end(JSON.stringify({ success: true, data: record }));
  }
  if (req.url?.startsWith("/api/shipments?")) {
    const query = new URL(req.url, "http://localhost").searchParams;
    if (query.has("exclude_status")) {
      return res.end(
        JSON.stringify({
          success: true,
          data: [shipment],
          pagination: { page: 1, limit: 3, total: 1, totalPages: 1 },
        }),
      );
    }
    const search = (query.get("search") ?? "").toLowerCase();
    const status = query.get("status");
    const filtered = manifestRecords.filter((record) => {
      const matchesSearch =
        `${record.id} ${record.trackingNumber} ${record.customer.name} ${record.receiver.name} ${record.origin} ${record.destination}`
          .toLowerCase()
          .includes(search);
      return matchesSearch && (!status || record.status === status);
    });
    const page = Number(query.get("page") || 1);
    const limit = Number(query.get("limit") || 10);
    const start = (page - 1) * limit;
    return res.end(
      JSON.stringify({
        success: true,
        data: filtered.slice(start, start + limit),
        pagination: {
          page,
          limit,
          total: filtered.length,
          totalPages: Math.max(1, Math.ceil(filtered.length / limit)),
        },
      }),
    );
  }
  if (req.url === "/api/hubs")
    return res.end(
      JSON.stringify({
        success: true,
        data: [
          {
            id: "HUB-001",
            name: "Bangkok Hub",
            carrierId: "THUN",
            address: "Bang Rak, Bangkok",
            coords: { lat: 13.7563, lng: 100.5018 },
            capacity: 1200,
            currentUtilization: 840,
            status: "active",
            createdAt: "2026-01-01T00:00:00Z",
          },
          {
            id: "HUB-002",
            name: "Chonburi Hub",
            carrierId: "THUN",
            address: "Mueang Chonburi, Chonburi",
            coords: { lat: 13.3611, lng: 100.9847 },
            capacity: 800,
            currentUtilization: 0,
            status: "maintenance",
            createdAt: "2026-01-02T00:00:00Z",
          },
          {
            id: "HUB-003",
            name: "Chiang Mai Hub",
            carrierId: "THUN",
            address: "Mueang Chiang Mai, Chiang Mai",
            coords: { lat: 18.7883, lng: 98.9853 },
            capacity: 950,
            currentUtilization: 610,
            status: "active",
            createdAt: "2026-01-03T00:00:00Z",
          },
          {
            id: "HUB-004",
            name: "Phuket Hub",
            carrierId: "THUN",
            address: "Mueang Phuket, Phuket",
            coords: { lat: 7.8804, lng: 98.3923 },
            capacity: 700,
            currentUtilization: 430,
            status: "active",
            createdAt: "2026-01-04T00:00:00Z",
          },
          {
            id: "HUB-005",
            name: "Khon Kaen Hub",
            carrierId: "THUN",
            address: "Mueang Khon Kaen, Khon Kaen",
            coords: { lat: 16.4419, lng: 102.835 },
            capacity: 820,
            currentUtilization: 570,
            status: "active",
            createdAt: "2026-01-05T00:00:00Z",
          },
          {
            id: "HUB-006",
            name: "Nakhon Ratchasima Hub",
            carrierId: "THUN",
            address: "Mueang Nakhon Ratchasima, Nakhon Ratchasima",
            coords: { lat: 14.9799, lng: 102.0977 },
            capacity: 900,
            currentUtilization: 530,
            status: "active",
            createdAt: "2026-01-06T00:00:00Z",
          },
          {
            id: "HUB-007",
            name: "Hat Yai Hub",
            carrierId: "THUN",
            address: "Hat Yai, Songkhla",
            coords: { lat: 7.0084, lng: 100.4747 },
            capacity: 650,
            currentUtilization: 390,
            status: "active",
            createdAt: "2026-01-07T00:00:00Z",
          },
          {
            id: "HUB-008",
            name: "Udon Thani Hub",
            carrierId: "THUN",
            address: "Mueang Udon Thani, Udon Thani",
            coords: { lat: 17.4138, lng: 102.7872 },
            capacity: 680,
            currentUtilization: 410,
            status: "active",
            createdAt: "2026-01-08T00:00:00Z",
          },
          {
            id: "HUB-009",
            name: "Rayong Hub",
            carrierId: "THUN",
            address: "Mueang Rayong, Rayong",
            coords: { lat: 12.6814, lng: 101.2816 },
            capacity: 600,
            currentUtilization: 350,
            status: "active",
            createdAt: "2026-01-09T00:00:00Z",
          },
          {
            id: "HUB-010",
            name: "Ayutthaya Hub",
            carrierId: "THUN",
            address: "Phra Nakhon Si Ayutthaya, Ayutthaya",
            coords: { lat: 14.3532, lng: 100.5689 },
            capacity: 720,
            currentUtilization: 470,
            status: "active",
            createdAt: "2026-01-10T00:00:00Z",
          },
          {
            id: "HUB-011",
            name: "Surat Thani Hub",
            carrierId: "THUN",
            address: "Mueang Surat Thani, Surat Thani",
            coords: { lat: 9.1382, lng: 99.3217 },
            capacity: 540,
            currentUtilization: 260,
            status: "active",
            createdAt: "2026-01-11T00:00:00Z",
          },
          {
            id: "HUB-012",
            name: "Lampang Hub",
            carrierId: "THUN",
            address: "Mueang Lampang, Lampang",
            coords: { lat: 18.2888, lng: 99.4909 },
            capacity: 500,
            currentUtilization: 300,
            status: "active",
            createdAt: "2026-01-12T00:00:00Z",
          },
        ],
      }),
    );
  if (req.method === "PUT" && req.url === "/api/shipments/ORD-TEST1") {
    if (!user) {
      res.statusCode = 401;
      return res.end(JSON.stringify({ success: false, error: "not authenticated" }));
    }
    if (user.role !== "staff") {
      res.statusCode = 403;
      return res.end(JSON.stringify({ success: false, error: "staff role required" }));
    }
    if (shipmentUpdateFailure) {
      res.statusCode = 503;
      return res.end(JSON.stringify({ success: false, error: "shipment update unavailable" }));
    }
    let body = "";
    for await (const chunk of req) body += chunk;
    const update = JSON.parse(body);
    shipment = {
      ...shipment,
      ...update,
      origin: `${update.customer.subDistrict}, ${update.customer.district}, ${update.customer.province}`,
      destination: `${update.receiver.subDistrict}, ${update.receiver.district}, ${update.receiver.province}`,
      currentCoords: update.customer.coords,
    };
    manifestRecords = manifestRecords.map((record) =>
      record.id === shipment.id ? shipment : record,
    );
    return res.end(JSON.stringify({ success: true, data: shipment }));
  }
  const createdShipment = manifestRecords.find(
    (record) => req.url === `/api/shipments/${encodeURIComponent(record.id)}`,
  );
  if (createdShipment) return res.end(JSON.stringify({ success: true, data: createdShipment }));
  if (req.url === "/api/shipments/ORD-TEST1")
    return res.end(JSON.stringify({ success: true, data: shipment }));
  if (req.url === "/api/shipments/missing") {
    res.statusCode = 404;
    return res.end(JSON.stringify({ success: false, error: "shipment not found" }));
  }
  if (req.url === "/api/track/TH-TEST-1")
    return res.end(
      JSON.stringify({
        success: true,
        data: {
          shipment,
          events: [
            {
              status: "In Transit",
              location: { name: "Bangkok Hub", lat: 13.5, lng: 100.7 },
              description: "On the way",
              timestamp: "2026-11-07T12:00:00Z",
            },
          ],
        },
      }),
    );
  res.statusCode = 404;
  res.end(JSON.stringify({ success: false, error: "not found" }));
}).listen(8081, "127.0.0.1");
