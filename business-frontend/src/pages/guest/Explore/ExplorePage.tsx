import { useState } from "react";
import { Link } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";
import { Loader2, MapPin, Search } from "lucide-react";
import { apiClient } from "@/api/client";
import type { ExploreVenueItem } from "@/types/api";
import { Input } from "@/components/ui/input";
import { Card, CardContent } from "@/components/ui/card";
import { PageHeader } from "@/components/layout/PageHeader";

export function ExplorePage() {
  const [searchQuery, setSearchQuery] = useState("");

  const { data: venues, isLoading } = useQuery({
    queryKey: ["explore", "nearby"],
    queryFn: () => apiClient.getExploreNearby(50.4501, 30.5234, 20),
  });

  const filteredVenues = (venues as ExploreVenueItem[] | undefined)?.filter((v) =>
    (v.name || "").toLowerCase().includes(searchQuery.toLowerCase())
  ) || [];

  return (
    <div className="px-6 py-8 lg:px-10">
      <PageHeader
        title="Explore"
        description="Discover trending places near you"
      />

      <div className="relative mx-auto mb-10 max-w-xl">
        <Search className="absolute left-4 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
        <Input
          placeholder="Search for amazing venues…"
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="h-12 rounded-full pl-11"
        />
      </div>

      {isLoading ? (
        <div className="flex justify-center py-16">
          <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
        </div>
      ) : filteredVenues.length === 0 ? (
        <Card className="mx-auto max-w-lg rounded-3xl bg-card-solid">
          <CardContent className="py-16 text-center text-muted-foreground">
            No venues found matching &quot;{searchQuery}&quot;
          </CardContent>
        </Card>
      ) : (
        <div className="mx-auto grid max-w-5xl grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
          {filteredVenues.map((venue) => (
              <Link key={venue.venue_id} to={`/restaurant/${venue.venue_id}`}>
                <Card className="group overflow-hidden rounded-3xl border-border bg-card-solid transition-all hover:-translate-y-1 hover:shadow-xl">
                  <div className="relative h-40 bg-muted">
                    {venue.avatar ? (
                      <img
                        src={venue.avatar}
                        alt=""
                        className="h-full w-full object-cover transition-transform group-hover:scale-105"
                      />
                    ) : (
                      <div className="flex h-full items-center justify-center text-primary">
                        <MapPin className="h-8 w-8" />
                      </div>
                    )}
                  </div>
                  <CardContent className="p-4">
                    <h3 className="truncate font-bold text-foreground group-hover:text-primary">
                      {venue.name || `Venue #${venue.venue_id}`}
                    </h3>
                    <p className="mt-2 flex items-center gap-1 text-xs font-semibold uppercase tracking-wider text-muted-foreground">
                      <MapPin className="h-3 w-3" /> Nearby
                    </p>
                  </CardContent>
                </Card>
              </Link>
            ))}
        </div>
      )}
    </div>
  );
}
