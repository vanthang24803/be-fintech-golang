import { useQuery } from "@tanstack/react-query";
import { apiRequest } from "@/lib/api/client";
import { MonthlyTrendItem, CategorySummaryItem } from "@/lib/api/types";

export function useMonthlyTrend(months: number = 6) {
  return useQuery({
    queryKey: ["reports", "monthly-trend", months],
    queryFn: () =>
      apiRequest<MonthlyTrendItem[]>("reports/monthly-trend", {
        body: JSON.stringify({ months }),
      }),
  });
}

export function useCategorySummary(startDate?: string, endDate?: string) {
  return useQuery({
    queryKey: ["reports", "category-summary", startDate, endDate],
    queryFn: () =>
      apiRequest<CategorySummaryItem[]>("reports/category-summary", {
        body: JSON.stringify({ start_date: startDate, end_date: endDate }),
      }),
  });
}
