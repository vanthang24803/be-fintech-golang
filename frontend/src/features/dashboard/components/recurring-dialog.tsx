"use client";

import React, { useState } from "react";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { useCategories } from "@/features/categories/hooks/use-categories";
import { getCategoryIcon, getCategoryColor } from "@/lib/category-icon";
import { cn } from "@/lib/utils";
import { RefreshCw } from "lucide-react";
import { toast } from "sonner";

const FREQUENCIES = [
  { key: "daily",   label: "Hàng ngày" },
  { key: "weekly",  label: "Hàng tuần" },
  { key: "monthly", label: "Hàng tháng" },
  { key: "yearly",  label: "Hàng năm" },
];

function fmt(n: string) {
  const num = Number(n.replace(/\D/g, ""));
  return isNaN(num) ? "0" : num.toLocaleString("vi-VN");
}

export function RecurringDialog({ trigger }: { trigger?: React.ReactNode }) {
  const [open, setOpen] = useState(false);
  const [type, setType] = useState<"expense" | "income">("expense");
  const [amount, setAmount] = useState("");
  const [frequency, setFrequency] = useState("monthly");
  const [categoryId, setCategoryId] = useState("");
  const [startDate, setStartDate] = useState(() => new Date().toISOString().slice(0, 10));
  const [note, setNote] = useState("");

  const { data: categories } = useCategories();
  const filtered = (categories ?? []).filter((c) => c.type === type);

  function handleSubmit() {
    if (!amount || !categoryId) return;
    toast.success("Đã lưu giao dịch định kỳ!");
    setOpen(false);
    setAmount("");
    setCategoryId("");
    setNote("");
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        {trigger ?? <button>Định kỳ</button>}
      </DialogTrigger>
      <DialogContent className="max-w-sm rounded-2xl p-0 overflow-hidden">
        <DialogHeader className="px-5 pt-5 pb-0">
          <div className="flex items-center gap-2">
            <RefreshCw className="w-4 h-4 text-blue-500" />
            <DialogTitle className="text-base font-semibold">Giao dịch định kỳ</DialogTitle>
          </div>
        </DialogHeader>

        <div className="px-5 pb-5 space-y-4 max-h-[85vh] overflow-y-auto mt-3">
          {/* Type toggle */}
          <div className="flex rounded-xl border border-gray-200 p-1">
            <button
              onClick={() => { setType("expense"); setCategoryId(""); }}
              className={cn(
                "flex-1 py-1.5 text-sm font-medium rounded-lg transition-colors",
                type === "expense" ? "bg-red-50 text-red-500 border border-red-200" : "text-gray-400"
              )}
            >
              Chi tiêu
            </button>
            <button
              onClick={() => { setType("income"); setCategoryId(""); }}
              className={cn(
                "flex-1 py-1.5 text-sm font-medium rounded-lg transition-colors",
                type === "income" ? "bg-green-50 text-green-600 border border-green-200" : "text-gray-400"
              )}
            >
              Thu nhập
            </button>
          </div>

          {/* Amount */}
          <div>
            <p className="text-xs text-gray-500 mb-1">Số tiền</p>
            <div className="flex items-end justify-between border-b border-gray-200 pb-1">
              <input
                type="text"
                inputMode="numeric"
                value={fmt(amount)}
                onChange={(e) => setAmount(e.target.value.replace(/\D/g, ""))}
                className="text-2xl font-semibold text-gray-900 bg-transparent outline-none w-full"
                placeholder="0"
              />
              <span className="text-sm text-gray-400 ml-2 mb-1">VND</span>
            </div>
          </div>

          {/* Frequency */}
          <div>
            <p className="text-xs text-gray-500 mb-2">Tần suất</p>
            <div className="grid grid-cols-2 gap-2">
              {FREQUENCIES.map((f) => (
                <button
                  key={f.key}
                  onClick={() => setFrequency(f.key)}
                  className={cn(
                    "py-2 rounded-xl text-sm font-medium border transition-colors",
                    frequency === f.key
                      ? "text-white border-transparent"
                      : "text-gray-500 border-gray-200 hover:bg-gray-50"
                  )}
                  style={frequency === f.key ? { background: "#001BB7" } : {}}
                >
                  {f.label}
                </button>
              ))}
            </div>
          </div>

          {/* Categories */}
          <div>
            <p className="text-xs text-gray-500 mb-2">Danh mục</p>
            <div className="grid grid-cols-4 gap-3 max-h-48 overflow-y-auto pr-1">
              {filtered.map((cat) => {
                const Icon = getCategoryIcon(cat.icon ?? "");
                const color = getCategoryColor(cat.icon ?? "");
                const selected = categoryId === cat.id;
                return (
                  <button
                    key={cat.id}
                    onClick={() => setCategoryId(cat.id)}
                    className="flex flex-col items-center gap-1"
                  >
                    <div
                      className={cn(
                        "w-12 h-12 rounded-full flex items-center justify-center transition-all",
                        selected ? "ring-2 ring-offset-2" : "opacity-90 hover:opacity-100"
                      )}
                      style={{ background: color }}
                    >
                      <Icon className="w-5 h-5 text-white" />
                    </div>
                    <span className="text-[11px] text-gray-600 text-center leading-tight">{cat.name}</span>
                  </button>
                );
              })}
            </div>
          </div>

          {/* Start date */}
          <div>
            <p className="text-xs text-gray-500 mb-1 flex items-center gap-1">
              <span>📅</span> Ngày bắt đầu
            </p>
            <input
              type="date"
              value={startDate}
              onChange={(e) => setStartDate(e.target.value)}
              className="w-full border border-gray-200 rounded-xl px-3 py-2.5 text-sm text-gray-700 outline-none focus:border-gray-400"
            />
          </div>

          {/* Note */}
          <div>
            <p className="text-xs text-gray-500 mb-1">Ghi chú</p>
            <input
              type="text"
              value={note}
              onChange={(e) => setNote(e.target.value)}
              placeholder="Ví dụ: Tiền phòng, Lương cứng..."
              className="w-full border border-gray-200 rounded-xl px-3 py-2.5 text-sm text-gray-700 outline-none focus:border-gray-400"
            />
          </div>

          {/* Submit */}
          <button
            onClick={handleSubmit}
            disabled={!amount || !categoryId}
            className="w-full py-3 rounded-xl text-white text-sm font-semibold disabled:opacity-50"
            style={{ background: "#001BB7" }}
          >
            Lưu thiết lập
          </button>
        </div>
      </DialogContent>
    </Dialog>
  );
}
