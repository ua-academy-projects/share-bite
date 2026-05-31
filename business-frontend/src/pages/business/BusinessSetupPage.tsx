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
    venueName: "",
    venueDescription: "",
    latitude: "1",
    longitude: "2",
  });

  const { data: context, isLoading } = useQuery({
    queryKey: ["businessOnboarding"],
    queryFn: () => businessApi.getMyOnboardingContext(token),
    enabled: !!token,
    retry: false,
  });

  useEffect(() => {
    if (context?.venueId) {
      setBusinessOrgId(context.venueId);
      navigate(getDefaultHomePath(), { replace: true });
    }
  }, [context?.venueId, navigate]);

  const setupMutation = useMutation({
    mutationFn: async () => {
      let brandId = context?.brandId ?? null;
      if (!brandId) {
        const brand = await businessApi.createBrand(
          {
            name: form.brandName.trim(),
            description: form.brandDescription.trim() || undefined,
          },
          token
        );
        brandId = brand.id;
      }

      const venue = await businessApi.createLocation(
        brandId,
        {
          name: form.venueName.trim(),
          description: form.venueDescription.trim() || undefined,
          latitude: Number(form.latitude),
          longitude: Number(form.longitude),
        },
        token
      );

      return venue.id;
    },
    onSuccess: (venueId) => {
      setBusinessOrgId(venueId);
      void queryClient.invalidateQueries({ queryKey: ["onboardingStatus"] });
      void queryClient.invalidateQueries({ queryKey: ["businessOnboarding"] });
      toast.success("Venue profile created!");
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

  const needsBrand = !context?.brandId;

  return (
    <PageLayout>
      <div className="mx-auto max-w-lg space-y-8">
        <div>
          <h1 className="mb-3 text-4xl font-bold tracking-tight text-[#1A3C34] dark:text-white md:text-5xl">
            Set Up Your Venue
          </h1>
          <p className="text-lg text-gray-600 dark:text-gray-400">
            Create your business profile to list rescue boxes and post updates.
          </p>
        </div>

        <form onSubmit={handleSubmit} className={cn(pagePanel, "space-y-5 p-8")}>
          {needsBrand ? (
            <>
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
                  placeholder="My Restaurant Group"
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
                />
              </div>
            </>
          ) : null}

          <div className="space-y-2">
            <label htmlFor="venueName" className={pageLabel}>
              Venue name
            </label>
            <input
              id="venueName"
              required
              value={form.venueName}
              onChange={(e) => setForm((f) => ({ ...f, venueName: e.target.value }))}
              className={pageInput}
              placeholder="Downtown location"
            />
          </div>
          <div className="space-y-2">
            <label htmlFor="venueDescription" className={pageLabel}>
              Venue description
            </label>
            <textarea
              id="venueDescription"
              value={form.venueDescription}
              onChange={(e) =>
                setForm((f) => ({ ...f, venueDescription: e.target.value }))
              }
              className={cn(pageInput, "min-h-[80px] resize-y py-3")}
            />
          </div>
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <label htmlFor="latitude" className={pageLabel}>
                Latitude
              </label>
              <input
                id="latitude"
                required
                type="number"
                step="any"
                value={form.latitude}
                onChange={(e) => setForm((f) => ({ ...f, latitude: e.target.value }))}
                className={pageInput}
              />
            </div>
            <div className="space-y-2">
              <label htmlFor="longitude" className={pageLabel}>
                Longitude
              </label>
              <input
                id="longitude"
                required
                type="number"
                step="any"
                value={form.longitude}
                onChange={(e) => setForm((f) => ({ ...f, longitude: e.target.value }))}
                className={pageInput}
              />
            </div>
          </div>

          <Button
            type="submit"
            className={cn(pageBtnPrimary, "w-full")}
            disabled={setupMutation.isPending}
          >
            {setupMutation.isPending ? "Creating…" : "Create venue profile"}
          </Button>
        </form>
      </div>
    </PageLayout>
  );
}
