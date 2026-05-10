import { useEffect, useRef } from "react";
import { useToast } from "@/context/ToastContext";
import { buildNotificationsStreamUrl, NotificationItem } from "@/api/business/notifications";
import { useQueryClient } from "@tanstack/react-query";

export function useRealtimeNotifications() {
  const eventSourceRef = useRef<EventSource | null>(null);
  const queryClient = useQueryClient();
  const { addToast } = useToast();

  useEffect(() => {
    const token = localStorage.getItem("token");
    if (!token) return;

    const connect = () => {
      if (eventSourceRef.current) {
        eventSourceRef.current.close();
      }

      const url = buildNotificationsStreamUrl(token);
      const es = new EventSource(url);

      es.onopen = () => {
        console.log("[SSE] Connection established successfully");
      };

      es.onmessage = (event) => {
        try {
          console.log("[SSE] Message received:", event.data);
          const notification = JSON.parse(event.data) as NotificationItem;
          
          const type = String(notification.type || "").toLowerCase();
          const metadata = notification.metadata || {};
          const actorUsername = String(metadata.actor_username || "someone");
          const actorAvatar = metadata.actor_avatar as string;
          
          let title = 'Interaction';
          let message = 'interacted with your story';
          let toastType: any = 'info';

          if (type === 'post_liked' || type === 'like') {
            title = 'New Like';
            message = 'liked your recent post';
            toastType = 'like';
          } else if (type === 'comment_added' || type === 'comment') {
            title = 'Comment';
            message = 'left a comment on your post';
            toastType = 'comment';
          }

          addToast(title, message, toastType, actorAvatar, actorUsername);

          // Update query cache for dropdown/history
          queryClient.invalidateQueries({ queryKey: ['notifications-history'] });
        } catch (err) {
          console.error("Failed to parse notification event", err);
        }
      };

      es.onerror = () => {
        console.error("SSE connection error, retrying in 5s...");
        es.close();
        setTimeout(connect, 5000);
      };

      eventSourceRef.current = es;
    };

    connect();

    return () => {
      eventSourceRef.current?.close();
    };
  }, [queryClient]);
}
