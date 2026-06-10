import { useState } from "react";
import { Link } from "react-router-dom";
import {
  AlertCircle,
  Building2,
  FilePlus2,
  Loader2,
  MapPin,
  PackagePlus,
  Search,
  Tag,
} from "lucide-react";

import { businessApi, type VenueSearchItem } from "@/api/business";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";

const PAGE_LIMIT = 12;

export function MyVenuesPage() {
  const [venues, setVenues] = useState<VenueSearchItem[]>([]);
  const [skip, setSkip] = useState(0);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const currentPage = Math.floor(skip / PAGE_LIMIT) + 1;
  const totalPages = Math.max(1, Math.ceil(total / PAGE_LIMIT));

  const loadVenues = async (nextSkip = 0) => {
    setLoading(true);
    setError(null);

    try {
      const data = await businessApi.listCurrentBusinessVenues({
        skip: nextSkip,
        limit: PAGE_LIMIT,
      });
      setVenues(data.items);
      setTotal(data.total);
      setSkip(nextSkip);
    } catch (e) {
      setError(e instanceof Error ? e.message : "Failed to load business venues");
      setVenues([]);
      setTotal(0);
      setSkip(0);
    } finally {
      setLoading(false);
    }
  };

  const handlePage = (nextSkip: number) => {
    void loadVenues(nextSkip);
  };

  return (
    <div className="min-h-screen bg-[#F9F7F2] dark:bg-[#0d241d] p-8 md:p-12 w-full transition-colors duration-300">
      <div className="max-w-7xl mx-auto space-y-8">
        <div className="flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
          <div>
            <p className="text-sm font-semibold uppercase tracking-wide text-emerald-700 dark:text-[#98FF98]">
              Owner workspace
            </p>
            <h1 className="text-4xl md:text-5xl font-bold text-[#1A3C34] dark:text-white tracking-tight mt-2">
              My Venues
            </h1>
            <p className="text-gray-600 dark:text-gray-400 text-lg mt-3 max-w-2xl">
              Choose one of your business locations before creating posts or magic boxes.
            </p>
          </div>

          <div className="flex flex-col gap-3 sm:flex-row">
            <Button
              type="button"
              disabled={loading}
              onClick={() => void loadVenues(0)}
              className="h-11 bg-[#FFD700] text-[#1A3C34] hover:bg-[#e6c200] font-bold rounded-xl px-6"
            >
              {loading ? <Loader2 className="w-4 h-4 animate-spin" /> : "Load my venues"}
            </Button>
            <Button asChild className="h-11 bg-[#163d32] text-white hover:bg-[#1A3C34] rounded-xl px-5">
              <Link to="/venues/search">
                <Search className="w-4 h-4 mr-2" />
                Public search
              </Link>
            </Button>
          </div>
        </div>

        <Card className="bg-white dark:bg-[#163d32] border border-gray-200 dark:border-[#2f5e50] rounded-2xl shadow-sm">
          <CardContent className="p-5">
            <div className="flex items-center gap-3">
              <div className="w-11 h-11 rounded-xl bg-[#FFD700] text-[#1A3C34] flex items-center justify-center">
                <Building2 className="w-5 h-5" />
              </div>
              <div>
                <h2 className="text-xl font-bold text-[#1A3C34] dark:text-white">
                  Current business account
                </h2>
                <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
                  This workspace shows locations that belong to the signed-in business.
                </p>
              </div>
              <div className="ml-auto hidden sm:block">
              <Button
                type="button"
                disabled={loading}
                onClick={() => void loadVenues(0)}
                className="bg-[#163d32] text-white hover:bg-[#1A3C34] rounded-xl"
              >
                Refresh venues
              </Button>
              </div>
            </div>
          </CardContent>
        </Card>

        {error && (
          <div className="rounded-xl border border-red-500/30 bg-red-50 dark:bg-red-500/10 px-4 py-3 text-sm text-red-700 dark:text-red-400 font-medium flex items-start gap-2">
            <AlertCircle className="w-4 h-4 mt-0.5 flex-shrink-0" />
            <span>{error}</span>
          </div>
        )}

        <div className="flex items-center justify-between">
          <p className="text-gray-600 dark:text-gray-300 text-sm">
            Venues: <span className="font-semibold">{total}</span>
          </p>
          <p className="text-gray-500 dark:text-gray-400 text-sm">
            Page {currentPage} / {totalPages}
          </p>
        </div>

        {loading && venues.length === 0 ? (
          <div className="flex justify-center items-center h-56">
            <Loader2 className="w-10 h-10 text-emerald-500 dark:text-[#98FF98] animate-spin" />
          </div>
        ) : venues.length === 0 ? (
          <div className="text-center bg-white dark:bg-[#163d32] border border-gray-200 dark:border-[#2f5e50] rounded-3xl p-16 shadow-sm">
            <p className="text-[#1A3C34] dark:text-gray-200 text-xl font-bold">No owned venues loaded</p>
            <p className="text-gray-500 dark:text-gray-400 mt-2">
              Load your venues to choose where the next post or magic box should be created.
            </p>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
            {venues.map((venue) => (
              <Card
                key={venue.id}
                className="bg-white dark:bg-[#163d32] border border-gray-200 dark:border-[#2f5e50] rounded-3xl shadow-sm hover:shadow-lg transition-all"
              >
                <CardContent className="p-5 space-y-4">
                  <div className="flex gap-4">
                    <img
                      src={venue.avatar || "https://placehold.co/96x96/163d32/FFF?text=SB"}
                      alt={venue.name}
                      className="w-16 h-16 rounded-2xl object-cover border border-gray-200 dark:border-[#2f5e50]"
                    />
                    <div>
                      <h3 className="text-xl font-bold text-[#1A3C34] dark:text-white">{venue.name}</h3>
                      <p className="text-sm text-gray-600 dark:text-gray-300 mt-1">
                        {venue.description || "No description yet."}
                      </p>
                    </div>
                  </div>

                  {venue.tags.length > 0 && (
                    <div className="flex flex-wrap gap-2">
                      {venue.tags.map((tag) => (
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

                  <div className="inline-flex items-center gap-1 text-sm text-gray-500 dark:text-gray-300">
                    <MapPin className="w-4 h-4" />
                    Venue ID: {venue.id}
                  </div>

                  <div className="grid grid-cols-1 gap-2">
                    <Button asChild className="bg-[#163d32] text-white hover:bg-[#1A3C34] rounded-xl">
                      <Link to={`/venue/${venue.id}?manage=1`}>Manage venue</Link>
                    </Button>
                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
                      <Button asChild className="bg-[#FFD700] text-[#1A3C34] hover:bg-[#e6c200] font-bold rounded-xl">
                        <Link to={`/venue/${venue.id}/create-post`}>
                          <FilePlus2 className="w-4 h-4 mr-2" />
                          Post
                        </Link>
                      </Button>
                      <Button asChild className="bg-[#0d241d] text-white hover:bg-[#1A3C34] rounded-xl">
                        <Link to={`/venue/${venue.id}/create-box`}>
                          <PackagePlus className="w-4 h-4 mr-2" />
                          Box
                        </Link>
                      </Button>
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        )}

        <div className="flex justify-end gap-3">
          <Button
            type="button"
            variant="outline"
            disabled={skip === 0 || loading}
            onClick={() => handlePage(Math.max(0, skip - PAGE_LIMIT))}
            className="rounded-full"
          >
            Previous
          </Button>
          <Button
            type="button"
            variant="outline"
            disabled={skip + PAGE_LIMIT >= total || loading}
            onClick={() => handlePage(skip + PAGE_LIMIT)}
            className="rounded-full"
          >
            Next
          </Button>
        </div>
      </div>
    </div>
  );
}
