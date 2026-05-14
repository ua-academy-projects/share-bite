import { useMemo, useState } from "react";
import { useParams, useSearchParams } from "react-router-dom";
import { Loader2 } from "lucide-react";

import { BoxCard } from "@/components/ui/BoxCard";
import { Button } from "@/components/ui/button";
import { EmptyState } from "@/features/brand/components/EmptyState";
import { BrandHeader } from "@/features/brand/components/BrandHeader";
import { BrandLocationCard } from "@/features/brand/components/BrandLocationCard";
import { BrandPostCard } from "@/features/brand/components/BrandPostCard";
import { BrandSection } from "@/features/brand/components/BrandSection";
import { BrandTabs } from "@/features/brand/components/BrandTabs";
import { EditBrandProfileModal } from "@/features/brand/components/EditBrandProfileModal";
import {
  useBrandBoxes,
  useBrandLocations,
  useBrandPosts,
  useBrandProfile,
} from "@/features/brand/hooks";

const TAB_POSTS = "posts";
const TAB_BOXES = "boxes";
const TAB_LOCATIONS = "locations";

function getCoordinates(locations: Array<{ latitude?: number | null; longitude?: number | null }>) {
  const withCoords = locations.filter(
    (loc) => Number.isFinite(loc.latitude) && Number.isFinite(loc.longitude)
  );
  if (withCoords.length === 0) return null;

  const sum = withCoords.reduce(
    (acc, loc) => ({
      lat: acc.lat + (loc.latitude ?? 0),
      lon: acc.lon + (loc.longitude ?? 0),
    }),
    { lat: 0, lon: 0 }
  );

  return {
    lat: sum.lat / withCoords.length,
    lon: sum.lon / withCoords.length,
  };
}

function ErrorBanner({ message }: { message: string }) {
  return (
    <div className="rounded-3xl border border-red-500/30 bg-red-500/10 px-5 py-4 text-sm text-red-100">
      {message}
    </div>
  );
}

