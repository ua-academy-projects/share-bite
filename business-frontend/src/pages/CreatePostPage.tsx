import * as React from "react";
import { useForm, ControllerRenderProps } from "react-hook-form";
import { useParams } from "react-router-dom";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";

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

const createPostSchema = z.object({
  textData: z
    .string()
    .trim()
    .min(5, "Опис має містити мінімум 5 символів"),
  images: z
    .custom<FileList>()
    .refine((files) => files && files.length > 0, "Додайте хоча б одне зображення"),
});

type CreatePostFormValues = z.infer<typeof createPostSchema>;

export default function CreatePostPage() {
  const { id } = useParams<{ id: string }>();

  const [loading, setLoading] = React.useState(false);
  const [successMessage, setSuccessMessage] = React.useState<string | null>(null);
  const [errorMessage, setErrorMessage] = React.useState<string | null>(null);

  const form = useForm<CreatePostFormValues>({
    resolver: zodResolver(createPostSchema),
    defaultValues: { textData: "" },
    mode: "onTouched",
  });

  async function onSubmit(values: CreatePostFormValues) {
    setSuccessMessage(null);
    setErrorMessage(null);

    if (!id) {
      setErrorMessage("Помилка: ID бізнесу/організації не знайдено в URL.");
      return;
    }

    const token = localStorage.getItem("token");
    if (!token) {
      setErrorMessage("Немає токена. Увійдіть у систему ще раз.");
      return;
    }

    try {
      setLoading(true);

      const apiBase = import.meta.env.VITE_API_URL;
      const url = `${apiBase}/business/posts/${encodeURIComponent(id)}`;

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
        throw new Error(`Помилка запиту (${res.status})${details}`);
      }

      setSuccessMessage("Допис успішно створено.");

      form.reset({ textData: "" });
      const fileInput = document.querySelector('input[type="file"]') as HTMLInputElement | null;
      if (fileInput) fileInput.value = "";
    } catch (e) {
      setErrorMessage(e instanceof Error ? e.message : "Сталася невідома помилка");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="min-h-screen w-full flex justify-center p-4">
      <div className="w-full max-w-2xl mt-12">
        <Card className="w-full bg-[#0d241d] border-[#2f5e50] text-white shadow-2xl rounded-3xl overflow-hidden">
          <CardHeader className="bg-[#163d32]/50 border-b border-[#2f5e50] pb-6">
            <CardTitle className="text-2xl font-bold tracking-tight text-white mb-4">
              Створити допис
            </CardTitle>
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 bg-black rounded-full flex items-center justify-center text-white font-bold text-sm">
                SB
              </div>
              <div>
                <h2 className="text-white font-semibold">Share Bite</h2>
                <p className="text-gray-300 text-xs">ID закладу: {id ?? "—"}</p>
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
                      <FormLabel className="text-gray-200">Опис допису</FormLabel>
                      <FormControl>
                        <Textarea
                          placeholder="Розкажіть про вашу смачну пропозицію... (мінімум 5 символів)"
                          className="bg-[#163d32] border-transparent text-white focus-visible:ring-green-500 rounded-xl min-h-[120px] resize-none placeholder:text-gray-400"
                          {...field}
                          disabled={loading}
                        />
                      </FormControl>
                      <FormMessage className="text-red-400" />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="images"
                  render={({ field }: { field: ControllerRenderProps<CreatePostFormValues, "images"> }) => (
                    <FormItem>
                      <FormLabel className="text-gray-200">Фотографії</FormLabel>
                      <FormControl>
                        <div className="relative border-2 border-dashed border-[#2f5e50] rounded-xl p-6 text-center hover:bg-[#163d32]/50 transition cursor-pointer flex flex-col items-center justify-center gap-2 group">
                          <div className="w-12 h-12 rounded-full bg-[#2f5e50] group-hover:bg-[#3a7564] transition flex items-center justify-center text-green-400 mb-2">
                            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={2} stroke="currentColor" className="w-6 h-6">
                              <path strokeLinecap="round" strokeLinejoin="round" d="M6.827 6.175A2.31 2.31 0 0 1 5.186 7.23c-.38.054-.757.112-1.134.175C2.999 7.58 2.25 8.507 2.25 9.574V18a2.25 2.25 0 0 0 2.25 2.25h15A2.25 2.25 0 0 0 21.75 18V9.574c0-1.067-.75-1.994-1.802-2.169a47.865 47.865 0 0 0-1.134-.175 2.31 2.31 0 0 1-1.64-1.055l-.822-1.316a2.192 2.192 0 0 0-1.736-1.039 48.774 48.774 0 0 0-5.232 0 2.192 2.192 0 0 0-1.736 1.039l-.821 1.316Z" />
                              <path strokeLinecap="round" strokeLinejoin="round" d="M16.5 12.75a4.5 4.5 0 1 1-9 0 4.5 4.5 0 0 1 9 0Z" />
                            </svg>
                          </div>
                          <span className="text-sm text-gray-300 font-medium">Натисніть сюди, щоб обрати файли</span>
                          <Input
                            type="file"
                            accept="image/*"
                            multiple
                            disabled={loading}
                            className="absolute inset-0 w-full h-full opacity-0 cursor-pointer"
                            onChange={(e) => field.onChange(e.target.files)}
                          />
                        </div>
                      </FormControl>
                      <FormMessage className="text-red-400" />
                    </FormItem>
                  )}
                />

                {successMessage && (
                  <div className="rounded-xl border border-green-500/30 bg-green-500/10 px-4 py-3 text-sm text-green-400 font-medium">
                    {successMessage}
                  </div>
                )}

                {errorMessage && (
                  <div className="rounded-xl border border-red-500/30 bg-red-500/10 px-4 py-3 text-sm text-red-400 font-medium">
                    {errorMessage}
                  </div>
                )}

                <Button 
                  type="submit" 
                  disabled={loading} 
                  className="w-full bg-green-500 text-black hover:bg-green-400 font-bold rounded-full py-6 text-lg mt-2 transition-all shadow-lg shadow-green-500/20"
                >
                  {loading ? "Відправка..." : "Створити допис"}
                </Button>
              </form>
            </Form>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}