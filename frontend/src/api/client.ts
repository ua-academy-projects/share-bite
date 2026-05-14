import axios from 'axios';
import type { 
  LoginRequest, 
  RegisterRequest, 
  AuthResponse, 
  PostResponse, 
  ExploreVenueItem,
  CreatePostInput,
  PaginatedComments,
  CommentResponse,
  CustomerResponse,
  PaginatedAdminUsers,
  FullUserDetails,
  AdminUsersParams,
  NotificationItem,
  CollectionItem,
} from '../types/api';

const apiRoot = axios.create({ baseURL: '/api' });
const API_BASE = import.meta.env.VITE_API_BASE_URL || '';
const authApi = axios.create({ baseURL: `${API_BASE}/api/auth`, withCredentials: true });
const guestApi = axios.create({ baseURL: `${API_BASE}/api/guest`, withCredentials: true });
const businessApi = axios.create({ baseURL: `${API_BASE}/api/business`, withCredentials: true });

const saveSession = (data: AuthResponse) => {
  if (data.access_token) localStorage.setItem('token', data.access_token);
  if (data.refresh_token) localStorage.setItem('refresh_token', data.refresh_token);
};

interface RawPost {
  id: string;
  customer_id: string;
  user_name: string;
  avatar_url?: string;
  venue_id: number;
  text: string;
  rating: number;
  status: string;
  likes_count: number;
  is_liked_by_me: boolean;
  images: string[];
  created_at: string;
  updated_at: string;
  published_at?: string;
}

const mapRawPostToPost = (p: RawPost): PostResponse => ({
  id: p.id,
  customerId: p.customer_id,
  venueId: p.venue_id,
  userName: p.user_name,
  avatarURL: p.avatar_url,
  text: p.text,
  rating: p.rating,
  status: p.status,
  likesCount: p.likes_count,
  isLikedByMe: p.is_liked_by_me,
  images: p.images,
  createdAt: p.created_at,
  updatedAt: p.updated_at,
  publishedAt: p.published_at,
});

guestApi.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  console.log('[Auth Interceptor] Sending token:', !!token);
  if (token && config.headers) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Add JWT interceptor for other instances
const tokenInterceptor = (config: any) => {
  const token = localStorage.getItem('token');
  if (token && config.headers) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
};

apiRoot.interceptors.request.use(tokenInterceptor);
authApi.interceptors.request.use(tokenInterceptor);
businessApi.interceptors.request.use(tokenInterceptor);

// Add 401 response interceptor
const responseErrorInterceptor = (error: any) => {
  if (error.response && error.response.status === 401) {
    localStorage.removeItem('token');
    if (window.location.pathname !== '/auth') {
      window.location.href = '/auth';
    }
  }
  return Promise.reject(error);
};

authApi.interceptors.response.use((res) => res, responseErrorInterceptor);
guestApi.interceptors.response.use((res) => res, responseErrorInterceptor);
businessApi.interceptors.response.use((res) => res, responseErrorInterceptor);

