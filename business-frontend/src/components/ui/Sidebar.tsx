import { useState } from "react";
import { Link, NavLink } from "react-router-dom";
import { Moon, Sun, LogOut } from "lucide-react";
import { Button } from "@/components/ui/button";
import { useTheme } from "@/components/theme-provider";
import { useCurrentCustomer } from "@/hooks/useCurrentCustomer";
import { NotificationBell } from "@/components/Notifications/NotificationBell";
import { LogoutDialog } from "@/components/LogoutDialog";
import { isAdminOrModerator, isBusinessRole } from "@/utils/auth";

export function Sidebar() {
  const { theme, setTheme } = useTheme();
  const { data: currentCustomer } = useCurrentCustomer();
  const [logoutOpen, setLogoutOpen] = useState(false);
  const token = localStorage.getItem("token");
  const isAuthenticated = !!token;

  const linkClass = ({ isActive }: { isActive: boolean }) =>
    `px-3 py-2 rounded-lg text-sm transition-colors duration-200 flex items-center gap-2 ${
      isActive
        ? "bg-[#2f5e50] text-white"
        : "text-gray-300 hover:bg-[#2f5e50]/50 hover:text-white"
    }`;

  const sectionLabel = "text-gray-400 px-3 py-1 text-xs font-semibold uppercase tracking-wider";

  const postHref = !isAuthenticated
    ? "/auth"
    : isBusinessRole()
      ? `/venue/${localStorage.getItem("business_org_id") || "1"}/create-box`
      : "/post/create";

  const profileHref = currentCustomer?.userName
    ? `/user/${currentCustomer.userName}`
    : "/profile";

  return (
    <>
      <aside className="flex h-screen w-64 shrink-0 flex-col border-r border-[#2f5e50] bg-[#163d32] p-5">
        <div className="flex min-h-0 flex-1 flex-col gap-4 overflow-y-auto">
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-full bg-[#0b0f0e] text-sm font-bold text-[#98FF98]">
              SB
            </div>
            <div>
              <h2 className="font-semibold text-white">Share Bite</h2>
              <p className="text-xs text-gray-400">The Art of Dining</p>
            </div>
          </div>

          {isAuthenticated && (
            <div className="flex items-center gap-2 rounded-xl border border-[#2f5e50]/60 bg-background/40 p-2">
              <Link
                to={profileHref}
                className="flex min-w-0 flex-1 items-center gap-2 hover:opacity-90"
              >
                <img
                  src={
                    currentCustomer?.avatarURL ||
                    "https://via.placeholder.com/40"
                  }
                  alt=""
                  className="h-9 w-9 shrink-0 rounded-full border border-[#2f5e50] object-cover"
                />
                <span className="truncate text-sm font-semibold text-foreground">
                  @{currentCustomer?.userName || "user"}
                </span>
              </Link>
              <NotificationBell variant="compact" />
            </div>
          )}

          <Button
            variant="ghost"
            className="w-full justify-start px-3 text-gray-300 hover:bg-[#2f5e50]/50 hover:text-white"
            onClick={() => setTheme(theme === "dark" ? "light" : "dark")}
          >
            {theme === "dark" ? (
              <Sun className="mr-2 h-4 w-4" />
            ) : (
              <Moon className="mr-2 h-4 w-4" />
            )}
            {theme === "dark" ? "Light Mode" : "Dark Mode"}
          </Button>

          <Button
            asChild
            className="w-full rounded-full bg-[#FFD700] font-bold text-[#1A3C34] hover:bg-[#FFD700]/80"
          >
            <Link to={postHref}>+ Share a Bite</Link>
          </Button>

          <nav className="flex flex-col gap-1">
            <NavLink to="/" end className={linkClass}>
              Home Feed
            </NavLink>
            <NavLink to="/boxes" className={linkClass}>
              Magic Boxes
            </NavLink>
            <NavLink to="/discover" className={linkClass}>
              Discover
            </NavLink>
            <NavLink to="/venues/search" className={linkClass}>
              Venue Search
            </NavLink>

            <span className={`${sectionLabel} mt-3`}>Social Bites</span>
            <NavLink to="/explore" className={linkClass}>
              Explore
            </NavLink>
            <NavLink to="/collections" className={linkClass}>
              Collections
            </NavLink>

            <span className={`${sectionLabel} mt-3`}>Settings</span>
            <NavLink to="/profile" className={linkClass}>
              Profile
            </NavLink>
            <NavLink to="/profile/edit" className={linkClass}>
              Edit Profile
            </NavLink>
            <NavLink to="/settings/security" className={linkClass}>
              Security
            </NavLink>

            {isAdminOrModerator() && (
              <>
                <span className={`${sectionLabel} mt-3`}>Admin</span>
                <NavLink to="/admin" className={linkClass}>
                  Users
                </NavLink>
              </>
            )}
          </nav>
        </div>

        <div className="mt-4 flex flex-col gap-3 border-t border-[#2f5e50]/60 pt-4">
          {isAuthenticated ? (
            <button
              type="button"
              onClick={() => setLogoutOpen(true)}
              className="flex items-center gap-2 px-3 py-2 text-sm text-gray-300 transition-colors hover:text-white"
            >
              <LogOut className="h-4 w-4" />
              Log out
            </button>
          ) : (
            <NavLink to="/auth" className={linkClass}>
              Log in
            </NavLink>
          )}
          <div className="flex gap-4 px-3 text-xs text-gray-400">
            <span className="cursor-pointer hover:text-white">Support</span>
            <span className="cursor-pointer hover:text-white">Privacy</span>
          </div>
        </div>
      </aside>

      <LogoutDialog open={logoutOpen} onOpenChange={setLogoutOpen} />
    </>
  );
}
