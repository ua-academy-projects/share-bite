import { Navigate, Routes, Route } from "react-router-dom";
import { AppShell } from "@/components/AppShell";
import { RequireAuth } from "@/components/RequireAuth/RequireAuth";
import { RequireAdmin } from "@/components/RequireAdmin/RequireAdmin";
import { RoleBasedHome } from "@/components/RoleBasedHome";
import { HomeFeedPage } from "@/pages/HomeFeedPage";
import { QRCodeModalProvider } from "@/contexts/QRCodeModalContext";
import { QRCodeModalContainer } from "@/components/ui/QRCodeModal";
import { BoxesPage } from "@/pages/BoxesPage";
import CreatePostPage from "@/pages/CreatePostPage";
import { CreateBoxPage } from "@/pages/CreateBoxPage";
import { VenueSearchPage } from "@/pages/VenueSearchPage";
import { VenueProfilePage } from "@/pages/VenueProfilePage";
import { Auth } from "@/pages/guest/Auth/Auth";
import { OAuthCallback } from "@/pages/guest/OAuthCallback/OAuthCallback";
import { GitHubSuccess } from "@/pages/guest/GitHubSuccess/GitHubSuccess";
import { HomeFeed } from "@/pages/guest/HomeFeed/HomeFeed";
import { ExplorePage } from "@/pages/guest/Explore/ExplorePage";
import { CollectionsPage } from "@/pages/guest/Collections/CollectionsPage";
import { NotificationsPage } from "@/pages/guest/Notifications/NotificationsPage";
import { SecurityPage } from "@/pages/guest/Settings/SecurityPage";
import { UserProfile } from "@/pages/guest/UserProfile/UserProfile";
import { ProfileCreatePage } from "@/pages/guest/UserProfile/ProfileCreatePage";
import { ProfileEditPage } from "@/pages/guest/UserProfile/ProfileEditPage";
import { RestaurantProfile } from "@/pages/guest/RestaurantProfile/RestaurantProfile";
import { CreatePost } from "@/pages/guest/CreatePost/CreatePost";
import { AdminUsersPage } from "@/pages/guest/Admin/AdminUsersPage";
import { AdminUserDetailPage } from "@/pages/guest/Admin/AdminUserDetailPage";

function App() {
  return (
    <QRCodeModalProvider>
      <AppShell>
        <Routes>
          <Route path="/" element={<RoleBasedHome />} />
          <Route path="/boxes" element={<BoxesPage />} />
          <Route path="/discover" element={<VenueSearchPage />} />
          <Route path="/venues/search" element={<VenueSearchPage />} />
          <Route path="/venue/:id/create-post" element={<CreatePostPage />} />
          <Route path="/venue/:id/create-box" element={<CreateBoxPage />} />
          <Route path="/venue/:id" element={<VenueProfilePage />} />

          <Route path="/auth" element={<Auth />} />
          <Route path="/oauth/google/callback" element={<OAuthCallback />} />
          <Route path="/oauth/github/success" element={<GitHubSuccess />} />

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
            element={
              <RequireAuth>
                <ProfileEditPage />
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
          <Route path="/restaurant/:id" element={<RestaurantProfile />} />

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
        </Routes>
      </AppShell>
      <QRCodeModalContainer />
    </QRCodeModalProvider>
  );
}

export default App;
