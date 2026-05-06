import { useEffect, useState } from "react";
import { businessApi, Box } from "@/api/business";
import { BoxCard } from "@/components/ui/BoxCard";

export function BoxesPage() {
  const [boxes, setBoxes] = useState<Box[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const loadBoxes = async () => {
      const data = await businessApi.getNearbyBoxes(49.8397, 24.0297);
      setBoxes(data);
      setLoading(false);
    };
    loadBoxes();
  }, []);

  if (loading) return <div className="p-8 text-white">Loading...</div>;

  return (
    <div className="p-8 text-white">
      <h1 className="text-4xl font-semibold mb-2">Curated Rescues</h1>
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mt-6">
        {boxes.map((box) => (
          // Виправлено: передаємо об'єкт box без лапок
          <BoxCard key={box.id} box={box} />
        ))}
      </div>
    </div>
  );
}