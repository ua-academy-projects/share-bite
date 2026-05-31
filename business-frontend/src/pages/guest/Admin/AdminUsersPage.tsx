import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { Calendar, Filter, Loader2, Mail, Shield, User } from "lucide-react";
import { apiClient } from "@/api/client";
import type { AdminUserListItem } from "@/types/api";
import { PageLayout } from "@/components/layout/PageLayout";
import {
  pageBtnPrimary,
  pageEmpty,
  pageFilterControl,
  pageLoader,
} from "@/components/layout/pageStyles";
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

const PAGE_SIZE = 10;
const ROLE_OPTIONS = ["user", "business", "moderator", "admin"];
const STATUS_OPTIONS = ["active", "muted", "suspended"];

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

function userInitial(email: string) {
  return email.charAt(0).toUpperCase();
}

export function AdminUsersPage() {
  const [users, setUsers] = useState<AdminUserListItem[]>([]);
  const [totalCount, setTotalCount] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [search, setSearch] = useState("");
  const [debouncedSearch, setDebouncedSearch] = useState("");
  const [roleFilter, setRoleFilter] = useState("all");
  const [statusFilter, setStatusFilter] = useState("all");
  const [page, setPage] = useState(0);
  const [sortOrder, setSortOrder] = useState<"asc" | "desc">("desc");

  useEffect(() => {
    const t = setTimeout(() => setDebouncedSearch(search), 350);
    return () => clearTimeout(t);
  }, [search]);

  useEffect(() => {
    setPage(0);
  }, [debouncedSearch, roleFilter, statusFilter, sortOrder]);

  useEffect(() => {
    let mounted = true;
    const load = async () => {
      setLoading(true);
      setError("");
      try {
        const data = await apiClient.adminGetUsers({
          limit: PAGE_SIZE,
          offset: page * PAGE_SIZE,
          search_email: debouncedSearch || undefined,
          role: roleFilter === "all" ? undefined : roleFilter,
          status: statusFilter === "all" ? undefined : statusFilter,
          sort_order: sortOrder,
        });
        if (mounted) {
          setUsers(data.items ?? []);
          setTotalCount(data.total_count ?? 0);
        }
      } catch (err: unknown) {
        const e = err as { response?: { data?: { error?: string } }; message?: string };
        if (mounted) {
          setError(e?.response?.data?.error || e?.message || "Failed to load users.");
        }
      } finally {
        if (mounted) setLoading(false);
      }
    };
    void load();
    return () => {
      mounted = false;
    };
  }, [page, debouncedSearch, roleFilter, statusFilter, sortOrder]);

  const totalPages = Math.max(1, Math.ceil(totalCount / PAGE_SIZE));
  const currentPage = page + 1;
  const filterControl = cn(pageFilterControl, "data-[size=default]:h-11");

  return (
    <PageLayout className="space-y-8">
      <div>
        <h1 className="mb-3 text-4xl font-bold tracking-tight text-[#1A3C34] dark:text-white md:text-5xl">
          Admin — Users{" "}
          <span className="text-emerald-500 dark:text-[#98FF98]">🛡️</span>
        </h1>
        <p className="text-lg text-gray-600 dark:text-gray-400">
          Manage accounts and permissions
        </p>
      </div>

      <Card className="rounded-3xl border border-gray-200 bg-white shadow-sm dark:border-[#2f5e50] dark:bg-[#163d32]">
        <CardContent className="space-y-4 p-6">
          <div className="flex items-center gap-2 font-semibold text-[#1A3C34] dark:text-white">
            <Filter size={20} className="text-emerald-500 dark:text-[#98FF98]" />
            <span>Filters</span>
          </div>
          <div className="grid gap-3 md:grid-cols-3">
            <input
              placeholder="Search by email…"
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className={filterControl}
            />
            <Select value={roleFilter} onValueChange={setRoleFilter}>
              <SelectTrigger className={cn(filterControl, "cursor-pointer")}>
                <SelectValue placeholder="All roles" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All roles</SelectItem>
                {ROLE_OPTIONS.map((r) => (
                  <SelectItem key={r} value={r}>
                    {r}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            <Select value={statusFilter} onValueChange={setStatusFilter}>
              <SelectTrigger className={cn(filterControl, "cursor-pointer")}>
                <SelectValue placeholder="All statuses" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All statuses</SelectItem>
                {STATUS_OPTIONS.map((s) => (
                  <SelectItem key={s} value={s}>
                    {s}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
          <div className="flex items-center justify-between text-sm">
            <p className="text-gray-600 dark:text-gray-300">
              Results: <span className="font-semibold">{totalCount}</span>
            </p>
            <Button
              type="button"
              variant="ghost"
              size="sm"
              className="text-gray-500 hover:text-emerald-600 dark:hover:text-[#98FF98]"
              onClick={() => setSortOrder((prev) => (prev === "asc" ? "desc" : "asc"))}
            >
              Sort by created {sortOrder === "asc" ? "↑" : "↓"}
            </Button>
          </div>
        </CardContent>
      </Card>

      {error ? (
        <div className="rounded-xl border border-red-500/30 bg-red-500/10 px-4 py-3 text-sm font-medium text-red-400">
          {error}
        </div>
      ) : null}

      {loading ? (
        <div className="flex h-44 items-center justify-center">
          <Loader2 className={cn(pageLoader, "h-10 w-10")} />
        </div>
      ) : users.length === 0 ? (
        <div className={pageEmpty}>
          <p className="text-xl font-bold text-[#1A3C34] dark:text-gray-200">No users found</p>
          <p className="mt-2 text-gray-500">Try adjusting your filters.</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 gap-6 md:grid-cols-2">
          {users.map((user) => (
            <Card
              key={user.id}
              className="rounded-3xl border border-gray-200 bg-white shadow-sm transition-all hover:shadow-lg dark:border-[#2f5e50] dark:bg-[#163d32]"
            >
              <CardContent className="space-y-4 p-5">
                <div className="flex gap-4">
                  <div className="flex h-16 w-16 shrink-0 items-center justify-center rounded-2xl border border-gray-200 bg-[#163d32] text-2xl font-bold text-[#98FF98] dark:border-[#2f5e50]">
                    {userInitial(user.email)}
                  </div>
                  <div className="min-w-0 flex-1">
                    <h3 className="truncate text-lg font-bold text-[#1A3C34] dark:text-white">
                      {user.email}
                    </h3>
                    <p className="mt-1 inline-flex items-center gap-1 text-sm text-gray-500 dark:text-gray-400">
                      <Calendar className="h-3.5 w-3.5" />
                      Joined {new Date(user.created_at).toLocaleDateString()}
                    </p>
                  </div>
                </div>

                <div className="flex flex-wrap gap-2">
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

                <div className="flex items-center justify-between">
                  <div className="inline-flex items-center gap-1 text-sm text-gray-500 dark:text-gray-300">
                    <Mail className="h-4 w-4" />
                    User ID: {user.id.slice(0, 8)}…
                  </div>
                  <Button
                    asChild
                    className="rounded-xl bg-[#163d32] px-5 font-semibold text-white ring-1 ring-[#FFD700]/40 hover:bg-[#1A3C34]"
                  >
                    <Link to={`/admin/users/${user.id}`}>View details</Link>
                  </Button>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}

      {totalPages > 1 ? (
        <div className="flex items-center justify-between border-t border-gray-200 py-6 dark:border-[#2f5e50]">
          <Button
            type="button"
            variant="outline"
            disabled={page === 0}
            onClick={() => setPage((p) => Math.max(0, p - 1))}
            className="rounded-full"
          >
            Previous
          </Button>
          <p className="text-sm text-gray-500 dark:text-gray-400">
            Page {currentPage} of {totalPages}
          </p>
          <Button
            type="button"
            className={pageBtnPrimary}
            disabled={page >= totalPages - 1}
            onClick={() => setPage((p) => Math.min(totalPages - 1, p + 1))}
          >
            Next →
          </Button>
        </div>
      ) : null}
    </PageLayout>
  );
}
