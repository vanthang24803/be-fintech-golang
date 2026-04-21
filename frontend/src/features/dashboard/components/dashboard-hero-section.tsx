"use client";

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { useMonthlyTrend } from "@/features/dashboard/hooks/use-reports";
import { TrendingDown, TrendingUp, Wallet, Loader2 } from "lucide-react";

export function DashboardHeroSection() {
  const { data: trend, isLoading } = useMonthlyTrend(6);

  if (isLoading) {
    return (
      <div className="flex h-48 items-center justify-center">
        <Loader2 className="size-8 animate-spin text-primary" />
      </div>
    );
  }

  const currentMonth = trend?.[trend.length - 1];
  const lastMonth = trend?.[trend.length - 2];

  const metrics = [
    {
      title: "Net income",
      value: currentMonth ? `${currentMonth.net_profit.toLocaleString()} VND` : "0 VND",
      detail: currentMonth?.month || "No data",
      icon: (currentMonth?.net_profit || 0) >= 0 ? TrendingUp : TrendingDown,
    },
    {
      title: "Total spending",
      value: currentMonth ? `${currentMonth.expense.toLocaleString()} VND` : "0 VND",
      detail: "This month",
      icon: TrendingDown,
    },
    {
      title: "Total income",
      value: currentMonth ? `${currentMonth.income.toLocaleString()} VND` : "0 VND",
      detail: "This month",
      icon: Wallet,
    },
  ];

  return (
    <section className="grid gap-4 xl:grid-cols-[1.2fr_0.8fr]">
      <Card className="overflow-hidden border-border/70 bg-card/95">
        <CardHeader className="space-y-4">
          <div className="space-y-2">
            <p className="text-xs uppercase tracking-[0.3em] text-primary">Financial health</p>
            <CardTitle className="text-3xl leading-tight md:text-4xl">
              Monthly Overview
            </CardTitle>
          </div>
          <CardDescription className="max-w-2xl text-sm leading-6">
            Real-time insights into your spending and income trends.
          </CardDescription>
        </CardHeader>
        <CardContent className="grid gap-4 md:grid-cols-3">
          {metrics.map((metric) => (
            <div
              key={metric.title}
              className="rounded-2xl border border-border/70 bg-background/80 p-4"
            >
              <metric.icon className="size-5 text-primary" />
              <p className="mt-4 text-sm text-muted-foreground">{metric.title}</p>
              <p className="mt-1 text-2xl font-semibold">{metric.value}</p>
              <p className="mt-2 text-xs text-muted-foreground">{metric.detail}</p>
            </div>
          ))}
        </CardContent>
      </Card>
      <Card className="border-border/70 bg-sidebar text-sidebar-foreground">
        <CardHeader>
          <CardTitle className="text-sidebar-foreground">Control Center</CardTitle>
          <CardDescription className="text-sidebar-foreground/70">
            Keep track of your financial goals and budgets.
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-3 text-sm text-sidebar-foreground/80">
          <p>Quick Summary:</p>
          <ul className="space-y-2">
            <li>Income is {((currentMonth?.income || 0) > (lastMonth?.income || 0)) ? "up" : "down"} compared to last month.</li>
            <li>Spending is {((currentMonth?.expense || 0) > (lastMonth?.expense || 0)) ? "higher" : "lower"} than last month.</li>
          </ul>
        </CardContent>
      </Card>
    </section>
  );
}

