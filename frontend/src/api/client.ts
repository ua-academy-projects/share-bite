import axios, {
  type AxiosResponse,
  type InternalAxiosRequestConfig,
} from "axios";
import type {
  LoginRequest,
  RegisterRequest,
  AuthResponse,
  PostResponse,
  CreatePostInput,
  PaginatedComments,
  CommentResponse,
  CustomerResponse,
  PaginatedAdminUsers,
  FullUserDetails,
  AdminUsersParams,
  CollectionItem,
} from "@/types/api";

const API_BASE = import.meta.env.VITE_API_BASE_URL || "";
const apiRoot = axios.create({ baseURL: `${API_BASE}/api` });
const authApi = axios.create({
  baseURL: `${API_BASE}/api/auth`,
  withCredentials: true,
});
const guestApi = axios.create({
  baseURL: `${API_BASE}/api/guest`,
  withCredentials: true,
});
const businessAxios = axios.create({
  baseURL: `${API_BASE}/api/business`,
  withCredentials: true,
});
const refreshInstance = axios.create({
  baseURL: `${API_BASE}/api/auth`,
  withCredentials: true,
});

let isRefreshing = false;
let failedQueue: Array<{
  resolve: (token: string | null) => void;
  reject: (error: unknown) => void;
}> = [];

const processQueue = (error: unknown, token: string | null = null) => {
  failedQueue.forEach((prom) => {
    if (error) prom.reject(error);
    else prom.resolve(token);
  });
  failedQueue = [];
};

const saveSession = (data: AuthResponse) => {
  if (data.access_token) localStorage.setItem("token", data.access_token);
  if (data.refresh_token) localStorage.setItem("refresh_token", data.refresh_token);
};

const markGuestHasCustomer = () => {
  localStorage.setItem("guest_has_customer", "1");
};

type RawCustomerResponse = CustomerResponse & { avatarUrl?: string | null };

const mapCustomer = (c: RawCustomerResponse): CustomerResponse => ({
  ...c,
  avatarURL: c.avatarURL ?? c.avatarUrl ?? null,
});

type RawPostImage =
  | string
  | {
      url?: string;
      thumbnailURL?: string;
      thumbnail_url?: string;
      processingStatus?: string;
      processing_status?: string;
    };

interface RawPost {
  id: string;
  customer_id?: string;
  customerId?: string;
  customer_username?: string;
  user_name?: string;
  userName?: string;
  avatar_url?: string;
  avatarURL?: string;
  venue_id?: number;
  venueId?: number;
  text: string;
  rating: number;
  status: string;
  likes_count?: number;
  likesCount?: number;
  is_liked_by_me?: boolean;
  isLikedByMe?: boolean;
  images?: RawPostImage[];
  created_at?: string;
  createdAt?: string;
  updated_at?: string;
  updatedAt?: string;
  published_at?: string;
  publishedAt?: string;
}

const mapPostImages = (images?: RawPostImage[]): string[] => {
  if (!images?.length) return [];

  return images
    .map((image) => {
      if (typeof image === "string") return image;
      return image.url || image.thumbnailURL || image.thumbnail_url || "";
    })
    .filter((url) => url.length > 0);
};

const mapRawPostToPost = (p: RawPost): PostResponse => ({
  id: p.id,
  customerId: p.customerId || p.customer_id || "",
  customerUsername: p.customer_username || p.user_name || p.userName || "",
  userName: p.user_name || p.userName || "",
  avatarURL: p.avatar_url || p.avatarURL,
  venueId: p.venueId ?? p.venue_id ?? 0,
  text: p.text,
  rating: p.rating,
  status: p.status,
  likesCount: p.likes_count ?? p.likesCount ?? 0,
  isLikedByMe: p.is_liked_by_me ?? p.isLikedByMe ?? false,
  images: mapPostImages(p.images),
  createdAt: p.created_at || p.createdAt || new Date().toISOString(),
  updatedAt: p.updated_at || p.updatedAt || new Date().toISOString(),
  publishedAt: p.published_at || p.publishedAt,
});

interface RawCollection {
  id: string | number;
  name: string;
  description?: string;
  isPublic?: boolean;
  is_public?: boolean;
  createdAt?: string;
  created_at?: string;
  updatedAt?: string;
  updated_at?: string;
}

const mapRawCollection = (c: RawCollection): CollectionItem => ({
  id: String(c.id),
  name: c.name,
  description: c.description || "",
  isPublic: c.isPublic ?? c.is_public ?? false,
  createdAt: c.createdAt || c.created_at || new Date().toISOString(),
  updatedAt: c.updatedAt || c.updated_at || new Date().toISOString(),
});

