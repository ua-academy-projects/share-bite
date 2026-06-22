import { useEffect } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { buildNotificationsStreamUrl } from "@/api/notifications";

/**
 * Opens an SSE connection to the notifications stream and refreshes the
 * ["notifications"] query cache whenever a new notification is pushed, so the
 * bell badge and the notifications list update live — no page reload needed.
 *
 * Mount once for an authenticated session (e.g. in the Sidebar).
 */
export function useRealtimeNotifications() {
  const queryClient = useQueryClient();
  const token = localStorage.getItem("token");

  useEffect(() => {
    if (!token) return;

    let es: EventSource | null = null;
    let reconnectTimer: number | null = null;
    let closed = false;

    const connect = () => {
      es = new EventSource(buildNotificationsStreamUrl(token));

      // Default ("message") events carry notifications; "ping" heartbeats are
      // named events and intentionally ignored here.
      es.onmessage = () => {
        queryClient.invalidateQueries({ queryKey: ["notifications"] });
      };

      es.onerror = () => {
        es?.close();
        if (closed || reconnectTimer) return;
        reconnectTimer = window.setTimeout(() => {
          reconnectTimer = null;
          connect();
        }, 5000);
      };
    };

    connect();

    return () => {
      closed = true;
      if (reconnectTimer) clearTimeout(reconnectTimer);
      es?.close();
    };
  }, [queryClient, token]);
}
