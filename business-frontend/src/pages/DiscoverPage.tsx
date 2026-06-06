import { useCallback, useEffect, useState } from "react";
import { Link } from "react-router-dom";
import {
  AlertCircle,
  Compass,
  Flame,
  Loader2,
  MapPin,
  MapPinned,
  Package,
  RefreshCw,
  Search,
  Sparkles,
  Tag,
} from "lucide-react";

import { businessApi, type RecommendedVenue } from "@/api/business";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

const DISCOVER_LAT = 37.77351509723814;
const DISCOVER_LON = -122.4182710369247;
const PAGE_LIMIT = 12;

const discoveryCards = [
  {
    title: "Venue search",
    description: "Search by keyword, tags, or intent when you know what kind of place you need.",
    to: "/venues/search",
    icon: Search,
    cta: "Search venues",
  },
  {
    title: "Magic boxes",
    description: "Browse discounted food boxes and reserve available offers.",
    to: "/boxes",
    icon: Package,
    cta: "View boxes",
  },
  {
    title: "Recommended feed",
    description: "Return to the feed for recommended business posts and fresh updates.",
    to: "/",
    icon: Flame,
    cta: "Open feed",
  },
];

const intentLinks = [
  "Morning coffee",
  "Work-friendly spots",
  "Family lunch",
  "Dinner plans",
  "Groceries nearby",
  "Dessert rescue",
];

const formatDistance = (distance?: number | null) => {
  if (!Number.isFinite(distance)) return "Nearby";
  return `${Number(distance).toFixed(1)} km away`;
};

