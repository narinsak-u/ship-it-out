import { toast } from "vue-sonner";
import { useQuery, useMutation, useQueryClient } from "@tanstack/vue-query";
import { fetchHubs, createHub, updateHub, deleteHub } from "@/lib/api/carriers";
import type { Hub } from "@/lib/carriers";

export function useHubs() {
  return useQuery({
    queryKey: ["hubs"],
    queryFn: fetchHubs,
  });
}

export function useCreateHub() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: Omit<Hub, "id">) => createHub(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["hubs"] });
    },
  });
}

export function useUpdateHub() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<Hub> }) => updateHub(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["hubs"] });
    },
  });
}

export function useDeleteHub() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => deleteHub(id),
    onSuccess: () => {
      toast.success("Hub deleted");
      queryClient.invalidateQueries({ queryKey: ["hubs"] });
    },
  });
}
