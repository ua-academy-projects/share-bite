import type { BrandLocation } from "@/api/business";
import { MapPin, ChevronRight, Settings2 } from "lucide-react";
import { Link } from "react-router-dom";

function getInitials(value: string) {
  const parts = value.trim().split(/\s+/);
  if (parts.length === 0) return "SB";
  const first = parts[0]?.[0] ?? "S";
  const second = parts.length > 1 ? parts[1]?.[0] ?? "" : "";
  return `${first}${second}`.toUpperCase();
}

type BrandLocationCardProps = {
  location: BrandLocation;
};

export function BrandLocationCard({ location }: BrandLocationCardProps) {
  const hasCoords = Number.isFinite(location.latitude) && Number.isFinite(location.longitude);

  return (
    <Link 
      to={`/venue/${location.id}`}
      className="group block transition-all duration-300 hover:-translate-y-1.5 active:scale-[0.98]"
    >
      <article className="relative overflow-hidden rounded-[32px] border border-white/10 bg-[#163d32]/40 backdrop-blur-md p-6 flex flex-col gap-5 h-full transition-all group-hover:bg-[#163d32]/60 group-hover:border-[#98FF98]/30 group-hover:shadow-[0_20px_40px_-15px_rgba(0,0,0,0.3)]">
        {/* Glow effect on hover */}
        <div className="absolute -right-16 -top-16 h-32 w-32 rounded-full bg-[#98FF98]/5 blur-3xl transition-opacity opacity-0 group-hover:opacity-100" />
        
        <div className="flex items-start gap-5">
          <div className="relative">
            {location.avatar ? (
              <img
                src={location.avatar}
                alt={location.name}
                className="h-16 w-16 rounded-[22px] object-cover border border-white/15 shadow-inner"
              />
            ) : (
              <div className="h-16 w-16 rounded-[22px] border border-white/10 bg-[#0b0f0e] flex items-center justify-center text-sm font-bold text-[#98FF98] shadow-inner">
                {getInitials(location.name)}
              </div>
            )}
            {/* Online status or similar indicator could go here */}
            <div className="absolute -bottom-1 -right-1 h-5 w-5 rounded-full bg-[#0b0f0e] p-1 shadow-lg">
               <div className="h-full w-full rounded-full bg-[#98FF98] animate-pulse" />
            </div>
          </div>

          <div className="flex-1 min-w-0">
            <div className="flex items-start justify-between gap-2">
              <h3 className="text-xl font-bold text-[#F9F7F2] truncate tracking-tight group-hover:text-white transition-colors">
                {location.name}
              </h3>
              <div className="flex h-8 w-8 items-center justify-center rounded-full bg-white/5 border border-white/10 text-[#9fb2a7] transition-all group-hover:bg-[#98FF98] group-hover:text-[#0b0f0e] group-hover:border-transparent group-hover:rotate-12">
                <Settings2 className="h-4 w-4" />
              </div>
            </div>
            {location.description ? (
              <p className="mt-1.5 text-sm text-[#cbd5cf] leading-relaxed line-clamp-2">
                {location.description}
              </p>
            ) : (
              <p className="mt-1.5 text-sm text-[#9fb2a7]/60 italic">
                No description provided.
              </p>
            )}
          </div>
        </div>

        {location.tags?.length > 0 ? (
          <div className="flex flex-wrap gap-2">
            {location.tags.map((tag) => (
              <span
                key={tag}
                className="rounded-full border border-white/10 bg-white/5 px-3 py-1 text-[11px] font-medium text-[#cbd5cf] transition-colors group-hover:border-white/20 group-hover:text-white"
              >
                {tag}
              </span>
            ))}
          </div>
        ) : null}

        <div className="mt-auto pt-4 border-t border-white/5 flex items-center justify-between">
          {hasCoords ? (
            <div className="flex items-center gap-2 text-xs font-medium text-[#9fb2a7] group-hover:text-[#cbd5cf] transition-colors">
              <div className="flex h-6 w-6 items-center justify-center rounded-lg bg-[#98FF98]/10">
                <MapPin className="h-3.5 w-3.5 text-[#98FF98]" />
              </div>
              <span>
                {location.latitude?.toFixed(4)}, {location.longitude?.toFixed(4)}
              </span>
            </div>
          ) : <div />}

          <div className="flex items-center gap-1 text-[11px] font-bold text-[#98FF98] uppercase tracking-widest opacity-0 group-hover:opacity-100 transition-all transform translate-x-2 group-hover:translate-x-0">
            Manage <ChevronRight className="h-3 w-3" />
          </div>
        </div>
      </article>
    </Link>
  );
}
