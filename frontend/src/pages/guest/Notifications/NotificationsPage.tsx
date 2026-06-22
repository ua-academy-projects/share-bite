import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Navigate } from "react-router-dom";
import { Bell, Check, CheckCheck, Loader2 } from "lucide-react";
import { fetchNotifications, markNotificationsRead } from "@/api/notifications";
import { PageHeader } from "@/components/layout/PageHeader";
import { PageLayout } from "@/components/layout/PageLayout";
import { pageEmpty, pageLoader, pagePanel } from "@/components/layout/pageStyles";
import { cn } from "@/lib/utils";

export function NotificationsPage() {
  const token = localStorage.getItem("token");

  const queryClient = useQueryClient();

  const { data: notifications = [], isLoading } = useQuery({
    queryKey: ["notifications"],
    queryFn: () => {
      if (!token) return Promise.resolve([]);
      return fetchNotifications(token, 50);
    },
    enabled: !!token,
  });

  const markRead = useMutation({
    mutationFn: (ids: string[]) => markNotificationsRead(token!, ids),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["notifications"] }),
  });

  const unreadIds = notifications.filter((n) => !n.read).map((n) => n.id);

  if (!token) {
    return <Navigate to="/auth" replace />;
  }

  return (
    <PageLayout>
      <PageHeader
        title="Notifications"
        description="Stay up to date with your Share Bite activity"
      />

      {unreadIds.length > 0 && (
        <div className="mb-3 flex justify-end">
          <button
            type="button"
            onClick={() => markRead.mutate(unreadIds)}
            disabled={markRead.isPending}
            className="inline-flex items-center gap-1.5 rounded-full border border-emerald-500/40 bg-emerald-500/10 px-3 py-1.5 text-xs font-medium text-emerald-600 transition-colors hover:bg-emerald-500/20 disabled:opacity-50 dark:text-emerald-300"
          >
            <CheckCheck className="h-4 w-4" />
            Mark all as read
          </button>
        </div>
      )}

      {isLoading ? (
        <div className="flex h-64 items-center justify-center">
          <Loader2 className={cn(pageLoader, "h-12 w-12")} />
        </div>
      ) : notifications.length === 0 ? (
        <div className={pageEmpty}>
          <Bell className="mx-auto mb-4 h-12 w-12 opacity-20" />
          <p className="text-xl font-bold text-[#1A3C34] dark:text-gray-300">
            No notifications yet
          </p>
          <p className="mt-2 text-gray-500">
            When you get notifications, they&apos;ll show up here.
          </p>
        </div>
      ) : (
        <div className="flex flex-col gap-3">
          {notifications.map((notification) => (
            <div
              key={notification.id}
              className={cn(
                pagePanel,
                "p-4",
                !notification.read && "border-emerald-500/30 bg-emerald-500/5 dark:border-[#98FF98]/30"
              )}
            >
              <div className="flex items-start justify-between gap-4">
                <p className="text-sm leading-relaxed text-[#1A3C34] dark:text-white">
                  {notification.message ||
                    (notification.metadata &&
                    typeof notification.metadata.message === "string"
                      ? notification.metadata.message
                      : "You have a new notification")}
                </p>
                <div className="flex shrink-0 flex-col items-end gap-1.5">
                  <span className="whitespace-nowrap text-xs text-gray-500 dark:text-gray-400">
                    {new Date(notification.createdAt).toLocaleDateString()}
                  </span>
                  {!notification.read && (
                    <button
                      type="button"
                      onClick={() => markRead.mutate([notification.id])}
                      disabled={markRead.isPending}
                      className="inline-flex items-center gap-1 text-xs font-medium text-emerald-600 transition-colors hover:text-emerald-700 disabled:opacity-50 dark:text-emerald-300 dark:hover:text-emerald-200"
                    >
                      <Check className="h-3.5 w-3.5" />
                      Mark as read
                    </button>
                  )}
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </PageLayout>
  );
}
