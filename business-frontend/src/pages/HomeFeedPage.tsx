import { useEffect, useState } from "react";
import { businessApi, RecommendedPost } from "@/api/business";
import { PostCard, PostData } from "@/components/ui/PostCard";
import { Loader2, AlertCircle, RefreshCw } from "lucide-react";
import { Button } from "@/components/ui/button";

const HARDCODED_LAT = 37.77351509723814;
const HARDCODED_LON = -122.4182710369247;

const PAGE_LIMIT = 10;

const mapRecommendedPostToPostData = (post: RecommendedPost): PostData => {
  return {
    id: post.id,
    content: post.content,
    created_at: post.created_at,
    org: {
      id: post.org_id,
      name: `Business ${post.org_id}`,
      profileType: post.post_type,
    },
    images: [],
  };
};

export function HomeFeedPage() {
  const [posts, setPosts] = useState<PostData[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [skip, setSkip] = useState(0);
  const [total, setTotal] = useState(0);

  const currentPage = Math.floor(skip / PAGE_LIMIT) + 1;
  const totalPages = Math.max(1, Math.ceil(total / PAGE_LIMIT));

  const loadPosts = async (skipValue = 0) => {
    setLoading(true);
    setError(null);
    try {
      const token = localStorage.getItem("token") || undefined;
      const data = await businessApi.recommendPosts(
        {
          lat: HARDCODED_LAT,
          lon: HARDCODED_LON,
          skip: skipValue,
          limit: PAGE_LIMIT,
        },
        token,
      );
      // Map API response to PostData format
      const mappedPosts = data.items.map(mapRecommendedPostToPostData);
      setPosts(mappedPosts);
      setTotal(data.total);
      setSkip(skipValue);
    } catch (e) {
      const errorMessage =
        e instanceof Error ? e.message : "Failed to load recommendations";
      setError(errorMessage);
      console.error("Error loading posts:", e);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    void loadPosts(0);
  }, []);

  const handleNextPage = () => {
    const nextSkip = skip + PAGE_LIMIT;
    if (nextSkip < total) {
      void loadPosts(nextSkip);
      window.scrollTo({ top: 0, behavior: "smooth" });
    }
  };

  const handlePrevPage = () => {
    const prevSkip = Math.max(0, skip - PAGE_LIMIT);
    void loadPosts(prevSkip);
    window.scrollTo({ top: 0, behavior: "smooth" });
  };

  const handleRefresh = () => {
    void loadPosts(0);
  };

  return (
    <div className="min-h-screen bg-[#F9F7F2] dark:bg-[#0d241d] p-8 md:p-12 w-full transition-colors duration-300">
      <div className="max-w-7xl mx-auto space-y-8">
        {/* Header */}
        <div className="flex items-start justify-between">
          <div>
            <h1 className="text-4xl md:text-5xl font-bold text-[#1A3C34] dark:text-white tracking-tight mb-3">
              Home Feed{" "}
              <span className="text-emerald-500 dark:text-[#98FF98]">🔥</span>
            </h1>
            <p className="text-gray-600 dark:text-gray-400 text-lg">
              Discover posts recommended for you based on your preferences
            </p>
            <p className="text-gray-500 dark:text-gray-500 text-sm mt-2">
              📍 Location: San Francisco (37.77°N, 122.42°W)
            </p>
          </div>
          <Button
            onClick={handleRefresh}
            disabled={loading}
            className="bg-[#FFD700] text-[#1A3C34] hover:bg-[#e6c200] dark:hover:bg-[#FFD700]/80 font-semibold rounded-xl px-6 py-6 shadow-md dark:shadow-lg dark:shadow-[#FFD700]/20 transition-all disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
          >
            <RefreshCw size={18} className={loading ? "animate-spin" : ""} />
            Refresh
          </Button>
        </div>

        {/* Error State */}
        {error && (
          <div className="rounded-2xl border border-red-500/30 bg-red-50 dark:bg-red-500/10 px-6 py-4 flex items-start gap-3 shadow-sm">
            <AlertCircle
              className="text-red-600 dark:text-red-400 flex-shrink-0 mt-0.5"
              size={20}
            />
            <div>
              <p className="font-bold text-red-700 dark:text-red-300">
                Error loading posts
              </p>
              <p className="text-red-600 dark:text-red-400 text-sm mt-1">
                {error}
              </p>
            </div>
          </div>
        )}

        {/* Loading State */}
        {loading && posts.length === 0 ? (
          <div className="flex justify-center items-center h-96 w-full">
            <div className="flex flex-col items-center gap-4">
              <Loader2 className="w-16 h-16 text-emerald-500 dark:text-[#98FF98] animate-spin" />
              <p className="text-gray-600 dark:text-gray-400 font-medium">
                Loading recommendations...
              </p>
            </div>
          </div>
        ) : posts.length > 0 ? (
          <>
            {/* Posts Feed */}
            <div className="flex flex-col items-center w-full gap-6">
              <div className="w-full max-w-2xl flex flex-col gap-6">
                {posts.map((post) => (
                  <PostCard key={post.id} post={post} />
                ))}
              </div>
            </div>

            {/* Pagination */}
            <div className="flex justify-between items-center py-12 border-t border-gray-200 dark:border-[#2f5e50]">
              <Button
                onClick={handlePrevPage}
                disabled={skip === 0}
                className="bg-white dark:bg-[#163d32] border border-gray-200 dark:border-[#2f5e50] text-[#1A3C34] dark:text-white hover:bg-gray-50 dark:hover:bg-[#244f42] font-semibold rounded-xl px-6 py-3 transition-all disabled:opacity-50 disabled:cursor-not-allowed"
              >
                ← Previous
              </Button>

              <div className="flex flex-col items-center gap-2">
                <p className="text-gray-600 dark:text-gray-400 font-medium">
                  Page {currentPage} of {totalPages}
                </p>
                <p className="text-gray-500 dark:text-gray-500 text-sm">
                  Showing {skip + 1}–{Math.min(skip + PAGE_LIMIT, total)} of{" "}
                  {total} posts
                </p>
              </div>

              <Button
                onClick={handleNextPage}
                disabled={skip + PAGE_LIMIT >= total}
                className="bg-[#FFD700] text-[#1A3C34] hover:bg-[#e6c200] dark:hover:bg-[#FFD700]/80 font-semibold rounded-xl px-6 py-3 shadow-md dark:shadow-lg dark:shadow-[#FFD700]/20 transition-all disabled:opacity-50 disabled:cursor-not-allowed"
              >
                Next →
              </Button>
            </div>
          </>
        ) : (
          <div className="text-center bg-white dark:bg-[#163d32] border border-gray-200 dark:border-[#2f5e50] rounded-3xl p-20 shadow-sm dark:shadow-none transition-colors duration-300">
            <p className="text-[#1A3C34] dark:text-gray-300 text-2xl font-bold">
              No recommendations found 😢
            </p>
            <p className="text-gray-500 dark:text-gray-400 mt-3">
              Try checking back later or enable notifications to get notified
              about new posts.
            </p>
          </div>
        )}
      </div>
    </div>
  );
}