export const apiClient = {
  // Auth
  login: async (data: LoginRequest) => {
    const res = await authApi.post<AuthResponse>('/login', data);
    saveSession(res.data);
    return res.data;
  },

  register: async (data: RegisterRequest) => {
    const res = await authApi.post<AuthResponse>('/register', data);
    saveSession(res.data);
    return res.data;
  },
  
  // Guest
  getPosts: async (limit = 10, offset = 0, venueId?: string | number, customerId?: string) => {
    const params: any = { limit, offset };
    if (venueId) params.venue_id = venueId;
    if (customerId) params.customer_id = customerId;

    const res = await guestApi.get<any>('/posts/', {
      params
    });
    
    console.log('[API getPosts] Raw backend response:', res.data.posts?.map((p: any) => ({ id: p.id, snake_case_liked: p.is_liked_by_me, camelCase_liked: p.isLikedByMe })));

    const rawPosts = (res.data.posts || []) as RawPost[];
    const mappedPosts = rawPosts.map(mapRawPostToPost);

    return {
      Posts: mappedPosts,
      Total: res.data.total || 0
    };
  },
  getExploreNearby: async (lat: number, lon: number, limit = 10) => {
    const res = await guestApi.get<ExploreVenueItem[]>('/posts/explore', {
      params: { lat, lon, limit }
    });
    return res.data;
  },
  createPost: async (data: CreatePostInput) => {
    // Step 1: Create draft
    const createFormData = new FormData();
    createFormData.append('venueId', data.venueId.toString());
    createFormData.append('text', data.text);
    createFormData.append('rating', data.rating.toString());
    if (data.images && data.images.length > 0) {
      data.images.forEach(img => createFormData.append('images', img));
    }
    const createRes = await guestApi.post<{post: RawPost}>('/posts/', createFormData);
    
    const newPostId = createRes.data.post.id;
    const initialPost = mapRawPostToPost(createRes.data.post);
    
    // Step 2: Publish with a simple retry loop
    const patchFormData = new FormData();
    patchFormData.append('status', 'published');
    
    let attempts = 0;
    const maxAttempts = 2;
    while (attempts < maxAttempts) {
      try {
        const patchRes = await guestApi.patch<{post: RawPost}>(`/posts/${newPostId}`, patchFormData);
        return mapRawPostToPost(patchRes.data.post);
      } catch (err) {
        attempts++;
        console.warn(`Publish attempt ${attempts} failed for post ${newPostId}:`, err);
        if (attempts >= maxAttempts) {
          throw Object.assign(
            new Error(`Post was created as draft, but publishing failed after ${attempts} attempts.`),
            { draftPost: createRes.data.post }
          );
        }
      }
    }
    
    return initialPost;
  },
  updatePost: async (postId: string | number, data: { text?: string; rating?: number; venueId?: number; status?: string; images?: File[] }) => {
    const formData = new FormData();
    if (data.text !== undefined) formData.append('text', data.text);
    if (data.rating !== undefined) formData.append('rating', data.rating.toString());
    if (data.venueId !== undefined) formData.append('venue_id', data.venueId.toString());
    if (data.status !== undefined) formData.append('status', data.status);
    if (data.images !== undefined) {
      if (data.images.length > 0) {
        data.images.forEach(img => formData.append('images', img));
      } else {
        formData.append('images_cleared', 'true');
      }
    }
    const res = await guestApi.patch<{post: RawPost}>(`/posts/${postId}`, formData);
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
    const res = await businessApi.get(`/org-units/${id}`);
    return res.data;
  },
  getNearbyVenues: async (lat: number = 50.4501, lon: number = 30.5234) => {
    const res = await businessApi.get<{Items: {id: number, name: string, avatar?: string}[], Total: number}>('/nearby-boxes', {
      params: { lat, lon, limit: 50 }
    });
    return res.data;
  },
  getUser: async (username: string) => {
    const res = await guestApi.get(`customers/${username}`);
    return res.data.customer;
  },
  createCustomer: async (data: { userName: string; firstName: string; lastName: string; bio?: string }) => {
    const res = await guestApi.post<{customerId: string}>('customers/', data);
    return res.data;
  },
  getCurrentCustomer: async () => {
    const res = await guestApi.get<{customer: CustomerResponse}>('customers/');
    return res.data.customer;
  },
  updateCustomer: async (data: { userName?: string; firstName?: string; lastName?: string; bio?: string }) => {
    const res = await guestApi.patch<{customer: CustomerResponse}>('customers/', data);
    return res.data.customer;
  },
  
  // Comments
  getComments: async (postId: string | number, limit = 20, offset = 0) => {
    const res = await guestApi.get<PaginatedComments>(`/posts/${postId}/comments`, {
      params: { take: limit, skip: offset }
    });
    return res.data;
  },
  createComment: async (postId: string | number, text: string) => {
    const res = await guestApi.post<{comment: CommentResponse}>(`/posts/${postId}/comments/`, { text });
    return res.data.comment;
  },
  updateComment: async (postId: string | number, commentId: number, text: string) => {
    const res = await guestApi.patch<{comment: CommentResponse}>(`/posts/${postId}/comments/${commentId}`, { text });
    return res.data.comment;
  },
  deleteComment: async (postId: string | number, commentId: number) => {
    await guestApi.delete(`/posts/${postId}/comments/${commentId}`);
  },

  // OAuth
  oauthCallback: async (provider: string, code: string, slug: string) => {
    const res = await authApi.post<AuthResponse>(`/oauth/${provider}/callback`, { code, slug });
    saveSession(res.data);
    return res.data;
  },

  // Admin
  adminGetUsers: async (params: AdminUsersParams = {}) => {
    const res = await apiRoot.get<PaginatedAdminUsers>('/admin/users', { params });
    return res.data;
  },

  adminGetUserDetails: async (id: string) => {
    const res = await apiRoot.get<FullUserDetails>(`/admin/users/${id}`);
    return res.data;
  },

  adminChangeUserRole: async (id: string, roleSlug: string) => {
    const res = await apiRoot.patch<{ message: string }>(`/admin/users/${id}/role`, { role_slug: roleSlug });
    return res.data;
  },

  updateUserStatus: async (userId: string, status: string) => {
    const res = await apiRoot.put<{ message: string }>(`/users/${userId}/status`, { status });
    return res.data;
  },
  logout: async () => {
    const refreshToken = localStorage.getItem('refresh_token');

    if (!refreshToken) return;

    try {
      const res = await apiRoot.post<{ message: string }>('/user/logout', {
        refresh_token: refreshToken
      });
      return res.data;
    } catch (error) {
      console.error("Logout API failed", error);
      throw error;
    }
  },

  revokeAllSessions: async () => {
    const res = await apiRoot.post<{ message: string }>('/user/sessions/revoke-all');
    return res.data;
  },

  // Notifications (Mocked/Stubbed for now based on typical implementation)
  getNotifications: async (limit = 20, offset = 0) => {
    try {
      const res = await guestApi.get<{ items: NotificationItem[], total: number }>('/notifications', {
        params: { limit, offset }
      });
      return res.data;
    } catch {
      // Mock fallback if endpoint doesn't exist
      return { items: [], total: 0 };
    }
  },
  markNotificationAsRead: async (id: number) => {
    try {
      await guestApi.patch(`/notifications/${id}/read`);
    } catch { }
  },

  // Collections
  getCollections: async () => {
    try {
      const res = await guestApi.get<{ collections: CollectionItem[] }>('/collections');
      return res.data.collections;
    } catch {
      return [];
    }
  },
  createCollection: async (data: { name: string, description?: string, isPublic: boolean }) => {
    const res = await guestApi.post<{ collection: CollectionItem }>('/collections', data);
    return res.data.collection;
  },
  getCollection: async (id: number) => {
    const res = await guestApi.get<{ collection: CollectionItem, items: any[] }>(`/collections/${id}`);
    return res.data;
  },
  savePostToCollection: async (collectionId: number, postId: number | string) => {
    await guestApi.post(`/collections/${collectionId}/venues`, { venue_id: postId }); // Assuming backend uses venue_id for posts/venues
  },
  inviteCollaborator: async (collectionId: number, email: string) => {
    await guestApi.post(`/collections/${collectionId}/collaborators/invite`, { email });
  },
};
