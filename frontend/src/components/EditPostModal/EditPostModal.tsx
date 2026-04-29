import React, { useState, useEffect, useRef } from 'react';
import { X, ImagePlus, ChevronLeft, ChevronRight } from 'lucide-react';
import { clsx } from 'clsx';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { apiClient } from '../../api/client';
import type { PostResponse } from '../../types/api';
import styles from './EditPostModal.module.css';

interface EditPostModalProps {
  post: PostResponse;
  isOpen: boolean;
  onClose: () => void;
}

interface GalleryItem {
  id: string;
  url: string;
  file?: File;
}

export const EditPostModal: React.FC<EditPostModalProps> = ({ post, isOpen, onClose }) => {
  const queryClient = useQueryClient();
  const [text, setText] = useState(post.text);
  const [rating, setRating] = useState(post.rating);
  const [venueId, setVenueId] = useState(post.venueId.toString());
  const [status, setStatus] = useState(post.status);
  const [gallery, setGallery] = useState<GalleryItem[]>(
    (post.images || []).map((url, idx) => ({ id: `old-${idx}`, url }))
  );
  const [isImagesModified, setIsImagesModified] = useState(false);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);

  const galleryRef = useRef(gallery);
  useEffect(() => {
    galleryRef.current = gallery;
  }, [gallery]);

  useEffect(() => {
    return () => {
      galleryRef.current.forEach(item => {
        if (item.file) {
          URL.revokeObjectURL(item.url);
        }
      });
    };
  }, []);

  const urlToFile = async (url: string, filename: string): Promise<File> => {
    const res = await fetch(url);
    const blob = await res.blob();
    return new File([blob], filename, { type: blob.type });
  };

  const updateMutation = useMutation({
    mutationFn: async () => {
      let filesToSend: File[] | undefined = undefined;

      if (isImagesModified) {
        filesToSend = [];
        let conversionFailed = false;
        for (let i = 0; i < gallery.length; i++) {
          const item = gallery[i];
          if (item.file) {
            filesToSend.push(item.file);
          } else {
            try {
              const file = await urlToFile(item.url, `image-${i}.jpg`);
              filesToSend.push(file);
            } catch (err) {
              console.error('Failed to convert URL to file:', item.url, err);
              conversionFailed = true;
            }
          }
        }
        if (conversionFailed) {
          setErrorMsg('Failed to process some existing images due to CORS issues. Please remove them to proceed.');
          throw new Error('Failed to process images');
        }
      }

      return await apiClient.updatePost(post.id, {
        text,
        rating,
        venueId: parseInt(venueId, 10),
        status,
        images: filesToSend
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['posts'] });
      onClose();
    },
    onError: (err: any) => {
      if (err.message !== 'Failed to process images') {
        setErrorMsg(err.response?.data?.error || err.response?.data?.message || 'Failed to update post.');
      }
    }
  });

  const deleteMutation = useMutation({
    mutationFn: async () => {
      await apiClient.deletePost(post.id);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['posts'] });
      onClose();
    },
    onError: (err: any) => {
      setErrorMsg(err.response?.data?.error || err.response?.data?.message || 'Failed to delete post.');
    }
  });

  const handleImageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files.length > 0) {
      setIsImagesModified(true);
      const files = Array.from(e.target.files);
      const newItems = files.map(file => ({
        id: `new-${Math.random().toString(36).substr(2, 9)}`,
        url: URL.createObjectURL(file),
        file
      }));
      setGallery(prev => [...prev, ...newItems].slice(0, 5));
      e.currentTarget.value = '';
    }
  };

  const removeImage = (id: string) => {
    setIsImagesModified(true);
    setGallery(prev => {
      const itemToRemove = prev.find(item => item.id === id);
      if (itemToRemove?.file) {
        URL.revokeObjectURL(itemToRemove.url);
      }
      return prev.filter(item => item.id !== id);
    });
  };

  const clearImages = () => {
    setIsImagesModified(true);
    setGallery(prev => {
      prev.forEach(item => {
        if (item.file) URL.revokeObjectURL(item.url);
      });
      return [];
    });
  };

  const moveImage = (index: number, direction: 'left' | 'right') => {
    setIsImagesModified(true);
    const newGallery = [...gallery];
    if (direction === 'left' && index > 0) {
      [newGallery[index - 1], newGallery[index]] = [newGallery[index], newGallery[index - 1]];
    } else if (direction === 'right' && index < gallery.length - 1) {
      [newGallery[index], newGallery[index + 1]] = [newGallery[index + 1], newGallery[index]];
    }
    setGallery(newGallery);
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setErrorMsg(null);
    if (!text.trim()) {
      setErrorMsg('Review text cannot be empty');
      return;
    }
    const parsedVenueId = parseInt(venueId, 10);
    if (isNaN(parsedVenueId)) {
      setErrorMsg('Invalid Venue ID');
      return;
    }
    updateMutation.mutate();
  };

  const handleDelete = () => {
    if (window.confirm('Are you sure you want to delete this post?')) {
      deleteMutation.mutate();
    }
  };

  return (
    <div className={styles.overlay} onClick={onClose}>
      <div className={clsx(styles.modal, 'glass-panel')} onClick={e => e.stopPropagation()}>
        <button className={styles.closeBtn} onClick={onClose} aria-label="Close modal">
          <X size={24} />
        </button>
        
        <div className={styles.headerRow}>
          <h2 className={styles.title}>Edit Post</h2>
          <button 
            type="button"
            className={styles.deleteBtn}
            onClick={handleDelete}
            disabled={deleteMutation.isPending || updateMutation.isPending}
          >
            Delete Post
          </button>
        </div>
        
        <form className={styles.form} onSubmit={handleSubmit}>
          {/* Image Gallery */}
          <div className={styles.imageSectionHeader}>
            <label className={styles.label}>Photos</label>
            {gallery.length > 0 && (
              <button 
                type="button" 
                className={styles.clearImagesBtn}
                onClick={clearImages}
              >
                Clear & Replace Gallery
              </button>
            )}
          </div>

          <div className={styles.imagePreviews}>
            {gallery.map((item, idx) => (
              <div key={item.id} className={styles.previewWrapper}>
                <img src={item.url} alt={`Preview ${idx}`} className={styles.imagePreview} />
                <div className={styles.previewOverlay}>
                  {idx > 0 && (
                    <button type="button" className={styles.moveImgBtn} onClick={() => moveImage(idx, 'left')}>
                      <ChevronLeft size={16} />
                    </button>
                  )}
                  {idx < gallery.length - 1 && (
                    <button type="button" className={styles.moveImgBtn} onClick={() => moveImage(idx, 'right')}>
                      <ChevronRight size={16} />
                    </button>
                  )}
                </div>
                <button type="button" className={styles.removeImgBtn} onClick={() => removeImage(item.id)}>&times;</button>
              </div>
            ))}
          </div>
            
          {gallery.length < 5 && (
            <div className={styles.imageUploadWrapper}>
              <label htmlFor={`edit-image-upload-${post.id}`} className={styles.imageUploadLabel}>
                <div className={styles.uploadPlaceholder}>
                  <ImagePlus size={32} className={styles.uploadIcon} />
                  <span>Add photo</span>
                </div>
              </label>
              <input 
                id={`edit-image-upload-${post.id}`} 
                type="file" 
                accept="image/*" 
                multiple
                className={styles.hiddenInput} 
                onChange={handleImageChange}
              />
            </div>
          )}
          

          {/* Input Fields */}
          <div className={styles.inputGroup}>
            <label className={styles.label}>Review</label>
            <textarea 
              className={clsx(styles.input, styles.textarea)} 
              rows={4}
              value={text}
              onChange={e => setText(e.target.value)}
              placeholder="Edit your review..."
            />
          </div>

          <div className={styles.row}>
            <div className={styles.inputGroup}>
              <label className={styles.label}>Rating</label>
              <select 
                className={styles.input}
                value={rating}
                onChange={e => setRating(parseInt(e.target.value, 10))}
              >
                {[1, 2, 3, 4, 5].map(num => (
                  <option key={num} value={num}>{num} ★</option>
                ))}
              </select>
            </div>

            <div className={styles.inputGroup}>
              <label className={styles.label}>Venue ID</label>
              <input 
                type="text"
                className={styles.input}
                value={venueId}
                onChange={e => setVenueId(e.target.value)}
              />
            </div>
          </div>

          <div className={styles.inputGroup}>
            <label className={styles.label}>Status</label>
            <select 
              className={styles.input}
              value={status}
              onChange={e => setStatus(e.target.value as any)}
            >
              <option value="draft">Draft</option>
              <option value="published">Published</option>
              <option value="archived">Archived</option>
            </select>
          </div>

          {errorMsg && (
            <div className={styles.errorMessage}>
              {errorMsg}
            </div>
          )}

          <div className={styles.actions}>
            <button 
              type="button" 
              className={styles.cancelBtn} 
              onClick={onClose}
              disabled={updateMutation.isPending}
            >
              Cancel
            </button>
            <button 
              type="submit" 
              className={styles.saveBtn}
              disabled={updateMutation.isPending}
            >
              {updateMutation.isPending ? 'Saving...' : 'Save Changes'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};
