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
  Customer,
  UpdateCustomerRequest,
  Collection,
  CreateCollectionRequest,
  UpdateCollectionRequest,
  ListCollectionsResponse,
  ListVenuesResponse,
  ReorderVenueRequest
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
  getCustomerByUsername: async (username: string) => {
    const res = await guestApi.get<{customer: Customer}>(`/customers/${username}`);
    return res.data.customer;
  },
  createCustomer: async (data: { userName: string; firstName: string; lastName: string; bio?: string }) => {
    const res = await guestApi.post<{customerId: string}>('/customers/', data);
    return res.data;
  },
  getCurrentCustomer: async () => {
    const res = await guestApi.get<{customer: Customer}>('/customers/');
    return res.data.customer;
  },
  updateCustomer: async (data: UpdateCustomerRequest) => {
    const res = await guestApi.patch<{customer: Customer}>('/customers/', data);
    return res.data.customer;
  },
  uploadAvatar: async (image: File) => {
    const formData = new FormData();
    formData.append('image', image);
    const res = await guestApi.post<Customer>('/customers/avatar', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    });
    return res.data;
  },

  // Collections
  createCollection: async (data: CreateCollectionRequest) => {
    const res = await guestApi.post<{collection: Collection}>('/collections/', data);
    return res.data.collection;
  },
  listMyCollections: async (pageSize = 20, pageToken?: string) => {
    const res = await guestApi.get<ListCollectionsResponse>('/collections/me', {
      params: { page_size: pageSize, page_token: pageToken }
    });
    return res.data;
  },
  getCollection: async (id: string) => {
    const res = await guestApi.get<{collection: Collection}>(`/collections/${id}`);
    return res.data.collection;
  },
  updateCollection: async (id: string, data: UpdateCollectionRequest) => {
    const res = await guestApi.patch<{collection: Collection}>(`/collections/${id}`, data);
    return res.data.collection;
  },
  deleteCollection: async (id: string) => {
    await guestApi.delete(`/collections/${id}`);
  },
  addVenueToCollection: async (collectionId: string, venueId: number) => {
    await guestApi.post(`/collections/${collectionId}/venues/${venueId}`);
  },
  removeVenueFromCollection: async (collectionId: string, venueId: number) => {
    await guestApi.delete(`/collections/${collectionId}/venues/${venueId}`);
  },
  reorderVenueInCollection: async (collectionId: string, venueId: number, data: ReorderVenueRequest) => {
    await guestApi.post(`/collections/${collectionId}/venues/${venueId}/reorder`, data);
  },
  listVenuesInCollection: async (collectionId: string) => {
    const res = await guestApi.get<ListVenuesResponse>(`/collections/${collectionId}/venues`);
    return res.data.venues;
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
