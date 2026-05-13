import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { Loader2, AlertTriangle, CheckCircle2 } from "lucide-react";

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
};

export function EditBrandProfileForm({ brand, onSuccess, onCancel }: EditBrandProfileFormProps) {
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
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to upload avatar");
    }
  };

  const handleUploadBanner = async (file: File) => {
    const token = localStorage.getItem("token") || "";
    try {
      const updated = await businessApi.uploadBanner(brand.id, file, token);
      setCurrentBrand(updated);
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
          <div className="h-20 w-20 rounded-full bg-emerald-500/10 flex items-center justify-center border border-emerald-500/20">
            <CheckCircle2 className="h-10 w-10 text-[#98FF98]" />
          </div>
          <div className="space-y-2">
            <h3 className="text-2xl font-black text-white tracking-tight">Ready to publish?</h3>
            <p className="text-[#cbd5cf] max-w-md text-lg">
              Review your changes below. Once confirmed, your updated brand identity will be live for all guests.
            </p>
          </div>
        </div>

        <div className="overflow-hidden rounded-[32px] border border-white/10 bg-white/5 backdrop-blur-xl">
          <div className="p-8 space-y-6">
            <div className="grid grid-cols-2 gap-8">
              <div className="space-y-1">
                <span className="text-xs font-black uppercase tracking-widest text-[#9fb2a7]">Brand Identity</span>
                <p className="text-xl font-bold text-white truncate">{form.getValues().name}</p>
              </div>
              <div className="space-y-1">
                <span className="text-xs font-black uppercase tracking-widest text-[#9fb2a7]">Status</span>
                <span className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-[#98FF98]/10 text-[#98FF98] text-xs font-bold border border-[#98FF98]/20">
                  <div className="h-1.5 w-1.5 rounded-full bg-[#98FF98] animate-pulse" />
                  Live Sync
                </span>
              </div>
            </div>
            
            <div className="space-y-1">
              <span className="text-xs font-black uppercase tracking-widest text-[#9fb2a7]">Public Bio</span>
              <p className="text-[#cbd5cf] leading-relaxed line-clamp-3">
                {form.getValues().description || "No description provided."}
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
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-10">
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-10">
          {/* Identity Section */}
          <div className="lg:col-span-7 space-y-8">
            <div className="space-y-2">
              <h4 className="text-xs font-black uppercase tracking-[0.2em] text-[#98FF98]">Core Identity</h4>
              <p className="text-sm text-[#9fb2a7]">Define how your brand appears in search and lists.</p>
            </div>
            
            <div className="space-y-6">
              <FormField
                control={form.control}
                name="name"
                render={({ field }) => (
                  <FormItem className="space-y-3">
                    <FormLabel className="text-sm font-bold text-[#F9F7F2]">Brand Name</FormLabel>
                    <FormControl>
                      <Input className={inputClass} placeholder="e.g. ShareBite Kitchen" {...field} />
                    </FormControl>
                    <FormMessage className="text-red-400 font-bold text-xs" />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="description"
                render={({ field }) => (
                  <FormItem className="space-y-3">
                    <FormLabel className="text-sm font-bold text-[#F9F7F2]">Public Description</FormLabel>
                    <FormControl>
                      <Textarea 
                        className={cn(inputClass, "min-h-[180px] py-4 resize-none leading-relaxed")} 
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

          {/* Visuals Section */}
          <div className="lg:col-span-5 space-y-8">
            <div className="space-y-2">
              <h4 className="text-xs font-black uppercase tracking-[0.2em] text-[#98FF98]">Visual Presence</h4>
              <p className="text-sm text-[#9fb2a7]">High-quality images increase guest trust.</p>
            </div>

            <div className="space-y-8 p-6 rounded-[32px] bg-white/5 border border-white/10">
              <ImageUploadField
                label="Avatar (1:1)"
                value={currentBrand.avatar}
                onUpload={handleUploadAvatar}
                aspectRatio="square"
              />

              <div className="h-px bg-white/5 w-full" />

              <ImageUploadField
                label="Cover Banner (3:1)"
                value={currentBrand.banner}
                onUpload={handleUploadBanner}
                aspectRatio="video"
              />
            </div>
          </div>
        </div>

        {error && (
          <div className="rounded-2xl border border-red-500/30 bg-red-500/10 px-6 py-4 text-sm text-red-100 font-bold flex items-center gap-3 animate-in slide-in-from-top-2 duration-300">
            <AlertTriangle className="h-5 w-5 text-red-400" />
            {error}
          </div>
        )}

        <div className="flex items-center justify-between pt-8 border-t border-white/5">
          <p className="text-xs text-[#9fb2a7] font-medium hidden md:block">
            All changes are saved as drafts until you publish.
          </p>
          <div className="flex gap-4 w-full md:w-auto">
            <Button
              type="button"
              variant="ghost"
              onClick={onCancel}
              disabled={loading}
              className="flex-1 md:flex-none text-[#cbd5cf] hover:text-white hover:bg-white/5 h-12 px-8 rounded-2xl font-bold"
            >
              Cancel
            </Button>
            <Button
              type="submit"
              disabled={loading}
              className="flex-1 md:flex-none bg-[#98FF98] text-[#0d241d] hover:bg-[#7cfc7c] font-black px-10 h-12 rounded-2xl shadow-lg shadow-[#98FF98]/10 transition-all hover:-translate-y-0.5"
            >
              Review & Publish
            </Button>
          </div>
        </div>
      </form>
    </Form>
  );
}

