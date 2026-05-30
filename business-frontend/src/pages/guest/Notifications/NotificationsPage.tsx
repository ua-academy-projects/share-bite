import { useQuery } from "@tanstack/react-query";
import { Bell, Loader2 } from "lucide-react";
import { fetchNotifications } from "@/api/notifications";
import { Card, CardContent } from "@/components/ui/card";
import { PageHeader } from "@/components/layout/PageHeader";

export function NotificationsPage() {
  const token = localStorage.getItem("token")!;

  const { data: notifications = [], isLoading } = useQuery({
    queryKey: ["notifications"],
    queryFn: () => fetchNotifications(token, 50),
  });

  return (
    <div className="px-6 py-8 lg:px-10">
      <PageHeader title="Notifications" icon={Bell} />

      {isLoading ? (
        <div className="flex justify-center py-16">
          <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
        </div>
      ) : notifications.length === 0 ? (
        <Card className="mx-auto max-w-2xl rounded-3xl bg-card-solid">
          <CardContent className="flex flex-col items-center py-16 text-muted-foreground">
            <Bell className="mb-4 h-12 w-12 opacity-20" />
            <p className="font-medium">No notifications yet</p>
            <p className="mt-1 text-sm">When you get notifications, they&apos;ll show up here.</p>
          </CardContent>
        </Card>
      ) : (
        <div className="mx-auto flex max-w-2xl flex-col gap-3">
          {notifications.map((notification) => {
            const message =
              notification.message ||
              String(
                (notification.metadata as Record<string, unknown>)?.message || ""
              ) ||
              "New notification";
            return (
              <Card
                key={String(notification.id)}
                className={`rounded-2xl bg-card-solid ${
                  notification.read
                    ? "border-border"
                    : "border-primary/30 bg-primary/5"
                }`}
              >
                <CardContent className="flex items-start justify-between gap-4 p-4">
                  <p className="text-sm leading-relaxed text-foreground">
                    {message}
                  </p>
                  <span className="shrink-0 text-xs text-muted-foreground">
                    {notification.created_at
                      ? new Date(notification.created_at).toLocaleDateString()
                      : ""}
                  </span>
                </CardContent>
              </Card>
            );
          })}
        </div>
      )}
    </div>
  );
}
