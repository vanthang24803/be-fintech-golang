"use client";

import React, { useState } from "react";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { useCreateCategory } from "@/features/categories/hooks/use-categories";
import { ICON_MAP, ICON_KEYS, encodeIcon } from "@/lib/category-icon";
import { cn } from "@/lib/utils";
import { Loader2, Check } from "lucide-react";

const COLOR_PALETTE = [
  "#ef4444", "#f97316", "#f59e0b", "#eab308", "#84cc16",
  "#22c55e", "#10b981", "#14b8a6", "#06b6d4", "#0ea5e9",
  "#3b82f6", "#6366f1", "#8b5cf6", "#a855f7", "#d946ef",
  "#ec4899", "#f43f5e", "#64748b", "#475569", "#1e293b",
];

export function AddCategoryDialog({ trigger }: { trigger?: React.ReactNode }) {
  const [open, setOpen] = useState(false);
  const [type, setType] = useState<"expense" | "income">("expense");
  const [name, setName] = useState("");
  const [icon, setIcon] = useState("more-horizontal");
  const [color, setColor] = useState("#ef4444");

  const { mutate: createCategory, isPending } = useCreateCategory();

  function handleSubmit() {
    if (!name.trim()) return;
    createCategory(
      { name: name.trim(), type, icon: encodeIcon(icon, color) },
      {
        onSuccess: () => {
          setOpen(false);
          setName("");
          setIcon("more-horizontal");
          setColor("#ef4444");
        },
      }
    );
  }

  const PreviewIcon = ICON_MAP[icon] ?? ICON_MAP["more-horizontal"];

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        {trigger ?? <button>Thêm danh mục</button>}
      </DialogTrigger>
      <DialogContent className="max-w-sm rounded-2xl p-0 overflow-hidden">
        <DialogHeader className="px-5 pt-5 pb-0">
          <DialogTitle className="text-base font-semibold">Thêm danh mục mới</DialogTitle>
        </DialogHeader>

        <div className="px-5 pb-5 space-y-4 max-h-[85vh] overflow-y-auto mt-3">
          {/* Type toggle */}
          <div className="flex rounded-xl border border-gray-200 p-1">
            <button
              onClick={() => setType("expense")}
              className={cn(
                "flex-1 py-1.5 text-sm font-medium rounded-lg transition-colors",
                type === "expense" ? "bg-red-50 text-red-500 border border-red-200" : "text-gray-400"
              )}
            >
              Chi tiêu
            </button>
            <button
              onClick={() => setType("income")}
              className={cn(
                "flex-1 py-1.5 text-sm font-medium rounded-lg transition-colors",
                type === "income" ? "bg-green-50 text-green-600 border border-green-200" : "text-gray-400"
              )}
            >
              Thu nhập
            </button>
          </div>

          {/* Name */}
          <div>
            <p className="text-xs text-gray-500 mb-1">Tên danh mục</p>
            <input
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="Ví dụ: Ăn vặt, Tiền điện..."
              className="w-full border border-gray-200 rounded-xl px-3 py-2.5 text-sm text-gray-700 outline-none focus:border-gray-400"
            />
          </div>

          {/* Icon picker */}
          <div>
            <p className="text-xs text-gray-500 mb-2">Chọn biểu tượng</p>
            <div className="grid grid-cols-6 gap-2">
              {ICON_KEYS.map((key) => {
                const Icon = ICON_MAP[key];
                return (
                  <button
                    key={key}
                    onClick={() => setIcon(key)}
                    className={cn(
                      "w-10 h-10 rounded-xl flex items-center justify-center transition-all border",
                      icon === key
                        ? "border-gray-400 bg-gray-100"
                        : "border-transparent hover:bg-gray-50"
                    )}
                  >
                    <Icon className="w-5 h-5 text-gray-600" />
                  </button>
                );
              })}
            </div>
          </div>

          {/* Color picker */}
          <div>
            <p className="text-xs text-gray-500 mb-2">Chọn màu sắc</p>
            <div className="grid grid-cols-5 gap-2">
              {COLOR_PALETTE.map((c) => (
                <button
                  key={c}
                  onClick={() => setColor(c)}
                  className="w-10 h-10 rounded-full flex items-center justify-center transition-transform hover:scale-110"
                  style={{ background: c }}
                >
                  {color === c && <Check className="w-4 h-4 text-white" />}
                </button>
              ))}
            </div>
          </div>

          {/* Preview */}
          <div className="flex items-center gap-3 bg-gray-50 rounded-xl p-3">
            <div
              className="w-11 h-11 rounded-full flex items-center justify-center flex-shrink-0"
              style={{ background: color }}
            >
              <PreviewIcon className="w-5 h-5 text-white" />
            </div>
            <div>
              <p className="text-sm font-medium text-gray-800">{name || "Tên danh mục"}</p>
              <p className="text-xs text-gray-400">{type === "expense" ? "Chi tiêu" : "Thu nhập"}</p>
            </div>
          </div>

          {/* Submit */}
          <button
            onClick={handleSubmit}
            disabled={isPending || !name.trim()}
            className="w-full py-3 rounded-xl text-white text-sm font-semibold disabled:opacity-50 flex items-center justify-center gap-2"
            style={{ background: "#001BB7" }}
          >
            {isPending && <Loader2 className="w-4 h-4 animate-spin" />}
            Lưu danh mục
          </button>
        </div>
      </DialogContent>
    </Dialog>
  );
}
