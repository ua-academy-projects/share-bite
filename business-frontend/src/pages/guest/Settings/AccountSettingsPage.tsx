import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import {
  AlertTriangle,
  Loader2,
  Mail,
  Shield,
  Smartphone,
  Tag,
  User,
} from "lucide-react";
import { toast } from "sonner";
import { apiClient } from "@/api/client";
import { useCurrentCustomer } from "@/hooks/useCurrentCustomer";
import { PageHeader } from "@/components/layout/PageHeader";
import { PageLayout } from "@/components/layout/PageLayout";
import {
  pageBtnPrimary,
  pageInput,
  pageLabel,
  pageLinkAccent,
  pageLoader,
} from "@/components/layout/pageStyles";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import {
  getBusinessOrgId,
  getTokenPayload,
  isBusinessRole,
} from "@/utils/auth";
import { cn } from "@/lib/utils";

const settingsCardClass =
  "rounded-3xl border border-gray-200 bg-white shadow-sm dark:border-[#2f5e50] dark:bg-[#163d32]";

const pageBtnDangerOutline =
  "rounded-xl border border-red-500/50 bg-transparent px-6 py-2.5 font-bold text-red-400 shadow-sm transition-all hover:border-red-500 hover:bg-red-500/10 disabled:opacity-70";

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

