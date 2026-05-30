import { useLocation } from "react-router-dom";
import { Navbar } from "@/components/Navbar/Navbar";
import { Sidebar } from "@/components/ui/Sidebar";

const MINIMAL_CHROME_PREFIXES = ["/auth", "/oauth/"];

export function AppShell({ children }: { children: React.ReactNode }) {
  const { pathname } = useLocation();
  const minimalChrome = MINIMAL_CHROME_PREFIXES.some((p) =>
    pathname.startsWith(p)
  );

  if (minimalChrome) {
    return (
      <div className="flex min-h-screen flex-col bg-background text-foreground">
        <Navbar />
        <main className="flex flex-1 flex-col">{children}</main>
      </div>
    );
  }

  return (
    <div className="flex min-h-screen flex-col bg-background text-foreground">
      <Navbar />
      <div className="flex min-h-0 flex-1">
        <Sidebar />
        <main className="flex-1 overflow-auto">{children}</main>
      </div>
    </div>
  );
}
