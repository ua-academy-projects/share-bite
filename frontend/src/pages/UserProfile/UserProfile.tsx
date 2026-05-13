import React, { useState } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import { PostCard } from '../../components/PostCard/PostCard';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { apiClient } from '../../api/client';
import { Button } from '../../components/ui/button';
import { Input } from '../../components/ui/input';
import { Label } from '../../components/ui/label';
import { toast } from 'sonner';
import { Edit3, Grid, Bookmark, Plus } from 'lucide-react';
import { clsx } from 'clsx';

import { useCurrentCustomer } from '../../hooks/useCurrentCustomer';

export const UserProfile: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { data: currentCustomer } = useCurrentCustomer();
  
  // If we have an ID in the URL, it's actually the username.
  // Otherwise, use the current logged-in customer's userName
  const username = id || currentCustomer?.userName;
  const isOwnProfile = !id || (currentCustomer && currentCustomer.userName === id);

  const [newProfile, setNewProfile] = useState({ firstName: '', lastName: '', userName: '', bio: '' });
  const [isEditing, setIsEditing] = useState(false);
  const [editProfile, setEditProfile] = useState({ firstName: '', lastName: '', userName: '', bio: '' });
  const [activeTab, setActiveTab] = useState<'posts' | 'collections'>('posts');
  const [isCreatingCollection, setIsCreatingCollection] = useState(false);
  const [newCollection, setNewCollection] = useState({ name: '', description: '', isPublic: false });

  const createProfileMutation = useMutation({
    mutationFn: (data: typeof newProfile) => apiClient.createCustomer(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['currentCustomer'] });
      queryClient.invalidateQueries({ queryKey: ['user', username] });
      toast.success("Profile created successfully!");
    },
    onError: (error: any) => {
      toast.error(error?.response?.data?.error || "Failed to create profile");
    }
  });

  const handleCreateProfile = (e: React.FormEvent) => {
    e.preventDefault();
    createProfileMutation.mutate(newProfile);
  };

  const updateProfileMutation = useMutation({
    mutationFn: (data: typeof editProfile) => apiClient.updateCustomer(data),
    onSuccess: (updatedUser) => {
      queryClient.invalidateQueries({ queryKey: ['currentCustomer'] });
      queryClient.invalidateQueries({ queryKey: ['user', updatedUser.userName] });
      setIsEditing(false);
      toast.success("Profile updated successfully!");
    },
    onError: (error: any) => {
      toast.error(error?.response?.data?.error || "Failed to update profile");
    }
  });

  const handleUpdateProfile = (e: React.FormEvent) => {
    e.preventDefault();
    updateProfileMutation.mutate(editProfile);
  };

  const handleOpenEdit = () => {
    if (!user) return;
    setEditProfile({
      firstName: user.firstName,
      lastName: user.lastName,
      userName: user.userName,
      bio: user.bio || ''
    });
    setIsEditing(true);
  };

  const createCollectionMutation = useMutation({
    mutationFn: (data: typeof newCollection) => apiClient.createCollection(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['collections'] });
      setIsCreatingCollection(false);
      setNewCollection({ name: '', description: '', isPublic: false });
      toast.success("Collection created successfully!");
    },
    onError: (error: any) => {
      toast.error(error?.response?.data?.error || "Failed to create collection");
    }
  });

  const handleCreateCollection = (e: React.FormEvent) => {
    e.preventDefault();
    createCollectionMutation.mutate(newCollection);
  };

  const { data: user, isLoading: userLoading, isError: userError, error: fetchError } = useQuery({
    queryKey: ['user', username],
    queryFn: () => apiClient.getUser(username!),
    enabled: !!username,
    retry: false,
  });

  // We must use the user.id for the posts endpoint if it requires customer_id, 
  // or we can pass the username if the backend supports it. But according to apiClient.getPosts it takes customerId.
  // We'll wait until we fetch the `user` to get their actual UUID for the posts call.
  const resolvedCustomerId = user?.id;

  const { data: postsData, isLoading: postsLoading } = useQuery({
    queryKey: ['userPosts', resolvedCustomerId],
    queryFn: () => apiClient.getPosts(20, 0, undefined, resolvedCustomerId),
    enabled: !!resolvedCustomerId
  });

  const { data: collections, isLoading: collectionsLoading } = useQuery({
    queryKey: ['collections'],
    queryFn: () => apiClient.getCollections(),
    enabled: isOwnProfile && activeTab === 'collections'
  });

  if (userLoading || (!username && !currentCustomer)) {
    return (
      <div className="flex items-center justify-center min-h-[50vh] text-muted-foreground font-serif text-lg">
        Loading user...
      </div>
    );
  }

  if (userError || !user) {
    const status = (fetchError as any)?.response?.status;
    if (status === 401) {
      localStorage.removeItem('token');
      localStorage.removeItem('refresh_token');
      if (window.location.pathname !== '/auth') {
        window.location.href = '/auth';
      }
      return null;
    }

    if (isOwnProfile) {
      return (
        <div className="flex flex-col items-center justify-center w-full min-h-screen bg-background pt-8 pb-16 px-4">
          <div className="max-w-md w-full bg-card p-10 rounded-3xl border border-border shadow-lg">
            <div className="text-center mb-8">
              <h2 className="text-3xl font-serif font-bold text-foreground mb-2">Create Profile</h2>
              <p className="text-muted-foreground">You don't have a profile yet. Let's set one up!</p>
            </div>
            
            <form onSubmit={handleCreateProfile} className="space-y-5">
              <div className="space-y-2">
                <Label htmlFor="userName">Username</Label>
                <Input 
                  id="userName" 
                  required 
                  placeholder="foodie123" 
                  value={newProfile.userName} 
                  onChange={(e) => setNewProfile({...newProfile, userName: e.target.value})}
                  className="bg-muted/50 border-border h-12 rounded-xl focus-visible:ring-primary"
                />
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="firstName">First Name</Label>
                  <Input 
                    id="firstName" 
                    required 
                    placeholder="John" 
                    value={newProfile.firstName} 
                    onChange={(e) => setNewProfile({...newProfile, firstName: e.target.value})}
                    className="bg-muted/50 border-border h-12 rounded-xl focus-visible:ring-primary"
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="lastName">Last Name</Label>
                  <Input 
                    id="lastName" 
                    required 
                    placeholder="Doe" 
                    value={newProfile.lastName} 
                    onChange={(e) => setNewProfile({...newProfile, lastName: e.target.value})}
                    className="bg-muted/50 border-border h-12 rounded-xl focus-visible:ring-primary"
                  />
                </div>
              </div>
              <div className="space-y-2">
                <Label htmlFor="bio">Bio</Label>
                <Input 
                  id="bio" 
                  placeholder="I love trying new food!" 
                  value={newProfile.bio} 
                  onChange={(e) => setNewProfile({...newProfile, bio: e.target.value})}
                  className="bg-muted/50 border-border h-12 rounded-xl focus-visible:ring-primary"
                />
              </div>
              <Button type="submit" disabled={createProfileMutation.isPending} className="w-full h-12 rounded-xl bg-accent text-accent-foreground hover:bg-[#e6c200] font-bold text-lg shadow-lg hover:shadow-xl transition-all mt-4">
                {createProfileMutation.isPending ? 'Creating...' : 'Create Profile'}
              </Button>
            </form>
          </div>
        </div>
      );
    }

    return (
      <div className="flex flex-col items-center justify-center min-h-[50vh] bg-background pt-8 px-4">
        <div className="max-w-md w-full bg-card p-10 rounded-3xl border border-border shadow-lg text-center">
          <h2 className="text-3xl font-serif font-bold text-foreground mb-4">User Not Found</h2>
          <p className="text-muted-foreground mb-6">The profile you are looking for doesn't exist or has been removed.</p>
          <Button onClick={() => navigate('/')} className="rounded-xl px-6 bg-accent text-accent-foreground hover:bg-[#e6c200] font-bold">
            Back to HomeFeed
          </Button>
        </div>
      </div>
    );
  }

  const userPosts = postsData?.Posts || [];

  return (
    <div className="flex flex-col items-center w-full min-h-screen bg-background pt-8 pb-16 px-4">
      <div className="max-w-2xl w-full">
        {/* Profile Header */}
        <div className="flex flex-col md:flex-row items-center md:items-start gap-8 bg-card p-8 rounded-3xl border border-border/50 shadow-lg dark:shadow-primary/5 mb-10 group relative overflow-hidden">
          {/* Subtle gradient background effect */}
          <div className="absolute inset-0 bg-gradient-to-br from-primary/5 via-transparent to-transparent opacity-50"></div>
          
          <img 
            src={user.avatarURL || 'https://via.placeholder.com/150'} 
            alt={user.userName} 
            className="w-28 h-28 md:w-36 md:h-36 rounded-full object-cover border-4 border-background shadow-xl relative z-10" 
          />
          <div className="flex flex-col items-center md:items-start flex-1 relative z-10">
            <h1 className="text-4xl font-serif font-bold text-foreground tracking-tight mb-1">{user.firstName} {user.lastName}</h1>
            <div className="flex items-center gap-2 mb-4">
              <span className="text-sm font-semibold text-muted-foreground bg-muted/50 px-3 py-1 rounded-full">@{user.userName}</span>
              <span className="text-[10px] font-black tracking-widest text-[#98FF98] bg-[#1A3C34] px-2 py-1 rounded-full uppercase">Guest</span>
            </div>
            
            {user.bio && (
              <p className="text-foreground/80 text-sm mb-6 max-w-md text-center md:text-left leading-relaxed">
                {user.bio}
              </p>
            )}

            <div className="flex gap-8 w-full justify-center md:justify-start mt-2">
              <div className="flex flex-col items-center group/stat">
                <span className="text-2xl font-bold text-foreground group-hover/stat:text-primary transition-colors">{userPosts.length}</span>
                <span className="text-xs font-semibold text-muted-foreground uppercase tracking-wider">Posts</span>
              </div>
              <div className="flex flex-col items-center group/stat">
                <span className="text-2xl font-bold text-foreground group-hover/stat:text-primary transition-colors">{user.followers || 0}</span>
                <span className="text-xs font-semibold text-muted-foreground uppercase tracking-wider">Followers</span>
              </div>
              <div className="flex flex-col items-center group/stat">
                <span className="text-2xl font-bold text-foreground group-hover/stat:text-primary transition-colors">{user.following || 0}</span>
                <span className="text-xs font-semibold text-muted-foreground uppercase tracking-wider">Following</span>
              </div>
            </div>
          </div>
          
          {isOwnProfile && (
            <div className="absolute top-8 right-8 z-30">
              <Button onClick={handleOpenEdit} variant="outline" className="rounded-xl border-border bg-background/50 backdrop-blur-sm hover:bg-muted font-bold shadow-sm cursor-pointer pointer-events-auto">
                <Edit3 size={16} className="mr-2" />
                Edit Profile
              </Button>
            </div>
          )}
        </div>

        {/* Tabs System */}
        <div className="flex items-center gap-6 border-b border-border/50 mb-8">
          <button 
            onClick={() => setActiveTab('posts')}
            className={clsx(
              "flex items-center gap-2 pb-4 text-sm font-bold uppercase tracking-wider transition-all border-b-2",
              activeTab === 'posts' ? "text-foreground border-primary" : "text-muted-foreground border-transparent hover:text-foreground"
            )}
          >
            <Grid size={18} />
            Recent Stories
          </button>
          
          {isOwnProfile && (
            <button 
              onClick={() => setActiveTab('collections')}
              className={clsx(
                "flex items-center gap-2 pb-4 text-sm font-bold uppercase tracking-wider transition-all border-b-2",
                activeTab === 'collections' ? "text-foreground border-primary" : "text-muted-foreground border-transparent hover:text-foreground"
              )}
            >
              <Bookmark size={18} />
              Saved Collections
            </button>
          )}
        </div>

        {/* Tab Content */}
        {activeTab === 'posts' ? (
          <div className="flex flex-col gap-6">
            {postsLoading ? (
              <p className="text-center text-muted-foreground py-8">Loading posts...</p>
            ) : (
              userPosts.length > 0 ? (
                userPosts.map(post => (
                  <PostCard 
                    key={post.id} 
                    post={post} 
                  />
                ))
              ) : (
                <div className="text-center py-12 text-muted-foreground bg-card rounded-2xl border border-border">
                  <p>This user hasn't posted anything yet.</p>
                </div>
              )
            )}
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {collectionsLoading ? (
              <p className="text-center text-muted-foreground py-8 col-span-2">Loading collections...</p>
            ) : (
              collections?.length ? (
                collections.map(collection => (
                  <Link to={`/collections/${collection.id}`} key={collection.id} className="group flex flex-col bg-card border border-border/50 rounded-3xl p-6 shadow-lg hover:-translate-y-1 transition-all overflow-hidden relative">
                    <div className="absolute inset-0 bg-gradient-to-br from-[#98FF98]/5 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition-opacity"></div>
                    <div className="flex justify-between items-start mb-4 relative z-10">
                      <div className="p-3 bg-[#1A3C34] text-[#98FF98] rounded-2xl border border-[#98FF98]/20 group-hover:bg-[#98FF98] group-hover:text-[#0b0f0e] transition-colors">
                        <Bookmark size={24} className="fill-current" />
                      </div>
                      <span className="text-xs font-bold uppercase tracking-widest text-muted-foreground group-hover:text-foreground transition-colors">
                        {collection.isPublic ? 'Public' : 'Private'}
                      </span>
                    </div>
                    <h3 className="text-2xl font-serif font-bold text-foreground mb-2 group-hover:text-primary transition-colors relative z-10">{collection.name}</h3>
                    {collection.description && <p className="text-muted-foreground text-sm line-clamp-2 mb-4 relative z-10">{collection.description}</p>}
                    <div className="mt-auto pt-4 border-t border-border/50 flex justify-between items-center relative z-10">
                      <span className="text-xs font-semibold text-muted-foreground">Updated {new Date(collection.updatedAt).toLocaleDateString()}</span>
                      <span className="text-xs font-bold text-primary group-hover:underline">View Collection &rarr;</span>
                    </div>
                  </Link>
                ))
              ) : (
                <div className="col-span-2 text-center py-12 bg-card rounded-3xl border border-border flex flex-col items-center">
                  <Bookmark className="mx-auto mb-4 text-muted-foreground opacity-30" size={48} />
                  <h3 className="text-2xl font-serif font-bold mb-2">No Collections Yet</h3>
                  <p className="text-muted-foreground mb-6">Create your first collection to start saving your favorite venues.</p>
                  <Button onClick={() => setIsCreatingCollection(true)} className="rounded-xl px-6 bg-accent text-accent-foreground hover:bg-[#e6c200] font-bold cursor-pointer">
                    <Plus size={18} className="mr-2" /> New Collection
                  </Button>
                </div>
              )
            )}
          </div>
        )}
      </div>

      {/* Edit Profile Modal */}
      {isEditing && (
        <div className="fixed inset-0 z-[100] flex items-center justify-center bg-black/60 backdrop-blur-sm" onClick={() => setIsEditing(false)}>
          <div className="bg-card text-card-foreground border border-border p-8 rounded-3xl max-w-md w-full shadow-2xl animate-in fade-in zoom-in-95" onClick={(e) => e.stopPropagation()}>
            <div className="flex justify-between items-center mb-6">
              <h3 className="text-2xl font-serif font-bold">Edit Profile</h3>
              <button onClick={() => setIsEditing(false)} className="text-muted-foreground hover:text-foreground transition-colors p-2 rounded-full hover:bg-muted/50">
                &times;
              </button>
            </div>
            
            <form onSubmit={handleUpdateProfile} className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="editUserName">Username</Label>
                <Input 
                  id="editUserName" 
                  required 
                  value={editProfile.userName} 
                  onChange={(e) => setEditProfile({...editProfile, userName: e.target.value})}
                  className="bg-muted/50 border-border h-11 rounded-xl focus-visible:ring-primary"
                />
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="editFirstName">First Name</Label>
                  <Input 
                    id="editFirstName" 
                    required 
                    value={editProfile.firstName} 
                    onChange={(e) => setEditProfile({...editProfile, firstName: e.target.value})}
                    className="bg-muted/50 border-border h-11 rounded-xl focus-visible:ring-primary"
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="editLastName">Last Name</Label>
                  <Input 
                    id="editLastName" 
                    required 
                    value={editProfile.lastName} 
                    onChange={(e) => setEditProfile({...editProfile, lastName: e.target.value})}
                    className="bg-muted/50 border-border h-11 rounded-xl focus-visible:ring-primary"
                  />
                </div>
              </div>
              <div className="space-y-2">
                <Label htmlFor="editBio">Bio</Label>
                <Input 
                  id="editBio" 
                  value={editProfile.bio} 
                  onChange={(e) => setEditProfile({...editProfile, bio: e.target.value})}
                  className="bg-muted/50 border-border h-11 rounded-xl focus-visible:ring-primary"
                />
              </div>
              
              <div className="flex justify-end gap-3 pt-6">
                <Button type="button" variant="ghost" onClick={() => setIsEditing(false)} className="h-11 px-6 rounded-full font-bold cursor-pointer hover:bg-muted/50">
                  Cancel
                </Button>
                <Button type="submit" disabled={updateProfileMutation.isPending} className="bg-accent text-accent-foreground rounded-full h-11 px-6 font-bold shadow-lg cursor-pointer hover:bg-[#e6c200] transition-colors">
                  {updateProfileMutation.isPending ? 'Saving...' : 'Save Changes'}
                </Button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* New Collection Modal */}
      {isCreatingCollection && (
        <div className="fixed inset-0 z-[100] flex items-center justify-center bg-black/60 backdrop-blur-sm" onClick={() => setIsCreatingCollection(false)}>
          <div className="bg-card text-card-foreground border border-border p-8 rounded-3xl max-w-md w-full shadow-2xl animate-in fade-in zoom-in-95" onClick={(e) => e.stopPropagation()}>
            <div className="flex justify-between items-center mb-6">
              <h3 className="text-2xl font-serif font-bold">New Collection</h3>
              <button onClick={() => setIsCreatingCollection(false)} className="text-muted-foreground hover:text-foreground transition-colors p-2 rounded-full hover:bg-muted/50 cursor-pointer">
                &times;
              </button>
            </div>
            
            <form onSubmit={handleCreateCollection} className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="collectionName">Collection Name</Label>
                <Input 
                  id="collectionName" 
                  required 
                  placeholder="e.g. Best Coffee Shops" 
                  value={newCollection.name} 
                  onChange={(e) => setNewCollection({...newCollection, name: e.target.value})}
                  className="bg-muted/50 border-border h-11 rounded-xl focus-visible:ring-primary"
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="collectionDesc">Description (Optional)</Label>
                <Input 
                  id="collectionDesc" 
                  placeholder="Places to visit this summer" 
                  value={newCollection.description} 
                  onChange={(e) => setNewCollection({...newCollection, description: e.target.value})}
                  className="bg-muted/50 border-border h-11 rounded-xl focus-visible:ring-primary"
                />
              </div>
              <div className="flex items-center gap-3 mt-4 mb-2">
                <input 
                  type="checkbox" 
                  id="isPublic" 
                  checked={newCollection.isPublic}
                  onChange={(e) => setNewCollection({...newCollection, isPublic: e.target.checked})}
                  className="w-4 h-4 rounded border-border bg-muted/50 accent-primary"
                />
                <Label htmlFor="isPublic" className="font-bold">Make Collection Public</Label>
              </div>
              
              <div className="flex justify-end gap-3 pt-6">
                <Button type="button" variant="ghost" onClick={() => setIsCreatingCollection(false)} className="h-11 px-6 rounded-full font-bold cursor-pointer hover:bg-muted/50">
                  Cancel
                </Button>
                <Button type="submit" disabled={createCollectionMutation.isPending} className="bg-accent text-accent-foreground rounded-full h-11 px-6 font-bold shadow-lg cursor-pointer hover:bg-[#e6c200] transition-colors">
                  {createCollectionMutation.isPending ? 'Creating...' : 'Create Collection'}
                </Button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};
