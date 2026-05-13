import React, { useState } from 'react';
import { Heart, MessageCircle, Share2, Bookmark, Trash2, Edit2, ChevronLeft, ChevronRight, Pencil, Plus } from 'lucide-react';
import type { PostResponse } from '../../types/api';
import { clsx } from 'clsx';
import { Link } from 'react-router-dom';
import { useMutation, useQueryClient, useQuery } from '@tanstack/react-query';
import { apiClient } from '../../api/client';
import { useCurrentCustomer } from '../../hooks/useCurrentCustomer';
import { EditPostModal } from '../EditPostModal/EditPostModal';
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '../ui/dialog';
import { Button } from '../ui/button';
import { toast } from 'sonner';

interface PostCardProps {
  post: PostResponse;
  restaurantName?: string;
}

export const PostCard: React.FC<PostCardProps> = ({ post, restaurantName }) => {
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

  const [isEditModalOpen, setIsEditModalOpen] = useState(false);
  const [isCollectionModalOpen, setIsCollectionModalOpen] = useState(false);

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
    queryKey: ['collections'],
    queryFn: () => apiClient.getCollections(),
    enabled: isCollectionModalOpen,
  });

  const saveToCollectionMutation = useMutation({
    mutationFn: (collectionId: number) => apiClient.savePostToCollection(collectionId, post.venueId),
    onSuccess: () => {
      toast.success("Saved to collection");
      setIsCollectionModalOpen(false);
    },
    onError: () => toast.error("Failed to save to collection")
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
      await queryClient.cancelQueries({ queryKey: ['posts'] });
      await queryClient.cancelQueries({ queryKey: ['userPosts'] });
      await queryClient.cancelQueries({ queryKey: ['venuePosts'] });

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
                likesCount: (p.likesCount || 0) + increment
              };
            }
            return p;
          })
        };
      };

      queryClient.setQueriesData({ queryKey: ['posts'] }, updater);
      queryClient.setQueriesData({ queryKey: ['userPosts'] }, updater);
      queryClient.setQueriesData({ queryKey: ['venuePosts'] }, updater);
      
      return { nextLiked };
    },
    onError: (err, nextLiked) => {
      const rollbackLiked = !nextLiked;
      const updater = (oldData: any) => {
        if (!oldData || !oldData.Posts) return oldData;
        return {
          ...oldData,
          Posts: oldData.Posts.map((p: any) => {
            if (p.id === post.id) {
              const increment = rollbackLiked && !p.isLikedByMe ? 1 : (!rollbackLiked && p.isLikedByMe ? -1 : 0);
              return {
                ...p,
                isLikedByMe: rollbackLiked,
                likesCount: (p.likesCount || 0) + increment
              };
            }
            return p;
          })
        };
      };
      queryClient.setQueriesData({ queryKey: ['posts'] }, updater);
      queryClient.setQueriesData({ queryKey: ['userPosts'] }, updater);
      queryClient.setQueriesData({ queryKey: ['venuePosts'] }, updater);
    }
  });

  const handleLike = () => {
    const nextLiked = !isLiked;
    toggleLikeMutation.mutate(nextLiked);
  };

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
          
          <div className="absolute top-4 left-4 bg-black/60 backdrop-blur-md text-white text-xs font-semibold px-3 py-1.5 rounded-full border border-white/10 shadow-sm flex items-center gap-1">
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
              <Link to={`/user/${post.customerId}`} className="text-foreground font-bold hover:text-primary transition-colors text-[16px]">
                {post.userName}
              </Link>
              <span className="text-xs text-muted-foreground font-medium">{new Date(post.createdAt).toLocaleDateString()}</span>
            </div>
          </div>
          
          {!(post.images && post.images.length > 0) && (
            <div className="flex flex-col items-end gap-1">
              <Link to={`/restaurant/${post.venueId}`} className="text-xs font-semibold text-primary hover:underline bg-primary/10 px-2 py-1 rounded-full">
                📍 {restaurantName || `Venue #${post.venueId}`}
              </Link>
              <span className="text-xs font-bold text-accent bg-accent/10 px-2 py-1 rounded-full">★ {post.rating}</span>
            </div>
          )}
        </div>

        <p className="text-foreground/90 text-[15px] leading-relaxed mb-6">
          {post.text}
        </p>

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

            <button className="text-muted-foreground hover:text-foreground transition-all hover:scale-110">
              <Share2 size={22} />
            </button>
          </div>
          
          <div className="flex items-center gap-2">
            {currentCustomer?.id === post.customerId && (
              <button 
                className="p-2 text-muted-foreground hover:text-primary hover:bg-muted/50 rounded-full transition-colors" 
                onClick={() => setIsEditModalOpen(true)}
              >
                <Pencil size={18} />
              </button>
            )}
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
      {isEditModalOpen && (
        <EditPostModal 
          post={post} 
          isOpen={isEditModalOpen} 
          onClose={() => setIsEditModalOpen(false)} 
        />
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
