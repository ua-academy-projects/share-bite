import React, { useState, useEffect, useRef } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";
import {
  ImagePlus,
  ChevronLeft,
  ChevronRight,
  Loader2,
  MapPin,
  Search,
  X,
} from "lucide-react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/api/client";
import { businessApi, type VenueSearchItem } from "@/api/business";
import { PageHeader } from "@/components/layout/PageHeader";
import { PageLayout } from "@/components/layout/PageLayout";
import {
  pageBtnPrimary,
  pageBtnSecondary,
  pageInput,
  pageLabel,
  pageLinkAccent,
  pagePanel,
} from "@/components/layout/pageStyles";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

const ALLOWED_IMAGE_TYPES = ["image/jpeg", "image/png", "image/jpg"];

const isValidImageType = (file: File) => ALLOWED_IMAGE_TYPES.includes(file.type);

function venueFromProfile(profile: {
  id: number;
  name: string;
  avatar?: string | null;
  description?: string | null;
  latitude?: number | null;
  longitude?: number | null;
  tags: string[];
}): VenueSearchItem {
  return {
    id: profile.id,
    name: profile.name,
    avatar: profile.avatar,
    description: profile.description,
    latitude: profile.latitude,
    longitude: profile.longitude,
    tags: profile.tags,
  };
}