export function DiscoverPage() {
  const [venues, setVenues] = useState<RecommendedVenue[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [skip, setSkip] = useState(0);
  const [total, setTotal] = useState(0);

  const currentPage = Math.floor(skip / PAGE_LIMIT) + 1;
  const totalPages = Math.max(1, Math.ceil(total / PAGE_LIMIT));

  const loadVenues = useCallback(async (skipValue = 0) => {
    setLoading(true);
    setError(null);

    try {
      const data = await businessApi.recommendVenues({
        lat: DISCOVER_LAT,
        lon: DISCOVER_LON,
        skip: skipValue,
        limit: PAGE_LIMIT,
      });

      setVenues(data.items);
      setTotal(data.total);
      setSkip(skipValue);
    } catch (e) {
      setError(e instanceof Error ? e.message : "Failed to load venues");
      setVenues([]);
      setTotal(0);
      setSkip(0);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    const timer = window.setTimeout(() => {
      void loadVenues(0);
    }, 0);

    return () => window.clearTimeout(timer);
  }, [loadVenues]);

  return (
    <div className="min-h-screen bg-[#F9F7F2] dark:bg-[#0d241d] p-8 md:p-12 w-full transition-colors duration-300">
      <div className="max-w-7xl mx-auto space-y-8">
        <section className="grid grid-cols-1 xl:grid-cols-[1.1fr_0.9fr] gap-6 items-stretch">
          <div className="rounded-3xl bg-[#163d32] text-white p-8 md:p-10 shadow-xl overflow-hidden relative">
            <div className="relative z-10 max-w-2xl">
              <div className="inline-flex items-center gap-2 rounded-full bg-white/10 border border-white/15 px-4 py-2 text-sm text-[#98FF98]">
                <Compass className="w-4 h-4" />
                Discover workspace
              </div>
              <h1 className="text-4xl md:text-5xl font-bold tracking-tight mt-6">
                Discover recommended venues nearby
              </h1>
              <p className="text-gray-200 text-lg mt-4 leading-relaxed">
                Browse venues from the business API, then jump into focused search, boxes,
                or venue profiles when you want to take action.
              </p>
              <div className="flex flex-col sm:flex-row gap-3 mt-8">
                <Button asChild className="bg-[#FFD700] text-[#1A3C34] hover:bg-[#e6c200] rounded-xl font-bold">
                  <Link to="/venues/search">
                    <Search className="w-4 h-4 mr-2" />
                    Search venues
                  </Link>
                </Button>
                <Button asChild className="bg-white/10 text-white hover:bg-white/15 border border-white/15 rounded-xl">
                  <Link to="/boxes">
                    <Package className="w-4 h-4 mr-2" />
                    Browse boxes
                  </Link>
                </Button>
              </div>
            </div>
          </div>

          <Card className="bg-white dark:bg-[#163d32] border border-gray-200 dark:border-[#2f5e50] rounded-3xl shadow-sm">
            <CardHeader>
              <div className="w-12 h-12 rounded-2xl bg-[#FFD700] text-[#1A3C34] flex items-center justify-center mb-2">
                <Sparkles className="w-6 h-6" />
              </div>
              <CardTitle className="text-2xl font-bold text-[#1A3C34] dark:text-white">
                Quick intents
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="flex flex-wrap gap-2">
                {intentLinks.map((intent) => (
                  <Link
                    key={intent}
                    to={`/venues/search?intent=${encodeURIComponent(intent)}`}
                    className="rounded-full border border-gray-200 dark:border-[#2f5e50] bg-gray-50 dark:bg-[#0d241d] px-3 py-2 text-sm text-gray-700 dark:text-gray-200 hover:border-[#FFD700] hover:text-[#1A3C34] dark:hover:text-white transition-colors"
                  >
                    {intent}
                  </Link>
                ))}
              </div>
            </CardContent>
          </Card>
        </section>

        <section className="space-y-5">
          <div className="flex flex-col gap-4 md:flex-row md:items-end md:justify-between">
            <div>
              <p className="text-sm font-semibold uppercase tracking-wide text-emerald-700 dark:text-[#98FF98]">
                From business-api
              </p>
              <h2 className="text-3xl font-bold text-[#1A3C34] dark:text-white mt-1">
                Recommended venues
              </h2>
              <p className="text-gray-600 dark:text-gray-300 mt-2">
                Showing venues near San Francisco coordinates used by the current demo data.
              </p>
            </div>
            <Button
              type="button"
              onClick={() => void loadVenues(0)}
              disabled={loading}
              className="bg-[#FFD700] text-[#1A3C34] hover:bg-[#e6c200] rounded-xl font-bold"
            >
              <RefreshCw className={`w-4 h-4 mr-2 ${loading ? "animate-spin" : ""}`} />
              Refresh
            </Button>
          </div>

          {error && (
            <div className="rounded-2xl border border-red-500/30 bg-red-50 dark:bg-red-500/10 px-5 py-4 flex items-start gap-3">
              <AlertCircle className="w-5 h-5 text-red-600 dark:text-red-400 mt-0.5 flex-shrink-0" />
              <div>
                <p className="font-bold text-red-700 dark:text-red-300">Could not load venues</p>
                <p className="text-sm text-red-600 dark:text-red-400 mt-1">{error}</p>
              </div>
            </div>
          )}

          {loading && venues.length === 0 ? (
            <div className="flex h-64 items-center justify-center rounded-3xl border border-gray-200 dark:border-[#2f5e50] bg-white dark:bg-[#163d32]">
              <div className="flex flex-col items-center gap-3 text-gray-600 dark:text-gray-300">
                <Loader2 className="w-10 h-10 animate-spin text-emerald-600 dark:text-[#98FF98]" />
                <span className="font-medium">Loading venues...</span>
              </div>
            </div>
          ) : venues.length === 0 ? (
            <div className="rounded-3xl border border-gray-200 dark:border-[#2f5e50] bg-white dark:bg-[#163d32] p-12 text-center">
              <p className="text-xl font-bold text-[#1A3C34] dark:text-white">No venues found</p>
              <p className="text-gray-500 dark:text-gray-300 mt-2">
                Try again after business-api has venue data for the demo location.
              </p>
            </div>
          ) : (
            <>
              <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-5">
                {venues.map((venue) => (
                  <Card
                    key={venue.id}
                    className="bg-white dark:bg-[#163d32] border border-gray-200 dark:border-[#2f5e50] rounded-2xl shadow-sm hover:shadow-lg transition-all"
                  >
                    <CardContent className="p-5 space-y-4">
                      <div className="flex gap-4">
                        <img
                          src={venue.avatar || "https://placehold.co/96x96/163d32/FFF?text=SB"}
                          alt={venue.name}
                          className="w-16 h-16 rounded-2xl object-cover border border-gray-200 dark:border-[#2f5e50]"
                        />
                        <div className="min-w-0">
                          <h3 className="text-xl font-bold text-[#1A3C34] dark:text-white truncate">
                            {venue.name}
                          </h3>
                          <p className="inline-flex items-center gap-1 text-sm text-gray-500 dark:text-gray-300 mt-1">
                            <MapPin className="w-4 h-4" />
                            {formatDistance(venue.distance)}
                          </p>
                        </div>
                      </div>

                      <p className="text-sm text-gray-600 dark:text-gray-300 min-h-10">
                        {venue.description || "Open the venue profile to view details and business actions."}
                      </p>

                      {venue.tags && venue.tags.length > 0 && (
                        <div className="flex flex-wrap gap-2">
                          {venue.tags.slice(0, 4).map((tag) => (
                            <span
                              key={`${venue.id}-${tag}`}
                              className="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-xs bg-gray-100 dark:bg-[#0d241d] text-gray-700 dark:text-gray-200 border border-gray-200 dark:border-[#2f5e50]"
                            >
                              <Tag className="w-3 h-3" />
                              {tag}
                            </span>
                          ))}
                        </div>
                      )}

                      <Button asChild className="w-full bg-[#163d32] text-white hover:bg-[#1A3C34] rounded-xl">
                        <Link to={`/venue/${venue.id}`}>Open venue</Link>
                      </Button>
                    </CardContent>
                  </Card>
                ))}
              </div>

              <div className="flex flex-col gap-3 border-t border-gray-200 dark:border-[#2f5e50] pt-5 sm:flex-row sm:items-center sm:justify-between">
                <p className="text-sm text-gray-600 dark:text-gray-300">
                  Page {currentPage} of {totalPages} - {total} venues
                </p>
                <div className="flex gap-3">
                  <Button
                    type="button"
                    variant="outline"
                    disabled={skip === 0 || loading}
                    onClick={() => void loadVenues(Math.max(0, skip - PAGE_LIMIT))}
                    className="rounded-xl"
                  >
                    Previous
                  </Button>
                  <Button
                    type="button"
                    variant="outline"
                    disabled={skip + PAGE_LIMIT >= total || loading}
                    onClick={() => void loadVenues(skip + PAGE_LIMIT)}
                    className="rounded-xl"
                  >
                    Next
                  </Button>
                </div>
              </div>
            </>
          )}
        </section>

        <section className="grid grid-cols-1 md:grid-cols-3 gap-5">
          {discoveryCards.map((card) => {
            const Icon = card.icon;
            return (
              <Card
                key={card.title}
                className="bg-white dark:bg-[#163d32] border border-gray-200 dark:border-[#2f5e50] rounded-2xl shadow-sm"
              >
                <CardHeader>
                  <div className="w-11 h-11 rounded-xl bg-gray-100 dark:bg-[#0d241d] text-emerald-700 dark:text-[#98FF98] flex items-center justify-center mb-2">
                    <Icon className="w-5 h-5" />
                  </div>
                  <CardTitle className="text-xl font-bold text-[#1A3C34] dark:text-white">
                    {card.title}
                  </CardTitle>
                </CardHeader>
                <CardContent className="space-y-5">
                  <p className="text-gray-600 dark:text-gray-300 leading-relaxed">
                    {card.description}
                  </p>
                  <Button asChild className="w-full bg-[#163d32] text-white hover:bg-[#1A3C34] rounded-xl">
                    <Link to={card.to}>{card.cta}</Link>
                  </Button>
                </CardContent>
              </Card>
            );
          })}
        </section>

        <section className="rounded-2xl border border-emerald-200 dark:border-[#2f5e50] bg-white dark:bg-[#163d32] p-5 flex flex-col md:flex-row gap-4 md:items-center md:justify-between">
          <div className="flex items-start gap-3">
            <MapPinned className="w-5 h-5 text-emerald-700 dark:text-[#98FF98] mt-1" />
            <div>
              <h2 className="font-bold text-[#1A3C34] dark:text-white">Focused search lives separately</h2>
              <p className="text-sm text-gray-600 dark:text-gray-300 mt-1">
                Keep using Venue Search for filters and exact lookups. Discover is now the broader entry point.
              </p>
            </div>
          </div>
          <Button asChild variant="outline" className="rounded-xl">
            <Link to="/venues/search">Open Venue Search</Link>
          </Button>
        </section>
      </div>
    </div>
  );
}
