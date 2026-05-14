import React, { useRef, useState } from "react";
import { Loader2, Upload } from "lucide-react";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

interface ImageUploadFieldProps {
  label?: string;
  value?: string | null;
  onUpload: (file: File) => Promise<void>;
  aspectRatio?: "square" | "video";
  className?: string;
  variant?: "default" | "profile";
}

export function ImageUploadField({
  label,
  value,
  onUpload,
  aspectRatio = "square",
  className,
  variant = "default",
}: ImageUploadFieldProps) {
  const [uploading, setUploading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    if (!file.type.startsWith("image/")) {
      setError("Image required");
      return;
    }

    if (file.size > 5 * 1024 * 1024) {
      setError("Max 5MB");
      return;
    }

    setUploading(true);
    setError(null);

    try {
      await onUpload(file);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Upload failed");
    } finally {
      setUploading(false);
      if (fileInputRef.current) {
        fileInputRef.current.value = "";
      }
    }
  };

  const triggerUpload = () => {
    fileInputRef.current?.click();
  };

  return (
    <div className={cn("group relative", className)}>
      {label && <label className="text-xs font-black uppercase tracking-widest text-[#9fb2a7] mb-2 block">{label}</label>}
      
      <div 
        className={cn(
          "relative overflow-hidden transition-all duration-300",
          aspectRatio === "square" 
            ? "aspect-square rounded-[32px] border-2 border-white/10" 
            : "aspect-[3/1] rounded-[40px] border border-white/10",
          variant === "profile" && aspectRatio === "square" && "shadow-2xl border-4 border-[#0d241d]",
          !value && "bg-white/5 border-dashed",
          uploading && "opacity-50 pointer-events-none"
        )}
      >
        {value ? (
          <>
            <img 
              src={value} 
              alt={label || "Upload"} 
              className="h-full w-full object-cover transition-transform duration-500 group-hover:scale-105"
            />
            <div className="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center backdrop-blur-sm">
              <Button
                type="button"
                variant="ghost"
                size="sm"
                onClick={triggerUpload}
                className="text-white hover:bg-white/20 font-bold border border-white/20 rounded-xl"
              >
                Change
              </Button>
            </div>
          </>
        ) : (
          <button
            type="button"
            onClick={triggerUpload}
            className="absolute inset-0 flex flex-col items-center justify-center gap-3 text-gray-500 hover:text-emerald-400 hover:bg-white/5 transition-all"
          >
            <Upload className="h-8 w-8" />
            <span className="text-xs font-bold uppercase tracking-widest">Upload {label}</span>
          </button>
        )}

        {uploading && (
          <div className="absolute inset-0 flex items-center justify-center bg-black/20 backdrop-blur-sm">
            <Loader2 className="h-8 w-8 animate-spin text-[#98FF98]" />
          </div>
        )}
      </div>

      <input
        type="file"
        ref={fileInputRef}
        onChange={handleFileChange}
        accept="image/*"
        className="hidden"
      />

      {error && (
        <p className="text-[10px] font-black uppercase text-red-400 mt-2 tracking-wider bg-red-400/10 px-2 py-1 rounded inline-block">{error}</p>
      )}
    </div>
  );
}
