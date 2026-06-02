export type NotificationMetadata = Record<string, unknown>;

export type NotificationItem = {
  id: string;
  type: string;
  entity_id: string;
  metadata?: NotificationMetadata;
  created_at: string;
};

const API_BASE_URL = import.meta.env.VITE_NOTIFICATIONS_API_URL || "http://localhost:4005";

function authHeaders(token: string) {
  return {
    Authorization: `Bearer ${token}`,
    "Content-Type": "application/json",
  };
}

export async function fetchNotifications(token: string, limit = 20): Promise<NotificationItem[]> {
  const response = await fetch(`${API_BASE_URL}/notifications?limit=${limit}`, {
    headers: authHeaders(token),
  });

  if (!response.ok) {
    throw new Error(`Failed to load notifications: ${response.status}`);
  }

  return response.json();
}

export async function markNotificationsRead(token: string, notificationIds: string[]): Promise<void> {
  const response = await fetch(`${API_BASE_URL}/notifications/mark-read`, {
    method: "POST",
    headers: authHeaders(token),
    body: JSON.stringify({ notification_ids: notificationIds }),
  });

  if (!response.ok) {
    throw new Error(`Failed to mark notifications read: ${response.status}`);
  }
}

export function buildNotificationsStreamUrl(token: string) {
  return `${API_BASE_URL}/notifications/stream?access_token=${encodeURIComponent(token)}`;
}
