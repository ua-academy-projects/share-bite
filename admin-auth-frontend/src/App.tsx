import { Routes, Route, Outlet } from "react-router-dom";
import { Sidebar } from "@/components/ui/Sidebar";
import { ProtectedRoute } from "@/components/ProtectedRoute";
import { LoginPage } from "@/pages/LoginPage";
import { RegisterPage } from "@/pages/RegisterPage";
import { RecoverAccessPage } from "@/pages/RecoverAccessPage";
import { ResetPasswordPage } from "@/pages/ResetPasswordPage";
import { DashboardPage } from "@/pages/DashboardPage";
import { LandingPage } from "@/pages/LandingPage";
import { UsersPage } from "@/pages/UsersPage";
import { UserStatusPage } from "@/pages/UserStatusPage";
import { SettingsPage } from "@/pages/SettingsPage";
import { OAuthCallbackPage } from "@/pages/OAuthCallbackPage";

function AdminLayout() {
  return (
    <div className="flex min-h-screen bg-background text-foreground">
      <Sidebar />
      <main className="flex-1">
        <Outlet />
      </main>
    </div>
  );
}

function App() {
  return (
    <Routes>
      {/* Public routes */}
      <Route path="/login" element={<LoginPage />} />
      <Route path="/register" element={<RegisterPage />} />
      <Route path="/recover-access" element={<RecoverAccessPage />} />
      <Route path="/reset-password" element={<ResetPasswordPage />} />
      <Route path="/oauth/:provider/callback" element={<OAuthCallbackPage />} />

      {/* Protected routes with sidebar */}
      <Route element={<ProtectedRoute />}>
        <Route element={<AdminLayout />}>
          <Route path="/" element={<LandingPage />} />
          <Route path="/dashboard" element={<DashboardPage />} />
          <Route path="/users" element={<UsersPage />} />
          <Route path="/users/:userId/status" element={<UserStatusPage />} />
          <Route path="/settings" element={<SettingsPage />} />
        </Route>
      </Route>
    </Routes>
  );
}

export default App;
