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

export type LocationTag = {
  id: number;
  name: string;
  slug: string;
};

export type VenueSearchItem = {
  id: number;
  name: string;
  avatar?: string | null;
  description?: string | null;
  latitude?: number | null;
  longitude?: number | null;
  tags: string[];
};

export type SearchVenuesRequest = {
  query?: string;
  tags?: string[];
  skip?: number;
  limit?: number;
};

export type SearchVenuesResponse = {
  items: VenueSearchItem[];
  total: number;
};

export type ReserveBoxResponse = {
  image: string;
  price_full: string | number;
  price_discount: string | number;
  box_code: string;
};

export type VenueBrand = {
  id: number;
  name: string;
  avatar?: string | null;
};

export type VenueProfile = {
  id: number;
  name: string;
  avatar?: string | null;
  banner?: string | null;
  description?: string | null;
  latitude?: number | null;
  longitude?: number | null;
  tags: string[];
  brand?: VenueBrand | null;
};

const API_BASE_URL = import.meta.env.VITE_API_URL || "http://localhost:3999";

export const businessApi = {
  // Тепер ми приймаємо categoryId і кидаємо його в запит
  getNearbyBoxes: async (lat: number, lon: number, categoryId?: string): Promise<Box[]> => {
    try {
      let url = `${API_BASE_URL}/business/nearby-boxes?lat=${lat}&lon=${lon}&limit=50`;
      if (categoryId && categoryId !== "all") {
        url += `&category_id=${categoryId}`;
      }
      
      const response = await fetch(url);
      if (!response.ok) throw new Error("Network response was not ok");
      const data = await response.json();
      return data.items || [];
    } catch (error) {
      console.error("API error:", error);
      return [];
    }
  },

  getLocationTags: async (): Promise<LocationTag[]> => {
    const response = await fetch(`${API_BASE_URL}/business/location-tags`);
    if (!response.ok) {
      throw new Error(`Failed to load location tags (${response.status})`);
    }
    return response.json();
  },

  searchVenues: async (params: SearchVenuesRequest): Promise<SearchVenuesResponse> => {
    const query = (params.query || "").trim();
    const tags = (params.tags || [])
      .map((tag) => tag.trim().toLowerCase())
      .filter(Boolean);

    if (!query && tags.length === 0) {
      throw new Error("At least one filter is required: query or tags.");
    }

    const search = new URLSearchParams();
    if (query) search.set("q", query);
    if (tags.length > 0) search.set("tags", tags.join(","));
    search.set("skip", String(params.skip ?? 0));
    search.set("limit", String(params.limit ?? 10));

    const response = await fetch(`${API_BASE_URL}/business/venues/search?${search.toString()}`);
    if (!response.ok) {
      const err = await response.json().catch(() => ({}));
      throw new Error(err.error || err.message || `Search failed (${response.status})`);
    }

    const data = await response.json();
    return {
      items: data.items || [],
      total: data.total || 0,
    };
  },

  getVenueProfile: async (id: number): Promise<VenueProfile> => {
    const response = await fetch(`${API_BASE_URL}/business/org-units/${id}`);
    if (!response.ok) {
      const err = await response.json().catch(() => ({}));
      throw new Error(err.error || err.message || `Failed to load venue (${response.status})`);
    }
    return response.json();
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
  },

  reserveBox: async (boxID: number, token: string): Promise<ReserveBoxResponse> => {
    const response = await fetch(`${API_BASE_URL}/business/boxes/${boxID}/reserve`, {
      method: "PATCH",
      headers: {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${token}`
      },
    });

    if (!response.ok) {
      const err = await response.json().catch(() => ({}));
      throw new Error(err.error || `Помилка резервування боксу (${response.status})`);
    }

    return response.json();
  },
};