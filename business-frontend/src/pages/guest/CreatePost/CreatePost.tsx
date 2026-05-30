import { useEffect, useRef, useState } from "react";
import { useNavigate } from "react-router-dom";
import { ImagePlus, MapPin, Star } from "lucide-react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { apiClient } from "@/api/client";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Card, CardContent } from "@/components/ui/card";
import { PageHeader } from "@/components/layout/PageHeader";

const ALLOWED = ["image/jpeg", "image/png", "image/jpg"];

export function CreatePost() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [images, setImages] = useState<{ file: File; previewUrl: string }[]>(
    []
  );
  const [venueId, setVenueId] = useState("");
  const [text, setText] = useState("");
  const [rating, setRating] = useState(5);
  const [validationError, setValidationError] = useState<string | null>(null);
  const imagesRef = useRef(images);

  useEffect(() => {
    imagesRef.current = images;
  }, [images]);

  useEffect(
    () => () => {
      imagesRef.current.forEach((img) => URL.revokeObjectURL(img.previewUrl));
    },
    []
  );

  const createMutation = useMutation({
    mutationFn: async () => {
      const finalVenueId = parseInt(venueId.trim(), 10);
      return apiClient.createPost({
        venueId: finalVenueId,
        text,
        rating,
        images: images.map((img) => img.file),
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["posts"] });
      toast.success("Post published!");
      navigate("/");
    },
    onError: () => toast.error("Failed to create post"),
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const finalVenueId = parseInt(venueId.trim(), 10);
    if (Number.isNaN(finalVenueId) || finalVenueId <= 0) {
      setValidationError("Enter a valid venue ID.");
      return;
    }
    if (!text.trim()) {
      setValidationError("Write something about your bite.");
      return;
    }
    setValidationError(null);
    createMutation.mutate();
  };

  return (
    <div className="px-6 py-8 lg:px-10">
      <PageHeader
        title="Share a Bite"
        description="Tell the community about your dining experience"
      />

      <Card className="mx-auto max-w-xl rounded-3xl bg-card-solid">
        <CardContent className="p-6">
          <form onSubmit={handleSubmit} className="space-y-5">
            <div className="space-y-2">
              <Label htmlFor="venueId">
                <MapPin className="mr-1 inline h-4 w-4" /> Venue ID
              </Label>
              <Input
                id="venueId"
                value={venueId}
                onChange={(e) => setVenueId(e.target.value)}
                placeholder="e.g. 1"
                className="rounded-xl"
              />
            </div>

            <div className="space-y-2">
              <Label>Rating</Label>
              <div className="flex gap-1">
                {[1, 2, 3, 4, 5].map((star) => (
                  <button
                    key={star}
                    type="button"
                    onClick={() => setRating(star)}
                  >
                    <Star
                      size={24}
                      className={
                        star <= rating ? "text-accent" : "text-muted-foreground/40"
                      }
                      fill={star <= rating ? "currentColor" : "none"}
                    />
                  </button>
                ))}
              </div>
            </div>

            <div className="space-y-2">
              <Label htmlFor="text">Your review</Label>
              <Textarea
                id="text"
                value={text}
                onChange={(e) => setText(e.target.value)}
                className="min-h-[120px] rounded-xl"
                placeholder="What did you love about this place?"
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="images">Photos (optional)</Label>
              <label className="flex cursor-pointer items-center gap-2 rounded-xl border border-dashed border-border p-4 text-sm text-muted-foreground hover:bg-muted/20">
                <ImagePlus className="h-5 w-5" />
                Add up to 5 images
                <input
                  id="images"
                  type="file"
                  accept={ALLOWED.join(",")}
                  multiple
                  className="hidden"
                  onChange={(e) => {
                    const files = Array.from(e.target.files || []).filter((f) =>
                      ALLOWED.includes(f.type)
                    );
                    setImages((prev) =>
                      [...prev, ...files.map((file) => ({
                        file,
                        previewUrl: URL.createObjectURL(file),
                      }))].slice(0, 5)
                    );
                  }}
                />
              </label>
              {images.length > 0 && (
                <div className="flex flex-wrap gap-2">
                  {images.map((img, i) => (
                    <img
                      key={i}
                      src={img.previewUrl}
                      alt=""
                      className="h-20 w-20 rounded-xl object-cover"
                    />
                  ))}
                </div>
              )}
            </div>

            {validationError && (
              <p className="text-sm text-destructive">{validationError}</p>
            )}

            <Button
              type="submit"
              disabled={createMutation.isPending}
              className="w-full rounded-xl bg-accent font-bold text-accent-foreground"
            >
              {createMutation.isPending ? "Publishing…" : "Publish post"}
            </Button>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
