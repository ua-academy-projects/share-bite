
import { Box } from "@/api/business";
import { Button } from "@/components/ui/button";

interface BoxCardProps {
  box: Box;
}

export function BoxCard({ box }: BoxCardProps) {
  const formatImageUrl = (url: string) => {
    if (url.includes("amazonaws.com/")) {
      return url.split("amazonaws.com/")[1];
    }
    return url;
  };

  return (
    <div className="bg-[#111] rounded-2xl overflow-hidden shadow-lg hover:scale-[1.02] transition">
      <img
        src={formatImageUrl(box.image)}
        alt="Surprise culinary box" // Додали alt (порада тіммейта)
        className="w-full h-40 object-cover"
      />
      <div className="p-4">
        <h3 className="text-white text-lg font-semibold">Surprise Box</h3>
        <div className="flex justify-between items-center mt-3">
          <span className="text-white font-semibold">${box.discount_price}</span>
          <Button variant="default" className="bg-green-500 text-black hover:bg-green-400">
            Reserve
          </Button>
        </div>
      </div>
    </div>
  );
}