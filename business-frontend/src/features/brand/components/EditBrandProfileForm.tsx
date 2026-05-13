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
    const updated = await businessApi.uploadAvatar(brand.id, file, token);
    setCurrentBrand(updated);
  };

  const handleUploadBanner = async (file: File) => {
    const token = localStorage.getItem("token") || "";
    const updated = await businessApi.uploadBanner(brand.id, file, token);
    setCurrentBrand(updated);
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

  const inputClass = "bg-[#163d32] border-white/5 text-white focus-visible:ring-emerald-500 rounded-xl placeholder:text-gray-500 h-11 transition-all focus:bg-[#1a4a3d]";

  if (showConfirmation) {
    return (
      <div className="space-y-6 py-4 animate-in fade-in zoom-in duration-300">
        <div className="flex flex-col items-center text-center gap-4">
          <div className="h-16 w-16 rounded-full bg-yellow-500/10 flex items-center justify-center">
            <AlertTriangle className="h-8 w-8 text-yellow-500" />
          </div>
          <div className="space-y-2">
            <h3 className="text-xl font-bold text-white">Review Changes</h3>
            <p className="text-[#cbd5cf] max-w-md">
              Are you sure you want to publish these updates? The changes will be visible to all customers immediately.
            </p>
          </div>
        </div>

        <div className="rounded-2xl bg-white/5 p-6 border border-white/10 space-y-3">
          <div className="flex justify-between text-sm">
            <span className="text-gray-400">Update Type</span>
            <span className="text-white font-medium">Public Profile</span>
          </div>
          <div className="flex justify-between text-sm">
            <span className="text-gray-400">Visibility</span>
            <span className="text-[#98FF98] font-medium flex items-center gap-1">
              <CheckCircle2 className="h-3 w-3" /> Live
            </span>
          </div>
        </div>

        <div className="flex flex-col gap-3 pt-4">
          <Button
            onClick={form.handleSubmit(onSubmit)}
            disabled={loading}
            className="w-full bg-[#98FF98] text-[#0d241d] hover:bg-[#7cfc7c] font-bold h-12 rounded-xl text-lg shadow-lg shadow-[#98FF98]/10"
          >
            {loading ? <Loader2 className="h-5 w-5 animate-spin" /> : "Confirm & Save"}
          </Button>
          <Button
            type="button"
            variant="ghost"
            onClick={() => setShowConfirmation(false)}
            disabled={loading}
            className="w-full text-[#cbd5cf] hover:text-white hover:bg-white/5 h-12 rounded-xl"
          >
            Go Back
          </Button>
        </div>
      </div>
    );
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
          <div className="space-y-6">
            <FormField
              control={form.control}
              name="name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel className="text-[#cbd5cf]">Brand Name</FormLabel>
                  <FormControl>
                    <Input className={inputClass} {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="description"
              render={({ field }) => (
                <FormItem>
                  <FormLabel className="text-[#cbd5cf]">Description</FormLabel>
                  <FormControl>
                    <Textarea 
                      className={cn(inputClass, "min-h-[150px] py-3 resize-none")} 
                      {...field} 
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
          </div>

          <div className="space-y-6">
            <ImageUploadField
              label="Avatar"
              value={currentBrand.avatar}
              onUpload={handleUploadAvatar}
              aspectRatio="square"
            />

            <ImageUploadField
              label="Banner"
              value={currentBrand.banner}
              onUpload={handleUploadBanner}
              aspectRatio="video"
            />
          </div>
        </div>

        {error && (
          <div className="rounded-xl border border-red-500/30 bg-red-500/10 px-4 py-3 text-sm text-red-400 font-medium flex items-center gap-2">
            <AlertTriangle className="h-4 w-4" />
            {error}
          </div>
        )}

        <div className="flex justify-end gap-3 pt-6 border-t border-white/5">
          <Button
            type="button"
            variant="ghost"
            onClick={onCancel}
            disabled={loading}
            className="text-[#cbd5cf] hover:text-white hover:bg-white/5 h-11 px-6 rounded-xl"
          >
            Cancel
          </Button>
          <Button
            type="submit"
            disabled={loading}
            className="bg-[#98FF98] text-[#0d241d] hover:bg-[#7cfc7c] font-bold px-8 h-11 rounded-xl shadow-lg shadow-[#98FF98]/10"
          >
            Review Changes
          </Button>
        </div>
      </form>
    </Form>
  );
}

