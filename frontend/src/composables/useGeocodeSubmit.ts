import { ref } from "vue";
import { geocodeAddress } from "@/lib/geocode";
import type { ContactInfo } from "@/lib/orders";

interface AddressState {
  name: string;
  zipcode: string;
  subDistrict: string;
  district: string;
  province: string;
}

export function useGeocodeSubmit() {
  const geocoding = ref(false);
  const geocodeErrors = ref<Record<string, string>>({});

  async function geocodeBoth(
    sender: AddressState,
    receiver: AddressState,
  ): Promise<{
    senderCoords: { lat: number; lng: number };
    receiverCoords: { lat: number; lng: number };
  } | null> {
    geocodeErrors.value = {};
    geocoding.value = true;

    try {
      const [senderRes, receiverRes] = await Promise.allSettled([
        geocodeAddress(sender.subDistrict, sender.district, sender.province),
        geocodeAddress(receiver.subDistrict, receiver.district, receiver.province),
      ]);

      if (senderRes.status === "rejected") {
        geocodeErrors.value.sender =
          senderRes.reason instanceof Error
            ? senderRes.reason.message
            : "Could not resolve address.";
      }
      if (receiverRes.status === "rejected") {
        geocodeErrors.value.receiver =
          receiverRes.reason instanceof Error
            ? receiverRes.reason.message
            : "Could not resolve address.";
      }
      if (senderRes.status === "rejected" || receiverRes.status === "rejected") {
        return null;
      }

      return {
        senderCoords: senderRes.value,
        receiverCoords: receiverRes.value,
      };
    } catch (e) {
      const msg = e instanceof Error ? e.message : "Could not resolve address.";
      geocodeErrors.value = { sender: msg, receiver: msg };
      return null;
    } finally {
      geocoding.value = false;
    }
  }

  return { geocoding, geocodeErrors, geocodeBoth };
}
