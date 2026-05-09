import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { useAuth } from "@/hooks/useAuth";
import { authApi, type UserListItem } from "@/api/auth";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Loader2, Users, ChevronRight } from "lucide-react";

export function UsersPage() {
  const { accessToken } = useAuth();
  const [users, setUsers] = useState<UserListItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    if (!accessToken) return;
    const load = async () => {
      try {
        const data = await authApi.listUsers(accessToken);
        setUsers(Array.isArray(data) ? data : []);
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to load users");
      } finally {
        setLoading(false);
      }
    };
    load();
  }, [accessToken]);

  const statusColor = (status: string) => {
    switch (status) {
      case "active": return "bg-emerald-100 text-emerald-700 dark:bg-emerald-900/30 dark:text-[#98FF98]";
      case "muted": return "bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400";
      case "suspended": return "bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400";
      default: return "bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-400";
    }
  };

  return (
    <div className="min-h-screen bg-[#F9F7F2] dark:bg-[#0d241d] p-8 md:p-12 w-full transition-colors duration-300">
      <div className="max-w-5xl mx-auto">
        <div className="mb-10">
          <h1 className="text-4xl md:text-5xl font-bold text-[#1A3C34] dark:text-white tracking-tight mb-3">
            Users <span className="text-emerald-500 dark:text-[#98FF98]">
              <Users className="inline w-10 h-10" />
            </span>
          </h1>
          <p className="text-gray-600 dark:text-gray-400 text-lg">Manage user accounts and statuses.</p>
        </div>

        {loading ? (
          <div className="flex justify-center items-center h-64">
            <Loader2 className="w-12 h-12 text-emerald-500 dark:text-[#98FF98] animate-spin" />
          </div>
        ) : error ? (
          <Card className="bg-white dark:bg-[#163d32] border-gray-200 dark:border-[#2f5e50]">
            <CardContent className="py-8 text-center">
              <p className="text-destructive">{error}</p>
            </CardContent>
          </Card>
        ) : users.length === 0 ? (
          <div className="text-center bg-white dark:bg-[#163d32] border border-gray-200 dark:border-[#2f5e50] rounded-3xl p-16 shadow-sm dark:shadow-none">
            <p className="text-[#1A3C34] dark:text-gray-300 text-xl font-bold">No users found</p>
          </div>
        ) : (
          <Card className="bg-white dark:bg-[#163d32] border-gray-200 dark:border-[#2f5e50] overflow-hidden">
            <CardHeader>
              <CardTitle className="text-[#1A3C34] dark:text-white">All Users</CardTitle>
              <CardDescription>{users.length} total</CardDescription>
            </CardHeader>
            <CardContent className="p-0">
              <div className="divide-y divide-gray-100 dark:divide-[#2f5e50]">
                {users.map((user) => (
                  <Link
                    key={user.id}
                    to={`/users/${user.id}/status`}
                    className="flex items-center justify-between px-4 py-3 hover:bg-gray-50 dark:hover:bg-[#1a4a3d] transition-colors"
                  >
                    <div className="flex-1 min-w-0">
                      <p className="text-sm font-medium text-[#1A3C34] dark:text-white truncate">
                        {user.email}
                      </p>
                      <p className="text-xs text-muted-foreground">ID: {user.id}</p>
                    </div>
                    <div className="flex items-center gap-3">
                      <span className={`px-2 py-1 rounded-full text-xs font-medium ${statusColor(user.status)}`}>
                        {user.status}
                      </span>
                      <ChevronRight className="w-4 h-4 text-muted-foreground" />
                    </div>
                  </Link>
                ))}
              </div>
            </CardContent>
          </Card>
        )}
      </div>
    </div>
  );
}
