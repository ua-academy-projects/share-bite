import React, { useState, useEffect, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
import { Button } from '../../components/Button/Button';
import styles from './CreatePost.module.css';
import { clsx } from 'clsx';
import { ImagePlus, ChevronLeft, ChevronRight } from 'lucide-react';
import { useMutation, useQueryClient, useQuery } from '@tanstack/react-query';
import { apiClient } from '../../api/client';

export const CreatePost: React.FC = () => {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [images, setImages] = useState<{ file: File; previewUrl: string }[]>([]);
  const [venueId, setVenueId] = useState('');
  const [text, setText] = useState('');
  const [rating, setRating] = useState(5);
  const [validationError, setValidationError] = useState<string | null>(null);

  const imagesRef = useRef(images);
  useEffect(() => {
    imagesRef.current = images;
  }, [images]);

  useEffect(() => {
    return () => {
      imagesRef.current.forEach(img => URL.revokeObjectURL(img.previewUrl));
    };
  }, []);

  const handleImageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files.length > 0) {
      const files = Array.from(e.target.files).slice(0, 5 - images.length);
      const newImages = files.map(file => ({
        file,
        previewUrl: URL.createObjectURL(file)
      }));
      setImages(prev => [...prev, ...newImages].slice(0, 5));
      e.currentTarget.value = '';
    }
  };

  const removeImage = (index: number) => {
    setImages(prev => {
      const newImages = [...prev];
      URL.revokeObjectURL(newImages[index].previewUrl);
      newImages.splice(index, 1);
      return newImages;
    });
  };

  const moveImage = (index: number, direction: 'left' | 'right') => {
    setImages(prev => {
      const newImages = [...prev];
      if (direction === 'left' && index > 0) {
        [newImages[index - 1], newImages[index]] = [newImages[index], newImages[index - 1]];
      } else if (direction === 'right' && index < prev.length - 1) {
        [newImages[index], newImages[index + 1]] = [newImages[index + 1], newImages[index]];
      }
      return newImages;
    });
  };

  const createMutation = useMutation({
    mutationFn: async ({ parsedVenueId, parsedRating }: { parsedVenueId: number; parsedRating: number }) => {
      return await apiClient.createPost({
        venueId: parsedVenueId,
        text,
        rating: parsedRating,
        images: images.length > 0 ? images.map(img => img.file) : undefined
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
    setValidationError(null);

    const parsedVenueId = parseInt(venueId, 10);
    if (Number.isNaN(parsedVenueId)) {
      setValidationError('Please enter a valid Venue ID');
      return;
    }
    const parsedRating = parseInt(rating.toString(), 10);
    if (Number.isNaN(parsedRating) || parsedRating < 1 || parsedRating > 5) {
      setValidationError('Rating must be between 1 and 5');
      return;
    }
    createMutation.mutate({ parsedVenueId, parsedRating });
  };

  return (
    <div className={styles.container}>
      <div className={clsx(styles.card, 'glass-panel')}>
        <h1 className={styles.title}>Create a Post</h1>
        <p className={styles.subtitle}>Share your latest culinary experience</p>

        <form className={styles.form} onSubmit={handleSubmit}>
          
          <div className={styles.imageUploadContainer}>
            <div className={styles.imagePreviews}>
              {images.map((img, idx) => (
                <div key={img.previewUrl} className={styles.previewWrapper}>
                  <img src={img.previewUrl} alt={`Preview ${idx}`} className={styles.imagePreview} />
                  <div className={styles.previewOverlay}>
                    {idx > 0 && (
                      <button type="button" className={styles.moveImgBtn} onClick={() => moveImage(idx, 'left')}>
                        <ChevronLeft size={16} />
                      </button>
                    )}
                    {idx < images.length - 1 && (
                      <button type="button" className={styles.moveImgBtn} onClick={() => moveImage(idx, 'right')}>
                        <ChevronRight size={16} />
                      </button>
                    )}
                  </div>
                  <button type="button" className={styles.removeImgBtn} onClick={() => removeImage(idx)}>&times;</button>
                </div>
              ))}
            </div>
            
            {images.length < 5 && (
              <div className={styles.imageUploadWrapper}>
                <label htmlFor="image-upload" className={styles.imageUploadLabel}>
                  <div className={styles.uploadPlaceholder}>
                    <ImagePlus size={32} className={styles.uploadIcon} />
                    <span>Upload photo</span>
                  </div>
                </label>
                <input 
                  id="image-upload" 
                  type="file" 
                  accept="image/*" 
                  multiple
                  className={styles.hiddenInput} 
                  onChange={handleImageChange}
                />
              </div>
            )}
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

          {validationError && (
            <div className={styles.errorMessage}>
              {validationError}
            </div>
          )}

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
