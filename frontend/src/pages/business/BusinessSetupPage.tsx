import { useEffect, useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useNavigate } from "react-router-dom";
import { toast } from "sonner";
import { MapPin } from "lucide-react";
import { businessApi, type LocationTag } from "@/api/business";
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

type SetupForm = {
  brandName: string;
  brandDescription: string;
  venueName: string;
  venueDescription: string;
  latitude: string;
  longitude: string;
  tagIds: number[];
};

const emptyForm: SetupForm = {
  brandName: "",
  brandDescription: "",
  venueName: "",
  venueDescription: "",
  latitude: "",
  longitude: "",
  tagIds: [],
};

function parseCoordinate(value: string): number | undefined {
  const trimmed = value.trim();
  if (!trimmed) return undefined;
  const parsed = Number(trimmed);
  return Number.isFinite(parsed) ? parsed : undefined;
}

export function BusinessSetupPage() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const token = localStorage.getItem("token") || "";

  const [form, setForm] = useState<SetupForm>(emptyForm);

  const { data: context, isLoading: contextLoading } = useQuery({
    queryKey: ["businessOnboarding"],
    queryFn: () => businessApi.getMyOnboardingContext(token),
    enabled: !!token,
    retry: false,
  });

  const existingBrandId = context?.brandId && context.brandId > 0 ? context.brandId : null;
  const hasVenue = context?.venueId != null && context.venueId > 0;

  const { data: existingBrand, isLoading: brandLoading } = useQuery({
    queryKey: ["businessBrand", existingBrandId],
    queryFn: () => businessApi.getBrand(existingBrandId!, token),
    enabled: !!token && !!existingBrandId,
    retry: false,
  });

  const { data: locationTags = [], isLoading: tagsLoading } = useQuery({
    queryKey: ["locationTags"],
    queryFn: () => businessApi.getLocationTags(),
    staleTime: 5 * 60_000,
  });

  useEffect(() => {
    if (hasVenue && context?.venueId) {
      setBusinessOrgId(context.venueId);
      navigate(getDefaultHomePath(), { replace: true });
    }
  }, [context?.venueId, hasVenue, navigate]);

  useEffect(() => {
    if (!existingBrand) return;
    setForm((current) => ({
      ...current,
      brandName: existingBrand.name,
      brandDescription: existingBrand.description || "",
    }));
  }, [existingBrand]);

  const setupMutation = useMutation({
    mutationFn: async () => {
      let brandId = existingBrandId;

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

      const latitude = parseCoordinate(form.latitude);
      const longitude = parseCoordinate(form.longitude);

      const venue = await businessApi.createLocation(
        brandId,
        {
          name: form.venueName.trim(),
          description: form.venueDescription.trim() || undefined,
          latitude,
          longitude,
          tagIds: form.tagIds.length > 0 ? form.tagIds : undefined,
        },
        token
      );

      return venue.id;
    },
    onSuccess: (venueId) => {
      setBusinessOrgId(venueId);
      void queryClient.invalidateQueries({ queryKey: ["onboardingStatus"] });
      void queryClient.invalidateQueries({ queryKey: ["businessOnboarding"] });
      toast.success("Business and first venue created!");
      navigate(getDefaultHomePath(), { replace: true });
    },
    onError: (error: unknown) => {
      const message = error instanceof Error ? error.message : "Setup failed";
      toast.error(message);
    },
  });

  const toggleTag = (tag: LocationTag) => {
    setForm((current) => {
      const hasTag = current.tagIds.includes(tag.id);
      if (hasTag) {
        return {
          ...current,
          tagIds: current.tagIds.filter((id) => id !== tag.id),
        };
      }
      if (current.tagIds.length >= 5) {
        toast.error("You can select at most 5 tags.");
        return current;
      }
      return {
        ...current,
        tagIds: [...current.tagIds, tag.id],
      };
    });
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!existingBrandId && !form.brandName.trim()) {
      toast.error("Brand name is required.");
      return;
    }
    if (!form.venueName.trim()) {
      toast.error("Venue name is required.");
      return;
    }
    setupMutation.mutate();
  };

  const isLoading = contextLoading || (existingBrandId != null && brandLoading);

  if (isLoading) {
    return (
      <PageLayout center>
        <p className="text-gray-500">Loading…</p>
      </PageLayout>
    );
  }

  const needsVenueOnly = existingBrandId != null;

  return (
    <PageLayout>
      <div className="mx-auto max-w-lg space-y-8">
        <div>
          <h1 className="mb-3 text-4xl font-bold tracking-tight text-[#1A3C34] dark:text-white md:text-5xl">
            {needsVenueOnly ? "Add Your First Venue" : "Set Up Your Business"}
          </h1>
          <p className="text-lg text-gray-600 dark:text-gray-400">
            {needsVenueOnly
              ? "Your brand is ready. Add your first location so you can post and list rescue boxes."
              : "Create your brand and first venue to get started on Share Bite."}
          </p>
        </div>

        <form onSubmit={handleSubmit} className={cn(pagePanel, "space-y-8 p-8")}>
          {!needsVenueOnly ? (
            <div className="space-y-5">
              <h2 className="text-lg font-semibold text-[#1A3C34] dark:text-white">Brand</h2>
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
            </div>
          ) : (
            <div className="rounded-xl border border-gray-200 bg-gray-50 p-4 dark:border-[#2f5e50] dark:bg-[#0d241d]">
              <p className="text-xs font-semibold uppercase tracking-wide text-emerald-700 dark:text-[#98FF98]">
                Your brand
              </p>
              <p className="mt-1 text-lg font-bold text-[#1A3C34] dark:text-white">
                {existingBrand?.name || form.brandName}
              </p>
              {(existingBrand?.description || form.brandDescription) && (
                <p className="mt-1 text-sm text-gray-600 dark:text-gray-400">
                  {existingBrand?.description || form.brandDescription}
                </p>
              )}
            </div>
          )}

          <div className="space-y-5">
            <div className="flex items-center gap-2">
              <MapPin className="h-5 w-5 text-emerald-600 dark:text-[#98FF98]" />
              <h2 className="text-lg font-semibold text-[#1A3C34] dark:text-white">
                First venue
              </h2>
            </div>

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
                placeholder="Main Street location"
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
                placeholder="Hours, neighborhood, or what makes this location special"
              />
            </div>

            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
              <div className="space-y-2">
                <label htmlFor="latitude" className={pageLabel}>
                  Latitude <span className="font-normal text-gray-500">(optional)</span>
                </label>
                <input
                  id="latitude"
                  type="number"
                  step="any"
                  min={-90}
                  max={90}
                  value={form.latitude}
                  onChange={(e) => setForm((f) => ({ ...f, latitude: e.target.value }))}
                  className={pageInput}
                  placeholder="37.7749"
                />
              </div>
              <div className="space-y-2">
                <label htmlFor="longitude" className={pageLabel}>
                  Longitude <span className="font-normal text-gray-500">(optional)</span>
                </label>
                <input
                  id="longitude"
                  type="number"
                  step="any"
                  min={-180}
                  max={180}
                  value={form.longitude}
                  onChange={(e) => setForm((f) => ({ ...f, longitude: e.target.value }))}
                  className={pageInput}
                  placeholder="-122.4194"
                />
              </div>
            </div>

            <div className="space-y-2">
              <label className={pageLabel}>
                Tags <span className="font-normal text-gray-500">(optional, up to 5)</span>
              </label>
              {tagsLoading ? (
                <p className="text-sm text-gray-500">Loading tags…</p>
              ) : locationTags.length === 0 ? (
                <p className="text-sm text-gray-500">No tags available yet.</p>
              ) : (
                <div className="flex flex-wrap gap-2">
                  {locationTags.map((tag) => {
                    const selected = form.tagIds.includes(tag.id);
                    return (
                      <button
                        key={tag.id}
                        type="button"
                        onClick={() => toggleTag(tag)}
                        className={cn(
                          "rounded-full border px-3 py-1.5 text-sm transition-colors",
                          selected
                            ? "border-[#FFD700] bg-[#FFD700] text-[#1A3C34]"
                            : "border-gray-200 bg-white text-gray-700 hover:bg-gray-50 dark:border-[#2f5e50] dark:bg-[#0d241d] dark:text-gray-200 dark:hover:bg-[#244f42]"
                        )}
                      >
                        {tag.name}
                      </button>
                    );
                  })}
                </div>
              )}
            </div>
          </div>

          <Button
            type="submit"
            className={cn(pageBtnPrimary, "w-full")}
            disabled={setupMutation.isPending}
          >
            {setupMutation.isPending
              ? "Saving…"
              : needsVenueOnly
                ? "Create first venue"
                : "Create business profile"}
          </Button>
        </form>
      </div>
    </PageLayout>
  );
}
