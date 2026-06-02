export interface JwtPayload {
  sub: string;
  role: string;
  status: string;
  exp: number;
  email?: string;
  name?: string;
}

export function parseJwt(token: string): JwtPayload | null {
  const parts = token.split('.');
  if (parts.length !== 3 || parts[1] == null || parts[1] === '') {
    return null;
  }

  try {
    const base64 = parts[1];
    const json = atob(base64.replace(/-/g, '+').replace(/_/g, '/'));
    return JSON.parse(json) as JwtPayload;
  } catch {
    return null;
  }
}

export function getTokenRole(): string | null {
  const token = localStorage.getItem('token');
  if (!token) return null;
  const payload = parseJwt(token);
  if (!payload) return null;
  if (Date.now() >= payload.exp * 1000) return null;
  return payload.role;
}

export function isAdminOrModerator(): boolean {
  const role = getTokenRole();
  return role === 'admin' || role === 'moderator';
}
