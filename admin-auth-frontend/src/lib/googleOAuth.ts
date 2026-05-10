const DEFAULT_GOOGLE_REDIRECT_URI = "http://localhost:5147/";
const GOOGLE_OAUTH_URL = "https://accounts.google.com/o/oauth2/v2/auth";

export type OAuthRoleSlug = "user" | "business";

export function isOAuthRoleSlug(value: string): value is OAuthRoleSlug {
  return value === "user" || value === "business";
}

export function getGoogleOAuthRedirectUri(): string {
  return import.meta.env.VITE_GOOGLE_REDIRECT_URI || DEFAULT_GOOGLE_REDIRECT_URI;
}

function encodeState(payload: Record<string, string>): string {
  const json = JSON.stringify(payload);
  const bytes = new TextEncoder().encode(json);
  const binary = Array.from(bytes, (byte) => String.fromCharCode(byte)).join("");
  return btoa(binary);
}

function decodeState(encoded: string): Record<string, string> | null {
  try {
    const binary = atob(encoded);
    const bytes = Uint8Array.from(binary, (char) => char.charCodeAt(0));
    const decoded = new TextDecoder().decode(bytes);
    const payload = JSON.parse(decoded) as Record<string, unknown>;
    if (!payload || typeof payload !== "object") {
      return null;
    }

    return Object.entries(payload).reduce<Record<string, string>>((acc, [key, value]) => {
      if (typeof value === "string") {
        acc[key] = value;
      }
      return acc;
    }, {});
  } catch {
    return null;
  }
}

export function buildGoogleAuthorizationUrl(clientId: string, slug: OAuthRoleSlug): string {
  const redirectUri = getGoogleOAuthRedirectUri();
  const scope = "email profile";
  const params = new URLSearchParams({
    client_id: clientId,
    redirect_uri: redirectUri,
    response_type: "code",
    scope,
  });

  if (slug) {
    params.set("state", encodeState({ slug }));
  }

  return `${GOOGLE_OAUTH_URL}?${params.toString()}`;
}

export function parseGoogleOAuthState(state: string | null): OAuthRoleSlug | null {
  if (!state) {
    return null;
  }

  if (isOAuthRoleSlug(state)) {
    return state;
  }

  const parsed = decodeState(state);
  const slug = parsed?.slug;
  if (!slug || !isOAuthRoleSlug(slug)) {
    return null;
  }
  return slug;
}

/**
 * Accepts a full redirect URL (e.g. after Google sends you to localhost:5147/?code=...)
 * or a raw code string. Returns the decoded authorization code (e.g. starts with "4/").
 */
export function parseGoogleAuthorizationCode(input: string): string {
  const trimmed = input.trim();
  if (!trimmed) return "";

  const tryParseAsUrl = (href: string): string | null => {
    try {
      const u = new URL(href);
      const c = u.searchParams.get("code");
      return c ? decodeURIComponent(c) : null;
    } catch {
      return null;
    }
  };

  if (trimmed.includes("code=")) {
    if (trimmed.startsWith("http://") || trimmed.startsWith("https://")) {
      const fromFull = tryParseAsUrl(trimmed);
      if (fromFull) return fromFull;
    }
    // Paste like "localhost:5147/?code=..." without scheme
    if (trimmed.startsWith("localhost")) {
      const withScheme = `http://${trimmed.replace(/^\/*/, "")}`;
      const fromLocal = tryParseAsUrl(withScheme);
      if (fromLocal) return fromLocal;
    }
    const qs = trimmed.includes("?") ? trimmed.slice(trimmed.indexOf("?")) : `?${trimmed}`;
    const fromQs = tryParseAsUrl(`http://local.invalid${qs}`);
    if (fromQs) return fromQs;
  }

  try {
    return decodeURIComponent(trimmed);
  } catch {
    return trimmed;
  }
}
