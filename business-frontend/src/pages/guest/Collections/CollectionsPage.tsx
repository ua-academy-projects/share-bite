import { useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { FolderHeart, Plus, UserPlus, Users } from "lucide-react";
import { toast } from "sonner";
import { apiClient } from "@/api/client";
import { PageHeader } from "@/components/layout/PageHeader";
import { PageLayout } from "@/components/layout/PageLayout";
import {
  pageBtnPrimary,
  pageBtnSecondary,
  pageEmpty,
  pageInput,
  pageLabel,
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
  DialogTrigger,
} from "@/components/ui/dialog";
import { Label } from "@/components/ui/label";
import { cn } from "@/lib/utils";

export function CollectionsPage() {
  const queryClient = useQueryClient();
  const [isCreateOpen, setIsCreateOpen] = useState<boolean>(false);
  const [isInviteOpen, setIsInviteOpen] = useState(false);
  const [selectedCollectionId, setSelectedCollectionId] = useState<string | null>(null);
  const [newCollectionName, setNewCollectionName] = useState("");
  const [inviteEmail, setInviteEmail] = useState("");

  const { data: collections, isLoading } = useQuery({
    queryKey: ["collections"],
    queryFn: () => apiClient.getCollections(),
  });

  const createMutation = useMutation({
    mutationFn: (name: string) => apiClient.createCollection({ name, isPublic: false }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["collections"] });
      setIsCreateOpen(false);
      setNewCollectionName("");
      toast.success("Collection created successfully");
    },
    onError: () => toast.error("Failed to create collection"),
  });

  const inviteMutation = useMutation({
    mutationFn: ({ id, email }: { id: string; email: string }) =>
      apiClient.inviteCollaborator(id, email),
    onSuccess: () => {
      setIsInviteOpen(false);
      setInviteEmail("");
      toast.success("Invitation sent successfully");
    },
    onError: () => toast.error("Failed to send invitation"),
  });

  const handleCreate = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (newCollectionName.trim()) {
      createMutation.mutate(newCollectionName);
    }
  };

  const handleInvite = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (inviteEmail.trim() && selectedCollectionId) {
      inviteMutation.mutate({ id: selectedCollectionId, email: inviteEmail });
    }
  };

  const openInviteModal = (id: string) => {
    setSelectedCollectionId(id);
    setIsInviteOpen(true);
  };

  return (
    <PageLayout maxWidth="5xl">
      <div className="mb-8 flex flex-col items-start justify-between gap-4 sm:flex-row sm:items-center">
        <PageHeader
          title="My Collections"
          description="Organize and share your favorite places"
          className="mb-0"
        />
        <Dialog open={isCreateOpen} onOpenChange={setIsCreateOpen}>
          <DialogTrigger asChild>
            <Button className={cn(pageBtnPrimary, "gap-2")}>
              <Plus size={18} /> New Collection
            </Button>
          </DialogTrigger>
          <DialogContent className={cn(pagePanel, "border-0 p-6 sm:max-w-md")}>
            <DialogHeader>
              <DialogTitle className="text-[#1A3C34] dark:text-white">
                Create New Collection
              </DialogTitle>
            </DialogHeader>
            <form onSubmit={handleCreate} className="space-y-5 py-4">
              <div className="space-y-2">
                <Label htmlFor="name" className={pageLabel}>
                  Name
                </Label>
                <input
                  id="name"
                  placeholder="e.g. Best Coffee Shops"
                  value={newCollectionName}
                  onChange={(e) => setNewCollectionName(e.target.value)}
                  autoFocus
                  className={pageInput}
                />
              </div>
              <DialogFooter className="gap-2">
                <Button
                  type="button"
                  className={pageBtnSecondary}
                  onClick={() => setIsCreateOpen(false)}
                >
                  Cancel
                </Button>
                <Button
                  type="submit"
                  className={pageBtnPrimary}
                  disabled={!newCollectionName.trim() || createMutation.isPending}
                >
                  Create
                </Button>
              </DialogFooter>
            </form>
          </DialogContent>
        </Dialog>
      </div>

      {isLoading ? (
        <div className={cn(pageEmpty, "py-16")}>
          <p className="text-gray-500">Loading collections...</p>
        </div>
      ) : !collections?.length ? (
        <div className={pageEmpty}>
          <FolderHeart size={56} className="mx-auto mb-6 opacity-20" />
          <p className="mb-2 text-xl font-bold text-[#1A3C34] dark:text-gray-300">
            No collections yet
          </p>
          <p className="mb-8 text-gray-500">
            Create one to start saving your favorite posts.
          </p>
          <Button className={pageBtnPrimary} onClick={() => setIsCreateOpen(true)}>
            Create Collection
          </Button>
        </div>
      ) : (
        <div className="grid grid-cols-1 gap-8 md:grid-cols-2">
          {collections.map((collection) => (
            <div
              key={collection.id}
              className={cn(
                pagePanel,
                "group relative overflow-hidden transition-transform hover:-translate-y-1"
              )}
            >
              <div className="flex h-full flex-col p-6">
                <div className="mb-4 flex items-start justify-between">
                  <h3 className="truncate pr-4 text-xl font-semibold text-[#1A3C34] transition-colors group-hover:text-emerald-600 dark:text-white dark:group-hover:text-[#98FF98]">
                    {collection.name}
                  </h3>
                  <Badge
                    variant="outline"
                    className="gap-1 border-gray-200 dark:border-[#2f5e50]"
                  >
                    <Users size={12} /> Private
                  </Badge>
                </div>
                <p className="mb-8 line-clamp-2 flex-1 text-sm leading-relaxed text-gray-500 dark:text-gray-400">
                  {collection.description || "No description provided."}
                </p>

                <div className="mt-auto flex items-center justify-between border-t border-gray-200 pt-5 dark:border-[#2f5e50]">
                  <span className="text-xs text-gray-500 dark:text-gray-400">
                    Created {new Date(collection.createdAt).toLocaleDateString()}
                  </span>
                  <Button
                    variant="ghost"
                    size="sm"
                    className="gap-2 rounded-xl text-emerald-600 hover:bg-emerald-500/10 dark:text-[#98FF98]"
                    onClick={() => openInviteModal(collection.id)}
                  >
                    <UserPlus size={16} /> Invite
                  </Button>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}

      <Dialog open={isInviteOpen} onOpenChange={setIsInviteOpen}>
        <DialogContent className={cn(pagePanel, "border-0 p-6 sm:max-w-md")}>
          <DialogHeader>
            <DialogTitle className="text-[#1A3C34] dark:text-white">
              Invite Collaborator
            </DialogTitle>
          </DialogHeader>
          <form onSubmit={handleInvite} className="space-y-4 py-4">
            <div className="space-y-2">
              <Label htmlFor="email" className={pageLabel}>
                Email address
              </Label>
              <input
                id="email"
                type="email"
                placeholder="friend@example.com"
                value={inviteEmail}
                onChange={(e) => setInviteEmail(e.target.value)}
                autoFocus
                className={pageInput}
              />
            </div>
            <DialogFooter className="gap-2">
              <Button
                type="button"
                className={pageBtnSecondary}
                onClick={() => setIsInviteOpen(false)}
              >
                Cancel
              </Button>
              <Button
                type="submit"
                className={pageBtnPrimary}
                disabled={!inviteEmail.trim() || inviteMutation.isPending}
              >
                Send Invite
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>
    </PageLayout>
  );
}
