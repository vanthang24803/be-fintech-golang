"use client";

import { useState, useMemo, useEffect, useRef } from "react";
import { useBudgets, useCreateBudget } from "@/features/budgets/hooks/use-budgets";
import { useCategorySummary } from "@/features/dashboard/hooks/use-reports";
import { useCategories } from "@/features/categories/hooks/use-categories";
import { getCategoryIcon, getCategoryColor } from "@/lib/category-icon";
import { Plus, Target, CalendarDays, Loader2, X } from "lucide-react";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { VisuallyHidden } from "@radix-ui/react-visually-hidden";
import { cn } from "@/lib/utils";

function fmt(n: number) {
  return n.toLocaleString("vi-VN") + "đ";
}

// ── Time filter ──────────────────────────────────────────────
type Range = "all" | "week" | "month" | "year";
const RANGES: { key: Range; label: string }[] = [
  { key: "all",   label: "Tất cả" },
  { key: "week",  label: "Tuần này" },
  { key: "month", label: "Tháng này" },
  { key: "year",  label: "Năm này" },
];

function getDateRange(range: Range): { start?: string; end?: string } {
  const now = new Date();
  if (range === "all") return {};
  const end = now.toISOString().slice(0, 10);
  if (range === "week") {
    const d = new Date(now); d.setDate(d.getDate() - d.getDay());
    return { start: d.toISOString().slice(0, 10), end };
  }
  if (range === "month") {
    return { start: `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, "0")}-01`, end };
  }
  return { start: `${now.getFullYear()}-01-01`, end };
}

// ── Budget setup dialog ──────────────────────────────────────
function BudgetSetupDialog() {
  const [open, setOpen] = useState(false);
  const [categoryId, setCategoryId] = useState("");
  const [amount, setAmount] = useState("");
  const { data: categories } = useCategories();
  const { mutate: createBudget, isPending } = useCreateBudget();

  const expenseCategories = (categories ?? []).filter((c) => c.type === "expense");

  function handleSave() {
    if (!amount || !categoryId) return;
    createBudget(
      { category_id: Number(categoryId), amount: Number(amount.replace(/\D/g, "")), period: "monthly" },
      { onSuccess: () => { setOpen(false); setAmount(""); setCategoryId(""); } }
    );
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <button className="flex items-center gap-2 px-4 py-2 rounded-xl text-white text-sm font-medium" style={{ background: "#1a1a1a" }}>
          <Plus className="w-4 h-4" />
          <span className="hidden sm:inline">Thiết lập</span>
        </button>
      </DialogTrigger>
      <DialogContent className="max-w-sm rounded-2xl p-0 overflow-hidden bg-white">
        <DialogHeader className="px-5 pt-5 pb-0">
          <DialogTitle className="text-base font-semibold">Thiết lập ngân sách</DialogTitle>
        </DialogHeader>
        <div className="px-5 pb-5 space-y-4 mt-3">
          <div>
            <p className="text-xs text-gray-500 mb-2">Danh mục</p>
            <div className="grid grid-cols-2 gap-2 max-h-48 overflow-y-auto">
              {expenseCategories.map((c) => {
                const Icon = getCategoryIcon(c.icon ?? "");
                const color = getCategoryColor(c.icon ?? "");
                const selected = categoryId === c.id;
                return (
                  <button
                    key={c.id}
                    type="button"
                    onClick={() => setCategoryId(c.id)}
                    className={cn(
                      "flex items-center gap-2 px-3 py-2 rounded-xl border text-sm transition-colors",
                      selected ? "border-blue-400 bg-blue-50 text-blue-700" : "border-gray-200 hover:bg-gray-50 text-gray-700"
                    )}
                  >
                    <div className="w-7 h-7 rounded-lg flex items-center justify-center flex-shrink-0" style={{ background: color }}>
                      <Icon className="w-3.5 h-3.5 text-white" />
                    </div>
                    <span className="truncate">{c.name}</span>
                  </button>
                );
              })}
            </div>
          </div>
          <div>
            <p className="text-xs text-gray-500 mb-1">Hạn mức (VND)</p>
            <input
              type="text"
              inputMode="numeric"
              value={amount ? Number(amount).toLocaleString("vi-VN") : ""}
              onChange={(e) => setAmount(e.target.value.replace(/\D/g, ""))}
              placeholder="5.000.000"
              className="w-full border border-gray-200 rounded-xl px-3 py-2.5 text-sm text-gray-700 outline-none focus:border-gray-400"
            />
          </div>
          <button
            onClick={handleSave}
            disabled={isPending || !categoryId || !amount}
            className="w-full py-3 rounded-xl text-white text-sm font-semibold disabled:opacity-50 flex items-center justify-center gap-2"
            style={{ background: "#001BB7" }}
          >
            {isPending && <Loader2 className="w-4 h-4 animate-spin" />}
            Lưu ngân sách
          </button>
        </div>
      </DialogContent>
    </Dialog>
  );
}

