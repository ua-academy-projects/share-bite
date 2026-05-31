import { useState } from "react";
import { NavLink, useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { useTheme } from "@/components/theme-provider";
import { NotificationBell } from "@/components/Notifications/NotificationBell";
import { LogoutDialog } from "@/components/LogoutDialog";
import { useCurrentCustomer } from "@/hooks/useCurrentCustomer";
import { useOnboardingStatus } from "@/hooks/useOnboardingStatus";
import {
  getBusinessOrgId,
  getTokenRole,
  isAdminOrModerator,
  isBusinessRole,
  isUserRole,
} from "@/utils/auth";
import { Moon, Sun, LogOut } from "lucide-react";
import { cn } from "@/lib/utils";

function NavSection({ label }: { label: string }) {
  return (
    <span className="mt-4 px-3 py-1 text-xs font-semibold uppercase tracking-wider text-gray-500">
      {label}
    </span>
  );
}

export function Sidebar() {
  const { theme, setTheme } = useTheme();
  const navigate = useNavigate();
  const token = localStorage.getItem("token");
  const { data: customer } = useCurrentCustomer();
  useOnboardingStatus(!!token);
  const [logoutOpen, setLogoutOpen] = useState(false);
  const businessOrgId = getBusinessOrgId();

  const linkClass = ({ isActive }: { isActive: boolean }) =>
    cn(
      "flex items-center gap-2 rounded-lg px-3 py-2 transition-colors duration-200",
      isActive
        ? "bg-[#2f5e50] text-white"
        : "text-gray-300 hover:bg-[#2f5e50]/50 hover:text-white"
    );

  const primaryCta = () => {
    if (!token) return { label: "Sign in to post", path: "/auth" };
    if (isBusinessRole()) {
      if (businessOrgId) {
        return {
          label: "List a rescue box",
          path: `/venue/${businessOrgId}/create-box`,
        };
      }
      return { label: "Set up your venue", path: "/business/setup" };
    }
    return { label: "Write a review", path: "/post/create" };
  };

  const cta = primaryCta();

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
                  {customer?.userName
                    ? `@${customer.userName}`
                    : isBusinessRole()
                      ? "Business"
                      : "Guest"}
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
            onClick={() => navigate(cta.path)}
          >
            + {cta.label}
          </Button>

          <nav className="flex flex-col gap-2">
            <NavSection label="Feeds" />
            <NavLink to="/feed/users" className={linkClass}>
              Users Feed
            </NavLink>
            <NavLink to="/feed/business" className={linkClass}>
              Business Feed
            </NavLink>

            <NavSection label="Discover" />
            <NavLink to="/discover" className={linkClass}>
              Find Venues
            </NavLink>
            <NavLink to="/boxes" className={linkClass}>
              Rescue Boxes
            </NavLink>

            {token && isUserRole() ? (
              <>
                <NavSection label="For you" />
                <NavLink to="/collections" className={linkClass}>
                  Collections
                </NavLink>
                <NavLink to="/profile" className={linkClass}>
                  Profile
                </NavLink>
                <NavLink to="/profile/edit" className={linkClass}>
                  Edit Profile
                </NavLink>
              </>
            ) : null}

            {token && isBusinessRole() ? (
              <>
                <NavSection label="Your venue" />
                {businessOrgId ? (
                  <NavLink to={`/venue/${businessOrgId}`} className={linkClass}>
                    Venue profile
                  </NavLink>
                ) : (
                  <NavLink to="/business/setup" className={linkClass}>
                    Set up venue
                  </NavLink>
                )}
              </>
            ) : null}

            {token ? (
              <>
                <NavSection label="Settings" />
                <NavLink to="/notifications" className={linkClass}>
                  Notifications
                </NavLink>
                <NavLink to="/settings/security" className={linkClass}>
                  Security
                </NavLink>
              </>
            ) : null}

            {isAdminOrModerator() ? (
              <>
                <NavSection label="Admin" />
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
