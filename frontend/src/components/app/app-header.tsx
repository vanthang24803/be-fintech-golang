"use client";

import { Menu, RefreshCw, Tag, Plus } from "lucide-react";
import { Button } from "@/components/ui/button";
import { QuickCaptureDialog } from "@/features/dashboard/components/quick-capture-dialog";
import { RecurringDialog } from "@/features/dashboard/components/recurring-dialog";
import { AddCategoryDialog } from "@/features/categories/components/add-category-dialog";

type AppHeaderProps = {
  onOpenMobileNav: () => void;
};

export function AppHeader({ onOpenMobileNav }: AppHeaderProps) {
  return (
    <header className="sticky top-0 z-20 flex h-16 items-center bg-white border-b border-gray-100 px-4 md:px-6 gap-3">
      {/* Mobile menu */}
      <Button variant="ghost" size="icon" className="lg:hidden" onClick={onOpenMobileNav}>
        <Menu className="w-5 h-5" />
      </Button>

      {/* Brand — hidden on mobile */}
      <div className="flex-1 hidden md:flex flex-col items-center justify-center">
        <p className="text-base font-extrabold text-gray-900 leading-tight tracking-wide">FINANSMART</p>
        <p className="text-xs text-gray-400">Làm chủ tài chính, kiến tạo tương lai!</p>
      </div>
      <div className="flex-1 md:hidden" />

      {/* Actions */}
      <div className="flex items-center gap-2">
        {/* Desktop: text + icon */}
        <RecurringDialog
          trigger={
            <Button variant="outline" size="sm" className="hidden md:flex gap-1.5 text-gray-600 text-xs h-8 rounded-lg border-gray-200">
              <RefreshCw className="w-3.5 h-3.5" />
              Định kỳ
            </Button>
          }
        />
        <AddCategoryDialog
          trigger={
            <Button variant="outline" size="sm" className="hidden md:flex gap-1.5 text-gray-600 text-xs h-8 rounded-lg border-gray-200">
              <Tag className="w-3.5 h-3.5" />
              Danh mục
            </Button>
          }
        />
        <QuickCaptureDialog
          trigger={
            <Button size="sm" className="hidden md:flex gap-1.5 text-xs h-8 rounded-lg text-white" style={{ background: "#001BB7" }}>
              <Plus className="w-3.5 h-3.5" />
              Thêm mới
            </Button>
          }
        />

        {/* Mobile: icon only */}
        <RecurringDialog
          trigger={
            <Button variant="outline" size="icon" className="md:hidden h-8 w-8 rounded-lg border-gray-200 text-gray-600">
              <RefreshCw className="w-4 h-4" />
            </Button>
          }
        />
        <AddCategoryDialog
          trigger={
            <Button variant="outline" size="icon" className="md:hidden h-8 w-8 rounded-lg border-gray-200 text-gray-600">
              <Tag className="w-4 h-4" />
            </Button>
          }
        />
        <QuickCaptureDialog
          trigger={
            <Button size="icon" className="md:hidden h-8 w-8 rounded-lg text-white" style={{ background: "#001BB7" }}>
              <Plus className="w-4 h-4" />
            </Button>
          }
        />
      </div>
    </header>
  );
}
