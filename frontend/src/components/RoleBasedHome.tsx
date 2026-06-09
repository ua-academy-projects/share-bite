import { Navigate, useLocation } from "react-router-dom";
import { getDefaultHomePath } from "@/utils/navigation";

export function RoleBasedHome() {
  const location = useLocation();
  const params = new URLSearchParams(location.search);
  const hasOAuthPayload = params.has("code") || params.has("error");

  if (hasOAuthPayload) {
    return (
      <Navigate to={`/oauth/google/callback${location.search}`} replace />
    );
  }

  return <Navigate to={getDefaultHomePath()} replace />;
}
