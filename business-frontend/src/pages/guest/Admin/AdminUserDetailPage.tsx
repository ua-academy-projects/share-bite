import { useEffect, useState } from "react";
import { Link, useParams } from "react-router-dom";
import { ArrowLeft, Loader2 } from "lucide-react";
import { apiClient } from "@/api/client";
import { getTokenRole } from "@/utils/auth";
import type { FullUserDetails } from "@/types/api";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

const STATUSES = ["active", "muted", "suspended"] as const;
const ROLES = ["user", "business", "moderator", "admin"] as const;

export function AdminUserDetailPage() {
  const { id } = useParams<{ id: string }>();
  const [user, setUser] = useState<FullUserDetails | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [selectedRole, setSelectedRole] = useState("");
  const [statusMsg, setStatusMsg] = useState("");
  const [roleMsg, setRoleMsg] = useState("");
  const isAdmin = getTokenRole() === "admin";

  useEffect(() => {
    if (!id) return;
    void (async () => {
      try {
        const data = await apiClient.adminGetUserDetails(id);
        setUser(data);
        setSelectedRole(data.role_slug);
      } catch (err: unknown) {
        const e = err as { response?: { data?: { error?: string } }; message?: string };
        setError(e?.response?.data?.error || e?.message || "Failed to load user.");
      } finally {
        setLoading(false);
      }
    })();
  }, [id]);

  if (loading) {
    return (
      <div className="flex justify-center py-24">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (error || !user) {
    return (
      <div className="px-6 py-8">
        <p className="text-destructive">{error || "User not found"}</p>
        <Button asChild variant="outline" className="mt-4 rounded-xl">
          <Link to="/admin">Back to users</Link>
        </Button>
      </div>
    );
  }

  return (
    <div className="px-6 py-8 lg:px-10">
      <Button asChild variant="ghost" className="mb-6 rounded-xl">
        <Link to="/admin">
          <ArrowLeft className="mr-2 h-4 w-4" /> Back to users
        </Link>
      </Button>

      <Card className="mx-auto max-w-2xl rounded-3xl bg-card-solid">
        <CardContent className="space-y-6 p-8">
          <div>
            <h1 className="text-2xl font-bold">{user.email}</h1>
            <p className="text-sm text-muted-foreground">
              Joined {new Date(user.created_at).toLocaleDateString()}
            </p>
            <div className="mt-3 flex gap-2">
              <Badge variant="outline">{user.role_slug}</Badge>
              <Badge variant="secondary">{user.status}</Badge>
            </div>
          </div>

          <div className="space-y-3 border-t border-border pt-6">
            <h2 className="font-semibold">Status</h2>
            <div className="flex flex-wrap gap-2">
              {STATUSES.map((status) => (
                <Button
                  key={status}
                  variant={user.status === status ? "default" : "outline"}
                  size="sm"
                  className="rounded-xl capitalize"
                  onClick={async () => {
                    if (!id) return;
                    try {
                      await apiClient.updateUserStatus(id, status);
                      setUser({ ...user, status });
                      setStatusMsg(`Status updated to ${status}.`);
                    } catch {
                      setStatusMsg("Failed to update status.");
                    }
                  }}
                >
                  {status}
                </Button>
              ))}
            </div>
            {statusMsg && (
              <p className="text-sm text-muted-foreground">{statusMsg}</p>
            )}
          </div>

          {isAdmin && (
            <div className="space-y-3 border-t border-border pt-6">
              <h2 className="font-semibold">Role</h2>
              <div className="flex flex-wrap gap-3">
                <Select value={selectedRole} onValueChange={setSelectedRole}>
                  <SelectTrigger className="w-40 rounded-xl">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    {ROLES.map((role) => (
                      <SelectItem key={role} value={role}>
                        {role}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                <Button
                  className="rounded-xl"
                  disabled={selectedRole === user.role_slug}
                  onClick={async () => {
                    if (!id) return;
                    try {
                      await apiClient.adminChangeUserRole(id, selectedRole);
                      setUser({ ...user, role_slug: selectedRole });
                      setRoleMsg(`Role updated to ${selectedRole}.`);
                    } catch {
                      setRoleMsg("Failed to update role.");
                    }
                  }}
                >
                  Save role
                </Button>
              </div>
              {roleMsg && (
                <p className="text-sm text-muted-foreground">{roleMsg}</p>
              )}
            </div>
          )}

          {user.customer_profile && (
            <div className="border-t border-border pt-6">
              <h2 className="mb-2 font-semibold">Customer profile</h2>
              <p className="text-sm text-muted-foreground">
                @{user.customer_profile.username} —{" "}
                {user.customer_profile.first_name}{" "}
                {user.customer_profile.last_name}
              </p>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
