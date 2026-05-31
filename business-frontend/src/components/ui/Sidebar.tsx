import { useState } from "react";
import { NavLink, useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { useTheme } from "@/components/theme-provider";
import { NotificationBell } from "@/components/Notifications/NotificationBell";
import { LogoutDialog } from "@/components/LogoutDialog";
import { useCurrentCustomer } from "@/hooks/useCurrentCustomer";
import {
  getTokenRole,
  isAdminOrModerator,
  isBusinessRole,
} from "@/utils/auth";
import { Moon, Sun, LogOut } from "lucide-react";
import { cn } from "@/lib/utils";

export function Sidebar() {
  const { theme, setTheme } = useTheme();
  const navigate = useNavigate();
  const token = localStorage.getItem("token");
  const { data: customer } = useCurrentCustomer();
  const [logoutOpen, setLogoutOpen] = useState(false);

  const linkClass = ({ isActive }: { isActive: boolean }) =>
    cn(
      "flex items-center gap-2 rounded-lg px-3 py-2 transition-colors duration-200",
      isActive
        ? "bg-[#2f5e50] text-white"
        : "text-gray-300 hover:bg-[#2f5e50]/50 hover:text-white"
    );

  const shareCtaPath = () => {
    if (!token) return "/auth";
    if (isBusinessRole()) {
      const orgId = localStorage.getItem("business_org_id");
      return orgId ? `/venue/${orgId}/create-box` : "/discover";
    }
    return "/post/create";
  };

  const avatarInitial =
    customer?.userName?.charAt(0)?.toUpperCase() ||
    getTokenRole()?.charAt(0)?.toUpperCase() ||
    "?";

  return (
    <>
      <aside className="flex w-64 flex-col justify-between border-r border-[#2f5e50] bg-[#163d32] p-6">
        <div>
          <div className="mb-4 flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-full bg-[#0b0f0e] font-bold text-[#98FF98]">
              SB
            </div>
            <div>
              <h2 className="font-semibold text-white">Share Bite</h2>
              <p className="text-xs text-gray-400">The Art of Dining</p>
            </div>
          </div>

          {token ? (
            <div className="mb-4 flex items-center gap-2 rounded-2xl border border-[#2f5e50] bg-[#0b0f0e]/40 p-2">
              <div className="flex h-9 w-9 shrink-0 items-center justify-center overflow-hidden rounded-full bg-[#2f5e50] text-sm font-semibold text-white">
                {customer?.avatarURL ? (
                  <img
                    src={customer.avatarURL}
                    alt=""
                    className="h-full w-full object-cover"
                  />
                ) : (
                  avatarInitial
                )}
              </div>
              <div className="min-w-0 flex-1">
                <p className="truncate text-sm font-medium text-white">
                  {customer?.userName ? `@${customer.userName}` : "Guest"}
                </p>
              </div>
              <NotificationBell variant="compact" />
            </div>
          ) : null}

          <div className="mb-4">
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
          </div>

          <Button
            className="mb-8 w-full rounded-full bg-[#FFD700] font-bold text-[#1A3C34] hover:bg-[#FFD700]/80"
            onClick={() => navigate(shareCtaPath())}
          >
            + Share a Bite
          </Button>

          <nav className="flex flex-col gap-2">
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

            <span className="mt-4 px-3 py-1 text-xs font-semibold uppercase tracking-wider text-gray-500">
              Social
            </span>
            <NavLink to="/explore" className={linkClass}>
              Explore
            </NavLink>
            <NavLink to="/collections" className={linkClass}>
              Collections
            </NavLink>

            {token ? (
              <>
                <span className="mt-4 px-3 py-1 text-xs font-semibold uppercase tracking-wider text-gray-500">
                  Settings
                </span>
                <NavLink to="/profile" className={linkClass}>
                  Profile
                </NavLink>
                <NavLink to="/profile/edit" className={linkClass}>
                  Edit Profile
                </NavLink>
                <NavLink to="/settings/security" className={linkClass}>
                  Security
                </NavLink>
              </>
            ) : null}

            {isAdminOrModerator() ? (
              <>
                <span className="mt-4 px-3 py-1 text-xs font-semibold uppercase tracking-wider text-gray-500">
                  Admin
                </span>
                <NavLink to="/admin" className={linkClass}>
                  Admin Users
                </NavLink>
              </>
            ) : null}
          </nav>
        </div>

        <div className="flex flex-col gap-4">
          {token ? (
            <Button
              variant="ghost"
              className="justify-start px-3 text-gray-300 hover:bg-[#2f5e50]/50 hover:text-white"
              onClick={() => setLogoutOpen(true)}
            >
              <LogOut className="mr-2 h-4 w-4" />
              Log out
            </Button>
          ) : (
            <NavLink to="/auth" className={linkClass}>
              Sign in
            </NavLink>
          )}
          <div className="flex gap-4 px-3 text-xs text-gray-400">
            <span className="cursor-pointer transition-colors hover:text-white">
              Support
            </span>
            <span className="cursor-pointer transition-colors hover:text-white">
              Privacy
            </span>
          </div>
        </div>
      </aside>

      <LogoutDialog open={logoutOpen} onOpenChange={setLogoutOpen} />
    </>
  );
}
