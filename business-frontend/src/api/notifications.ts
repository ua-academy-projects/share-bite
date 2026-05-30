export type NotificationMetadata = Record<string, unknown>;

export type NotificationItem = {
  id: string;
  type: string;
  entityID: string;
  metadata?: NotificationMetadata;
  isRead: boolean;
  createdAt: string;
  readAt?: string | null;
};

function getNotificationsApiBase(): string {
  const configured = import.meta.env.VITE_NOTIFICATIONS_API_URL;
  if (configured) {
    return configured.replace(/\/$/, "");
  }
  return "/api/notifications";
}

const API_BASE_URL = getNotificationsApiBase();

function authHeaders(token: string) {
  return {
    Authorization: `Bearer ${token}`,
    "Content-Type": "application/json",
  };
}

export async function fetchNotifications(
  token: string,
  limit = 20,
  offset = 0
): Promise<NotificationItem[]> {
  const params = new URLSearchParams({
    limit: String(limit),
    offset: String(offset),
  });
  const response = await fetch(`${API_BASE_URL}/history?${params}`, {
    headers: authHeaders(token),
  });

  if (!response.ok) {
    throw new Error(`Failed to load notifications: ${response.status}`);
  }

  return response.json();
}

export async function markNotificationsRead(
  token: string,
  notificationIds: string[]
): Promise<void> {
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
  return `${API_BASE_URL}/stream?access_token=${encodeURIComponent(token)}`;
}

export function formatNotificationMessage(item: NotificationItem): string {
  const meta = item.metadata ?? {};
  if (typeof meta.message === "string") return meta.message;
  if (typeof meta.content === "string") return meta.content;
  return `New ${item.type} notification`;
}
