import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiRequest } from "@/lib/api/client";
import { SourcePayment } from "@/lib/api/types";

export function useSources() {
  return useQuery({
    queryKey: ["sources", "list"],
    queryFn: () => apiRequest<SourcePayment[]>("sources/list"),
  });
}

export function useCreateSource() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: Partial<SourcePayment>) =>
      apiRequest<SourcePayment>("sources/create", {
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["sources"] });
    },
  });
}

export function useUpdateSource(id: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: Partial<SourcePayment>) =>
      apiRequest<SourcePayment>(`sources/update/${id}`, {
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["sources"] });
    },
  });
}

export function useDeleteSource() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) =>
      apiRequest(`sources/delete/${id}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["sources"] });
    },
  });
}
