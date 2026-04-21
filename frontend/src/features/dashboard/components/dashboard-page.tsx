"use client";

import { StatsRow } from "./stats-row";
import { SpendingBarChart } from "./spending-bar-chart";
import { SpendingDonutChart } from "./spending-donut-chart";
import { RecentTransactions } from "./recent-transactions";

export function DashboardPage() {
  return (
    <div className="p-5 space-y-5">
      <StatsRow />

      <div className="grid gap-5 lg:grid-cols-[1.1fr_0.9fr]">
        <SpendingBarChart />
        <SpendingDonutChart />
      </div>

      <RecentTransactions />
    </div>
  );
}
