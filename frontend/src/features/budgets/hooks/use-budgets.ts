import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiRequest } from "@/lib/api/client";

export interface Budget {
  id: string;
  user_id: string;
  category_id: string | null;
  category_name: string;
  amount: number;
  period: "monthly" | "weekly" | "custom";
  start_date: string;
  end_date: string | null;
  is_active: boolean;
  current_spending: number;
  remaining_amount: number;
  progress_percent: number;
  created_at: string;
  updated_at: string;
}

export function useBudgets() {
  return useQuery({
    queryKey: ["budgets", "list"],
    queryFn: () => apiRequest<Budget[]>("budgets/list"),
  });
}

export function useCreateBudget() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: { category_id?: number; amount: number; period: string }) =>
      apiRequest<Budget>("budgets/create", { body: JSON.stringify(data) }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["budgets"] });
    },
  });
}

export function useDeleteBudget() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => apiRequest(`budgets/delete/${id}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["budgets"] });
    },
  });
}
