"use client";

import { useMemo, useState, type FormEvent } from "react";
import { getDistrictNames, getProvinceName, getSubDistrictNames } from "thai-data";
import { useRouter } from "next/navigation";
import type { Shipment } from "./types";
import { saveShipment } from "./browser-api";
import { useAuthStore } from "@/features/auth/auth-store";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

type AddressDraft = Omit<Shipment["customer"], "coords">;
const emptyAddress: AddressDraft = {
  name: "",
  zipcode: "",
  subDistrict: "",
  district: "",
  province: "",
};

function addressDraft(contact?: Shipment["customer"]): AddressDraft {
  return contact
    ? {
        name: contact.name,
        zipcode: contact.zipcode,
        subDistrict: contact.subDistrict,
        district: contact.district,
        province: contact.province,
      }
    : { ...emptyAddress };
}

function AddressFields({
  title,
  value,
  onChange,
  errors,
}: {
  title: string;
  value: AddressDraft;
  onChange: (value: AddressDraft) => void;
  errors: Record<string, string>;
}) {
  const subDistricts = useMemo(
    () => (value.zipcode.length === 5 ? getSubDistrictNames(value.zipcode) : []),
    [value.zipcode],
  );
  function patch(field: keyof AddressDraft, next: string) {
    if (field === "zipcode")
      onChange({ ...value, zipcode: next, subDistrict: "", district: "", province: "" });
    else onChange({ ...value, [field]: next });
  }
  function selectSubDistrict(subDistrict: string) {
    onChange({
      ...value,
      subDistrict,
      district: getDistrictNames(value.zipcode)[0] ?? "",
      province: getProvinceName(value.zipcode) ?? "",
    });
  }

  return (
    <fieldset className="rounded-xl border border-border bg-card p-5 sm:p-6">
      <legend className="px-2 text-lg font-semibold">{title}</legend>
      <div className="grid gap-4 sm:grid-cols-2">
        <label>
          Name
          <Input
            className="mt-1"
            aria-label={`${title} name`}
            value={value.name}
            onChange={(event) => patch("name", event.target.value)}
            required
          />
        </label>
        {errors.name && <p role="alert">{errors.name}</p>}
        <label>
          Zipcode
          <Input
            className="mt-1"
            aria-label={`${title} zipcode`}
            inputMode="numeric"
            maxLength={5}
            value={value.zipcode}
            onChange={(event) => patch("zipcode", event.target.value)}
            required
          />
        </label>
        {errors.zipcode && <p role="alert">{errors.zipcode}</p>}
        <label>
          Sub-district
          <select
            className="mt-1 flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50"
            aria-label={`${title} sub-district`}
            value={value.subDistrict}
            onChange={(event) => selectSubDistrict(event.target.value)}
            disabled={subDistricts.length === 0}
            required
          >
            <option value="">
              {value.zipcode.length === 5 ? "Select sub-district" : "Enter zipcode first"}
            </option>
            {subDistricts.map((subDistrict) => (
              <option key={subDistrict} value={subDistrict}>
                {subDistrict}
              </option>
            ))}
          </select>
        </label>
        {errors.subDistrict && <p role="alert">{errors.subDistrict}</p>}
        <label>
          District
          <Input
            className="mt-1"
            aria-label={`${title} district`}
            value={value.district}
            readOnly
          />
        </label>
        <label>
          Province
          <Input
            className="mt-1"
            aria-label={`${title} province`}
            value={value.province}
            readOnly
          />
        </label>
      </div>
    </fieldset>
  );
}

