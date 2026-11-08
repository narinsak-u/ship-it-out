import type { Order } from "@/lib/orders";
import styles from "./shipment-form-page.module.css";
import { ShipmentForm } from "./ShipmentForm";

export function ShipmentFormPage({ initial }: { initial?: Order }) {
  const editing = !!initial;
  return (
    <main>
      <section className={styles.hero}>
        <div>
          <span>/ orders</span>
          <h1>{editing ? "Edit Shipment" : "Create Shipment"}</h1>
          <p>{editing ? "Update shipment details." : "Register a new shipment in the system."}</p>
        </div>
      </section>
      <section className={styles.content}>
        <ShipmentForm initial={initial} />
      </section>
    </main>
  );
}
