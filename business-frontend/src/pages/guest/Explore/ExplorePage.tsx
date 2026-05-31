import { useMemo, useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { Link } from "react-router-dom";
import { Loader2, MapPin, Search } from "lucide-react";
import { apiClient } from "@/api/client";
import { PageHeader } from "@/components/layout/PageHeader";
import { PageLayout } from "@/components/layout/PageLayout";
import {
  pageEmpty,
  pageFilterBar,
  pageInput,
  pageLoader,
  pagePanel,
} from "@/components/layout/pageStyles";
import { cn } from "@/lib/utils";

type ExploreVenue = {
  venue_id: number;
  name?: string;
  avatar?: string;
};

export function ExplorePage() {
  const [searchQuery, setSearchQuery] = useState("");

  const { data: venues, isLoading } = useQuery({
    queryKey: ["explore", "nearby"],
    queryFn: () => apiClient.getExploreNearby(50.4501, 30.5234, 20),
  });

  const filteredVenues = useMemo(() => {
    const list = (venues || []) as ExploreVenue[];
    return list.filter((venue) =>
      (venue.name || "").toLowerCase().includes(searchQuery.toLowerCase())
    );
  }, [searchQuery, venues]);

  return (
    <PageLayout>
      <PageHeader
        title="Explore"
        description="Discover trending places near you"
      />

      <div className={cn(pageFilterBar, "mb-8 flex items-center gap-3")}>
        <Search className="h-5 w-5 shrink-0 text-emerald-500 dark:text-[#98FF98]" />
        <input
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          placeholder="Search for amazing venues..."
          className={cn(pageInput, "border-0 bg-transparent px-0 focus:ring-0")}
        />
      </div>

      {isLoading ? (
        <div className="flex h-64 items-center justify-center">
          <Loader2 className={cn(pageLoader, "h-12 w-12")} />
        </div>
      ) : filteredVenues.length === 0 ? (
        <div className={pageEmpty}>
          <p className="text-xl font-bold text-[#1A3C34] dark:text-gray-300">
            No venues found
          </p>
          <p className="mt-2 text-gray-500">
            {searchQuery
              ? `Nothing matches "${searchQuery}". Try a different search.`
              : "No venues nearby right now."}
          </p>
        </div>
      ) : (
        <div className="grid grid-cols-1 gap-8 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
          {filteredVenues.map((venue) => (
            <Link key={venue.venue_id} to={`/restaurant/${venue.venue_id}`}>
              <div
                className={cn(
                  pagePanel,
                  "group h-full overflow-hidden transition-transform hover:-translate-y-1"
                )}
              >
                <div className="relative h-40 overflow-hidden bg-gray-100 dark:bg-[#0d241d]">
                  {venue.avatar ? (
                    <img
                      src={venue.avatar}
                      alt={venue.name || "Venue"}
                      className="h-full w-full object-cover transition-transform group-hover:scale-105"
                    />
                  ) : (
                    <div className="flex h-full w-full items-center justify-center text-gray-400">
                      <MapPin className="h-8 w-8" />
                    </div>
                  )}
                </div>
                <div className="flex items-end justify-between gap-3 p-5">
                  <h3 className="line-clamp-1 text-lg font-semibold text-[#1A3C34] dark:text-white">
                    {venue.name || `Venue #${venue.venue_id}`}
                  </h3>
                  <span className="inline-flex items-center gap-1 text-xs text-gray-500 dark:text-gray-400">
                    <MapPin className="h-3 w-3" />
                    Nearby
                  </span>
                </div>
              </div>
            </Link>
          ))}
        </div>
      )}
    </PageLayout>
  );
}
