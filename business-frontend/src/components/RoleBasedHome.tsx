import { Navigate, useLocation } from "react-router-dom";
import { isAdminOrModerator, isBusinessRole } from "@/utils/auth";
import { HomeFeedPage } from "@/pages/HomeFeedPage";
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

  if (isAdminOrModerator()) {
    return <Navigate to="/admin" replace />;
  }

  if (isBusinessRole()) {
    return <HomeFeedPage />;
  }

  return <HomeFeed />;
}
