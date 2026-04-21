"use client";

import Link from "next/link";
import {
  Utensils, TrendingUp, Car, ShoppingBag, Home, BookOpen,
  Zap, Heart, MoreHorizontal, CircleDollarSign,
} from "lucide-react";
import { useTransactions } from "@/features/transactions/hooks/use-transactions";

const CATEGORY_CONFIG: Record<string, { icon: React.ReactNode; bg: string }> = {
  "ăn uống":    { icon: <Utensils className="w-4 h-4 text-white" />,          bg: "#f97316" },
  "đầu tư":     { icon: <TrendingUp className="w-4 h-4 text-white" />,         bg: "#22c55e" },
  "di chuyển":  { icon: <Car className="w-4 h-4 text-white" />,                bg: "#001BB7" },
  "mua sắm":    { icon: <ShoppingBag className="w-4 h-4 text-white" />,        bg: "#f59e0b" },
  "nhà cửa":    { icon: <Home className="w-4 h-4 text-white" />,               bg: "#14b8a6" },
  "giáo dục":   { icon: <BookOpen className="w-4 h-4 text-white" />,           bg: "#a855f7" },
  "hóa đơn":    { icon: <Zap className="w-4 h-4 text-white" />,                bg: "#6366f1" },
  "sức khỏe":   { icon: <Heart className="w-4 h-4 text-white" />,              bg: "#ec4899" },
  "giải trí":   { icon: <MoreHorizontal className="w-4 h-4 text-white" />,     bg: "#84cc16" },
};

function getCategoryMeta(name: string) {
  const key = name?.toLowerCase().trim();
  for (const [k, v] of Object.entries(CATEGORY_CONFIG)) {
    if (key?.includes(k)) return v;
  }
  return { icon: <CircleDollarSign className="w-4 h-4 text-white" />, bg: "#9ca3af" };
}

function fmtDate(iso: string) {
  return new Date(iso).toLocaleDateString("vi-VN", { day: "2-digit", month: "2-digit", year: "numeric" });
}

export function RecentTransactions() {
  const { data: transactions, isLoading } = useTransactions();
  const recent = transactions?.slice(0, 5) ?? [];

  return (
    <div className="bg-white rounded-2xl border border-gray-100 shadow-sm p-5">
      <div className="flex items-center justify-between mb-5">
        <h3 className="text-sm font-semibold text-gray-800">Giao dịch gần đây</h3>
        <Link href="/transactions" className="text-xs font-medium" style={{ color: "#001BB7" }}>
          Xem tất cả
        </Link>
      </div>

      {isLoading ? (
        <div className="space-y-4">
          {[0, 1, 2, 3, 4].map((i) => (
            <div key={i} className="flex items-center gap-3">
              <div className="w-10 h-10 rounded-full bg-gray-100 animate-pulse flex-shrink-0" />
              <div className="flex-1 space-y-1.5">
                <div className="h-3 bg-gray-100 rounded animate-pulse w-24" />
                <div className="h-2.5 bg-gray-100 rounded animate-pulse w-16" />
              </div>
              <div className="h-3 bg-gray-100 rounded animate-pulse w-20" />
            </div>
          ))}
        </div>
      ) : recent.length === 0 ? (
        <p className="text-sm text-gray-400 text-center py-8">Chưa có giao dịch nào</p>
      ) : (
        <div className="space-y-4">
          {recent.map((tx) => {
            const meta = getCategoryMeta(tx.category_name);
            const isIncome = tx.type === "income";
            return (
              <div key={tx.id} className="flex items-center gap-3">
                <div
                  className="w-10 h-10 rounded-full flex items-center justify-center flex-shrink-0"
                  style={{ background: meta.bg }}
                >
                  {meta.icon}
                </div>
                <div className="flex-1 min-w-0">
                  <p className="text-sm font-medium text-gray-800 truncate">{tx.category_name}</p>
                  <p className="text-xs text-gray-400">{fmtDate(tx.transaction_date)}</p>
                </div>
                <span
                  className="text-sm font-semibold flex-shrink-0"
                  style={{ color: isIncome ? "#22c55e" : "#374151" }}
                >
                  {isIncome ? "+" : "-"}{tx.amount.toLocaleString("vi-VN")} đ
                </span>
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
}
