import { useEffect, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { apiClient } from "@/api/client";
import { PageHeader } from "@/components/layout/PageHeader";
import { PageLayout } from "@/components/layout/PageLayout";
import {
  pageBtnPrimary,
  pageInput,
  pageLabel,
  pagePanel,
} from "@/components/layout/pageStyles";
import { Button } from "@/components/ui/button";
import { getDefaultHomePath, resolvePostAuthDestination } from "@/utils/navigation";
import { cn } from "@/lib/utils";

export function ProfileCreatePage() {
  const navigate = useNavigate();
  const location = useLocation();
  const queryClient = useQueryClient();
  const [form, setForm] = useState({
    userName: "",
    firstName: "",
    lastName: "",
    bio: "",
  });

  const { data: hasCustomer, isLoading } = useQuery({
    queryKey: ["hasCustomer"],
    queryFn: () => apiClient.hasCurrentCustomer(),
    retry: false,
  });

  useEffect(() => {
    if (hasCustomer) {
      const from = (location.state as { from?: { pathname?: string } } | null)?.from
        ?.pathname;
      navigate(resolvePostAuthDestination(from), { replace: true });
    }
  }, [hasCustomer, location.state, navigate]);

  const createMutation = useMutation({
    mutationFn: () => apiClient.createCustomer(form),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["currentCustomer"] });
      void queryClient.invalidateQueries({ queryKey: ["onboardingStatus"] });
      void queryClient.invalidateQueries({ queryKey: ["hasCustomer"] });
      toast.success("Profile created successfully!");
      navigate(getDefaultHomePath(), { replace: true });
    },
    onError: (error: unknown) => {
      const e = error as { response?: { data?: { error?: string } } };
      toast.error(e?.response?.data?.error || "Failed to create profile");
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    createMutation.mutate();
  };

  if (isLoading || hasCustomer) {
    return (
      <PageLayout center>
        <p className="text-gray-500">Loading…</p>
      </PageLayout>
    );
  }

  return (
    <PageLayout>
      <PageHeader
        title="Create Profile"
        description="Set up your guest profile to start sharing bites"
      />
      <div className={cn(pagePanel, "mx-auto max-w-lg p-8")}>
        <form onSubmit={handleSubmit} className="space-y-5">
          <div className="space-y-2">
            <label htmlFor="userName" className={pageLabel}>
              Username
            </label>
            <input
              id="userName"
              required
              value={form.userName}
              onChange={(e) => setForm((f) => ({ ...f, userName: e.target.value }))}
              className={pageInput}
            />
          </div>
          <div className="space-y-2">
            <label htmlFor="firstName" className={pageLabel}>
              First name
            </label>
            <input
              id="firstName"
              required
              value={form.firstName}
              onChange={(e) => setForm((f) => ({ ...f, firstName: e.target.value }))}
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
              onChange={(e) => setForm((f) => ({ ...f, lastName: e.target.value }))}
              className={pageInput}
            />
          </div>
          <div className="space-y-2">
            <label htmlFor="bio" className={pageLabel}>
              Bio
            </label>
            <textarea
              id="bio"
              value={form.bio}
              onChange={(e) => setForm((f) => ({ ...f, bio: e.target.value }))}
              className={cn(pageInput, "min-h-[100px] resize-y py-3")}
            />
          </div>
          <Button
            type="submit"
            className={cn(pageBtnPrimary, "w-full")}
            disabled={createMutation.isPending}
          >
            {createMutation.isPending ? "Creating…" : "Create Profile"}
          </Button>
        </form>
      </div>
    </PageLayout>
  );
}
