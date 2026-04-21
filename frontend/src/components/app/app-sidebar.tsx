"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { ChevronLeft, ChevronRight, LogOut, Wallet } from "lucide-react";
import { primaryNavigation } from "@/lib/navigation/config";
import { isNavItemActive } from "@/lib/navigation/matchers";
import { useAuth } from "@/features/auth/context/auth-context";
import { cn } from "@/lib/utils";

type AppSidebarProps = {
  collapsed?: boolean;
  onToggle?: () => void;
  mobile?: boolean;
  onNavigate?: () => void;
};

export function AppSidebar({ collapsed = false, onToggle, mobile = false, onNavigate }: AppSidebarProps) {
  const pathname = usePathname();
  const { logout } = useAuth();

  return (
    <aside
      className={cn(
        "flex h-full flex-col bg-white border-r border-gray-100 transition-all duration-200",
        collapsed && !mobile ? "w-[68px]" : "w-[200px]",
        mobile && "w-full border-r-0",
      )}
    >
      {/* Logo */}
      <div className={cn("flex items-center gap-3 px-4 py-5", collapsed && !mobile && "justify-center px-0")}>
        <div
          className="w-9 h-9 rounded-xl flex items-center justify-center flex-shrink-0"
          style={{ background: "#001BB7" }}
        >
          <Wallet className="w-4 h-4 text-white" />
        </div>
        {(!collapsed || mobile) && (
          <span className="text-sm font-bold text-gray-800 tracking-wide">FINANSMART</span>
        )}
      </div>

      {/* Nav */}
      <nav className="flex-1 px-3 space-y-0.5 mt-1">
        {primaryNavigation.map((item) => {
          const active = isNavItemActive(pathname, item.href);
          return (
            <Link
              key={item.href}
              href={item.href}
              onClick={onNavigate}
              className={cn(
                "flex items-center gap-3 rounded-xl px-3 py-2.5 text-sm font-medium transition-colors",
                active ? "text-[#001BB7] bg-blue-50" : "text-gray-500 hover:text-gray-800 hover:bg-gray-50",
                collapsed && !mobile && "justify-center px-0",
              )}
            >
              <item.icon className={cn("w-4 h-4 flex-shrink-0", active && "text-[#001BB7]")} />
              {(!collapsed || mobile) && <span>{item.title}</span>}
            </Link>
          );
        })}
      </nav>

      {/* Bottom */}
      <div className={cn("px-3 pb-5 space-y-0.5", collapsed && !mobile && "px-1")}>
        <button
          onClick={() => logout()}
          className={cn(
            "flex items-center gap-3 w-full rounded-xl px-3 py-2.5 text-sm font-medium text-gray-400 hover:text-red-500 hover:bg-red-50 transition-colors",
            collapsed && !mobile && "justify-center px-0",
          )}
        >
          <LogOut className="w-4 h-4 flex-shrink-0" />
          {(!collapsed || mobile) && <span>Đăng xuất</span>}
        </button>

        {!mobile && onToggle && (
          <button
            onClick={onToggle}
            className={cn(
              "flex items-center gap-2 w-full rounded-xl px-3 py-2 text-xs text-gray-400 hover:text-gray-600 hover:bg-gray-50 transition-colors",
              collapsed && "justify-center px-0",
            )}
          >
            {collapsed ? (
              <ChevronRight className="w-4 h-4" />
            ) : (
              <>
                <ChevronLeft className="w-4 h-4" />
                <span>Thu gọn</span>
              </>
            )}
          </button>
        )}
      </div>
    </aside>
  );
}
