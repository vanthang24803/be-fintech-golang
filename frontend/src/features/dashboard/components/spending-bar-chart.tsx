"use client";

import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from "recharts";
import { useMonthlyTrend } from "@/features/dashboard/hooks/use-reports";

function shortMonth(monthStr: string) {
  const d = new Date(monthStr + "-01");
  return d.toLocaleDateString("vi-VN", { day: "2-digit", month: "2-digit" });
}

function fmtK(v: number) {
  if (v >= 1_000_000) return (v / 1_000_000).toFixed(1) + "M";
  if (v >= 1_000) return Math.round(v / 1_000) + "k";
  return String(v);
}

export function SpendingBarChart() {
  const { data: trend, isLoading } = useMonthlyTrend(7);

  const chartData = (trend ?? []).map((item) => ({
    label: shortMonth(item.month),
    income: item.income,
    expense: item.expense,
  }));

  return (
    <div className="bg-white rounded-2xl border border-gray-100 shadow-sm p-5">
      <div className="flex items-center justify-between mb-5">
        <h3 className="text-sm font-semibold text-gray-800">Biểu đồ 7 ngày qua</h3>
        <span className="text-xs font-semibold text-gray-400 tracking-widest uppercase">
          Xu hướng thu chi
        </span>
      </div>

      {isLoading ? (
        <div className="h-52 flex items-center justify-center">
          <div className="w-8 h-8 border-2 border-[#001BB7] border-t-transparent rounded-full animate-spin" />
        </div>
      ) : (
        <ResponsiveContainer width="100%" height={210}>
          <BarChart data={chartData} barGap={4} barCategoryGap="30%">
            <CartesianGrid strokeDasharray="3 3" stroke="#f3f4f6" vertical={false} />
            <XAxis
              dataKey="label"
              tick={{ fontSize: 11, fill: "#9ca3af" }}
              axisLine={false}
              tickLine={false}
            />
            <YAxis
              tickFormatter={fmtK}
              tick={{ fontSize: 11, fill: "#9ca3af" }}
              axisLine={false}
              tickLine={false}
              width={36}
            />
            <Tooltip
              formatter={(v: unknown, name: unknown) => [
                Number(v).toLocaleString("vi-VN") + " đ",
                name === "income" ? "Thu nhập" : "Chi tiêu" as string,
              ]}
              contentStyle={{ borderRadius: 10, border: "none", boxShadow: "0 4px 20px rgba(0,0,0,.08)", fontSize: 12 }}
            />
            <Bar dataKey="income" fill="#22c55e" radius={[4, 4, 0, 0]} maxBarSize={20} />
            <Bar dataKey="expense" fill="#f87171" radius={[4, 4, 0, 0]} maxBarSize={20} />
          </BarChart>
        </ResponsiveContainer>
      )}
    </div>
  );
}
