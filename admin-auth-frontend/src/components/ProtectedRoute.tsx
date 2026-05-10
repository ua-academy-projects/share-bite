import { Navigate, Outlet, useLocation } from "react-router-dom";
import { useAuth } from "@/hooks/useAuth";

export function ProtectedRoute() {
  const { isAuthenticated } = useAuth();
  const location = useLocation();

  if (!isAuthenticated) {
    const params = new URLSearchParams(location.search);
    const hasOAuthCode = params.has("code");
    const hasOAuthError = params.has("error");
    if (location.pathname === "/" && (hasOAuthCode || hasOAuthError)) {
      return <Navigate to={`/oauth/google/callback${location.search}`} replace />;
    }

    return <Navigate to="/login" replace />;
  }

  return <Outlet />;
}
