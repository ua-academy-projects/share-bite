import { useState } from "react";
import { Link } from "react-router-dom";
import {
  Heart,
  MessageCircle,
  Share2,
  Bookmark,
  Trash2,
  Edit2,
  ChevronLeft,
  ChevronRight,
  Plus,
  MoreHorizontal,
  Star,
} from "lucide-react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { clsx } from "clsx";
import type { PostResponse } from "@/types/api";
import { apiClient } from "@/api/client";
import { useCurrentCustomer } from "@/hooks/useCurrentCustomer";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { Input } from "@/components/ui/input";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

type GuestPostCardProps = {
  post: PostResponse;
  restaurantName?: string;
};

export function GuestPostCard({ post, restaurantName }: GuestPostCardProps) {
  const queryClient = useQueryClient();
  const { data: currentCustomer } = useCurrentCustomer();
  const isLiked = post.isLikedByMe || false;
  const likesCount = post.likesCount || 0;
  const [showComments, setShowComments] = useState(false);
  const [commentText, setCommentText] = useState("");
  const [currentImageIndex, setCurrentImageIndex] = useState(0);
  const [editingCommentId, setEditingCommentId] = useState<number | null>(null);
  const [editCommentText, setEditCommentText] = useState("");
  const [isCollectionModalOpen, setIsCollectionModalOpen] = useState(false);
  const [isEditing, setIsEditing] = useState(false);
  const [editValue, setEditValue] = useState(post.text);
  const [editRating, setEditRating] = useState(post.rating);
  const isOwner = currentCustomer?.id === post.customerId;

  const { data: commentsData, isLoading: isLoadingComments } = useQuery({
    queryKey: ["comments", post.id],
    queryFn: () => apiClient.getComments(post.id),
    enabled: showComments,
  });

  const { data: collections } = useQuery({
    queryKey: ["collections", "mine"],
    queryFn: () => apiClient.getCollections(),
    enabled: isCollectionModalOpen,
  });

  const invalidatePosts = () => {
    ["posts", "userPosts", "venuePosts"].forEach((key) =>
      queryClient.invalidateQueries({ queryKey: [key] })
    );
  };

  const saveToCollectionMutation = useMutation({
    mutationFn: (collectionId: string | number) => {
      if (!post.venueId) throw new Error("Missing venue ID");
      return apiClient.savePostToCollection(collectionId, post.venueId);
    },
    onSuccess: () => {
      toast.success("Saved to collection");
      setIsCollectionModalOpen(false);
    },
    onError: (error: unknown) => {
      const err = error as { response?: { status?: number; data?: { error?: string } } };
      if (err?.response?.status === 409) toast.info("Already in collection");
      else toast.error(err?.response?.data?.error || "Failed to save");
      setIsCollectionModalOpen(false);
    },
  });

  const createCommentMutation = useMutation({
    mutationFn: (text: string) => apiClient.createComment(post.id, text),
    onSuccess: () => {
      setCommentText("");
      queryClient.invalidateQueries({ queryKey: ["comments", post.id] });
    },
  });

  const deleteCommentMutation = useMutation({
    mutationFn: (commentId: number) =>
      apiClient.deleteComment(post.id, commentId),
    onSuccess: () =>
      queryClient.invalidateQueries({ queryKey: ["comments", post.id] }),
  });

  const updateCommentMutation = useMutation({
    mutationFn: ({ commentId, text }: { commentId: number; text: string }) =>
      apiClient.updateComment(post.id, commentId, text),
    onSuccess: () => {
      setEditingCommentId(null);
      queryClient.invalidateQueries({ queryKey: ["comments", post.id] });
    },
  });

  const toggleLikeMutation = useMutation({
    mutationFn: async (nextLiked: boolean) => {
      if (nextLiked) await apiClient.likePost(post.id);
      else await apiClient.unlikePost(post.id);
    },
    onSettled: invalidatePosts,
  });

  const deletePostMutation = useMutation({
    mutationFn: () => apiClient.deletePost(post.id),
    onSuccess: () => {
      toast.success("Post deleted");
      invalidatePosts();
    },
    onError: (error: unknown) => {
      const err = error as { response?: { data?: { error?: string } } };
      toast.error(err?.response?.data?.error || "Failed to delete post");
    },
  });

  const updatePostMutation = useMutation({
    mutationFn: ({ text, rating }: { text: string; rating: number }) =>
      apiClient.updatePost(post.id, { text, rating }),
    onSuccess: () => {
      toast.success("Post updated");
      setIsEditing(false);
      invalidatePosts();
    },
    onError: (error: unknown) => {
      const err = error as { response?: { data?: { error?: string } } };
      toast.error(err?.response?.data?.error || "Failed to update post");
    },
  });

  const venueLabel = restaurantName || `Venue #${post.venueId}`;

  return (
    <Card className="mx-auto w-full max-w-[520px] overflow-hidden rounded-[2rem] border border-gray-200 bg-white shadow-xl dark:border-[#2f5e50]/60 dark:bg-[#112f26]">
      <CardContent className="p-0">
        {post.images?.length > 0 && (
          <div className="relative h-64 overflow-hidden bg-muted">
            <img
              src={post.images[currentImageIndex]}
              alt=""
              className="h-full w-full object-cover"
            />
            <div className="absolute inset-0 bg-gradient-to-t from-black/60 to-transparent" />
            <Link
              to={`/restaurant/${post.venueId}`}
              className="absolute left-4 top-4 rounded-full border border-primary/30 bg-primary/20 px-3 py-1 text-xs font-semibold text-primary-foreground backdrop-blur-md"
            >
              {venueLabel}
            </Link>
            <span className="absolute right-4 top-4 rounded-full bg-accent px-3 py-1 text-xs font-bold text-accent-foreground">
              ★ {post.rating}
            </span>
            {post.images.length > 1 && (
              <>
                {currentImageIndex > 0 && (
                  <button
                    type="button"
                    className="absolute left-3 top-1/2 -translate-y-1/2 rounded-full bg-black/40 p-2 text-white"
                    onClick={() => setCurrentImageIndex((i) => i - 1)}
                  >
                    <ChevronLeft size={18} />
                  </button>
                )}
                {currentImageIndex < post.images.length - 1 && (
                  <button
                    type="button"
                    className="absolute right-3 top-1/2 -translate-y-1/2 rounded-full bg-black/40 p-2 text-white"
                    onClick={() => setCurrentImageIndex((i) => i + 1)}
                  >
                    <ChevronRight size={18} />
                  </button>
                )}
              </>
            )}
          </div>
        )}

        <div className="p-5">
          <div className="mb-4 flex items-start justify-between gap-3">
            <div className="flex items-center gap-3">
              <img
                src={post.avatarURL || "https://via.placeholder.com/40"}
                alt=""
                className="h-10 w-10 rounded-full border border-border object-cover"
              />
              <div>
                <Link
                  to={`/user/${post.customerUsername || post.userName}`}
                  className="font-bold text-foreground hover:text-primary"
                >
                  {post.userName}
                </Link>
                <p className="text-xs text-muted-foreground">
                  {new Date(post.createdAt).toLocaleDateString()}
                </p>
              </div>
            </div>
            {isOwner && (
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button variant="ghost" size="icon-sm" className="rounded-full">
                    <MoreHorizontal className="h-4 w-4" />
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end">
                  <DropdownMenuItem onClick={() => setIsEditing(true)}>
                    <Edit2 className="h-4 w-4" /> Edit
                  </DropdownMenuItem>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem
                    variant="destructive"
                    onClick={() => {
                      if (window.confirm("Delete this post?"))
                        deletePostMutation.mutate();
                    }}
                  >
                    <Trash2 className="h-4 w-4" /> Delete
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            )}
          </div>

          {!post.images?.length && (
            <div className="mb-3 flex flex-wrap gap-2">
              <Link
                to={`/restaurant/${post.venueId}`}
                className="rounded-full bg-primary/10 px-2 py-1 text-xs font-semibold text-primary"
              >
                {venueLabel}
              </Link>
              <span className="rounded-full bg-accent/15 px-2 py-1 text-xs font-bold text-accent">
                ★ {post.rating}
              </span>
            </div>
          )}

          {isEditing ? (
            <div className="mb-4 space-y-3">
              <div className="flex gap-1">
                {[1, 2, 3, 4, 5].map((star) => (
                  <button
                    key={star}
                    type="button"
                    onClick={() => setEditRating(star)}
                  >
                    <Star
                      size={22}
                      className={clsx(
                        star <= editRating ? "text-accent" : "text-muted-foreground/40"
                      )}
                      fill={star <= editRating ? "currentColor" : "none"}
                    />
                  </button>
                ))}
              </div>
              <Textarea
                value={editValue}
                onChange={(e) => setEditValue(e.target.value)}
                className="min-h-[100px] rounded-xl bg-muted/30"
              />
              <div className="flex justify-end gap-2">
                <Button variant="ghost" size="sm" onClick={() => setIsEditing(false)}>
                  Cancel
                </Button>
                <Button
                  size="sm"
                  className="rounded-full bg-primary text-primary-foreground"
                  disabled={updatePostMutation.isPending}
                  onClick={() =>
                    updatePostMutation.mutate({ text: editValue, rating: editRating })
                  }
                >
                  Save
                </Button>
              </div>
            </div>
          ) : (
            <p className="mb-4 text-[15px] leading-relaxed text-foreground/90">
              {post.text}
            </p>
          )}

          <div className="flex items-center justify-between border-t border-border/60 pt-4">
            <div className="flex items-center gap-4">
              <button
                type="button"
                className={clsx(
                  "flex items-center gap-1.5 transition-colors",
                  isLiked ? "text-secondary" : "text-muted-foreground"
                )}
                onClick={() => toggleLikeMutation.mutate(!isLiked)}
              >
                <Heart size={20} fill={isLiked ? "currentColor" : "none"} />
                <span className="text-sm font-semibold">{likesCount}</span>
              </button>
              <button
                type="button"
                className="flex items-center gap-1.5 text-muted-foreground hover:text-foreground"
                onClick={() => setShowComments(!showComments)}
              >
                <MessageCircle size={20} />
                <span className="text-sm font-semibold">
                  {commentsData?.total ?? "…"}
                </span>
              </button>
              <button
                type="button"
                className="text-muted-foreground hover:text-foreground"
                onClick={() => {
                  navigator.clipboard.writeText(
                    `${window.location.origin}/post/${post.id}`
                  );
                  toast.success("Link copied");
                }}
              >
                <Share2 size={20} />
              </button>
            </div>
            <button
              type="button"
              className="text-muted-foreground hover:text-primary"
              onClick={() => setIsCollectionModalOpen(true)}
            >
              <Bookmark size={20} />
            </button>
          </div>

          {showComments && (
            <div className="mt-4 space-y-3 border-t border-border/40 pt-4">
              {isLoadingComments ? (
                <p className="text-center text-sm text-muted-foreground">
                  Loading comments…
                </p>
              ) : commentsData?.entities.length === 0 ? (
                <p className="text-center text-sm text-muted-foreground">
                  No comments yet
                </p>
              ) : (
                commentsData?.entities.map((comment) => (
                  <div key={comment.id} className="text-sm">
                    <div className="flex items-center gap-2">
                      <span className="font-semibold">
                        {comment.customer.userName}
                      </span>
                      <span className="text-xs text-muted-foreground">
                        {new Date(comment.createdAt).toLocaleDateString()}
                      </span>
                    </div>
                    {editingCommentId === comment.id ? (
                      <div className="mt-1 flex gap-2">
                        <Input
                          value={editCommentText}
                          onChange={(e) => setEditCommentText(e.target.value)}
                          className="h-8"
                        />
                        <Button
                          size="sm"
                          onClick={() =>
                            updateCommentMutation.mutate({
                              commentId: comment.id,
                              text: editCommentText,
                            })
                          }
                        >
                          Save
                        </Button>
                      </div>
                    ) : (
                      <>
                        <p className="mt-0.5 text-foreground/90">{comment.text}</p>
                        {currentCustomer?.id === comment.customer.id && (
                          <div className="mt-1 flex gap-2 text-xs">
                            <button
                              type="button"
                              className="text-muted-foreground hover:text-primary"
                              onClick={() => {
                                setEditingCommentId(comment.id);
                                setEditCommentText(comment.text);
                              }}
                            >
                              Edit
                            </button>
                            <button
                              type="button"
                              className="text-muted-foreground hover:text-destructive"
                              onClick={() =>
                                deleteCommentMutation.mutate(comment.id)
                              }
                            >
                              Delete
                            </button>
                          </div>
                        )}
                      </>
                    )}
                  </div>
                ))
              )}
              <form
                className="flex gap-2"
                onSubmit={(e) => {
                  e.preventDefault();
                  if (commentText.trim())
                    createCommentMutation.mutate(commentText);
                }}
              >
                <Input
                  placeholder="Add a comment…"
                  value={commentText}
                  onChange={(e) => setCommentText(e.target.value)}
                  className="rounded-full"
                />
                <Button type="submit" size="sm" className="rounded-full">
                  Post
                </Button>
              </form>
            </div>
          )}
        </div>
      </CardContent>

      <Dialog open={isCollectionModalOpen} onOpenChange={setIsCollectionModalOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Save to collection</DialogTitle>
          </DialogHeader>
          <div className="flex max-h-60 flex-col gap-2 overflow-y-auto">
            {!collections?.length ? (
              <p className="text-sm text-muted-foreground">No collections yet.</p>
            ) : (
              collections.map((c) => (
                <Button
                  key={c.id}
                  variant="outline"
                  className="justify-start"
                  onClick={() => saveToCollectionMutation.mutate(c.id)}
                >
                  <Plus className="mr-2 h-4 w-4" />
                  {c.name}
                </Button>
              ))
            )}
          </div>
        </DialogContent>
      </Dialog>
    </Card>
  );
}
