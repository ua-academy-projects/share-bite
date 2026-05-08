/* eslint-disable @typescript-eslint/no-explicit-any */
import * as React from "react";
import { useForm } from "react-hook-form";
import { useParams } from "react-router-dom";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { format } from "date-fns";
import { CalendarIcon } from "lucide-react";

import { businessApi } from "@/api/business";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Calendar } from "@/components/ui/calendar";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { cn } from "@/lib/utils";

const CATEGORIES = [
  { id: "1", name: "Випічка та десерти" },
  { id: "2", name: "Готові страви" },
  { id: "3", name: "Продукти" },
];

const createBoxSchema = z.object({
  categoryId: z.string().min(1, "Оберіть категорію"),
  priceFull: z.coerce.number().min(0.01, "Вкажіть ціну"),
  priceDiscount: z.coerce.number().min(0, "Вкажіть ціну"),
  quantity: z.coerce.number().min(1, "Мінімум 1"),
  expiresAt: z.date({ message: "Оберіть час" }),
  images: z.custom<FileList>().refine((files) => files && files.length > 0, "Додайте фото"),
}).refine((data) => data.priceDiscount <= data.priceFull, {
  message: "Знижка не може бути більшою за повну ціну",
  path: ["priceDiscount"],
});

type CreateBoxFormValues = z.infer<typeof createBoxSchema>;

export function CreateBoxPage() {
  const { id } = useParams<{ id: string }>();
  const [fileInputKey, setFileInputKey] = React.useState(0);
  const [loading, setLoading] = React.useState(false);
  const [successMessage, setSuccessMessage] = React.useState<string | null>(null);
  const [errorMessage, setErrorMessage] = React.useState<string | null>(null);

  const form = useForm<CreateBoxFormValues>({
    resolver: zodResolver(createBoxSchema) as any,
    defaultValues: {
      categoryId: "",
      priceFull: 0,
      priceDiscount: 0,
      quantity: 1,
    } as any,
  });

  const fileToBase64 = (file: File): Promise<string> => {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.readAsDataURL(file);
      reader.onload = () => resolve(reader.result as string);
      reader.onerror = (error) => reject(error);
    });
  };

  const onSubmit = async (values: CreateBoxFormValues) => {
    setSuccessMessage(null);
    setErrorMessage(null);
    if (!id) return setErrorMessage("ID закладу не знайдено.");

    const token = localStorage.getItem("token") || "";

    try {
      setLoading(true);
      const base64Image = await fileToBase64(values.images[0]);

      await businessApi.createBox({
        venue_id: parseInt(id),
        category_id: parseInt(values.categoryId),
        image: base64Image,
        price_full: values.priceFull,
        price_discount: values.priceDiscount,
        quantity: values.quantity,
        expires_at: values.expiresAt.toISOString(),
      }, token);

      setSuccessMessage("Magic Box успішно створено.");
      form.reset();
      setFileInputKey(prev => prev + 1);
      
    } catch (e) {
      setErrorMessage(e instanceof Error ? e.message : "Помилка");
    } finally {
      setLoading(false);
    }
  };

  // Базовий клас для всіх інпутів
  const inputClass = "bg-[#163d32] border-transparent text-white focus-visible:ring-green-500 rounded-xl placeholder:text-gray-400";
  // Клас для вимкнення стрілочок у number-інпутах
  const noSpinnersClass = "[appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none";

  return (
    <div className="min-h-screen w-full flex justify-center p-4">
      <div className="w-full max-w-2xl mt-12">
        <Card className="w-full bg-[#0d241d] border-[#2f5e50] text-white shadow-2xl rounded-3xl overflow-hidden">
          <CardHeader className="bg-[#163d32]/50 border-b border-[#2f5e50] pb-6">
            <CardTitle className="text-2xl font-bold tracking-tight text-white mb-4">
              Створити Magic Box 🪄
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
                  control={form.control as any}
                  name="categoryId"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel className="text-gray-200">Категорія</FormLabel>
                      <Select onValueChange={field.onChange} value={field.value}>
                        <FormControl>
                          <SelectTrigger className={cn(inputClass, "w-full")}>
                            <SelectValue placeholder="Оберіть категорію" />
                          </SelectTrigger>
                        </FormControl>
                        {/* Стилізований випадаючий список */}
                        <SelectContent className="bg-[#163d32] border-[#2f5e50] text-white rounded-xl">
                          {CATEGORIES.map((cat) => (
                            <SelectItem 
                              key={cat.id} 
                              value={cat.id}
                              className="focus:bg-[#2f5e50] focus:text-white cursor-pointer rounded-lg m-1"
                            >
                              {cat.name}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                      <FormMessage className="text-red-400" />
                    </FormItem>
                  )}
                />

                <div className="grid grid-cols-2 gap-4">
                  <FormField
                    control={form.control as any}
                    name="priceFull"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel className="text-gray-200">Ціна (₴)</FormLabel>
                        <FormControl>
                          <Input type="number" step="0.01" className={cn(inputClass, noSpinnersClass)} {...field} />
                        </FormControl>
                        <FormMessage className="text-red-400" />
                      </FormItem>
                    )}
                  />
                  <FormField
                    control={form.control as any}
                    name="priceDiscount"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel className="text-gray-200">Зі знижкою (₴)</FormLabel>
                        <FormControl>
                          <Input type="number" step="0.01" className={cn(inputClass, noSpinnersClass)} {...field} />
                        </FormControl>
                        <FormMessage className="text-red-400" />
                      </FormItem>
                    )}
                  />
                </div>

                <div className="grid grid-cols-2 gap-4">
                  <FormField
                    control={form.control as any}
                    name="quantity"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel className="text-gray-200">Кількість</FormLabel>
                        <FormControl>
                          <Input type="number" className={cn(inputClass, noSpinnersClass)} {...field} />
                        </FormControl>
                        <FormMessage className="text-red-400" />
                      </FormItem>
                    )}
                  />
                  <FormField
                    control={form.control as any}
                    name="expiresAt"
                    render={({ field }) => (
                      <FormItem className="flex flex-col">
                        <FormLabel className="mb-2 text-gray-200">Діє до</FormLabel>
                        <Popover>
                          <PopoverTrigger asChild>
                            <FormControl>
                              <Button
                                variant={"outline"}
                                className={cn(
                                  inputClass,
                                  "w-full pl-3 text-left font-normal border-0",
                                  !field.value && "text-gray-400"
                                )}
                              >
                                {field.value ? format(field.value, "PPP") : <span>Оберіть дату</span>}
                                <CalendarIcon className="ml-auto h-4 w-4 opacity-50" />
                              </Button>
                            </FormControl>
                          </PopoverTrigger>
                          {/* Стилізований календар */}
                          <PopoverContent className="w-auto p-0 bg-[#163d32] border-[#2f5e50] text-white rounded-xl" align="start">
                            <Calendar
                              mode="single"
                              selected={field.value}
                              onSelect={field.onChange}
                              disabled={(date) => date < new Date()}
                              className="text-white"
                            />
                          </PopoverContent>
                        </Popover>
                        <FormMessage className="text-red-400" />
                      </FormItem>
                    )}
                  />
                </div>

                <FormField
                  control={form.control as any}
                  name="images"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel className="text-gray-200">Фотографія</FormLabel>
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
                            key={fileInputKey}
                            type="file"
                            accept="image/*"
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
                  {loading ? "Відправка..." : "Створити Magic Box"}
                </Button>
              </form>
            </Form>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}