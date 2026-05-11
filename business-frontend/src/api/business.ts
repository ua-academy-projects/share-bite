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
  image: File;
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
    const formData = new FormData();
    
    const formattedDate = new Date(data.expires_at).toISOString();
    
    formData.append("venue_id", String(data.venue_id));
    formData.append("category_id", String(data.category_id || 1));
    formData.append("price_full", String(data.price_full));
    formData.append("price_discount", String(data.price_discount));
    formData.append("expires_at", formattedDate);
    formData.append("quantity", String(data.quantity));
    
    // Передаємо файл зображення
    if (data.image) {
      formData.append("image", data.image, data.image.name);
    }

    const response = await fetch(`${API_BASE_URL}/business/boxes`, {
      method: "POST",
      headers: {
        "Authorization": `Bearer ${token}`
      },
      body: formData,
    });

    if (!response.ok) {
      const err = await response.json().catch(() => ({}));
      throw new Error(err.error || err.message || `Помилка: ${response.status}`);
    }

    return response.json();
  }
};