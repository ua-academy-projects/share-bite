import { useQuery } from "@tanstack/react-query";
import { Bell, Loader2 } from "lucide-react";
import { fetchNotifications } from "@/api/notifications";
import { PageHeader } from "@/components/layout/PageHeader";
import { PageLayout } from "@/components/layout/PageLayout";
import { pageEmpty, pageLoader, pagePanel } from "@/components/layout/pageStyles";
import { cn } from "@/lib/utils";

export function NotificationsPage() {
  const token = localStorage.getItem("token")!;

  const { data: notifications = [], isLoading } = useQuery({
    queryKey: ["notifications"],
    queryFn: () => fetchNotifications(token, 50),
    enabled: !!token,
  });

  return (
    <PageLayout>
      <PageHeader
        title="Notifications"
        description="Stay up to date with your Share Bite activity"
      />

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
                <span className="whitespace-nowrap text-xs text-gray-500 dark:text-gray-400">
                  {new Date(notification.createdAt).toLocaleDateString()}
                </span>
              </div>
            </div>
          ))}
        </div>
      )}
    </PageLayout>
  );
}
