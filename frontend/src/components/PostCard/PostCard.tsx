import React, { useState } from 'react';
import { Heart, MessageCircle, Share2, Bookmark } from 'lucide-react';
import type { PostResponse } from '../../types/api';
import styles from './PostCard.module.css';
import { clsx } from 'clsx';
import { Link } from 'react-router-dom';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { apiClient } from '../../api/client';

interface PostCardProps {
  post: PostResponse;
  restaurantName?: string;
}

export const PostCard: React.FC<PostCardProps> = ({ post, restaurantName }) => {
  const queryClient = useQueryClient();
  const [liked, setLiked] = useState(post.isLikedByMe);
  const [saved, setSaved] = useState(false);

  const toggleLikeMutation = useMutation({
    mutationFn: async () => {
      if (liked) {
        await apiClient.unlikePost(post.id);
      } else {
        await apiClient.likePost(post.id);
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['posts'] });
    }
  });

  const handleLike = () => {
    setLiked(!liked);
    toggleLikeMutation.mutate();
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

      {/* Image */}
      {post.images && post.images.length > 0 && (
        <div className={styles.imageContainer}>
          <img src={post.images[0]} alt="Food" className={styles.image} />
        </div>
      )}

      {/* Actions */}
      <div className={styles.actions}>
        <div className={styles.leftActions}>
          <button 
            className={clsx(styles.actionBtn, liked && styles.liked)} 
            onClick={handleLike}
            disabled={toggleLikeMutation.isPending}
          >
            <Heart size={24} fill={liked ? "currentColor" : "none"} />
          </button>
          <button className={styles.actionBtn}>
            <MessageCircle size={24} />
          </button>
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
        <div className={styles.likes}>{post.likesCount + (liked && !post.isLikedByMe ? 1 : (!liked && post.isLikedByMe ? -1 : 0))} likes</div>
        <div className={styles.text}>
          <span className={styles.authorNameInline}>{post.userName}</span> {post.text}
        </div>
      </div>
    </div>
  );
};
