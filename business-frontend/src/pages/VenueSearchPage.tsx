import { useEffect, useMemo, useState, type FormEvent } from "react";
import { Link, useSearchParams } from "react-router-dom";
import { Loader2, MapPin, Search, Tag } from "lucide-react";

import { businessApi, type LocationTag, type VenueSearchItem } from "@/api/business";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { cn } from "@/lib/utils";

const PAGE_LIMIT = 10;

type IntentPreset = {
  id: string;
  label: string;
  query: string;
  tags: string[];
};

const INTENT_PRESETS: IntentPreset[] = [
  { id: "morning-coffee", label: "Morning Coffee", query: "coffee", tags: ["coffee", "breakfast"] },
  { id: "work-session", label: "Work Session", query: "workspace", tags: ["wifi", "work-friendly", "power-outlets"] },
  { id: "romantic-evening", label: "Romantic Evening", query: "dinner", tags: ["dinner", "romantic", "cozy"] },
  { id: "family-lunch", label: "Family Lunch", query: "lunch", tags: ["lunch", "family-friendly", "kid-friendly"] },
];

const normalizeIntent = (value: string) => value.trim().toLowerCase();

export function VenueSearchPage() {
  const [searchParams] = useSearchParams();
  const initialIntent = searchParams.get("intent") || "";
  const initialPreset = INTENT_PRESETS.find(
    (preset) => normalizeIntent(preset.label) === normalizeIntent(initialIntent),
  );

  const [query, setQuery] = useState(initialPreset?.query ?? initialIntent);
  const [allTags, setAllTags] = useState<LocationTag[]>([]);
  const [selectedTags, setSelectedTags] = useState<string[]>(initialPreset?.tags ?? []);
  const [activeIntentId, setActiveIntentId] = useState<string | null>(initialPreset?.id ?? null);

  const [results, setResults] = useState<VenueSearchItem[]>([]);
  const [skip, setSkip] = useState(0);
  const [total, setTotal] = useState(0);

  const [loadingTags, setLoadingTags] = useState(true);
  const [loadingSearch, setLoadingSearch] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const totalPages = Math.max(1, Math.ceil(total / PAGE_LIMIT));
  const currentPage = Math.floor(skip / PAGE_LIMIT) + 1;
  const hasFilter = query.trim().length > 0 || selectedTags.length > 0;

  useEffect(() => {
    const loadTags = async () => {
      try {
        const tags = await businessApi.getLocationTags();
        setAllTags(tags);
      } catch (e) {
        setError(e instanceof Error ? e.message : "Failed to load tags");
      } finally {
        setLoadingTags(false);
      }
    };
    void loadTags();
  }, []);

  const runSearch = async (params?: { nextQuery?: string; nextTags?: string[]; nextSkip?: number }) => {
    const nextQuery = (params?.nextQuery ?? query).trim();
    const nextTags = params?.nextTags ?? selectedTags;
    const nextSkip = params?.nextSkip ?? 0;

    setError(null);

    if (!nextQuery && nextTags.length === 0) {
      setResults([]);
      setTotal(0);
      setSkip(0);
      setError("Use at least one filter: query or tags.");
      return;
    }

    try {
      setLoadingSearch(true);
      const data = await businessApi.searchVenues({
        query: nextQuery,
        tags: nextTags,
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
    setActiveIntentId(null);
    await runSearch({ nextSkip: 0 });
  };

  const toggleTag = (slug: string) => {
    setSelectedTags((prev) =>
      prev.includes(slug) ? prev.filter((tag) => tag !== slug) : [...prev, slug]
    );
  };

  const applyIntent = async (preset: IntentPreset) => {
    setActiveIntentId(preset.id);
    setQuery(preset.query);
    setSelectedTags(preset.tags);
    await runSearch({ nextQuery: preset.query, nextTags: preset.tags, nextSkip: 0 });
  };

  const clearAll = () => {
    setActiveIntentId(null);
    setQuery("");
    setSelectedTags([]);
    setResults([]);
    setTotal(0);
    setSkip(0);
    setError(null);
  };

  const selectedTagSet = useMemo(() => new Set(selectedTags), [selectedTags]);

  return (
    <div className="min-h-screen bg-[#F9F7F2] dark:bg-[#0d241d] p-8 md:p-12 w-full transition-colors duration-300">
      <div className="max-w-7xl mx-auto space-y-8">
        <div>
          <h1 className="text-4xl md:text-5xl font-bold text-[#1A3C34] dark:text-white tracking-tight mb-3">
            Discover Venues <span className="text-emerald-500 dark:text-[#98FF98]">🔎</span>
          </h1>
          <p className="text-gray-600 dark:text-gray-400 text-lg">
            Search for your favorite venues and explore nearby spots.
          </p>
        </div>

        <Card className="bg-white dark:bg-[#163d32] border border-gray-200 dark:border-[#2f5e50] rounded-3xl shadow-sm transition-colors duration-300">
          <CardHeader className="pb-2">
            <CardTitle className="text-[#1A3C34] dark:text-white text-xl">Quick Intents</CardTitle>
          </CardHeader>
          <CardContent className="space-y-6">
            <div className="flex flex-wrap gap-3">
              {INTENT_PRESETS.map((preset) => (
                <Button
                  key={preset.id}
                  type="button"
                  onClick={() => void applyIntent(preset)}
                  className={cn(
                    "rounded-full px-5 bg-white dark:bg-[#0d241d] text-[#1A3C34] dark:text-white border border-gray-200 dark:border-[#2f5e50] hover:bg-[#f4efe5] dark:hover:bg-[#244f42]",
                    activeIntentId === preset.id && "bg-[#FFD700] text-[#1A3C34] border-[#FFD700]"
                  )}
                >
                  {preset.label}
                </Button>
              ))}
              <Button
                type="button"
                variant="ghost"
                className="rounded-full px-5 bg-[#2f5e50] text-white hover:bg-[#3a7564] border border-[#3a7564] font-medium"
                onClick={clearAll}
              >
                Clear
              </Button>
            </div>

            <form onSubmit={onSubmit} className="space-y-4">
              <div className="flex gap-3">
                <Input
                  value={query}
                  onChange={(e) => setQuery(e.target.value)}
                  placeholder="Search by keyword (name/description)"
                  className="h-11 bg-gray-50 dark:bg-[#0d241d] border-gray-200 dark:border-[#2f5e50]"
                />
                <Button
                  type="submit"
                  disabled={!hasFilter || loadingSearch}
                  className="h-11 px-6 bg-[#163d32] text-white hover:bg-[#1A3C34] dark:bg-green-500 dark:text-black dark:hover:bg-green-400"
                >
                  {loadingSearch ? <Loader2 className="w-4 h-4 animate-spin" /> : <Search className="w-4 h-4" />}
                </Button>
              </div>

              <div className="flex flex-wrap gap-2">
                {loadingTags ? (
                  <div className="text-sm text-gray-500 dark:text-gray-300">Loading tags...</div>
                ) : (
                  allTags.map((tag) => (
                    <button
                      key={tag.id}
                      type="button"
                      onClick={() => toggleTag(tag.slug)}
                      className={cn(
                        "px-3 py-1.5 rounded-full text-sm border transition-colors",
                        selectedTagSet.has(tag.slug)
                          ? "bg-[#FFD700] text-[#1A3C34] border-[#FFD700]"
                          : "bg-white dark:bg-[#0d241d] text-gray-700 dark:text-gray-200 border-gray-200 dark:border-[#2f5e50]"
                      )}
                    >
                      {tag.name}
                    </button>
                  ))
                )}
              </div>
            </form>
          </CardContent>
        </Card>

        {error && (
          <div className="rounded-xl border border-red-500/30 bg-red-50 dark:bg-red-500/10 px-4 py-3 text-sm text-red-700 dark:text-red-400 font-medium">
            {error}
          </div>
        )}

        <div className="flex items-center justify-between">
          <p className="text-gray-600 dark:text-gray-300 text-sm">
            Results: <span className="font-semibold">{total}</span>
          </p>
          <p className="text-gray-500 dark:text-gray-400 text-sm">
            Page {currentPage} / {totalPages}
          </p>
        </div>

        {loadingSearch ? (
          <div className="flex justify-center items-center h-44">
            <Loader2 className="w-10 h-10 text-emerald-500 dark:text-[#98FF98] animate-spin" />
          </div>
        ) : results.length === 0 ? (
          <div className="text-center bg-white dark:bg-[#163d32] border border-gray-200 dark:border-[#2f5e50] rounded-3xl p-16 shadow-sm">
            <p className="text-[#1A3C34] dark:text-gray-200 text-xl font-bold">No venues found</p>
            <p className="text-gray-500 mt-2">Try another intent, keyword, or tag combination.</p>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {results.map((venue) => (
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
                        {venue.description || "No description"}
                      </p>
                    </div>
                  </div>

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

                  <div className="flex items-center justify-between">
                    <div className="inline-flex items-center gap-1 text-sm text-gray-500 dark:text-gray-300">
                      <MapPin className="w-4 h-4" />
                      Venue ID: {venue.id}
                    </div>
                    <Button asChild className="bg-[#163d32] text-white hover:bg-[#1A3C34] font-semibold rounded-xl px-5 ring-1 ring-[#FFD700]/40">
                    <Link to={`/venue/${venue.id}`}>Visit Venue</Link>
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
    </div>
  );
}