export function BrandProfilePage() {
  const { id } = useParams<{ id: string }>();
  const [searchParams, setSearchParams] = useSearchParams();
  
  const brandId = Number(id);
  const isValidId = Number.isFinite(brandId) && brandId > 0;

  const activeTab = searchParams.get("tab") || TAB_POSTS;

  const setActiveTab = (tabId: string) => {
    setSearchParams({ tab: tabId }, { replace: true });
  };

  const [isEditModalOpen, setIsEditModalOpen] = useState(false);

  const brandState = useBrandProfile(isValidId ? brandId : undefined);
  const locationsState = useBrandLocations(isValidId ? brandId : undefined, 50);

  const orgIds = useMemo(() => {
    if (!isValidId) return [];
    const ids = new Set<number>();
    ids.add(brandId);
    locationsState.items.forEach((loc) => ids.add(loc.id));
    return Array.from(ids);
  }, [brandId, isValidId, locationsState.items]);

  const postsState = useBrandPosts(orgIds, 9);
  const coordinates = useMemo(() => getCoordinates(locationsState.items), [locationsState.items]);
  const boxesState = useBrandBoxes({ orgId: isValidId ? brandId : undefined, coordinates, pageSize: 9 });

  const tabs = [
    {
      id: TAB_POSTS,
      label: "Posts",
      count: postsState.totalLoaded,
      loading: postsState.loading && postsState.totalLoaded === 0,
    },
    {
      id: TAB_BOXES,
      label: "Magic Boxes",
      count: locationsState.error || boxesState.missingCoordinates ? 0 : boxesState.total,
      loading: boxesState.loading && boxesState.total === 0 && !boxesState.missingCoordinates,
    },
    {
      id: TAB_LOCATIONS,
      label: "Locations",
      count: locationsState.total,
      loading: locationsState.loading && locationsState.total === 0,
    },
  ];

  if (!isValidId) {
    return (
      <div className="min-h-screen w-full px-6 py-10 md:px-10">
        <ErrorBanner message="Invalid brand id." />
      </div>
    );
  }

  return (
    <div className="min-h-screen w-full bg-background text-foreground">
      <div className="relative px-6 py-8 md:px-10 md:py-12">
        <div className="absolute inset-0 -z-10 bg-[radial-gradient(circle_at_top,_#1f4a3f_0%,_#0b0f0e_55%,_#0b0f0e_100%)]" />

        <div className="space-y-8">
          <BrandHeader
            brand={brandState.data}
            loading={brandState.loading}
            error={brandState.error}
            onEdit={() => setIsEditModalOpen(true)}
          />

          {brandState.data && (
            <EditBrandProfileModal
              brand={brandState.data}
              isOpen={isEditModalOpen}
              onOpenChange={setIsEditModalOpen}
              onSuccess={() => brandState.refresh()}
              onRefresh={() => brandState.refresh()}
            />
          )}

          <BrandTabs tabs={tabs} activeTab={activeTab} onChange={setActiveTab} />

          {activeTab === TAB_POSTS ? (
            <BrandSection title="Posts" subtitle="Latest updates from your brand">
              {postsState.error ? <ErrorBanner message={postsState.error} /> : null}
              {postsState.loading && postsState.items.length === 0 ? (
                <div className="flex justify-center py-12">
                  <Loader2 className="h-8 w-8 animate-spin text-[#98FF98]" />
                </div>
              ) : postsState.items.length === 0 ? (
                <EmptyState
                  title="No posts yet"
                  description="Start sharing updates to build your brand story."
                />
              ) : (
                <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
                  {postsState.items.map((post) => (
                    <BrandPostCard key={post.id} post={post} />
                  ))}
                </div>
              )}

              {postsState.hasMore && postsState.items.length > 0 ? (
                <div className="flex justify-center">
                  <Button
                    variant="outline"
                    onClick={postsState.loadMore}
                    disabled={postsState.loading}
                    className="rounded-full border-white/15 text-[#F9F7F2] hover:bg-white/10"
                  >
                    {postsState.loading ? "Loading..." : "Load more"}
                  </Button>
                </div>
              ) : null}
            </BrandSection>
          ) : null}

          {activeTab === TAB_BOXES ? (
            <BrandSection title="Magic Boxes" subtitle="Active boxes from your locations">
              {locationsState.error ? (
                <ErrorBanner message="Locations failed to load. Boxes cannot be resolved." />
              ) : boxesState.error ? (
                <ErrorBanner message={boxesState.error} />
              ) : null}
              {!locationsState.error && boxesState.missingCoordinates ? (
                <EmptyState
                  title="Locations need coordinates"
                  description="Add location coordinates to surface nearby boxes."
                />
              ) : boxesState.loading && boxesState.items.length === 0 ? (
                <div className="flex justify-center py-12">
                  <Loader2 className="h-8 w-8 animate-spin text-[#98FF98]" />
                </div>
              ) : boxesState.items.length === 0 ? (
                <EmptyState
                  title="No boxes available"
                  description="Create a magic box to start rescuing food."
                />
              ) : (
                <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
                  {boxesState.items.map((box) => (
                    <BoxCard key={box.id} box={box} />
                  ))}
                </div>
              )}

              {boxesState.hasMore && boxesState.items.length > 0 ? (
                <div className="flex justify-center">
                  <Button
                    variant="outline"
                    onClick={boxesState.loadMore}
                    disabled={boxesState.loading}
                    className="rounded-full border-white/15 text-[#F9F7F2] hover:bg-white/10"
                  >
                    {boxesState.loading ? "Loading..." : "Load more"}
                  </Button>
                </div>
              ) : null}
            </BrandSection>
          ) : null}

          {activeTab === TAB_LOCATIONS ? (
            <BrandSection title="Locations" subtitle="All venues attached to your brand">
              {locationsState.error ? <ErrorBanner message={locationsState.error} /> : null}
              {locationsState.loading && locationsState.items.length === 0 ? (
                <div className="flex justify-center py-12">
                  <Loader2 className="h-8 w-8 animate-spin text-[#98FF98]" />
                </div>
              ) : locationsState.items.length === 0 ? (
                <EmptyState
                  title="No locations yet"
                  description="Create a venue to connect your brand with guests."
                />
              ) : (
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  {locationsState.items.map((location) => (
                    <BrandLocationCard key={location.id} location={location} />
                  ))}
                </div>
              )}
            </BrandSection>
          ) : null}
        </div>
      </div>
    </div>
  );
}
