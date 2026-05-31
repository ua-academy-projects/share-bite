import { useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { Bookmark, Edit3, Grid, Loader2, Plus } from "lucide-react";
import { apiClient } from "@/api/client";
import { useCurrentCustomer } from "@/hooks/useCurrentCustomer";
import { GuestPostCard } from "@/components/PostCard/PostCard";
import { PageHeader } from "@/components/layout/PageHeader";
import { PageLayout } from "@/components/layout/PageLayout";
import {
  pageBtnPrimary,
  pageBtnSecondary,
  pageEmpty,
  pageInput,
  pageLabel,
  pageLoader,
  pagePanel,
} from "@/components/layout/pageStyles";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { cn } from "@/lib/utils";

export function UserProfile() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { data: currentCustomer } = useCurrentCustomer();

  const username = id || currentCustomer?.userName;
  const isOwnProfile = !id || currentCustomer?.userName === id;

  const [activeTab, setActiveTab] = useState<"posts" | "collections">("posts");
  const [isEditing, setIsEditing] = useState(false);
  const [isCreatingCollection, setIsCreatingCollection] = useState(false);
  const [editProfile, setEditProfile] = useState({
    firstName: "",
    lastName: "",
    userName: "",
    bio: "",
  });
  const [newCollection, setNewCollection] = useState({
    name: "",
    description: "",
    visibility: "private",
  });

  const {
    data: user,
    isLoading: userLoading,
    isError: userError,
    error: fetchError,
  } = useQuery({
    queryKey: ["user", username],
    queryFn: () => apiClient.getUser(username!),
    enabled: !!username,
    retry: false,
  });

  const resolvedCustomerId = user?.id;

  const { data: postsData, isLoading: postsLoading } = useQuery({
    queryKey: ["userPosts", resolvedCustomerId],
    queryFn: () => apiClient.getPosts(100, 0, undefined, resolvedCustomerId),
    enabled: !!resolvedCustomerId,
  });

  const { data: collections, isLoading: collectionsLoading } = useQuery({
    queryKey: ["collections", "mine"],
    queryFn: () => apiClient.getCollections(),
    enabled: activeTab === "collections" && isOwnProfile,
  });

  const updateProfileMutation = useMutation({
    mutationFn: (data: typeof editProfile) => apiClient.updateCustomer(data),
    onSuccess: (updatedUser) => {
      queryClient.invalidateQueries({ queryKey: ["currentCustomer"] });
      queryClient.invalidateQueries({ queryKey: ["user", updatedUser.userName] });
      setIsEditing(false);
      toast.success("Profile updated successfully!");
    },
    onError: (error: any) => {
      toast.error(error?.response?.data?.error || "Failed to update profile");
    },
  });

  const createCollectionMutation = useMutation({
    mutationFn: (data: { name: string; description: string; isPublic: boolean }) =>
      apiClient.createCollection(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["collections", "mine"] });
      setIsCreatingCollection(false);
      setNewCollection({ name: "", description: "", visibility: "private" });
      toast.success("Collection created successfully!");
    },
    onError: (error: any) => {
      toast.error(error?.response?.data?.error || "Failed to create collection");
    },
  });

  const handleOpenEdit = () => {
    if (!user) return;
    setEditProfile({
      firstName: user.firstName,
      lastName: user.lastName,
      userName: user.userName,
      bio: user.bio || "",
    });
    setIsEditing(true);
  };

  const handleUpdateProfile = (e: React.FormEvent) => {
    e.preventDefault();
    updateProfileMutation.mutate(editProfile);
  };

  const handleCreateCollection = (e: React.FormEvent) => {
    e.preventDefault();
    createCollectionMutation.mutate({
      name: newCollection.name,
      description: newCollection.description,
      isPublic: newCollection.visibility === "public",
    });
  };

  if (userLoading || (!username && !currentCustomer)) {
    return (
      <PageLayout>
        <div className="flex h-64 items-center justify-center">
          <Loader2 className={cn(pageLoader, "h-12 w-12")} />
        </div>
      </PageLayout>
    );
  }

  if (userError || !user) {
    const status = (fetchError as any)?.response?.status;
    if (status === 401) {
      localStorage.removeItem("token");
      localStorage.removeItem("refresh_token");
      if (window.location.pathname !== "/auth") {
        window.location.href = "/auth";
      }
      return null;
    }

    if (isOwnProfile) {
      navigate("/profile/create", { replace: true });
      return null;
    }

    return (
      <PageLayout>
        <div className={cn(pageEmpty, "mx-auto max-w-lg p-8")}>
          <h2 className="text-2xl font-bold text-[#1A3C34] dark:text-white">User Not Found</h2>
          <p className="mt-2 text-gray-500">
            The profile you are looking for doesn't exist or has been removed.
          </p>
          <Button onClick={() => navigate("/feed/users")} className={cn(pageBtnPrimary, "mt-6")}>
            Back to Users Feed
          </Button>
        </div>
      </PageLayout>
    );
  }

  const userPosts = postsData?.Posts || [];

  return (
    <PageLayout>
      <PageHeader title="Profile" description={`@${user.userName}`} />

      <div className={cn(pagePanel, "relative mb-8 overflow-hidden p-6 md:p-8")}>
        <div className="flex flex-col gap-6 md:flex-row md:items-start">
          <img
            src={user.avatarURL || "https://via.placeholder.com/150"}
            alt={user.userName}
            className="h-24 w-24 rounded-full border border-gray-200 object-cover md:h-28 md:w-28 dark:border-[#2f5e50]"
          />
          <div className="flex-1 space-y-3">
            <div className="flex flex-wrap items-center gap-2">
              <h2 className="text-2xl font-semibold text-[#1A3C34] dark:text-white">
                {user.firstName} {user.lastName}
              </h2>
              <Badge variant="outline" className="border-gray-200 dark:border-[#2f5e50]">
                Guest
              </Badge>
            </div>
            {user.bio ? (
              <p className="text-sm leading-relaxed text-gray-500 dark:text-gray-400">
                {user.bio}
              </p>
            ) : null}
            <div className="flex items-center gap-6 text-sm">
              <span className="text-gray-500 dark:text-gray-400">
                <strong className="text-[#1A3C34] dark:text-white">{userPosts.length}</strong>{" "}
                posts
              </span>
              <span className="text-gray-500 dark:text-gray-400">
                <strong className="text-[#1A3C34] dark:text-white">{user.followers || 0}</strong>{" "}
                followers
              </span>
              <span className="text-gray-500 dark:text-gray-400">
                <strong className="text-[#1A3C34] dark:text-white">{user.following || 0}</strong>{" "}
                following
              </span>
            </div>
          </div>
          {isOwnProfile ? (
            <Button onClick={handleOpenEdit} className={pageBtnSecondary}>
              <Edit3 className="mr-2 h-4 w-4" />
              Edit Profile
            </Button>
          ) : null}
        </div>
      </div>

      <div className="mb-6 flex items-center gap-2 border-b border-gray-200 pb-2 dark:border-[#2f5e50]">
        <Button
          variant="ghost"
          onClick={() => setActiveTab("posts")}
          className={cn(
            "rounded-full",
            activeTab === "posts" &&
              "bg-gray-100 text-[#1A3C34] dark:bg-[#0d241d] dark:text-white"
          )}
        >
          <Grid className="mr-2 h-4 w-4" />
          Recent Stories
        </Button>
        {isOwnProfile ? (
          <Button
            variant="ghost"
            onClick={() => setActiveTab("collections")}
            className={cn(
              "rounded-full",
              activeTab === "collections" &&
                "bg-gray-100 text-[#1A3C34] dark:bg-[#0d241d] dark:text-white"
            )}
          >
            <Bookmark className="mr-2 h-4 w-4" />
            Saved Collections
          </Button>
        ) : null}
      </div>

        {activeTab === "posts" ? (
          <div className="mx-auto max-w-3xl space-y-6">
            {postsLoading ? (
              <div className={cn(pageEmpty, "py-12")}>
                <p className="text-gray-500">Loading posts...</p>
              </div>
            ) : userPosts.length > 0 ? (
              userPosts.map((post) => <GuestPostCard key={post.id} post={post} />)
            ) : (
              <div className={cn(pageEmpty, "py-12")}>
                <p className="text-gray-500">This user hasn't posted anything yet.</p>
              </div>
            )}
          </div>
        ) : (
          <div className="grid grid-cols-1 gap-8 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
            {collectionsLoading ? (
              <div className={cn(pageEmpty, "py-12")}>
                <p className="text-gray-500">Loading collections...</p>
              </div>
            ) : collections && collections.length > 0 ? (
              collections.map((collection) => (
                <div key={collection.id} className={cn(pagePanel, "p-6")}>
                  <div className="flex items-start justify-between gap-4">
                    <div>
                      <h3 className="text-lg font-semibold text-[#1A3C34] dark:text-white">
                        {collection.name}
                      </h3>
                      <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
                        {collection.description || "No description provided."}
                      </p>
                    </div>
                    <Badge variant="outline" className="border-gray-200 dark:border-[#2f5e50]">
                      {collection.isPublic ? "Public" : "Private"}
                    </Badge>
                  </div>
                </div>
              ))
            ) : (
              <div className={cn(pageEmpty, "py-12")}>
                <Bookmark className="mx-auto mb-3 h-10 w-10 text-gray-300" />
                <p className="text-lg font-semibold text-[#1A3C34] dark:text-white">
                  No Collections Yet
                </p>
                <p className="mt-1 text-sm text-gray-500">
                  Create your first collection to start saving favorite venues.
                </p>
                <Button
                  onClick={() => setIsCreatingCollection(true)}
                  className={cn(pageBtnPrimary, "mt-5")}
                >
                  <Plus className="mr-2 h-4 w-4" />
                  New Collection
                </Button>
              </div>
            )}
          </div>
        )}

      <Dialog open={isEditing} onOpenChange={setIsEditing}>
        <DialogContent className={cn(pagePanel, "border-0 p-6 sm:max-w-md")}>
          <DialogHeader>
            <DialogTitle className="text-[#1A3C34] dark:text-white">Edit Profile</DialogTitle>
          </DialogHeader>
          <form onSubmit={handleUpdateProfile} className="space-y-4">
            <div className="space-y-2">
              <label htmlFor="editUserName" className={pageLabel}>
                Username
              </label>
              <input
                id="editUserName"
                required
                value={editProfile.userName}
                onChange={(e) =>
                  setEditProfile({ ...editProfile, userName: e.target.value })
                }
                className={pageInput}
              />
            </div>
            <div className="grid grid-cols-2 gap-3">
              <div className="space-y-2">
                <label htmlFor="editFirstName" className={pageLabel}>
                  First Name
                </label>
                <input
                  id="editFirstName"
                  required
                  value={editProfile.firstName}
                  onChange={(e) =>
                    setEditProfile({ ...editProfile, firstName: e.target.value })
                  }
                  className={pageInput}
                />
              </div>
              <div className="space-y-2">
                <label htmlFor="editLastName" className={pageLabel}>
                  Last Name
                </label>
                <input
                  id="editLastName"
                  required
                  value={editProfile.lastName}
                  onChange={(e) =>
                    setEditProfile({ ...editProfile, lastName: e.target.value })
                  }
                  className={pageInput}
                />
              </div>
            </div>
            <div className="space-y-2">
              <label htmlFor="editBio" className={pageLabel}>
                Bio
              </label>
              <textarea
                id="editBio"
                value={editProfile.bio}
                onChange={(e) =>
                  setEditProfile({ ...editProfile, bio: e.target.value })
                }
                className={cn(pageInput, "min-h-[100px] resize-y py-3")}
              />
            </div>
            <DialogFooter className="gap-2">
              <Button type="button" className={pageBtnSecondary} onClick={() => setIsEditing(false)}>
                Cancel
              </Button>
              <Button type="submit" className={pageBtnPrimary} disabled={updateProfileMutation.isPending}>
                {updateProfileMutation.isPending ? "Saving..." : "Save Changes"}
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>

      <Dialog open={isCreatingCollection} onOpenChange={setIsCreatingCollection}>
        <DialogContent className={cn(pagePanel, "border-0 p-6 sm:max-w-md")}>
          <DialogHeader>
            <DialogTitle className="text-[#1A3C34] dark:text-white">New Collection</DialogTitle>
          </DialogHeader>
          <form onSubmit={handleCreateCollection} className="space-y-4">
            <div className="space-y-2">
              <label htmlFor="collectionName" className={pageLabel}>
                Collection Name
              </label>
              <input
                id="collectionName"
                required
                placeholder="e.g. Best Coffee Shops"
                value={newCollection.name}
                onChange={(e) =>
                  setNewCollection({ ...newCollection, name: e.target.value })
                }
                className={pageInput}
              />
            </div>
            <div className="space-y-2">
              <label htmlFor="collectionDesc" className={pageLabel}>
                Description (Optional)
              </label>
              <textarea
                id="collectionDesc"
                placeholder="Places to visit this summer"
                value={newCollection.description}
                onChange={(e) =>
                  setNewCollection({
                    ...newCollection,
                    description: e.target.value,
                  })
                }
                className={cn(pageInput, "min-h-[90px] resize-y py-3")}
              />
            </div>
            <div className="space-y-2">
              <span className={pageLabel}>Visibility</span>
              <Select
                value={newCollection.visibility}
                onValueChange={(value) =>
                  setNewCollection({ ...newCollection, visibility: value })
                }
              >
                <SelectTrigger className={cn(pageInput, "h-auto cursor-pointer")}>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="private">Private</SelectItem>
                  <SelectItem value="public">Public</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <DialogFooter className="gap-2">
              <Button
                type="button"
                className={pageBtnSecondary}
                onClick={() => setIsCreatingCollection(false)}
              >
                Cancel
              </Button>
              <Button type="submit" className={pageBtnPrimary} disabled={createCollectionMutation.isPending}>
                {createCollectionMutation.isPending ? "Creating..." : "Create Collection"}
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>
    </PageLayout>
  );
}
