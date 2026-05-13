import { useEffect, useMemo, useState } from "react";
import { Link, useParams, useSearchParams, useNavigate } from "react-router-dom";
import { Loader2, MapPin, Tag, Plus } from "lucide-react";

import { businessApi, type VenueProfile } from "@/api/business";
import { Button } from "@/components/ui/button";
import { BrandTabs } from "@/features/brand/components/BrandTabs";
import { BrandSection } from "@/features/brand/components/BrandSection";
import { BrandPostCard } from "@/features/brand/components/BrandPostCard";
import { BoxCard } from "@/components/ui/BoxCard";
import { EmptyState } from "@/features/brand/components/EmptyState";
import { useBrandPosts, useBrandBoxes } from "@/features/brand/hooks";

const TAB_POSTS = "posts";
const TAB_BOXES = "boxes";

export function VenueProfilePage() {
  const { id } = useParams<{ id: string }>();
  const [searchParams, setSearchParams] = useSearchParams();
  const navigate = useNavigate();
  
  const venueId = Number(id);
  const isValidId = Number.isFinite(venueId) && venueId > 0;

  const [venue, setVenue] = useState<VenueProfile | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const activeTab = searchParams.get("tab") || TAB_POSTS;

  const setActiveTab = (tabId: string) => {
    setSearchParams({ tab: tabId }, { replace: true });
  };

  useEffect(() => {
    const loadVenue = async () => {
      if (!isValidId) {
        setError("Invalid venue id");
        setLoading(false);
        return;
      }

      try {
        setError(null);
        const data = await businessApi.getVenueProfile(venueId);
        setVenue(data);
      } catch (e) {
        setError(e instanceof Error ? e.message : "Failed to load venue");
      } finally {
        setLoading(false);
      }
    };

    void loadVenue();
  }, [venueId, isValidId]);

  const orgIds = useMemo(() => (isValidId ? [venueId] : []), [venueId, isValidId]);
  const postsState = useBrandPosts(orgIds, 9);
  
  const coordinates = useMemo(() => {
    if (venue && Number.isFinite(venue.latitude) && Number.isFinite(venue.longitude)) {
      return { lat: venue.latitude!, lon: venue.longitude! };
    }
    return null;
  }, [venue]);

  const boxesState = useBrandBoxes({ 
    orgId: isValidId ? venueId : undefined, 
    coordinates, 
    pageSize: 9 
  });

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
      count: boxesState.total,
      loading: boxesState.loading && boxesState.total === 0,
    },
  ];

  if (loading) {
    return (
      <div className="min-h-screen w-full p-8 md:p-12 bg-background flex items-center justify-center">
        <Loader2 className="w-10 h-10 text-[#98FF98] animate-spin" />
      </div>
    );
  }

  if (error || !venue) {
    return (
      <div className="min-h-screen w-full p-8 md:p-12 bg-background">
        <div className="max-w-4xl mx-auto bg-[#163d32]/50 border border-white/10 rounded-3xl p-8 backdrop-blur-md">
          <h1 className="text-2xl font-bold text-white mb-2">Venue not found</h1>
          <p className="text-red-400">{error || "Unknown error"}</p>
          <Button asChild className="mt-6 bg-white/10 text-white hover:bg-white/20 rounded-full">
            <Link to="/venue/search">Back to Search</Link>
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen w-full p-6 md:p-10 bg-background text-foreground">
      <div className="max-w-6xl mx-auto space-y-8">
        <div className="relative rounded-[28px] border border-white/10 bg-[#163d32]/40 backdrop-blur-xl overflow-hidden">
          <div className="h-44 md:h-64 relative">
            <img
              src={venue.banner || venue.avatar || "https://placehold.co/1200x400/163d32/FFF?text=ShareBite"}
              alt={venue.name}
              className="w-full h-full object-cover"
            />
            <div className="absolute inset-0 bg-gradient-to-t from-[#0b0f0e]/90 via-[#0b0f0e]/40 to-transparent" />
          </div>

          <div className="relative p-6 md:p-10">
            <div className="flex flex-col md:flex-row items-start gap-6">
              <img
                src={venue.avatar || "https://placehold.co/120x120/163d32/FFF?text=SB"}
                alt={venue.name}
                className="w-20 h-20 md:w-24 md:h-24 rounded-2xl object-cover border border-white/20 shadow-lg -mt-16 md:-mt-20"
              />
              <div className="flex-1">
                <div className="flex flex-wrap items-center gap-3">
                  <h1 className="text-2xl md:text-3xl font-bold text-white tracking-tight">{venue.name}</h1>
                  {venue.brand && (
                    <Link 
                      to={`/brand/${venue.brand.id}`}
                      className="text-xs bg-[#98FF98]/10 text-[#98FF98] border border-[#98FF98]/20 px-2 py-1 rounded-full hover:bg-[#98FF98]/20 transition"
                    >
                      Brand: {venue.brand.name}
                    </Link>
                  )}
                </div>
                <p className="text-[#cbd5cf] mt-2 max-w-2xl">
                  {venue.description || "No description yet."}
                </p>

                <div className="mt-4 flex flex-wrap items-center gap-4 text-sm text-[#9fb2a7]">
                  <span className="inline-flex items-center gap-1.5">
                    <MapPin className="w-4 h-4 text-[#98FF98]" />
                    {venue.latitude && venue.longitude ? (
                      `${venue.latitude.toFixed(4)}, ${venue.longitude.toFixed(4)}`
                    ) : "No coordinates set"}
                  </span>
                  <span className="px-2 py-0.5 rounded-md bg-white/5 border border-white/10 text-[10px] tracking-wider uppercase font-semibold">
                    ID: {venue.id}
                  </span>
                </div>
              </div>

              <div className="flex gap-2">
                 <Button asChild variant="outline" className="rounded-full border-white/15 hover:bg-white/10">
                   <Link to="/venue/search">Back to Search</Link>
                 </Button>
              </div>
            </div>

            {venue.tags?.length > 0 && (
              <div className="flex flex-wrap gap-2 mt-6">
                {venue.tags.map((tag) => (
                  <span
                    key={tag}
                    className="inline-flex items-center gap-1 px-3 py-1 rounded-full text-xs bg-white/5 text-[#cbd5cf] border border-white/10"
                  >
                    <Tag className="w-3 h-3" />
                    {tag}
                  </span>
                ))}
              </div>
            )}
          </div>
        </div>

        <BrandTabs tabs={tabs} activeTab={activeTab} onChange={setActiveTab} />

        {activeTab === TAB_POSTS && (
          <BrandSection 
            title="Venue Posts" 
            subtitle="Recent updates from this location"
            action={
              <Button 
                className="bg-[#98FF98] text-[#0b0f0e] hover:bg-[#98FF98]/90 rounded-full gap-2 font-semibold"
                onClick={() => navigate(`/venue/${venueId}/create-post`)}
              >
                <Plus className="h-4 w-4" /> Create Post
              </Button>
            }
          >
            {postsState.error ? <div className="text-red-400">{postsState.error}</div> : null}
            {postsState.loading && postsState.items.length === 0 ? (
              <div className="flex justify-center py-12">
                <Loader2 className="h-8 w-8 animate-spin text-[#98FF98]" />
              </div>
            ) : postsState.items.length === 0 ? (
              <EmptyState
                title="No posts yet"
                description="Share what's happening at this venue."
              />
            ) : (
              <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
                {postsState.items.map((post) => (
                  <BrandPostCard key={post.id} post={post} />
                ))}
              </div>
            )}
            
            {postsState.hasMore && postsState.items.length > 0 ? (
              <div className="flex justify-center mt-6">
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
        )}

        {activeTab === TAB_BOXES && (
          <BrandSection 
            title="Magic Boxes" 
            subtitle="Surplus food rescues available here"
            action={
              <Button 
                className="bg-[#FFD700] text-[#1A3C34] hover:bg-[#FFD700]/90 rounded-full gap-2 font-semibold"
                onClick={() => navigate(`/venue/${venueId}/create-box`)}
              >
                <Plus className="h-4 w-4" /> Create Box
              </Button>
            }
          >
            {boxesState.error ? <div className="text-red-400">{boxesState.error}</div> : null}
            {boxesState.missingCoordinates ? (
              <EmptyState
                title="Coordinates required"
                description="Please set coordinates for this venue to manage boxes."
              />
            ) : boxesState.loading && boxesState.items.length === 0 ? (
              <div className="flex justify-center py-12">
                <Loader2 className="h-8 w-8 animate-spin text-[#98FF98]" />
              </div>
            ) : boxesState.items.length === 0 ? (
              <EmptyState
                title="No boxes available"
                description="Add a magic box to rescue food at this location."
              />
            ) : (
              <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
                {boxesState.items.map((box) => (
                  <BoxCard key={box.id} box={box} />
                ))}
              </div>
            )}

            {boxesState.hasMore && boxesState.items.length > 0 ? (
              <div className="flex justify-center mt-6">
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
        )}
      </div>
    </div>
  );
}
