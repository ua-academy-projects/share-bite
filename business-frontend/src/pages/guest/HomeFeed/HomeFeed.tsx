import { useQuery } from "@tanstack/react-query";
import { Loader2, RefreshCw } from "lucide-react";
import { apiClient } from "@/api/client";
import { GuestPostCard } from "@/components/PostCard/PostCard";
import { PageLayout } from "@/components/layout/PageLayout";
import { pageEmpty, pageLoader } from "@/components/layout/pageStyles";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

export function HomeFeed() {
  const { data, isLoading, isFetching, error, refetch } = useQuery({
    queryKey: ["posts"],
    queryFn: () => apiClient.getPosts(20, 0),
  });

  return (
    <PageLayout className="space-y-8">
      <div className="flex w-full items-start justify-between gap-4">
        <div>
          <h1 className="mb-3 text-4xl font-bold tracking-tight text-[#1A3C34] dark:text-white md:text-5xl">
            Users Feed{" "}
            <span className="text-emerald-500 dark:text-[#98FF98]">👥</span>
          </h1>
          <p className="text-lg text-gray-600 dark:text-gray-400">
            Discover what your friends are eating
          </p>
        </div>
        <Button
          onClick={() => void refetch()}
          disabled={isFetching}
          className="flex shrink-0 items-center gap-2 rounded-xl bg-[#FFD700] px-6 py-6 font-semibold text-[#1A3C34] shadow-md transition-all hover:bg-[#e6c200] disabled:cursor-not-allowed disabled:opacity-50 dark:shadow-lg dark:shadow-[#FFD700]/20 dark:hover:bg-[#FFD700]/80"
        >
          <RefreshCw size={18} className={isFetching ? "animate-spin" : ""} />
          Refresh
        </Button>
      </div>

      {isLoading ? (
        <div className="flex h-64 items-center justify-center">
          <Loader2 className={cn(pageLoader, "h-12 w-12")} />
        </div>
      ) : error ? (
        <div className={cn(pageEmpty, "mx-auto max-w-lg p-8 text-destructive")}>
          Error loading posts. Try refreshing the page.
        </div>
      ) : (
        <div className="mx-auto flex max-w-3xl flex-col gap-8">
          {data?.Posts?.map((post) => (
            <GuestPostCard key={post.id} post={post} />
          ))}
          {(!data?.Posts || data.Posts.length === 0) && (
            <div className={pageEmpty}>
              <p className="text-xl font-bold text-[#1A3C34] dark:text-gray-300">
                No posts yet
              </p>
              <p className="mt-2 text-gray-500">
                Be the first to share a bite!
              </p>
            </div>
          )}
        </div>
      )}
    </PageLayout>
  );
}
