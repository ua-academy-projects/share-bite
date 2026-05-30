import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Bookmark, Edit3, Grid, Loader2 } from "lucide-react";
import { toast } from "sonner";
import { clsx } from "clsx";
import { apiClient } from "@/api/client";
import { useCurrentCustomer } from "@/hooks/useCurrentCustomer";
import { GuestPostCard } from "@/components/PostCard/PostCard";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent } from "@/components/ui/card";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Badge } from "@/components/ui/badge";

type UserProfileProps = {
  mode?: "create" | "edit";
};

export function UserProfile({ mode }: UserProfileProps) {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { data: currentCustomer } = useCurrentCustomer();

  const username = id || currentCustomer?.userName;
  const isOwnProfile =
    !id || (currentCustomer && currentCustomer.userName === id);

  const [activeTab, setActiveTab] = useState<"posts" | "collections">("posts");
  const [showEditForm, setShowEditForm] = useState(mode === "edit");
  const [newProfile, setNewProfile] = useState({
    firstName: "",
    lastName: "",
    userName: "",
    bio: "",
  });
  const [editProfile, setEditProfile] = useState({
    firstName: "",
    lastName: "",
    userName: "",
    bio: "",
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

  const { data: postsData, isLoading: postsLoading } = useQuery({
    queryKey: ["userPosts", user?.id],
    queryFn: () => apiClient.getPosts(100, 0, undefined, user!.id),
    enabled: !!user?.id,
  });

  const { data: collections, isLoading: collectionsLoading } = useQuery({
    queryKey: ["collections", isOwnProfile ? "mine" : user?.id],
    queryFn: () =>
      apiClient.getCollections(isOwnProfile ? undefined : user?.id),
    enabled:
      activeTab === "collections" && (isOwnProfile || !!user?.id),
  });

  useEffect(() => {
    if (mode === "edit" && user) {
      setEditProfile({
        firstName: user.firstName,
        lastName: user.lastName,
        userName: user.userName,
        bio: user.bio || "",
      });
      setShowEditForm(true);
    }
  }, [mode, user]);

  const createProfileMutation = useMutation({
    mutationFn: () => apiClient.createCustomer(newProfile),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["currentCustomer"] });
      toast.success("Profile created!");
      navigate("/profile", { replace: true });
    },
    onError: (error: unknown) => {
      const err = error as { response?: { data?: { error?: string } } };
      toast.error(err?.response?.data?.error || "Failed to create profile");
    },
  });

  const updateProfileMutation = useMutation({
    mutationFn: () => apiClient.updateCustomer(editProfile),
    onSuccess: (updated) => {
      queryClient.invalidateQueries({ queryKey: ["currentCustomer"] });
      queryClient.invalidateQueries({ queryKey: ["user", updated.userName] });
      setShowEditForm(false);
      toast.success("Profile updated!");
      if (mode === "edit") navigate("/profile", { replace: true });
    },
    onError: (error: unknown) => {
      const err = error as { response?: { data?: { error?: string } } };
      toast.error(err?.response?.data?.error || "Failed to update profile");
    },
  });

  if (userLoading || (!username && !currentCustomer)) {
    return (
      <div className="flex justify-center py-24">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (userError || !user) {
    const status = (fetchError as { response?: { status?: number } })?.response
      ?.status;
    if (status === 401) {
      localStorage.removeItem("token");
      navigate("/auth", { replace: true });
      return null;
    }

    if (isOwnProfile || mode === "create") {
      return (
        <div className="px-6 py-8 lg:px-10">
          <Card className="mx-auto max-w-md rounded-3xl bg-card-solid">
            <CardContent className="p-8">
              <h2 className="mb-2 text-2xl font-bold">Create profile</h2>
              <p className="mb-6 text-muted-foreground">
                Set up your guest profile to post and interact.
              </p>
              <form
                onSubmit={(e) => {
                  e.preventDefault();
                  createProfileMutation.mutate();
                }}
                className="space-y-4"
              >
                <div className="space-y-2">
                  <Label htmlFor="userName">Username</Label>
                  <Input
                    id="userName"
                    required
                    value={newProfile.userName}
                    onChange={(e) =>
                      setNewProfile({ ...newProfile, userName: e.target.value })
                    }
                    className="rounded-xl"
                  />
                </div>
                <div className="grid grid-cols-2 gap-3">
                  <div className="space-y-2">
                    <Label htmlFor="firstName">First name</Label>
                    <Input
                      id="firstName"
                      required
                      value={newProfile.firstName}
                      onChange={(e) =>
                        setNewProfile({
                          ...newProfile,
                          firstName: e.target.value,
                        })
                      }
                      className="rounded-xl"
                    />
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="lastName">Last name</Label>
                    <Input
                      id="lastName"
                      required
                      value={newProfile.lastName}
                      onChange={(e) =>
                        setNewProfile({
                          ...newProfile,
                          lastName: e.target.value,
                        })
                      }
                      className="rounded-xl"
                    />
                  </div>
                </div>
                <div className="space-y-2">
                  <Label htmlFor="bio">Bio</Label>
                  <Input
                    id="bio"
                    value={newProfile.bio}
                    onChange={(e) =>
                      setNewProfile({ ...newProfile, bio: e.target.value })
                    }
                    className="rounded-xl"
                  />
                </div>
                <Button
                  type="submit"
                  disabled={createProfileMutation.isPending}
                  className="w-full rounded-xl bg-accent font-bold text-accent-foreground"
                >
                  Create profile
                </Button>
              </form>
            </CardContent>
          </Card>
        </div>
      );
    }

    return (
      <div className="flex min-h-[50vh] items-center justify-center px-4">
        <Card className="max-w-md rounded-3xl bg-card-solid">
          <CardContent className="p-8 text-center">
            <h2 className="text-xl font-bold">User not found</h2>
            <Button className="mt-6 rounded-xl" onClick={() => navigate("/")}>
              Back home
            </Button>
          </CardContent>
        </Card>
      </div>
    );
  }

  const userPosts = postsData?.Posts || [];

  return (
    <div className="px-6 py-8 lg:px-10">
      <Card className="relative mx-auto mb-8 max-w-2xl overflow-hidden rounded-3xl bg-card-solid">
        <CardContent className="flex flex-col items-center gap-6 p-8 md:flex-row md:items-start">
          <img
            src={user.avatarURL || "https://via.placeholder.com/150"}
            alt=""
            className="h-32 w-32 rounded-full border-4 border-background object-cover shadow-xl"
          />
          <div className="flex-1 text-center md:text-left">
            <h1 className="text-3xl font-bold">
              {user.firstName} {user.lastName}
            </h1>
            <div className="mt-2 flex flex-wrap items-center justify-center gap-2 md:justify-start">
              <Badge variant="outline">@{user.userName}</Badge>
              <Badge variant="secondary">Guest</Badge>
            </div>
            {user.bio && (
              <p className="mt-4 text-sm text-muted-foreground">{user.bio}</p>
            )}
            <p className="mt-4 text-sm text-muted-foreground">
              {userPosts.length} posts
            </p>
          </div>
          {isOwnProfile && (
            <Button
              variant="outline"
              className="rounded-xl"
              onClick={() => {
                setEditProfile({
                  firstName: user.firstName,
                  lastName: user.lastName,
                  userName: user.userName,
                  bio: user.bio || "",
                });
                setShowEditForm(true);
              }}
            >
              <Edit3 className="mr-2 h-4 w-4" /> Edit profile
            </Button>
          )}
        </CardContent>
      </Card>

      <div className="mx-auto mb-6 flex max-w-2xl gap-6 border-b border-border">
        <button
          type="button"
          className={clsx(
            "flex items-center gap-2 border-b-2 pb-3 text-sm font-semibold",
            activeTab === "posts"
              ? "border-primary text-foreground"
              : "border-transparent text-muted-foreground"
          )}
          onClick={() => setActiveTab("posts")}
        >
          <Grid className="h-4 w-4" /> Posts
        </button>
        {isOwnProfile && (
          <button
            type="button"
            className={clsx(
              "flex items-center gap-2 border-b-2 pb-3 text-sm font-semibold",
              activeTab === "collections"
                ? "border-primary text-foreground"
                : "border-transparent text-muted-foreground"
            )}
            onClick={() => setActiveTab("collections")}
          >
            <Bookmark className="h-4 w-4" /> Collections
          </button>
        )}
      </div>

      <div className="mx-auto max-w-xl">
        {activeTab === "posts" ? (
          postsLoading ? (
            <div className="flex justify-center py-12">
              <Loader2 className="h-6 w-6 animate-spin" />
            </div>
          ) : userPosts.length ? (
            <div className="flex flex-col gap-8">
              {userPosts.map((post) => (
                <GuestPostCard key={post.id} post={post} />
              ))}
            </div>
          ) : (
            <p className="py-12 text-center text-muted-foreground">
              No posts yet.
            </p>
          )
        ) : collectionsLoading ? (
          <div className="flex justify-center py-12">
            <Loader2 className="h-6 w-6 animate-spin" />
          </div>
        ) : (
          <div className="grid gap-4 sm:grid-cols-2">
            {collections?.map((c) => (
              <Card key={c.id} className="rounded-2xl bg-card-solid">
                <CardContent className="p-4">
                  <h3 className="font-bold">{c.name}</h3>
                  <p className="mt-1 text-sm text-muted-foreground">
                    {c.description || "No description"}
                  </p>
                </CardContent>
              </Card>
            ))}
          </div>
        )}
      </div>

      <Dialog open={showEditForm} onOpenChange={setShowEditForm}>
        <DialogContent className="rounded-3xl">
          <DialogHeader>
            <DialogTitle>Edit profile</DialogTitle>
          </DialogHeader>
          <form
            onSubmit={(e) => {
              e.preventDefault();
              updateProfileMutation.mutate();
            }}
            className="space-y-4"
          >
            <div className="space-y-2">
              <Label>Username</Label>
              <Input
                value={editProfile.userName}
                onChange={(e) =>
                  setEditProfile({ ...editProfile, userName: e.target.value })
                }
                className="rounded-xl"
              />
            </div>
            <div className="grid grid-cols-2 gap-3">
              <div className="space-y-2">
                <Label>First name</Label>
                <Input
                  value={editProfile.firstName}
                  onChange={(e) =>
                    setEditProfile({
                      ...editProfile,
                      firstName: e.target.value,
                    })
                  }
                  className="rounded-xl"
                />
              </div>
              <div className="space-y-2">
                <Label>Last name</Label>
                <Input
                  value={editProfile.lastName}
                  onChange={(e) =>
                    setEditProfile({
                      ...editProfile,
                      lastName: e.target.value,
                    })
                  }
                  className="rounded-xl"
                />
              </div>
            </div>
            <div className="space-y-2">
              <Label>Bio</Label>
              <Input
                value={editProfile.bio}
                onChange={(e) =>
                  setEditProfile({ ...editProfile, bio: e.target.value })
                }
                className="rounded-xl"
              />
            </div>
            <Button
              type="submit"
              disabled={updateProfileMutation.isPending}
              className="w-full rounded-xl bg-accent font-bold text-accent-foreground"
            >
              Save changes
            </Button>
          </form>
        </DialogContent>
      </Dialog>
    </div>
  );
}
