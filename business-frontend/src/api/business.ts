import axios from "axios";

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

const API_BASE = import.meta.env.VITE_API_BASE_URL || "";

const businessAxios = axios.create({
  baseURL: `${API_BASE}/api/business`,
});

businessAxios.interceptors.request.use((config) => {
  const token = localStorage.getItem("token");
  if (token && config.headers) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export const businessApi = {
  getNearbyBoxes: async (
    lat: number,
    lon: number,
    categoryId?: string
  ): Promise<Box[]> => {
    try {
      const params: Record<string, string | number> = {
        lat,
        lon,
        limit: 50,
      };
      if (categoryId && categoryId !== "all") {
        params.category_id = categoryId;
      }
      const res = await businessAxios.get<{ items?: Box[] }>("nearby-boxes", {
        params,
      });
      return res.data.items || [];
    } catch (error) {
      console.error("API error:", error);
      return [];
    }
  },

  getLocationTags: async (): Promise<LocationTag[]> => {
    const res = await businessAxios.get<LocationTag[]>("location-tags");
    return res.data;
  },

  searchVenues: async (
    params: SearchVenuesRequest
  ): Promise<SearchVenuesResponse> => {
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

    const res = await businessAxios.get<SearchVenuesResponse>(
      `venues/search?${search.toString()}`
    );
    return {
      items: res.data.items || [],
      total: res.data.total || 0,
    };
  },

  getVenueProfile: async (id: number): Promise<VenueProfile> => {
    const res = await businessAxios.get<VenueProfile>(`org-units/${id}`);
    return res.data;
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

    const res = await businessAxios.post("boxes", formData, {
      headers: { Authorization: `Bearer ${token}` },
    });
    return res.data;
  },

  reserveBox: async (
    boxID: number,
    token: string
  ): Promise<ReserveBoxResponse> => {
    const res = await businessAxios.patch<ReserveBoxResponse>(
      `boxes/${boxID}/reserve`,
      {},
      {
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
      }
    );
    return res.data;
  },

  recommendPosts: async (
    params: RecommendPostsRequest,
    token?: string
  ): Promise<RecommendedPostsResponse> => {
    const search = new URLSearchParams();
    search.set("lat", String(params.lat));
    search.set("lon", String(params.lon));
    search.set("skip", String(params.skip ?? 0));
    search.set("limit", String(params.limit ?? 24));

    const headers: Record<string, string> = {
      "Content-Type": "application/json",
    };
    if (token) {
      headers.Authorization = `Bearer ${token}`;
    }

    const res = await businessAxios.get<RecommendedPostsResponse>(
      `posts/recommend?${search.toString()}`,
      { headers }
    );
    return {
      items: res.data.items || [],
      total: res.data.total || 0,
    };
  },
};
