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
      <div className="max-w-7xl mx-auto space-y-12">
        {/* Immersive Header */}
        <div className="relative rounded-[48px] border border-white/10 bg-[#163d32]/20 backdrop-blur-2xl overflow-hidden shadow-2xl">
          {/* Banner */}
          <div className="h-64 md:h-[400px] relative overflow-hidden group">
            <img
              src={venue.banner || venue.avatar || "https://placehold.co/1200x400/163d32/FFF?text=ShareBite"}
              alt={venue.name}
              className="w-full h-full object-cover transition-transform duration-1000 group-hover:scale-105"
            />
            <div className="absolute inset-0 bg-gradient-to-t from-[#0b0f0e] via-[#0b0f0e]/30 to-transparent" />
            <div className="absolute inset-0 bg-black/5" />
          </div>

          <div className="relative p-8 md:p-16 pt-0">
            <div className="flex flex-col md:flex-row items-end gap-10">
              {/* Prominent Avatar */}
              <div className="-mt-20 md:-mt-32 relative group">
                <div className="relative h-32 w-32 md:h-48 md:w-48 rounded-[40px] p-1.5 bg-gradient-to-br from-white/20 to-white/5 backdrop-blur-xl shadow-2xl transition-transform duration-500 group-hover:scale-[1.02]">
                  <img
                    src={venue.avatar || "https://placehold.co/120x120/163d32/FFF?text=SB"}
                    alt={venue.name}
                    className="w-full h-full rounded-[34px] object-cover border border-white/20 shadow-inner"
                  />
                </div>
              </div>

              <div className="flex-1 space-y-4 pb-2">
                <div className="flex flex-wrap items-center justify-between gap-6">
                  <div className="space-y-1">
                    <div className="flex items-center gap-4">
                      <h1 className="text-4xl md:text-6xl font-black text-white tracking-tighter">
                        {venue.name}
                      </h1>
                      {venue.brand && (
                        <Link 
                          to={`/brand/${venue.brand.id}`}
                          className="mt-2 text-xs bg-[#98FF98]/10 text-[#98FF98] border border-[#98FF98]/20 px-3 py-1.5 rounded-full hover:bg-[#98FF98]/20 transition-all backdrop-blur-md font-bold uppercase tracking-wider"
                        >
                          Part of {venue.brand.name}
                        </Link>
                      )}
                    </div>
                    <div className="flex items-center gap-4 text-sm font-bold uppercase tracking-[0.2em] text-[#98FF98]/80">
                      <span className="flex items-center gap-1.5">
                        <MapPin className="h-4 w-4" />
                        {venue.latitude && venue.longitude ? (
                          `${venue.latitude.toFixed(4)}, ${venue.longitude.toFixed(4)}`
                        ) : "Location Pending"}
                      </span>
                      <span className="h-1 w-1 rounded-full bg-white/20" />
                      <span>ID: {venue.id}</span>
                    </div>
                  </div>

                  <div className="flex gap-3">
                     <Button asChild variant="outline" className="rounded-2xl border-white/10 bg-white/5 px-6 py-6 text-sm font-bold text-white backdrop-blur-md transition-all hover:bg-white/10 hover:border-white/20">
                       <Link to="/venue/search">Back to Search</Link>
                     </Button>
                  </div>
                </div>

                {venue.description && (
                  <p className="max-w-3xl text-lg md:text-xl font-medium leading-relaxed text-[#cbd5cf]">
                    {venue.description}
                  </p>
                )}

                {venue.tags?.length > 0 && (
                  <div className="flex flex-wrap gap-3 pt-2">
                    {venue.tags.map((tag) => (
                      <span
                        key={tag}
                        className="inline-flex items-center gap-2 px-4 py-1.5 rounded-full text-[11px] font-black uppercase tracking-widest bg-white/5 text-[#cbd5cf] border border-white/10 backdrop-blur-md"
                      >
                        <Tag className="w-3 h-3 text-[#98FF98]" />
                        {tag}
                      </span>
                    ))}
                  </div>
                )}
              </div>
            </div>
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
