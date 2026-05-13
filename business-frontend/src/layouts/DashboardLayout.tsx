import { Suspense } from "react";
import { Outlet } from "react-router-dom";
import { Loader2 } from "lucide-react";

import { Sidebar } from "@/components/ui/Sidebar";

export function DashboardLayout() {
  return (
    <div className="flex min-h-screen bg-background text-foreground">
      <Sidebar />
      <main className="flex-1 min-w-0">
        <Suspense fallback={<RouteLoader />}>
          <Outlet />
        </Suspense>
      </main>
    </div>
  );
}

function RouteLoader() {
  return (
    <div className="min-h-screen w-full flex items-center justify-center bg-background">
      <Loader2 className="h-10 w-10 animate-spin text-emerald-400" />
    </div>
  );
}
