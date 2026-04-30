import React, { useState, useRef } from 'react';
import { X, Camera, Loader2 } from 'lucide-react';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { apiClient } from '../../../api/client';
import styles from './EditProfileModal.module.css';
import { clsx } from 'clsx';
import type { Customer } from '../../../types/api';

interface EditProfileModalProps {
  profile: Customer;
  isOpen: boolean;
  onClose: () => void;
}

export const EditProfileModal: React.FC<EditProfileModalProps> = ({ 
  profile, 
  isOpen, 
  onClose 
}) => {
  const queryClient = useQueryClient();
  const fileInputRef = useRef<HTMLInputElement>(null);
  
  const [formData, setFormData] = useState({
    firstName: profile.firstName,
    lastName: profile.lastName,
    userName: profile.userName,
    bio: profile.bio || ''
  });
  
  const [selectedImage, setSelectedImage] = useState<File | null>(null);
  const [imagePreview, setImagePreview] = useState<string | null>(profile.avatarUrl || null);
  const [error, setError] = useState<string | null>(null);

  const updateMutation = useMutation({
    mutationFn: async () => {
      // 1. Update text fields
      await apiClient.updateCustomer(formData);
      
      // 2. Upload avatar if selected
      if (selectedImage) {
        await apiClient.uploadAvatar(selectedImage);
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['currentCustomer'] });
      queryClient.invalidateQueries({ queryKey: ['user', formData.userName] });
      onClose();
    },
    onError: (err: any) => {
      setError(err.response?.data?.error || 'Failed to update profile');
    }
  });

  const handleImageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      if (file.size > 5 * 1024 * 1024) {
        setError('Image size must be less than 5MB');
        return;
      }
      setSelectedImage(file);
      const reader = new FileReader();
      reader.onloadend = () => {
        setImagePreview(reader.result as string);
      };
      reader.readAsDataURL(file);
      setError(null);
    }
  };

  if (!isOpen) return null;

  return (
    <div className={styles.overlay} onClick={onClose}>
      <div className={clsx(styles.modal, 'glass-panel')} onClick={e => e.stopPropagation()}>
        <div className={styles.header}>
          <h2 className={styles.title}>Edit Profile</h2>
          <button className={styles.closeBtn} onClick={onClose}>
            <X size={24} />
          </button>
        </div>

        <form 
          className={styles.form} 
          onSubmit={(e) => {
            e.preventDefault();
            updateMutation.mutate();
          }}
        >
          {/* Avatar Upload */}
          <div className={styles.avatarSection}>
            <div className={styles.avatarWrapper} onClick={() => fileInputRef.current?.click()}>
              <img 
                src={imagePreview || 'https://via.placeholder.com/100'} 
                alt="Profile Preview" 
                className={styles.avatarPreview} 
              />
              <div className={styles.avatarOverlay}>
                <Camera size={24} />
              </div>
            </div>
            <input 
              type="file" 
              ref={fileInputRef} 
              className={styles.hiddenInput} 
              accept="image/jpeg,image/png"
              onChange={handleImageChange}
            />
            <p className={styles.avatarHint}>Click to change photo</p>
          </div>

          {error && <div className={styles.errorMsg}>{error}</div>}

          <div className={styles.inputGroup}>
            <label className={styles.label}>Username</label>
            <input 
              type="text" 
              className={styles.input}
              value={formData.userName}
              onChange={e => setFormData({ ...formData, userName: e.target.value })}
              required
            />
          </div>

          <div className={styles.row}>
            <div className={styles.inputGroup}>
              <label className={styles.label}>First Name</label>
              <input 
                type="text" 
                className={styles.input}
                value={formData.firstName}
                onChange={e => setFormData({ ...formData, firstName: e.target.value })}
                required
              />
            </div>
            <div className={styles.inputGroup}>
              <label className={styles.label}>Last Name</label>
              <input 
                type="text" 
                className={styles.input}
                value={formData.lastName}
                onChange={e => setFormData({ ...formData, lastName: e.target.value })}
                required
              />
            </div>
          </div>

          <div className={styles.inputGroup}>
            <label className={styles.label}>Bio</label>
            <textarea 
              className={clsx(styles.input, styles.textarea)}
              value={formData.bio}
              onChange={e => setFormData({ ...formData, bio: e.target.value })}
              maxLength={500}
              placeholder="Tell us about yourself..."
            />
            <span className={styles.charCount}>{formData.bio.length}/500</span>
          </div>

          <div className={styles.footer}>
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
              {updateMutation.isPending ? (
                <>
                  <Loader2 size={18} className={styles.spinner} />
                  <span>Saving...</span>
                </>
              ) : 'Save Changes'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};
