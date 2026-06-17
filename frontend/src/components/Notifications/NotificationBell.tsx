import { Link } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";
import { Bell } from "lucide-react";
import { fetchNotifications } from "@/api/notifications";
import { Button } from "@/components/ui/button";
import {
  Popover,
  PopoverContent,
  PopoverHeader,
  PopoverTitle,
  PopoverTrigger,
} from "@/components/ui/popover";
import { pageLinkAccent } from "@/components/layout/pageStyles";
import { cn } from "@/lib/utils";

type NotificationBellProps = {
  variant?: "default" | "compact";
};

export function NotificationBell({ variant = "default" }: NotificationBellProps) {
  const token = localStorage.getItem("token");

  const { data: notifications = [] } = useQuery({
    queryKey: ["notifications"],
    queryFn: () => fetchNotifications(token!, 20),
    refetchInterval: 30000,
    refetchIntervalInBackground: false,
    enabled: !!token,
  });

  const unreadCount = notifications.filter((n) => !n.read).length;

  const formatMessage = (n: (typeof notifications)[0]) =>
    n.message ||
    (n.metadata && typeof n.metadata === "object"
      ? String((n.metadata as Record<string, unknown>).message || "")
      : "") ||
    "You have a new notification";

  const formatDate = (n: (typeof notifications)[0]) => {
    const raw = n.createdAt;
    return raw ? new Date(raw).toLocaleDateString() : "";
  };

  return (
    <Popover>
      <PopoverTrigger asChild>
        <Button
          variant="ghost"
          size={variant === "compact" ? "icon-sm" : "sm"}
          className={cn(
            "relative shrink-0",
            variant === "compact"
              ? "h-8 w-8 rounded-full text-gray-400 hover:bg-[#2f5e50]/40 hover:text-[#98FF98]"
              : "rounded-full bg-[#FFD700]/15 text-[#FFD700] hover:bg-[#FFD700]/25"
          )}
          aria-label="Notifications"
        >
          <Bell className="h-4 w-4" />
          {unreadCount > 0 && (
            <span className="absolute -right-0.5 -top-0.5 flex h-4 w-4 items-center justify-center rounded-full bg-red-500 text-[10px] font-bold text-white">
              {unreadCount > 99 ? "99+" : unreadCount}
            </span>
          )}
        </Button>
      </PopoverTrigger>
      <PopoverContent
        side="right"
        align="start"
        sideOffset={12}
        collisionPadding={16}
        className="w-80 gap-0 overflow-hidden rounded-3xl border border-gray-200 bg-white p-0 shadow-lg dark:border-[#2f5e50] dark:bg-[#163d32]"
      >
        <PopoverHeader className="flex flex-row items-center justify-between gap-2 px-5 py-4">
          <div className="flex items-center gap-2">
            <Bell size={18} className="text-emerald-500 dark:text-[#98FF98]" />
            <PopoverTitle className="text-sm font-semibold text-[#1A3C34] dark:text-white">
              Notifications
            </PopoverTitle>
          </div>
          <span className="rounded-full border border-emerald-500/40 bg-emerald-500/10 px-2.5 py-0.5 text-xs font-medium text-emerald-600 dark:text-emerald-300">
            {unreadCount} new
          </span>
        </PopoverHeader>

        <div className="max-h-80 overflow-y-auto border-t border-gray-200 dark:border-[#2f5e50]">
          {notifications.length === 0 ? (
            <div className="px-5 py-10 text-center">
              <Bell className="mx-auto mb-3 h-10 w-10 text-gray-400 opacity-30 dark:text-gray-500" />
              <p className="text-sm font-semibold text-[#1A3C34] dark:text-gray-200">
                No notifications yet
              </p>
              <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">
                When you get notifications, they&apos;ll show up here.
              </p>
            </div>
          ) : (
            notifications.map((notification) => (
              <Link
                key={String(notification.id)}
                to="/notifications"
                className={cn(
                  "block border-b border-gray-200 px-5 py-3 transition-colors last:border-b-0 hover:bg-gray-50 dark:border-[#2f5e50] dark:hover:bg-[#0d241d]/60",
                  !notification.read && "bg-emerald-500/5 dark:bg-[#98FF98]/5"
                )}
              >
                <p className="text-sm leading-snug text-[#1A3C34] dark:text-white">
                  {formatMessage(notification)}
                </p>
                <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">
                  {formatDate(notification)}
                </p>
              </Link>
            ))
          )}
        </div>

        <Link
          to="/notifications"
          className={cn(
            pageLinkAccent,
            "block border-t border-gray-200 px-5 py-3 text-center text-sm dark:border-[#2f5e50]"
          )}
        >
          View all notifications
        </Link>
      </PopoverContent>
    </Popover>
  );
}
