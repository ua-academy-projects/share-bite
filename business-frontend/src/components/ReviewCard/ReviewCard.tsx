import type { ReviewResponse } from "@/types/api";
import { Card, CardContent } from "@/components/ui/card";
import { Star } from "lucide-react";

type ReviewCardProps = {
  review: ReviewResponse;
};

export function ReviewCard({ review }: ReviewCardProps) {
  return (
    <Card className="rounded-2xl border-border bg-card-solid">
      <CardContent className="flex gap-4 p-4">
        <img
          src={review.avatarURL || "https://via.placeholder.com/48"}
          alt=""
          className="h-12 w-12 rounded-full border border-border object-cover"
        />
        <div className="min-w-0 flex-1">
          <div className="flex items-center justify-between gap-2">
            <span className="font-semibold text-foreground">
              {review.userName}
            </span>
            <span className="flex items-center gap-0.5 text-accent">
              {Array.from({ length: review.rating }).map((_, i) => (
                <Star key={i} size={14} fill="currentColor" />
              ))}
            </span>
          </div>
          <p className="mt-2 text-sm leading-relaxed text-foreground/90">
            {review.text}
          </p>
          <p className="mt-2 text-xs text-muted-foreground">
            {new Date(review.createdAt).toLocaleDateString()}
          </p>
        </div>
      </CardContent>
    </Card>
  );
}
