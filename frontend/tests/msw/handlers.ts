import { http, HttpResponse } from "msw";

const BASE = "http://localhost:8080/api";

const mockUser = { id: 1, name: "Admin User", email: "admin@harborops.io", role: "admin", created_at: "2026-01-01T00:00:00Z" };

const mockShipments = [
  {
    id: "ORD-001", trackingNumber: "TH202600001",
    customer: { name: "John Doe", zipcode: "10100", subDistrict: "Bang Rak", district: "Bang Rak", province: "Bangkok", coords: { lat: 13.7279, lng: 100.5242 } },
    receiver: { name: "Jane Doe", zipcode: "50000", subDistrict: "Sri Phum", district: "Mueang", province: "Chiang Mai", coords: { lat: 18.7883, lng: 98.9853 } },
    origin: "Bang Rak, Bangkok", destination: "Sri Phum, Chiang Mai",
    currentCoords: { lat: 16.0, lng: 99.5 },
    status: "in_transit", carrier: "Pacific Freight", weight: 12.4, items: 3,
    estimatedDelivery: "2026-06-01T00:00:00Z", createdAt: "2026-05-28T00:00:00Z", progress: 62,
  },
];

const mockHubs = [
  { id: "hub-1", name: "Bangkok Hub", carrierId: "carrier-1", address: "123 Bangkok Rd", coords: { lat: 13.75, lng: 100.5 }, capacity: 1000, currentUtilization: 450, status: "active", createdAt: "2026-01-01T00:00:00Z" },
];

const mockAnalytics = {
  total: 100, active: 45, delivered: 55,
  by_status: [{ status: "in_transit", count: 30 }, { status: "delivered", count: 55 }],
  by_region: [{ name: "Bangkok", total: 40 }, { name: "Chiang Mai", total: 20 }],
};

const mockTimeseries = {
  by_month: [{ month: "2026-01", count: 10 }, { month: "2026-02", count: 15 }],
  by_day_of_week: [{ day: "Monday", count: 20 }, { day: "Tuesday", count: 18 }],
};

export const handlers = [
  http.get(`${BASE}/auth/me`, () => HttpResponse.json({ data: mockUser })),
  http.post(`${BASE}/auth/login`, () => HttpResponse.json({ data: { user: mockUser } })),
  http.post(`${BASE}/auth/register`, () => HttpResponse.json({ data: { user: mockUser } })),
  http.post(`${BASE}/auth/logout`, () => HttpResponse.json({ data: { success: true } })),
  http.get(`${BASE}/shipments`, ({ request }) => {
    const url = new URL(request.url);
    const limit = url.searchParams.get("limit");
    const status = url.searchParams.get("exclude_status");
    let data = mockShipments;
    if (status === "delivered") data = mockShipments.filter((s) => s.status !== "delivered");
    if (limit === "-1") return HttpResponse.json({ data, pagination: { page: 1, limit: -1, total: data.length, totalPages: 1 } });
    return HttpResponse.json({ data, pagination: { page: 1, limit: 10, total: data.length, totalPages: 1 } });
  }),
  http.get(`${BASE}/shipments/:id`, () => HttpResponse.json({ data: mockShipments[0] })),
  http.post(`${BASE}/shipments`, async ({ request }) => {
    const body = await request.json();
    return HttpResponse.json({ data: { ...mockShipments[0], ...body } }, { status: 201 });
  }),
  http.put(`${BASE}/shipments/:id`, async ({ request }) => {
    const body = await request.json();
    return HttpResponse.json({ data: { ...mockShipments[0], ...body } });
  }),
  http.delete(`${BASE}/shipments/:id`, () => HttpResponse.json({ data: { success: true } })),
  http.patch(`${BASE}/shipments/:id/status`, async ({ request }) => {
    const body = await request.json();
    return HttpResponse.json({ data: { ...mockShipments[0], ...body } });
  }),
  http.get(`${BASE}/track/:trackingNumber`, () =>
    HttpResponse.json({ data: { shipment: { id: "ORD-001" }, events: [] } }),
  ),
  http.get(`${BASE}/hubs`, () => HttpResponse.json({ data: mockHubs })),
  http.post(`${BASE}/hubs`, async ({ request }) => {
    const body = await request.json();
    return HttpResponse.json({ data: { ...mockHubs[0], ...body } }, { status: 201 });
  }),
  http.put(`${BASE}/hubs/:id`, async ({ request }) => {
    const body = await request.json();
    return HttpResponse.json({ data: { ...mockHubs[0], ...body } });
  }),
  http.delete(`${BASE}/hubs/:id`, () => HttpResponse.json({ data: { success: true } })),
  http.get(`${BASE}/analytics/overview`, () => HttpResponse.json({ data: mockAnalytics })),
  http.get(`${BASE}/analytics/timeseries`, () => HttpResponse.json({ data: mockTimeseries })),
];
