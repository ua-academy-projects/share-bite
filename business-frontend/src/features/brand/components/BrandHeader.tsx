import type { BrandProfile } from "@/api/business";
import { Button } from "@/components/ui/button";
import { Settings2, ShieldCheck } from "lucide-react";

function getInitials(value: string) {
  const parts = value.trim().split(/\s+/);
  if (parts.length === 0) return "SB";
  const first = parts[0]?.[0] ?? "S";
  const second = parts.length > 1 ? parts[1]?.[0] ?? "" : "";
  return `${first}${second}`.toUpperCase();
}

type BrandHeaderProps = {
  brand: BrandProfile | null;
  loading: boolean;
  error: string | null;
  onEdit?: () => void;
};

export function BrandHeader({ brand, loading, error, onEdit }: BrandHeaderProps) {
  if (error) {
    return (
      <div className="rounded-[40px] border border-red-500/30 bg-red-500/10 px-8 py-10 text-red-200 backdrop-blur-md">
        <p className="text-lg font-bold">Brand failed to load</p>
        <p className="mt-2 text-sm text-red-100/70">{error}</p>
      </div>
    );
  }

  return (
    <div className="relative rounded-[48px] border border-white/10 bg-[#163d32]/20 backdrop-blur-2xl overflow-hidden shadow-2xl">
      {/* Immersive Banner */}
      <div className="relative h-64 md:h-[400px] w-full overflow-hidden">
        {brand?.banner ? (
          <img
            src={brand.banner}
            alt={brand.name}
            className="h-full w-full object-cover transition-transform duration-1000 hover:scale-105"
          />
        ) : (
          <div className="h-full w-full bg-[radial-gradient(circle_at_top,_#1f4a3f_0%,_#0b0f0e_60%,_#050706_100%)]" />
        )}
        {/* Multi-layered Gradient for Depth */}
        <div className="absolute inset-0 bg-gradient-to-t from-[#0b0f0e] via-[#0b0f0e]/30 to-transparent" />
        <div className="absolute inset-0 bg-black/10" />
      </div>

      <div className="relative px-8 pb-10 pt-0 md:px-16">
        <div className="flex flex-col gap-6 md:flex-row md:items-end md:gap-10">
          {/* Prominent Avatar */}
          <div className="-mt-20 md:-mt-32 relative group">
            <div className="relative h-32 w-32 md:h-48 md:w-48 rounded-[40px] p-1.5 bg-gradient-to-br from-white/20 to-white/5 backdrop-blur-xl shadow-2xl transition-transform duration-500 group-hover:scale-[1.02]">
              {brand?.avatar ? (
                <img
                  src={brand.avatar}
                  alt={brand.name}
                  className="h-full w-full rounded-[34px] object-cover border border-white/10 shadow-inner"
                />
              ) : (
                <div className="h-full w-full rounded-[34px] border border-white/10 bg-[#0f1b17] flex items-center justify-center text-3xl font-black text-[#98FF98] shadow-inner">
                  {brand?.name ? getInitials(brand.name) : "SB"}
                </div>
              )}
            </div>
          </div>

          <div className="flex-1 space-y-4 pb-2">
            <div className="flex flex-wrap items-center justify-between gap-6">
              <div className="space-y-1">
                <div className="flex items-center gap-3">
                  <h1 className="text-4xl md:text-6xl font-black text-white tracking-tighter">
                    {brand?.name ?? (loading ? "Loading..." : "")}
                  </h1>
                </div>
              </div>

              {!loading && brand && onEdit && (
                <Button
                  onClick={onEdit}
                  className="rounded-2xl bg-white/5 border border-white/10 px-6 py-6 text-sm font-bold text-white backdrop-blur-md transition-all hover:bg-white/10 hover:border-white/20 active:scale-95"
                >
                  <Settings2 className="mr-2 h-4 w-4" />
                  Edit Profile
                </Button>
              )}
            </div>
            
            {brand?.description && (
              <p className="max-w-3xl text-lg md:text-xl font-medium leading-relaxed text-[#cbd5cf]">
                {brand.description}
              </p>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
