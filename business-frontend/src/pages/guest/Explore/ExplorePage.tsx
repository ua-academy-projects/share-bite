import { useState, type FormEvent } from "react";
import { Link } from "react-router-dom";
import { Loader2, MapPin, Search, Tag } from "lucide-react";

import { businessApi, type VenueSearchItem } from "@/api/business";
import { PageLayout } from "@/components/layout/PageLayout";
import { pageEmpty, pageLoader } from "@/components/layout/pageStyles";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { cn } from "@/lib/utils";

const PAGE_LIMIT = 10;

export function ExplorePage() {
  const [query, setQuery] = useState("");
  const [results, setResults] = useState<VenueSearchItem[]>([]);
  const [skip, setSkip] = useState(0);
  const [total, setTotal] = useState(0);
  const [loadingSearch, setLoadingSearch] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const totalPages = Math.max(1, Math.ceil(total / PAGE_LIMIT));
  const currentPage = Math.floor(skip / PAGE_LIMIT) + 1;
  const hasQuery = query.trim().length > 0;

  const runSearch = async (params?: { nextQuery?: string; nextSkip?: number }) => {
    const nextQuery = (params?.nextQuery ?? query).trim();
    const nextSkip = params?.nextSkip ?? 0;

    setError(null);

    if (!nextQuery) {
      setResults([]);
      setTotal(0);
      setSkip(0);
      setError("Enter a keyword to search.");
      return;
    }

    try {
      setLoadingSearch(true);
      const data = await businessApi.searchVenues({
        query: nextQuery,
        tags: [],
        skip: nextSkip,
        limit: PAGE_LIMIT,
      });
      setResults(data.items);
      setTotal(data.total);
      setSkip(nextSkip);
    } catch (e) {
      setError(e instanceof Error ? e.message : "Search failed");
    } finally {
      setLoadingSearch(false);
    }
  };

  const onSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    await runSearch({ nextSkip: 0 });
  };

  return (
    <PageLayout>
      <div className="space-y-8">
        <div>
          <h1 className="mb-3 text-4xl font-bold tracking-tight text-[#1A3C34] dark:text-white md:text-5xl">
            Explore <span className="text-emerald-500 dark:text-[#98FF98]">🌿</span>
          </h1>
          <p className="text-lg text-gray-600 dark:text-gray-400">
            Discover trending places near you.
          </p>
        </div>

        <Card className="rounded-3xl border border-gray-200 bg-white shadow-sm transition-colors duration-300 dark:border-[#2f5e50] dark:bg-[#163d32]">
          <CardContent className="p-6">
            <form onSubmit={onSubmit} className="flex gap-3">
              <Input
                value={query}
                onChange={(e) => setQuery(e.target.value)}
                placeholder="Search for amazing venues..."
                className="h-11 border-gray-200 bg-gray-50 dark:border-[#2f5e50] dark:bg-[#0d241d]"
              />
              <Button
                type="submit"
                disabled={!hasQuery || loadingSearch}
                className="h-11 bg-[#163d32] px-6 text-white hover:bg-[#1A3C34] dark:bg-green-500 dark:text-black dark:hover:bg-green-400"
              >
                {loadingSearch ? (
                  <Loader2 className="h-4 w-4 animate-spin" />
                ) : (
                  <Search className="h-4 w-4" />
                )}
              </Button>
            </form>
          </CardContent>
        </Card>

        {error ? (
          <div className="rounded-xl border border-red-500/30 bg-red-50 px-4 py-3 text-sm font-medium text-red-700 dark:bg-red-500/10 dark:text-red-400">
            {error}
          </div>
        ) : null}

        <div className="flex items-center justify-between">
          <p className="text-sm text-gray-600 dark:text-gray-300">
            Results: <span className="font-semibold">{total}</span>
          </p>
          <p className="text-sm text-gray-500 dark:text-gray-400">
            Page {currentPage} / {totalPages}
          </p>
        </div>

        {loadingSearch ? (
          <div className="flex h-44 items-center justify-center">
            <Loader2 className={cn(pageLoader, "h-10 w-10")} />
          </div>
        ) : results.length === 0 ? (
          <div className={pageEmpty}>
            <p className="text-xl font-bold text-[#1A3C34] dark:text-gray-200">No venues found</p>
            <p className="mt-2 text-gray-500">Try a different keyword.</p>
          </div>
        ) : (
          <div className="grid grid-cols-1 gap-6 md:grid-cols-2">
            {results.map((venue) => (
              <Card
                key={venue.id}
                className="rounded-3xl border border-gray-200 bg-white shadow-sm transition-all hover:shadow-lg dark:border-[#2f5e50] dark:bg-[#163d32]"
              >
                <CardContent className="space-y-4 p-5">
                  <div className="flex gap-4">
                    <img
                      src={venue.avatar || "https://placehold.co/96x96/163d32/FFF?text=SB"}
                      alt={venue.name}
                      className="h-16 w-16 rounded-2xl border border-gray-200 object-cover dark:border-[#2f5e50]"
                    />
                    <div>
                      <h3 className="text-xl font-bold text-[#1A3C34] dark:text-white">{venue.name}</h3>
                      <p className="mt-1 text-sm text-gray-600 dark:text-gray-300">
                        {venue.description || "No description"}
                      </p>
                    </div>
                  </div>

                  <div className="flex flex-wrap gap-2">
                    {venue.tags.map((tag) => (
                      <span
                        key={`${venue.id}-${tag}`}
                        className="inline-flex items-center gap-1 rounded-full border border-gray-200 bg-gray-100 px-2.5 py-1 text-xs text-gray-700 dark:border-[#2f5e50] dark:bg-[#0d241d] dark:text-gray-200"
                      >
                        <Tag className="h-3 w-3" />
                        {tag}
                      </span>
                    ))}
                  </div>

                  <div className="flex items-center justify-between">
                    <div className="inline-flex items-center gap-1 text-sm text-gray-500 dark:text-gray-300">
                      <MapPin className="h-4 w-4" />
                      Venue ID: {venue.id}
                    </div>
                    <Button
                      asChild
                      className="rounded-xl bg-[#163d32] px-5 font-semibold text-white ring-1 ring-[#FFD700]/40 hover:bg-[#1A3C34]"
                    >
                      <Link to={`/restaurant/${venue.id}`}>View Profile</Link>
                    </Button>
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
            disabled={skip === 0 || loadingSearch}
            onClick={() => void runSearch({ nextSkip: Math.max(0, skip - PAGE_LIMIT) })}
            className="rounded-full"
          >
            Previous
          </Button>
          <Button
            type="button"
            variant="outline"
            disabled={skip + PAGE_LIMIT >= total || loadingSearch}
            onClick={() => void runSearch({ nextSkip: skip + PAGE_LIMIT })}
            className="rounded-full"
          >
            Next
          </Button>
        </div>
      </div>
    </PageLayout>
  );
}
