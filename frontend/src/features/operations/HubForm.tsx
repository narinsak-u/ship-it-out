"use client";

import { useEffect, useRef, useState, type FormEvent } from "react";
import type { Hub, HubStatus } from "@/lib/hubs";
import styles from "./HubsPanel.module.css";

const provinces = [
  "Amnat Charoen", "Bueng Kan", "Chaiyaphum", "Kalasin", "Loei", "Maha Sarakham",
  "Mukdahan", "Nakhon Phanom", "Nong Bua Lamphu", "Nong Khai", "Roi Et", "Sakon Nakhon",
  "Sisaket", "Surin", "Yasothon",
];
const statuses: HubStatus[] = ["active", "maintenance", "full", "closed"];

export function HubForm({
  hub,
  onClose,
  onSaved,
}: {
  hub: Hub | null;
  onClose: () => void;
  onSaved: () => Promise<void>;
}) {
  const dialogRef = useRef<HTMLDialogElement>(null);
  useEffect(() => {
    if (!dialogRef.current?.open) dialogRef.current?.showModal();
  }, []);

  const [name, setName] = useState(hub?.name ?? "");
  const [location, setLocation] = useState(hub?.address ?? "");
  const [capacity, setCapacity] = useState(hub?.capacity ?? 1000);
  const [utilization, setUtilization] = useState(hub?.currentUtilization ?? 0);
  const [status, setStatus] = useState<HubStatus>(hub?.status ?? "active");
  const [error, setError] = useState("");
  const [pending, setPending] = useState(false);

  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setError("");
    setPending(true);
    try {
      let coords = hub?.coords;
      if (!coords || hub?.address !== location.trim()) {
        const geocodeResponse = await fetch("/api/geocode", {
          method: "POST",
          headers: { "content-type": "application/json" },
          body: JSON.stringify({ subDistrict: "", district: "", province: location }),
        });
        const geocode = await geocodeResponse.json();
        if (!geocodeResponse.ok) throw new Error(geocode.error || "Could not resolve location.");
        coords = geocode.data;
      }

      const body = {
        name: name.trim(), carrierId: "THUN", address: location.trim(), coords,
        capacity, currentUtilization: utilization, status,
      };
      const response = await fetch(hub ? `/api/hubs/${encodeURIComponent(hub.id)}` : "/api/hubs", {
        method: hub ? "PUT" : "POST",
        headers: { "content-type": "application/json" },
        body: JSON.stringify(body),
      });
      const result = await response.json();
      if (!response.ok) throw new Error(result.error || `Failed to save hub (${response.status}).`);
      await onSaved();
    } catch (cause) {
      setError(cause instanceof Error ? cause.message : "Failed to save hub.");
    } finally {
      setPending(false);
    }
  }

  return (
    <dialog ref={dialogRef} aria-labelledby="hub-form-title" className={styles.dialog} onCancel={onClose}>
      <form onSubmit={submit}>
        <h2 id="hub-form-title">{hub ? "Edit Hub" : "Add Hub"}</h2>
        <label>Name<input aria-label="Name" required value={name} onChange={(event) => setName(event.target.value)} /></label>
        <label>Province / address
          {hub ? (
            <input aria-label="Province" required value={location} onChange={(event) => setLocation(event.target.value)} />
          ) : (
            <select aria-label="Province" required value={location} onChange={(event) => setLocation(event.target.value)}>
              <option value="">Select province…</option>
              {provinces.map((province) => <option key={province} value={province}>{province}</option>)}
            </select>
          )}
        </label>
        <label>Capacity (units)<input aria-label="Capacity (units)" type="number" min="1" required value={capacity} onChange={(event) => setCapacity(Number(event.target.value))} /></label>
        <label>Current Utilization<input aria-label="Current Utilization" type="number" min="0" max={capacity} required value={utilization} onChange={(event) => setUtilization(Number(event.target.value))} /></label>
        {utilization > capacity && <p role="alert">Cannot exceed capacity ({capacity})</p>}
        <label>Status<select aria-label="Status" value={status} onChange={(event) => setStatus(event.target.value as HubStatus)}>{statuses.map((value) => <option key={value} value={value}>{value}</option>)}</select></label>
        {error && <p role="alert">{error}</p>}
        <div className={styles.formActions}>
          <button type="button" onClick={onClose} disabled={pending}>Cancel</button>
          <button type="submit" disabled={!name.trim() || !location.trim() || capacity <= 0 || utilization < 0 || utilization > capacity || pending}>{pending ? "Resolving location…" : hub ? "Update Hub" : "Create Hub"}</button>
        </div>
      </form>
    </dialog>
  );
}
