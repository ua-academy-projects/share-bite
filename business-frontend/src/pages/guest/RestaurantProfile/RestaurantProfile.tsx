import { Link, useParams } from "react-router-dom";
import { MapPin, Star, Loader2 } from "lucide-react";
import { useQuery } from "@tanstack/react-query";
import { apiClient } from "@/api/client";
import { GuestPostCard } from "@/components/PostCard/PostCard";
import { PageHeader } from "@/components/layout/PageHeader";
import { PageLayout } from "@/components/layout/PageLayout";
import {
  pageBtnPrimary,
  pageEmpty,
  pageLoader,
  pagePanel,
} from "@/components/layout/pageStyles";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

export function RestaurantProfile() {
  const { id } = useParams<{ id: string }>();

  const { data: restaurant, isLoading: resLoading, error: resError } = useQuery({
    queryKey: ["restaurant", id],
    queryFn: () => apiClient.getRestaurant(id!),
    enabled: !!id,
    retry: false,
  });

  const { data: postsData, isLoading: postsLoading } = useQuery({
    queryKey: ["restaurantPosts", id],
    queryFn: () => apiClient.getPosts(20, 0, id),
    enabled: !!id,
  });

  if (resLoading) {
    return (
      <PageLayout maxWidth="xl">
        <div className="flex h-64 items-center justify-center">
          <Loader2 className={cn(pageLoader, "h-12 w-12")} />
        </div>
      </PageLayout>
    );
  }

  if (resError || !restaurant) {
    return (
      <PageLayout maxWidth="lg" center>
        <div className={cn(pageEmpty, "max-w-lg p-8")}>
          <h2 className="text-xl font-bold text-[#1A3C34] dark:text-white">
            Restaurant not found
          </h2>
          <p className="mt-2 text-gray-500">This location could not be found.</p>
          <Button asChild className={cn(pageBtnPrimary, "mt-6")}>
            <Link to="/explore">Back to Explore</Link>
          </Button>
        </div>
      </PageLayout>
    );
  }

  const restaurantPosts = postsData?.Posts || [];
  const venue = restaurant as {
    name?: string;
    avatar?: string;
    banner?: string;
    description?: string;
    tags?: string[];
    latitude?: number;
    longitude?: number;
  };

  return (
    <PageLayout maxWidth="xl">
      <div className={cn(pagePanel, "mb-8 overflow-hidden")}>
        <div className="relative h-48 bg-gray-100 md:h-64 dark:bg-[#0d241d]">
          <img
            src={
              venue.banner ||
              venue.avatar ||
              "https://images.unsplash.com/photo-1514933651103-005eec06c04b"
            }
            alt={venue.name || "Restaurant"}
            className="h-full w-full object-cover"
          />
          <div className="absolute inset-0 bg-gradient-to-t from-[#0d241d]/90 to-transparent" />
        </div>
        <div className="relative -mt-16 space-y-4 p-6 md:p-8">
          <div className="flex flex-wrap items-start justify-between gap-4">
            <div>
              <h1 className="text-3xl font-bold text-[#1A3C34] dark:text-white md:text-4xl">
                {venue.name}
              </h1>
              {venue.tags?.length ? (
                <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
                  {venue.tags.join(" · ")}
                </p>
              ) : null}
            </div>
            <div className="flex items-center gap-1 rounded-full bg-emerald-500/20 px-3 py-1 text-emerald-600 dark:text-[#98FF98]">
              <Star className="h-4 w-4 fill-current" />
              <span className="text-sm font-semibold">Venue</span>
            </div>
          </div>
          {venue.description ? (
            <p className="max-w-2xl text-gray-600 dark:text-gray-400">{venue.description}</p>
          ) : null}
          <div className="flex items-center gap-2 text-sm text-gray-500 dark:text-gray-400">
            <MapPin className="h-4 w-4" />
            <span>
              {venue.latitude != null && venue.longitude != null
                ? `${venue.latitude}, ${venue.longitude}`
                : "Location TBA"}
            </span>
          </div>
        </div>
      </div>

      <PageHeader title="Recent Posts" />
      <div className="flex flex-col gap-8">
        {postsLoading ? (
          <div className="flex justify-center py-12">
            <Loader2 className={cn(pageLoader, "h-12 w-12")} />
          </div>
        ) : restaurantPosts.length > 0 ? (
          restaurantPosts.map((post) => (
            <GuestPostCard key={post.id} post={post} restaurantName={venue.name} />
          ))
        ) : (
          <div className={pageEmpty}>
            <p className="text-gray-500">No posts yet.</p>
          </div>
        )}
      </div>
    </PageLayout>
  );
}
