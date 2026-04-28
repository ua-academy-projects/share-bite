import React from 'react';
import { useParams } from 'react-router-dom';
import { ReviewCard } from '../../components/ReviewCard/ReviewCard';
import { PostCard } from '../../components/PostCard/PostCard';
import styles from './RestaurantProfile.module.css';
import { MapPin, Star } from 'lucide-react';
import { clsx } from 'clsx';
import { useQuery } from '@tanstack/react-query';
import { apiClient } from '../../api/client';
import axios from 'axios';

export const RestaurantProfile: React.FC = () => {
  const { id } = useParams<{ id: string }>();

  const { data: restaurant, isLoading: resLoading, error: resError } = useQuery({
    queryKey: ['restaurant', id],
    queryFn: async () => {
      const res = await axios.get(`/api/business/org-units/${id}`);
      return res.data;
    },
    enabled: !!id,
    retry: false
  });

  const { data: postsData, isLoading: postsLoading } = useQuery({
    queryKey: ['restaurantPosts', id],
    queryFn: () => apiClient.getPosts(20, 0), // Mocking endpoint logic for now as specific venue post endpoint not defined
    enabled: !!id
  });

  if (resLoading) return <div className={styles.notFound}>Loading venue...</div>;
  
  if (resError || !restaurant) {
    const errorMsg = (resError as any)?.response?.data?.error || "This location could not be found.";
    return (
      <div className={styles.notFoundContainer}>
        <div className={clsx(styles.notFoundCard, 'glass-panel')}>
          <h2>Oops! Restaurant not found</h2>
          <p>{errorMsg}</p>
          <button 
            onClick={() => window.location.href = '/explore'} 
            className={styles.backButton}
          >
            Back to Explore
          </button>
        </div>
      </div>
    );
  }

  const restaurantPosts = postsData?.Posts?.filter(p => p.venueId === Number(id)) || [];

  return (
    <div className={styles.container}>
      {/* Hero Section */}
      <div className={clsx(styles.hero, 'glass-panel')}>
        <div className={styles.imageContainer}>
          <img src={restaurant.image || 'https://images.unsplash.com/photo-1514933651103-005eec06c04b'} alt={restaurant.name} className={styles.image} />
        </div>
        <div className={styles.heroContent}>
          <div className={styles.headerRow}>
            <h1 className={styles.title}>{restaurant.name}</h1>
            <div className={styles.ratingBadge}>
              <Star size={20} fill="currentColor" />
              <span>{restaurant.rating || '--'}</span>
            </div>
          </div>
          <p className={styles.category}>{restaurant.category || 'Restaurant'}</p>
          <p className={styles.description}>{restaurant.description}</p>
          <div className={styles.location}>
            <MapPin size={18} />
            <span>{restaurant.location || 'Location TBA'}</span>
          </div>
        </div>
      </div>

      <div className={styles.contentGrid}>
        {/* Right Column - Posts */}
        <div className={styles.column}>
          <h2 className={styles.sectionTitle}>Recent Posts</h2>
          <div className={styles.postsList}>
            {postsLoading ? <p>Loading posts...</p> : (
              restaurantPosts.length > 0 ? (
                restaurantPosts.map(post => (
                  <PostCard 
                    key={post.id} 
                    post={post} 
                    restaurantName={restaurant.name} 
                  />
                ))
              ) : (
                <p className={styles.emptyState}>No posts yet.</p>
              )
            )}
          </div>
        </div>
      </div>
    </div>
  );
};
