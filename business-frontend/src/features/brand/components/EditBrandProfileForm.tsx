import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { Loader2, AlertTriangle, Image as ImageIcon } from "lucide-react";

import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { businessApi, type BrandProfile } from "@/api/business";
import { cn } from "@/lib/utils";
import { ImageUploadField } from "./ImageUploadField";

const editProfileSchema = z.object({
  name: z.string().min(2, "Name must be at least 2 characters."),
  description: z.string().max(500, "Description must be less than 500 characters.").optional(),
});

type EditProfileFormValues = z.infer<typeof editProfileSchema>;

type EditBrandProfileFormProps = {
  brand: BrandProfile;
  onSuccess: (updatedBrand: BrandProfile) => void;
  onCancel: () => void;
  onRefresh?: () => void;
};

export function EditBrandProfileForm({ brand, onSuccess, onCancel, onRefresh }: EditBrandProfileFormProps) {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [showConfirmation, setShowConfirmation] = useState(false);
  const [currentBrand, setCurrentBrand] = useState(brand);

  const form = useForm<EditProfileFormValues>({
    resolver: zodResolver(editProfileSchema),
    defaultValues: {
      name: brand.name,
      description: brand.description ?? "",
    },
  });

  const handleUploadAvatar = async (file: File) => {
    const token = localStorage.getItem("token") || "";
    try {
      const updated = await businessApi.uploadAvatar(brand.id, file, token);
      setCurrentBrand(updated);
      onRefresh?.();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to upload avatar");
    }
  };

  const handleUploadBanner = async (file: File) => {
    const token = localStorage.getItem("token") || "";
    try {
      const updated = await businessApi.uploadBanner(brand.id, file, token);
      setCurrentBrand(updated);
      onRefresh?.();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to upload banner");
    }
  };

  const onSubmit = async (values: EditProfileFormValues) => {
    if (!showConfirmation) {
      setShowConfirmation(true);
      return;
    }

    setLoading(true);
    setError(null);

    const token = localStorage.getItem("token") || "";

    try {
      const updated = await businessApi.updateBrandProfile(brand.id, {
        name: values.name,
        description: values.description,
      }, token);
      
      onSuccess(updated);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to update profile");
      setShowConfirmation(false);
    } finally {
      setLoading(false);
    }
  };

  const inputClass = "bg-[#163d32]/60 border-white/5 text-white focus-visible:ring-emerald-500 rounded-2xl placeholder:text-gray-500 h-12 transition-all focus:bg-[#1a4a3d] border-white/10";

  if (showConfirmation) {
    return (
      <div className="space-y-8 py-6 animate-in fade-in zoom-in-95 duration-500">
        <div className="flex flex-col items-center text-center gap-5">
          <div className="space-y-2">
            <h3 className="text-2xl font-black text-white tracking-tight">Ready to publish?</h3>
            <p className="text-[#9fb2a7] max-w-md text-lg">
              Review your changes below. Once confirmed, your updated brand identity will be live for all guests.
            </p>
          </div>
        </div>

        <div className="overflow-hidden rounded-[40px] border border-white/10 bg-white/5 backdrop-blur-xl shadow-2xl">
          {/* Visual Preview */}
          <div className="relative h-48 w-full bg-[#163d32]/40">
            {currentBrand.banner ? (
              <img src={currentBrand.banner} className="h-full w-full object-cover opacity-60" alt="Banner preview" />
            ) : (
              <div className="h-full w-full bg-gradient-to-br from-[#1f4a3f] to-[#0b0f0e] flex items-center justify-center">
                <ImageIcon className="h-10 w-10 text-white/10" />
              </div>
            )}
            <div className="absolute -bottom-12 left-8">
              <div className="h-24 w-24 rounded-3xl p-1 bg-gradient-to-br from-white/20 to-white/5 backdrop-blur-xl shadow-2xl">
                {currentBrand.avatar ? (
                  <img src={currentBrand.avatar} className="h-full w-full rounded-[22px] object-cover border border-white/10" alt="Avatar preview" />
                ) : (
                  <div className="h-full w-full rounded-[22px] bg-[#0f1b17] border border-white/10" />
                )}
              </div>
            </div>
          </div>

          <div className="p-8 pt-16 space-y-6">
            <div className="space-y-1">
              <span className="text-xs font-black uppercase tracking-widest text-[#9fb2a7]">Brand Identity</span>
              <p className="text-xl font-bold text-white truncate">{form.getValues().name}</p>
            </div>
            
            <div className="space-y-1">
              <span className="text-xs font-black uppercase tracking-widest text-[#9fb2a7]">Public Bio</span>
              <p className="text-[#cbd5cf] leading-relaxed line-clamp-3 italic">
                "{form.getValues().description || "No description provided."}"
              </p>
            </div>
          </div>
        </div>

        <div className="flex flex-col gap-4 pt-4">
          <Button
            onClick={form.handleSubmit(onSubmit)}
            disabled={loading}
            className="w-full bg-[#98FF98] text-[#0d241d] hover:bg-[#7cfc7c] font-black h-14 rounded-2xl text-lg shadow-[0_20px_40px_-10px_rgba(152,255,152,0.2)] transition-all hover:-translate-y-1 active:translate-y-0"
          >
            {loading ? <Loader2 className="h-6 w-6 animate-spin" /> : "Publish Changes"}
          </Button>
          <Button
            type="button"
            variant="ghost"
            onClick={() => setShowConfirmation(false)}
            disabled={loading}
            className="w-full text-[#cbd5cf] hover:text-white hover:bg-white/5 h-14 rounded-2xl font-bold"
          >
            Keep Editing
          </Button>
        </div>
      </div>
    );
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-12">
        <div className="space-y-16">
          {/* Visuals Section */}
          <div className="relative">
            {/* Banner Area */}
            <ImageUploadField
              value={currentBrand.banner}
              onUpload={handleUploadBanner}
              aspectRatio="video"
              className="w-full"
            />
            
            {/* Avatar Area - Overlapping */}
            <div className="absolute -bottom-10 left-8 md:left-12">
              <ImageUploadField
                value={currentBrand.avatar}
                onUpload={handleUploadAvatar}
                aspectRatio="square"
                variant="profile"
                className="w-24 h-24 md:w-32 md:h-32"
              />
            </div>
          </div>

          {/* Identity Section */}
          <div className="px-4 md:px-8 space-y-8">
            <FormField
              control={form.control}
              name="name"
              render={({ field }) => (
                <FormItem className="space-y-2">
                  <FormLabel className="text-xs font-black uppercase tracking-widest text-[#9fb2a7]">Brand Name</FormLabel>
                  <FormControl>
                    <Input className={cn(inputClass, "text-xl h-14 font-bold")} placeholder="Name" {...field} />
                  </FormControl>
                  <FormMessage className="text-red-400 font-bold text-xs" />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="description"
              render={({ field }) => (
                <FormItem className="space-y-2">
                  <FormLabel className="text-xs font-black uppercase tracking-widest text-[#9fb2a7]">Bio</FormLabel>
                  <FormControl>
                    <Textarea 
                      className={cn(inputClass, "min-h-[140px] py-4 resize-none leading-relaxed text-lg")} 
                      placeholder="Tell guests what makes your brand special..."
                      {...field} 
                    />
                  </FormControl>
                  <FormMessage className="text-red-400 font-bold text-xs" />
                </FormItem>
              )}
            />
          </div>
        </div>

        {error && (
          <div className="mx-8 rounded-2xl border border-red-500/30 bg-red-500/10 px-6 py-4 text-sm text-red-100 font-bold flex items-center gap-3">
            <AlertTriangle className="h-5 w-5 text-red-400" />
            {error}
          </div>
        )}

        <div className="flex items-center justify-end gap-4 px-8 pt-8 border-t border-white/5">
          <Button
            type="button"
            variant="ghost"
            onClick={onCancel}
            disabled={loading}
            className="text-[#cbd5cf] hover:text-white hover:bg-white/5 h-12 px-8 rounded-2xl font-bold"
          >
            Discard
          </Button>
          <Button
            type="submit"
            disabled={loading}
            className="bg-[#98FF98] text-[#0d241d] hover:bg-[#7cfc7c] font-black px-12 h-12 rounded-2xl shadow-lg shadow-[#98FF98]/10 transition-all hover:-translate-y-0.5"
          >
            Review Changes
          </Button>
        </div>
      </form>
    </Form>
  );
}
