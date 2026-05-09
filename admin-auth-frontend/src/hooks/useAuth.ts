import { createContext, useContext, useCallback, useMemo } from "react";

const ACCESS_TOKEN_KEY = "sharebite-access-token";
const REFRESH_TOKEN_KEY = "sharebite-refresh-token";

export type AuthContextType = {
  accessToken: string | null;
  refreshToken: string | null;
  isAuthenticated: boolean;
  userPayload: JwtPayload | null;
  saveTokens: (accessToken: string, refreshToken: string) => void;
  clearTokens: () => void;
};

export type JwtPayload = {
  sub: string;
  role: string;
  status: string;
  exp: number;
  name?: string;
  email?: string;
  preferred_username?: string;
  given_name?: string;
  family_name?: string;
};

function parseJwt(token: string): JwtPayload | null {
  try {
    const base64 = token.split(".")[1];
    const json = atob(base64.replace(/-/g, "+").replace(/_/g, "/"));
    return JSON.parse(json);
  } catch {
    return null;
  }
}

function isTokenExpired(payload: JwtPayload): boolean {
  return Date.now() >= payload.exp * 1000;
}

export const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function useAuthValue(): AuthContextType {
  const accessToken = localStorage.getItem(ACCESS_TOKEN_KEY);
  const refreshToken = localStorage.getItem(REFRESH_TOKEN_KEY);

  const userPayload = useMemo(() => {
    if (!accessToken) return null;
    const payload = parseJwt(accessToken);
    if (!payload || isTokenExpired(payload)) return null;
    return payload;
  }, [accessToken]);

  const isAuthenticated = !!userPayload;

  const saveTokens = useCallback((at: string, rt: string) => {
    localStorage.setItem(ACCESS_TOKEN_KEY, at);
    localStorage.setItem(REFRESH_TOKEN_KEY, rt);
    window.dispatchEvent(new Event("storage"));
  }, []);

  const clearTokens = useCallback(() => {
    localStorage.removeItem(ACCESS_TOKEN_KEY);
    localStorage.removeItem(REFRESH_TOKEN_KEY);
    window.dispatchEvent(new Event("storage"));
  }, []);

  return { accessToken, refreshToken, isAuthenticated, userPayload, saveTokens, clearTokens };
}

export function useAuth(): AuthContextType {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be used within AuthProvider");
  return ctx;
}
