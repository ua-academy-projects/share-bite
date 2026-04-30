import React from 'react';
import { useParams, Link, useNavigate } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { apiClient } from '../../api/client';
import styles from './CollectionDetail.module.css';
import { ArrowLeft, MapPin, ChevronUp, ChevronDown, Trash2, ExternalLink } from 'lucide-react';
import { clsx } from 'clsx';
import { useCurrentCustomer } from '../../hooks/useCurrentCustomer';

export const CollectionDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { data: currentCustomer } = useCurrentCustomer();

  const { data: collection, isLoading: isCollectionLoading } = useQuery({
    queryKey: ['collection', id],
    queryFn: () => apiClient.getCollection(id!),
    enabled: !!id
  });

  const { data: venues, isLoading: isVenuesLoading } = useQuery({
    queryKey: ['collectionVenues', id],
    queryFn: () => apiClient.listVenuesInCollection(id!),
    enabled: !!id
  });

  const reorderMutation = useMutation({
    mutationFn: ({ venueId, prevVenueId, nextVenueId }: { venueId: number; prevVenueId?: number; nextVenueId?: number }) => 
      apiClient.reorderVenueInCollection(id!, venueId, { prevVenueId, nextVenueId }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['collectionVenues', id] });
    }
  });

  const removeMutation = useMutation({
    mutationFn: (venueId: number) => apiClient.removeVenueFromCollection(id!, venueId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['collectionVenues', id] });
    }
  });

  const isOwner = currentCustomer?.id && collection?.id; // Simplified owner check

  const handleMoveUp = (index: number) => {
    if (index === 0 || !venues) return;
    const venueId = venues[index].id;
    const targetPrevId = index > 1 ? venues[index - 2].id : undefined;
    const targetNextId = venues[index - 1].id;
    
    reorderMutation.mutate({ 
      venueId, 
      prevVenueId: targetPrevId,
      nextVenueId: targetNextId
    });
  };

  const handleMoveDown = (index: number) => {
    if (!venues || index === venues.length - 1) return;
    const venueId = venues[index].id;
    const targetPrevId = venues[index + 1].id;
    const targetNextId = index < venues.length - 2 ? venues[index + 2].id : undefined;

    reorderMutation.mutate({ 
      venueId, 
      prevVenueId: targetPrevId,
      nextVenueId: targetNextId
    });
  };

  if (isCollectionLoading) return <div className={styles.loading}>Loading collection...</div>;
  if (!collection) return <div className={styles.notFound}>Collection not found</div>;

  return (
    <div className={styles.container}>
      <header className={styles.header}>
        <button className={styles.backBtn} onClick={() => navigate(-1)}>
          <ArrowLeft size={20} />
          <span>Back</span>
        </button>
        <div className={styles.headerContent}>
          <h1 className={styles.title}>{collection.name}</h1>
          {collection.description && <p className={styles.description}>{collection.description}</p>}
          <div className={styles.meta}>
            <span className={styles.tag}>{collection.isPublic ? 'Public' : 'Private'}</span>
            <span className={styles.count}>{venues?.length || 0} venues</span>
          </div>
        </div>
      </header>

      <div className={styles.venueList}>
        {isVenuesLoading ? (
          <p>Loading venues...</p>
        ) : venues && venues.length > 0 ? (
          venues.map((venue, index) => (
            <div key={venue.id} className={clsx(styles.venueItem, 'glass-panel')}>
              <div className={styles.venueMain}>
                <img 
                  src={venue.avatarUrl || 'https://via.placeholder.com/60'} 
                  alt={venue.name} 
                  className={styles.venueAvatar} 
                />
                <div className={styles.venueInfo}>
                  <h3 className={styles.venueName}>{venue.name}</h3>
                  {venue.description && <p className={styles.venueDesc}>{venue.description}</p>}
                  <div className={styles.venueActions}>
                    <Link to={`/restaurant/${venue.id}`} className={styles.venueLink}>
                      <ExternalLink size={16} />
                      <span>View Profile</span>
                    </Link>
                  </div>
                </div>
              </div>

              {isOwner && (
                <div className={styles.reorderControls}>
                  <button 
                    className={styles.reorderBtn} 
                    onClick={() => handleMoveUp(index)}
                    disabled={index === 0 || reorderMutation.isPending}
                    title="Move up"
                  >
                    <ChevronUp size={20} />
                  </button>
                  <button 
                    className={styles.reorderBtn} 
                    onClick={() => handleMoveDown(index)}
                    disabled={index === venues.length - 1 || reorderMutation.isPending}
                    title="Move down"
                  >
                    <ChevronDown size={20} />
                  </button>
                  <button 
                    className={clsx(styles.reorderBtn, styles.removeBtn)} 
                    onClick={() => removeMutation.mutate(venue.id)}
                    disabled={removeMutation.isPending}
                    title="Remove from collection"
                  >
                    <Trash2 size={18} />
                  </button>
                </div>
              )}
            </div>
          ))
        ) : (
          <div className={styles.emptyState}>
            <MapPin size={48} className={styles.emptyIcon} />
            <p>No venues in this collection yet.</p>
            <Link to="/explore" className={styles.exploreBtn}>Explore Venues</Link>
          </div>
        )}
      </div>
    </div>
  );
};
