export type NotificationMetadata = Record<string, unknown>;

export type NotificationItem = {
  id: string;
  type: string;
  entity_id: string;
  metadata?: NotificationMetadata;
  created_at: string;
  read?: boolean;
  message?: string;
};

const API_BASE = import.meta.env.VITE_API_BASE_URL || "";

function authHeaders(token: string) {
  return {
    Authorization: `Bearer ${token}`,
    "Content-Type": "application/json",
  };
}

export async function fetchNotifications(
  token: string,
  limit = 20
): Promise<NotificationItem[]> {
  const response = await fetch(
    `${API_BASE}/api/notifications/history?limit=${limit}`,
    { headers: authHeaders(token) }
  );

  if (!response.ok) {
    throw new Error(`Failed to load notifications: ${response.status}`);
  }

  const data = await response.json();
  return Array.isArray(data) ? data : data.items || [];
}

export async function markNotificationsRead(
  token: string,
  notificationIds: string[]
): Promise<void> {
  const response = await fetch(`${API_BASE}/api/notifications/mark-read`, {
    method: "POST",
    headers: authHeaders(token),
    body: JSON.stringify({ notification_ids: notificationIds }),
  });

  if (!response.ok) {
    throw new Error(`Failed to mark notifications read: ${response.status}`);
  }
}

export function buildNotificationsStreamUrl(token: string) {
  return `${API_BASE}/api/notifications/stream?access_token=${encodeURIComponent(token)}`;
}
