import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { Loader2, AlertTriangle } from "lucide-react";

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
      setError(err instanceof Error ? err.message : "Failed to update profile");
    } finally {
      setLoading(false);
    }
  };

  const inputClass = "bg-[#163d32] border-white/5 text-white focus-visible:ring-emerald-500 rounded-xl placeholder:text-gray-500 h-11 transition-all focus:bg-[#1a4a3d]";

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
            <FormField
              control={form.control}
              name="avatar"
              render={({ field }) => (
                <FormItem>
                  <FormLabel className="text-[#cbd5cf]">Avatar URL</FormLabel>
                  <FormControl>
                    <Input className={inputClass} placeholder="https://..." {...field} />
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
                  <FormLabel className="text-[#cbd5cf]">Banner URL</FormLabel>
                  <FormControl>
                    <Input className={inputClass} placeholder="https://..." {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
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
            {loading ? (
              <>
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                Saving...
              </>
            ) : (
              "Save Changes"
            )}
          </Button>
        </div>
      </form>
    </Form>
  );
}
