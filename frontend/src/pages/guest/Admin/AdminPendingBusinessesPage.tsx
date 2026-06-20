import { useEffect, useState } from "react";
import { Building2, CheckCircle2, Loader2, XCircle } from "lucide-react";
import { toast } from "sonner";
import { apiClient } from "@/api/client";
import type { PendingBusinessListItem } from "@/types/api";
import { PageLayout } from "@/components/layout/PageLayout";
import { pageBtnPrimary, pageEmpty, pageLoader } from "@/components/layout/pageStyles";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Textarea } from "@/components/ui/textarea";
import { cn } from "@/lib/utils";

const PAGE_SIZE = 10;

function businessInitial(name: string) {
  return name.trim().charAt(0).toUpperCase() || "?";
}

export function AdminPendingBusinessesPage() {
  const [businesses, setBusinesses] = useState<PendingBusinessListItem[]>([]);
  const [totalCount, setTotalCount] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [page, setPage] = useState(0);
  const [actingId, setActingId] = useState<number | null>(null);

  const [rejectTarget, setRejectTarget] = useState<PendingBusinessListItem | null>(null);
  const [rejectComment, setRejectComment] = useState("");

  useEffect(() => {
    let mounted = true;
    const load = async () => {
      setLoading(true);
      setError("");
      try {
        const data = await apiClient.adminGetPendingBusinesses({
          limit: PAGE_SIZE,
          offset: page * PAGE_SIZE,
        });
        if (mounted) {
          setBusinesses(data.items ?? []);
          setTotalCount(data.total_count ?? 0);
        }
      } catch (err: unknown) {
        const e = err as { response?: { data?: { error?: string } }; message?: string };
        if (mounted) {
          setError(
            e?.response?.data?.error || e?.message || "Failed to load pending businesses."
          );
        }
      } finally {
        if (mounted) setLoading(false);
      }
    };
    void load();
    return () => {
      mounted = false;
    };
  }, [page]);

  const removeFromList = (id: number) => {
    setBusinesses((prev) => prev.filter((b) => b.id !== id));
    setTotalCount((prev) => Math.max(0, prev - 1));
  };

  const handleApprove = async (business: PendingBusinessListItem) => {
    setActingId(business.id);
    try {
      await apiClient.adminReviewBusiness(business.id, "verified");
      toast.success(`"${business.name}" has been verified.`);
      removeFromList(business.id);
    } catch (err: unknown) {
      const e = err as { response?: { data?: { error?: string } }; message?: string };
      toast.error(e?.response?.data?.error || e?.message || "Failed to verify business.");
    } finally {
      setActingId(null);
    }
  };

  const handleReject = async () => {
    if (!rejectTarget) return;
    const comment = rejectComment.trim();
    if (!comment) return;
    setActingId(rejectTarget.id);
    try {
      await apiClient.adminReviewBusiness(rejectTarget.id, "rejected", comment);
      toast.success(`"${rejectTarget.name}" has been rejected.`);
      removeFromList(rejectTarget.id);
      setRejectTarget(null);
      setRejectComment("");
    } catch (err: unknown) {
      const e = err as { response?: { data?: { error?: string } }; message?: string };
      toast.error(e?.response?.data?.error || e?.message || "Failed to reject business.");
    } finally {
      setActingId(null);
    }
  };

  const totalPages = Math.max(1, Math.ceil(totalCount / PAGE_SIZE));
  const currentPage = page + 1;

  return (
    <PageLayout className="space-y-8">
      <div>
        <h1 className="mb-3 text-4xl font-bold tracking-tight text-[#1A3C34] dark:text-white md:text-5xl">
          Admin — Verify Businesses{" "}
          <span className="text-emerald-500 dark:text-[#98FF98]">🏢</span>
        </h1>
        <p className="text-lg text-gray-600 dark:text-gray-400">
          Review and approve business accounts awaiting verification
        </p>
      </div>

      <div className="flex items-center gap-2 text-sm text-gray-600 dark:text-gray-300">
        <span className="inline-flex items-center gap-1 rounded-full border border-[#FFD700]/40 bg-[#FFD700]/10 px-3 py-1 font-semibold text-[#FFD700]">
          {totalCount} pending
        </span>
      </div>

      {error ? (
        <div className="rounded-xl border border-red-500/30 bg-red-500/10 px-4 py-3 text-sm font-medium text-red-400">
          {error}
        </div>
      ) : null}

      {loading ? (
        <div className="flex h-44 items-center justify-center">
          <Loader2 className={cn(pageLoader, "h-10 w-10")} />
        </div>
      ) : businesses.length === 0 ? (
        <div className={pageEmpty}>
          <p className="text-xl font-bold text-[#1A3C34] dark:text-gray-200">
            No businesses awaiting review
          </p>
          <p className="mt-2 text-gray-500">All caught up — nothing to verify right now.</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 gap-6 md:grid-cols-2">
          {businesses.map((business) => {
            const busy = actingId === business.id;
            return (
              <Card
                key={business.id}
                className="rounded-3xl border border-gray-200 bg-white shadow-sm transition-all hover:shadow-lg dark:border-[#2f5e50] dark:bg-[#163d32]"
              >
                <CardContent className="space-y-4 p-5">
                  <div className="flex gap-4">
                    <div className="flex h-16 w-16 shrink-0 items-center justify-center overflow-hidden rounded-2xl border border-gray-200 bg-[#163d32] text-2xl font-bold text-[#98FF98] dark:border-[#2f5e50]">
                      {business.avatar ? (
                        <img
                          src={business.avatar}
                          alt=""
                          className="h-full w-full object-cover"
                        />
                      ) : (
                        businessInitial(business.name)
                      )}
                    </div>
                    <div className="min-w-0 flex-1">
                      <h3 className="truncate text-lg font-bold text-[#1A3C34] dark:text-white">
                        {business.name}
                      </h3>
                      <p className="mt-1 inline-flex items-center gap-1 text-sm text-gray-500 dark:text-gray-400">
                        <Building2 className="h-3.5 w-3.5" />
                        Org #{business.id}
                      </p>
                    </div>
                  </div>

                  {business.description ? (
                    <p className="line-clamp-3 text-sm text-gray-600 dark:text-gray-300">
                      {business.description}
                    </p>
                  ) : (
                    <p className="text-sm italic text-gray-400">No description provided.</p>
                  )}

                  <div className="flex items-center gap-3 border-t border-gray-100 pt-4 dark:border-[#2f5e50]">
                    <Button
                      type="button"
                      disabled={busy}
                      onClick={() => void handleApprove(business)}
                      className="flex-1 rounded-xl bg-emerald-600 font-semibold text-white hover:bg-emerald-700"
                    >
                      {busy ? (
                        <Loader2 className="h-4 w-4 animate-spin" />
                      ) : (
                        <>
                          <CheckCircle2 className="h-4 w-4" />
                          Approve
                        </>
                      )}
                    </Button>
                    <Button
                      type="button"
                      variant="outline"
                      disabled={busy}
                      onClick={() => {
                        setRejectTarget(business);
                        setRejectComment("");
                      }}
                      className="flex-1 rounded-xl border-red-500/40 font-semibold text-red-500 hover:bg-red-500/10 hover:text-red-400"
                    >
                      <XCircle className="h-4 w-4" />
                      Reject
                    </Button>
                  </div>
                </CardContent>
              </Card>
            );
          })}
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

      <Dialog
        open={rejectTarget !== null}
        onOpenChange={(open) => {
          if (!open) {
            setRejectTarget(null);
            setRejectComment("");
          }
        }}
      >
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>Reject business</DialogTitle>
            <DialogDescription>
              Rejecting <span className="font-medium text-foreground">{rejectTarget?.name}</span>{" "}
              requires a reason. It will be recorded with the review.
            </DialogDescription>
          </DialogHeader>
          <div className="space-y-1.5">
            <Textarea
              placeholder="Explain why this business is being rejected…"
              value={rejectComment}
              onChange={(e) => setRejectComment(e.target.value)}
              maxLength={1000}
              rows={4}
              autoFocus
              className="min-h-24 resize-none"
            />
            <p className="text-right text-xs text-muted-foreground">
              {rejectComment.length}/1000
            </p>
          </div>
          <DialogFooter className="sm:items-center sm:justify-evenly">
            <Button
              type="button"
              variant="outline"
              className="min-w-24"
              onClick={() => {
                setRejectTarget(null);
                setRejectComment("");
              }}
            >
              Cancel
            </Button>
            <Button
              type="button"
              className="min-w-36 bg-red-600 text-white hover:bg-red-700 disabled:bg-red-600/40 disabled:text-white/70"
              disabled={!rejectComment.trim() || actingId === rejectTarget?.id}
              onClick={() => void handleReject()}
            >
              {actingId === rejectTarget?.id ? (
                <Loader2 className="h-4 w-4 animate-spin" />
              ) : (
                "Confirm rejection"
              )}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </PageLayout>
  );
}
