import { useQuery } from "@tanstack/react-query";
import { Loader2 } from "lucide-react";
import { apiClient } from "@/api/client";
import { GuestPostCard } from "@/components/PostCard/PostCard";
import { PageHeader } from "@/components/layout/PageHeader";
import { PageLayout } from "@/components/layout/PageLayout";
import { pageEmpty, pageLoader } from "@/components/layout/pageStyles";
import { cn } from "@/lib/utils";

export function HomeFeed() {
  const { data, isLoading, error } = useQuery({
    queryKey: ["posts"],
    queryFn: () => apiClient.getPosts(20, 0),
  });

  return (
    <PageLayout maxWidth="xl">
      <PageHeader
        title="The Home Feed"
        description="Discover what your friends are eating"
        className="justify-center text-center [&>div]:mx-auto"
      />

      {isLoading ? (
        <div className="flex h-64 items-center justify-center">
          <Loader2 className={cn(pageLoader, "h-12 w-12")} />
        </div>
      ) : error ? (
        <div className={cn(pageEmpty, "mx-auto max-w-lg p-8 text-destructive")}>
          Error loading posts. Try refreshing the page.
        </div>
      ) : (
        <div className="mx-auto flex flex-col gap-8">
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
