import { Navigate, Route, Routes, useLocation } from "react-router-dom";
import { AppShell } from "@/components/AppShell";
import { RequireAuth } from "@/components/RequireAuth/RequireAuth";
import { RequireAdmin } from "@/components/RequireAdmin/RequireAdmin";
import { QRCodeModalProvider } from "@/contexts/QRCodeModalContext";
import { QRCodeModalContainer } from "@/components/ui/QRCodeModal";
import { getTokenRole } from "@/utils/auth";
import { BoxesPage } from "@/pages/BoxesPage";
import CreatePostPage from "@/pages/CreatePostPage";
import { CreateBoxPage } from "@/pages/CreateBoxPage";
import { VenueSearchPage } from "@/pages/VenueSearchPage";
import { VenueProfilePage } from "@/pages/VenueProfilePage";
import { HomeFeedPage } from "@/pages/HomeFeedPage";
import { Auth } from "@/pages/guest/Auth/Auth";
import { HomeFeed } from "@/pages/guest/HomeFeed/HomeFeed";
import { ExplorePage } from "@/pages/guest/Explore/ExplorePage";
import { CollectionsPage } from "@/pages/guest/Collections/CollectionsPage";
import { NotificationsPage } from "@/pages/guest/Notifications/NotificationsPage";
import { SecurityPage } from "@/pages/guest/Settings/SecurityPage";
import { UserProfile } from "@/pages/guest/UserProfile/UserProfile";
import { CreatePost } from "@/pages/guest/CreatePost/CreatePost";
import { RestaurantProfile } from "@/pages/guest/RestaurantProfile/RestaurantProfile";
import { OAuthCallback } from "@/pages/guest/OAuthCallback/OAuthCallback";
import { GitHubSuccess } from "@/pages/guest/GitHubSuccess/GitHubSuccess";
import { AdminUsersPage } from "@/pages/guest/Admin/AdminUsersPage";
import { AdminUserDetailPage } from "@/pages/guest/Admin/AdminUserDetailPage";

function RoleBasedHome() {
  const location = useLocation();
  const params = new URLSearchParams(location.search);
  if (params.has("code") || params.has("error")) {
    return <Navigate to={`/oauth/google/callback${location.search}`} replace />;
  }

  const role = getTokenRole();
  if (role === "admin" || role === "moderator") {
    return <Navigate to="/admin" replace />;
  }
  if (role === "business") {
    return <HomeFeedPage />;
  }
  return <HomeFeed />;
}

function App() {
  return (
    <QRCodeModalProvider>
      <AppShell>
        <Routes>
          <Route path="/auth" element={<Auth />} />
          <Route path="/oauth/google/callback" element={<OAuthCallback />} />
          <Route path="/oauth/github/success" element={<GitHubSuccess />} />

          <Route path="/" element={<RoleBasedHome />} />
          <Route path="/explore" element={<ExplorePage />} />
          <Route
            path="/collections"
            element={
              <RequireAuth>
                <CollectionsPage />
              </RequireAuth>
            }
          />
          <Route
            path="/notifications"
            element={
              <RequireAuth>
                <NotificationsPage />
              </RequireAuth>
            }
          />
          <Route
            path="/settings/security"
            element={
              <RequireAuth>
                <SecurityPage />
              </RequireAuth>
            }
          />
          <Route path="/restaurant/:id" element={<RestaurantProfile />} />
          <Route
            path="/profile"
            element={
              <RequireAuth>
                <UserProfile />
              </RequireAuth>
            }
          />
          <Route
            path="/profile/create"
            element={
              <RequireAuth>
                <UserProfile mode="create" />
              </RequireAuth>
            }
          />
          <Route
            path="/profile/edit"
            element={
              <RequireAuth>
                <UserProfile mode="edit" />
              </RequireAuth>
            }
          />
          <Route
            path="/user/:id"
            element={
              <RequireAuth>
                <UserProfile />
              </RequireAuth>
            }
          />
          <Route
            path="/post/create"
            element={
              <RequireAuth>
                <CreatePost />
              </RequireAuth>
            }
          />

          <Route path="/boxes" element={<BoxesPage />} />
          <Route path="/discover" element={<VenueSearchPage />} />
          <Route path="/venues/search" element={<VenueSearchPage />} />
          <Route path="/venue/:id/create-post" element={<CreatePostPage />} />
          <Route path="/venue/:id/create-box" element={<CreateBoxPage />} />
          <Route path="/venue/:id" element={<VenueProfilePage />} />

          <Route
            path="/admin"
            element={
              <RequireAdmin>
                <AdminUsersPage />
              </RequireAdmin>
            }
          />
          <Route
            path="/admin/users/:id"
            element={
              <RequireAdmin>
                <AdminUserDetailPage />
              </RequireAdmin>
            }
          />
        </Routes>
      </AppShell>
      <QRCodeModalContainer />
    </QRCodeModalProvider>
  );
}

export default App;
