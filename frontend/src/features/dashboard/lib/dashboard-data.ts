import { TrendingDown, TrendingUp, Wallet } from "lucide-react";

import type { DashboardActivityRow, DashboardMetric } from "@/features/dashboard/types";

export const dashboardMetricCards: DashboardMetric[] = [
  {
    title: "Net cash flow",
    value: "+12.4%",
    detail: "Compared to the last 30 days",
    icon: TrendingUp,
  },
  {
    title: "Spending drift",
    value: "-3.1%",
    detail: "Dining and transport are cooling off",
    icon: TrendingDown,
  },
  {
    title: "Liquid reserves",
    value: "$18,420",
    detail: "Across wallet, bank, and emergency fund",
    icon: Wallet,
  },
];

export const dashboardRecentRows: DashboardActivityRow[] = [
  { category: "Groceries", amount: "$214.20", source: "Visa", when: "2h ago" },
  { category: "Subscriptions", amount: "$42.00", source: "Debit", when: "6h ago" },
  { category: "Transport", amount: "$18.40", source: "Wallet", when: "Today" },
  { category: "Salary", amount: "+$2,950.00", source: "Bank", when: "Yesterday" },
];

