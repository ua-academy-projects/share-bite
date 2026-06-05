import * as React from "react";
import { useForm, ControllerRenderProps } from "react-hook-form";
import { Link } from "react-router-dom";
import { useParams } from "react-router-dom";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { CheckCircle, ArrowLeft } from "lucide-react";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";

// Імпортуємо наш новий уніфікований компонент та його тип
import { PostCard, PostData } from "@/components/ui/PostCard";

const createPostSchema = z.object({
  textData: z
    .string()
    .trim()
    .min(5, "Description must contain at least 5 characters"),
  images: z
    .custom<FileList>()
    .refine((files) => files && files.length > 0, "Please add at least one image"),
});

type CreatePostFormValues = z.infer<typeof createPostSchema>;

export default function CreatePostPage() {
  const { id } = useParams<{ id: string }>();
  const fileInputRef = React.useRef<HTMLInputElement | null>(null);

  const [loading, setLoading] = React.useState(false);
  const [errorMessage, setErrorMessage] = React.useState<string | null>(null);
  
  const [createdPost, setCreatedPost] = React.useState<PostData | null>(null);

  const form = useForm<CreatePostFormValues>({
    resolver: zodResolver(createPostSchema),
    defaultValues: { textData: "" },
    mode: "onTouched",
  });

  async function onSubmit(values: CreatePostFormValues) {
    setErrorMessage(null);
    setCreatedPost(null);

    if (!id) {
      setErrorMessage("Error: Venue ID not found in URL.");
      return;
    }

    const token = localStorage.getItem("token");
    if (!token) {
      setErrorMessage("Token missing. Please log in again.");
      return;
    }

    try {
      setLoading(true);

      const apiBase = import.meta.env.VITE_API_URL;
      if (!apiBase) {
        setErrorMessage("VITE_API_URL is not configured.");
        setLoading(false);
        return;
      }
      
      const url = new URL(`/business/posts/${encodeURIComponent(id)}`, apiBase).toString();

      const fd = new FormData();
      fd.append("content", values.textData);
      Array.from(values.images).forEach((file) => fd.append("photos", file));

      const res = await fetch(url, {
        method: "POST",
        headers: { Authorization: `Bearer ${token}` },
        body: fd,
      });

      if (!res.ok) {
        let details = "";
        try {
          const ct = res.headers.get("content-type") ?? "";
          if (ct.includes("application/json")) {
            const data = await res.json();
            details = data?.message ? `: ${data.message}` : "";
          } else {
            const text = await res.text();
            details = text ? `: ${text}` : "";
          }
        } catch {
          // ignore
        }
        throw new Error(`Request failed (${res.status})${details}`);
      }

      const postData: PostData = await res.json();

      // ✨ Optimistic UI: Створюємо локальні прев'ю
      if (values.images && values.images.length > 0) {
        postData.images = Array.from(values.images).map(file => URL.createObjectURL(file));
      }

      setCreatedPost(postData);

      form.reset({ textData: "" });
      if (fileInputRef.current) {
        fileInputRef.current.value = "";
      }
    } catch (e) {
      setErrorMessage(e instanceof Error ? e.message : "An unknown error occurred");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="min-h-screen w-full flex justify-center p-4 md:p-8 bg-[#F9F7F2] dark:bg-[#0d241d] transition-colors duration-300">
      <div className="w-full max-w-xl mt-6 md:mt-12 animate-in fade-in slide-in-from-bottom-4 duration-500">
        {id && (
          <Button asChild variant="ghost" className="mb-4 text-[#1A3C34] dark:text-gray-200 hover:bg-white/70 dark:hover:bg-[#163d32] rounded-xl">
            <Link to={`/venue/${id}`}>
              <ArrowLeft size={18} />
              Back to venue
            </Link>
          </Button>
        )}
        
        {createdPost ? (
          <div className="space-y-8 flex flex-col items-center">
            
            {/* Success Message залишається на сторінці створення */}
            <div className="bg-emerald-100 dark:bg-emerald-500/10 border border-emerald-500/30 rounded-3xl p-6 flex flex-col items-center justify-center text-center w-full shadow-sm">
              <div className="w-16 h-16 bg-emerald-500 rounded-full flex items-center justify-center text-white mb-4 shadow-lg shadow-emerald-500/30">
                <CheckCircle size={32} strokeWidth={2.5} />
              </div>
              <h2 className="text-2xl font-extrabold text-[#1A3C34] dark:text-emerald-400 mb-1">Awesome! Post is live.</h2>
              <p className="text-emerald-700 dark:text-emerald-200/70 font-medium">Your followers can now see this delicious offer.</p>
            </div>

            {/* Використовуємо наш новий уніфікований компонент! */}
            <PostCard post={createdPost} />

            {/* Кнопка "Створити ще" */}
            <Button 
              onClick={() => {
                if (createdPost?.images) {
                  createdPost.images.forEach(url => {
                    if (url.startsWith('blob:')) URL.revokeObjectURL(url);
                  });
                }
                setCreatedPost(null);
              }} 
              className="w-full max-w-[460px] bg-white dark:bg-[#112f26] text-[#1A3C34] dark:text-white border border-gray-200 dark:border-[#2f5e50]/60 hover:bg-gray-50 dark:hover:bg-[#163d32] font-bold rounded-full py-6 text-lg transition-all flex items-center justify-center gap-2 shadow-sm"
            >
              <ArrowLeft size={20} />
              Create Another Post
            </Button>
          </div>
        ) : (
          <Card className="w-full bg-white dark:bg-[#0d241d] border border-gray-200 dark:border-[#2f5e50] shadow-xl rounded-3xl overflow-hidden transition-colors duration-300">
            <CardHeader className="bg-gray-50 dark:bg-[#163d32]/50 border-b border-gray-200 dark:border-[#2f5e50] pb-6 transition-colors duration-300">
              <CardTitle className="text-2xl font-bold tracking-tight text-[#1A3C34] dark:text-white mb-4">
                Create Post
              </CardTitle>
              <div className="flex items-center gap-3">
                <div className="w-10 h-10 bg-black rounded-full flex items-center justify-center text-white font-bold text-sm shadow-md">
                  SB
                </div>
                <div>
                  <h2 className="text-[#1A3C34] dark:text-white font-semibold">Share Bite</h2>
                  <p className="text-gray-500 dark:text-gray-300 text-xs">Venue ID: {id ?? "—"}</p>
                </div>
              </div>
            </CardHeader>

            <CardContent className="pt-6">
              <Form {...form}>
                <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
                  <FormField
                    control={form.control}
                    name="textData"
                    render={({ field }: { field: ControllerRenderProps<CreatePostFormValues, "textData"> }) => (
                      <FormItem>
                        <FormLabel className="text-gray-700 dark:text-gray-200 font-medium">Post Description</FormLabel>
                        <FormControl>
                          <Textarea
                            placeholder="Tell us about your delicious offer... (min 5 characters)"
                            className="bg-gray-50 dark:bg-[#163d32] border border-gray-200 dark:border-transparent text-[#1A3C34] dark:text-white focus-visible:ring-emerald-500 dark:focus-visible:ring-green-500 rounded-xl min-h-[120px] resize-none placeholder:text-gray-400 dark:placeholder:text-gray-400 transition-colors duration-300"
                            {...field}
                            disabled={loading}
                          />
                        </FormControl>
                        <FormMessage className="text-red-500 dark:text-red-400" />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={form.control}
                    name="images"
                    render={({ field }: { field: ControllerRenderProps<CreatePostFormValues, "images"> }) => {
                      const selectedFiles = form.watch("images");

                      return (
                        <FormItem>
                          <FormLabel className="text-gray-700 dark:text-gray-200 font-medium">Photos</FormLabel>
                          <FormControl>
                            <div className="space-y-4">
                              <div className="relative border-2 border-dashed border-gray-300 dark:border-[#2f5e50] rounded-xl p-6 text-center hover:bg-gray-50 dark:hover:bg-[#163d32]/50 transition cursor-pointer flex flex-col items-center justify-center gap-2 group">
                                <div className="w-12 h-12 rounded-full bg-gray-100 dark:bg-[#2f5e50] group-hover:bg-gray-200 dark:group-hover:bg-[#3a7564] transition flex items-center justify-center text-emerald-600 dark:text-green-400 mb-2 shadow-sm">
                                  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={2} stroke="currentColor" className="w-6 h-6">
                                    <path strokeLinecap="round" strokeLinejoin="round" d="M6.827 6.175A2.31 2.31 0 0 1 5.186 7.23c-.38.054-.757.112-1.134.175C2.999 7.58 2.25 8.507 2.25 9.574V18a2.25 2.25 0 0 0 2.25 2.25h15A2.25 2.25 0 0 0 21.75 18V9.574c0-1.067-.75-1.994-1.802-2.169a47.865 47.865 0 0 0-1.134-.175 2.31 2.31 0 0 1-1.64-1.055l-.822-1.316a2.192 2.192 0 0 0-1.736-1.039 48.774 48.774 0 0 0-5.232 0 2.192 2.192 0 0 0-1.736 1.039l-.821 1.316Z" />
                                    <path strokeLinecap="round" strokeLinejoin="round" d="M16.5 12.75a4.5 4.5 0 1 1-9 0 4.5 4.5 0 0 1 9 0Z" />
                                  </svg>
                                </div>
                                <span className="text-sm text-gray-600 dark:text-gray-300 font-medium">
                                  {selectedFiles && selectedFiles.length > 0 
                                    ? `${selectedFiles.length} file(s) selected. Click to change.` 
                                    : "Click here to select files"}
                                </span>
                                <Input
                                  ref={fileInputRef}
                                  type="file"
                                  accept="image/*"
                                  multiple
                                  disabled={loading}
                                  className="absolute inset-0 w-full h-full opacity-0 cursor-pointer"
                                  onChange={(e) => field.onChange(e.target.files)}
                                />
                              </div>
                              
                              {selectedFiles && selectedFiles.length > 0 && (
                                <div className="grid grid-cols-3 sm:grid-cols-4 gap-4 mt-4">
                                  {Array.from(selectedFiles).map((file, idx) => (
                                    <div key={idx} className="relative aspect-square rounded-xl overflow-hidden border border-gray-200 dark:border-[#2f5e50] shadow-sm bg-gray-50 dark:bg-[#0d241d]">
                                      <img
                                        src={URL.createObjectURL(file)}
                                        alt={`preview-${idx}`}
                                        className="w-full h-full object-cover transition-transform hover:scale-105 duration-300"
                                      />
                                    </div>
                                  ))}
                                </div>
                              )}
                            </div>
                          </FormControl>
                          <FormMessage className="text-red-500 dark:text-red-400" />
                        </FormItem>
                      );
                    }}
                  />

                  {errorMessage && (
                    <div className="rounded-xl border border-red-500/30 bg-red-50 dark:bg-red-500/10 px-4 py-3 text-sm text-red-700 dark:text-red-400 font-medium animate-in fade-in">
                      {errorMessage}
                    </div>
                  )}

                  <Button 
                    type="submit" 
                    disabled={loading} 
                    className="w-full bg-[#163d32] text-white hover:bg-[#1A3C34] dark:bg-emerald-500 dark:text-black dark:hover:bg-emerald-400 font-bold rounded-full py-6 text-lg mt-2 transition-all shadow-lg dark:shadow-emerald-500/20"
                  >
                    {loading ? "Sending..." : "Create Post"}
                  </Button>
                </form>
              </Form>
            </CardContent>
          </Card>
        )}
      </div>
    </div>
  );
}
