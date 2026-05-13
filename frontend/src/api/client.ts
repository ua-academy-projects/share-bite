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
} from '../types/api';

const apiRoot = axios.create({ baseURL: '/api' });
const authApi = axios.create({ baseURL: '/api/auth' });
const guestApi = axios.create({ baseURL: '/api/guest' });
const businessApi = axios.create({ baseURL: '/api/business' });
const refreshInstance = axios.create({ baseURL: '/api/auth' });

let isRefreshing = false;
let failedQueue: any[] = [];

const processQueue = (error: any, token: string | null = null) => {
  failedQueue.forEach(prom => {
    if (error) prom.reject(error);
    else prom.resolve(token);
  });
  failedQueue = [];
};

const saveSession = (data: AuthResponse) => {
  if (data.access_token) localStorage.setItem('token', data.access_token);
  if (data.refresh_token) localStorage.setItem('refresh_token', data.refresh_token);
};

// Add JWT interceptor
const tokenInterceptor = (config: any) => {
  const token = localStorage.getItem('token');
  if (token && config.headers) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
};

apiRoot.interceptors.request.use(tokenInterceptor);
authApi.interceptors.request.use(tokenInterceptor);
guestApi.interceptors.request.use(tokenInterceptor);
businessApi.interceptors.request.use(tokenInterceptor);

const responseInterceptor = async (error: any) => {
  const originalRequest = error.config;
  if (error.response?.status === 401 && !originalRequest._retry) {

    if (isRefreshing) {
      return new Promise((resolve, reject) => {
        failedQueue.push({ resolve, reject });
      }).then(token => {
        originalRequest.headers.Authorization = `Bearer ${token}`;
        return axios(originalRequest);
      }).catch(err => Promise.reject(err));
    }

    originalRequest._retry = true;
    isRefreshing = true;

    try {
      const refreshToken = localStorage.getItem('refresh_token');

      if (!refreshToken) {
        processQueue(new Error('No refresh token'), null);
        localStorage.removeItem('token');
        localStorage.removeItem('refresh_token');
        if (window.location.pathname !== '/auth') {
          window.location.href = '/auth';
        }
        return Promise.reject(error);
      }
      const res = await refreshInstance.post<AuthResponse>('/refresh', {
        refresh_token: refreshToken
      });

      const newToken = res.data.access_token;
      saveSession(res.data);

      originalRequest.headers.Authorization = `Bearer ${newToken}`;
      processQueue(null, newToken);
      return axios(originalRequest);

    } catch (refreshError) {
      processQueue(refreshError, null);
      localStorage.removeItem('token');
      localStorage.removeItem('refresh_token');
      if (window.location.pathname !== '/auth') {
        window.location.href = '/auth';
      }
      return Promise.reject(refreshError);
    } finally {
      isRefreshing = false;
    }
  }

  return Promise.reject(error);
};

[apiRoot, authApi, guestApi, businessApi].forEach(api => {
  api.interceptors.request.use(tokenInterceptor);
});

[apiRoot, guestApi, businessApi].forEach(api => {
  api.interceptors.response.use(response => response, responseInterceptor);
});

authApi.interceptors.response.use(
    response => response,
    error => Promise.reject(error)
);

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

    const res = await guestApi.get<{posts: PostResponse[], total: number}>('/posts/', {
      params
    });
    return {
      Posts: res.data.posts || [],
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
    createFormData.append('venue_id', data.venueId.toString());
    createFormData.append('text', data.text);
    createFormData.append('rating', data.rating.toString());
    if (data.images && data.images.length > 0) {
      data.images.forEach(img => createFormData.append('images', img));
    }
    const createRes = await guestApi.post<{post: PostResponse}>('/posts/', createFormData);
    
    const newPostId = createRes.data.post.id;
    
    // Step 2: Publish with a simple retry loop
    const patchFormData = new FormData();
    patchFormData.append('status', 'published');
    
    let attempts = 0;
    const maxAttempts = 2;
    while (attempts < maxAttempts) {
      try {
        const patchRes = await guestApi.patch<{post: PostResponse}>(`/posts/${newPostId}`, patchFormData);
        return patchRes.data.post;
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
    
    return createRes.data.post;
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
    const res = await guestApi.patch<{post: PostResponse}>(`/posts/${postId}`, formData);
    return res.data.post;
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
    const res = await guestApi.get(`/customers/${username}`);
    return res.data;
  },
  createCustomer: async (data: { userName: string; firstName: string; lastName: string; bio?: string }) => {
    const res = await guestApi.post<{customerId: string}>('/customers/', data);
    return res.data;
  },
  getCurrentCustomer: async () => {
    const res = await guestApi.get<{customer: CustomerResponse}>('/customers/');
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
};
