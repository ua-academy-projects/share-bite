import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { Filter, Loader2 } from "lucide-react";
import { apiClient } from "@/api/client";
import type { AdminUserListItem } from "@/types/api";
import { PageHeader } from "@/components/layout/PageHeader";
import { PageLayout } from "@/components/layout/PageLayout";
import {
  pageBtnPrimary,
  pageBtnSecondary,
  pageEmpty,
  pageFilterBar,
  pageInput,
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

const PAGE_SIZE = 10;
const ROLE_OPTIONS = ["user", "business", "moderator", "admin"];
const STATUS_OPTIONS = ["active", "muted", "suspended"];

function roleBadgeVariant(role: string) {
  if (role === "admin") return "destructive";
  if (role === "moderator") return "secondary";
  return "outline";
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
    load();
    return () => {
      mounted = false;
    };
  }, [page, debouncedSearch, roleFilter, statusFilter, sortOrder]);

  const totalPages = Math.max(1, Math.ceil(totalCount / PAGE_SIZE));

  return (
    <PageLayout maxWidth="7xl">
      <PageHeader title="Admin — Users" description={`${totalCount} total users`} />

      <div className={cn(pageFilterBar, "mb-8 flex flex-col gap-4 xl:flex-row xl:items-center")}>
        <div className="flex items-center gap-2 whitespace-nowrap font-semibold text-[#1A3C34] dark:text-white">
          <Filter size={20} className="text-emerald-500 dark:text-[#98FF98]" />
          <span>Filters:</span>
        </div>
        <div className="grid flex-1 gap-3 md:grid-cols-3">
          <input
            placeholder="Search by email…"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className={pageInput}
          />
          <Select value={roleFilter} onValueChange={setRoleFilter}>
            <SelectTrigger className={cn(pageInput, "h-auto cursor-pointer")}>
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
            <SelectTrigger className={cn(pageInput, "h-auto cursor-pointer")}>
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
      </div>

      {error ? <p className="mb-4 text-sm text-destructive">{error}</p> : null}

      {loading ? (
        <div className="flex h-64 items-center justify-center">
          <Loader2 className={cn(pageLoader, "h-12 w-12")} />
        </div>
      ) : users.length === 0 ? (
        <div className={pageEmpty}>
          <p className="text-xl font-bold text-[#1A3C34] dark:text-gray-300">No users found</p>
        </div>
      ) : (
        <div className={cn(pagePanel, "overflow-hidden")}>
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead className="border-b border-gray-200 bg-gray-50 dark:border-[#2f5e50] dark:bg-[#0d241d]/50">
                <tr>
                  <th className="px-4 py-3 text-left font-medium text-[#1A3C34] dark:text-white">
                    Email
                  </th>
                  <th className="px-4 py-3 text-left font-medium text-[#1A3C34] dark:text-white">
                    Role
                  </th>
                  <th className="px-4 py-3 text-left font-medium text-[#1A3C34] dark:text-white">
                    Status
                  </th>
                  <th className="px-4 py-3 text-left font-medium text-[#1A3C34] dark:text-white">
                    <Button
                      type="button"
                      variant="ghost"
                      size="sm"
                      className="h-auto px-0 text-left hover:text-emerald-600 dark:hover:text-[#98FF98]"
                      onClick={() =>
                        setSortOrder((prev) => (prev === "asc" ? "desc" : "asc"))
                      }
                    >
                      Created {sortOrder === "asc" ? "↑" : "↓"}
                    </Button>
                  </th>
                  <th className="px-4 py-3 text-right font-medium text-[#1A3C34] dark:text-white">
                    Actions
                  </th>
                </tr>
              </thead>
              <tbody>
                {users.map((user) => (
                  <tr
                    key={user.id}
                    className="border-b border-gray-200 dark:border-[#2f5e50]/60"
                  >
                    <td className="px-4 py-3 text-[#1A3C34] dark:text-gray-200">
                      {user.email}
                    </td>
                    <td className="px-4 py-3">
                      <Badge variant={roleBadgeVariant(user.role_slug)}>
                        {user.role_slug}
                      </Badge>
                    </td>
                    <td className="px-4 py-3">
                      <Badge variant="outline">{user.status}</Badge>
                    </td>
                    <td className="px-4 py-3 text-gray-500 dark:text-gray-400">
                      {new Date(user.created_at).toLocaleDateString()}
                    </td>
                    <td className="px-4 py-3 text-right">
                      <Link to={`/admin/users/${user.id}`} className={pageLinkAccent}>
                        Details →
                      </Link>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}

      {totalPages > 1 ? (
        <div className="mt-6 flex items-center justify-center gap-4">
          <Button
            className={pageBtnSecondary}
            onClick={() => setPage((p) => Math.max(0, p - 1))}
            disabled={page === 0}
          >
            ← Prev
          </Button>
          <span className="text-sm text-gray-500 dark:text-gray-400">
            Page {page + 1} of {totalPages}
          </span>
          <Button
            className={pageBtnPrimary}
            onClick={() => setPage((p) => Math.min(totalPages - 1, p + 1))}
            disabled={page >= totalPages - 1}
          >
            Next →
          </Button>
        </div>
      ) : null}
    </PageLayout>
  );
}
