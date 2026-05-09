const API_BASE_URL = import.meta.env.VITE_API_URL || "http://localhost:3850";

export type TokensResponse = {
  access_token: string;
  refresh_token: string;
};

export type MessageResponse = {
  message: string;
};

export type ErrorResponse = {
  error: string;
};

export type UserStatusResponse = {
  status: "active" | "muted" | "suspended";
};

export type UserListItem = {
  id: string;
  email: string;
  status: string;
  created_at: string;
  updated_at: string;
};

async function request<T>(
  path: string,
  options: RequestInit = {}
): Promise<T> {
  const res = await fetch(`${API_BASE_URL}${path}`, {
    headers: {
      "Content-Type": "application/json",
      ...options.headers,
    },
    ...options,
  });

  if (!res.ok) {
    const body = await res.json().catch(() => ({}));
    throw new Error(body.error || `Request failed (${res.status})`);
  }

  return res.json();
}

function authHeaders(token: string): HeadersInit {
  return { Authorization: `Bearer ${token}` };
}

export const authApi = {
  login: (email: string, password: string) =>
    request<TokensResponse>("/auth/login", {
      method: "POST",
      body: JSON.stringify({ email, password }),
    }),

  register: (email: string, password: string, slug: string) =>
    request<TokensResponse>("/auth/register", {
      method: "POST",
      body: JSON.stringify({ email, password, slug }),
    }),

  refresh: (refreshToken: string) =>
    request<TokensResponse>("/auth/refresh", {
      method: "POST",
      body: JSON.stringify({ refresh_token: refreshToken }),
    }),

  recoverAccess: (email: string) =>
    request<MessageResponse>("/auth/recover-access", {
      method: "POST",
      body: JSON.stringify({ email }),
    }),

  resetPassword: (token: string, newPassword: string) =>
    request<MessageResponse>("/auth/reset-password", {
      method: "POST",
      body: JSON.stringify({ token, new_password: newPassword }),
    }),

  oauthCallback: (provider: string, code: string, slug: string) =>
    request<TokensResponse>(`/auth/oauth/${provider}/callback`, {
      method: "POST",
      body: JSON.stringify({ code, slug }),
    }),

  logout: (token: string, refreshToken: string) =>
    request<MessageResponse>("/user/logout", {
      method: "POST",
      headers: authHeaders(token),
      body: JSON.stringify({ refresh_token: refreshToken }),
    }),

  revokeAllSessions: (token: string) =>
    request<MessageResponse>("/user/sessions/revoke-all", {
      method: "POST",
      headers: authHeaders(token),
    }),

  linkProvider: (token: string, provider: string, code: string) =>
    request<MessageResponse>(`/user/link/${provider}`, {
      method: "POST",
      headers: authHeaders(token),
      body: JSON.stringify({ code }),
    }),

  getUserStatus: (token: string, userId: string) =>
    request<UserStatusResponse>(`/users/${userId}/status`, {
      headers: authHeaders(token),
    }),

  updateUserStatus: (
    token: string,
    userId: string,
    status: string
  ) =>
    request<MessageResponse>(`/users/${userId}/status`, {
      method: "PUT",
      headers: authHeaders(token),
      body: JSON.stringify({ status }),
    }),

  listUsers: (token: string) =>
    request<UserListItem[]>("/users", {
      headers: authHeaders(token),
    }),

  getGitHubLoginUrl: () => `${API_BASE_URL}/auth/github`,
};
