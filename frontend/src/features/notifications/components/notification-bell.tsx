"use client";

import { useNotifications, useUnreadNotificationsCount, useMarkNotificationsRead } from "@/features/notifications/hooks/use-notifications";
import { Bell, Loader2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { cn } from "@/lib/utils";

export function NotificationBell() {
  const { data: notifications, isLoading } = useNotifications();
  const { data: unreadCount } = useUnreadNotificationsCount();
  const { mutate: markRead } = useMarkNotificationsRead();

  const handleMarkAllRead = () => {
    const unreadIds = notifications?.filter(n => !n.is_read).map(n => n.id) || [];
    if (unreadIds.length > 0) {
      markRead(unreadIds);
    }
  };

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="ghost" size="icon" className="relative">
          <Bell className="size-5" />
          {unreadCount && unreadCount.count > 0 && (
            <span className="absolute right-2 top-2 flex size-2 items-center justify-center rounded-full bg-destructive" />
          )}
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="w-80">
        <DropdownMenuLabel className="flex items-center justify-between">
          Notifications
          {unreadCount && unreadCount.count > 0 && (
            <Button variant="ghost" size="sm" className="h-auto p-0 text-xs font-normal" onClick={handleMarkAllRead}>
              Mark all as read
            </Button>
          )}
        </DropdownMenuLabel>
        <DropdownMenuSeparator />
        <div className="max-h-96 overflow-y-auto">
          {isLoading ? (
            <div className="flex h-20 items-center justify-center">
              <Loader2 className="size-4 animate-spin text-muted-foreground" />
            </div>
          ) : notifications?.length === 0 ? (
            <div className="p-4 text-center text-sm text-muted-foreground">
              No notifications yet.
            </div>
          ) : (
            notifications?.map((n) => (
              <DropdownMenuItem 
                key={n.id} 
                className={cn(
                  "flex flex-col items-start gap-1 p-4 cursor-default focus:bg-accent/50",
                  !n.is_read && "bg-primary/5"
                )}
              >
                <div className="flex w-full items-center justify-between gap-2">
                  <span className="font-semibold text-sm">{n.title}</span>
                  <span className="text-[10px] text-muted-foreground whitespace-nowrap">
                    {new Date(n.created_at).toLocaleDateString()}
                  </span>
                </div>
                <p className="text-xs text-muted-foreground line-clamp-2">{n.body}</p>
                {!n.is_read && (
                  <Button 
                    variant="ghost" 
                    size="sm" 
                    className="h-auto p-0 text-[10px] font-normal mt-1 hover:bg-transparent"
                    onClick={() => markRead([n.id])}
                  >
                    Mark as read
                  </Button>
                )}
              </DropdownMenuItem>
            ))
          )}
        </div>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
