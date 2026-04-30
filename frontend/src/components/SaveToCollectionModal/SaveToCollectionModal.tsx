import React, { useState } from 'react';
import { Bookmark, Plus, X, Check } from 'lucide-react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { apiClient } from '../../api/client';
import styles from './SaveToCollectionModal.module.css';
import { clsx } from 'clsx';

interface SaveToCollectionModalProps {
  venueId: number;
  isOpen: boolean;
  onClose: () => void;
}

export const SaveToCollectionModal: React.FC<SaveToCollectionModalProps> = ({ 
  venueId, 
  isOpen, 
  onClose 
}) => {
  const queryClient = useQueryClient();
  const [showCreate, setShowCreate] = useState(false);
  const [newCollectionName, setNewCollectionName] = useState('');

  const { data: collectionsData, isLoading } = useQuery({
    queryKey: ['myCollections'],
    queryFn: () => apiClient.listMyCollections(100),
    enabled: isOpen
  });

  const { data: venueCollections } = useQuery({
    queryKey: ['venueCollections', venueId],
    queryFn: async () => {
      // Since there's no direct endpoint to check which collections a venue is in,
      // we'd ideally have one. For now, we'll have to fetch venues for each collection
      // OR the backend should provide this. 
      // ASSUMPTION: We'll just manage it via state or the user will see it update.
      // For this PR, I'll assume we fetch all my collections and then we'll need to know
      // which ones contain this venue.
      // BUT the backend doesn't easily give "collections containing venue X".
      // I'll skip the "already in" check for now or assume it's handled by add/remove feedback.
      return [] as string[]; 
    },
    enabled: false // Disabled for now due to lack of endpoint
  });

  const addMutation = useMutation({
    mutationFn: (collectionId: string) => apiClient.addVenueToCollection(collectionId, venueId),
    onSuccess: () => {
      // Invalidate venues for this collection if we were viewing it
      queryClient.invalidateQueries({ queryKey: ['collectionVenues'] });
    }
  });

  const removeMutation = useMutation({
    mutationFn: (collectionId: string) => apiClient.removeVenueFromCollection(collectionId, venueId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['collectionVenues'] });
    }
  });

  const createMutation = useMutation({
    mutationFn: (name: string) => apiClient.createCollection({ name, isPublic: true }),
    onSuccess: (newCollection) => {
      setNewCollectionName('');
      setShowCreate(false);
      queryClient.invalidateQueries({ queryKey: ['myCollections'] });
      // Automatically add venue to the new collection
      addMutation.mutate(newCollection.id);
    }
  });

  if (!isOpen) return null;

  return (
    <div className={styles.overlay} onClick={onClose}>
      <div className={clsx(styles.modal, 'glass-panel')} onClick={e => e.stopPropagation()}>
        <div className={styles.header}>
          <div className={styles.titleGroup}>
            <Bookmark className={styles.titleIcon} size={20} />
            <h2 className={styles.title}>Save to Collection</h2>
          </div>
          <button className={styles.closeBtn} onClick={onClose}>
            <X size={20} />
          </button>
        </div>

        <div className={styles.content}>
          {isLoading ? (
            <div className={styles.loading}>Loading collections...</div>
          ) : (
            <div className={styles.list}>
              {collectionsData?.collections.map(collection => (
                <button 
                  key={collection.id} 
                  className={styles.item}
                  onClick={() => addMutation.mutate(collection.id)}
                  disabled={addMutation.isPending}
                >
                  <span className={styles.itemName}>{collection.name}</span>
                  {addMutation.variables === collection.id && addMutation.isPending ? (
                    <div className={styles.spinner} />
                  ) : addMutation.isSuccess && addMutation.variables === collection.id ? (
                    <Check size={18} className={styles.checkIcon} />
                  ) : (
                    <Plus size={18} className={styles.plusIcon} />
                  )}
                </button>
              ))}

              {collectionsData?.collections.length === 0 && !showCreate && (
                <p className={styles.empty}>You don't have any collections yet.</p>
              )}
            </div>
          )}
        </div>

        <div className={styles.footer}>
          {showCreate ? (
            <div className={styles.createForm}>
              <input 
                type="text" 
                placeholder="Collection name" 
                className={styles.input}
                value={newCollectionName}
                onChange={e => setNewCollectionName(e.target.value)}
                autoFocus
              />
              <div className={styles.createActions}>
                <button 
                  className={styles.cancelBtn} 
                  onClick={() => setShowCreate(false)}
                >
                  Cancel
                </button>
                <button 
                  className={styles.submitBtn}
                  onClick={() => createMutation.mutate(newCollectionName)}
                  disabled={!newCollectionName.trim() || createMutation.isPending}
                >
                  {createMutation.isPending ? 'Creating...' : 'Create'}
                </button>
              </div>
            </div>
          ) : (
            <button 
              className={styles.addBtn}
              onClick={() => setShowCreate(true)}
            >
              <Plus size={18} />
              <span>Create New Collection</span>
            </button>
          )}
        </div>
      </div>
    </div>
  );
};