export function AccountSettingsPage() {
  const queryClient = useQueryClient();
  const payload = getTokenPayload();
  const { data: customer, isLoading: customerLoading, isError: noCustomer } =
    useCurrentCustomer();
  const businessOrgId = getBusinessOrgId();

  const [form, setForm] = useState({
    userName: "",
    firstName: "",
    lastName: "",
  });
  const [isRevoking, setIsRevoking] = useState(false);

  useEffect(() => {
    if (customer) {
      setForm({
        userName: customer.userName,
        firstName: customer.firstName,
        lastName: customer.lastName,
      });
    }
  }, [customer]);

  const updateMutation = useMutation({
    mutationFn: () =>
      apiClient.updateCustomer({
        userName: form.userName,
        firstName: form.firstName,
        lastName: form.lastName,
      }),
    onSuccess: (updated) => {
      void queryClient.invalidateQueries({ queryKey: ["currentCustomer"] });
      void queryClient.invalidateQueries({ queryKey: ["user", updated.userName] });
      toast.success("Profile updated successfully");
    },
    onError: (error: unknown) => {
      const e = error as { response?: { data?: { error?: string } } };
      toast.error(e?.response?.data?.error || "Failed to update profile");
    },
  });

  const email = payload?.email ?? "—";
  const role = payload?.role ?? "—";
  const status = payload?.status ?? "—";
  const userId = payload?.sub ?? "—";

  const displayName =
    customer?.firstName && customer?.lastName
      ? `${customer.firstName} ${customer.lastName}`
      : payload?.name ?? email;

  const handleRevokeSessions = async () => {
    setIsRevoking(true);
    try {
      await apiClient.revokeAllSessions();
      toast.success(
        "All sessions revoked successfully. You will be logged out of other devices."
      );
    } catch (error: unknown) {
      const e = error as { response?: { data?: { error?: string } } };
      toast.error(e?.response?.data?.error || "Failed to revoke sessions.");
    } finally {
      setIsRevoking(false);
    }
  };

  return (
    <PageLayout className="space-y-8">
      <PageHeader
        title={
          <>
            Account Settings{" "}
            <span className="text-emerald-500 dark:text-[#98FF98]">⚙️</span>
          </>
        }
        description="View your account details and manage profile settings"
      />

      <div className="mx-auto flex w-full max-w-3xl flex-col gap-6">
        <Card className={settingsCardClass}>
          <CardContent className="space-y-6 p-6">
            <div className="flex items-center gap-2 font-semibold text-[#1A3C34] dark:text-white">
              <User size={20} className="text-emerald-500 dark:text-[#98FF98]" />
              <span>Account</span>
            </div>

            <div className="flex items-start gap-4">
              <div className="flex h-14 w-14 shrink-0 items-center justify-center rounded-2xl border border-gray-200 bg-[#0d241d] text-xl font-bold text-[#98FF98] dark:border-[#2f5e50]">
                {customer?.userName?.charAt(0)?.toUpperCase() ||
                  displayName.charAt(0).toUpperCase()}
              </div>
              <div className="min-w-0 flex-1 space-y-3">
                <p className="text-sm font-medium text-[#1A3C34] dark:text-white">
                  {displayName}
                </p>
                <div className="flex flex-wrap gap-2">
                  <span
                    className={cn(
                      "inline-flex items-center gap-1 rounded-full border px-2.5 py-1 text-xs capitalize",
                      rolePillClass(role)
                    )}
                  >
                    <Shield className="h-3 w-3" />
                    {role}
                  </span>
                  <span
                    className={cn(
                      "inline-flex items-center gap-1 rounded-full border px-2.5 py-1 text-xs capitalize",
                      statusPillClass(status)
                    )}
                  >
                    {status}
                  </span>
                </div>
                <div className="flex flex-wrap gap-4 text-sm text-gray-600 dark:text-gray-300">
                  <span className="inline-flex items-center gap-1">
                    <Mail className="h-4 w-4" />
                    {email}
                  </span>
                  <span className="font-mono text-xs text-gray-500">ID: {userId}</span>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card className={settingsCardClass}>
          <CardContent className="space-y-6 p-6">
            <div className="flex items-center gap-2 font-semibold text-[#1A3C34] dark:text-white">
              <Tag size={20} className="text-emerald-500 dark:text-[#98FF98]" />
              <span>Profile information</span>
            </div>

            {customerLoading ? (
              <div className="flex h-32 items-center justify-center">
                <Loader2 className={cn(pageLoader, "h-8 w-8")} />
              </div>
            ) : customer ? (
              <form
                onSubmit={(e) => {
                  e.preventDefault();
                  updateMutation.mutate();
                }}
                className="space-y-5"
              >
                <div className="space-y-2">
                  <label htmlFor="userName" className={pageLabel}>
                    Username
                  </label>
                  <input
                    id="userName"
                    required
                    value={form.userName}
                    onChange={(e) =>
                      setForm((f) => ({ ...f, userName: e.target.value }))
                    }
                    className={pageInput}
                    placeholder="foodie"
                  />
                </div>
                <div className="grid gap-5 sm:grid-cols-2">
                  <div className="space-y-2">
                    <label htmlFor="firstName" className={pageLabel}>
                      First name
                    </label>
                    <input
                      id="firstName"
                      required
                      value={form.firstName}
                      onChange={(e) =>
                        setForm((f) => ({ ...f, firstName: e.target.value }))
                      }
                      className={pageInput}
                    />
                  </div>
                  <div className="space-y-2">
                    <label htmlFor="lastName" className={pageLabel}>
                      Last name
                    </label>
                    <input
                      id="lastName"
                      required
                      value={form.lastName}
                      onChange={(e) =>
                        setForm((f) => ({ ...f, lastName: e.target.value }))
                      }
                      className={pageInput}
                    />
                  </div>
                </div>
                <div className="flex flex-wrap gap-3 border-t border-gray-200 pt-6 dark:border-[#2f5e50]">
                  <Button
                    type="submit"
                    className={pageBtnPrimary}
                    disabled={updateMutation.isPending}
                  >
                    {updateMutation.isPending ? "Saving…" : "Save profile"}
                  </Button>
                  <Button asChild variant="outline" className="rounded-xl">
                    <Link to="/profile">View public profile</Link>
                  </Button>
                </div>
              </form>
            ) : (
              <div className="space-y-4 text-sm text-gray-600 dark:text-gray-300">
                <p>
                  {noCustomer
                    ? "Create a guest profile to set your username and display name."
                    : "Profile information is unavailable right now."}
                </p>
                {noCustomer ? (
                  <Button asChild className={pageBtnPrimary}>
                    <Link to="/profile/create">Create profile</Link>
                  </Button>
                ) : null}
              </div>
            )}
          </CardContent>
        </Card>

        <Card className={settingsCardClass}>
          <CardContent className="space-y-6 p-6">
            <div className="flex items-center gap-2 font-semibold text-[#1A3C34] dark:text-white">
              <Smartphone size={20} className="text-emerald-500 dark:text-[#98FF98]" />
              <span>Active Sessions</span>
            </div>

            <div className="flex items-start gap-4">
              <div className="flex h-14 w-14 shrink-0 items-center justify-center rounded-2xl border border-gray-200 bg-[#0d241d] dark:border-[#2f5e50]">
                <Shield className="h-7 w-7 text-[#98FF98]" />
              </div>
              <div className="min-w-0 flex-1">
                <p className="text-sm leading-relaxed text-gray-600 dark:text-gray-300">
                  Manage the devices that are currently logged into your account.
                  Revoking all sessions will log out all other devices.
                </p>
              </div>
            </div>

            <div className="flex flex-col gap-4 border-t border-gray-200 pt-6 sm:flex-row sm:items-center sm:justify-between dark:border-[#2f5e50]">
              <div className="inline-flex items-center gap-2 text-sm font-medium text-red-400">
                <AlertTriangle className="h-4 w-4 shrink-0" />
                This action cannot be undone.
              </div>
              <Button
                variant="outline"
                onClick={() => void handleRevokeSessions()}
                disabled={isRevoking}
                className={cn(pageBtnDangerOutline, "h-auto sm:shrink-0")}
              >
                {isRevoking ? "Revoking…" : "Revoke All Sessions"}
              </Button>
            </div>
          </CardContent>
        </Card>

        {isBusinessRole() ? (
          <div className="px-1">
            <Link
              to={businessOrgId ? `/venue/${businessOrgId}` : "/business/setup"}
              className={pageLinkAccent}
            >
              {businessOrgId ? "Venue profile" : "Set up venue"} →
            </Link>
          </div>
        ) : null}
      </div>
    </PageLayout>
  );
}
