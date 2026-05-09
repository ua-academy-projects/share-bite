import { useState, useEffect, useCallback, useMemo } from "react";
import { AuthContext, type AuthContextType, type JwtPayload } from "@/hooks/useAuth";

const ACCESS_TOKEN_KEY = "sharebite-access-token";
const REFRESH_TOKEN_KEY = "sharebite-refresh-token";

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

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [, setTick] = useState(0);

  useEffect(() => {
    const handler = () => setTick((t) => t + 1);
    window.addEventListener("storage", handler);
    return () => window.removeEventListener("storage", handler);
  }, []);

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
    setTick((t) => t + 1);
  }, []);

  const clearTokens = useCallback(() => {
    localStorage.removeItem(ACCESS_TOKEN_KEY);
    localStorage.removeItem(REFRESH_TOKEN_KEY);
    setTick((t) => t + 1);
  }, []);

  const value: AuthContextType = useMemo(
    () => ({ accessToken, refreshToken, isAuthenticated, userPayload, saveTokens, clearTokens }),
    [accessToken, refreshToken, isAuthenticated, userPayload, saveTokens, clearTokens]
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}
