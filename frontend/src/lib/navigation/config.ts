import {
  LayoutGrid,
  ArrowLeftRight,
  BarChart2,
  Tag,
  User,
  type LucideIcon,
} from "lucide-react";

export type NavigationItem = {
  title: string;
  href: string;
  icon: LucideIcon;
};

export const primaryNavigation: NavigationItem[] = [
  { title: "Tổng quan",  href: "/dashboard",    icon: LayoutGrid },
  { title: "Giao dịch",  href: "/transactions", icon: ArrowLeftRight },
  { title: "Báo cáo",    href: "/budgets",       icon: BarChart2 },
  { title: "Danh mục",   href: "/categories",    icon: Tag },
  { title: "Tài khoản",  href: "/settings",      icon: User },
];
