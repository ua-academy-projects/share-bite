export type Box = {
  id: number;
  image: string;
  full_price: string;
  discount_price: string;
};

const API_BASE_URL = import.meta.env.VITE_API_URL || "http://localhost:3999";

export const businessApi = {
  getNearbyBoxes: async (lat: number, lon: number): Promise<Box[]> => {
    try {
      const response = await fetch(`${API_BASE_URL}/business/nearby-boxes?lat=${lat}&lon=${lon}`);
      if (!response.ok) throw new Error("Network response was not ok");
      const data = await response.json();
      return data.items || [];
    } catch (error) {
      console.error("API error:", error);
      return [];
    }
  }
};