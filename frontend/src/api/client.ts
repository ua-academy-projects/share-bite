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
  CustomerResponse
} from '../types/api';

const authApi = axios.create({ baseURL: '/api/auth' });
const guestApi = axios.create({ baseURL: '/api/guest' });
const businessApi = axios.create({ baseURL: '/api/business' });

// Add JWT interceptor
const tokenInterceptor = (config: any) => {
  const token = localStorage.getItem('token');
  if (token && config.headers) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
};

authApi.interceptors.request.use(tokenInterceptor);
guestApi.interceptors.request.use(tokenInterceptor);
businessApi.interceptors.request.use(tokenInterceptor);

export const apiClient = {
  // Auth
  login: async (data: LoginRequest) => {
    const res = await authApi.post<AuthResponse>('/login', data);
    return res.data;
  },
  register: async (data: RegisterRequest) => {
    const res = await authApi.post<AuthResponse>('/register', data);
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
    const createRes = await guestApi.post<{post: PostResponse}>('/posts/', createFormData, {
      headers: { 
        'Content-Type': 'multipart/form-data'
      }
    });
    
    const newPostId = createRes.data.post.id;
    
    // Step 2: Publish
    const patchFormData = new FormData();
    patchFormData.append('status', 'published');
    
    const patchRes = await guestApi.patch<{post: PostResponse}>(`/posts/${newPostId}`, patchFormData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    });
    
    return patchRes.data.post;
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
  }
};
