import { notFound } from "next/navigation";
import { getShipment, getShipmentEvents } from "@/lib/api/preview";
import RouteMap from "./route-map";

function date(value: string) {
  return new Date(value).toLocaleDateString("en-US", {
    month: "short",
    day: "numeric",
    year: "numeric",
  });
}

export default async function ShipmentPage({
  params,
  searchParams,
}: {
  params: Promise<{ orderId: string }>;
  searchParams: Promise<{ updated?: string }>;
}) {
  const [{ orderId }, { updated }] = await Promise.all([params, searchParams]);
  const shipment = await getShipment(orderId);
  if (!shipment) notFound();
  const events = await getShipmentEvents(shipment.trackingNumber);

  return (
    <main className="detail-layout">
      <section className="details">
        {updated === "1" && <p role="status">Shipment updated.</p>}
        <a className="back" href="/orders">
          ← <span>All orders</span>
        </a>
        <div className="title-row">
          <h1>{shipment.trackingNumber}</h1>
          <span className={`status status-${shipment.status}`}>
            <i aria-hidden="true" />
            {shipment.status.replaceAll("_", " ")}
          </span>
        </div>
        <p className="order-id">
          Order ID <strong>{shipment.id}</strong> · {shipment.customer.name}
        </p>
        <article className="route-card">
          <div className="route-points">
            <span className="point-origin" />
            <span className="point-line" />
            <span className="point-destination" />
          </div>
          <div className="route-names">
            <small>FROM</small>
            <div>{shipment.origin}</div>
            <small>TO</small>
            <div>{shipment.destination}</div>
          </div>
          <span className="carrier">{shipment.carrier}</span>
          <div className="progress-meta">
            <div>
              <small>PROGRESS</small>
              <div>{shipment.progress}%</div>
            </div>
            <div>
              <small>CREATED</small>
              <div>{date(shipment.createdAt)}</div>
            </div>
            <div>
              <small>ETA</small>
              <div className="primary">{date(shipment.estimatedDelivery)}</div>
            </div>
          </div>
          <div className="progress">
            <span style={{ width: `${shipment.progress}%` }} />
          </div>
        </article>
        <div className="meta-grid">
          {[
            ["Tracking #", shipment.trackingNumber],
            ["Customer", shipment.customer.name],
            ["Carrier", shipment.carrier],
            ["Weight/Kg", `${shipment.weight} kg`],
            ["Created", date(shipment.createdAt)],
          ].map(([label, value]) => (
            <article className="meta-card" key={label}>
              <small>{label}</small>
              <div>{value}</div>
            </article>
          ))}
        </div>
        <section className="timeline">
          <h2>Shipment status</h2>
          <ol>
            {[...events].reverse().map((event, index) => (
              <li key={`${event.timestamp}-${index}`}>
                <span className={`timeline-dot ${index === 0 ? "current" : ""}`} />
                <div className="event-heading">
                  <strong>{event.status}</strong>
                  <time>
                    {new Date(event.timestamp).toLocaleString("en-US", {
                      month: "short",
                      day: "numeric",
                      hour: "2-digit",
                      minute: "2-digit",
                      hour12: false,
                    })}
                  </time>
                </div>
                <p className="event-location">⌖ {event.location.name}</p>
                {event.description && <p className="event-description">{event.description}</p>}
              </li>
            ))}
          </ol>
        </section>
      </section>
      <aside className="map-panel">
        <RouteMap
          origin={shipment.customer.coords}
          destination={shipment.receiver.coords}
          current={shipment.currentCoords}
          originLabel={shipment.origin}
          destinationLabel={shipment.destination}
          carrier={shipment.carrier}
        />
        <div className="telemetry">
          <small>
            <i aria-hidden="true" /> LIVE TELEMETRY
          </small>
          <strong>{shipment.carrier}</strong>
          <span>
            {shipment.currentCoords.lat.toFixed(2)}°, {shipment.currentCoords.lng.toFixed(2)}°
          </span>
        </div>
        <div className="geo-route">⌗ GEO ROUTE</div>
      </aside>
    </main>
  );
}
