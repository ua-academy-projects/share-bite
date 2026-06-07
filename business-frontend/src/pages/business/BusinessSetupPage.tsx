import { useEffect, useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useNavigate } from "react-router-dom";
import { toast } from "sonner";
import { businessApi } from "@/api/business";
import { PageLayout } from "@/components/layout/PageLayout";
import {
  pageBtnPrimary,
  pageInput,
  pageLabel,
  pagePanel,
} from "@/components/layout/pageStyles";
import { Button } from "@/components/ui/button";
import { setBusinessOrgId } from "@/utils/auth";
import { getDefaultHomePath } from "@/utils/navigation";
import { cn } from "@/lib/utils";

export function BusinessSetupPage() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const token = localStorage.getItem("token") || "";

  const [form, setForm] = useState({
    brandName: "",
    brandDescription: "",
  });

  const { data: context, isLoading } = useQuery({
    queryKey: ["businessOnboarding"],
    queryFn: () => businessApi.getMyOnboardingContext(token),
    enabled: !!token,
    retry: false,
  });

  useEffect(() => {
    if (context?.brandId) {
      setBusinessOrgId(context.brandId);
      navigate(getDefaultHomePath(), { replace: true });
    }
  }, [context?.brandId, navigate]);

  const setupMutation = useMutation({
    mutationFn: async () => {
      const brand = await businessApi.createBrand(
        {
          name: form.brandName.trim(),
          description: form.brandDescription.trim() || undefined,
        },
        token
      );
      return brand.id;
    },
    onSuccess: (brandId) => {
      setBusinessOrgId(brandId);
      void queryClient.invalidateQueries({ queryKey: ["onboardingStatus"] });
      void queryClient.invalidateQueries({ queryKey: ["businessOnboarding"] });
      toast.success("Business profile created!");
      navigate(getDefaultHomePath(), { replace: true });
    },
    onError: (error: unknown) => {
      const message = error instanceof Error ? error.message : "Setup failed";
      toast.error(message);
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setupMutation.mutate();
  };

  if (isLoading) {
    return (
      <PageLayout center>
        <p className="text-gray-500">Loading…</p>
      </PageLayout>
    );
  }

  return (
    <PageLayout>
      <div className="mx-auto max-w-lg space-y-8">
        <div>
          <h1 className="mb-3 text-4xl font-bold tracking-tight text-[#1A3C34] dark:text-white md:text-5xl">
            Set Up Your Business
          </h1>
          <p className="text-lg text-gray-600 dark:text-gray-400">
            Create your brand profile to get started on Share Bite.
          </p>
        </div>

        <form onSubmit={handleSubmit} className={cn(pagePanel, "space-y-5 p-8")}>
          <div className="space-y-2">
            <label htmlFor="brandName" className={pageLabel}>
              Brand name
            </label>
            <input
              id="brandName"
              required
              value={form.brandName}
              onChange={(e) => setForm((f) => ({ ...f, brandName: e.target.value }))}
              className={pageInput}
              placeholder="Downtown Bakery"
            />
          </div>
          <div className="space-y-2">
            <label htmlFor="brandDescription" className={pageLabel}>
              Brand description
            </label>
            <textarea
              id="brandDescription"
              value={form.brandDescription}
              onChange={(e) =>
                setForm((f) => ({ ...f, brandDescription: e.target.value }))
              }
              className={cn(pageInput, "min-h-[80px] resize-y py-3")}
              placeholder="Tell customers about your business"
            />
          </div>

          <Button
            type="submit"
            className={cn(pageBtnPrimary, "w-full")}
            disabled={setupMutation.isPending}
          >
            {setupMutation.isPending ? "Creating…" : "Create business profile"}
          </Button>
        </form>
      </div>
    </PageLayout>
  );
}
