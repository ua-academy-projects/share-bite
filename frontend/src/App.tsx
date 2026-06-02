import { Routes, Route, Navigate, useLocation } from 'react-router-dom';
import { Navbar } from './components/Navbar/Navbar';
import { HomeFeed } from './pages/HomeFeed/HomeFeed';
import { RestaurantProfile } from './pages/RestaurantProfile/RestaurantProfile';
import { UserProfile } from './pages/UserProfile/UserProfile';
import { Auth } from './pages/Auth/Auth';
import { CreatePost } from './pages/CreatePost/CreatePost';
import { OAuthCallback } from './pages/OAuthCallback/OAuthCallback';
import { GitHubSuccess } from './pages/GitHubSuccess/GitHubSuccess';
import { AdminUsersPage } from './pages/Admin/AdminUsersPage';
import { AdminUserDetailPage } from './pages/Admin/AdminUserDetailPage';
import { ExplorePage } from './pages/Explore/ExplorePage';
import { CollectionsPage } from './pages/Collections/CollectionsPage';
import { NotificationsPage } from './pages/Notifications/NotificationsPage';
import { SecurityPage } from './pages/Settings/SecurityPage';
import { RequireAuth } from './components/RequireAuth/RequireAuth';
import { RequireAdmin } from './components/RequireAdmin/RequireAdmin';
import { Toaster } from './components/ui/sonner';


function HomeEntry() {
  const location = useLocation();
  const params = new URLSearchParams(location.search);
  const hasOAuthPayload = params.has('code') || params.has('error');

  if (hasOAuthPayload) {
    return <Navigate to={`/oauth/google/callback${location.search}`} replace />;
  }

  return <HomeFeed />;
}

function App() {
  return (
    <>
      <div className="flex min-h-screen bg-background text-foreground flex-col">
      <Navbar />
        <main className="flex-1 flex flex-col">
          <Routes>
            <Route path="/" element={<HomeEntry />} />
            <Route path="/explore" element={<ExplorePage />} />
            <Route path="/collections" element={<RequireAuth><CollectionsPage /></RequireAuth>} />
            <Route path="/notifications" element={<RequireAuth><NotificationsPage /></RequireAuth>} />
            <Route path="/settings/security" element={<RequireAuth><SecurityPage /></RequireAuth>} />
            <Route path="/restaurant/:id" element={<RestaurantProfile />} />
            <Route path="/profile" element={<RequireAuth><UserProfile /></RequireAuth>} />
            <Route path="/user/:id" element={<RequireAuth><UserProfile /></RequireAuth>} />
            <Route path="/auth" element={<Auth />} />
            <Route path="/post/create" element={<RequireAuth><CreatePost /></RequireAuth>} />
            <Route path="/oauth/google/callback" element={<OAuthCallback />} />
            <Route path="/oauth/github/success" element={<GitHubSuccess />} />
            <Route path="/admin" element={<RequireAdmin><AdminUsersPage /></RequireAdmin>} />
            <Route path="/admin/users/:id" element={<RequireAdmin><AdminUserDetailPage /></RequireAdmin>} />
          </Routes>
        </main>
      </div>
      <Toaster />
    </>
  );
}

export default App;
