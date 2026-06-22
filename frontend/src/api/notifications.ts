export type NotificationMetadata = Record<string, unknown>;

export type NotificationItem = {
  id: string;
  type: string;
  entityID: string;
  metadata?: NotificationMetadata;
  createdAt: string;
  read?: boolean;
  message?: string;
};

const API_BASE_URL = "/api/notifications";

function authHeaders(token: string) {
  return {
    Authorization: `Bearer ${token}`,
    "Content-Type": "application/json",
  };
}

function mapNotification(raw: Record<string, unknown>): NotificationItem {
  const metadata = raw.metadata as NotificationMetadata | undefined;
  const message =
    (raw.message as string | undefined) ||
    (metadata && typeof metadata.message === "string" ? metadata.message : undefined);

  return {
    id: String(raw.id ?? ""),
    type: String(raw.type ?? ""),
    entityID: String(raw.entityID ?? raw.entity_id ?? ""),
    metadata,
    createdAt: String(raw.createdAt ?? raw.created_at ?? new Date().toISOString()),
    read: Boolean(raw.read),
    message,
  };
}

export async function fetchNotifications(
  token: string,
  limit = 20
): Promise<NotificationItem[]> {
  const response = await fetch(`${API_BASE_URL}/history?limit=${limit}`, {
    headers: authHeaders(token),
  });

  if (!response.ok) {
    throw new Error(`Failed to load notifications: ${response.status}`);
  }

  const data = await response.json();
  const items = Array.isArray(data) ? data : data.items || [];
  return items.map((item: Record<string, unknown>) => mapNotification(item));
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

export async function fetchNotificationPreferences(
  token: string
): Promise<Record<string, boolean>> {
  const response = await fetch(`${API_BASE_URL}/preferences`, {
    headers: authHeaders(token),
  });

  if (!response.ok) {
    throw new Error(`Failed to load notification preferences: ${response.status}`);
  }

  return response.json();
}

export async function updateNotificationPreferences(
  token: string,
  preferences: Record<string, boolean>
): Promise<void> {
  const response = await fetch(`${API_BASE_URL}/preferences`, {
    method: "PUT",
    headers: authHeaders(token),
    body: JSON.stringify(preferences),
  });

  if (!response.ok) {
    throw new Error(`Failed to update notification preferences: ${response.status}`);
  }
}

export function buildNotificationsStreamUrl(token: string) {
  return `${API_BASE_URL}/stream?access_token=${encodeURIComponent(token)}`;
}
