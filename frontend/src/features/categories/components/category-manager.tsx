"use client";

import { useCategories, useDeleteCategory } from "@/features/categories/hooks/use-categories";
import { AddCategoryDialog } from "@/features/categories/components/add-category-dialog";
import { getCategoryIcon, getCategoryColor } from "@/lib/category-icon";
import { Plus, Trash2, Loader2 } from "lucide-react";
import { cn } from "@/lib/utils";
import { Category } from "@/lib/api/types";

function CategoryCard({ cat, onDelete }: { cat: Category; onDelete: (id: string) => void }) {
  const Icon = getCategoryIcon(cat.icon ?? "");
  const color = getCategoryColor(cat.icon ?? "");

  return (
    <div className="flex items-center gap-3 px-4 py-3.5 border-b border-gray-100 last:border-0 hover:bg-gray-50 transition-colors group">
      <div
        className="w-10 h-10 rounded-xl flex items-center justify-center flex-shrink-0"
        style={{ background: color }}
      >
        <Icon className="w-5 h-5 text-white" />
      </div>
      <span className="flex-1 text-sm font-medium text-gray-800">{cat.name}</span>
      <button
        onClick={() => onDelete(cat.id)}
        className="opacity-0 group-hover:opacity-100 w-7 h-7 rounded-lg hover:bg-gray-100 flex items-center justify-center transition-all"
      >
        <Trash2 className="w-3.5 h-3.5 text-gray-400" />
      </button>
    </div>
  );
}

function CategoryPanel({ title, categories, onDelete }: {
  title: string;
  categories: Category[];
  onDelete: (id: string) => void;
}) {
  return (
    <div className="bg-white rounded-2xl border border-gray-100 shadow-sm overflow-hidden">
      <div className="px-4 py-3.5 border-b border-gray-200">
        <h3 className="text-sm font-semibold text-gray-700">{title}</h3>
      </div>
      {categories.length === 0 ? (
        <div className="px-4 py-8 text-center text-sm text-gray-400">Chưa có danh mục</div>
      ) : (
        categories.map((cat) => (
          <CategoryCard key={cat.id} cat={cat} onDelete={onDelete} />
        ))
      )}
    </div>
  );
}

export function CategoryManager() {
  const { data: categories, isLoading } = useCategories();
  const { mutate: deleteCategory } = useDeleteCategory();

  const expense = (categories ?? []).filter((c) => c.type === "expense");
  const income = (categories ?? []).filter((c) => c.type === "income");

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <h1 className="text-xl font-bold text-gray-900">Quản lý danh mục</h1>
        <AddCategoryDialog
          trigger={
            <button
              className="flex items-center gap-2 px-4 py-2 rounded-xl text-white text-sm font-medium"
              style={{ background: "#001BB7" }}
            >
              <Plus className="w-4 h-4" />
              Thêm danh mục
            </button>
          }
        />
      </div>

      {/* Panels */}
      {isLoading ? (
        <div className="flex items-center justify-center h-48">
          <Loader2 className="w-6 h-6 animate-spin text-gray-300" />
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-5">
          <CategoryPanel
            title="Chi tiêu"
            categories={expense}
            onDelete={(id) => deleteCategory(id)}
          />
          <CategoryPanel
            title="Thu nhập"
            categories={income}
            onDelete={(id) => deleteCategory(id)}
          />
        </div>
      )}
    </div>
  );
}