const isPublicGuestRead = (config: InternalAxiosRequestConfig) => {
  if ((config.method || "get").toLowerCase() !== "get") return false;
  const normalized = (config.url || "").replace(/^\//, "").split("?")[0] || "";
  if (normalized === "posts/" || normalized === "posts") return true;
  if (/^posts\/[^/]+$/.test(normalized)) return true;
  if (/^customers\/[^/]+$/.test(normalized)) return true;
  return false;
};

const tokenInterceptor = (config: InternalAxiosRequestConfig) => {
  const token = localStorage.getItem("token");
  if (!token) return config;

  const skipBearer =
    localStorage.getItem("guest_has_customer") !== "1" &&
    isPublicGuestRead(config);

  if (!skipBearer) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
};

const guestCustomerInterceptor = (response: AxiosResponse) => {
  const normalized = (response.config.url || "")
    .replace(/^\//, "")
    .split("?")[0];
  const m = (response.config.method || "get").toLowerCase();
  if (m === "post" && normalized === "customers/") {
    markGuestHasCustomer();
  }
  if (m === "get" && (normalized === "customers/" || normalized === "customers")) {
    markGuestHasCustomer();
  }
  return response;
};

[apiRoot, authApi, businessAxios].forEach((api) => {
  api.interceptors.request.use(tokenInterceptor);
});

guestApi.interceptors.request.use(tokenInterceptor);
guestApi.interceptors.response.use(guestCustomerInterceptor);

const responseInterceptor = async (error: {
  config?: InternalAxiosRequestConfig & { _retry?: boolean };
  response?: { status?: number };
}) => {
  const originalRequest = error.config;
  if (!originalRequest) return Promise.reject(error);

  if (error.response?.status === 401 && !originalRequest._retry) {
    if (isRefreshing) {
      return new Promise<string | null>((resolve, reject) => {
        failedQueue.push({ resolve, reject });
      })
        .then((token) => {
          if (originalRequest.headers) {
            originalRequest.headers.Authorization = `Bearer ${token}`;
          }
          return axios(originalRequest);
        })
        .catch((err) => Promise.reject(err));
    }

    originalRequest._retry = true;
    isRefreshing = true;

    try {
      const refreshToken = localStorage.getItem("refresh_token");
      if (!refreshToken) {
        processQueue(new Error("No refresh token"), null);
        localStorage.removeItem("token");
        localStorage.removeItem("refresh_token");
        if (window.location.pathname !== "/auth") {
          window.location.href = "/auth";
        }
        return Promise.reject(error);
      }

      const res = await refreshInstance.post<AuthResponse>("/refresh", {
        refresh_token: refreshToken,
      });
      const newToken = res.data.access_token;
      saveSession(res.data);
      if (originalRequest.headers) {
        originalRequest.headers.Authorization = `Bearer ${newToken}`;
      }
      processQueue(null, newToken);
      return axios(originalRequest);
    } catch (refreshError) {
      processQueue(refreshError, null);
      localStorage.removeItem("token");
      localStorage.removeItem("refresh_token");
      if (window.location.pathname !== "/auth") {
        window.location.href = "/auth";
      }
      return Promise.reject(refreshError);
    } finally {
      isRefreshing = false;
    }
  }

  return Promise.reject(error);
};

[apiRoot, guestApi, businessAxios].forEach((api) => {
  api.interceptors.response.use((response) => response, responseInterceptor);
});

authApi.interceptors.response.use(
  (response) => response,
  (error) => Promise.reject(error)
);

export const apiClient = {
  login: async (data: LoginRequest) => {
    const res = await authApi.post<AuthResponse>("/login", data);
    saveSession(res.data);
    return res.data;
  },

  register: async (data: RegisterRequest) => {
    const res = await authApi.post<AuthResponse>("/register", data);
    saveSession(res.data);
    return res.data;
  },

  getPosts: async (
    limit = 10,
    offset = 0,
    venueId?: string | number,
    customerId?: string
  ) => {
    const params: Record<string, string | number> = { limit, offset };
    if (venueId) params.venue_id = venueId;
    if (customerId) params.customer_id = customerId;

    const res = await guestApi.get<{ posts?: RawPost[]; total?: number }>(
      "/posts/",
      { params }
    );
    const rawPosts = (res.data.posts || []) as RawPost[];
    return {
      Posts: rawPosts.map(mapRawPostToPost),
      Total: res.data.total || 0,
    };
  },

  getExploreNearby: async (lat: number, lon: number, limit = 10) => {
    const res = await guestApi.get<Record<string, unknown>[]>("/posts/explore", {
      params: { lat, lon, limit },
    });

    return (res.data || []).map((item) => ({
      ...item,
      venue_id: item.venue_id || item.venueId,
      name: item.name as string | undefined,
      avatar: item.avatar as string | undefined,
      posts: ((item.posts as RawPost[]) || []).map((p) =>
        mapRawPostToPost({
          ...p,
          customer_id: p.customer_id || p.customerId,
          venueId: (item.venue_id || item.venueId) as number,
          text: (p as { content?: string }).content || p.text,
          likesCount: p.likes_count || p.likesCount,
          isLikedByMe: p.is_liked_by_me || p.isLikedByMe,
          createdAt: p.created_at || p.createdAt,
        })
      ),
    }));
  },

  createPost: async (data: CreatePostInput) => {
    const createFormData = new FormData();
    createFormData.append("venueId", data.venueId.toString());
    createFormData.append("text", data.text);
    createFormData.append("rating", data.rating.toString());
    if (data.images?.length) {
      data.images.forEach((img) => createFormData.append("images", img));
    }
    const createRes = await guestApi.post<{ post: RawPost }>(
      "/posts/",
      createFormData
    );
    const newPostId = createRes.data.post.id;
    const patchFormData = new FormData();
    patchFormData.append("status", "published");

    let attempts = 0;
    while (attempts < 2) {
      try {
        const patchRes = await guestApi.patch<{ post: RawPost }>(
          `/posts/${newPostId}`,
          patchFormData
        );
        return mapRawPostToPost(patchRes.data.post);
      } catch {
        attempts++;
        if (attempts >= 2) {
          throw Object.assign(
            new Error(
              `Post was created as draft, but publishing failed after ${attempts} attempts.`
            ),
            { draftPost: createRes.data.post }
          );
        }
      }
    }
    return mapRawPostToPost(createRes.data.post);
  },

  updatePost: async (
    postId: string | number,
    data: {
      text?: string;
      rating?: number;
      venueId?: number;
      status?: string;
      images?: File[];
    }
  ) => {
    const formData = new FormData();
    if (data.text !== undefined) formData.append("text", data.text);
    if (data.rating !== undefined)
      formData.append("rating", data.rating.toString());
    if (data.venueId !== undefined)
      formData.append("venueId", data.venueId.toString());
    if (data.status !== undefined) formData.append("status", data.status);
    if (data.images !== undefined) {
      if (data.images.length > 0) {
        data.images.forEach((img) => formData.append("images", img));
      } else {
        formData.append("images_cleared", "true");
      }
    }
    const res = await guestApi.patch<{ post: RawPost }>(
      `/posts/${postId}`,
      formData
    );
    return mapRawPostToPost(res.data.post);
  },

  deletePost: async (postId: string | number) => {
    await guestApi.delete(`/posts/${postId}`);
  },

  likePost: async (postId: string) => {
    await guestApi.post(`/posts/${postId}/like`);
  },

  unlikePost: async (postId: string) => {
    await guestApi.delete(`/posts/${postId}/like`);
  },

  getRestaurant: async (id: string | number) => {
    const res = await businessAxios.get(`org-units/${id}`);
    return res.data;
  },

  getNearbyVenues: async (lat = 50.4501, lon = 30.5234) => {
    const res = await businessAxios.get<{
      Items: { id: number; name: string; avatar?: string }[];
      Total: number;
    }>("nearby-boxes", { params: { lat, lon, limit: 50 } });
    return res.data;
  },

  getUser: async (username: string) => {
    const res = await guestApi.get<{ customer: RawCustomerResponse }>(
      `customers/${username}`
    );
    return mapCustomer(res.data.customer);
  },

  createCustomer: async (data: {
    email: string;
    userName: string;
    firstName: string;
    lastName: string;
    bio?: string;
  }) => {
    const res = await guestApi.post<{ customerId: string }>("customers/", data);
    markGuestHasCustomer();
    return res.data;
  },

  getCurrentCustomer: async () => {
    const res = await guestApi.get<{ customer: RawCustomerResponse }>("customers/");
    markGuestHasCustomer();
    return mapCustomer(res.data.customer);
  },

  hasCurrentCustomer: async (): Promise<boolean> => {
    try {
      await guestApi.get<{ customer: CustomerResponse }>("customers/");
      markGuestHasCustomer();
      return true;
    } catch (error) {
      if (axios.isAxiosError(error) && error.response?.status === 404) {
        return false;
      }
      throw error;
    }
  },

  updateCustomer: async (data: {
    userName?: string;
    firstName?: string;
    lastName?: string;
    bio?: string;
    avatarObjectKey?: string;
  }) => {
    const res = await guestApi.patch<{ customer: RawCustomerResponse }>(
      "customers/",
      data
    );
    markGuestHasCustomer();
    return mapCustomer(res.data.customer);
  },

  uploadCustomerAvatar: async (file: File) => {
    const formData = new FormData();
    formData.append("image", file, file.name);
    const res = await guestApi.post<RawCustomerResponse>("customers/avatar", formData);
    markGuestHasCustomer();
    return mapCustomer(res.data);
  },

  removeCustomerAvatar: async () => {
    const res = await guestApi.patch<{ customer: RawCustomerResponse }>(
      "customers/",
      { avatarObjectKey: "" }
    );
    markGuestHasCustomer();
    return mapCustomer(res.data.customer);
  },

  getComments: async (postId: string | number, limit = 20, offset = 0) => {
    const res = await guestApi.get<PaginatedComments>(
      `/posts/${postId}/comments`,
      { params: { take: limit, skip: offset } }
    );
    return res.data;
  },

  createComment: async (postId: string | number, text: string) => {
    const res = await guestApi.post<{ comment: CommentResponse }>(
      `/posts/${postId}/comments/`,
      { text }
    );
    return res.data.comment;
  },

  updateComment: async (
    postId: string | number,
    commentId: number,
    text: string
  ) => {
    const res = await guestApi.patch<{ comment: CommentResponse }>(
      `/posts/${postId}/comments/${commentId}`,
      { text }
    );
    return res.data.comment;
  },

  deleteComment: async (postId: string | number, commentId: number) => {
    await guestApi.delete(`/posts/${postId}/comments/${commentId}`);
  },

  oauthCallback: async (provider: string, code: string, slug: string) => {
    const res = await authApi.post<AuthResponse>(`/oauth/${provider}/callback`, {
      code,
      slug,
    });
    saveSession(res.data);
    return res.data;
  },

  adminGetUsers: async (params: AdminUsersParams = {}) => {
    const res = await apiRoot.get<PaginatedAdminUsers>("/admin/users", {
      params,
    });
    return res.data;
  },

  adminGetUserDetails: async (id: string) => {
    const res = await apiRoot.get<FullUserDetails>(`/admin/users/${id}`);
    return res.data;
  },

  adminChangeUserRole: async (id: string, roleSlug: string) => {
    const res = await apiRoot.patch<{ message: string }>(
      `/admin/users/${id}/role`,
      { role_slug: roleSlug }
    );
    return res.data;
  },

  updateUserStatus: async (userId: string, status: string) => {
    const res = await apiRoot.put<{ message: string }>(
      `/users/${userId}/status`,
      { status }
    );
    return res.data;
  },

  logout: async () => {
    const refreshToken = localStorage.getItem("refresh_token");
    if (!refreshToken) return;
    const res = await apiRoot.post<{ message: string }>("/user/logout", {
      refresh_token: refreshToken,
    });
    return res.data;
  },

  revokeAllSessions: async () => {
    const res = await apiRoot.post<{ message: string }>(
      "/user/sessions/revoke-all"
    );
    return res.data;
  },

  getCollections: async () => {
    try {
      const res = await guestApi.get<{ collections?: RawCollection[] } | RawCollection[]>(
        "/collections/me"
      );
      const collectionsArray = Array.isArray(res.data)
        ? res.data
        : res.data.collections || [];
      return (Array.isArray(collectionsArray) ? collectionsArray : []).map(
        mapRawCollection
      );
    } catch {
      return [];
    }
  },

  createCollection: async (data: {
    name: string;
    description?: string;
    isPublic: boolean;
  }) => {
    const res = await guestApi.post<{ collection: RawCollection }>(
      "/collections/",
      data
    );
    return mapRawCollection(res.data.collection);
  },

  getCollection: async (id: string | number) => {
    const res = await guestApi.get<{ collection: RawCollection }>(
      `/collections/${id}`
    );
    return {
      ...res.data,
      collection: mapRawCollection(res.data.collection),
    };
  },

  savePostToCollection: async (
    collectionId: number | string,
    venueId: number | string
  ) => {
    await guestApi.post(`/collections/${collectionId}/venues/${venueId}`);
  },

  inviteCollaborator: async (
    collectionId: number | string,
    email: string
  ) => {
    await guestApi.post(`/collections/${collectionId}/invitations`, { email });
  },
};
