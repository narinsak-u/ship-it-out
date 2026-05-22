import { useQuery, useMutation, useQueryClient } from "@tanstack/vue-query";
import { fetchDrivers, assignDriverToOrder } from "@/lib/api/carriers";

export function useDrivers() {
  return useQuery({
    queryKey: ["drivers"],
    queryFn: fetchDrivers,
  });
}

export function useAssignDriver() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ driverId, orderId }: { driverId: string; orderId: string }) =>
      assignDriverToOrder(driverId, orderId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["drivers"] });
      queryClient.invalidateQueries({ queryKey: ["deliveries"] });
    },
  });
}
