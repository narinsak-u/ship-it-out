import { useQuery } from "@tanstack/vue-query";
import { fetchTimeSeries } from "@/lib/api/analytics";

export function useTimeSeries() {
  return useQuery({
    queryKey: ["analytics", "timeseries"],
    queryFn: fetchTimeSeries,
  });
}
