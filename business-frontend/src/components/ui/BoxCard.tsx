import { Box } from "@/api/business";
import { Button } from "@/components/ui/button";
import { MapPin, Clock, Flame } from "lucide-react";

interface BoxCardProps {
  box: Box;
}

const CATEGORY_MAP: Record<number, string> = {
  1: "Bakery & Desserts",
  2: "Ready Meals",
  3: "Groceries",
};

export function BoxCard({ box }: BoxCardProps) {
  // Видалено небезпечний split по amazonaws, тепер нормальні URL не ламаються
  const formatImageUrl = (url: string) => {
    if (!url) return "https://placehold.co/600x400/163d32/FFF?text=Magic+Box";
    return url;
  };

  const categoryName = box.category_id ? CATEGORY_MAP[box.category_id] || "Secret Box" : "Secret Box";

  // Захист для дати (якщо прийде Invalid Date, покажемо "N/A")
  const expireDate = new Date(box.expires_at);
  const formattedTime = Number.isNaN(expireDate.getTime())
    ? "N/A"
    : expireDate.toLocaleTimeString("en-US", { hour: "2-digit", minute: "2-digit" });

  // Захист для цін
  const fullPrice = Number(box.full_price);
  const discountPrice = Number(box.discount_price);

  return (
    <div className="bg-white dark:bg-[#163d32] border border-gray-100 dark:border-[#2f5e50] rounded-3xl overflow-hidden shadow-sm hover:shadow-xl dark:shadow-lg dark:hover:shadow-[#98FF98]/10 hover:-translate-y-1 transition-all duration-300 flex flex-col h-full group">
      {/* Upper part: Photo and badges */}
      <div className="relative h-52 overflow-hidden bg-gray-100 dark:bg-[#0d241d]">
        <img
          src={formatImageUrl(box.image)}
          alt={categoryName}
          className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
          onError={(e) => {
            e.currentTarget.src = "https://placehold.co/600x400/163d32/FFF?text=ShareBite";
          }}
        />
        {/* Gradient overlay */}
        <div className="absolute inset-0 bg-gradient-to-t from-white dark:from-[#163d32] to-transparent opacity-90 dark:opacity-80"></div>
        
        {/* Category badge */}
        <div className="absolute top-4 left-4 bg-black/60 backdrop-blur-md text-white text-xs font-semibold px-3 py-1.5 rounded-full border border-white/10 shadow-sm">
          {categoryName}
        </div>

        {/* "Low" availability badge */}
        {box.availability_status === "running_low" && (
          <div className="absolute top-4 right-4 bg-red-500 text-white text-xs font-bold px-3 py-1.5 rounded-full flex items-center gap-1 shadow-md animate-pulse">
            <Flame size={14} />
            Low!
          </div>
        )}
      </div>

      {/* Lower part: Info */}
      <div className="p-5 flex flex-col flex-1 relative z-10">
        <h3 className="text-[#1A3C34] dark:text-white text-xl font-bold tracking-tight mb-3">
          Magic Box
        </h3>
        
        <div className="flex items-center gap-5 text-gray-500 dark:text-gray-300 text-sm mb-5">
          <div className="flex items-center gap-1.5">
            <MapPin size={16} className="text-emerald-500 dark:text-[#98FF98]" />
            <span>{box.distance ? box.distance.toFixed(1) : "0.0"} km</span>
          </div>
          <div className="flex items-center gap-1.5">
            <Clock size={16} className="text-yellow-600 dark:text-[#FFD700]" />
            <span>Until {formattedTime}</span>
          </div>
        </div>

        {/* Pricing and Button */}
        <div className="mt-auto pt-4 border-t border-gray-100 dark:border-[#2f5e50]/50 flex justify-between items-end">
          <div className="flex flex-col">
            <span className="text-gray-400 dark:text-gray-400 text-sm line-through mb-0.5 font-medium">
              {Number.isFinite(fullPrice) ? fullPrice.toFixed(2) : "—"} ₴
            </span>
            <span className="text-emerald-600 dark:text-[#98FF98] text-2xl font-bold leading-none">
              {Number.isFinite(discountPrice) ? discountPrice.toFixed(2) : "—"} ₴
            </span>
          </div>
          <Button className="bg-[#FFD700] text-[#1A3C34] hover:bg-[#e6c200] dark:hover:bg-[#FFD700]/80 font-bold rounded-xl px-6 py-5 shadow-md dark:shadow-lg dark:shadow-[#FFD700]/20 transition-all">
            Reserve
          </Button>
        </div>
      </div>
    </div>
  );
}