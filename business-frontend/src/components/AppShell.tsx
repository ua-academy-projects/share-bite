import { Link, useLocation } from "react-router-dom";
import { Sidebar } from "@/components/ui/Sidebar";

function MinimalHeader() {
  return (
    <header className="border-b border-border bg-[#163d32] px-6 py-4">
      <Link to="/" className="flex items-center gap-3">
        <div className="flex h-9 w-9 items-center justify-center rounded-full bg-[#0b0f0e] text-sm font-bold text-[#98FF98]">
          SB
        </div>
        <div>
          <span className="font-semibold text-white">Share Bite</span>
          <p className="text-xs text-gray-400">The Art of Dining</p>
        </div>
      </Link>
    </header>
  );
}

export function AppShell({ children }: { children: React.ReactNode }) {
  const location = useLocation();
  const isMinimalChrome =
    location.pathname === "/auth" ||
    location.pathname.startsWith("/oauth/");

  if (isMinimalChrome) {
    return (
      <div className="flex min-h-screen flex-col bg-background text-foreground">
        <MinimalHeader />
        <main className="flex-1">{children}</main>
      </div>
    );
  }

  return (
    <div className="flex min-h-screen bg-background text-foreground">
      <Sidebar />
      <main className="flex-1 overflow-auto">{children}</main>
    </div>
  );
}
