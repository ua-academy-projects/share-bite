import React, { useState } from 'react';
import { Heart, MessageCircle, Share2, Bookmark, Trash2, Edit2, ChevronLeft, ChevronRight, Plus, MoreHorizontal, Star } from 'lucide-react';
import type { PostResponse } from '@/types/api';
import { clsx } from 'clsx';
import { Link } from 'react-router-dom';
import { useMutation, useQueryClient, useQuery } from '@tanstack/react-query';
import { apiClient } from '@/api/client';
import { useCurrentCustomer } from '@/hooks/useCurrentCustomer';
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { toast } from 'sonner';
import { 
  DropdownMenu, 
  DropdownMenuContent, 
  DropdownMenuItem, 
  DropdownMenuTrigger,
  DropdownMenuSeparator
} from "@/components/ui/dropdown-menu";
import { Textarea } from "@/components/ui/textarea";

interface PostCardProps {
  post: PostResponse;
  restaurantName?: string;
}

export function GuestPostCard({ post, restaurantName }: PostCardProps) {
  const queryClient = useQueryClient();
  const { data: currentCustomer } = useCurrentCustomer();
  const isLiked = post.isLikedByMe || false;
  const likesCount = post.likesCount || 0;
  const [showComments, setShowComments] = useState(false);
  const [commentText, setCommentText] = useState('');
  const [currentImageIndex, setCurrentImageIndex] = useState(0);
  const [editingCommentId, setEditingCommentId] = useState<number | null>(null);
  const [editCommentText, setEditCommentText] = useState('');
  const [commentError, setCommentError] = useState<string | null>(null);

  const [isCollectionModalOpen, setIsCollectionModalOpen] = useState(false);
  const [isEditing, setIsEditing] = useState(false);
  const [editValue, setEditValue] = useState(post.text);
  const [editRating, setEditRating] = useState(post.rating);

  const isOwner = currentCustomer?.id === post.customerId;

  const handleNextImage = () => {
    if (post.images && currentImageIndex < post.images.length - 1) {
      setCurrentImageIndex(prev => prev + 1);
    }
  };

  const handlePrevImage = () => {
    if (currentImageIndex > 0) {
      setCurrentImageIndex(prev => prev - 1);
    }
  };

  const { data: commentsData, isLoading: isLoadingComments } = useQuery({
    queryKey: ['comments', post.id],
    queryFn: () => apiClient.getComments(post.id),
    enabled: showComments,
  });

  const { data: collections } = useQuery({
    queryKey: ['collections', 'mine'],
    queryFn: () => apiClient.getCollections(),
    enabled: isCollectionModalOpen,
  });

  const saveToCollectionMutation = useMutation({
    mutationFn: (collectionId: string | number) => {
      console.log('[PostCard] Saving post to collection:', {
        postId: post.id,
        collectionId,
        venueId: post.venueId,
        fullPost: post
      });
      if (post.venueId === undefined || post.venueId === 0) {
        toast.error("Cannot save: Venue ID is missing for this post.");
        throw new Error("Missing Venue ID");
      }
      return apiClient.savePostToCollection(collectionId, post.venueId);
    },
    onSuccess: () => {
      toast.success("Saved to collection");
      setIsCollectionModalOpen(false);
    },
    onError: (error: any) => {
      if (error?.response?.status === 409) {
        toast.info("This post is already in your collection");
      } else {
        toast.error(error?.response?.data?.error || "Failed to save to collection");
      }
      setIsCollectionModalOpen(false);
    }
  });

  const createCommentMutation = useMutation({
    mutationFn: async (text: string) => {
      return await apiClient.createComment(post.id, text);
    },
    onSuccess: () => {
      setCommentText('');
      queryClient.invalidateQueries({ queryKey: ['comments', post.id] });
    }
  });

  const deleteCommentMutation = useMutation({
    mutationFn: async (commentId: number) => {
      return await apiClient.deleteComment(post.id, commentId);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['comments', post.id] });
    },
    onError: (error: any) => {
      console.error("Failed to delete comment:", error);
      setCommentError("Failed to delete comment. Please try again.");
    }
  });

  const updateCommentMutation = useMutation({
    mutationFn: async ({ commentId, text }: { commentId: number; text: string }) => {
      return await apiClient.updateComment(post.id, commentId, text);
    },
    onSuccess: () => {
      setEditingCommentId(null);
      setEditCommentText('');
      setCommentError(null);
      queryClient.invalidateQueries({ queryKey: ['comments', post.id] });
    },
    onError: (error: any) => {
      console.error("Failed to update comment:", error);
      setCommentError("Failed to update comment. Please try again.");
    }
  });

  const toggleLikeMutation = useMutation({
    mutationFn: async (nextLiked: boolean) => {
      if (nextLiked) {
        await apiClient.likePost(post.id);
      } else {
        await apiClient.unlikePost(post.id);
      }
    },
    onMutate: async (nextLiked) => {
      const keys = [['posts'], ['userPosts'], ['venuePosts']] as const;
      await Promise.all(keys.map(k => queryClient.cancelQueries({ queryKey: k })));

      const snapshots = keys.map(k => [k, queryClient.getQueriesData({ queryKey: k })] as const);

      const updater = (oldData: any) => {
        if (!oldData || !oldData.Posts) return oldData;
        return {
          ...oldData,
          Posts: oldData.Posts.map((p: any) => {
            if (p.id === post.id) {
              const increment = nextLiked && !p.isLikedByMe ? 1 : (!nextLiked && p.isLikedByMe ? -1 : 0);
              return {
                ...p,
                isLikedByMe: nextLiked,
                likesCount: Math.max(0, (p.likesCount || 0) + increment)
              };
            }
            return p;
          })
        };
      };

      keys.forEach(k => queryClient.setQueriesData({ queryKey: k }, updater));
      
      return { snapshots };
    },
    onError: (_err, _vars, context) => {
      context?.snapshots.forEach(([_, entries]) => {
        entries.forEach(([qKey, data]) => {
          queryClient.setQueryData(qKey, data);
        });
      });
    },
  });

  const deletePostMutation = useMutation({
    mutationFn: () => apiClient.deletePost(post.id),
    onSuccess: () => {
      toast.success("Post deleted successfully");
      queryClient.invalidateQueries({ queryKey: ['posts'] });
      queryClient.invalidateQueries({ queryKey: ['userPosts'] });
      queryClient.invalidateQueries({ queryKey: ['venuePosts'] });
    },
    onError: (error: any) => {
      toast.error(error?.response?.data?.error || "Failed to delete post");
    }
  });

  const updatePostMutation = useMutation({
    mutationFn: ({ text, rating }: { text: string; rating: number }) => apiClient.updatePost(post.id, { text, rating }),
    onSuccess: () => {
      toast.success("Post updated successfully");
      setIsEditing(false);
      queryClient.invalidateQueries({ queryKey: ['posts'] });
      queryClient.invalidateQueries({ queryKey: ['userPosts'] });
      queryClient.invalidateQueries({ queryKey: ['venuePosts'] });
    },
    onError: (error: any) => {
      toast.error(error?.response?.data?.error || "Failed to update post");
    }
  });

  const handleDeletePost = () => {
    if (window.confirm("Are you sure you want to delete this post?")) {
      deletePostMutation.mutate();
    }
  };

  const handleUpdatePost = () => {
    if (!editValue.trim()) {
      toast.error("Post text cannot be empty");
      return;
    }
    updatePostMutation.mutate({ text: editValue, rating: editRating });
  };

  const handleLike = () => {
    const nextLiked = !isLiked;
    toggleLikeMutation.mutate(nextLiked);
  };

  console.log("DEBUG POST:", post.id, " | CurrentUser:", currentCustomer?.id, " | PostOwner:", post.customerId, " | isOwner:", currentCustomer?.id === post.customerId);

  return (
    <div className="bg-card dark:bg-card border border-border rounded-3xl overflow-hidden shadow-sm hover:shadow-xl dark:shadow-lg dark:hover:shadow-primary/10 hover:-translate-y-1 transition-all duration-300 flex flex-col mb-6 group">
      
      {/* Image Container with Gradient & Badges */}
      {post.images && post.images.length > 0 && (
        <div className="relative h-64 overflow-hidden bg-muted">
          <img 
            src={post.images[currentImageIndex]} 
            alt="Post content" 
            className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500" 
          />
          <div className="absolute inset-0 bg-gradient-to-t from-black/60 via-transparent to-transparent opacity-90"></div>
          
          <div className="absolute top-4 left-4 bg-primary/20 backdrop-blur-md text-primary-foreground text-xs font-semibold px-3 py-1.5 rounded-full border border-primary/30 shadow-sm flex items-center gap-1">
            <Link to={`/restaurant/${post.venueId}`} className="hover:text-primary transition-colors">
              📍 {restaurantName || `Venue #${post.venueId}`}
            </Link>
          </div>
          
          <div className="absolute top-4 right-4 bg-accent text-accent-foreground text-xs font-bold px-3 py-1.5 rounded-full flex items-center gap-1 shadow-md">
            ★ {post.rating}
          </div>

          {post.images.length > 1 && (
            <>
              {currentImageIndex > 0 && (
                <button 
                  className="absolute left-3 top-1/2 -translate-y-1/2 bg-black/40 backdrop-blur-sm text-white rounded-full p-2 opacity-0 group-hover:opacity-100 transition-opacity hover:bg-black/60 border border-white/10" 
                  onClick={handlePrevImage}
                >
                  <ChevronLeft size={20} />
                </button>
              )}
              {currentImageIndex < post.images.length - 1 && (
                <button 
                  className="absolute right-3 top-1/2 -translate-y-1/2 bg-black/40 backdrop-blur-sm text-white rounded-full p-2 opacity-0 group-hover:opacity-100 transition-opacity hover:bg-black/60 border border-white/10" 
                  onClick={handleNextImage}
                >
                  <ChevronRight size={20} />
                </button>
              )}
              <div className="absolute bottom-4 left-0 right-0 flex justify-center gap-2">
                {post.images.map((_, idx) => (
                  <button 
                    type="button"
                    key={idx} 
                    className={clsx(
                      "w-1.5 h-1.5 rounded-full transition-all", 
                      idx === currentImageIndex ? "bg-white w-3" : "bg-white/50 hover:bg-white/80"
                    )}
                    onClick={() => setCurrentImageIndex(idx)}
                  />
                ))}
              </div>
            </>
          )}
        </div>
      )}

      {/* Info Section */}
      <div className="p-5 flex flex-col flex-1 relative z-10">
        <div className="flex justify-between items-start mb-4">
          <div className="flex items-center gap-3">
            <img 
              src={post.avatarURL || 'https://via.placeholder.com/40'} 
              alt={post.userName} 
              className="w-11 h-11 rounded-full object-cover border-2 border-background shadow-sm" 
            />
            <div className="flex flex-col">
              <Link to={`/user/${post.customerUsername}`} className="text-foreground font-bold hover:text-primary transition-colors text-[16px]">
                {post.userName}
              </Link>
              <span className="text-xs text-muted-foreground font-medium">{new Date(post.createdAt).toLocaleDateString()}</span>
            </div>
          </div>
          
          <div className="flex flex-col items-end gap-2">
            {isOwner && (
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button variant="ghost" size="icon" className="h-8 w-8 rounded-full hover:bg-muted/20 cursor-pointer group/options">
                    <MoreHorizontal size={20} className="text-foreground/80 group-hover/options:text-primary transition-colors" />
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end" className="w-44">
                  <DropdownMenuItem onClick={() => setIsEditing(true)} className="gap-2 cursor-pointer py-2.5">
                    <Edit2 size={14} /> Edit Post
                  </DropdownMenuItem>
                  <DropdownMenuSeparator className="bg-[#2f5e50]" />
                  <DropdownMenuItem onClick={handleDeletePost} className="gap-2 text-red-400 focus:text-red-400 cursor-pointer py-2.5" variant="destructive">
                    <Trash2 size={14} /> Delete Post
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            )}
            
            {!(post.images && post.images.length > 0) && (
              <div className="flex flex-col items-end gap-1">
                <Link to={`/restaurant/${post.venueId}`} className="text-xs font-semibold text-primary hover:underline bg-primary/10 px-2 py-1 rounded-full">
                  📍 {restaurantName || `Venue #${post.venueId}`}
                </Link>
                <span className="text-xs font-bold text-accent bg-accent/10 px-2 py-1 rounded-full">★ {post.rating}</span>
              </div>
            )}
          </div>
        </div>

        {isEditing ? (
          <div className="flex flex-col gap-4 mb-6">
            <div className="flex flex-col gap-2">
              <label className="text-xs font-black uppercase tracking-widest text-muted-foreground">Rating</label>
              <div className="flex gap-2">
                {[1, 2, 3, 4, 5].map((star) => (
                  <button
                    key={star}
                    type="button"
                    onClick={() => setEditRating(star)}
                    className="transition-transform hover:scale-110 active:scale-95 outline-none"
                  >
                    <Star 
                      size={26} 
                      fill={star <= editRating ? "currentColor" : "none"}
                      className={clsx(
                        "transition-all duration-200",
                        star <= editRating ? "text-accent" : "text-muted-foreground/40 hover:text-accent/50"
                      )} 
                    />
                  </button>
                ))}
              </div>
            </div>

            <div className="flex flex-col gap-2">
              <label className="text-xs font-black uppercase tracking-widest text-muted-foreground">Review</label>
              <Textarea 
                value={editValue}
                onChange={(e) => setEditValue(e.target.value)}
                className="bg-muted/50 border-border min-h-[120px] rounded-xl focus:ring-primary text-[15px] leading-relaxed"
                placeholder="What's on your mind?"
              />
            </div>

            <div className="flex justify-end gap-2">
              <Button 
                variant="ghost" 
                size="sm" 
                onClick={() => {
                  setIsEditing(false);
                  setEditValue(post.text);
                  setEditRating(post.rating);
                }}
                className="rounded-full h-9 px-5 font-bold text-muted-foreground hover:text-foreground"
              >
                Cancel
              </Button>
              <Button 
                size="sm" 
                onClick={handleUpdatePost}
                disabled={updatePostMutation.isPending}
                className="bg-primary text-primary-foreground hover:bg-primary/90 rounded-full h-9 px-8 font-bold shadow-lg"
              >
                {updatePostMutation.isPending ? 'Saving...' : 'Save Changes'}
              </Button>
            </div>
          </div>
        ) : (
          <p className="text-foreground/90 text-[15px] leading-relaxed mb-6">
            {post.text}
          </p>
        )}

        {/* Action Bar */}
        <div className="mt-auto pt-4 border-t border-border/60 flex justify-between items-center">
          <div className="flex gap-5">
            <div className="flex items-center gap-2">
              <button 
                className={clsx("transition-all hover:scale-110", isLiked ? "text-secondary" : "text-muted-foreground hover:text-foreground")} 
                onClick={handleLike}
                disabled={toggleLikeMutation.isPending}
              >
                <Heart size={22} fill={isLiked ? "currentColor" : "none"} />
              </button>
              <span className="text-sm font-semibold text-muted-foreground">
                {likesCount}
              </span>
            </div>
            
            <div className="flex items-center gap-2">
              <button 
                className={clsx("transition-all hover:scale-110", showComments ? "text-primary" : "text-muted-foreground hover:text-foreground")}
                onClick={() => setShowComments(!showComments)}
              >
                <MessageCircle size={22} />
              </button>
              <span className="text-sm font-semibold text-muted-foreground">
                {commentsData ? commentsData.total : '...'}
              </span>
            </div>

            <button 
              className="text-muted-foreground hover:text-foreground transition-all hover:scale-110"
              onClick={() => {
                const url = window.location.origin + '/post/' + post.id;
                navigator.clipboard.writeText(url);
                toast.success("Link copied to clipboard!");
              }}
            >
              <Share2 size={22} />
            </button>
          </div>
          
          <div className="flex items-center gap-2">
            <button 
              className="p-2 text-muted-foreground hover:text-primary hover:bg-muted/50 rounded-full transition-colors"
              onClick={() => setIsCollectionModalOpen(true)}
            >
              <Bookmark size={22} />
            </button>
          </div>
        </div>
      </div>

      {/* Comments Section */}
      {showComments && (
        <div className="px-5 pb-5 bg-muted/10">
          <div className="mt-4 flex flex-col gap-4">
            {isLoadingComments ? (
              <div className="text-sm text-muted-foreground text-center py-2">Loading comments...</div>
            ) : commentsData?.entities.length === 0 ? (
              <div className="text-sm text-muted-foreground text-center py-2">No comments yet. Be the first!</div>
            ) : (
              commentsData?.entities.map(comment => (
                <div key={comment.id} className="flex flex-col gap-1">
                  <div className="flex items-center gap-2">
                    <img 
                      src={comment.customer.avatarURL || 'https://via.placeholder.com/24'} 
                      alt={comment.customer.userName} 
                      className="w-6 h-6 rounded-full object-cover border border-border" 
                    />
                    <span className="font-semibold text-sm">{comment.customer.userName}</span>
                    <span className="text-xs text-muted-foreground">{new Date(comment.createdAt).toLocaleDateString()}</span>
                  </div>
                  <div className="ml-8 text-sm text-foreground/90">
                    {editingCommentId === comment.id ? (
                      <div className="flex flex-col gap-2 mt-1">
                        <input 
                          type="text" 
                          className="flex h-9 w-full rounded-md border border-input bg-background px-3 py-1 text-sm shadow-sm focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring" 
                          value={editCommentText}
                          onChange={e => setEditCommentText(e.target.value)}
                          autoFocus
                        />
                        <div className="flex gap-2">
                          <button 
                            className="inline-flex items-center justify-center rounded-md text-xs font-medium bg-primary text-primary-foreground hover:bg-primary/90 h-7 px-3"
                            onClick={() => updateCommentMutation.mutate({ commentId: comment.id, text: editCommentText })}
                            disabled={!editCommentText.trim() || updateCommentMutation.isPending}
                          >
                            Save
                          </button>
                          <button 
                            className="inline-flex items-center justify-center rounded-md text-xs font-medium border border-input bg-background hover:bg-accent hover:text-accent-foreground h-7 px-3"
                            onClick={() => { setEditingCommentId(null); setEditCommentText(''); }}
                          >
                            Cancel
                          </button>
                        </div>
                      </div>
                    ) : (
                      <>
                        <p>{comment.text}</p>
                        {currentCustomer?.id === comment.customer.id && (
                          <div className="flex gap-3 mt-1 opacity-60 hover:opacity-100 transition-opacity">
                            <button 
                              className="text-xs hover:text-primary transition-colors flex items-center gap-1"
                              onClick={() => {
                                setEditingCommentId(comment.id);
                                setEditCommentText(comment.text);
                              }}
                              title="Edit comment"
                            >
                              <Edit2 size={12} /> Edit
                            </button>
                            <button 
                              className="text-xs hover:text-destructive transition-colors flex items-center gap-1"
                              onClick={() => deleteCommentMutation.mutate(comment.id)}
                              disabled={deleteCommentMutation.isPending}
                              title="Delete comment"
                            >
                              <Trash2 size={12} /> Delete
                            </button>
                          </div>
                        )}
                      </>
                    )}
                  </div>
                </div>
              ))
            )}
          </div>
          
          <form 
            className="mt-4 flex gap-2 items-center" 
            onSubmit={(e) => {
              e.preventDefault();
              if (commentText.trim()) {
                createCommentMutation.mutate(commentText);
              }
            }}
          >
            <input 
              type="text" 
              placeholder="Add a comment..." 
              className="flex h-9 w-full rounded-full border border-input bg-background px-4 py-1 text-sm shadow-sm focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
              value={commentText}
              onChange={e => setCommentText(e.target.value)}
              disabled={createCommentMutation.isPending}
            />
            <button 
              type="submit" 
              className="inline-flex items-center justify-center rounded-full text-sm font-medium bg-primary text-primary-foreground hover:bg-primary/90 h-9 px-4 disabled:opacity-50 disabled:cursor-not-allowed"
              disabled={!commentText.trim() || createCommentMutation.isPending}
            >
              Post
            </button>
          </form>
          {createCommentMutation.isError && (
            <div className="text-xs text-destructive mt-2 pl-2">Failed to post comment</div>
          )}
          {commentError && (
            <div className="text-xs text-destructive mt-2 pl-2">{commentError}</div>
          )}
        </div>
      )}
      <Dialog open={isCollectionModalOpen} onOpenChange={setIsCollectionModalOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Save to Collection</DialogTitle>
          </DialogHeader>
          <div className="py-4 flex flex-col gap-3 max-h-60 overflow-y-auto">
            {collections?.length === 0 ? (
              <p className="text-sm text-muted-foreground text-center">You don't have any collections yet.</p>
            ) : (
              collections?.map(c => (
                <Button 
                  key={c.id} 
                  variant="outline" 
                  className="justify-start"
                  onClick={() => saveToCollectionMutation.mutate(c.id)}
                  disabled={saveToCollectionMutation.isPending}
                >
                  <Plus size={16} className="mr-2" /> {c.name}
                </Button>
              ))
            )}
          </div>
        </DialogContent>
      </Dialog>
    </div>
  );
};
