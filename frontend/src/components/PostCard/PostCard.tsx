import React, { useState } from 'react';
import { Heart, MessageCircle, Share2, Bookmark, Trash2, Edit2, ChevronLeft, ChevronRight } from 'lucide-react';
import type { PostResponse } from '../../types/api';
import styles from './PostCard.module.css';
import { clsx } from 'clsx';
import { Link } from 'react-router-dom';
import { useMutation, useQueryClient, useQuery } from '@tanstack/react-query';
import { apiClient } from '../../api/client';
import { useCurrentCustomer } from '../../hooks/useCurrentCustomer';

interface PostCardProps {
  post: PostResponse;
  restaurantName?: string;
}

export const PostCard: React.FC<PostCardProps> = ({ post, restaurantName }) => {
  const queryClient = useQueryClient();
  const { data: currentCustomer } = useCurrentCustomer();
  const [liked, setLiked] = useState(post.isLikedByMe);
  const [saved, setSaved] = useState(false);
  const [showComments, setShowComments] = useState(false);
  const [commentText, setCommentText] = useState('');
  const [currentImageIndex, setCurrentImageIndex] = useState(0);
  const [editingCommentId, setEditingCommentId] = useState<number | null>(null);
  const [editCommentText, setEditCommentText] = useState('');

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
    }
  });

  const updateCommentMutation = useMutation({
    mutationFn: async ({ commentId, text }: { commentId: number; text: string }) => {
      return await apiClient.updateComment(post.id, commentId, text);
    },
    onSuccess: () => {
      setEditingCommentId(null);
      setEditCommentText('');
      queryClient.invalidateQueries({ queryKey: ['comments', post.id] });
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
      const previousLiked = liked;
      setLiked(nextLiked);
      return { previousLiked };
    },
    onError: (err, nextLiked, context) => {
      if (context?.previousLiked !== undefined) {
        setLiked(context.previousLiked);
      }
    },
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: ['posts'] });
    }
  });

  const handleLike = () => {
    const nextLiked = !liked;
    toggleLikeMutation.mutate(nextLiked);
  };

  return (
    <div className={clsx(styles.card, 'glass-panel')}>
      {/* Header */}
      <div className={styles.header}>
        <div className={styles.authorInfo}>
          <img src={post.avatarURL || 'https://via.placeholder.com/40'} alt={post.userName} className={styles.avatar} />
          <div>
            <Link to={`/user/${post.customerId}`} className={styles.authorName}>{post.userName}</Link>
            <div className={styles.meta}>
              <span>{new Date(post.createdAt).toLocaleDateString()}</span>
            </div>
          </div>
        </div>
      </div>

      {/* Restaurant Tag */}
      <div className={styles.restaurantTag}>
        <Link to={`/restaurant/${post.venueId}`} className={styles.restaurantLink}>
          📍 {restaurantName || `Venue #${post.venueId}`}
        </Link>
        <span className={styles.rating}>★ {post.rating}</span>
      </div>

      {/* Image Carousel */}
      {post.images && post.images.length > 0 && (
        <div className={styles.imageContainer}>
          <img src={post.images[currentImageIndex]} alt="Food" className={styles.image} />
          
          {post.images.length > 1 && (
            <>
              {currentImageIndex > 0 && (
                <button 
                  className={clsx(styles.carouselBtn, styles.carouselPrev)} 
                  onClick={handlePrevImage}
                >
                  <ChevronLeft size={24} />
                </button>
              )}
              {currentImageIndex < post.images.length - 1 && (
                <button 
                  className={clsx(styles.carouselBtn, styles.carouselNext)} 
                  onClick={handleNextImage}
                >
                  <ChevronRight size={24} />
                </button>
              )}
              
              <div className={styles.carouselDots}>
                {post.images.map((_, idx) => (
                  <span 
                    key={idx} 
                    className={clsx(styles.dot, idx === currentImageIndex && styles.activeDot)}
                    onClick={() => setCurrentImageIndex(idx)}
                  />
                ))}
              </div>
            </>
          )}
        </div>
      )}

      {/* Actions */}
      <div className={styles.actions}>
        <div className={styles.leftActions}>
          <div className={styles.actionGroup}>
            <button 
              className={clsx(styles.actionBtn, liked && styles.liked)} 
              onClick={handleLike}
              disabled={toggleLikeMutation.isPending}
            >
              <Heart size={24} fill={liked ? "currentColor" : "none"} />
            </button>
            <span className={styles.actionCount}>
              {post.likesCount + (liked && !post.isLikedByMe ? 1 : (!liked && post.isLikedByMe ? -1 : 0))}
            </span>
          </div>
          
          <div className={styles.actionGroup}>
            <button 
              className={clsx(styles.actionBtn, showComments && styles.active)}
              onClick={() => setShowComments(!showComments)}
            >
              <MessageCircle size={24} />
            </button>
            <span className={styles.actionCount}>
              {commentsData ? commentsData.total : '...'}
            </span>
          </div>

          <button className={styles.actionBtn}>
            <Share2 size={24} />
          </button>
        </div>
        <button 
          className={clsx(styles.actionBtn, saved && styles.saved)}
          onClick={() => setSaved(!saved)}
        >
          <Bookmark size={24} fill={saved ? "currentColor" : "none"} />
        </button>
      </div>

      {/* Content */}
      <div className={styles.content}>
        <div className={styles.text}>
          <span className={styles.authorNameInline}>{post.userName}</span> {post.text}
        </div>
      </div>

      {/* Comments Section */}
      {showComments && (
        <div className={styles.commentsSection}>
          <div className={styles.commentsList}>
            {isLoadingComments ? (
              <div className={styles.commentsLoading}>Loading comments...</div>
            ) : commentsData?.entities.length === 0 ? (
              <div className={styles.noComments}>No comments yet. Be the first!</div>
            ) : (
              commentsData?.entities.map(comment => (
                <div key={comment.id} className={styles.commentItem}>
                  <div className={styles.commentHeader}>
                    <img 
                      src={comment.customer.avatarURL || 'https://via.placeholder.com/24'} 
                      alt={comment.customer.userName} 
                      className={styles.commentAvatar} 
                    />
                    <span className={styles.commentAuthor}>{comment.customer.userName}</span>
                    <span className={styles.commentDate}>{new Date(comment.createdAt).toLocaleDateString()}</span>
                  </div>
                  <div className={styles.commentBody}>
                    {editingCommentId === comment.id ? (
                      <div className={styles.inlineEditForm}>
                        <input 
                          type="text" 
                          className={styles.commentInput} 
                          value={editCommentText}
                          onChange={e => setEditCommentText(e.target.value)}
                          autoFocus
                        />
                        <div className={styles.inlineEditActions}>
                          <button 
                            className={styles.commentSubmitBtn}
                            onClick={() => updateCommentMutation.mutate({ commentId: comment.id, text: editCommentText })}
                            disabled={!editCommentText.trim() || updateCommentMutation.isPending}
                          >
                            Save
                          </button>
                          <button 
                            className={styles.cancelEditBtn}
                            onClick={() => { setEditingCommentId(null); setEditCommentText(''); }}
                          >
                            Cancel
                          </button>
                        </div>
                      </div>
                    ) : (
                      <>
                        <p className={styles.commentText}>{comment.text}</p>
                        {currentCustomer?.id === comment.customer.id && (
                          <div className={styles.commentActions}>
                            <button 
                              className={styles.editCommentBtn}
                              onClick={() => {
                                setEditingCommentId(comment.id);
                                setEditCommentText(comment.text);
                              }}
                              title="Edit comment"
                            >
                              <Edit2 size={14} />
                            </button>
                            <button 
                              className={styles.deleteCommentBtn}
                              onClick={() => deleteCommentMutation.mutate(comment.id)}
                              disabled={deleteCommentMutation.isPending}
                              title="Delete comment"
                            >
                              <Trash2 size={14} />
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
            className={styles.commentForm} 
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
              className={styles.commentInput}
              value={commentText}
              onChange={e => setCommentText(e.target.value)}
              disabled={createCommentMutation.isPending}
            />
            <button 
              type="submit" 
              className={styles.commentSubmitBtn}
              disabled={!commentText.trim() || createCommentMutation.isPending}
            >
              Post
            </button>
          </form>
          {createCommentMutation.isError && (
            <div className={styles.commentError}>Failed to post comment</div>
          )}
        </div>
      )}
    </div>
  );
};
