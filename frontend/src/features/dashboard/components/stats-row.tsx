"use client";

import { TrendingDown, TrendingUp, Wallet } from "lucide-react";
import { useMonthlyTrend } from "@/features/dashboard/hooks/use-reports";
import { useSources } from "@/features/sources/hooks/use-sources";

function fmt(n: number) {
  return n.toLocaleString("vi-VN") + " đ";
}

export function StatsRow() {
  const { data: trend, isLoading: trendLoading } = useMonthlyTrend(6);
  const { data: sources, isLoading: sourcesLoading } = useSources();

  const isLoading = trendLoading || sourcesLoading;

  if (isLoading) {
    return (
      <div className="grid grid-cols-3 gap-4">
        {[0, 1, 2].map((i) => (
          <div key={i} className="h-24 rounded-2xl bg-gray-100 animate-pulse" />
        ))}
      </div>
    );
  }

  const totalBalance = (sources ?? []).reduce((sum, s) => sum + s.balance, 0);
  const cur = trend?.[trend.length - 1];
  const income = cur?.income ?? 0;
  const expense = cur?.expense ?? 0;

  return (
    <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
      {/* Tổng số dư */}
      <div
        className="relative rounded-2xl p-5 text-white overflow-hidden"
        style={{ background: "linear-gradient(135deg,#00137a 0%,#001BB7 100%)" }}
      >
        <div className="absolute right-4 top-1/2 -translate-y-1/2 opacity-20">
          <Wallet className="w-16 h-16" />
        </div>
        <div className="flex items-center gap-1.5 mb-3">
          <Wallet className="w-4 h-4 opacity-80" />
          <span className="text-xs font-medium opacity-80">Tổng số dư</span>
        </div>
        <p className="text-2xl font-bold">{fmt(totalBalance)}</p>
      </div>

      {/* Thu nhập */}
      <div className="bg-white rounded-2xl p-5 border border-gray-100 shadow-sm">
        <div className="flex items-center justify-between mb-3">
          <span className="text-xs font-semibold text-gray-400 tracking-widest uppercase">Thu nhập</span>
          <div className="w-8 h-8 rounded-full bg-green-50 flex items-center justify-center">
            <TrendingUp className="w-4 h-4 text-green-500" />
          </div>
        </div>
        <p className="text-xl font-bold text-gray-900">{fmt(income)}</p>
      </div>

      {/* Chi tiêu */}
      <div className="bg-white rounded-2xl p-5 border border-gray-100 shadow-sm">
        <div className="flex items-center justify-between mb-3">
          <span className="text-xs font-semibold text-gray-400 tracking-widest uppercase">Chi tiêu</span>
          <div className="w-8 h-8 rounded-full bg-red-50 flex items-center justify-center">
            <TrendingDown className="w-4 h-4 text-red-400" />
          </div>
        </div>
        <p className="text-xl font-bold text-gray-900">{fmt(expense)}</p>
      </div>
    </div>
  );
}
