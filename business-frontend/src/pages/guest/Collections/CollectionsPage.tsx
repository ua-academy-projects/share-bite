import { useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { FolderHeart, Loader2, Plus, UserPlus } from "lucide-react";
import { toast } from "sonner";
import { apiClient } from "@/api/client";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent } from "@/components/ui/card";
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { PageHeader } from "@/components/layout/PageHeader";

export function CollectionsPage() {
  const queryClient = useQueryClient();
  const [isCreateOpen, setIsCreateOpen] = useState(false);
  const [isInviteOpen, setIsInviteOpen] = useState(false);
  const [selectedCollectionId, setSelectedCollectionId] = useState<string | null>(
    null
  );
  const [newCollectionName, setNewCollectionName] = useState("");
  const [inviteEmail, setInviteEmail] = useState("");

  const { data: collections, isLoading } = useQuery({
    queryKey: ["collections"],
    queryFn: () => apiClient.getCollections(),
  });

  const createMutation = useMutation({
    mutationFn: (name: string) =>
      apiClient.createCollection({ name, isPublic: false }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["collections"] });
      setIsCreateOpen(false);
      setNewCollectionName("");
      toast.success("Collection created");
    },
    onError: () => toast.error("Failed to create collection"),
  });

  const inviteMutation = useMutation({
    mutationFn: ({ id, email }: { id: string; email: string }) =>
      apiClient.inviteCollaborator(id, email),
    onSuccess: () => {
      setIsInviteOpen(false);
      setInviteEmail("");
      toast.success("Invitation sent");
    },
    onError: () => toast.error("Failed to send invitation"),
  });

  return (
    <div className="px-6 py-8 lg:px-10">
      <PageHeader
        title="My Collections"
        description="Organize and share your favorite places"
        icon={FolderHeart}
      >
        <Dialog open={isCreateOpen} onOpenChange={setIsCreateOpen}>
          <DialogTrigger asChild>
            <Button className="rounded-full bg-accent font-bold text-accent-foreground">
              <Plus className="mr-2 h-4 w-4" /> New Collection
            </Button>
          </DialogTrigger>
          <DialogContent className="rounded-3xl">
            <DialogHeader>
              <DialogTitle>Create collection</DialogTitle>
            </DialogHeader>
            <form
              onSubmit={(e) => {
                e.preventDefault();
                if (newCollectionName.trim())
                  createMutation.mutate(newCollectionName);
              }}
              className="space-y-4"
            >
              <div className="space-y-2">
                <Label htmlFor="name">Name</Label>
                <Input
                  id="name"
                  value={newCollectionName}
                  onChange={(e) => setNewCollectionName(e.target.value)}
                  placeholder="Best coffee shops"
                  className="rounded-xl"
                />
              </div>
              <DialogFooter>
                <Button type="submit" disabled={createMutation.isPending}>
                  Create
                </Button>
              </DialogFooter>
            </form>
          </DialogContent>
        </Dialog>
      </PageHeader>

      {isLoading ? (
        <div className="flex justify-center py-16">
          <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
        </div>
      ) : !collections?.length ? (
        <Card className="mx-auto max-w-lg rounded-3xl bg-card-solid">
          <CardContent className="py-16 text-center text-muted-foreground">
            No collections yet. Create one to save venues you love.
          </CardContent>
        </Card>
      ) : (
        <div className="mx-auto grid max-w-4xl gap-4 sm:grid-cols-2">
          {collections.map((collection) => (
            <Card
              key={collection.id}
              className="rounded-2xl border-border bg-card-solid"
            >
              <CardContent className="flex items-start justify-between p-5">
                <div>
                  <h3 className="font-bold text-foreground">{collection.name}</h3>
                  {collection.description && (
                    <p className="mt-1 text-sm text-muted-foreground">
                      {collection.description}
                    </p>
                  )}
                </div>
                <Button
                  variant="outline"
                  size="sm"
                  className="rounded-full"
                  onClick={() => {
                    setSelectedCollectionId(collection.id);
                    setIsInviteOpen(true);
                  }}
                >
                  <UserPlus className="h-4 w-4" />
                </Button>
              </CardContent>
            </Card>
          ))}
        </div>
      )}

      <Dialog open={isInviteOpen} onOpenChange={setIsInviteOpen}>
        <DialogContent className="rounded-3xl">
          <DialogHeader>
            <DialogTitle>Invite collaborator</DialogTitle>
          </DialogHeader>
          <form
            onSubmit={(e) => {
              e.preventDefault();
              if (inviteEmail.trim() && selectedCollectionId) {
                inviteMutation.mutate({
                  id: selectedCollectionId,
                  email: inviteEmail,
                });
              }
            }}
            className="space-y-4"
          >
            <div className="space-y-2">
              <Label htmlFor="email">Email</Label>
              <Input
                id="email"
                type="email"
                value={inviteEmail}
                onChange={(e) => setInviteEmail(e.target.value)}
                className="rounded-xl"
              />
            </div>
            <DialogFooter>
              <Button type="submit" disabled={inviteMutation.isPending}>
                Send invite
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>
    </div>
  );
}
