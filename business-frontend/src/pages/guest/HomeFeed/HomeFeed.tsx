import { useQuery } from "@tanstack/react-query";
import { Loader2 } from "lucide-react";
import { apiClient } from "@/api/client";
import { GuestPostCard } from "@/components/PostCard/PostCard";
import { PageHeader } from "@/components/layout/PageHeader";

export function HomeFeed() {
  const { data, isLoading, error } = useQuery({
    queryKey: ["posts"],
    queryFn: () => apiClient.getPosts(20, 0),
  });

  return (
    <div className="px-6 py-8 lg:px-10">
      <PageHeader
        title="The Home Feed"
        description="Discover what your friends are eating"
        className="justify-center text-center [&>div]:mx-auto"
      />

      {isLoading ? (
        <div className="flex justify-center py-16">
          <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
        </div>
      ) : error ? (
        <div className="mx-auto max-w-lg rounded-2xl border border-destructive/30 bg-destructive/10 p-6 text-center text-destructive">
          Error loading posts. Try refreshing the page.
        </div>
      ) : (
        <div className="mx-auto flex max-w-xl flex-col gap-8">
          {data?.Posts?.map((post) => (
            <GuestPostCard key={post.id} post={post} />
          ))}
          {(!data?.Posts || data.Posts.length === 0) && (
            <p className="py-16 text-center text-muted-foreground">
              No posts yet. Be the first to share a bite!
            </p>
          )}
        </div>
      )}
    </div>
  );
}
