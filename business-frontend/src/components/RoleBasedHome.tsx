import { Navigate, useLocation } from "react-router-dom";
import { getTokenRole, isAdminOrModerator } from "@/utils/auth";
import { HomeFeed } from "@/pages/guest/HomeFeed/HomeFeed";
import { HomeFeedPage } from "@/pages/HomeFeedPage";

function HomeEntry() {
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

export function RoleBasedHome() {
  const role = getTokenRole();

  if (isAdminOrModerator()) {
    return <Navigate to="/admin" replace />;
  }

  if (role === "business") {
    return <HomeFeedPage />;
  }

  return <HomeEntry />;
}
