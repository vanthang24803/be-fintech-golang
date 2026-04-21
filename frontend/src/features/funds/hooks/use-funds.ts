import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiRequest } from "@/lib/api/client";

export interface Fund {
  id: string;
  user_id: string;
  name: string;
  description: string | null;
  target_amount: number;
  balance: number;
  currency: string;
  created_at: string;
  updated_at: string;
}

export function useFunds() {
  return useQuery({
    queryKey: ["funds", "list"],
    queryFn: () => apiRequest<Fund[]>("funds/list"),
  });
}

export function useCreateFund() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: any) =>
      apiRequest<Fund>("funds/create", {
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["funds"] });
    },
  });
}
