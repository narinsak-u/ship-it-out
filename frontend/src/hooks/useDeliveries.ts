import { useQuery, useMutation, useQueryClient } from "@tanstack/vue-query";
import { fetchActiveDeliveries, updateShipmentStatus } from "@/lib/api/orders";
import type { Order, ShipmentStatus } from "@/lib/orders";
import { deliveryKeys, orderKeys, eventKeys } from "@/lib/api/queryKeys";
import { toast } from "vue-sonner";

export function useActiveDeliveries() {
  return useQuery({
    queryKey: deliveryKeys.active(),
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
    onSuccess: (data, variables) => {
      queryClient.setQueryData<Order[]>(deliveryKeys.active(), (old) =>
        old?.map((d) => (d.id === data.id ? data : d)),
      );
      queryClient.setQueryData(orderKeys.detail(data.id), data);
      queryClient.invalidateQueries({ queryKey: deliveryKeys.all });
      queryClient.invalidateQueries({ queryKey: orderKeys.all });
      queryClient.invalidateQueries({ queryKey: eventKeys.all });
      toast.success("Delivery status updated");
    },
    onError: () => {
      toast.error("Failed to update delivery status");
    },
  });
}
