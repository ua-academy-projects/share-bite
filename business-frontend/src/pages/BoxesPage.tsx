import { useEffect, useState } from "react";
import { businessApi, Box } from "@/api/business";
import { BoxCard } from "@/components/ui/BoxCard";
import { Loader2 } from "lucide-react";

export function BoxesPage() {
  const [boxes, setBoxes] = useState<Box[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const loadBoxes = async () => {
      try {
        const data = await businessApi.getNearbyBoxes(49.8397, 24.0297);
        setBoxes(data);
      } catch (error) {
        console.error("Failed to load boxes", error);
      } finally {
        setLoading(false);
      }
    };
    loadBoxes();
  }, []);

  return (
    <div className="min-h-screen bg-[#F9F7F2] dark:bg-[#0d241d] p-8 md:p-12 w-full transition-colors duration-300">
      <div className="max-w-7xl mx-auto">
        <div className="mb-10">
          <h1 className="text-4xl md:text-5xl font-bold text-[#1A3C34] dark:text-white tracking-tight mb-3">
            Curated Rescues <span className="text-emerald-500 dark:text-[#98FF98]">🌿</span>
          </h1>
          <p className="text-gray-600 dark:text-gray-400 text-lg">Rescue fresh food at a discount near you.</p>
        </div>

        {loading ? (
          <div className="flex justify-center items-center h-64 w-full">
            <Loader2 className="w-12 h-12 text-emerald-500 dark:text-[#98FF98] animate-spin" />
          </div>
        ) : boxes.length > 0 ? (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-8">
            {boxes.map((box) => (
              <BoxCard key={box.id} box={box} />
            ))}
          </div>
        ) : (
          <div className="text-center bg-white dark:bg-[#163d32] border border-gray-200 dark:border-[#2f5e50] rounded-3xl p-16 shadow-sm dark:shadow-none transition-colors duration-300">
            <p className="text-[#1A3C34] dark:text-gray-300 text-xl font-bold">No available boxes nearby yet 😢</p>
            <p className="text-gray-500 mt-2">Try checking back here a little later.</p>
          </div>
        )}
      </div>
    </div>
  );
}