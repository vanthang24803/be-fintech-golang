import {
  Utensils, Car, ShoppingBag, Gamepad2, Home, HeartPulse,
  BookOpen, Receipt, MoreHorizontal, CreditCard, Gift, TrendingUp,
  CirclePlus, Coffee, Plane, Heart, Star, Music,
  Camera, Briefcase, Bus, Bike, Wifi, Phone,
  Monitor, Mail, Video, Trash2, Zap, Wallet,
} from "lucide-react";
import type { LucideIcon } from "lucide-react";

export const ICON_MAP: Record<string, LucideIcon> = {
  "utensils": Utensils,
  "car": Car,
  "shopping-bag": ShoppingBag,
  "gamepad-2": Gamepad2,
  "home": Home,
  "heart-pulse": HeartPulse,
  "book-open": BookOpen,
  "receipt": Receipt,
  "more-horizontal": MoreHorizontal,
  "credit-card": CreditCard,
  "gift": Gift,
  "trending-up": TrendingUp,
  "circle-plus": CirclePlus,
  "coffee": Coffee,
  "plane": Plane,
  "heart": Heart,
  "star": Star,
  "music": Music,
  "camera": Camera,
  "briefcase": Briefcase,
  "bus": Bus,
  "bike": Bike,
  "wifi": Wifi,
  "phone": Phone,
  "monitor": Monitor,
  "mail": Mail,
  "video": Video,
  "trash-2": Trash2,
  "zap": Zap,
  "wallet": Wallet,
};

export const ICON_KEYS = Object.keys(ICON_MAP);

// Default colors for known icon names (used when no color is encoded)
const ICON_DEFAULT_COLORS: Record<string, string> = {
  "utensils": "#ef4444",
  "car": "#f97316",
  "shopping-bag": "#f59e0b",
  "gamepad-2": "#22c55e",
  "home": "#22c55e",
  "heart-pulse": "#14b8a6",
  "book-open": "#8b5cf6",
  "receipt": "#a855f7",
  "more-horizontal": "#6b7280",
  "credit-card": "#3b82f6",
  "gift": "#ec4899",
  "trending-up": "#22c55e",
  "circle-plus": "#22c55e",
  "wallet": "#001BB7",
};

// Encode icon + color into single field
export function encodeIcon(iconName: string, color: string) {
  return `${iconName}|${color}`;
}

// Parse icon field → { iconName, color }
export function parseIcon(field: string): { iconName: string; color: string } {
  if (field?.includes("|")) {
    const [iconName, color] = field.split("|");
    return { iconName, color };
  }
  return {
    iconName: field ?? "more-horizontal",
    color: ICON_DEFAULT_COLORS[field] ?? "#6b7280",
  };
}

export function getCategoryIcon(iconField: string): LucideIcon {
  const { iconName } = parseIcon(iconField);
  return ICON_MAP[iconName] ?? MoreHorizontal;
}

export function getCategoryColor(iconField: string): string {
  return parseIcon(iconField).color;
}
