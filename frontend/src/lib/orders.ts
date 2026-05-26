export type ShipmentStatus =
  | "pending"
  | "picked_up"
  | "departed"
  | "in_transit"
  | "out_for_delivery"
  | "delivered"
  | "delayed";

export interface Location {
  name: string;
  lat: number;
  lng: number;
}

export interface TrackingEvent {
  timestamp: string;
  location: Location;
  status: string;
  description: string;
}

export interface GeoPoint {
  lat: number;
  lng: number;
}

export interface ContactInfo {
  name: string;
  zipcode: string;
  subDistrict: string;
  district: string;
  province: string;
  coords: GeoPoint;
}

export interface Order {
  id: string;
  trackingNumber: string;
  customer: ContactInfo;
  receiver: ContactInfo;
  origin: string;
  destination: string;
  currentCoords: GeoPoint;
  status: ShipmentStatus;
  carrier: string;
  hubId?: string;
  weight: string;
  items: number;
  estimatedDelivery: string;
  estimatedDeliveryRaw?: string;
  createdAt: string;
  progress: number;
  events: TrackingEvent[];
}

export interface PaginatedResponse<T> {
  data: T[];
  pagination: {
    page: number;
    limit: number;
    total: number;
    totalPages: number;
  };
}

export const orders: Order[] = [
  {
    id: "ORD-10245",
    trackingNumber: "TH202600100",
    customer: {
      name: "Aria Nakamura",
      zipcode: "3011",
      subDistrict: "Stadsdriehoek",
      district: "Centrum",
      province: "Zuid-Holland",
      coords: { lat: 51.9244, lng: 4.4777 },
    },
    receiver: {
      name: "James Mitchell",
      zipcode: "11201",
      subDistrict: "DUMBO",
      district: "Brooklyn",
      province: "New York",
      coords: { lat: 40.6782, lng: -73.9442 },
    },
    origin: "Stadsdriehoek, Centrum, Zuid-Holland",
    destination: "DUMBO, Brooklyn, New York",
    currentCoords: { lat: 46.5, lng: -34.5 },
    status: "in_transit",
    carrier: "Pacific Freight",
    weight: "12.4 kg",
    items: 3,
    estimatedDelivery: "May 24, 2026",
    createdAt: "May 18, 2026",
    progress: 62,
    events: [
      {
        timestamp: "May 21, 10:00",
        location: { name: "North Sea", lat: 62.5, lng: -2.5 },
        status: "Delayed",
        description: "Weather hold, 36h estimate.",
      },
      {
        timestamp: "May 18, 22:14",
        location: { name: "Oslo Port", lat: 59.9, lng: 10.75 },
        status: "Departed",
        description: "",
      },
      {
        timestamp: "May 16, 09:02",
        location: { name: "Oslo Warehouse", lat: 59.9139, lng: 10.7522 },
        status: "Picked up",
        description: "",
      },
    ],
  },
  {
    id: "ORD-10249",
    trackingNumber: "TH202600101",
    customer: {
      name: "Priya Anand",
      zipcode: "400001",
      subDistrict: "Fort",
      district: "Mumbai City",
      province: "Maharashtra",
      coords: { lat: 19.076, lng: 72.8777 },
    },
    receiver: {
      name: "Ahmed Al-Rashid",
      zipcode: "00000",
      subDistrict: "Downtown",
      district: "Dubai",
      province: "Dubai",
      coords: { lat: 25.2048, lng: 55.2708 },
    },
    origin: "Fort, Mumbai City, Maharashtra",
    destination: "Downtown, Dubai, Dubai",
    currentCoords: { lat: 19.076, lng: 72.8777 },
    status: "pending",
    carrier: "Gulf Logistics",
    weight: "0.9 kg",
    items: 1,
    estimatedDelivery: "May 25, 2026",
    createdAt: "May 20, 2026",
    progress: 8,
    events: [
      {
        timestamp: "May 20, 16:40",
        location: { name: "Mumbai Warehouse", lat: 19.076, lng: 72.8777 },
        status: "Label created",
        description: "Awaiting pickup.",
      },
    ],
  },
  {
    id: "ORD-10250",
    trackingNumber: "TH202600102",
    customer: {
      name: "Hugo Martín",
      zipcode: "08001",
      subDistrict: "El Raval",
      district: "Ciutat Vella",
      province: "Barcelona",
      coords: { lat: 41.3851, lng: 2.1734 },
    },
    receiver: {
      name: "Camille Dubois",
      zipcode: "13001",
      subDistrict: "Vieux-Port",
      district: "Marseille",
      province: "Provence-Alpes",
      coords: { lat: 43.2965, lng: 5.3698 },
    },
    origin: "El Raval, Ciutat Vella, Barcelona",
    destination: "Vieux-Port, Marseille, Provence-Alpes",
    currentCoords: { lat: 42.6986, lng: 2.8954 },
    status: "in_transit",
    carrier: "Mediterranean Freight",
    weight: "7.3 kg",
    items: 4,
    estimatedDelivery: "May 22, 2026",
    createdAt: "May 19, 2026",
    progress: 55,
    events: [
      {
        timestamp: "May 20, 13:22",
        location: { name: "Perpignan", lat: 42.69, lng: 2.89 },
        status: "In transit",
        description: "Crossing border.",
      },
      {
        timestamp: "May 19, 18:00",
        location: { name: "Barcelona Hub", lat: 41.3851, lng: 2.1734 },
        status: "Departed",
        description: "",
      },
    ],
  },
  {
    id: "ORD-10251",
    trackingNumber: "TH202600103",
    customer: {
      name: "สมหญิง ใจดี",
      zipcode: "10200",
      subDistrict: "Bang Rak",
      district: "Bang Rak",
      province: "Bangkok",
      coords: { lat: 13.7279, lng: 100.5242 },
    },
    receiver: {
      name: "สมศักดิ์ อินทร์แก้ว",
      zipcode: "50000",
      subDistrict: "Sri Phum",
      district: "Mueang",
      province: "Chiang Mai",
      coords: { lat: 18.7883, lng: 98.9853 },
    },
    origin: "Bang Rak, Bang Rak, Bangkok",
    destination: "Sri Phum, Mueang, Chiang Mai",
    currentCoords: { lat: 18.7883, lng: 98.9853 },
    status: "delivered",
    carrier: "Pacific Freight",
    weight: "2.1 kg",
    items: 1,
    estimatedDelivery: "May 20, 2026",
    createdAt: "May 16, 2026",
    progress: 100,
    events: [
      {
        timestamp: "May 20, 14:30",
        location: { name: "Sri Phum, Mueang, Chiang Mai", lat: 18.7883, lng: 98.9853 },
        status: "Delivered",
        description: "Delivered to recipient.",
      },
      {
        timestamp: "May 19, 10:00",
        location: { name: "Rayong Hub", lat: 12.6814, lng: 101.2817 },
        status: "Out for Delivery",
        description: "Out for delivery.",
      },
      {
        timestamp: "May 18, 16:00",
        location: { name: "Chanthaburi Hub", lat: 12.6096, lng: 102.1041 },
        status: "In Transit",
        description: "Transit to next hub.",
      },
      {
        timestamp: "May 17, 09:00",
        location: { name: "Laem Chabang Port Hub", lat: 13.0833, lng: 100.8833 },
        status: "Departed",
        description: "In transit to hub.",
      },
      {
        timestamp: "May 16, 15:00",
        location: { name: "Bang Rak, Bang Rak, Bangkok", lat: 13.7279, lng: 100.5242 },
        status: "Picked Up",
        description: "Parcel collected from sender.",
      },
      {
        timestamp: "May 16, 08:00",
        location: { name: "Bang Rak, Bang Rak, Bangkok", lat: 13.7279, lng: 100.5242 },
        status: "Label Created",
        description: "Awaiting pickup.",
      },
    ],
  },
  {
    id: "ORD-10252",
    trackingNumber: "TH202600104",
    customer: {
      name: "Lek Thongdee",
      zipcode: "10270",
      subDistrict: "Samrong Nua",
      district: "Mueang",
      province: "Samut Prakan",
      coords: { lat: 13.6509, lng: 100.6016 },
    },
    receiver: {
      name: "Paisan Saelaeo",
      zipcode: "83000",
      subDistrict: "Patong",
      district: "Kathu",
      province: "Phuket",
      coords: { lat: 7.8961, lng: 98.2966 },
    },
    origin: "Samrong Nua, Mueang, Samut Prakan",
    destination: "Patong, Kathu, Phuket",
    currentCoords: { lat: 13.65, lng: 100.6 },
    status: "picked_up",
    carrier: "Gulf Logistics",
    weight: "8.5 kg",
    items: 5,
    estimatedDelivery: "May 29, 2026",
    createdAt: "May 25, 2026",
    progress: 15,
    events: [
      {
        timestamp: "May 25, 14:30",
        location: { name: "Samrong Nua, Mueang, Samut Prakan", lat: 13.6509, lng: 100.6016 },
        status: "Picked Up",
        description: "Parcel collected from sender.",
      },
      {
        timestamp: "May 25, 09:00",
        location: { name: "Samrong Nua, Mueang, Samut Prakan", lat: 13.6509, lng: 100.6016 },
        status: "Label Created",
        description: "Awaiting pickup.",
      },
    ],
  },
  {
    id: "ORD-10253",
    trackingNumber: "TH202600105",
    customer: {
      name: "Rattana Pimthong",
      zipcode: "11000",
      subDistrict: "Tha Sai",
      district: "Mueang",
      province: "Nonthaburi",
      coords: { lat: 13.8548, lng: 100.5146 },
    },
    receiver: {
      name: "Somsak Kaewphon",
      zipcode: "57000",
      subDistrict: "Wiang",
      district: "Mueang",
      province: "Chiang Rai",
      coords: { lat: 19.9072, lng: 99.8325 },
    },
    origin: "Tha Sai, Mueang, Nonthaburi",
    destination: "Wiang, Mueang, Chiang Rai",
    currentCoords: { lat: 16.5, lng: 100.5 },
    status: "in_transit",
    carrier: "Skyline Express",
    weight: "4.2 kg",
    items: 2,
    estimatedDelivery: "May 28, 2026",
    createdAt: "May 24, 2026",
    progress: 35,
    events: [
      {
        timestamp: "May 25, 11:00",
        location: { name: "Nakhon Sawan", lat: 15.7167, lng: 100.1333 },
        status: "In Transit",
        description: "Transit to next hub.",
      },
      {
        timestamp: "May 24, 16:00",
        location: { name: "Pattaya Hub", lat: 12.9236, lng: 100.8825 },
        status: "Departed",
        description: "In transit to hub.",
      },
      {
        timestamp: "May 24, 10:00",
        location: { name: "Tha Sai, Mueang, Nonthaburi", lat: 13.8548, lng: 100.5146 },
        status: "Picked Up",
        description: "Parcel collected from sender.",
      },
      {
        timestamp: "May 24, 08:00",
        location: { name: "Tha Sai, Mueang, Nonthaburi", lat: 13.8548, lng: 100.5146 },
        status: "Label Created",
        description: "Awaiting pickup.",
      },
    ],
  },
  {
    id: "ORD-10254",
    trackingNumber: "TH202600106",
    customer: {
      name: "Suthee Sriwat",
      zipcode: "12000",
      subDistrict: "Prachathipat",
      district: "Thanyaburi",
      province: "Pathum Thani",
      coords: { lat: 13.9782, lng: 100.6147 },
    },
    receiver: {
      name: "Abdullah Samae",
      zipcode: "90110",
      subDistrict: "Kho Hong",
      district: "Hat Yai",
      province: "Songkhla",
      coords: { lat: 7.0088, lng: 100.4747 },
    },
    origin: "Prachathipat, Thanyaburi, Pathum Thani",
    destination: "Kho Hong, Hat Yai, Songkhla",
    currentCoords: { lat: 11.5, lng: 99.5 },
    status: "delayed",
    carrier: "Trans-Atlantic Cargo",
    weight: "6.7 kg",
    items: 4,
    estimatedDelivery: "May 24, 2026",
    createdAt: "May 20, 2026",
    progress: 60,
    events: [
      {
        timestamp: "May 25, 08:00",
        location: { name: "Chanthaburi Hub", lat: 12.6096, lng: 102.1041 },
        status: "Delayed",
        description: "Unexpected issue encountered.",
      },
      {
        timestamp: "May 23, 14:00",
        location: { name: "Chanthaburi Hub", lat: 12.6096, lng: 102.1041 },
        status: "In Transit",
        description: "Transit to next hub.",
      },
      {
        timestamp: "May 22, 09:00",
        location: { name: "Laem Chabang Port Hub", lat: 13.0833, lng: 100.8833 },
        status: "Departed",
        description: "In transit to hub.",
      },
      {
        timestamp: "May 21, 16:00",
        location: { name: "Prachathipat, Thanyaburi, Pathum Thani", lat: 13.9782, lng: 100.6147 },
        status: "Picked Up",
        description: "Parcel collected from sender.",
      },
      {
        timestamp: "May 21, 08:00",
        location: { name: "Prachathipat, Thanyaburi, Pathum Thani", lat: 13.9782, lng: 100.6147 },
        status: "Label Created",
        description: "Awaiting pickup.",
      },
    ],
  },
  {
    id: "ORD-10255",
    trackingNumber: "TH202600107",
    customer: {
      name: "Kanjana Sirichok",
      zipcode: "13000",
      subDistrict: "Pratu Chai",
      district: "Phra Nakhon Si Ayutthaya",
      province: "Phra Nakhon Si Ayutthaya",
      coords: { lat: 14.3542, lng: 100.5547 },
    },
    receiver: {
      name: "Thanaphon Chansri",
      zipcode: "40000",
      subDistrict: "Nai Mueang",
      district: "Mueang",
      province: "Khon Kaen",
      coords: { lat: 16.4322, lng: 102.8236 },
    },
    origin: "Pratu Chai, Phra Nakhon Si Ayutthaya, Phra Nakhon Si Ayutthaya",
    destination: "Nai Mueang, Mueang, Khon Kaen",
    currentCoords: { lat: 14.3542, lng: 100.5547 },
    status: "pending",
    carrier: "Nordic Lines",
    weight: "3.3 kg",
    items: 2,
    estimatedDelivery: "May 29, 2026",
    createdAt: "May 26, 2026",
    progress: 0,
    events: [
      {
        timestamp: "May 26, 10:00",
        location: {
          name: "Pratu Chai, Phra Nakhon Si Ayutthaya, Phra Nakhon Si Ayutthaya",
          lat: 14.3542,
          lng: 100.5547,
        },
        status: "Label Created",
        description: "Awaiting pickup.",
      },
    ],
  },
  {
    id: "ORD-10256",
    trackingNumber: "TH202600108",
    customer: {
      name: "Naruemon Sukkasem",
      zipcode: "10100",
      subDistrict: "Pathum Wan",
      district: "Pathum Wan",
      province: "Bangkok",
      coords: { lat: 13.7466, lng: 100.5326 },
    },
    receiver: {
      name: "Adun Klahan",
      zipcode: "30000",
      subDistrict: "Nai Mueang",
      district: "Mueang",
      province: "Nakhon Ratchasima",
      coords: { lat: 14.975, lng: 102.0825 },
    },
    origin: "Pathum Wan, Pathum Wan, Bangkok",
    destination: "Nai Mueang, Mueang, Nakhon Ratchasima",
    currentCoords: { lat: 14.975, lng: 102.0825 },
    status: "delivered",
    carrier: "Pacific Freight",
    weight: "1.8 kg",
    items: 1,
    estimatedDelivery: "May 23, 2026",
    createdAt: "May 19, 2026",
    progress: 100,
    events: [
      {
        timestamp: "May 23, 11:30",
        location: { name: "Nai Mueang, Mueang, Nakhon Ratchasima", lat: 14.975, lng: 102.0825 },
        status: "Delivered",
        description: "Delivered to recipient.",
      },
      {
        timestamp: "May 22, 15:00",
        location: { name: "Laem Chabang Port Hub", lat: 13.0833, lng: 100.8833 },
        status: "Out for Delivery",
        description: "Out for delivery.",
      },
      {
        timestamp: "May 21, 10:00",
        location: { name: "Chachoengsao Hub", lat: 13.6883, lng: 101.0719 },
        status: "In Transit",
        description: "Transit to next hub.",
      },
      {
        timestamp: "May 20, 14:00",
        location: { name: "Laem Chabang Port Hub", lat: 13.0833, lng: 100.8833 },
        status: "Departed",
        description: "In transit to hub.",
      },
      {
        timestamp: "May 19, 16:30",
        location: { name: "Pathum Wan, Pathum Wan, Bangkok", lat: 13.7466, lng: 100.5326 },
        status: "Picked Up",
        description: "Parcel collected from sender.",
      },
      {
        timestamp: "May 19, 08:00",
        location: { name: "Pathum Wan, Pathum Wan, Bangkok", lat: 13.7466, lng: 100.5326 },
        status: "Label Created",
        description: "Awaiting pickup.",
      },
    ],
  },
  {
    id: "ORD-10257",
    trackingNumber: "TH202600109",
    customer: {
      name: "Saroj Charoensuk",
      zipcode: "20110",
      subDistrict: "Bang Phra",
      district: "Si Racha",
      province: "Chonburi",
      coords: { lat: 13.1173, lng: 100.9256 },
    },
    receiver: {
      name: "Wasana Khongman",
      zipcode: "84000",
      subDistrict: "Talat",
      district: "Mueang",
      province: "Surat Thani",
      coords: { lat: 9.1382, lng: 99.3214 },
    },
    origin: "Bang Phra, Si Racha, Chonburi",
    destination: "Talat, Mueang, Surat Thani",
    currentCoords: { lat: 9.5, lng: 99.2 },
    status: "out_for_delivery",
    carrier: "Trans-Atlantic Cargo",
    weight: "10.0 kg",
    items: 6,
    estimatedDelivery: "May 26, 2026",
    createdAt: "May 22, 2026",
    progress: 80,
    events: [
      {
        timestamp: "May 26, 09:00",
        location: { name: "Chachoengsao Hub", lat: 13.6883, lng: 101.0719 },
        status: "Out for Delivery",
        description: "Out for delivery.",
      },
      {
        timestamp: "May 25, 14:00",
        location: { name: "Chachoengsao Hub", lat: 13.6883, lng: 101.0719 },
        status: "In Transit",
        description: "Transit to next hub.",
      },
      {
        timestamp: "May 24, 10:00",
        location: { name: "Laem Chabang Port Hub", lat: 13.0833, lng: 100.8833 },
        status: "Departed",
        description: "In transit to hub.",
      },
      {
        timestamp: "May 23, 14:00",
        location: { name: "Bang Phra, Si Racha, Chonburi", lat: 13.1173, lng: 100.9256 },
        status: "Picked Up",
        description: "Parcel collected from sender.",
      },
      {
        timestamp: "May 23, 08:00",
        location: { name: "Bang Phra, Si Racha, Chonburi", lat: 13.1173, lng: 100.9256 },
        status: "Label Created",
        description: "Awaiting pickup.",
      },
    ],
  },
  {
    id: "ORD-10258",
    trackingNumber: "TH202600110",
    customer: {
      name: "Prasert Wongdee",
      zipcode: "21160",
      subDistrict: "Nikhom Phatthana",
      district: "Nikhom Phatthana",
      province: "Rayong",
      coords: { lat: 12.7167, lng: 101.15 },
    },
    receiver: {
      name: "Buntham Phimpha",
      zipcode: "41000",
      subDistrict: "Mak Khaeng",
      district: "Mueang",
      province: "Udon Thani",
      coords: { lat: 17.4132, lng: 102.7856 },
    },
    origin: "Nikhom Phatthana, Nikhom Phatthana, Rayong",
    destination: "Mak Khaeng, Mueang, Udon Thani",
    currentCoords: { lat: 15.0, lng: 102.0 },
    status: "in_transit",
    carrier: "Gulf Logistics",
    weight: "7.1 kg",
    items: 3,
    estimatedDelivery: "May 28, 2026",
    createdAt: "May 24, 2026",
    progress: 40,
    events: [
      {
        timestamp: "May 25, 16:00",
        location: { name: "Nakhon Ratchasima", lat: 14.975, lng: 102.0825 },
        status: "In Transit",
        description: "Transit to next hub.",
      },
      {
        timestamp: "May 24, 18:00",
        location: { name: "Rayong Hub", lat: 12.6814, lng: 101.2817 },
        status: "Departed",
        description: "In transit to hub.",
      },
      {
        timestamp: "May 24, 11:00",
        location: { name: "Nikhom Phatthana, Nikhom Phatthana, Rayong", lat: 12.7167, lng: 101.15 },
        status: "Picked Up",
        description: "Parcel collected from sender.",
      },
      {
        timestamp: "May 24, 08:00",
        location: { name: "Nikhom Phatthana, Nikhom Phatthana, Rayong", lat: 12.7167, lng: 101.15 },
        status: "Label Created",
        description: "Awaiting pickup.",
      },
    ],
  },
  {
    id: "ORD-10259",
    trackingNumber: "TH202600111",
    customer: {
      name: "Suda Maksuk",
      zipcode: "24110",
      subDistrict: "Bang Khla",
      district: "Bang Khla",
      province: "Chachoengsao",
      coords: { lat: 13.7167, lng: 101.2 },
    },
    receiver: {
      name: "Somporn Rakdee",
      zipcode: "80000",
      subDistrict: "Tha Wang",
      district: "Mueang",
      province: "Nakhon Si Thammarat",
      coords: { lat: 8.4333, lng: 99.9667 },
    },
    origin: "Bang Khla, Bang Khla, Chachoengsao",
    destination: "Tha Wang, Mueang, Nakhon Si Thammarat",
    currentCoords: { lat: 13.7, lng: 101.1 },
    status: "departed",
    carrier: "Mediterranean Freight",
    weight: "5.5 kg",
    items: 3,
    estimatedDelivery: "May 28, 2026",
    createdAt: "May 25, 2026",
    progress: 20,
    events: [
      {
        timestamp: "May 25, 17:00",
        location: { name: "Chachoengsao Hub", lat: 13.6883, lng: 101.0719 },
        status: "Departed",
        description: "In transit to hub.",
      },
      {
        timestamp: "May 25, 13:00",
        location: { name: "Bang Khla, Bang Khla, Chachoengsao", lat: 13.7167, lng: 101.2 },
        status: "Picked Up",
        description: "Parcel collected from sender.",
      },
      {
        timestamp: "May 25, 08:00",
        location: { name: "Bang Khla, Bang Khla, Chachoengsao", lat: 13.7167, lng: 101.2 },
        status: "Label Created",
        description: "Awaiting pickup.",
      },
    ],
  },
  {
    id: "ORD-10260",
    trackingNumber: "TH202600112",
    customer: {
      name: "Pimon Sukjai",
      zipcode: "10250",
      subDistrict: "Khlong Toei",
      district: "Khlong Toei",
      province: "Bangkok",
      coords: { lat: 13.7122, lng: 100.5638 },
    },
    receiver: {
      name: "Wira Wongkham",
      zipcode: "34000",
      subDistrict: "Nai Mueang",
      district: "Mueang",
      province: "Ubon Ratchathani",
      coords: { lat: 15.2296, lng: 104.8603 },
    },
    origin: "Khlong Toei, Khlong Toei, Bangkok",
    destination: "Nai Mueang, Mueang, Ubon Ratchathani",
    currentCoords: { lat: 15.2296, lng: 104.8603 },
    status: "delivered",
    carrier: "Pacific Freight",
    weight: "4.5 kg",
    items: 2,
    estimatedDelivery: "May 24, 2026",
    createdAt: "May 20, 2026",
    progress: 100,
    events: [
      {
        timestamp: "May 24, 10:00",
        location: { name: "Nai Mueang, Mueang, Ubon Ratchathani", lat: 15.2296, lng: 104.8603 },
        status: "Delivered",
        description: "Delivered to recipient.",
      },
      {
        timestamp: "May 23, 14:00",
        location: { name: "Chanthaburi Hub", lat: 12.6096, lng: 102.1041 },
        status: "Out for Delivery",
        description: "Out for delivery.",
      },
      {
        timestamp: "May 22, 16:00",
        location: { name: "Chanthaburi Hub", lat: 12.6096, lng: 102.1041 },
        status: "In Transit",
        description: "Transit to next hub.",
      },
      {
        timestamp: "May 21, 14:00",
        location: { name: "Laem Chabang Port Hub", lat: 13.0833, lng: 100.8833 },
        status: "Departed",
        description: "In transit to hub.",
      },
      {
        timestamp: "May 20, 17:00",
        location: { name: "Khlong Toei, Khlong Toei, Bangkok", lat: 13.7122, lng: 100.5638 },
        status: "Picked Up",
        description: "Parcel collected from sender.",
      },
      {
        timestamp: "May 20, 08:00",
        location: { name: "Khlong Toei, Khlong Toei, Bangkok", lat: 13.7122, lng: 100.5638 },
        status: "Label Created",
        description: "Awaiting pickup.",
      },
    ],
  },
];

export const statusLabels: Record<ShipmentStatus, string> = {
  pending: "Pending",
  picked_up: "Picked Up",
  departed: "Departed",
  in_transit: "In Transit",
  out_for_delivery: "Out for Delivery",
  delivered: "Delivered",
  delayed: "Delayed",
};

export function getOrder(id: string) {
  return orders.find((o) => o.id === id);
}
