import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiRequest } from "@/lib/api/client";

export interface Notification {
  id: string;
  user_id: string;
  source: string;
  source_id: string | null;
  type: string;
  title: string;
  body: string;
  is_read: boolean;
  read_at: string | null;
  created_at: string;
}

export function useNotifications(filters: any = {}) {
  return useQuery({
    queryKey: ["notifications", "list", filters],
    queryFn: () => apiRequest<Notification[]>("notifications/list", {
      body: JSON.stringify(filters),
    }),
  });
}

export function useUnreadNotificationsCount() {
  return useQuery({
    queryKey: ["notifications", "unread-count"],
    queryFn: () => apiRequest<{ count: number }>("notifications/unread-count"),
  });
}

export function useMarkNotificationsRead() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (ids: string[]) =>
      apiRequest("notifications/mark-read", {
        body: JSON.stringify({ ids }),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["notifications"] });
    },
  });
}
