import axios from 'axios';
import type { 
  LoginRequest, 
  RegisterRequest, 
  AuthResponse, 
  PostResponse, 
  ExploreVenueItem,
  CreatePostInput
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
  getPosts: async (limit = 10, offset = 0) => {
    const res = await guestApi.get<{Posts: PostResponse[], Total: number}>('/posts/', {
      params: { limit, offset }
    });
    return res.data;
  },
  getExploreNearby: async (lat: number, lon: number, limit = 10) => {
    const res = await guestApi.get<ExploreVenueItem[]>('/posts/explore', {
      params: { lat, lon, limit }
    });
    return res.data;
  },
  createPost: async (data: CreatePostInput) => {
    const formData = new FormData();
    formData.append('venue_id', data.venueId.toString());
    formData.append('text', data.text);
    formData.append('rating', data.rating.toString());
    if (data.images) {
      data.images.forEach(img => formData.append('images', img));
    }
    const res = await axios.post<PostResponse>('/api/guest/posts/', formData, {
      headers: { 
        'Content-Type': 'multipart/form-data',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    });
    return res.data;
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
    const res = await businessApi.get<{Items: {id: number, name: string, avatar?: string}[], Total: number}>('/locations/nearby', {
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
  }
};
