import React from 'react';
import { PostCard } from '../../components/PostCard/PostCard';
import styles from './HomeFeed.module.css';
import { useQuery } from '@tanstack/react-query';
import { apiClient } from '../../api/client';

export const HomeFeed: React.FC = () => {
  const { data, isLoading, error } = useQuery({
    queryKey: ['posts'],
    queryFn: () => apiClient.getPosts(20, 0)
  });

  return (
    <div className={styles.container}>
      <header className={styles.header}>
        <h1 className={styles.title}>Your Feed</h1>
        <p className={styles.subtitle}>Discover what your friends are eating</p>
      </header>

      {isLoading ? (
        <div className={styles.loading}>Loading posts...</div>
      ) : error ? (
        <div className={styles.error}>Error loading posts. Make sure you are logged in.</div>
      ) : (
        <div className={styles.feed}>
          {data?.Posts?.map(post => (
            <PostCard
              key={post.id}
              post={post}
            />
          ))}
          {(!data?.Posts || data.Posts.length === 0) && (
            <p className={styles.empty}>No posts yet.</p>
          )}
        </div>
      )}
    </div>
  );
};