export function ShipmentForm({ initial }: { initial?: Shipment }) {
  const router = useRouter();
  const isStaff = useAuthStore((state) => state.isStaff);
  const editing = !!initial;
  const [sender, setSender] = useState(() => addressDraft(initial?.customer));
  const [receiver, setReceiver] = useState(() => addressDraft(initial?.receiver));
  const [weight, setWeight] = useState(initial?.weight ?? 0);
  const [items, setItems] = useState(initial?.items ?? 1);
  const [estimatedDelivery, setEstimatedDelivery] = useState(initial?.estimatedDeliveryRaw ?? "");
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [error, setError] = useState("");
  const [pending, setPending] = useState(false);

  function validate() {
    const next: Record<string, string> = {};
    for (const [label, address] of [
      ["sender", sender],
      ["receiver", receiver],
    ] as const) {
      for (const field of ["name", "zipcode", "subDistrict", "district", "province"] as const) {
        if (!address[field].trim()) next[`${label}.${field}`] = "Required";
      }
      if (address.zipcode && !/^\d{5}$/.test(address.zipcode))
        next[`${label}.zipcode`] = "Enter a 5-digit Thai zipcode";
      else if (address.zipcode && getSubDistrictNames(address.zipcode).length === 0)
        next[`${label}.zipcode`] = "No Thai address matches this zipcode";
    }
    if (!Number.isFinite(weight) || weight <= 0) next.weight = "Enter a weight greater than 0";
    if (!Number.isInteger(items) || items < 1) next.items = "Enter at least 1 item";
    if (editing && !estimatedDelivery) next.estimatedDelivery = "Required";
    setErrors(next);
    return Object.keys(next).length === 0;
  }

  async function geocode(address: AddressDraft) {
    const response = await fetch("/api/geocode", {
      method: "POST",
      headers: { "content-type": "application/json" },
      body: JSON.stringify(address),
    });
    const result = await response.json();
    if (!response.ok) throw new Error(result.error || `Address lookup failed (${response.status})`);
    return result.data as { lat: number; lng: number };
  }

  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!isStaff || pending || !validate()) return;
    setPending(true);
    setError("");
    try {
      const [senderCoords, receiverCoords] = await Promise.all([
        geocode(sender),
        geocode(receiver),
      ]);
      await saveShipment(initial?.id, {
        customer: { ...sender, coords: senderCoords },
        receiver: { ...receiver, coords: receiverCoords },
        carrier: initial?.carrier ?? "Thun-u-der Express",
        weight,
        items,
        ...(editing ? { estimatedDelivery: new Date(estimatedDelivery).toISOString() } : {}),
      });
      window.location.assign(
        initial ? `/orders/${encodeURIComponent(initial.id)}?updated=1` : "/orders?created=1",
      );
    } catch (cause) {
      setError(cause instanceof Error ? cause.message : "Shipment could not be saved.");
    } finally {
      setPending(false);
    }
  }

  if (!isStaff)
    return (
      <p
        className="rounded-lg border border-destructive/40 bg-destructive/10 p-4 text-sm text-destructive"
        role="status"
      >
        Staff access required to {editing ? "edit" : "create"} shipments. Sign in with a staff
        account.
      </p>
    );

  return (
    <form className="space-y-6" onSubmit={submit} noValidate>
      <AddressFields
        title="Sender Info"
        value={sender}
        onChange={setSender}
        errors={Object.fromEntries(
          Object.entries(errors)
            .filter(([key]) => key.startsWith("sender."))
            .map(([key, message]) => [key.slice(7), message]),
        )}
      />
      <AddressFields
        title="Receiver Info"
        value={receiver}
        onChange={setReceiver}
        errors={Object.fromEntries(
          Object.entries(errors)
            .filter(([key]) => key.startsWith("receiver."))
            .map(([key, message]) => [key.slice(9), message]),
        )}
      />
      <fieldset className="rounded-xl border border-border bg-card p-5 sm:p-6">
        <legend className="px-2 text-lg font-semibold">Parcel Info</legend>
        <div className="grid gap-4 sm:grid-cols-2">
          <label>
            Carrier
            <Input className="mt-1" value={initial?.carrier ?? "Thun-u-der Express"} readOnly />
          </label>
          <label>
            Weight/Kg
            <Input
              className="mt-1"
              aria-label="Weight/Kg"
              type="number"
              min="0.1"
              step="0.1"
              value={weight || ""}
              onChange={(event) => setWeight(Number(event.target.value))}
              required
            />
          </label>
          {errors.weight && <p role="alert">{errors.weight}</p>}
          <label>
            Items
            <Input
              className="mt-1"
              aria-label="Items"
              type="number"
              min="1"
              step="1"
              value={items || ""}
              onChange={(event) => setItems(Number(event.target.value))}
              required
            />
          </label>
          {errors.items && <p role="alert">{errors.items}</p>}
          {editing && (
            <label>
              Estimated Delivery
              <Input
                className="mt-1"
                aria-label="Estimated Delivery"
                type="date"
                value={estimatedDelivery.slice(0, 10)}
                onChange={(event) => setEstimatedDelivery(event.target.value)}
                required
              />
            </label>
          )}
          {errors.estimatedDelivery && <p role="alert">{errors.estimatedDelivery}</p>}
        </div>
      </fieldset>
      {error && (
        <p
          className="rounded-lg border border-destructive/40 bg-destructive/10 p-4 text-sm text-destructive"
          role="alert"
        >
          {error}
        </p>
      )}
      <div className="flex justify-end gap-3">
        <Button
          type="button"
          variant="outline"
          onClick={() =>
            router.push(initial ? `/orders/${encodeURIComponent(initial.id)}` : "/orders")
          }
        >
          Cancel
        </Button>
        <Button type="submit" disabled={pending}>
          {pending ? "Resolving addresses…" : editing ? "Save Changes" : "Create Shipment"}
        </Button>
      </div>
    </form>
  );
}
