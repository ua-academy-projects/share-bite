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
import { RequireAuth } from './components/RequireAuth/RequireAuth';
import { RequireAdmin } from './components/RequireAdmin/RequireAdmin';
import { ThemeProvider } from './context/ThemeContext';
import './styles/variables.css';

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
    <ThemeProvider>
      <div className="app-container">
        <Navbar />
        <main className="main-content">
          <Routes>
            <Route path="/" element={<HomeEntry />} />
            <Route path="/explore" element={<Navigate to="/" replace />} />
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
    </ThemeProvider>
  );
}

export default App;
