export interface TrackingEvent {
  timestamp: string;
  location: { name: string; lat: number; lng: number };
  status: string;
  description?: string;
}
