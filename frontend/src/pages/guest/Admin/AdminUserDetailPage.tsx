import { useEffect, useState } from "react";
import { Link, useParams } from "react-router-dom";
import { Building2, Calendar, Loader2, Mail, Shield, Tag, User } from "lucide-react";
import { apiClient } from "@/api/client";
import { getTokenRole } from "@/utils/auth";
import type { FullUserDetails } from "@/types/api";
import { PageLayout } from "@/components/layout/PageLayout";
import { pageInput, pageLoader } from "@/components/layout/pageStyles";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
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

function rolePillClass(role: string) {
  if (role === "admin") return "border-red-500/40 bg-red-500/10 text-red-300";
  if (role === "moderator") return "border-[#FFD700]/40 bg-[#FFD700]/10 text-[#FFD700]";
  if (role === "business") return "border-emerald-500/40 bg-emerald-500/10 text-emerald-300";
  return "border-[#2f5e50] bg-[#0d241d] text-gray-200";
}

function statusPillClass(status: string) {
  if (status === "active") return "border-emerald-500/40 bg-emerald-500/10 text-emerald-300";
  if (status === "suspended") return "border-red-500/40 bg-red-500/10 text-red-300";
  return "border-[#2f5e50] bg-[#0d241d] text-gray-300";
}

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
    void load();
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
      <PageLayout>
        <div className="flex h-64 items-center justify-center">
          <Loader2 className={cn(pageLoader, "h-12 w-12")} />
        </div>
      </PageLayout>
    );
  }

  if (error || !user) {
    return (
      <PageLayout>
        <Card className="max-w-lg rounded-3xl border border-gray-200 bg-white dark:border-[#2f5e50] dark:bg-[#163d32]">
          <CardContent className="space-y-4 p-8">
            <h1 className="text-2xl font-bold text-[#1A3C34] dark:text-white">User not found</h1>
            <p className="text-red-400">{error || "Unknown error"}</p>
            <Button asChild className="bg-[#163d32] text-white hover:bg-[#1A3C34]">
              <Link to="/admin">Back to users</Link>
            </Button>
          </CardContent>
        </Card>
      </PageLayout>
    );
  }

  return (
    <PageLayout className="max-w-5xl space-y-6">
      <Button
        asChild
        variant="ghost"
        className="px-0 text-emerald-600 hover:text-emerald-700 dark:text-[#98FF98] dark:hover:text-emerald-300"
      >
        <Link to="/admin">← Back to users</Link>
      </Button>

      <Card className="overflow-hidden rounded-3xl border border-gray-200 bg-white dark:border-[#2f5e50] dark:bg-[#163d32]">
        <div className="h-32 bg-[#0d241d] md:h-40" />
        <CardContent className="space-y-5 p-6 md:p-8">
          <div className="-mt-16 flex items-start gap-4 md:-mt-20">
            <div className="flex h-20 w-20 shrink-0 items-center justify-center rounded-2xl border-4 border-[#163d32] bg-[#163d32] text-3xl font-bold text-[#98FF98] dark:border-[#0d241d]">
              {user.email.charAt(0).toUpperCase()}
            </div>
            <div className="min-w-0 flex-1 pt-2 md:pt-6">
              <h1 className="truncate text-2xl font-bold text-[#1A3C34] dark:text-white md:text-3xl">
                {user.email}
              </h1>
              <div className="mt-3 flex flex-wrap gap-2">
                <span
                  className={cn(
                    "inline-flex items-center gap-1 rounded-full border px-2.5 py-1 text-xs capitalize",
                    rolePillClass(user.role_slug)
                  )}
                >
                  <Shield className="h-3 w-3" />
                  {user.role_slug}
                </span>
                <span
                  className={cn(
                    "inline-flex items-center gap-1 rounded-full border px-2.5 py-1 text-xs capitalize",
                    statusPillClass(user.status)
                  )}
                >
                  <User className="h-3 w-3" />
                  {user.status}
                </span>
              </div>
            </div>
          </div>

          <div className="flex flex-wrap gap-4 text-sm text-gray-600 dark:text-gray-300">
            <span className="inline-flex items-center gap-1">
              <Mail className="h-4 w-4" />
              {user.email}
            </span>
            <span className="inline-flex items-center gap-1">
              <Calendar className="h-4 w-4" />
              Joined {new Date(user.created_at).toLocaleString()}
            </span>
            <span className="inline-flex items-center gap-1 font-mono text-xs text-gray-500">
              ID: {user.id}
            </span>
          </div>
        </CardContent>
      </Card>

      {user.customer_profile ? (
        <Card className="rounded-3xl border border-gray-200 bg-white dark:border-[#2f5e50] dark:bg-[#163d32]">
          <CardContent className="space-y-4 p-6">
            <h2 className="text-lg font-semibold text-[#1A3C34] dark:text-white">
              Customer Profile
            </h2>
            <div className="flex flex-wrap gap-2">
              <span className="inline-flex items-center gap-1 rounded-full border border-[#2f5e50] bg-[#0d241d] px-2.5 py-1 text-xs text-gray-200">
                <Tag className="h-3 w-3" />@{user.customer_profile.username}
              </span>
            </div>
            <p className="text-sm text-gray-600 dark:text-gray-300">
              {user.customer_profile.first_name} {user.customer_profile.last_name}
            </p>
            {user.customer_profile.bio ? (
              <p className="text-sm text-gray-500 dark:text-gray-400">
                {user.customer_profile.bio}
              </p>
            ) : null}
          </CardContent>
        </Card>
      ) : null}

      {user.business_profile ? (
        <Card className="rounded-3xl border border-gray-200 bg-white dark:border-[#2f5e50] dark:bg-[#163d32]">
          <CardContent className="space-y-4 p-6">
            <h2 className="text-lg font-semibold text-[#1A3C34] dark:text-white">
              Business Profile
            </h2>
            <div className="flex flex-wrap gap-2">
              <span className="inline-flex items-center gap-1 rounded-full border border-[#2f5e50] bg-[#0d241d] px-2.5 py-1 text-xs text-gray-200">
                <Building2 className="h-3 w-3" />
                {user.business_profile.profile_type}
              </span>
            </div>
            <p className="text-sm font-medium text-gray-600 dark:text-gray-300">
              {user.business_profile.name}
            </p>
            {user.business_profile.description ? (
              <p className="text-sm text-gray-500 dark:text-gray-400">
                {user.business_profile.description}
              </p>
            ) : null}
          </CardContent>
        </Card>
      ) : null}

      <Card className="rounded-3xl border border-gray-200 bg-white dark:border-[#2f5e50] dark:bg-[#163d32]">
        <CardContent className="space-y-4 p-6">
          <h2 className="text-lg font-semibold text-[#1A3C34] dark:text-white">Account Status</h2>
          <div className="flex flex-wrap gap-2">
            {STATUS_OPTIONS.map((opt) => (
              <Button
                key={opt}
                type="button"
                size="sm"
                variant="outline"
                disabled={statusSaving || user.status === opt}
                onClick={() => void handleStatusChange(opt)}
                className={cn(
                  "rounded-full capitalize",
                  user.status === opt &&
                    "border-emerald-500/40 bg-emerald-500/10 text-emerald-300 hover:bg-emerald-500/20"
                )}
              >
                {opt}
              </Button>
            ))}
          </div>
          {statusMsg ? (
            <p className="text-sm text-emerald-600 dark:text-[#98FF98]">{statusMsg}</p>
          ) : null}
          {statusError ? <p className="text-sm text-red-400">{statusError}</p> : null}
        </CardContent>
      </Card>

      {isAdmin ? (
        <Card className="rounded-3xl border border-gray-200 bg-white dark:border-[#2f5e50] dark:bg-[#163d32]">
          <CardContent className="space-y-4 p-6">
            <h2 className="text-lg font-semibold text-[#1A3C34] dark:text-white">Change Role</h2>
            <div className="flex flex-wrap items-center gap-3">
              <Select value={selectedRole} onValueChange={setSelectedRole}>
                <SelectTrigger className={cn(pageInput, "h-auto w-48 cursor-pointer")}>
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
                className="rounded-xl bg-[#FFD700] font-bold text-[#1A3C34] hover:bg-[#e6c200]"
                onClick={() => void handleRoleChange()}
                disabled={roleSaving || selectedRole === user.role_slug}
              >
                {roleSaving ? "Saving…" : "Save role"}
              </Button>
            </div>
            {roleMsg ? (
              <p className="text-sm text-emerald-600 dark:text-[#98FF98]">{roleMsg}</p>
            ) : null}
            {roleError ? <p className="text-sm text-red-400">{roleError}</p> : null}
          </CardContent>
        </Card>
      ) : null}
    </PageLayout>
  );
}
