import { NavLink } from "react-router-dom";
import { isAdminOrModerator } from "@/utils/auth";

const sectionLabel =
  "text-[10px] font-black uppercase tracking-widest text-muted-foreground/60 px-3 py-1";

export function Sidebar() {
  const linkClass = ({ isActive }: { isActive: boolean }) =>
    `px-3 py-2 rounded-lg transition-colors duration-200 flex items-center gap-2 ${
      isActive
        ? "bg-secondary text-secondary-foreground"
        : "text-muted-foreground hover:bg-secondary/50 hover:text-foreground"
    }`;

  const hasToken = !!localStorage.getItem("token");

  return (
    <aside className="flex w-64 shrink-0 flex-col justify-between border-r border-border bg-card p-6">
      <nav className="flex flex-col gap-2">
        <p className={sectionLabel}>Home</p>
        <NavLink to="/" end className={linkClass}>
          Home
        </NavLink>

        <p className={`${sectionLabel} mt-4`}>Community</p>
        <NavLink to="/explore" className={linkClass}>
          Explore
        </NavLink>
        {hasToken && (
          <>
            <NavLink to="/collections" className={linkClass}>
              Collections
            </NavLink>
            <NavLink to="/notifications" className={linkClass}>
              Notifications
            </NavLink>
            <NavLink to="/profile" className={linkClass}>
              Profile
            </NavLink>
          </>
        )}

        <p className={`${sectionLabel} mt-4`}>Business</p>
        <NavLink to="/boxes" className={linkClass}>
          Magic Boxes
        </NavLink>
        <NavLink to="/discover" className={linkClass}>
          Discover
        </NavLink>
        <NavLink to="/venues/search" className={linkClass}>
          Venue Search
        </NavLink>

        <p className={`${sectionLabel} mt-4`}>Settings</p>
        {hasToken ? (
          <>
            <NavLink to="/profile/edit" className={linkClass}>
              Edit profile
            </NavLink>
            <NavLink to="/settings/security" className={linkClass}>
              Security
            </NavLink>
          </>
        ) : (
          <NavLink to="/auth" className={linkClass}>
            Sign in
          </NavLink>
        )}

        {isAdminOrModerator() && (
          <>
            <p className={`${sectionLabel} mt-4`}>Admin</p>
            <NavLink to="/admin" className={linkClass}>
              Users
            </NavLink>
          </>
        )}
      </nav>
    </aside>
  );
}
