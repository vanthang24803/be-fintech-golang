"use client";

import { useState, useMemo } from "react";
import { useTransactions } from "@/features/transactions/hooks/use-transactions";
import { getCategoryIcon, getCategoryColor } from "@/lib/category-icon";
import { useCategories } from "@/features/categories/hooks/use-categories";
import { Search, Settings2, Loader2 } from "lucide-react";
import { cn } from "@/lib/utils";

function fmtDate(dateStr: string) {
  const d = new Date(dateStr);
  return d.toLocaleDateString("vi-VN", { day: "2-digit", month: "long", year: "numeric" });
}

function fmtAmount(amount: number, type: string) {
  const sign = type === "income" ? "+" : "-";
  return `${sign}${amount.toLocaleString("vi-VN")}đ`;
}

export function TransactionTable() {
  const [search, setSearch] = useState("");
  const { data: transactions, isLoading } = useTransactions();
  const { data: categories } = useCategories();

  const catMap = useMemo(() => {
    const m: Record<string, string> = {};
    (categories ?? []).forEach((c) => { m[c.id] = c.icon ?? ""; });
    return m;
  }, [categories]);

  const filtered = useMemo(() => {
    if (!search.trim()) return transactions ?? [];
    const q = search.toLowerCase();
    return (transactions ?? []).filter(
      (t) =>
        t.description?.toLowerCase().includes(q) ||
        t.category_name?.toLowerCase().includes(q)
    );
  }, [transactions, search]);

  return (
    <div className="bg-white rounded-2xl border border-gray-100 shadow-sm">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center gap-3 px-6 py-4 border-b border-gray-100">
        <h2 className="text-lg font-semibold text-gray-900 flex-1">Lịch sử giao dịch</h2>
        <div className="relative">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
          <input
            type="text"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Tìm kiếm ghi chú..."
            className="pl-9 pr-4 py-2 text-sm border border-gray-200 rounded-xl bg-gray-50 outline-none focus:border-gray-400 w-60"
          />
        </div>
        <span className="text-sm text-gray-400 whitespace-nowrap">
          Tổng cộng: {filtered.length} giao dịch
        </span>
      </div>

      {/* Table header — desktop only */}
      <div className="hidden md:grid grid-cols-[160px_1fr_1fr_160px_40px] px-6 py-3 border-b border-gray-100">
        {["NGÀY", "HẠNG MỤC", "GHI CHÚ", "SỐ TIỀN", ""].map((h, i) => (
          <span key={i} className="text-[11px] font-semibold text-gray-400 tracking-widest">{h}</span>
        ))}
      </div>

      {/* Rows */}
      {isLoading ? (
        <div className="flex items-center justify-center h-48">
          <Loader2 className="w-6 h-6 animate-spin text-gray-300" />
        </div>
      ) : filtered.length === 0 ? (
        <div className="flex items-center justify-center h-48 text-sm text-gray-400">
          Không có giao dịch nào
        </div>
      ) : (
        <div className="divide-y divide-gray-50">
          {filtered.map((row) => {
            const iconField = catMap[row.category_id] ?? "";
            const Icon = getCategoryIcon(iconField);
            const color = getCategoryColor(iconField);
            const isIncome = row.type === "income";

            return (
              <div
                key={row.id}
                className="grid grid-cols-[1fr_auto] md:grid-cols-[160px_1fr_1fr_160px_40px] items-center px-6 py-4 hover:bg-gray-50 transition-colors"
              >
                {/* Date — desktop */}
                <span className="text-sm text-gray-500 hidden md:block">
                  {fmtDate(row.transaction_date)}
                </span>

                {/* Category icon + name */}
                <div className="flex items-center gap-3">
                  <div
                    className="w-10 h-10 rounded-xl flex items-center justify-center flex-shrink-0"
                    style={{ background: color }}
                  >
                    <Icon className="w-5 h-5 text-white" />
                  </div>
                  <div>
                    <p className="text-sm font-semibold text-gray-800">{row.category_name}</p>
                    {/* Date — mobile */}
                    <p className="text-xs text-gray-400 md:hidden">{fmtDate(row.transaction_date)}</p>
                  </div>
                </div>

                {/* Note */}
                <span className="text-sm text-gray-400 italic hidden md:block truncate pr-4">
                  {row.description || "—"}
                </span>

                {/* Amount */}
                <span className={cn("text-sm font-bold text-right", isIncome ? "text-green-500" : "text-gray-800")}>
                  {fmtAmount(row.amount, row.type)}
                </span>

                {/* Settings */}
                <button className="hidden md:flex items-center justify-center w-8 h-8 rounded-lg hover:bg-gray-100 transition-colors ml-auto">
                  <Settings2 className="w-4 h-4 text-gray-400" />
                </button>
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
}
