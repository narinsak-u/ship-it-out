import { HubsPanel } from "@/features/operations/HubsPanel";
import styles from "@/features/operations/CarrierShell.module.css";
import type { Hub } from "@/lib/hubs";

export const dynamic = "force-dynamic";

async function getHubs(): Promise<Hub[]> {
  const base = process.env.API_BASE_URL;
  if (!base) throw new Error("API_BASE_URL is required");

  const response = await fetch(`${base.replace(/\/$/, "")}/api/hubs`, { cache: "no-store" });
  if (!response.ok) throw new Error(`Hub API returned ${response.status}`);
  const result = await response.json();
  return result.data as Hub[];
}

export default async function CarriersPage() {
  let hubs: Hub[] = [];
  let error = false;

  try {
    hubs = await getHubs();
  } catch {
    error = true;
  }

  return (
    <main>
      <section className={styles.hero}>
        <div className={styles.heroInner}>
          <span className={styles.eyebrow}>/ carriers</span>
          <h1>Carrier operations</h1>
          <p>Manage carrier hubs and monitor active deliveries across the fleet.</p>
        </div>
      </section>
      <section className={styles.content}>
        <div className={styles.tabs} aria-label="Carrier operations">
          <span className={styles.tab} aria-current="page">
            Hubs
          </span>
        </div>
        <div className={styles.panel}>
          <HubsPanel hubs={hubs} initialError={error} />
        </div>
      </section>
    </main>
  );
}
