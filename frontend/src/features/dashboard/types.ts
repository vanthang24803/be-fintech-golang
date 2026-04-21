import type { LucideIcon } from "lucide-react";

export type DashboardMetric = {
  title: string;
  value: string;
  detail: string;
  icon: LucideIcon;
};

export type DashboardActivityRow = {
  category: string;
  amount: string;
  source: string;
  when: string;
};

