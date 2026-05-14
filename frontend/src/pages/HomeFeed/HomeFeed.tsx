import React from 'react';
import { PostCard } from '../../components/PostCard/PostCard';
import { useQuery } from '@tanstack/react-query';
import { apiClient } from '../../api/client';

export const HomeFeed: React.FC = () => {
  const { data, isLoading, error } = useQuery({
    queryKey: ['posts'],
    queryFn: () => apiClient.getPosts(20, 0)
  });

  return (
    <div className="flex flex-col items-center w-full min-h-screen bg-background pt-8 pb-16 px-4">
      <header className="max-w-2xl w-full mb-10 text-center">
        <h1 className="text-5xl font-serif font-bold tracking-tight text-foreground mb-3">The Home Feed</h1>
        <p className="text-muted-foreground text-lg font-medium">Discover what your friends are eating</p>
      </header>

      {isLoading ? (
        <div className="flex items-center justify-center py-12 text-muted-foreground">Loading posts...</div>
      ) : error ? (
        <div className="flex items-center justify-center py-12 text-destructive bg-destructive/10 rounded-lg p-4 max-w-2xl w-full">Error loading posts. Make sure you are logged in.</div>
      ) : (
        <div className="flex flex-col gap-8 max-w-2xl w-full">
          {data?.Posts?.map(post => (
            <PostCard
              key={post.id}
              post={post}
            />
          ))}
          {(!data?.Posts || data.Posts.length === 0) && (
            <p className="text-center text-muted-foreground py-12">No posts yet.</p>
          )}
        </div>
      )}
    </div>
  );
};
