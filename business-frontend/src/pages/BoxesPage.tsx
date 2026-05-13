import { useEffect, useState } from "react";
import { businessApi, Box, ReserveBoxResponse } from "@/api/business";
import { BoxCard } from "@/components/ui/BoxCard";
import { Loader2 } from "lucide-react";

type ReservationInfo = ReserveBoxResponse & { boxId: number };

export function BoxesPage() {
  const [boxes, setBoxes] = useState<Box[]>([]);
  const [loading, setLoading] = useState(true);

  const [reservingBoxId, setReservingBoxId] = useState<number | null>(null);
  const [reservation, setReservation] = useState<ReservationInfo | null>(null);
  const [reserveError, setReserveError] = useState<string | null>(null);

  useEffect(() => {
    const loadBoxes = async () => {
      try {
        const data = await businessApi.getNearbyBoxes({
          lat: 49.8397,
          lon: 24.0297,
          limit: 24,
        });
        setBoxes(data.items);
      } catch (error) {
        console.error("Failed to load boxes", error);
      } finally {
        setLoading(false);
      }
    };
    void loadBoxes();
  }, []);

  const handleReserve = async (box: Box) => {
    setReserveError(null);

    const token = localStorage.getItem("token");
    if (!token) {
      setReserveError("Token missing. Please log in.");
      return;
    }

    try {
      setReservingBoxId(box.id);
      const data = await businessApi.reserveBox(box.id, token);
      setReservation({ ...data, boxId: box.id });
    } catch (e) {
      setReserveError(e instanceof Error ? e.message : "Failed to reserve the box");
    } finally {
      setReservingBoxId(null);
    }
  };

  const discountPrice = Number(reservation?.price_discount);

  return (
    <div className="min-h-screen bg-[#F9F7F2] dark:bg-[#0d241d] p-8 md:p-12 w-full transition-colors duration-300">
      <div className="max-w-7xl mx-auto">
        <div className="mb-10">
          <h1 className="text-4xl md:text-5xl font-bold text-[#1A3C34] dark:text-white tracking-tight mb-3">
            Curated Rescues <span className="text-emerald-500 dark:text-[#98FF98]">🌿</span>
          </h1>
          <p className="text-gray-600 dark:text-gray-400 text-lg">Rescue fresh food at a discount near you.</p>
        </div>

        {reservation && (
          <div className="mb-8 rounded-2xl border border-green-500/30 bg-green-50 dark:bg-green-500/10 px-5 py-4 text-green-800 dark:text-green-300">
            <p className="font-bold text-base">Reservation confirmed</p>
            <p className="text-sm mt-1">
              Box code: <span className="font-semibold">{reservation.box_code}</span>
            </p>
            <p className="text-sm">
              Price: {Number.isFinite(discountPrice) ? discountPrice.toFixed(2) : reservation.price_discount} ₴
            </p>
          </div>
        )}

        {reserveError && (
          <div className="mb-8 rounded-2xl border border-red-500/30 bg-red-50 dark:bg-red-500/10 px-5 py-4 text-red-700 dark:text-red-400">
            {reserveError}
          </div>
        )}

        {loading ? (
          <div className="flex justify-center items-center h-64 w-full">
            <Loader2 className="w-12 h-12 text-emerald-500 dark:text-[#98FF98] animate-spin" />
          </div>
        ) : boxes.length > 0 ? (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-8">
            {boxes.map((box) => (
              <BoxCard
                key={box.id}
                box={box}
                onReserve={handleReserve}
                reserving={reservingBoxId === box.id}
                isReserved={reservation?.boxId === box.id}
              />
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
