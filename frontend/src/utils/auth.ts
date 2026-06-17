export interface JwtPayload {
  sub: string;
  role: string;
  status: string;
  exp: number;
  email?: string;
  name?: string;
}

export function parseJwt(token: string): JwtPayload | null {
  const parts = token.split(".");
  if (parts.length !== 3 || parts[1] == null || parts[1] === "") {
    return null;
  }

  try {
    let base64 = parts[1].replace(/-/g, "+").replace(/_/g, "/");
    const pad = (4 - (base64.length % 4)) % 4;
    base64 += "=".repeat(pad);
    const json = atob(base64);
    return JSON.parse(json) as JwtPayload;
  } catch {
    return null;
  }
}

export function getTokenRole(): string | null {
  const payload = getTokenPayload();
  return payload?.role ?? null;
}

export function getTokenPayload(): JwtPayload | null {
  const token = localStorage.getItem("token");
  if (!token) return null;
  const payload = parseJwt(token);
  if (!payload) return null;
  if (Date.now() >= payload.exp * 1000) return null;
  return payload;
}

export function isAdminOrModerator(): boolean {
  const role = getTokenRole();
  return role === "admin" || role === "moderator";
}

export function isBusinessRole(): boolean {
  const role = getTokenRole();
  return role === "business";
}

export function isUserRole(): boolean {
  const role = getTokenRole();
  return role === "user";
}

export function setBusinessOrgId(venueId: number) {
  localStorage.setItem("business_org_id", String(venueId));
}

export function getBusinessOrgId(): number | null {
  const raw = localStorage.getItem("business_org_id");
  if (!raw) return null;
  const id = Number(raw);
  return Number.isFinite(id) && id > 0 ? id : null;
}

export function clearSessionStorage() {
  localStorage.removeItem("token");
  localStorage.removeItem("refresh_token");
  localStorage.removeItem("guest_has_customer");
  localStorage.removeItem("business_org_id");
}
