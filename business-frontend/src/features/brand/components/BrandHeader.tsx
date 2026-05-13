import type { BrandProfile } from "@/api/business";
import { Button } from "@/components/ui/button";
import { Settings2 } from "lucide-react";

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
      <div className="rounded-3xl border border-red-500/30 bg-red-500/10 px-6 py-6 text-red-200">
        <p className="text-sm font-semibold">Brand failed to load</p>
        <p className="mt-1 text-sm text-red-100/80">{error}</p>
      </div>
    );
  }

  return (
    <div className="rounded-[28px] border border-white/10 bg-[#163d32]/50 backdrop-blur-xl overflow-hidden">
      <div className="relative h-44 md:h-60">
        {brand?.banner ? (
          <img
            src={brand.banner}
            alt={brand.name}
            className="h-full w-full object-cover"
          />
        ) : (
          <div className="h-full w-full bg-[radial-gradient(circle_at_top,_#1f4a3f_0%,_#0b0f0e_55%,_#0b0f0e_100%)]" />
        )}
        <div className="absolute inset-0 bg-gradient-to-t from-[#0b0f0e]/90 via-[#0b0f0e]/40 to-transparent" />
      </div>

      <div className="relative px-6 pb-6 pt-4 md:px-10 md:pb-8">
        <div className="flex flex-col gap-4 md:flex-row md:items-start md:gap-6">
          <div className="-mt-12 md:-mt-14">
            {brand?.avatar ? (
              <img
                src={brand.avatar}
                alt={brand.name}
                className="h-20 w-20 md:h-24 md:w-24 rounded-2xl object-cover border border-white/20 shadow-lg"
              />
            ) : (
              <div className="h-20 w-20 md:h-24 md:w-24 rounded-2xl border border-white/15 bg-[#0f1b17] flex items-center justify-center text-lg font-semibold text-[#98FF98]">
                {brand?.name ? getInitials(brand.name) : "SB"}
              </div>
            )}
          </div>

          <div className="flex-1">
            <div className="flex flex-wrap items-center justify-between gap-3">
              <div className="flex items-center gap-3">
                <h1 className="text-2xl md:text-3xl font-semibold text-[#F9F7F2] tracking-tight">
                  {brand?.name ?? (loading ? "" : "")}
                </h1>
                {loading ? (
                  <span className="h-5 w-20 rounded-full bg-white/10 animate-pulse" />
                ) : null}
              </div>

              {!loading && brand && onEdit && (
                <Button
                  onClick={onEdit}
                  variant="outline"
                  className="rounded-full border-white/15 bg-white/5 text-[#F9F7F2] hover:bg-white/10 gap-2 h-9 px-4 text-sm"
                >
                  <Settings2 className="h-4 w-4" />
                  Edit Profile
                </Button>
              )}
            </div>
            <p className="mt-2 text-sm md:text-base text-[#cbd5cf]">
              {loading ? (
                <span className="block h-4 w-3/4 bg-white/10 rounded-full animate-pulse" />
              ) : brand?.description ? (
                brand.description
              ) : (
                "Description not provided."
              )}
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}
