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

createServer((req, res) => {
  res.setHeader("content-type", "application/json");
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
