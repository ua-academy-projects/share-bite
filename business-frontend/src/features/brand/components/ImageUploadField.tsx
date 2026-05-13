import React, { useRef, useState } from "react";
import { Loader2, Upload, X, Image as ImageIcon } from "lucide-react";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

interface ImageUploadFieldProps {
  label: string;
  value?: string | null;
  onUpload: (file: File) => Promise<void>;
  aspectRatio?: "square" | "video";
  className?: string;
}

export function ImageUploadField({
  label,
  value,
  onUpload,
  aspectRatio = "square",
  className,
}: ImageUploadFieldProps) {
  const [uploading, setUploading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    // Basic validation
    if (!file.type.startsWith("image/")) {
      setError("Please select an image file");
      return;
    }

    if (file.size > 5 * 1024 * 1024) {
      setError("File size must be less than 5MB");
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
    <div className={cn("space-y-2", className)}>
      <label className="text-sm font-medium text-[#cbd5cf]">{label}</label>
      
      <div 
        className={cn(
          "relative group overflow-hidden rounded-xl border-2 border-dashed border-white/10 bg-[#163d32] transition-all hover:border-emerald-500/50",
          aspectRatio === "square" ? "aspect-square w-32" : "aspect-[3/1] w-full",
          uploading && "opacity-50 pointer-events-none"
        )}
      >
        {value ? (
          <>
            <img 
              src={value} 
              alt={label} 
              className="h-full w-full object-cover"
            />
            <div className="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center">
              <Button
                type="button"
                variant="ghost"
                size="sm"
                onClick={triggerUpload}
                className="text-white hover:bg-white/20"
              >
                Change
              </Button>
            </div>
          </>
        ) : (
          <button
            type="button"
            onClick={triggerUpload}
            className="absolute inset-0 flex flex-col items-center justify-center gap-2 text-gray-500 hover:text-gray-400"
          >
            <Upload className="h-6 w-6" />
            <span className="text-xs">Upload</span>
          </button>
        )}

        {uploading && (
          <div className="absolute inset-0 flex items-center justify-center bg-black/20">
            <Loader2 className="h-6 w-6 animate-spin text-emerald-500" />
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
        <p className="text-xs text-red-400 mt-1">{error}</p>
      )}
    </div>
  );
}
