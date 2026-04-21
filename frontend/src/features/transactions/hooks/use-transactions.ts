import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiRequest } from "@/lib/api/client";
import { Transaction } from "@/lib/api/types";

export function useTransactions(filters: {
  type?: "income" | "expense";
  category_id?: string;
  source_id?: string;
} = {}) {
  return useQuery({
    queryKey: ["transactions", "list", filters],
    queryFn: () =>
      apiRequest<Transaction[]>("transactions/list", {
        body: JSON.stringify(filters),
      }),
  });
}

export function useCreateTransaction() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: any) =>
      apiRequest<Transaction>("transactions/create", {
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["transactions"] });
      queryClient.invalidateQueries({ queryKey: ["reports"] });
    },
  });
}
