import { useEffect, useState, useMemo } from "react";
import { businessApi, Box } from "@/api/business";
import { BoxCard } from "@/components/ui/BoxCard";
import { Loader2, Filter, Search, RotateCcw } from "lucide-react";

const CATEGORIES = [
  { id: "all", name: "All Categories" },
  { id: "4", name: "Bakery & Desserts" },
  { id: "5", name: "Sushi & Rolls" },
  { id: "6", name: "Groceries" },
];

export function BoxesPage() {
  const [boxes, setBoxes] = useState<Box[]>([]);
  const [loading, setLoading] = useState(true);

  // Стейт "чернетки" (що відображається в інпутах)
  const [draftCategory, setDraftCategory] = useState<string>("all");
  const [draftMaxDistance, setDraftMaxDistance] = useState<number>(10);
  const [draftMaxPrice, setDraftMaxPrice] = useState<number>(500);

  // Стейт застосованих фільтрів
  const [activeCategory, setActiveCategory] = useState<string>("all");
  const [activeMaxDistance, setActiveMaxDistance] = useState<number>(10);
  const [activeMaxPrice, setActiveMaxPrice] = useState<number>(500);

  // Викликається при першому завантаженні та при зміні АКТИВНОЇ категорії
  useEffect(() => {
    const loadBoxes = async () => {
      setLoading(true);
      try {
        // Категорія тепер фільтрується РЕАЛЬНО на бекенді
        const data = await businessApi.getNearbyBoxes(49.8397, 24.0297, activeCategory);
        setBoxes(data);
      } catch (error) {
        console.error("Failed to load boxes", error);
      } finally {
        setLoading(false);
      }
    };
    loadBoxes();
  }, [activeCategory]);

  const handleApplyFilters = () => {
    setActiveCategory(draftCategory);
    setActiveMaxDistance(draftMaxDistance);
    setActiveMaxPrice(draftMaxPrice);
  };

  const handleResetFilters = () => {
    setDraftCategory("all");
    setDraftMaxDistance(10);
    setDraftMaxPrice(500);
    
    setActiveCategory("all");
    setActiveMaxDistance(10);
    setActiveMaxPrice(500);
  };

  const handleReserveBox = (box: Box) => {
    alert(`Reservation for box #${box.id} is coming soon!`);
  };

  // Фільтрація по ціні та відстані залишається на фронтенді
  const filteredBoxes = useMemo(() => {
    return boxes.filter((box) => {
      const matchDistance = box.distance ? box.distance <= activeMaxDistance : true;
      const matchPrice = Number(box.discount_price) <= activeMaxPrice;
      return matchDistance && matchPrice;
    });
  }, [boxes, activeMaxDistance, activeMaxPrice]);

  return (
    <div className="min-h-screen bg-[#F9F7F2] dark:bg-[#0d241d] p-8 md:p-12 w-full transition-colors duration-300">
      <div className="max-w-7xl mx-auto">
        <div className="mb-8">
          <h1 className="text-4xl md:text-5xl font-bold text-[#1A3C34] dark:text-white tracking-tight mb-3">
            Curated Rescues <span className="text-emerald-500 dark:text-[#98FF98]">🌿</span>
          </h1>
          <p className="text-gray-600 dark:text-gray-400 text-lg">Rescue fresh food at a discount near you.</p>
        </div>

        {/* Панель фільтрів */}
        <div className="bg-white dark:bg-[#163d32] border border-gray-200 dark:border-[#2f5e50] rounded-2xl p-5 mb-10 shadow-sm transition-colors duration-300 flex flex-col xl:flex-row gap-6 items-start xl:items-center">
          <div className="flex items-center gap-2 text-[#1A3C34] dark:text-white font-semibold whitespace-nowrap">
            <Filter size={20} className="text-emerald-500 dark:text-[#98FF98]" />
            <span>Filters:</span>
          </div>
          
          <div className="flex flex-col md:flex-row gap-6 w-full flex-1">
            <div className="flex flex-col gap-2 flex-1 min-w-[200px]">
              <label className="text-xs text-gray-500 dark:text-gray-400 font-medium uppercase tracking-wider">Category</label>
              <select
                value={draftCategory}
                onChange={(e) => setDraftCategory(e.target.value)}
                className="bg-gray-50 dark:bg-[#0d241d] border border-gray-200 dark:border-transparent text-[#1A3C34] dark:text-white rounded-xl px-4 py-2.5 focus:ring-2 focus:ring-emerald-500 dark:focus:ring-[#98FF98] outline-none transition-all cursor-pointer w-full"
              >
                {CATEGORIES.map((cat) => (
                  <option key={cat.id} value={cat.id}>
                    {cat.name}
                  </option>
                ))}
              </select>
            </div>

            <div className="flex flex-col gap-2 flex-1 min-w-[200px]">
              <div className="flex justify-between items-center text-xs font-medium uppercase tracking-wider">
                <label className="text-gray-500 dark:text-gray-400">Max Price</label>
                <span className="text-emerald-600 dark:text-[#98FF98] font-bold">Up to {draftMaxPrice} ₴</span>
              </div>
              <input
                type="range"
                min="50"
                max="1000"
                step="10"
                value={draftMaxPrice}
                onChange={(e) => setDraftMaxPrice(parseInt(e.target.value))}
                className="w-full h-2 bg-gray-200 dark:bg-[#0d241d] rounded-lg appearance-none cursor-pointer accent-emerald-500 dark:accent-[#98FF98] mt-2"
              />
            </div>

            <div className="flex flex-col gap-2 flex-1 min-w-[200px]">
              <div className="flex justify-between items-center text-xs font-medium uppercase tracking-wider">
                <label className="text-gray-500 dark:text-gray-400">Distance</label>
                <span className="text-emerald-600 dark:text-[#98FF98] font-bold">Up to {draftMaxDistance} km</span>
              </div>
              <input
                type="range"
                min="0.5"
                max="20"
                step="0.5"
                value={draftMaxDistance}
                onChange={(e) => setDraftMaxDistance(parseFloat(e.target.value))}
                className="w-full h-2 bg-gray-200 dark:bg-[#0d241d] rounded-lg appearance-none cursor-pointer accent-emerald-500 dark:accent-[#98FF98] mt-2"
              />
            </div>
          </div>

          {/* Кнопки дій */}
          <div className="flex items-end gap-3 w-full xl:w-auto h-full mt-2 xl:mt-0 pt-2 xl:pt-6">
            <button
              onClick={handleResetFilters}
              disabled={loading}
              className="flex-1 xl:flex-none flex items-center justify-center gap-2 bg-gray-100 text-gray-600 hover:bg-gray-200 dark:bg-[#0d241d] dark:text-gray-300 dark:hover:bg-[#2f5e50] font-bold rounded-xl px-6 py-2.5 transition-all shadow-sm disabled:opacity-70"
            >
              <RotateCcw size={18} />
              Reset
            </button>
            <button
              onClick={handleApplyFilters}
              disabled={loading}
              className="flex-1 xl:flex-none flex items-center justify-center gap-2 bg-[#163d32] text-white hover:bg-[#1A3C34] dark:bg-emerald-500 dark:text-black dark:hover:bg-emerald-400 font-bold rounded-xl px-6 py-2.5 transition-all shadow-md disabled:opacity-70"
            >
              <Search size={18} />
              Apply
            </button>
          </div>
        </div>

        {loading ? (
          <div className="flex justify-center items-center h-64 w-full">
            <Loader2 className="w-12 h-12 text-emerald-500 dark:text-[#98FF98] animate-spin" />
          </div>
        ) : filteredBoxes.length > 0 ? (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-8">
            {filteredBoxes.map((box) => (
              <BoxCard key={box.id} box={box} onReserve={handleReserveBox} />
            ))}
          </div>
        ) : (
          <div className="text-center bg-white dark:bg-[#163d32] border border-gray-200 dark:border-[#2f5e50] rounded-3xl p-16 shadow-sm dark:shadow-none transition-colors duration-300">
            <p className="text-[#1A3C34] dark:text-gray-300 text-xl font-bold">No boxes match your filters 😢</p>
            <p className="text-gray-500 mt-2">Try expanding the distance, price, or selecting a different category.</p>
            <button 
              onClick={handleResetFilters}
              className="mt-6 text-emerald-600 dark:text-[#98FF98] font-semibold hover:underline"
            >
              Reset filters
            </button>
          </div>
        )}
      </div>
    </div>
  );
}