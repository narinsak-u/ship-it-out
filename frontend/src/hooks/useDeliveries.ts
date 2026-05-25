import { useQuery, useMutation, useQueryClient } from "@tanstack/vue-query";
import { fetchActiveDeliveries, updateShipmentStatus } from "@/lib/api/orders";
import type { ShipmentStatus } from "@/lib/orders";

export function useActiveDeliveries() {
  return useQuery({
    queryKey: ["deliveries"],
    queryFn: fetchActiveDeliveries,
    refetchInterval: 15_000,
  });
}

export function useUpdateShipmentStatus() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({
      orderId,
      status,
      hubId,
    }: {
      orderId: string;
      status: ShipmentStatus;
      hubId?: string;
    }) => updateShipmentStatus(orderId, status, hubId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["deliveries"] });
    },
  });
}
