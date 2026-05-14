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

export type ListResponse<T> = {
  items: T[];
  total: number;
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

export type BrandProfile = {
  id: number;
  name: string;
  avatar?: string | null;
  banner?: string | null;
  description?: string | null;
};

export type BrandLocation = {
  id: number;
  name: string;
  avatar?: string | null;
  description?: string | null;
  latitude?: number | null;
  longitude?: number | null;
  tags: string[];
};

export type BrandPost = {
  id: number;
  content: string;
  created_at: string;
  images: string[];
  org: {
    id: number;
    name: string;
    profileType: string;
  };
};

export type BrandPostsResponse = ListResponse<BrandPost>;
export type BrandLocationsResponse = ListResponse<BrandLocation>;
export type NearbyBoxesResponse = ListResponse<Box>;

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

const API_BASE_URL = import.meta.env.VITE_API_URL || "http://localhost:3900";

export type UpdateBrandProfileRequest = {
  name?: string;
  avatar?: string;
  banner?: string;
  description?: string;
};

export const businessApi = {
  getNearbyBoxes: async (params: {
    lat: number;
    lon: number;
    skip?: number;
    limit?: number;
    orgId?: number;
    categoryId?: number;
  }): Promise<NearbyBoxesResponse> => {
    const search = new URLSearchParams();
    search.set("lat", String(params.lat));
    search.set("lon", String(params.lon));
    search.set("skip", String(params.skip ?? 0));
    search.set("limit", String(params.limit ?? 10));
    if (params.orgId) search.set("org_id", String(params.orgId));
    if (params.categoryId) search.set("category_id", String(params.categoryId));

    const response = await fetch(`${API_BASE_URL}/business/nearby-boxes?${search.toString()}`);
    if (!response.ok) {
      const err = await response.json().catch(() => ({}));
      throw new Error(err.error || err.message || `Failed to load boxes (${response.status})`);
    }

    const data = await response.json();
    return {
      items: data.items || [],
      total: data.total || 0,
    };
  },

  getLocationTags: async (): Promise<LocationTag[]> => {
    const response = await fetch(`${API_BASE_URL}/business/location-tags`);
    if (!response.ok) {
      throw new Error(`Failed to load location tags (${response.status})`);
    }
    return response.json();
  },

  getBrandProfile: async (id: number): Promise<BrandProfile> => {
    const response = await fetch(`${API_BASE_URL}/business/${id}`);
    if (!response.ok) {
      const err = await response.json().catch(() => ({}));
      throw new Error(err.error || err.message || `Failed to load brand (${response.status})`);
    }
    return response.json();
  },

  getBrandLocations: async (brandId: number, params?: {
    skip?: number;
    limit?: number;
    tags?: string[];
  }): Promise<BrandLocationsResponse> => {
    const search = new URLSearchParams();
    search.set("skip", String(params?.skip ?? 0));
    search.set("limit", String(params?.limit ?? 10));
    if (params?.tags && params.tags.length > 0) {
      search.set("tags", params.tags.map((tag) => tag.trim()).filter(Boolean).join(","));
    }

    const response = await fetch(`${API_BASE_URL}/business/org-units/${brandId}/locations?${search.toString()}`);
    if (!response.ok) {
      const err = await response.json().catch(() => ({}));
      throw new Error(err.error || err.message || `Failed to load locations (${response.status})`);
    }

    const data = await response.json();
    return {
      items: data.items || [],
      total: data.total || 0,
    };
  },

  getBrandPosts: async (params?: { skip?: number; limit?: number; orgIds?: number[] }): Promise<BrandPostsResponse> => {
    const search = new URLSearchParams();
    search.set("skip", String(params?.skip ?? 0));
    search.set("limit", String(params?.limit ?? 10));
    if (params?.orgIds && params.orgIds.length > 0) {
      search.set("org_id", params.orgIds.join(","));
    }

    const response = await fetch(`${API_BASE_URL}/business/posts?${search.toString()}`);
    if (!response.ok) {
      const err = await response.json().catch(() => ({}));
      throw new Error(err.error || err.message || `Failed to load posts (${response.status})`);
    }

    const data = await response.json();
    return {
      items: data.items || [],
      total: data.total || 0,
    };
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
    
    // Передаємо файл зображення
    if (data.image) {
      formData.append("image", data.image, data.image.name);
    }

    const authHeader = token.startsWith("Bearer ") ? token : `Bearer ${token}`;

    const response = await fetch(`${API_BASE_URL}/business/boxes`, {
      method: "POST",
      headers: {
        "Authorization": authHeader
      },
      body: formData,
    });

    if (!response.ok) {
      const err = await response.json().catch(() => ({}));
      throw new Error(err.error || err.message || `Помилка: ${response.status}`);
    }

    return response.json();
  },

  reserveBox: async (boxId: number, token: string): Promise<ReserveBoxResponse> => {
    const authHeader = token.startsWith("Bearer ") ? token : `Bearer ${token}`;

    const response = await fetch(`${API_BASE_URL}/business/boxes/${boxId}/reserve`, {
      method: "PATCH",
      headers: {
        Authorization: authHeader,
      },
    });

    if (!response.ok) {
      const err = await response.json().catch(() => ({}));
      throw new Error(err.error || err.message || `Reservation failed (${response.status})`);
    }

    return response.json();
  },

  updateBrandProfile: async (id: number, data: UpdateBrandProfileRequest, token: string): Promise<BrandProfile> => {
    const authHeader = token.startsWith("Bearer ") ? token : `Bearer ${token}`;
    
    const response = await fetch(`${API_BASE_URL}/business/${id}`, {
      method: "PATCH",
      headers: {
        "Content-Type": "application/json",
        Authorization: authHeader,
      },
      body: JSON.stringify(data),
    });

    if (!response.ok) {
      const err = await response.json().catch(() => ({}));
      // Handle business errors specifically if they exist
      const errorMessage = err.error || err.message || `Update failed (${response.status})`;
      throw new Error(errorMessage);
    }

    return response.json();
  },

  uploadAvatar: async (id: number, file: File, token: string): Promise<BrandProfile> => {
    const formData = new FormData();
    formData.append("image", file);

    const authHeader = token.startsWith("Bearer ") ? token : `Bearer ${token}`;
    const response = await fetch(`${API_BASE_URL}/business/${id}/avatar`, {
      method: "POST",
      headers: {
        Authorization: authHeader,
      },
      body: formData,
    });

    if (!response.ok) {
      const err = await response.json().catch(() => ({}));
      throw new Error(err.error || err.message || `Avatar upload failed (${response.status})`);
    }

    return response.json();
  },

  uploadBanner: async (id: number, file: File, token: string): Promise<BrandProfile> => {
    const formData = new FormData();
    formData.append("image", file);

    const authHeader = token.startsWith("Bearer ") ? token : `Bearer ${token}`;
    const response = await fetch(`${API_BASE_URL}/business/${id}/banner`, {
      method: "POST",
      headers: {
        Authorization: authHeader,
      },
      body: formData,
    });

    if (!response.ok) {
      const err = await response.json().catch(() => ({}));
      throw new Error(err.error || err.message || `Banner upload failed (${response.status})`);
    }

    return response.json();
  },
};
