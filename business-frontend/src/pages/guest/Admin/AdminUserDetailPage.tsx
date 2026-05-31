import { useEffect, useState } from "react";
import { Link, useParams } from "react-router-dom";
import { Loader2 } from "lucide-react";
import { apiClient } from "@/api/client";
import { getTokenRole } from "@/utils/auth";
import type { FullUserDetails } from "@/types/api";
import { PageHeader } from "@/components/layout/PageHeader";
import { PageLayout } from "@/components/layout/PageLayout";
import {
  pageBtnPrimary,
  pageLinkAccent,
  pageLoader,
  pagePanel,
} from "@/components/layout/pageStyles";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { cn } from "@/lib/utils";

const STATUS_OPTIONS = ["active", "muted", "suspended"] as const;
const ROLE_OPTIONS = ["user", "business", "moderator", "admin"] as const;

export function AdminUserDetailPage() {
  const { id } = useParams<{ id: string }>();
  const [user, setUser] = useState<FullUserDetails | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [statusSaving, setStatusSaving] = useState(false);
  const [statusMsg, setStatusMsg] = useState("");
  const [statusError, setStatusError] = useState("");
  const [roleSaving, setRoleSaving] = useState(false);
  const [roleMsg, setRoleMsg] = useState("");
  const [roleError, setRoleError] = useState("");
  const [selectedRole, setSelectedRole] = useState("");

  const isAdmin = getTokenRole() === "admin";

  useEffect(() => {
    if (!id) return;
    const load = async () => {
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
    };
    load();
  }, [id]);

  const handleStatusChange = async (newStatus: string) => {
    if (!user || !id) return;
    setStatusSaving(true);
    setStatusMsg("");
    setStatusError("");
    try {
      await apiClient.updateUserStatus(id, newStatus);
      setUser((prev) => (prev ? { ...prev, status: newStatus } : prev));
      setStatusMsg(`Status updated to "${newStatus}".`);
    } catch (err: unknown) {
      const e = err as { response?: { data?: { error?: string } }; message?: string };
      setStatusError(e?.response?.data?.error || e?.message || "Failed to update status.");
    } finally {
      setStatusSaving(false);
    }
  };

  const handleRoleChange = async () => {
    if (!user || !id || selectedRole === user.role_slug) return;
    setRoleSaving(true);
    setRoleMsg("");
    setRoleError("");
    try {
      await apiClient.adminChangeUserRole(id, selectedRole);
      setUser((prev) => (prev ? { ...prev, role_slug: selectedRole } : prev));
      setRoleMsg(`Role updated to "${selectedRole}".`);
    } catch (err: unknown) {
      const e = err as { response?: { data?: { error?: string } }; message?: string };
      setRoleError(e?.response?.data?.error || e?.message || "Failed to update role.");
    } finally {
      setRoleSaving(false);
    }
  };

  if (loading) {
    return (
      <PageLayout maxWidth="3xl">
        <div className="flex h-64 items-center justify-center">
          <Loader2 className={cn(pageLoader, "h-12 w-12")} />
        </div>
      </PageLayout>
    );
  }

  if (error || !user) {
    return (
      <PageLayout maxWidth="3xl">
        <Link to="/admin" className={cn(pageLinkAccent, "mb-4 inline-block text-sm")}>
          ← Back to users
        </Link>
        <p className="text-destructive">{error || "User not found."}</p>
      </PageLayout>
    );
  }

  return (
    <PageLayout maxWidth="3xl">
      <Link to="/admin" className={cn(pageLinkAccent, "mb-4 inline-block text-sm")}>
        ← Back to users
      </Link>

      <PageHeader title={user.email} />

      <div className="grid gap-6">
        <div className={cn(pagePanel, "p-6")}>
          <h2 className="mb-4 text-lg font-semibold text-[#1A3C34] dark:text-white">Account</h2>
          <div className="space-y-3">
            <div className="flex gap-2">
              <Badge>{user.role_slug}</Badge>
              <Badge variant="outline">{user.status}</Badge>
            </div>
            <p className="text-sm text-gray-500 dark:text-gray-400">
              ID: <code>{user.id}</code>
            </p>
            <p className="text-sm text-gray-500 dark:text-gray-400">
              Joined: {new Date(user.created_at).toLocaleString()}
            </p>
          </div>
        </div>

        {user.customer_profile ? (
          <div className={cn(pagePanel, "p-6")}>
            <h2 className="mb-4 text-lg font-semibold text-[#1A3C34] dark:text-white">
              Customer Profile
            </h2>
            <div className="grid gap-3 text-sm">
              <p>
                <span className="text-gray-500 dark:text-gray-400">Username: </span>
                {user.customer_profile.username}
              </p>
              <p>
                <span className="text-gray-500 dark:text-gray-400">Name: </span>
                {user.customer_profile.first_name} {user.customer_profile.last_name}
              </p>
              {user.customer_profile.bio ? (
                <p>
                  <span className="text-gray-500 dark:text-gray-400">Bio: </span>
                  {user.customer_profile.bio}
                </p>
              ) : null}
            </div>
          </div>
        ) : null}

        {user.business_profile ? (
          <div className={cn(pagePanel, "p-6")}>
            <h2 className="mb-4 text-lg font-semibold text-[#1A3C34] dark:text-white">
              Business Profile
            </h2>
            <div className="grid gap-3 text-sm">
              <p>
                <span className="text-gray-500 dark:text-gray-400">Name: </span>
                {user.business_profile.name}
              </p>
              <p>
                <span className="text-gray-500 dark:text-gray-400">Type: </span>
                {user.business_profile.profile_type}
              </p>
              {user.business_profile.description ? (
                <p>
                  <span className="text-gray-500 dark:text-gray-400">Description: </span>
                  {user.business_profile.description}
                </p>
              ) : null}
            </div>
          </div>
        ) : null}

        <div className={cn(pagePanel, "p-6")}>
          <h2 className="mb-4 text-lg font-semibold text-[#1A3C34] dark:text-white">
            Account Status
          </h2>
          <div className="space-y-3">
            <div className="flex flex-wrap gap-2">
              {STATUS_OPTIONS.map((opt) => (
                <Button
                  key={opt}
                  variant={user.status === opt ? "default" : "outline"}
                  size="sm"
                  className={cn(
                    "rounded-xl capitalize",
                    user.status === opt && pageBtnPrimary
                  )}
                  disabled={statusSaving || user.status === opt}
                  onClick={() => handleStatusChange(opt)}
                >
                  {opt}
                </Button>
              ))}
            </div>
            {statusMsg ? (
              <p className="text-sm text-emerald-600 dark:text-[#98FF98]">{statusMsg}</p>
            ) : null}
            {statusError ? <p className="text-sm text-destructive">{statusError}</p> : null}
          </div>
        </div>

        {isAdmin ? (
          <div className={cn(pagePanel, "p-6")}>
            <h2 className="mb-4 text-lg font-semibold text-[#1A3C34] dark:text-white">
              Change Role
            </h2>
            <div className="space-y-3">
              <div className="flex flex-wrap gap-2">
                <Select value={selectedRole} onValueChange={setSelectedRole}>
                  <SelectTrigger className="w-48 rounded-xl border-gray-200 bg-gray-50 dark:border-transparent dark:bg-[#0d241d]">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    {ROLE_OPTIONS.map((r) => (
                      <SelectItem key={r} value={r}>
                        {r}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                <Button
                  className={pageBtnPrimary}
                  onClick={handleRoleChange}
                  disabled={roleSaving || selectedRole === user.role_slug}
                >
                  {roleSaving ? "Saving…" : "Save Role"}
                </Button>
              </div>
              {roleMsg ? (
                <p className="text-sm text-emerald-600 dark:text-[#98FF98]">{roleMsg}</p>
              ) : null}
              {roleError ? <p className="text-sm text-destructive">{roleError}</p> : null}
            </div>
          </div>
        ) : null}
      </div>
    </PageLayout>
  );
}
