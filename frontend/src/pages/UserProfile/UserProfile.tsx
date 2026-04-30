import React, { useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import { PostCard } from '../../components/PostCard/PostCard';
import styles from './UserProfile.module.css';
import { useQuery } from '@tanstack/react-query';
import { apiClient } from '../../api/client';
import { useCurrentCustomer } from '../../hooks/useCurrentCustomer';
import { EditProfileModal } from './components/EditProfileModal';
import { Settings, Grid, Bookmark, FolderHeart } from 'lucide-react';
import { clsx } from 'clsx';

export const UserProfile: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const { data: currentCustomer, isLoading: currentCustomerLoading } = useCurrentCustomer();
  const [activeTab, setActiveTab] = useState<'posts' | 'collections'>('posts');
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);

  // If no ID, we are viewing "me"
  const isMe = !id || id === currentCustomer?.id;
  const username = id || currentCustomer?.userName;

  const { data: profile, isLoading: profileLoading, error: profileError } = useQuery({
    queryKey: ['user', username],
    queryFn: () => apiClient.getCustomerByUsername(username!),
    enabled: !!username,
    retry: (failureCount, error: any) => {
      // Don't retry on 404 as it's a valid state for new users
      if (error?.response?.status === 404) return false;
      return failureCount < 3;
    }
  });

  const isProfileNotFound = (profileError as any)?.response?.status === 404;
  const showOnboarding = isMe && (isProfileNotFound || (!username && !currentCustomerLoading));

  const { data: postsData, isLoading: postsLoading } = useQuery({
    queryKey: ['userPosts', profile?.id],
    queryFn: () => apiClient.getPosts(20, 0, undefined, profile?.id),
    enabled: !!profile?.id && activeTab === 'posts'
  });

  const { data: collectionsData, isLoading: collectionsLoading } = useQuery({
    queryKey: ['myCollections'],
    queryFn: () => apiClient.listMyCollections(),
    enabled: isMe && activeTab === 'collections'
  });

  if (profileLoading || currentCustomerLoading) return <div className={styles.loading}>Loading profile...</div>;

  if (showOnboarding) {
    return (
      <div className={styles.container}>
        <div className={clsx(styles.onboardingCard, 'glass-panel')}>
          <div className={styles.onboardingIcon}>
            <FolderHeart size={64} />
          </div>
          <h1 className={styles.onboardingTitle}>Welcome to ShareBite!</h1>
          <p className={styles.onboardingText}>
            It looks like you haven't set up your profile yet. 
            Let's get started so you can share your food experiences!
          </p>
          <button 
            className={styles.setupBtn}
            onClick={() => setIsEditModalOpen(true)}
          >
            Setup Profile
          </button>
        </div>

        {isEditModalOpen && (
          <EditProfileModal 
            profile={profile || {
              id: '',
              userId: '',
              userName: '',
              firstName: '',
              lastName: '',
              avatarUrl: null,
              bio: '',
              createdAt: new Date().toISOString()
            }}
            isOpen={isEditModalOpen}
            onClose={() => setIsEditModalOpen(false)}
          />
        )}
      </div>
    );
  }

  if (!profile) return <div className={styles.notFound}>User not found</div>;

  const userPosts = postsData?.Posts || [];
  const myCollections = collectionsData?.collections || [];

  return (
    <div className={styles.container}>
      {/* Profile Header */}
      <div className={clsx(styles.header, 'glass-panel')}>
        <div className={styles.headerTop}>
          <img 
            src={profile.avatarUrl || 'https://via.placeholder.com/150'} 
            alt={profile.userName} 
            className={styles.avatar} 
          />
          <div className={styles.info}>
            <div className={styles.nameRow}>
              <h1 className={styles.name}>{profile.firstName} {profile.lastName}</h1>
              {isMe && (
                <button 
                  className={styles.editBtn}
                  onClick={() => setIsEditModalOpen(true)}
                >
                  <Settings size={20} />
                  <span>Edit Profile</span>
                </button>
              )}
            </div>
            <p className={styles.handle}>@{profile.userName}</p>
            {profile.bio && <p className={styles.bio}>{profile.bio}</p>}
            
            <div className={styles.stats}>
              <div className={styles.stat}>
                <span className={styles.statValue}>{userPosts.length}</span>
                <span className={styles.statLabel}>Posts</span>
              </div>
              <div className={styles.stat}>
                <span className={styles.statValue}>--</span>
                <span className={styles.statLabel}>Followers</span>
              </div>
              <div className={styles.stat}>
                <span className={styles.statValue}>--</span>
                <span className={styles.statLabel}>Following</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Tabs */}
      <div className={styles.tabs}>
        <button 
          className={clsx(styles.tab, activeTab === 'posts' && styles.activeTab)}
          onClick={() => setActiveTab('posts')}
        >
          <Grid size={20} />
          <span>Posts</span>
        </button>
        {isMe && (
          <button 
            className={clsx(styles.tab, activeTab === 'collections' && styles.activeTab)}
            onClick={() => setActiveTab('collections')}
          >
            <FolderHeart size={20} />
            <span>Collections</span>
          </button>
        )}
      </div>

      {/* Tab Content */}
      <div className={styles.content}>
        {activeTab === 'posts' ? (
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
                <div className={styles.emptyState}>
                  <Grid size={48} className={styles.emptyIcon} />
                  <p>No posts yet.</p>
                </div>
              )
            )}
          </div>
        ) : (
          <div className={styles.collectionsGrid}>
            {collectionsLoading ? <p>Loading collections...</p> : (
              myCollections.length > 0 ? (
                myCollections.map(collection => (
                  <Link 
                    key={collection.id} 
                    to={`/collections/${collection.id}`}
                    className={clsx(styles.collectionCard, 'glass-panel')}
                  >
                    <div className={styles.collectionIcon}>
                      <Bookmark size={32} />
                    </div>
                    <div className={styles.collectionInfo}>
                      <h3 className={styles.collectionName}>{collection.name}</h3>
                      <p className={styles.collectionDesc}>
                        {collection.description || 'No description'}
                      </p>
                      <span className={styles.collectionTag}>
                        {collection.isPublic ? 'Public' : 'Private'}
                      </span>
                    </div>
                  </Link>
                ))
              ) : (
                <div className={styles.emptyState}>
                  <FolderHeart size={48} className={styles.emptyIcon} />
                  <p>No collections created yet.</p>
                </div>
              )
            )}
          </div>
        )}
      </div>

      {isEditModalOpen && (
        <EditProfileModal 
          profile={profile}
          isOpen={isEditModalOpen}
          onClose={() => setIsEditModalOpen(false)}
        />
      )}
    </div>
  );
};
