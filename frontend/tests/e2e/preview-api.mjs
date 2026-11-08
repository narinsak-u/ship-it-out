import { createServer } from "node:http";

const shipment = {
  id: "ORD-TEST1",
  trackingNumber: "TH-TEST-1",
  status: "in_transit",
  carrier: "Flash Express",
  customer: {
    name: "Sender",
    zipcode: "10100",
    subDistrict: "A",
    district: "B",
    province: "Bangkok",
    coords: { lat: 13.7, lng: 100.5 },
  },
  receiver: {
    name: "Receiver",
    zipcode: "20100",
    subDistrict: "C",
    district: "D",
    province: "Chonburi",
    coords: { lat: 13.3, lng: 100.9 },
  },
  origin: "Bangkok",
  destination: "Chonburi",
  currentCoords: { lat: 13.5, lng: 100.7 },
  weight: 2,
  estimatedDelivery: "2026-12-01T00:00:00Z",
  createdAt: "2026-11-01T00:00:00Z",
  progress: 55,
};

let analyticsMode = "normal";

createServer((req, res) => {
  res.setHeader("content-type", "application/json");
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
  if (req.url?.startsWith("/api/shipments?"))
    return res.end(
      JSON.stringify({
        success: true,
        data: [shipment],
        pagination: { page: 1, limit: 3, total: 1, totalPages: 1 },
      }),
    );
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
