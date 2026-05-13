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
  FormDescription,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { businessApi, type BrandProfile } from "@/api/business";
import { cn } from "@/lib/utils";

const editProfileSchema = z.object({
  name: z.string().min(2, "Name must be at least 2 characters."),
  description: z.string().max(500, "Description must be less than 500 characters.").optional(),
  avatar: z.string().url("Must be a valid URL").or(z.literal("")).optional(),
  banner: z.string().url("Must be a valid URL").or(z.literal("")).optional(),
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

  const form = useForm<EditProfileFormValues>({
    resolver: zodResolver(editProfileSchema),
    defaultValues: {
      name: brand.name,
      description: brand.description ?? "",
      avatar: brand.avatar ?? "",
      banner: brand.banner ?? "",
    },
  });

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
        avatar: values.avatar || undefined,
        banner: values.banner || undefined,
      }, token);
      
      onSuccess(updated);
    } catch (err) {
      // business error handling: the error message from businessApi is already extracted
      setError(err instanceof Error ? err.message : "Failed to update profile");
      setShowConfirmation(false); // Reset confirmation on error so user can fix and retry
    } finally {
      setLoading(false);
    }
  };

  const inputClass = "bg-[#163d32] border-white/5 text-white focus-visible:ring-emerald-500 rounded-2xl placeholder:text-gray-500 h-12 transition-all focus:bg-[#1a4a3d] text-base";

  if (showConfirmation) {
    return (
      <div className="space-y-8 py-10 animate-in fade-in zoom-in duration-300 max-w-2xl mx-auto">
        <div className="flex flex-col items-center text-center gap-6">
          <div className="h-20 w-20 rounded-3xl bg-yellow-500/10 flex items-center justify-center border border-yellow-500/20">
            <AlertTriangle className="h-10 w-10 text-yellow-500" />
          </div>
          <div className="space-y-3">
            <h3 className="text-3xl font-bold text-white">Review your changes</h3>
            <p className="text-[#cbd5cf] text-lg leading-relaxed">
              You are about to update your public business identity. These changes will be reflected across all customer-facing platforms immediately.
            </p>
          </div>
        </div>

        <div className="rounded-[32px] bg-white/5 p-8 border border-white/10 space-y-4">
          <div className="flex justify-between items-center pb-4 border-b border-white/5">
            <span className="text-gray-400 text-base">Entity Type</span>
            <span className="text-white font-semibold px-3 py-1 bg-white/10 rounded-full text-sm">Business Organization</span>
          </div>
          <div className="flex justify-between items-center pt-2">
            <span className="text-gray-400 text-base">Impact Level</span>
            <span className="text-[#98FF98] font-semibold flex items-center gap-2 text-base">
              <CheckCircle2 className="h-5 w-5" /> Global Update
            </span>
          </div>
        </div>

        <div className="flex flex-col gap-4 pt-4">
          <Button
            onClick={form.handleSubmit(onSubmit)}
            disabled={loading}
            className="w-full bg-[#98FF98] text-[#0d241d] hover:bg-[#7cfc7c] font-bold h-14 rounded-2xl text-xl shadow-2xl shadow-[#98FF98]/20 transition-all hover:scale-[1.02] active:scale-[0.98]"
          >
            {loading ? <Loader2 className="h-6 w-6 animate-spin" /> : "Publish Changes"}
          </Button>
          <Button
            type="button"
            variant="ghost"
            onClick={() => setShowConfirmation(false)}
            disabled={loading}
            className="w-full text-[#cbd5cf] hover:text-white hover:bg-white/5 h-14 rounded-2xl text-lg"
          >
            Back to editing
          </Button>
        </div>
      </div>
    );
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-10 animate-in fade-in duration-500">
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-x-12 gap-y-8">
          {/* Left Column: Basic Info */}
          <div className="space-y-8">
            <div className="space-y-1">
              <h4 className="text-[#98FF98] font-semibold text-sm uppercase tracking-wider">Basic Information</h4>
              <p className="text-gray-500 text-sm">Update your brand's core identity.</p>
            </div>

            <FormField
              control={form.control}
              name="name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel className="text-[#F9F7F2] text-base font-medium">Brand Name</FormLabel>
                  <FormControl>
                    <Input className={inputClass} placeholder="e.g. ShareBite Downtown" {...field} />
                  </FormControl>
                  <FormDescription className="text-gray-500 text-xs italic">
                    Maximum 50 characters allowed for display.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="description"
              render={({ field }) => (
                <FormItem>
                  <FormLabel className="text-[#F9F7F2] text-base font-medium">Brand Bio</FormLabel>
                  <FormControl>
                    <Textarea 
                      className={cn(inputClass, "min-h-[180px] py-4 resize-none")} 
                      placeholder="Tell customers about your mission, the food you save, and your story..." 
                      {...field} 
                    />
                  </FormControl>
                  <FormDescription className="text-gray-500 text-xs">
                    Explain what makes your food special.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
          </div>

          {/* Right Column: Visuals */}
          <div className="space-y-8">
            <div className="space-y-1">
              <h4 className="text-[#98FF98] font-semibold text-sm uppercase tracking-wider">Visual Assets</h4>
              <p className="text-gray-500 text-sm">Configure your brand's look and feel.</p>
            </div>

            <FormField
              control={form.control}
              name="avatar"
              render={({ field }) => (
                <FormItem>
                  <FormLabel className="text-[#F9F7F2] text-base font-medium">Avatar URL</FormLabel>
                  <FormControl>
                    <Input className={inputClass} placeholder="https://image-server.com/avatar.png" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="banner"
              render={({ field }) => (
                <FormItem>
                  <FormLabel className="text-[#F9F7F2] text-base font-medium">Banner Background URL</FormLabel>
                  <FormControl>
                    <Input className={inputClass} placeholder="https://image-server.com/banner.jpg" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <div className="p-6 rounded-3xl bg-emerald-500/5 border border-emerald-500/10 space-y-3">
              <h5 className="text-emerald-400 font-semibold text-sm flex items-center gap-2">
                <CheckCircle2 className="h-4 w-4" /> Pro Tip
              </h5>
              <p className="text-gray-400 text-sm leading-relaxed">
                High-quality visual assets increase the probability of customer reservations by up to 40%. Use clear, well-lit images for your banner.
              </p>
            </div>
          </div>
        </div>

        {error && (
          <div className="rounded-[20px] border border-red-500/30 bg-red-500/10 p-5 text-sm text-red-400 font-medium flex items-center gap-4 animate-in slide-in-from-top-2">
            <div className="h-10 w-10 rounded-full bg-red-500/20 flex items-center justify-center shrink-0">
              <AlertTriangle className="h-5 w-5" />
            </div>
            {error}
          </div>
        )}

        <div className="flex justify-end gap-4 pt-10 border-t border-white/5">
          <Button
            type="button"
            variant="ghost"
            onClick={onCancel}
            disabled={loading}
            className="text-[#cbd5cf] hover:text-white hover:bg-white/5 h-14 px-10 rounded-2xl text-lg transition-all"
          >
            Discard
          </Button>
          <Button
            type="submit"
            disabled={loading}
            className="bg-[#98FF98] text-[#0d241d] hover:bg-[#7cfc7c] font-bold px-12 h-14 rounded-2xl text-lg shadow-xl shadow-[#98FF98]/10 transition-all hover:scale-[1.02] active:scale-[0.98]"
          >
            Review Changes
          </Button>
        </div>
      </form>
    </Form>
  );
}
