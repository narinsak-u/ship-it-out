"use client";

import { useState, type FormEvent } from "react";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { saveHub } from "./api";
import type { Hub, HubStatus } from "./types";

const provinces = [
  "Amnat Charoen",
  "Bueng Kan",
  "Chaiyaphum",
  "Kalasin",
  "Loei",
  "Maha Sarakham",
  "Mukdahan",
  "Nakhon Phanom",
  "Nong Bua Lamphu",
  "Nong Khai",
  "Roi Et",
  "Sakon Nakhon",
  "Sisaket",
  "Surin",
  "Yasothon",
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
      if (!hub) {
        const geocodeResponse = await fetch("/api/geocode", {
          method: "POST",
          headers: { "content-type": "application/json" },
          body: JSON.stringify({ subDistrict: "", district: "", province: location }),
        });
        const geocode = await geocodeResponse.json();
        if (!geocodeResponse.ok) throw new Error(geocode.error || "Could not resolve location.");
        coords = geocode.data;
      }

      await saveHub(
        {
          name: name.trim(),
          carrierId: "THUN",
          address: location.trim(),
          coords: coords!,
          capacity,
          currentUtilization: utilization,
          status,
        },
        hub?.id,
      );
      await onSaved();
    } catch (cause) {
      setError(cause instanceof Error ? cause.message : "Failed to save hub.");
    } finally {
      setPending(false);
    }
  }

  return (
    <Dialog
      open
      onOpenChange={(isOpen) => {
        if (!isOpen) onClose();
      }}
    >
      <DialogContent className="max-h-[calc(100dvh-2rem)] overflow-y-auto sm:max-w-lg">
        <DialogHeader>
          <DialogTitle>{hub ? "Edit Hub" : "Add Hub"}</DialogTitle>
          <DialogDescription>Update hub capacity, location, and status.</DialogDescription>
        </DialogHeader>
        <form className="grid gap-4" onSubmit={submit}>
          <div className="grid gap-2">
            <Label htmlFor="hub-name">Name</Label>
            <Input
              id="hub-name"
              required
              value={name}
              onChange={(event) => setName(event.target.value)}
            />
          </div>
          <div className="grid gap-2">
            <Label htmlFor="hub-location">Province / address</Label>
            {hub ? (
              <Input
                id="hub-location"
                aria-label="Province"
                required
                value={location}
                onChange={(event) => setLocation(event.target.value)}
              />
            ) : (
              <Select value={location} onValueChange={setLocation} required>
                <SelectTrigger id="hub-location" aria-label="Province">
                  <SelectValue placeholder="Select province…" />
                </SelectTrigger>
                <SelectContent>
                  {provinces.map((province) => (
                    <SelectItem key={province} value={province}>
                      {province}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            )}
          </div>
          <div className="grid gap-2">
            <Label htmlFor="hub-capacity">Capacity (units)</Label>
            <Input
              id="hub-capacity"
              type="number"
              min="1"
              required
              value={capacity}
              onChange={(event) => setCapacity(Number(event.target.value))}
            />
          </div>
          <div className="grid gap-2">
            <Label htmlFor="hub-utilization">Current Utilization</Label>
            <Input
              id="hub-utilization"
              type="number"
              min="0"
              max={capacity}
              required
              value={utilization}
              onChange={(event) => setUtilization(Number(event.target.value))}
            />
          </div>
          {utilization > capacity && (
            <p className="text-sm text-destructive" role="alert">
              Cannot exceed capacity ({capacity})
            </p>
          )}
          <div className="grid gap-2">
            <Label htmlFor="hub-status">Status</Label>
            <Select value={status} onValueChange={(value) => setStatus(value as HubStatus)}>
              <SelectTrigger id="hub-status">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {statuses.map((value) => (
                  <SelectItem key={value} value={value}>
                    {value}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
          {error && (
            <p className="text-sm text-destructive" role="alert">
              {error}
            </p>
          )}
          <div className="flex justify-end gap-2">
            <Button type="button" variant="outline" onClick={onClose} disabled={pending}>
              Cancel
            </Button>
            <Button
              type="submit"
              disabled={
                !name.trim() ||
                !location.trim() ||
                capacity <= 0 ||
                utilization < 0 ||
                utilization > capacity ||
                pending
              }
            >
              {pending ? "Resolving location…" : hub ? "Update Hub" : "Create Hub"}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}
