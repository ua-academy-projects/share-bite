import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Button } from '../../components/Button/Button';
import styles from './CreatePost.module.css';
import { clsx } from 'clsx';
import { ImagePlus } from 'lucide-react';
import { useMutation, useQueryClient, useQuery } from '@tanstack/react-query';
import { apiClient } from '../../api/client';

export const CreatePost: React.FC = () => {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [imagePreview, setImagePreview] = useState<string | null>(null);
  const [imageFile, setImageFile] = useState<File | null>(null);
  const [venueId, setVenueId] = useState('');
  const [text, setText] = useState('');
  const [rating, setRating] = useState(5);

  const handleImageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      const file = e.target.files[0];
      setImageFile(file);
      const url = URL.createObjectURL(file);
      setImagePreview(url);
    }
  };

  const createMutation = useMutation({
    mutationFn: async () => {
      return await apiClient.createPost({
        venueId: parseInt(venueId, 10),
        text,
        rating,
        images: imageFile ? [imageFile] : undefined
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['posts'] });
      navigate('/');
    },
    onError: (error) => {
      console.error("Failed to create post:", error);
      // Inline error UI handles displaying the message
    }
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!venueId) {
      alert("Please select a venue");
      return;
    }
    createMutation.mutate();
  };

  return (
    <div className={styles.container}>
      <div className={clsx(styles.card, 'glass-panel')}>
        <h1 className={styles.title}>Create a Post</h1>
        <p className={styles.subtitle}>Share your latest culinary experience</p>

        <form className={styles.form} onSubmit={handleSubmit}>
          
          <div className={styles.imageUploadWrapper}>
            <label htmlFor="image-upload" className={styles.imageUploadLabel}>
              {imagePreview ? (
                <img src={imagePreview} alt="Preview" className={styles.imagePreview} />
              ) : (
                <div className={styles.uploadPlaceholder}>
                  <ImagePlus size={48} className={styles.uploadIcon} />
                  <span>Click to upload a photo</span>
                </div>
              )}
            </label>
            <input 
              id="image-upload" 
              type="file" 
              accept="image/*" 
              className={styles.hiddenInput} 
              onChange={handleImageChange}
            />
          </div>

          <div className={styles.inputGroup}>
            <label htmlFor="venueId" className={styles.label}>Venue ID</label>
            <input 
              id="venueId" 
              type="text"
              className={styles.input} 
              placeholder="Enter venue ID"
              required 
              value={venueId}
              onChange={e => setVenueId(e.target.value)}
            />
          </div>

          <div className={styles.inputGroup}>
            <label htmlFor="rating" className={styles.label}>Rating (1-5)</label>
            <input 
              id="rating" 
              type="number"
              min="1"
              max="5"
              className={styles.input} 
              required 
              value={rating}
              onChange={e => setRating(parseInt(e.target.value, 10))}
            />
          </div>

          <div className={styles.inputGroup}>
            <label htmlFor="review" className={styles.label}>Review</label>
            <textarea 
              id="review" 
              className={clsx(styles.input, styles.textarea)} 
              placeholder="What did you think of the food?"
              required
              rows={4}
              value={text}
              onChange={e => setText(e.target.value)}
            />
          </div>

          {createMutation.isError && (
            <div className={styles.errorMessage}>
              {(createMutation.error as any)?.response?.data?.error || 
                "Failed to create post. Please try again."}
            </div>
          )}

          <div className={styles.actions}>
            <Button type="button" variant="outline" onClick={() => navigate(-1)}>
              Cancel
            </Button>
            <Button type="submit" disabled={createMutation.isPending}>
              {createMutation.isPending ? 'Posting...' : 'Post Review'}
            </Button>
          </div>
        </form>
      </div>
    </div>
  );
};
