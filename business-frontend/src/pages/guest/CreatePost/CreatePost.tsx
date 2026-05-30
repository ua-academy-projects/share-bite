import React, { useState, useEffect, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
import { ImagePlus, ChevronLeft, ChevronRight, MapPin } from 'lucide-react';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { apiClient } from '@/api/client';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';

const ALLOWED_IMAGE_TYPES = ['image/jpeg', 'image/png', 'image/jpg'];

const isValidImageType = (file: File) => ALLOWED_IMAGE_TYPES.includes(file.type);

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
      setValidationError(null);
      const files = Array.from(e.target.files);
      
      const invalidFiles = files.filter(file => !isValidImageType(file));
      if (invalidFiles.length > 0) {
        setValidationError('Unsupported image type. Only JPEG and PNG are supported.');
        e.currentTarget.value = '';
        return;
      }

      const allowedFiles = files.slice(0, 5 - images.length);
      const newImages = allowedFiles.map(file => ({
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
    mutationFn: async ({ finalVenueId, parsedRating }: { finalVenueId: number; parsedRating: number }) => {
      return await apiClient.createPost({
        venueId: finalVenueId,
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
    }
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setValidationError(null);

    let finalVenueId = parseInt(venueId.trim(), 10);
    if (isNaN(finalVenueId) || finalVenueId <= 0) {
      // Require a valid numeric ID instead of silent fallback
      setValidationError("Please enter a valid numeric venue ID.");
      return;
    }

    const parsedRating = parseInt(rating.toString(), 10);
    if (Number.isNaN(parsedRating) || parsedRating < 1 || parsedRating > 5) {
      setValidationError('Rating must be between 1 and 5');
      return;
    }

    const invalidImage = images.find(img => !['image/jpeg', 'image/png', 'image/jpg'].includes(img.file.type));
    if (invalidImage) {
      setValidationError('Unsupported image type. Only JPEG and PNG are supported.');
      return;
    }

    createMutation.mutate({ finalVenueId, parsedRating });
  };

  return (
    <div className="max-w-3xl mx-auto pt-24 pb-12 px-4">
      <div className="bg-card text-card-foreground border border-border p-8 md:p-10 rounded-3xl shadow-2xl relative z-10 backdrop-blur-sm">
        <div className="mb-8">
          <h1 className="text-4xl font-serif font-bold text-foreground mb-2">Create a Post</h1>
          <p className="text-muted-foreground">Share your latest culinary experience</p>
        </div>

        <form className="space-y-6" onSubmit={handleSubmit}>
          
          <div className="space-y-4">
            <Label className="text-lg">Photos</Label>
            <div className="grid grid-cols-2 md:grid-cols-5 gap-3">
              {images.map((img, idx) => (
                <div key={img.previewUrl} className="relative aspect-square rounded-xl overflow-hidden group bg-muted border border-border">
                  <img src={img.previewUrl} alt={`Preview ${idx}`} className="w-full h-full object-cover" />
                  <div className="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center gap-2 backdrop-blur-sm">
                    {idx > 0 && (
                      <button type="button" className="p-1.5 bg-background/80 rounded-full hover:bg-background text-foreground" onClick={() => moveImage(idx, 'left')}>
                        <ChevronLeft size={16} />
                      </button>
                    )}
                    {idx < images.length - 1 && (
                      <button type="button" className="p-1.5 bg-background/80 rounded-full hover:bg-background text-foreground" onClick={() => moveImage(idx, 'right')}>
                        <ChevronRight size={16} />
                      </button>
                    )}
                  </div>
                  <button type="button" className="absolute top-2 right-2 p-1 bg-destructive/80 hover:bg-destructive text-destructive-foreground rounded-full shadow-sm" onClick={() => removeImage(idx)}>
                    &times;
                  </button>
                </div>
              ))}
              
              {images.length < 5 && (
                <label htmlFor="image-upload" className="aspect-square rounded-xl border-2 border-dashed border-border/60 bg-muted/30 hover:bg-muted/50 transition-colors flex flex-col items-center justify-center cursor-pointer text-muted-foreground hover:text-foreground">
                  <ImagePlus size={28} className="mb-2" />
                  <span className="text-xs font-semibold">Add Photo</span>
                  <input 
                    id="image-upload" 
                    type="file" 
                    accept="image/png, image/jpeg, image/jpg" 
                    multiple
                    className="hidden" 
                    onChange={handleImageChange}
                  />
                </label>
              )}
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="space-y-2">
              <Label htmlFor="venueId" className="text-sm font-semibold text-muted-foreground flex items-center gap-1">
                <MapPin size={16} /> Venue Selection
              </Label>
              <Input 
                id="venueId" 
                type="text"
                className="bg-muted/50 border-border h-12 rounded-xl focus-visible:ring-primary" 
                placeholder="Search or enter UUID..."
                value={venueId}
                onChange={e => setVenueId(e.target.value)}
              />
              <p className="text-xs text-muted-foreground/60 pl-1">Leave blank to use default test venue.</p>
            </div>

            <div className="space-y-2">
              <Label htmlFor="rating" className="text-sm font-semibold text-muted-foreground">Rating (1-5)</Label>
              <Input 
                id="rating" 
                type="number"
                min="1"
                max="5"
                className="bg-muted/50 border-border h-12 rounded-xl focus-visible:ring-primary" 
                required 
                value={rating}
                onChange={e => setRating(parseInt(e.target.value, 10))}
              />
            </div>
          </div>

          <div className="space-y-2">
            <Label htmlFor="review" className="text-sm font-semibold text-muted-foreground">Review</Label>
            <textarea 
              id="review" 
              className="flex w-full rounded-xl border border-input bg-muted/50 px-3 py-3 text-sm shadow-sm transition-colors file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50 min-h-[120px] resize-y" 
              placeholder="What did you think of the food and atmosphere?"
              required
              rows={5}
              value={text}
              onChange={e => setText(e.target.value)}
            />
          </div>

          {validationError && (
            <div className="p-4 bg-destructive/10 text-destructive text-sm font-bold rounded-xl border border-destructive/20">
              {validationError}
            </div>
          )}

          {createMutation.isError && (
            <div className="p-4 bg-destructive/10 text-destructive text-sm font-bold rounded-xl border border-destructive/20">
              {(createMutation.error as any)?.response?.data?.error || "Failed to create post. Please try again."}
            </div>
          )}

          <div className="flex justify-end gap-3 pt-6 mt-4 border-t border-border/50">
            <Button type="button" variant="ghost" onClick={() => navigate(-1)} className="h-12 px-6 rounded-full font-bold">
              Cancel
            </Button>
            <Button type="submit" disabled={createMutation.isPending} className="h-12 px-8 rounded-full bg-accent text-accent-foreground hover:bg-[#c9a900] font-bold shadow-lg">
              {createMutation.isPending ? 'Posting...' : 'Post Review'}
            </Button>
          </div>
        </form>
      </div>
    </div>
  );
};
