import React, { useState, useRef, useEffect } from "react";
import { Bell } from "lucide-react";
import { Link } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";
import { apiClient } from "@/api/client";
import type { NotificationItem } from "@/types/api";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";

export const NotificationBell: React.FC = () => {
  const [isOpen, setIsOpen] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);

  const { data: notifications } = useQuery({
    queryKey: ["notifications"],
    queryFn: () => apiClient.getNotifications(),
    refetchInterval: 30000,
    refetchIntervalInBackground: false,
    enabled: !!localStorage.getItem("token"),
  });

  const unreadCount =
    notifications?.items?.filter((n: { read?: boolean }) => !n.read).length ||
    0;

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        dropdownRef.current &&
        !dropdownRef.current.contains(event.target as Node)
      ) {
        setIsOpen(false);
      }
    };
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  return (
    <div className="relative" ref={dropdownRef}>
      <Button
        variant="ghost"
        size="icon"
        className="relative rounded-full"
        style={{ color: "var(--navbar-muted)" }}
        onClick={() => setIsOpen(!isOpen)}
        aria-label="Notifications"
      >
        <Bell size={20} />
        {unreadCount > 0 && (
          <Badge
            variant="destructive"
            className="absolute -right-1 -top-1 h-5 min-w-5 px-1 text-[10px]"
          >
            {unreadCount > 99 ? "99+" : unreadCount}
          </Badge>
        )}
      </Button>

      {isOpen && (
        <div className="absolute right-0 z-50 mt-3 w-80 overflow-hidden rounded-3xl border border-border bg-popover shadow-2xl">
          <div className="flex items-center justify-between border-b border-border bg-muted/30 p-5 backdrop-blur-md">
            <h3 className="font-serif text-2xl font-bold text-accent">
              Latest Stories
            </h3>
            <Badge variant="secondary" className="text-xs font-bold">
              {unreadCount} New
            </Badge>
          </div>
          <div className="max-h-80 overflow-y-auto">
            {!notifications?.items?.length ? (
              <div className="p-8 text-center text-muted-foreground">
                <Bell className="mx-auto mb-3 size-8 opacity-20" />
                <p className="text-sm font-semibold">No new notifications</p>
              </div>
            ) : (
              notifications.items.map((notification: NotificationItem) => (
                  <Link
                    key={String(notification.id)}
                    to="/notifications"
                    className="block border-b border-border p-4 transition-colors hover:bg-muted/40"
                    onClick={() => setIsOpen(false)}
                  >
                    <p className="text-sm font-medium leading-tight text-foreground">
                      {notification.message || "You have a new notification"}
                    </p>
                    <p className="mt-1 text-xs font-semibold text-muted-foreground">
                      {new Date(notification.createdAt).toLocaleDateString()}
                    </p>
                  </Link>
                )
              )
            )}
          </div>
          <Link
            to="/notifications"
            className="block w-full bg-muted/80 p-3 text-center text-sm font-bold text-accent hover:bg-accent hover:text-accent-foreground"
            onClick={() => setIsOpen(false)}
          >
            View All Notifications
          </Link>
        </div>
      )}
    </div>
  );
};
