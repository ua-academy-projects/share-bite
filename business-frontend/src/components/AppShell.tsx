import { useLocation } from "react-router-dom";
import { Sidebar } from "@/components/ui/Sidebar";
import { RequireOnboarding } from "@/components/RequireOnboarding/RequireOnboarding";

export function AppShell({ children }: { children: React.ReactNode }) {
  const location = useLocation();
  const isMinimalChrome =
    location.pathname === "/auth" ||
    location.pathname.startsWith("/oauth/");

  if (isMinimalChrome) {
    return (
      <div className="min-h-screen bg-background text-foreground">{children}</div>
    );
  }

  return (
    <div className="flex min-h-screen bg-background text-foreground">
      <Sidebar />
      <main className="flex-1 overflow-auto">
        <RequireOnboarding>{children}</RequireOnboarding>
      </main>
    </div>
  );
}
