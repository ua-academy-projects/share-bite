import { useAuth } from "@/hooks/useAuth";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Shield, Users, Activity } from "lucide-react";

export function DashboardPage() {
  const { userPayload } = useAuth();

  return (
    <div className="min-h-screen bg-[#F9F7F2] dark:bg-[#0d241d] p-8 md:p-12 w-full transition-colors duration-300">
      <div className="max-w-5xl mx-auto">
        <div className="mb-10">
          <h1 className="text-4xl md:text-5xl font-bold text-[#1A3C34] dark:text-white tracking-tight mb-3">
            Dashboard <span className="text-emerald-500 dark:text-[#98FF98]">
              <Shield className="inline w-10 h-10" />
            </span>
          </h1>
          <p className="text-gray-600 dark:text-gray-400 text-lg">Welcome to ShareBite</p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          <Card className="bg-white dark:bg-[#163d32] border-gray-200 dark:border-[#2f5e50]">
            <CardHeader>
              <CardTitle className="flex items-center gap-2 text-[#1A3C34] dark:text-white">
                <Users className="w-5 h-5 text-emerald-500 dark:text-[#98FF98]" />
                Your Role
              </CardTitle>
              <CardDescription>Current account role</CardDescription>
            </CardHeader>
            <CardContent>
              <p className="text-2xl font-bold text-[#1A3C34] dark:text-[#98FF98] capitalize">
                {userPayload?.role || "unknown"}
              </p>
            </CardContent>
          </Card>

          <Card className="bg-white dark:bg-[#163d32] border-gray-200 dark:border-[#2f5e50]">
            <CardHeader>
              <CardTitle className="flex items-center gap-2 text-[#1A3C34] dark:text-white">
                <Activity className="w-5 h-5 text-emerald-500 dark:text-[#98FF98]" />
                Account Status
              </CardTitle>
              <CardDescription>Current account status</CardDescription>
            </CardHeader>
            <CardContent>
              <p className={`text-2xl font-bold capitalize ${
                userPayload?.status === "active"
                  ? "text-emerald-600 dark:text-[#98FF98]"
                  : userPayload?.status === "muted"
                    ? "text-yellow-600 dark:text-yellow-400"
                    : "text-red-600 dark:text-red-400"
              }`}>
                {userPayload?.status || "unknown"}
              </p>
            </CardContent>
          </Card>

          <Card className="bg-white dark:bg-[#163d32] border-gray-200 dark:border-[#2f5e50]">
            <CardHeader>
              <CardTitle className="flex items-center gap-2 text-[#1A3C34] dark:text-white">
                <Shield className="w-5 h-5 text-emerald-500 dark:text-[#98FF98]" />
                Quick Actions
              </CardTitle>
              <CardDescription>Admin tools</CardDescription>
            </CardHeader>
            <CardContent>
              <a
                href="/users"
                className="text-primary hover:underline text-sm"
              >
                Manage Users →
              </a>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
