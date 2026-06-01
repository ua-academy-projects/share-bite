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

export type BrandVenuesResponse = {
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

export type VenueRating = {
  rating: number;
};

export type RecommendedPost = {
  id: number;
  org_id: number;
  content: string;
  post_type: string;
  created_at: string;
};

export type RecommendedPostsResponse = {
  items: RecommendedPost[];
  total: number;
};

export type RecommendPostsRequest = {
  lat: number;
  lon: number;
  skip?: number;
  limit?: number;
};

const API_BASE_URL = import.meta.env.VITE_API_URL || "http://localhost:3999";
const BUSINESS_ACCOUNT_ID = Number(import.meta.env.VITE_BUSINESS_ACCOUNT_ID || "1");

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

  listBrandVenues: async (
    brandId: number,
    params: Pick<SearchVenuesRequest, "tags" | "skip" | "limit"> = {},
  ): Promise<BrandVenuesResponse> => {
    const search = new URLSearchParams();
    search.set("skip", String(params.skip ?? 0));
    search.set("limit", String(params.limit ?? 24));

    const tags = (params.tags || [])
      .map((tag) => tag.trim().toLowerCase())
      .filter(Boolean);
    if (tags.length > 0) search.set("tags", tags.join(","));

    const response = await fetch(
      `${API_BASE_URL}/business/org-units/${brandId}/locations?${search.toString()}`
    );
    if (!response.ok) {
      const err = await response.json().catch(() => ({}));
      throw new Error(err.error || err.message || `Failed to load brand venues (${response.status})`);
    }

    const data = await response.json();
    return {
      items: data.items || [],
      total: data.total || 0,
    };
  },

  listCurrentBusinessVenues: async (
    params: Pick<SearchVenuesRequest, "tags" | "skip" | "limit"> = {},
  ): Promise<BrandVenuesResponse> => {
    if (!Number.isFinite(BUSINESS_ACCOUNT_ID) || BUSINESS_ACCOUNT_ID <= 0) {
      throw new Error("Business account is not configured.");
    }

    return businessApi.listBrandVenues(BUSINESS_ACCOUNT_ID, params);
  },

  getVenueRating: async (id: number): Promise<VenueRating> => {
    const response = await fetch(`${API_BASE_URL}/business/org-units/${id}/rating`);
    if (!response.ok) {
      const err = await response.json().catch(() => ({}));
      throw new Error(err.error || err.message || `Failed to load venue rating (${response.status})`);
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

  recommendPosts: async (params: RecommendPostsRequest, token?: string): Promise<RecommendedPostsResponse> => {
    const search = new URLSearchParams();
    search.set("lat", String(params.lat));
    search.set("lon", String(params.lon));
    search.set("skip", String(params.skip ?? 0));
    search.set("limit", String(params.limit ?? 24));

    const headers: Record<string, string> = {
      "Content-Type": "application/json",
    };

    if (token) {
      headers["Authorization"] = `Bearer ${token}`;
    }

    const response = await fetch(`${API_BASE_URL}/business/posts/recommend?${search.toString()}`, {
      method: "GET",
      headers,
    });
    if (!response.ok) {
      const err = await response.json().catch(() => ({}));
      throw new Error(err.error || err.message || `Failed to load recommendations (${response.status})`);
    }

    const data = await response.json();
    return {
      items: data.items || [],
      total: data.total || 0,
    };
  }
};
