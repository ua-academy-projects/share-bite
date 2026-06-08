import { Navigate, Routes, Route, useParams } from "react-router-dom";
import { AppShell } from "@/components/AppShell";
import { RequireAuth } from "@/components/RequireAuth/RequireAuth";
import { RequireAdmin } from "@/components/RequireAdmin/RequireAdmin";
import { RoleBasedHome } from "@/components/RoleBasedHome";
import { HomeFeedPage } from "@/pages/HomeFeedPage";
import { CreateHubPage } from "@/pages/CreateHubPage";
import { BusinessAccountPage } from "@/pages/BusinessAccountPage";
import { DiscoverPage } from "@/pages/DiscoverPage";
import { MyVenuesPage } from "@/pages/MyVenuesPage";
import { QRCodeModalProvider } from "@/contexts/QRCodeModalContext";
import { QRCodeModalContainer } from "@/components/ui/QRCodeModal";
import { BoxesPage } from "@/pages/BoxesPage";
import CreatePostPage from "@/pages/CreatePostPage";
import { CreateBoxPage } from "@/pages/CreateBoxPage";
import { VenueSearchPage } from "@/pages/VenueSearchPage";
import { VenueProfilePage } from "@/pages/VenueProfilePage";
import { BusinessSetupPage } from "@/pages/business/BusinessSetupPage";
import { Auth } from "@/pages/guest/Auth/Auth";
import { OAuthCallback } from "@/pages/guest/OAuthCallback/OAuthCallback";
import { GitHubSuccess } from "@/pages/guest/GitHubSuccess/GitHubSuccess";
import { HomeFeed } from "@/pages/guest/HomeFeed/HomeFeed";
import { CollectionsPage } from "@/pages/guest/Collections/CollectionsPage";
import { NotificationsPage } from "@/pages/guest/Notifications/NotificationsPage";
import { AccountSettingsPage } from "@/pages/guest/Settings/AccountSettingsPage";
import { UserProfile } from "@/pages/guest/UserProfile/UserProfile";
import { ProfileCreatePage } from "@/pages/guest/UserProfile/ProfileCreatePage";
import { CreatePost } from "@/pages/guest/CreatePost/CreatePost";
import { AdminUsersPage } from "@/pages/guest/Admin/AdminUsersPage";
import { AdminUserDetailPage } from "@/pages/guest/Admin/AdminUserDetailPage";
import { ForbiddenPage } from "@/pages/guest/Forbidden/ForbiddenPage";
import { isUserRole } from "@/utils/auth";

function RestaurantRedirect() {
  const { id } = useParams();
  return <Navigate to={`/venue/${id}`} replace />;
}

function RequireGuestUser({ children }: { children: React.ReactNode }) {
  const token = localStorage.getItem("token");
  if (!token) {
    return <Navigate to="/auth" replace />;
  }
  if (!isUserRole()) {
    return <Navigate to="/forbidden" replace />;
  }
  return <>{children}</>;
}

function App() {
  return (
    <QRCodeModalProvider>
      <AppShell>
        <Routes>
          <Route path="/" element={<RoleBasedHome />} />
          <Route path="/boxes" element={<BoxesPage />} />
          <Route path="/create" element={<CreateHubPage />} />
          <Route path="/account" element={<BusinessAccountPage />} />
          <Route path="/discover" element={<DiscoverPage />} />
          <Route path="/venues/mine" element={<MyVenuesPage />} />
          <Route path="/venues/search" element={<VenueSearchPage />} />
          <Route path="/explore" element={<Navigate to="/discover" replace />} />
          <Route
            path="/venue/:id/create-post"
            element={
              <RequireAuth>
                <CreatePostPage />
              </RequireAuth>
            }
          />
          <Route path="/venue/:id/create-box" element={<CreateBoxPage />} />
          <Route path="/venue/:id" element={<VenueProfilePage />} />

          <Route path="/auth" element={<Auth />} />
          <Route path="/oauth/google/callback" element={<OAuthCallback />} />
          <Route path="/oauth/github/success" element={<GitHubSuccess />} />

          <Route
            path="/collections"
            element={
              <RequireGuestUser>
                <CollectionsPage />
              </RequireGuestUser>
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
            path="/settings/account"
            element={
              <RequireAuth>
                <AccountSettingsPage />
              </RequireAuth>
            }
          />
          <Route
            path="/settings/security"
            element={<Navigate to="/settings/account" replace />}
          />
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
                <ProfileCreatePage />
              </RequireAuth>
            }
          />
          <Route
            path="/profile/edit"
            element={<Navigate to="/settings/account" replace />}
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
          <Route
            path="/venue/:id/post/create"
            element={
              <RequireAuth>
                <CreatePost />
              </RequireAuth>
            }
          />
          <Route path="/restaurant/:id" element={<RestaurantRedirect />} />

          <Route path="/forbidden" element={<ForbiddenPage />} />
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

          <Route path="/feed/users" element={<HomeFeed />} />
          <Route path="/feed/business" element={<HomeFeedPage />} />
          <Route path="/feed" element={<Navigate to="/feed/users" replace />} />

          <Route
            path="/business/setup"
            element={
              <RequireAuth>
                <BusinessSetupPage />
              </RequireAuth>
            }
          />
        </Routes>
      </AppShell>
      <QRCodeModalContainer />
    </QRCodeModalProvider>
  );
}

export default App;
