import React from 'react';
import { useParams } from 'react-router-dom';
import { PostCard } from '../../components/PostCard/PostCard';
import styles from './UserProfile.module.css';
import { useQuery } from '@tanstack/react-query';
import { apiClient } from '../../api/client';

export const UserProfile: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  // For the current user, we might get ID from auth token or have a dedicated endpoint
  // Defaulting to 1 for demo purposes if no ID
  const userId = id || '1'; 

  const { data: user, isLoading: userLoading } = useQuery({
    queryKey: ['user', userId],
    queryFn: () => apiClient.getUser(userId),
    enabled: !!userId
  });

  const { data: postsData, isLoading: postsLoading } = useQuery({
    queryKey: ['userPosts', userId],
    queryFn: () => apiClient.getPosts(20, 0), // Mock logic: filter locally since specific endpoint is not mapped yet
    enabled: !!userId
  });

  if (userLoading) return <div className={styles.notFound}>Loading user...</div>;
  if (!user) return <div className={styles.notFound}>User not found</div>;

  const userPosts = postsData?.Posts?.filter(p => p.customerId === userId) || [];

  return (
    <div className={styles.container}>
      {/* Profile Header */}
      <div className={styles.header}>
        <img src={user.avatar || 'https://via.placeholder.com/150'} alt={user.name} className={styles.avatar} />
        <div className={styles.info}>
          <h1 className={styles.name}>{user.name || 'User'}</h1>
          <p className={styles.handle}>@{user.id}</p>
          <div className={styles.stats}>
            <div className={styles.stat}>
              <span className={styles.statValue}>{userPosts.length}</span>
              <span className={styles.statLabel}>Posts</span>
            </div>
            <div className={styles.stat}>
              <span className={styles.statValue}>1.2k</span>
              <span className={styles.statLabel}>Followers</span>
            </div>
            <div className={styles.stat}>
              <span className={styles.statValue}>340</span>
              <span className={styles.statLabel}>Following</span>
            </div>
          </div>
        </div>
      </div>

      <div className={styles.divider}></div>

      {/* User's Posts Feed */}
      <h2 className={styles.feedTitle}>Recent Posts</h2>
      <div className={styles.feed}>
        {postsLoading ? <p>Loading posts...</p> : (
          userPosts.length > 0 ? (
            userPosts.map(post => (
              <PostCard 
                key={post.id} 
                post={post} 
              />
            ))
          ) : (
            <p className={styles.emptyState}>This user hasn't posted anything yet.</p>
          )
        )}
      </div>
    </div>
  );
};
