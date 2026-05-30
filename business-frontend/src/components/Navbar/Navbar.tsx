import React, { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { Moon, Sun, Search, Shield, LogOut, Plus } from "lucide-react";
import { useTheme } from "@/components/theme-provider";
import { isAdminOrModerator } from "@/utils/auth";
import { apiClient } from "@/api/client";
import { useCurrentCustomer } from "@/hooks/useCurrentCustomer";
import { NotificationBell } from "@/components/Notifications/NotificationBell";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";

export const Navbar: React.FC = () => {
  const { theme, setTheme } = useTheme();
  const navigate = useNavigate();
  const { data: currentCustomer } = useCurrentCustomer();
  const [showLogoutDialog, setShowLogoutDialog] = useState(false);
  const [logoutLoading, setLogoutLoading] = useState(false);
  const [logoutError, setLogoutError] = useState("");

  const clearLocalSession = () => {
    localStorage.removeItem("token");
    localStorage.removeItem("refresh_token");
  };

  const handleLogoutCurrentDevice = async () => {
    setLogoutLoading(true);
    setLogoutError("");
    try {
      await apiClient.logout();
    } catch (error: unknown) {
      const err = error as { response?: { data?: { error?: string } }; message?: string };
      setLogoutError(
        err?.response?.data?.error || err?.message || "Failed to logout."
      );
    } finally {
      clearLocalSession();
      navigate("/auth", { replace: true, state: { isLogin: true } });
      setShowLogoutDialog(false);
      setLogoutLoading(false);
    }
  };

  const handleLogoutAllDevices = async () => {
    setLogoutLoading(true);
    setLogoutError("");
    try {
      await apiClient.revokeAllSessions();
      clearLocalSession();
      navigate("/auth", { replace: true, state: { isLogin: true } });
      setShowLogoutDialog(false);
    } catch (error: unknown) {
      const err = error as { response?: { data?: { error?: string } }; message?: string };
      setLogoutError(
        err?.response?.data?.error ||
          err?.message ||
          "Failed to revoke all sessions."
      );
    } finally {
      setLogoutLoading(false);
    }
  };

  const isAuthenticated = !!localStorage.getItem("token");
  const currentDate = new Date()
    .toLocaleDateString("en-US", {
      weekday: "long",
      month: "short",
      day: "numeric",
    })
    .toUpperCase();

  return (
    <>
      <nav
        style={{
          backgroundColor: "var(--navbar-bg)",
          borderColor: "var(--navbar-border)",
        }}
        className="sticky top-0 z-50 w-full border-b px-6 py-3 backdrop-blur-md lg:px-8"
      >
        <div className="mx-auto flex w-full max-w-7xl items-center justify-between">
          <div className="flex items-center gap-6">
            <Link
              to="/"
              className="font-serif text-3xl font-bold tracking-tight"
              style={{ color: "var(--navbar-foreground)" }}
            >
              ShareBite
            </Link>
            <span
              className="hidden rounded-full border px-3 py-1.5 text-[11px] font-black tracking-[0.2em] md:inline-flex"
              style={{
                color: "var(--navbar-muted)",
                backgroundColor: "rgba(170,206,195,0.1)",
                borderColor: "rgba(170,206,195,0.2)",
              }}
            >
              {currentDate}
            </span>
          </div>

          <div className="mx-8 hidden max-w-lg flex-1 lg:flex">
            <div className="relative flex w-full items-center">
              <Search
                className="absolute left-4"
                size={18}
                style={{ color: "var(--navbar-muted)" }}
              />
              <Input
                type="text"
                placeholder="Search restaurants, users..."
                className="h-11 rounded-full border pl-12 shadow-inner"
                style={{
                  backgroundColor: "rgba(255,255,255,0.07)",
                  color: "var(--navbar-foreground)",
                  borderColor: "rgba(170,206,195,0.2)",
                }}
              />
            </div>
          </div>

          <div className="flex items-center gap-4">
            <Button
              variant="ghost"
              size="icon"
              className="rounded-full"
              style={{ color: "var(--navbar-muted)" }}
              onClick={() => setTheme(theme === "dark" ? "light" : "dark")}
            >
              {theme === "light" ? <Moon size={20} /> : <Sun size={20} />}
            </Button>

            {isAuthenticated ? (
              <>
                {isAdminOrModerator() && (
                  <Link
                    to="/admin"
                    className="rounded-full p-2.5 transition-colors"
                    style={{ color: "var(--navbar-muted)" }}
                  >
                    <Shield size={20} />
                  </Link>
                )}

                <Button
                  asChild
                  className="rounded-full bg-accent font-bold text-accent-foreground shadow-lg hover:bg-accent/90"
                >
                  <Link to="/post/create">
                    <Plus size={18} />
                    Post
                  </Link>
                </Button>

                <NotificationBell />

                <div className="mx-2 hidden h-8 w-px bg-border/50 sm:block" />

                <Link
                  to={
                    currentCustomer?.userName
                      ? `/user/${currentCustomer.userName}`
                      : "/profile"
                  }
                  className="group flex items-center gap-3 hover:opacity-80"
                >
                  <img
                    src={
                      currentCustomer?.avatarURL ||
                      "https://via.placeholder.com/40"
                    }
                    alt="Avatar"
                    className="h-10 w-10 rounded-full border-2 border-border object-cover transition-colors group-hover:border-primary"
                  />
                  <div className="hidden flex-col items-start sm:flex">
                    <span
                      className="text-sm font-bold"
                      style={{ color: "var(--navbar-foreground)" }}
                    >
                      @{currentCustomer?.userName || "user"}
                    </span>
                  </div>
                </Link>

                <Button
                  variant="ghost"
                  size="icon"
                  className="rounded-full"
                  style={{ color: "var(--navbar-muted)" }}
                  onClick={() => setShowLogoutDialog(true)}
                  aria-label="Logout"
                >
                  <LogOut size={20} />
                </Button>
              </>
            ) : (
              <div className="flex items-center gap-3">
                <Link
                  to="/auth"
                  state={{ isLogin: true }}
                  className="px-4 text-sm font-bold"
                  style={{ color: "var(--navbar-muted)" }}
                >
                  Log in
                </Link>
                <Link
                  to="/auth"
                  state={{ isLogin: false }}
                  className="rounded-full px-6 py-2.5 text-sm font-bold shadow-md"
                  style={{
                    backgroundColor: "var(--navbar-foreground)",
                    color: "var(--navbar-bg)",
                  }}
                >
                  Sign Up
                </Link>
              </div>
            )}
          </div>
        </div>
      </nav>

      <Dialog open={showLogoutDialog} onOpenChange={setShowLogoutDialog}>
        <DialogContent className="max-w-sm rounded-3xl">
          <DialogHeader>
            <DialogTitle className="font-serif text-2xl">Log out</DialogTitle>
            <DialogDescription>
              Do you want to log out only on this device, or on all devices?
            </DialogDescription>
          </DialogHeader>
          {logoutError && (
            <p className="text-sm text-destructive">{logoutError}</p>
          )}
          <div className="flex flex-col gap-3">
            <Button
              className="w-full rounded-xl bg-accent font-bold text-accent-foreground"
              onClick={handleLogoutAllDevices}
              disabled={logoutLoading}
            >
              {logoutLoading ? "Logging out…" : "All devices"}
            </Button>
            <Button
              variant="secondary"
              className="w-full rounded-xl font-bold"
              onClick={handleLogoutCurrentDevice}
              disabled={logoutLoading}
            >
              This device
            </Button>
          </div>
        </DialogContent>
      </Dialog>
    </>
  );
};
