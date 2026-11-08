import { Activity, ArrowRight, Boxes, Globe2, Truck } from "lucide-react";

import Link from "next/link";
import { getHomeData } from "@/features/shipments/home-data";
import TrackingSearch from "@/features/shipments/TrackingSearch";
import styles from "@/features/shipments/home.module.css";

export const dynamic = "force-dynamic";

const metrics = [
  { label: "Active shipments", icon: Truck },
  { label: "Delivered", icon: Boxes },
  { label: "On-time rate", icon: Activity },
  { label: "Provinces served", icon: Globe2 },
];

export default async function HomePage() {
  const { analytics, shipments } = await getHomeData();
  const values = [
    analytics ? analytics.total - analytics.delivered : "—",
    analytics?.delivered ?? "—",
    "99.9%",
    "77",
  ];

  return (
    <main>
      <section className={styles.hero}>
        <div className={styles.heroInner}>
          <div className={styles.heroCopy}>
            <span className={styles.liveLabel}>
              <i className={styles.liveDot} /> Live ops console
            </span>
            <h1 className={styles.heroTitle}>
              <span className={styles.heroFirstLine}>
                <span>Move fast</span>
                <Truck className={styles.heroTruck} aria-hidden="true" />
              </span>
              <span>Break nothing.</span>
            </h1>
            <p className={styles.heroText}>
              Trace every parcel from origin warehouse to doorstep. Realtime telemetry, port-side
              timelines, and exception alerts in one console.
            </p>
            <TrackingSearch />
          </div>
        </div>
      </section>
      <section className={styles.metrics} aria-label="Key metrics">
        <div className={`${styles.sectionInner} ${styles.metricGrid}`}>
          {metrics.map((metric, index) => {
            const Icon = metric.icon;
            return (
              <div className={styles.metric} key={metric.label}>
                <span className={styles.metricIcon} aria-hidden="true">
                  <Icon size={20} />
                </span>
                <div>
                  <div className={styles.metricValue}>{values[index]}</div>
                  <div className={styles.metricLabel}>{metric.label}</div>
                </div>
              </div>
            );
          })}
        </div>
      </section>
      <section className={`${styles.sectionInner} ${styles.recent}`}>
        <div className={styles.recentHeading}>
          <div>
            <h2>Recent shipments</h2>
            <p>Latest activity from the fleet.</p>
          </div>
          <Link className={styles.viewAll} href="/orders">
            View all <ArrowRight size={16} aria-hidden="true" />
          </Link>
        </div>
        <div className={styles.shipmentGrid}>
          {shipments === null ? (
            <p className={styles.feedback} role="status">
              Recent shipments are unavailable right now.
            </p>
          ) : shipments.length === 0 ? (
            <p className={styles.feedback}>No active shipments yet.</p>
          ) : (
            shipments.map((shipment) => (
              <Link
                className={styles.shipment}
                href={`/orders/${encodeURIComponent(shipment.id)}`}
                key={shipment.id}
              >
                <div className={styles.shipmentTop}>
                  <span className={styles.shipmentId}>{shipment.id}</span>
                  <span className={styles.status}>{shipment.status.replaceAll("_", " ")}</span>
                </div>
                <div className={styles.tracking}>{shipment.trackingNumber}</div>
                <div className={styles.customer}>{shipment.customer.name}</div>
                <div className={styles.route}>
                  <span>{shipment.origin}</span>
                  <ArrowRight className={styles.routeArrow} size={14} aria-hidden="true" />
                  <span>{shipment.destination}</span>
                </div>
                <div className={styles.progress} aria-label={`${shipment.progress}% complete`}>
                  <span style={{ width: `${shipment.progress}%` }} />
                </div>
              </Link>
            ))
          )}
        </div>
        {analytics === null && (
          <p className={styles.feedback} role="status">
            Dashboard metrics are unavailable right now.
          </p>
        )}
      </section>
    </main>
  );
}
