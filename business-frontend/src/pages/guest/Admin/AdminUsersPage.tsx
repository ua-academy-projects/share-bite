import { useCallback, useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { Loader2, Shield } from "lucide-react";
import { apiClient } from "@/api/client";
import type { AdminUserListItem } from "@/types/api";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Card, CardContent } from "@/components/ui/card";
import { PageHeader } from "@/components/layout/PageHeader";

const PAGE_SIZE = 10;
const ROLES = ["all", "user", "business", "moderator", "admin"];
const STATUSES = ["all", "active", "muted", "suspended"];

function roleVariant(role: string) {
  if (role === "admin") return "destructive" as const;
  if (role === "moderator") return "accent" as const;
  if (role === "business") return "secondary" as const;
  return "outline" as const;
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

  const load = useCallback(async () => {
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
      setUsers(data.items ?? []);
      setTotalCount(data.total_count ?? 0);
    } catch (err: unknown) {
      const e = err as { response?: { data?: { error?: string } }; message?: string };
      setError(e?.response?.data?.error || e?.message || "Failed to load users.");
    } finally {
      setLoading(false);
    }
  }, [page, debouncedSearch, roleFilter, statusFilter, sortOrder]);

  useEffect(() => {
    setPage(0);
  }, [debouncedSearch, roleFilter, statusFilter, sortOrder]);

  useEffect(() => {
    void load();
  }, [load]);

  const totalPages = Math.max(1, Math.ceil(totalCount / PAGE_SIZE));

  return (
    <div className="px-6 py-8 lg:px-10">
      <PageHeader
        title="Admin — Users"
        description={`${totalCount} total users`}
        icon={Shield}
      />

      <Card className="mb-6 rounded-2xl bg-card-solid">
        <CardContent className="flex flex-wrap gap-3 p-4">
          <Input
            placeholder="Search by email…"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="max-w-xs rounded-xl"
          />
          <Select value={roleFilter} onValueChange={setRoleFilter}>
            <SelectTrigger className="w-36 rounded-xl">
              <SelectValue placeholder="Role" />
            </SelectTrigger>
            <SelectContent>
              {ROLES.map((r) => (
                <SelectItem key={r} value={r}>
                  {r === "all" ? "All roles" : r}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          <Select value={statusFilter} onValueChange={setStatusFilter}>
            <SelectTrigger className="w-36 rounded-xl">
              <SelectValue placeholder="Status" />
            </SelectTrigger>
            <SelectContent>
              {STATUSES.map((s) => (
                <SelectItem key={s} value={s}>
                  {s === "all" ? "All statuses" : s}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          <Button
            variant="outline"
            className="rounded-xl"
            onClick={() => setSortOrder((o) => (o === "asc" ? "desc" : "asc"))}
          >
            Sort: {sortOrder}
          </Button>
        </CardContent>
      </Card>

      {loading ? (
        <div className="flex justify-center py-16">
          <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
        </div>
      ) : error ? (
        <p className="text-destructive">{error}</p>
      ) : (
        <div className="space-y-3">
          {users.map((user) => (
            <Card key={user.id} className="rounded-2xl bg-card-solid">
              <CardContent className="flex flex-wrap items-center justify-between gap-3 p-4">
                <div>
                  <p className="font-medium">{user.email}</p>
                  <p className="text-xs text-muted-foreground">
                    {new Date(user.created_at).toLocaleDateString()}
                  </p>
                </div>
                <div className="flex items-center gap-2">
                  <Badge variant={roleVariant(user.role_slug)}>
                    {user.role_slug}
                  </Badge>
                  <Badge variant="outline">{user.status}</Badge>
                  <Button asChild variant="outline" size="sm" className="rounded-xl">
                    <Link to={`/admin/users/${user.id}`}>View</Link>
                  </Button>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}

      <div className="mt-6 flex items-center justify-center gap-3">
        <Button
          variant="outline"
          disabled={page === 0}
          onClick={() => setPage((p) => p - 1)}
        >
          Previous
        </Button>
        <span className="text-sm text-muted-foreground">
          Page {page + 1} of {totalPages}
        </span>
        <Button
          variant="outline"
          disabled={page + 1 >= totalPages}
          onClick={() => setPage((p) => p + 1)}
        >
          Next
        </Button>
      </div>
    </div>
  );
}
