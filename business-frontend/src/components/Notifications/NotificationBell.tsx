import { Link } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";
import { Bell } from "lucide-react";
import { fetchNotifications } from "@/api/notifications";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Popover,
  PopoverContent,
  PopoverHeader,
  PopoverTitle,
  PopoverTrigger,
} from "@/components/ui/popover";
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
              ? "h-8 w-8 rounded-full text-muted-foreground hover:text-foreground"
              : "rounded-full bg-accent/15 text-accent hover:bg-accent/25"
          )}
          aria-label="Notifications"
        >
          <Bell className="h-4 w-4" />
          {unreadCount > 0 && (
            <span className="absolute -right-0.5 -top-0.5 flex h-4 w-4 items-center justify-center rounded-full bg-destructive text-[10px] font-bold text-white">
              {unreadCount > 99 ? "99+" : unreadCount}
            </span>
          )}
        </Button>
      </PopoverTrigger>
      <PopoverContent
        side="right"
        align="end"
        className="w-80 gap-0 overflow-hidden border-border bg-card-solid p-0 shadow-2xl"
      >
        <PopoverHeader className="flex flex-row items-center justify-between border-b border-border bg-card-solid px-4 py-3">
          <PopoverTitle className="text-sm font-semibold tracking-tight">
            Latest stories
          </PopoverTitle>
          <Badge variant="default">{unreadCount} new</Badge>
        </PopoverHeader>

        <div className="max-h-80 overflow-y-auto bg-card-solid">
          {notifications.length === 0 ? (
            <div className="bg-card-solid px-4 py-10 text-center text-muted-foreground">
              <Bell className="mx-auto mb-3 h-8 w-8 opacity-30" />
              <p className="text-sm font-medium">No new notifications</p>
              <p className="mt-1 text-xs">You&apos;re all caught up!</p>
            </div>
          ) : (
            notifications.map((notification) => (
              <Link
                key={String(notification.id)}
                to="/notifications"
                className="block border-b border-border/60 bg-card-solid px-4 py-3 transition-colors hover:bg-[#2f5e50]/30"
              >
                <p className="text-sm leading-snug text-foreground">
                  {formatMessage(notification)}
                </p>
                <p className="mt-1 text-xs text-muted-foreground">
                  {formatDate(notification)}
                </p>
              </Link>
            ))
          )}
        </div>

        <Link
          to="/notifications"
          className="block bg-card-solid px-4 py-3 text-center text-sm font-semibold text-accent hover:bg-[#2f5e50]/40"
        >
          View all notifications
        </Link>
      </PopoverContent>
    </Popover>
  );
}
