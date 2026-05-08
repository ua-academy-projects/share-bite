export type Box = {
  id: number;
  venue_id: number;
  category_id?: number;
  image: string;
  full_price: string | number;
  discount_price: string | number;
  created_at: string;
  expires_at: string;
  availability_status: string;
  distance: number;
};

export type CreateBoxRequest = {
  venue_id: number;
  category_id?: number;
  image: string;
  price_full: number;
  price_discount: number;
  expires_at: string;
  quantity: number;
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
  },

  createBox: async (data: CreateBoxRequest, token: string) => {
    const response = await fetch(`${API_BASE_URL}/business/boxes`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${token}`
      },
      body: JSON.stringify(data),
    });

    if (!response.ok) {
      const err = await response.json().catch(() => ({}));
      throw new Error(err.message || `Помилка створення боксу (${response.status})`);
    }

    return response.json();
  }
};