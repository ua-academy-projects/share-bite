import { useEffect, useState } from "react";
import { Link, useParams, useSearchParams } from "react-router-dom";
import { Building2, FilePlus2, Loader2, MapPin, PackagePlus, Star, Tag } from "lucide-react";

import { businessApi, type VenueProfile } from "@/api/business";
import { Button } from "@/components/ui/button";

export function VenueProfilePage() {
  const { id } = useParams<{ id: string }>();
  const [searchParams] = useSearchParams();
  const venueId = Number(id);
  const isManagementMode = searchParams.get("manage") === "1";

  const [venue, setVenue] = useState<VenueProfile | null>(null);
  const [rating, setRating] = useState<number | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const loadVenue = async () => {
      if (!Number.isFinite(venueId) || venueId <= 0) {
        setError("Invalid venue id");
        setLoading(false);
        return;
      }

      try {
        setError(null);
        const [profile, ratingData] = await Promise.all([
          businessApi.getVenueProfile(venueId),
          businessApi.getVenueRating(venueId).catch(() => null),
        ]);
        setVenue(profile);
        setRating(ratingData?.rating ?? null);
      } catch (e) {
        setError(e instanceof Error ? e.message : "Failed to load venue");
      } finally {
        setLoading(false);
      }
    };

    void loadVenue();
  }, [venueId]);

  if (loading) {
    return (
      <div className="min-h-screen w-full p-8 md:p-12 bg-[#F9F7F2] dark:bg-[#0d241d] flex items-center justify-center">
        <Loader2 className="w-10 h-10 text-emerald-500 dark:text-[#98FF98] animate-spin" />
      </div>
    );
  }

  if (error || !venue) {
    return (
      <div className="min-h-screen w-full p-8 md:p-12 bg-[#F9F7F2] dark:bg-[#0d241d]">
        <div className="max-w-4xl mx-auto bg-white dark:bg-[#163d32] border border-gray-200 dark:border-[#2f5e50] rounded-3xl p-8">
          <h1 className="text-2xl font-bold text-[#1A3C34] dark:text-white mb-2">Venue not found</h1>
          <p className="text-red-600 dark:text-red-400">{error || "Unknown error"}</p>
          <Button asChild className="mt-6 bg-[#163d32] text-white hover:bg-[#1A3C34]">
            <Link to="/venues/search">Back to Search</Link>
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen w-full p-8 md:p-12 bg-[#F9F7F2] dark:bg-[#0d241d] transition-colors duration-300">
      <div className="max-w-5xl mx-auto space-y-6">
        <div className="bg-white dark:bg-[#163d32] border border-gray-200 dark:border-[#2f5e50] rounded-3xl overflow-hidden">
          <div className="h-52 bg-gray-200 dark:bg-[#0d241d]">
            <img
              src={venue.banner || venue.avatar || "https://placehold.co/1200x400/163d32/FFF?text=ShareBite"}
              alt={venue.name}
              className="w-full h-full object-cover"
            />
          </div>

          <div className="p-6 md:p-8">
            <div className="flex flex-col gap-5 lg:flex-row lg:items-start lg:justify-between">
              <div className="flex items-start gap-4">
              <img
                src={venue.avatar || "https://placehold.co/120x120/163d32/FFF?text=SB"}
                alt={venue.name}
                className="w-16 h-16 rounded-2xl object-cover border border-gray-200 dark:border-[#2f5e50]"
              />
              <div className="flex-1">
                <h1 className="text-3xl font-bold text-[#1A3C34] dark:text-white">{venue.name}</h1>
                <p className="text-gray-600 dark:text-gray-300 mt-2">
                  {venue.description || "No description yet."}
                </p>

                {venue.brand && (
                  <p className="text-sm text-gray-500 dark:text-gray-400 mt-2">
                    Brand: {venue.brand.name}
                  </p>
                )}
              </div>
            </div>

              <div className="flex flex-col gap-3 min-w-full lg:min-w-[240px]">
                <div className="inline-flex items-center gap-2 rounded-xl border border-gray-200 dark:border-[#2f5e50] bg-gray-50 dark:bg-[#0d241d] px-4 py-3 text-[#1A3C34] dark:text-white">
                  <Star className="w-4 h-4 text-[#FFD700] fill-[#FFD700]" />
                  <span className="font-semibold">
                    {rating === null ? "Rating unavailable" : `${rating.toFixed(1)} rating`}
                  </span>
                </div>
                {isManagementMode ? (
                  <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-1 gap-2">
                    <Button asChild className="bg-[#FFD700] text-[#1A3C34] hover:bg-[#e6c200] font-bold rounded-xl">
                      <Link to={`/venue/${venue.id}/create-post`}>
                        <FilePlus2 className="w-4 h-4 mr-2" />
                        Create Post
                      </Link>
                    </Button>
                    <Button asChild className="bg-[#163d32] text-white hover:bg-[#1A3C34] rounded-xl">
                      <Link to={`/venue/${venue.id}/create-box`}>
                        <PackagePlus className="w-4 h-4 mr-2" />
                        Create Box
                      </Link>
                    </Button>
                  </div>
                ) : (
                  <div className="rounded-xl border border-gray-200 dark:border-[#2f5e50] bg-gray-50 dark:bg-[#0d241d] p-4 text-sm text-gray-600 dark:text-gray-300">
                    <p>Business actions are available from owned venues.</p>
                    <Button asChild className="mt-3 w-full bg-[#163d32] text-white hover:bg-[#1A3C34] rounded-xl">
                      <Link to="/venues/mine">
                        <Building2 className="w-4 h-4 mr-2" />
                        Open My Venues
                      </Link>
                    </Button>
                  </div>
                )}
              </div>
            </div>

            {venue.tags?.length > 0 && (
              <div className="flex flex-wrap gap-2 mt-5">
                {venue.tags.map((tag) => (
                  <span
                    key={tag}
                    className="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-xs bg-gray-100 dark:bg-[#0d241d] text-gray-700 dark:text-gray-200 border border-gray-200 dark:border-[#2f5e50]"
                  >
                    <Tag className="w-3 h-3" />
                    {tag}
                  </span>
                ))}
              </div>
            )}

            <div className="mt-6 flex flex-wrap items-center gap-4 text-sm text-gray-600 dark:text-gray-300">
              <span className="inline-flex items-center gap-1">
                <MapPin className="w-4 h-4" />
                Venue ID: {venue.id}
              </span>
              <span>
                Lat: {venue.latitude ?? "—"}, Lon: {venue.longitude ?? "—"}
              </span>
            </div>
          </div>
        </div>

        <div className="flex gap-3">
          <Button asChild className="bg-[#FFD700] text-[#1A3C34] hover:bg-[#e6c200] font-bold">
            <Link to={isManagementMode ? "/venues/mine" : "/venues/search"}>
              {isManagementMode ? "Back to My Venues" : "Back to Search"}
            </Link>
          </Button>
          <Button asChild className="bg-[#163d32] text-white hover:bg-[#1A3C34]">
            <Link to="/boxes">View Boxes</Link>
          </Button>
        </div>
      </div>
    </div>
  );
}
