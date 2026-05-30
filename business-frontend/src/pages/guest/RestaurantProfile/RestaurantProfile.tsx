import { Link, useParams } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";
import axios from "axios";
import { Loader2, MapPin, Star } from "lucide-react";
import { apiClient } from "@/api/client";
import { GuestPostCard } from "@/components/PostCard/PostCard";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";

type VenueData = {
  id?: number;
  name?: string;
  avatar?: string | null;
  banner?: string | null;
  description?: string | null;
  tags?: string[];
};

export function RestaurantProfile() {
  const { id } = useParams<{ id: string }>();

  const { data: restaurant, isLoading, error } = useQuery({
    queryKey: ["restaurant", id],
    queryFn: async () => {
      const res = await axios.get<VenueData>(`/api/business/org-units/${id}`);
      return res.data;
    },
    enabled: !!id,
    retry: false,
  });

  const { data: postsData, isLoading: postsLoading } = useQuery({
    queryKey: ["venuePosts", id],
    queryFn: () => apiClient.getPosts(20, 0, id),
    enabled: !!id,
  });

  if (isLoading) {
    return (
      <div className="flex justify-center py-24">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (error || !restaurant) {
    return (
      <div className="flex min-h-[50vh] items-center justify-center px-4">
        <Card className="max-w-md rounded-3xl bg-card-solid">
          <CardContent className="p-8 text-center">
            <h2 className="text-xl font-bold">Venue not found</h2>
            <p className="mt-2 text-muted-foreground">
              This location could not be loaded.
            </p>
            <Button asChild className="mt-6 rounded-xl">
              <Link to="/explore">Back to Explore</Link>
            </Button>
          </CardContent>
        </Card>
      </div>
    );
  }

  const posts = postsData?.Posts || [];

  return (
    <div className="px-6 py-8 lg:px-10">
      <Card className="mx-auto mb-10 max-w-4xl overflow-hidden rounded-3xl border-border bg-card-solid">
        <div className="relative h-48 bg-muted sm:h-56">
          <img
            src={
              restaurant.banner ||
              restaurant.avatar ||
              "https://images.unsplash.com/photo-1514933651103-005eec06c04b"
            }
            alt=""
            className="h-full w-full object-cover"
          />
          <div className="absolute inset-0 bg-gradient-to-t from-[#0b0f0e]/90 to-transparent" />
        </div>
        <CardContent className="relative -mt-12 p-6 pt-0">
          <div className="flex flex-wrap items-end justify-between gap-4">
            <div>
              <h1 className="text-3xl font-bold text-foreground">
                {restaurant.name}
              </h1>
              {restaurant.tags?.length ? (
                <div className="mt-2 flex flex-wrap gap-2">
                  {restaurant.tags.map((tag) => (
                    <Badge key={tag} variant="outline">
                      {tag}
                    </Badge>
                  ))}
                </div>
              ) : null}
              {restaurant.description && (
                <p className="mt-3 max-w-2xl text-muted-foreground">
                  {restaurant.description}
                </p>
              )}
              <p className="mt-2 flex items-center gap-1 text-sm text-muted-foreground">
                <MapPin className="h-4 w-4" /> Venue profile
              </p>
            </div>
            <Badge variant="accent" className="gap-1 px-3 py-1">
              <Star className="h-3 w-3" fill="currentColor" /> Guest reviews
            </Badge>
          </div>
        </CardContent>
      </Card>

      <h2 className="mb-6 text-xl font-bold">Recent posts</h2>
      {postsLoading ? (
        <div className="flex justify-center py-12">
          <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
        </div>
      ) : posts.length === 0 ? (
        <p className="text-center text-muted-foreground">No posts yet.</p>
      ) : (
        <div className="mx-auto flex max-w-xl flex-col gap-8">
          {posts.map((post) => (
            <GuestPostCard
              key={post.id}
              post={post}
              restaurantName={restaurant.name}
            />
          ))}
        </div>
      )}
    </div>
  );
}
