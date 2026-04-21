import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiRequest } from "@/lib/api/client";

export interface SavingsGoal {
  id: string;
  user_id: string;
  name: string;
  target_amount: number;
  saved_amount: number;
  progress_percent: number;
  deadline: string | null;
  is_reached: boolean;
  created_at: string;
}

export function useGoals() {
  return useQuery({
    queryKey: ["goals", "list"],
    queryFn: () => apiRequest<SavingsGoal[]>("goals/list"),
  });
}

export function useCreateGoal() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: any) =>
      apiRequest<SavingsGoal>("goals/create", {
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["goals"] });
    },
  });
}

export function useContributeGoal() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: { goal_id: string; fund_id: string; amount: number }) =>
      apiRequest<SavingsGoal>("goals/contribute", {
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["goals"] });
      queryClient.invalidateQueries({ queryKey: ["funds"] });
    },
  });
}
