import React, { useState, useEffect } from 'react';
import styles from './Explore.module.css';
import { Link } from 'react-router-dom';
import { clsx } from 'clsx';
import { Star, MapPin } from 'lucide-react';
import { useQuery } from '@tanstack/react-query';
import { apiClient } from '../../api/client';

export const Explore: React.FC = () => {
  const [coords, setCoords] = useState<{lat: number, lng: number} | null>(null);

  useEffect(() => {
    if ('geolocation' in navigator) {
      navigator.geolocation.getCurrentPosition(
        (position) => {
          setCoords({ lat: position.coords.latitude, lng: position.coords.longitude });
        },
        (error) => {
          console.warn('Geolocation error:', error);
          setCoords({ lat: 50.4501, lng: 30.5234 }); // fallback to Kyiv
        }
      );
    } else {
      setCoords({ lat: 50.4501, lng: 30.5234 }); // fallback
    }
  }, []);

  const { data: venues, isLoading, error } = useQuery({
    queryKey: ['explore', coords?.lat, coords?.lng],
    queryFn: () => apiClient.getExploreNearby(coords!.lat, coords!.lng, 20),
    enabled: !!coords
  });

  return (
    <div className={styles.container}>
      <header className={styles.header}>
        <h1 className={styles.title}>Explore</h1>
        <p className={styles.subtitle}>Discover trending places around you</p>
      </header>

      {isLoading ? (
        <div className={styles.loading}>Finding venues...</div>
      ) : error ? (
        <div className={styles.error}>Failed to load venues.</div>
      ) : (
        <div className={styles.grid}>
          {venues?.map(item => (
            <Link key={item.venue_id} to={`/restaurant/${item.venue_id}`} className={clsx(styles.card, 'glass-panel')}>
              <div className={styles.imageContainer}>
                <img src={item.posts[0]?.images[0] || 'https://images.unsplash.com/photo-1514933651103-005eec06c04b'} alt={`Venue ${item.venue_id}`} className={styles.image} />
              </div>
              <div className={styles.content}>
                <div className={styles.row}>
                  <h2 className={styles.name}>Venue #{item.venue_id}</h2>
                  <div className={styles.rating}>
                    <Star size={16} fill="currentColor" className={styles.starIcon} />
                    <span>--</span>
                  </div>
                </div>
                <p className={styles.category}>Explore the latest posts</p>
                <div className={styles.location}>
                  <MapPin size={16} className={styles.locationIcon} />
                  <span>{item.posts.length} posts here</span>
                </div>
              </div>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
};
