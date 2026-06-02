import React, { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { apiClient } from '../../api/client';
import { FolderHeart, Plus, Users, UserPlus } from 'lucide-react';
import { Button } from '../../components/ui/button';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger, DialogFooter } from '../../components/ui/dialog';
import { Input } from '../../components/ui/input';
import { Label } from '../../components/ui/label';
import { toast } from 'sonner';

export const CollectionsPage: React.FC = () => {
  const queryClient = useQueryClient();
  const [isCreateOpen, setIsCreateOpen] = useState(false);
  const [isInviteOpen, setIsInviteOpen] = useState(false);
  const [selectedCollectionId, setSelectedCollectionId] = useState<string | null>(null);
  const [newCollectionName, setNewCollectionName] = useState('');
  const [inviteEmail, setInviteEmail] = useState('');

  const { data: collections, isLoading } = useQuery({
    queryKey: ['collections'],
    queryFn: () => apiClient.getCollections(),
  });

  const createMutation = useMutation({
    mutationFn: (name: string) => apiClient.createCollection({ name, isPublic: false }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['collections'] });
      setIsCreateOpen(false);
      setNewCollectionName('');
      toast.success("Collection created successfully");
    },
    onError: () => toast.error("Failed to create collection")
  });

  const inviteMutation = useMutation({
    mutationFn: ({ id, email }: { id: string, email: string }) => apiClient.inviteCollaborator(id, email),
    onSuccess: () => {
      setIsInviteOpen(false);
      setInviteEmail('');
      toast.success("Invitation sent successfully");
    },
    onError: () => toast.error("Failed to send invitation")
  });

  const handleCreate = (e: React.FormEvent) => {
    e.preventDefault();
    if (newCollectionName.trim()) {
      createMutation.mutate(newCollectionName);
    }
  };

  const handleInvite = (e: React.FormEvent) => {
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
    <div className="flex flex-col items-center w-full min-h-screen bg-background pt-8 pb-16 px-4">
      <div className="max-w-4xl w-full">
        <header className="mb-10 flex flex-col sm:flex-row items-center justify-between gap-6">
          <div className="flex items-center gap-4">
            <div className="p-4 bg-primary/10 text-primary rounded-2xl shadow-sm">
              <FolderHeart size={28} />
            </div>
            <div>
              <h1 className="text-4xl font-serif font-bold tracking-tight text-foreground">My Collections</h1>
              <p className="text-muted-foreground text-base mt-1 font-medium">Organize and share your favorite places</p>
            </div>
          </div>
          
          <Dialog open={isCreateOpen} onOpenChange={setIsCreateOpen}>
            <DialogTrigger asChild>
              <Button className="gap-2 rounded-full h-11 px-6 bg-accent text-accent-foreground hover:bg-[#e6c200] font-bold shadow-md hover:shadow-lg transition-all hover:-translate-y-0.5">
                <Plus size={18} /> New Collection
              </Button>
            </DialogTrigger>
            <DialogContent className="rounded-3xl sm:rounded-3xl p-6">
              <DialogHeader>
                <DialogTitle className="text-2xl font-serif">Create New Collection</DialogTitle>
              </DialogHeader>
              <form onSubmit={handleCreate} className="space-y-5 py-4">
                <div className="space-y-2">
                  <Label htmlFor="name">Name</Label>
                  <Input 
                    id="name" 
                    placeholder="e.g. Best Coffee Shops" 
                    value={newCollectionName}
                    onChange={(e) => setNewCollectionName(e.target.value)}
                    autoFocus
                    className="h-12 rounded-xl bg-muted/50 focus-visible:ring-primary"
                  />
                </div>
                <DialogFooter>
                  <Button type="button" variant="outline" className="rounded-xl h-11" onClick={() => setIsCreateOpen(false)}>Cancel</Button>
                  <Button type="submit" className="rounded-xl h-11" disabled={!newCollectionName.trim() || createMutation.isPending}>
                    Create
                  </Button>
                </DialogFooter>
              </form>
            </DialogContent>
          </Dialog>
        </header>

        {isLoading ? (
          <div className="flex items-center justify-center py-16 text-muted-foreground font-serif text-lg">Loading collections...</div>
        ) : !collections?.length ? (
          <div className="flex flex-col items-center justify-center py-20 text-muted-foreground bg-card rounded-3xl border border-border shadow-sm">
            <FolderHeart size={56} className="mb-6 opacity-20" />
            <p className="text-2xl font-serif font-medium mb-2">No collections yet</p>
            <p className="text-base mb-8">Create one to start saving your favorite posts.</p>
            <Button variant="outline" className="rounded-full h-11 px-6 font-bold" onClick={() => setIsCreateOpen(true)}>Create Collection</Button>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
            {collections.map((collection) => (
              <div key={collection.id} className="group bg-card rounded-3xl border border-border/50 p-6 shadow-sm flex flex-col hover:shadow-xl dark:hover:shadow-primary/5 hover:-translate-y-1 transition-all duration-300 relative overflow-hidden">
                <div className="absolute inset-0 bg-gradient-to-br from-primary/5 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition-opacity"></div>
                <div className="relative z-10 flex flex-col h-full">
                  <div className="flex justify-between items-start mb-4">
                    <h3 className="font-serif font-bold text-2xl text-foreground truncate pr-4 group-hover:text-primary transition-colors">{collection.name}</h3>
                    <div className="flex items-center gap-1.5 bg-muted/80 backdrop-blur-sm px-3 py-1.5 rounded-full text-[10px] font-black uppercase tracking-wider text-muted-foreground">
                      <Users size={12} /> Private
                    </div>
                  </div>
                  <p className="text-muted-foreground text-sm line-clamp-2 mb-8 flex-1 leading-relaxed">
                    {collection.description || "No description provided."}
                  </p>
                  
                  <div className="flex justify-between items-center pt-5 border-t border-border/50 mt-auto">
                    <span className="text-xs font-semibold text-muted-foreground uppercase tracking-wider">
                      Created {new Date(collection.createdAt).toLocaleDateString()}
                    </span>
                    <Button variant="ghost" size="sm" className="gap-2 text-primary font-bold hover:bg-primary/10 rounded-xl" onClick={() => openInviteModal(collection.id)}>
                      <UserPlus size={16} /> Invite
                    </Button>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}

        {/* Invite Dialog */}
        <Dialog open={isInviteOpen} onOpenChange={setIsInviteOpen}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Invite Collaborator</DialogTitle>
            </DialogHeader>
            <form onSubmit={handleInvite} className="space-y-4 py-4">
              <div className="space-y-2">
                <Label htmlFor="email">Email address</Label>
                <Input 
                  id="email" 
                  type="email"
                  placeholder="friend@example.com" 
                  value={inviteEmail}
                  onChange={(e) => setInviteEmail(e.target.value)}
                  autoFocus
                />
              </div>
              <DialogFooter>
                <Button type="button" variant="outline" onClick={() => setIsInviteOpen(false)}>Cancel</Button>
                <Button type="submit" disabled={!inviteEmail.trim() || inviteMutation.isPending}>
                  Send Invite
                </Button>
              </DialogFooter>
            </form>
          </DialogContent>
        </Dialog>
      </div>
    </div>
  );
};
