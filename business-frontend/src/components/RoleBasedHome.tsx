import { Navigate, useLocation } from "react-router-dom";
import { HomeFeed } from "@/pages/guest/HomeFeed/HomeFeed";

export function RoleBasedHome() {
  const location = useLocation();
  const params = new URLSearchParams(location.search);
  const hasOAuthPayload = params.has("code") || params.has("error");

  if (hasOAuthPayload) {
    return (
      <Navigate to={`/oauth/google/callback${location.search}`} replace />
    );
  }

  return <HomeFeed />;
}
