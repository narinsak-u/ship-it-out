import { useQuery } from "@tanstack/vue-query";
import { fetchCarriers } from "@/lib/api/carriers";

export function useCarriers() {
  return useQuery({
    queryKey: ["carriers"],
    queryFn: fetchCarriers,
  });
}