// ── Budget section ───────────────────────────────────────────
function BudgetSection() {
  const { data: budgets, isLoading } = useBudgets();
  const { data: categories } = useCategories();

  const catIconMap = useMemo(() => {
    const m: Record<string, string> = {};
    (categories ?? []).forEach((c) => { m[String(c.id)] = c.icon ?? ""; });
    return m;
  }, [categories]);

  const active = (budgets ?? []).filter((b) => b.is_active);

  return (
    <div className="bg-white rounded-2xl border border-gray-100 shadow-sm p-5">
      <div className="flex items-center justify-between mb-5">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 rounded-xl bg-blue-50 flex items-center justify-center">
            <Target className="w-5 h-5 text-blue-500" />
          </div>
          <div>
            <h2 className="text-base font-semibold text-gray-900">Ngân sách tháng này</h2>
            <p className="text-xs text-gray-400">Tiến độ chi tiêu so với hạn mức</p>
          </div>
        </div>
        <BudgetSetupDialog />
      </div>

      {isLoading ? (
        <div className="flex justify-center py-8"><Loader2 className="w-5 h-5 animate-spin text-gray-300" /></div>
      ) : active.length === 0 ? (
        <p className="text-sm text-gray-400 text-center py-6">Chưa có ngân sách nào. Nhấn Thiết lập để bắt đầu.</p>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          {active.map((b) => {
            const iconField = catIconMap[b.category_id ?? ""] ?? "";
            const Icon = getCategoryIcon(iconField);
            const color = getCategoryColor(iconField);
            const pct = Math.min(b.progress_percent ?? 0, 100);
            const over = pct >= 90;
            return (
              <div key={b.id} className="border border-gray-100 rounded-xl p-4">
                <div className="flex items-center gap-2 mb-3">
                  <div className="w-8 h-8 rounded-lg flex items-center justify-center flex-shrink-0" style={{ background: color }}>
                    <Icon className="w-4 h-4 text-white" />
                  </div>
                  <span className="text-sm font-semibold text-gray-800">{b.category_name}</span>
                </div>
                <div className="flex items-center justify-between text-xs text-gray-500 mb-1.5">
                  <span className={cn("font-medium", over ? "text-red-500" : "text-gray-600")}>{Math.round(pct)}%</span>
                  <span>{fmt(b.current_spending)} / {fmt(b.amount)}</span>
                </div>
                <div className="h-1.5 bg-gray-100 rounded-full overflow-hidden">
                  <div
                    className="h-full rounded-full transition-all"
                    style={{ width: `${pct}%`, background: over ? "#ef4444" : "#001BB7" }}
                  />
                </div>
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
}

// ── Horizontal bar chart ─────────────────────────────────────
function SpendingBarChart({ startDate, endDate }: { startDate?: string; endDate?: string }) {
  const { data: summary, isLoading } = useCategorySummary(startDate, endDate);
  const { data: categories } = useCategories();
  const [visible, setVisible] = useState(true);
  const [animated, setAnimated] = useState(false);
  const prevKey = useRef("");
  const key = `${startDate}-${endDate}`;

  useEffect(() => {
    if (prevKey.current === key) return;
    prevKey.current = key;
    // fade out → reset bars → fade in → animate bars
    setVisible(false);
    setAnimated(false);
    const t1 = setTimeout(() => setVisible(true), 150);
    const t2 = setTimeout(() => setAnimated(true), 200);
    return () => { clearTimeout(t1); clearTimeout(t2); };
  }, [key]);

  // also animate on first data load
  useEffect(() => {
    if (!isLoading && summary) {
      setAnimated(false);
      const t = setTimeout(() => setAnimated(true), 50);
      return () => clearTimeout(t);
    }
  }, [isLoading]);

  const catMap = useMemo(() => {
    const m: Record<string, string> = {};
    (categories ?? []).forEach((c) => { m[c.id] = c.icon ?? ""; });
    return m;
  }, [categories]);

  const items = useMemo(
    () => (summary ?? []).filter((c) => c.total_amount > 0).sort((a, b) => b.total_amount - a.total_amount),
    [summary]
  );

  const total = items.reduce((s, c) => s + c.total_amount, 0);

  const chartData = items.map((c) => ({
    name: c.category_name,
    value: c.total_amount,
    color: getCategoryColor(catMap[c.category_id] ?? ""),
    icon: catMap[c.category_id] ?? "",
  }));

  if (isLoading) return <div className="flex justify-center py-12"><Loader2 className="w-5 h-5 animate-spin text-gray-300" /></div>;
  if (items.length === 0) return <p className="text-sm text-gray-400 text-center py-8">Chưa có dữ liệu</p>;

  return (
    <div
      className="grid grid-cols-1 lg:grid-cols-2 gap-5 transition-opacity duration-150"
      style={{ opacity: visible ? 1 : 0 }}
    >
      {/* Horizontal bar chart — custom div-based */}
      <div className="bg-white rounded-2xl border border-gray-100 shadow-sm p-5">
        <h3 className="text-base font-semibold text-gray-900 mb-4">Phân tích chi tiêu</h3>
        <div className="space-y-3">
          {chartData.map((item, i) => {
            const Icon = getCategoryIcon(item.icon);
            const pct = total > 0 ? (item.value / total) * 100 : 0;
            return (
              <div key={item.name} className="flex items-center gap-3">
                <div className="flex items-center gap-1.5 w-24 flex-shrink-0 justify-end">
                  <div className="w-5 h-5 rounded-md flex items-center justify-center flex-shrink-0" style={{ background: item.color }}>
                    <Icon className="w-3 h-3 text-white" />
                  </div>
                  <span className="text-xs text-gray-500 truncate">{item.name}</span>
                </div>
                <div className="flex-1 h-5 bg-gray-100 rounded-full overflow-hidden">
                  <div
                    className="h-full rounded-full"
                    style={{
                      width: animated ? `${pct}%` : "0%",
                      background: item.color,
                      transition: `width 0.5s cubic-bezier(0.4,0,0.2,1) ${i * 40}ms`,
                    }}
                  />
                </div>
              </div>
            );
          })}
        </div>
      </div>

      {/* Summary + Stats */}
      <div className="flex flex-col gap-5">
        {/* Top list */}
        <div className="bg-white rounded-2xl border border-gray-100 shadow-sm p-5 flex-1">
          <h3 className="text-sm font-semibold text-gray-700 mb-4">Tổng quan lượt chi</h3>
          <div className="space-y-3">
            {items.slice(0, 5).map((c) => {
              const Icon = getCategoryIcon(catMap[c.category_id] ?? "");
              const color = getCategoryColor(catMap[c.category_id] ?? "");
              return (
                <div key={c.category_id} className="flex items-center gap-3">
                  <div className="w-8 h-8 rounded-lg flex items-center justify-center flex-shrink-0" style={{ background: color }}>
                    <Icon className="w-4 h-4 text-white" />
                  </div>
                  <span className="flex-1 text-sm text-gray-700">{c.category_name}</span>
                  <span className="text-sm font-semibold text-gray-800">{fmt(c.total_amount)}</span>
                </div>
              );
            })}
          </div>
        </div>

        {/* Stats card */}
        <div className="rounded-2xl p-5 text-white" style={{ background: "linear-gradient(135deg,#00137a 0%,#001BB7 100%)" }}>
          <p className="text-sm font-semibold mb-1">Thống kê</p>
          <p className="text-xs opacity-70 mb-4">Dựa trên {items.length} giao dịch chi tiêu.</p>
          <p className="text-3xl font-bold mb-1">{fmt(total)}</p>
          <p className="text-xs font-semibold tracking-widest opacity-70 uppercase">Tổng chi tiêu đã ghi nhận</p>
        </div>
      </div>
    </div>
  );
}

// ── Main ─────────────────────────────────────────────────────
export function ReportsPage() {
  const [range, setRange] = useState<Range>("month");
  const { start, end } = getDateRange(range);

  return (
    <div className="space-y-5">
      <BudgetSection />

      {/* Time filter */}
      <div className="bg-white rounded-2xl border border-gray-100 shadow-sm px-5 py-4 flex flex-col sm:flex-row sm:items-center gap-3">
        <div className="flex items-center gap-2 flex-1">
          <CalendarDays className="w-4 h-4 text-blue-500" />
          <span className="text-sm font-semibold text-gray-700">Thời gian báo cáo</span>
        </div>
        <div className="flex gap-1 bg-gray-50 rounded-xl p-1">
          {RANGES.map((r) => (
            <button
              key={r.key}
              onClick={() => setRange(r.key)}
              className={cn(
                "px-3 py-1.5 text-xs font-medium rounded-lg transition-colors",
                range === r.key ? "bg-white text-blue-600 shadow-sm" : "text-gray-500 hover:text-gray-700"
              )}
            >
              {r.label}
            </button>
          ))}
        </div>
      </div>

      <SpendingBarChart startDate={start} endDate={end} />
    </div>
  );
}
