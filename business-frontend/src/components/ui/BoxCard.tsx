import { Box } from "@/api/business";
import { Button } from "@/components/ui/button";
import { MapPin, Clock, Flame, ChevronRight, Sparkles } from "lucide-react";
import { Link } from "react-router-dom";

interface BoxCardProps {
  box: Box;
  onReserve?: (box: Box) => void;
  reserving?: boolean;
  isReserved?: boolean;
}

const CATEGORY_MAP: Record<number, string> = {
  1: "Bakery & Desserts",
  2: "Ready Meals",
  3: "Groceries",
};

export function BoxCard({ box,  onReserve, reserving = false, isReserved = false  }: BoxCardProps) {
  const formatImageUrl = (url: string) => {
    if (!url) return "https://placehold.co/600x400/163d32/FFF?text=Magic+Box";
    return url;
  };

  const categoryName = box.category_id ? CATEGORY_MAP[box.category_id] || "Secret Box" : "Secret Box";

  const expireDate = new Date(box.expires_at);
  const formattedTime = Number.isNaN(expireDate.getTime())
    ? "N/A"
    : expireDate.toLocaleTimeString("en-US", { hour: "2-digit", minute: "2-digit" });

  const fullPrice = Number(box.full_price);
  const discountPrice = Number(box.discount_price);

  const isSoldOut = box.availability_status === "sold_out";
  const showAction = Boolean(onReserve);
  const disabled = !showAction || reserving || isReserved || isSoldOut;

  const buttonLabel = isReserved
    ? "Reserved"
    : reserving
      ? "Reserving..."
      : isSoldOut
        ? "Sold Out"
        : "Reserve";
  
  const content = (
    <div className="group relative flex h-full flex-col overflow-hidden rounded-[32px] border border-white/10 bg-white dark:bg-[#163d32]/40 backdrop-blur-md transition-all duration-300 hover:-translate-y-1.5 hover:border-[#98FF98]/30 hover:bg-[#163d32]/60 hover:shadow-[0_20px_40px_-15px_rgba(0,0,0,0.3)]">
      {/* Photo Section */}
      <div className="relative aspect-[16/10] overflow-hidden">
        <img
          src={formatImageUrl(box.image)}
          alt={categoryName}
          className="h-full w-full object-cover transition-transform duration-700 group-hover:scale-110"
          onError={(e) => {
            e.currentTarget.src = "https://placehold.co/600x400/163d32/FFF?text=ShareBite";
          }}
        />
        <div className="absolute inset-0 bg-gradient-to-t from-white dark:from-[#0b0f0e] via-transparent to-transparent opacity-60 dark:opacity-80" />
        
        {/* Badges */}
        <div className="absolute left-4 top-4 flex flex-wrap gap-2">
          <div className="rounded-full bg-black/40 backdrop-blur-md border border-white/10 px-3 py-1.5 text-[10px] font-bold uppercase tracking-widest text-white">
            {categoryName}
          </div>
          {box.availability_status === "running_low" && (
            <div className="flex items-center gap-1.5 rounded-full bg-red-500 px-3 py-1.5 text-[10px] font-bold uppercase tracking-widest text-white shadow-lg animate-pulse">
              <Flame size={12} /> Low
            </div>
          )}
        </div>

        {!showAction && (
          <div className="absolute bottom-4 right-4 flex h-10 w-10 items-center justify-center rounded-full bg-white/10 backdrop-blur-md border border-white/20 text-white opacity-0 transition-all transform scale-75 group-hover:opacity-100 group-hover:scale-100">
            <ChevronRight className="h-5 w-5" />
          </div>
        )}
      </div>

      {/* Info Section */}
      <div className="flex flex-1 flex-col p-6">
        <div className="mb-4 flex items-center gap-2">
          <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-[#98FF98]/10 text-[#98FF98]">
            <Sparkles size={16} />
          </div>
          <h3 className="text-xl font-bold tracking-tight text-[#1A3C34] dark:text-[#F9F7F2]">
            Magic Box
          </h3>
        </div>
        
        <div className="flex items-center gap-6 text-sm font-medium text-gray-500 dark:text-[#9fb2a7]">
          <div className="flex items-center gap-1.5">
            <MapPin size={16} className="text-[#98FF98]" />
            <span>{box.distance ? box.distance.toFixed(1) : "0.0"} km</span>
          </div>
          <div className="flex items-center gap-1.5">
            <Clock size={16} className="text-[#FFD700]" />
            <span>{formattedTime}</span>
          </div>
        </div>

        {/* Pricing & CTA */}
        <div className="mt-auto pt-6 flex items-end justify-between border-t border-gray-100 dark:border-white/5">
          <div className="flex flex-col">
            <span className="text-xs font-bold text-gray-400 dark:text-[#9fb2a7]/60 line-through tracking-wider">
              {Number.isFinite(fullPrice) ? fullPrice.toFixed(0) : "—"} ₴
            </span>
            <div className="flex items-baseline gap-1">
              <span className="text-3xl font-black text-[#1A3C34] dark:text-[#98FF98] tracking-tighter">
                {Number.isFinite(discountPrice) ? discountPrice.toFixed(0) : "—"}
              </span>
              <span className="text-sm font-bold text-[#1A3C34] dark:text-[#98FF98]">₴</span>
            </div>
          </div>

          {showAction ? (
            <Button
              onClick={(e) => {
                e.preventDefault();
                e.stopPropagation();
                onReserve?.(box);
              }}
              disabled={disabled}
              className="rounded-2xl bg-[#FFD700] px-6 py-6 text-sm font-black uppercase tracking-widest text-[#1A3C34] shadow-[0_8px_20px_-6px_rgba(255,215,0,0.4)] transition-all hover:scale-105 hover:bg-[#FFD700] hover:shadow-[0_12px_24px_-8px_rgba(255,215,0,0.6)] active:scale-95 disabled:opacity-50"
            >
              {buttonLabel}
            </Button>
          ) : (
            <div className="text-[10px] font-bold uppercase tracking-[0.2em] text-[#98FF98] opacity-0 group-hover:opacity-100 transition-opacity">
               Manage Box
            </div>
          )}
        </div>
      </div>
    </div>
  );

  if (!showAction) {
    return (
      <Link to={`/venue/${box.venue_id}`} className="block h-full transition-all">
        {content}
      </Link>
    );
  }

  return content;
}