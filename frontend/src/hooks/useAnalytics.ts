import { useQuery } from "@tanstack/vue-query";
import { fetchAnalytics } from "@/lib/api/analytics";

export function useAnalytics() {
  return useQuery({
    queryKey: ["analytics"],
    queryFn: fetchAnalytics,
  });
}
