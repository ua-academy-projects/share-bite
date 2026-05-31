import { useEffect, useState, useMemo } from "react";
import { businessApi, Box } from "@/api/business";
import { BoxCard } from "@/components/ui/BoxCard";
import { PageHeader } from "@/components/layout/PageHeader";
import { PageLayout } from "@/components/layout/PageLayout";
import {
  pageBtnPrimary,
  pageBtnSecondary,
  pageEmpty,
  pageFilterBar,
  pageInput,
  pageLabel,
  pageLoader,
} from "@/components/layout/pageStyles";
import { Loader2, Filter, Search, RotateCcw } from "lucide-react";
import { cn } from "@/lib/utils";

const CATEGORIES = [
  { id: "all", name: "All Categories" },
  { id: "4", name: "Bakery & Desserts" },
  { id: "5", name: "Sushi & Rolls" },
  { id: "6", name: "Groceries" },
];

export function BoxesPage() {
  const [boxes, setBoxes] = useState<Box[]>([]);
  const [loading, setLoading] = useState(true);

  const [draftCategory, setDraftCategory] = useState<string>("all");
  const [draftMaxDistance, setDraftMaxDistance] = useState<number>(10);
  const [draftMaxPrice, setDraftMaxPrice] = useState<number>(500);

  const [activeCategory, setActiveCategory] = useState<string>("all");
  const [activeMaxDistance, setActiveMaxDistance] = useState<number>(10);
  const [activeMaxPrice, setActiveMaxPrice] = useState<number>(500);

  useEffect(() => {
    const loadBoxes = async () => {
      setLoading(true);
      try {
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

  const filteredBoxes = useMemo(() => {
    return boxes.filter((box) => {
      const matchDistance = box.distance ? box.distance <= activeMaxDistance : true;
      const matchPrice = Number(box.discount_price) <= activeMaxPrice;
      return matchDistance && matchPrice;
    });
  }, [boxes, activeMaxDistance, activeMaxPrice]);

  return (
    <PageLayout>
      <PageHeader
        title={
          <>
            Curated Rescues <span className="text-emerald-500 dark:text-[#98FF98]">🌿</span>
          </>
        }
        description="Rescue fresh food at a discount near you."
      />

      <div
        className={cn(
          pageFilterBar,
          "mb-10 flex flex-col items-start gap-6 xl:flex-row xl:items-center"
        )}
      >
        <div className="flex items-center gap-2 whitespace-nowrap font-semibold text-[#1A3C34] dark:text-white">
          <Filter size={20} className="text-emerald-500 dark:text-[#98FF98]" />
          <span>Filters:</span>
        </div>

        <div className="flex w-full flex-1 flex-col gap-6 md:flex-row">
          <div className="flex min-w-[200px] flex-1 flex-col gap-2">
            <label className={pageLabel}>Category</label>
            <select
              value={draftCategory}
              onChange={(e) => setDraftCategory(e.target.value)}
              className={cn(pageInput, "cursor-pointer py-2.5")}
            >
              {CATEGORIES.map((cat) => (
                <option key={cat.id} value={cat.id}>
                  {cat.name}
                </option>
              ))}
            </select>
          </div>

          <div className="flex min-w-[200px] flex-1 flex-col gap-2">
            <div className="flex items-center justify-between text-xs font-medium uppercase tracking-wider">
              <label className={pageLabel}>Max Price</label>
              <span className="font-bold text-emerald-600 dark:text-[#98FF98]">
                Up to {draftMaxPrice} ₴
              </span>
            </div>
            <input
              type="range"
              min="50"
              max="1000"
              step="10"
              value={draftMaxPrice}
              onChange={(e) => setDraftMaxPrice(parseInt(e.target.value))}
              className="mt-2 h-2 w-full cursor-pointer appearance-none rounded-lg bg-gray-200 accent-emerald-500 dark:bg-[#0d241d] dark:accent-[#98FF98]"
            />
          </div>

          <div className="flex min-w-[200px] flex-1 flex-col gap-2">
            <div className="flex items-center justify-between text-xs font-medium uppercase tracking-wider">
              <label className={pageLabel}>Distance</label>
              <span className="font-bold text-emerald-600 dark:text-[#98FF98]">
                Up to {draftMaxDistance} km
              </span>
            </div>
            <input
              type="range"
              min="0.5"
              max="20"
              step="0.5"
              value={draftMaxDistance}
              onChange={(e) => setDraftMaxDistance(parseFloat(e.target.value))}
              className="mt-2 h-2 w-full cursor-pointer appearance-none rounded-lg bg-gray-200 accent-emerald-500 dark:bg-[#0d241d] dark:accent-[#98FF98]"
            />
          </div>
        </div>

        <div className="mt-2 flex h-full w-full items-end gap-3 pt-2 xl:mt-0 xl:w-auto xl:pt-6">
          <button
            onClick={handleResetFilters}
            disabled={loading}
            className={cn(pageBtnSecondary, "flex flex-1 items-center justify-center gap-2 xl:flex-none")}
          >
            <RotateCcw size={18} />
            Reset
          </button>
          <button
            onClick={handleApplyFilters}
            disabled={loading}
            className={cn(pageBtnPrimary, "flex flex-1 items-center justify-center gap-2 xl:flex-none")}
          >
            <Search size={18} />
            Apply
          </button>
        </div>
      </div>

      {loading ? (
        <div className="flex h-64 w-full items-center justify-center">
          <Loader2 className={cn(pageLoader, "h-12 w-12")} />
        </div>
      ) : filteredBoxes.length > 0 ? (
        <div className="grid grid-cols-1 gap-8 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
          {filteredBoxes.map((box) => (
            <BoxCard key={box.id} box={box} onReserve={handleReserveBox} />
          ))}
        </div>
      ) : (
        <div className={pageEmpty}>
          <p className="text-xl font-bold text-[#1A3C34] dark:text-gray-300">
            No boxes match your filters 😢
          </p>
          <p className="mt-2 text-gray-500">
            Try expanding the distance, price, or selecting a different category.
          </p>
          <button
            onClick={handleResetFilters}
            className="mt-6 font-semibold text-emerald-600 hover:underline dark:text-[#98FF98]"
          >
            Reset filters
          </button>
        </div>
      )}
    </PageLayout>
  );
}