export const CreatePost: React.FC = () => {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { id: venueIdParam } = useParams<{ id?: string }>();
  const presetVenueId = venueIdParam ? Number(venueIdParam) : null;

  const [images, setImages] = useState<{ file: File; previewUrl: string }[]>([]);
  const [selectedVenue, setSelectedVenue] = useState<VenueSearchItem | null>(null);
  const [searchQuery, setSearchQuery] = useState("");
  const [activeSearchQuery, setActiveSearchQuery] = useState("");
  const [text, setText] = useState("");
  const [rating, setRating] = useState(5);
  const [validationError, setValidationError] = useState<string | null>(null);

  const imagesRef = useRef(images);
  useEffect(() => {
    imagesRef.current = images;
  }, [images]);

  useEffect(() => {
    return () => {
      imagesRef.current.forEach((img) => URL.revokeObjectURL(img.previewUrl));
    };
  }, []);

  const { data: presetVenue, isLoading: presetVenueLoading } = useQuery({
    queryKey: ["venueProfile", presetVenueId],
    queryFn: () => businessApi.getVenueProfile(presetVenueId!),
    enabled: !!presetVenueId && presetVenueId > 0,
    retry: false,
  });

  useEffect(() => {
    if (presetVenue) {
      setSelectedVenue(venueFromProfile(presetVenue));
    }
  }, [presetVenue]);

  const {
    data: searchData,
    isFetching: searching,
    error: searchError,
  } = useQuery({
    queryKey: ["venueSearch", "post-create", activeSearchQuery],
    queryFn: () =>
      businessApi.searchVenues({
        query: activeSearchQuery,
        limit: 20,
      }),
    enabled: activeSearchQuery.trim().length > 0,
    retry: false,
  });

  const searchResults = searchData?.items ?? [];

  const handleImageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files.length > 0) {
      setValidationError(null);
      const files = Array.from(e.target.files);

      const invalidFiles = files.filter((file) => !isValidImageType(file));
      if (invalidFiles.length > 0) {
        setValidationError("Unsupported image type. Only JPEG and PNG are supported.");
        e.currentTarget.value = "";
        return;
      }

      const allowedFiles = files.slice(0, 5 - images.length);
      const newImages = allowedFiles.map((file) => ({
        file,
        previewUrl: URL.createObjectURL(file),
      }));
      setImages((prev) => [...prev, ...newImages].slice(0, 5));
      e.currentTarget.value = "";
    }
  };

  const removeImage = (index: number) => {
    setImages((prev) => {
      const newImages = [...prev];
      URL.revokeObjectURL(newImages[index].previewUrl);
      newImages.splice(index, 1);
      return newImages;
    });
  };

  const moveImage = (index: number, direction: "left" | "right") => {
    setImages((prev) => {
      const newImages = [...prev];
      if (direction === "left" && index > 0) {
        [newImages[index - 1], newImages[index]] = [newImages[index], newImages[index - 1]];
      } else if (direction === "right" && index < prev.length - 1) {
        [newImages[index], newImages[index + 1]] = [newImages[index + 1], newImages[index]];
      }
      return newImages;
    });
  };

  const handleVenueSearch = () => {
    const query = searchQuery.trim();
    if (!query) {
      setValidationError("Enter a venue name to search.");
      return;
    }
    setValidationError(null);
    setActiveSearchQuery(query);
  };

  const createMutation = useMutation({
    mutationFn: async ({
      finalVenueId,
      parsedRating,
    }: {
      finalVenueId: number;
      parsedRating: number;
    }) => {
      return await apiClient.createPost({
        venueId: finalVenueId,
        text,
        rating: parsedRating,
        images: images.length > 0 ? images.map((img) => img.file) : undefined,
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["posts"] });
      navigate("/");
    },
    onError: (error) => {
      console.error("Failed to create post:", error);
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setValidationError(null);

    if (!selectedVenue) {
      setValidationError("Please select a venue for your review.");
      return;
    }

    const parsedRating = parseInt(rating.toString(), 10);
    if (Number.isNaN(parsedRating) || parsedRating < 1 || parsedRating > 5) {
      setValidationError("Rating must be between 1 and 5");
      return;
    }

    const invalidImage = images.find(
      (img) => !["image/jpeg", "image/png", "image/jpg"].includes(img.file.type)
    );
    if (invalidImage) {
      setValidationError("Unsupported image type. Only JPEG and PNG are supported.");
      return;
    }

    createMutation.mutate({ finalVenueId: selectedVenue.id, parsedRating });
  };

  return (
    <PageLayout>
      <PageHeader
        title="Create a Post"
        description="Share your latest culinary experience"
      />

      <div className={cn(pagePanel, "p-8 md:p-10")}>
        <form className="space-y-6" onSubmit={handleSubmit}>
          <div className="space-y-4">
            <span className={pageLabel}>Photos</span>
            <div className="grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-5 lg:grid-cols-5">
              {images.map((img, idx) => (
                <div
                  key={img.previewUrl}
                  className="group relative aspect-square overflow-hidden rounded-xl border border-gray-200 bg-gray-100 dark:border-[#2f5e50] dark:bg-[#0d241d]"
                >
                  <img
                    src={img.previewUrl}
                    alt={`Preview ${idx}`}
                    className="h-full w-full object-cover"
                  />
                  <div className="absolute inset-0 flex items-center justify-center gap-2 bg-black/40 opacity-0 backdrop-blur-sm transition-opacity group-hover:opacity-100">
                    {idx > 0 ? (
                      <Button
                        type="button"
                        size="icon-xs"
                        variant="secondary"
                        className="rounded-full"
                        onClick={() => moveImage(idx, "left")}
                      >
                        <ChevronLeft size={16} />
                      </Button>
                    ) : null}
                    {idx < images.length - 1 ? (
                      <Button
                        type="button"
                        size="icon-xs"
                        variant="secondary"
                        className="rounded-full"
                        onClick={() => moveImage(idx, "right")}
                      >
                        <ChevronRight size={16} />
                      </Button>
                    ) : null}
                  </div>
                  <Button
                    type="button"
                    size="icon-xs"
                    variant="destructive"
                    className="absolute top-2 right-2 rounded-full"
                    onClick={() => removeImage(idx)}
                  >
                    <X className="h-3 w-3" />
                  </Button>
                </div>
              ))}

              {images.length < 5 ? (
                <label
                  htmlFor="image-upload"
                  className="flex aspect-square cursor-pointer flex-col items-center justify-center rounded-xl border-2 border-dashed border-gray-300 bg-gray-50 text-gray-500 transition-colors hover:bg-gray-100 hover:text-[#1A3C34] dark:border-[#2f5e50] dark:bg-[#0d241d] dark:hover:text-white"
                >
                  <ImagePlus size={28} className="mb-2" />
                  <span className="text-xs font-semibold">Add Photo</span>
                  <input
                    id="image-upload"
                    type="file"
                    accept="image/png, image/jpeg, image/jpg"
                    multiple
                    className="hidden"
                    onChange={handleImageChange}
                  />
                </label>
              ) : null}
            </div>
          </div>

          <div className="space-y-4">
            <label className={cn(pageLabel, "flex items-center gap-1")}>
              <MapPin size={14} /> Venue
            </label>

            {presetVenueLoading ? (
              <div className="flex h-20 items-center justify-center">
                <Loader2 className="h-6 w-6 animate-spin text-emerald-500 dark:text-[#98FF98]" />
              </div>
            ) : selectedVenue ? (
              <div className="flex items-center gap-4 rounded-xl border border-emerald-500/30 bg-emerald-500/5 p-4 dark:border-[#98FF98]/30 dark:bg-[#98FF98]/5">
                <img
                  src={
                    selectedVenue.avatar ||
                    "https://placehold.co/56x56/163d32/FFF?text=SB"
                  }
                  alt={selectedVenue.name}
                  className="h-14 w-14 rounded-xl border border-gray-200 object-cover dark:border-[#2f5e50]"
                />
                <div className="min-w-0 flex-1">
                  <p className="font-semibold text-[#1A3C34] dark:text-white">
                    {selectedVenue.name}
                  </p>
                  <p className="truncate text-sm text-gray-500 dark:text-gray-400">
                    {selectedVenue.description || "Selected venue"}
                  </p>
                </div>
                <Button
                  type="button"
                  variant="outline"
                  className="rounded-xl"
                  onClick={() => setSelectedVenue(null)}
                >
                  Change
                </Button>
              </div>
            ) : (
              <div className="space-y-3">
                <div className="flex flex-col gap-3 sm:flex-row">
                  <input
                    type="text"
                    className={cn(pageInput, "h-12 flex-1")}
                    placeholder="Search by venue name"
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    onKeyDown={(e) => {
                      if (e.key === "Enter") {
                        e.preventDefault();
                        handleVenueSearch();
                      }
                    }}
                  />
                  <Button
                    type="button"
                    className={cn(pageBtnPrimary, "h-12 px-6")}
                    onClick={handleVenueSearch}
                    disabled={searching}
                  >
                    {searching ? (
                      <Loader2 className="h-4 w-4 animate-spin" />
                    ) : (
                      <>
                        <Search className="mr-2 h-4 w-4" />
                        Search
                      </>
                    )}
                  </Button>
                </div>

                {searchError ? (
                  <p className="text-sm text-destructive">
                    {searchError instanceof Error
                      ? searchError.message
                      : "Venue search failed."}
                  </p>
                ) : null}

                {activeSearchQuery && !searching && searchResults.length === 0 ? (
                  <p className="text-sm text-gray-500 dark:text-gray-400">
                    No venues found for &quot;{activeSearchQuery}&quot;. Try another name or{" "}
                    <Link to="/venues/search" className={pageLinkAccent}>
                      browse venues
                    </Link>
                    .
                  </p>
                ) : null}

                {searchResults.length > 0 ? (
                  <div className="max-h-64 space-y-2 overflow-y-auto rounded-xl border border-gray-200 p-2 dark:border-[#2f5e50]">
                    {searchResults.map((venue) => (
                      <button
                        key={venue.id}
                        type="button"
                        onClick={() => {
                          setSelectedVenue(venue);
                          setValidationError(null);
                        }}
                        className="flex w-full items-center gap-3 rounded-xl border border-transparent px-3 py-2 text-left transition-colors hover:border-emerald-500/30 hover:bg-emerald-500/5 dark:hover:border-[#98FF98]/30 dark:hover:bg-[#98FF98]/5"
                      >
                        <img
                          src={
                            venue.avatar ||
                            "https://placehold.co/40x40/163d32/FFF?text=SB"
                          }
                          alt={venue.name}
                          className="h-10 w-10 rounded-lg border border-gray-200 object-cover dark:border-[#2f5e50]"
                        />
                        <div className="min-w-0">
                          <p className="font-medium text-[#1A3C34] dark:text-white">
                            {venue.name}
                          </p>
                          {venue.description ? (
                            <p className="truncate text-xs text-gray-500 dark:text-gray-400">
                              {venue.description}
                            </p>
                          ) : null}
                        </div>
                      </button>
                    ))}
                  </div>
                ) : null}

                <p className="text-xs text-gray-500 dark:text-gray-400">
                  Can&apos;t find it?{" "}
                  <Link to="/venues/search" className={pageLinkAccent}>
                    Browse all venues
                  </Link>
                </p>
              </div>
            )}
          </div>

          <div className="space-y-2">
            <label htmlFor="rating" className={pageLabel}>
              Rating (1-5)
            </label>
            <input
              id="rating"
              type="number"
              min="1"
              max="5"
              className={cn(pageInput, "h-12 max-w-xs")}
              required
              value={rating}
              onChange={(e) => setRating(parseInt(e.target.value, 10))}
            />
          </div>

          <div className="space-y-2">
            <label htmlFor="review" className={pageLabel}>
              Review
            </label>
            <textarea
              id="review"
              className={cn(pageInput, "min-h-[120px] resize-y py-3")}
              placeholder="What did you think of the food and atmosphere?"
              required
              rows={5}
              value={text}
              onChange={(e) => setText(e.target.value)}
            />
          </div>

          {validationError ? (
            <div className="rounded-xl border border-destructive/20 bg-destructive/10 p-4 text-sm font-bold text-destructive">
              {validationError}
            </div>
          ) : null}

          {createMutation.isError ? (
            <div className="rounded-xl border border-destructive/20 bg-destructive/10 p-4 text-sm font-bold text-destructive">
              {(createMutation.error as { response?: { data?: { error?: string } } })
                ?.response?.data?.error || "Failed to create post. Please try again."}
            </div>
          ) : null}

          <div className="mt-4 flex justify-end gap-3 border-t border-gray-200 pt-6 dark:border-[#2f5e50]">
            <Button type="button" className={pageBtnSecondary} onClick={() => navigate(-1)}>
              Cancel
            </Button>
            <Button
              type="submit"
              disabled={createMutation.isPending}
              className={pageBtnPrimary}
            >
              {createMutation.isPending ? "Posting..." : "Post Review"}
            </Button>
          </div>
        </form>
      </div>
    </PageLayout>
  );
};
