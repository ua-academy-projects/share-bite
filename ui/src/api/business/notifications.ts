export type NotificationMetadata = Record<string, unknown>;

export type NotificationItem = {
  id: string;
  type: string;
  entityID: string;
  metadata?: NotificationMetadata;
  createdAt: string;
};

const API_BASE_URL = "/api/notifications";

function authHeaders(token: string) {
  return {
    Authorization: `Bearer ${token}`,
    "Content-Type": "application/json",
  };
}

export async function fetchNotifications(token: string, limit = 20): Promise<NotificationItem[]> {
  const response = await fetch(`${API_BASE_URL}/?limit=${limit}`, {
    headers: authHeaders(token),
  });

  if (!response.ok) {
    throw new Error(`Failed to load notifications: ${response.status}`);
  }

  return response.json();
}

export async function markNotificationsRead(token: string, notificationIds: string[]): Promise<void> {
  const response = await fetch(`${API_BASE_URL}/mark-read`, {
    method: "POST",
    headers: authHeaders(token),
    body: JSON.stringify({ notificationIDs: notificationIds }),
  });

  if (!response.ok) {
    throw new Error(`Failed to mark notifications read: ${response.status}`);
  }
}

export function buildNotificationsStreamUrl(token: string) {
  const url = `/api/notifications/stream?access_token=${encodeURIComponent(token)}`;
  console.log("[NotificationsAPI] Building stream URL");
  return url;
}
