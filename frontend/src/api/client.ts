import axios from 'axios';
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
    NotificationItem,
    CollectionItem,
} from '../types/api';

const apiRoot = axios.create({ baseURL: '/api' });
const API_BASE = import.meta.env.VITE_API_BASE_URL || '';
const authApi = axios.create({ baseURL: `${API_BASE}/api/auth`, withCredentials: true });
const guestApi = axios.create({ baseURL: `${API_BASE}/api/guest`, withCredentials: true });
const businessApi = axios.create({ baseURL: `${API_BASE}/api/business`, withCredentials: true });
const refreshInstance = axios.create({ baseURL: `${API_BASE}/api/auth`, withCredentials: true });

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
    images: string[];
    created_at?: string;
    createdAt?: string;
    updated_at?: string;
    updatedAt?: string;
    published_at?: string;
    publishedAt?: string;
}

const mapRawPostToPost = (p: RawPost): PostResponse => ({
    id: p.id,
    customerId: p.customerId || p.customer_id || '',
    customerUsername: p.customer_username || p.user_name || p.userName || '',
    userName: p.user_name || p.userName || '',
    avatarURL: p.avatar_url || p.avatarURL,
    venueId: p.venueId ?? p.venue_id ?? 0,
    text: p.text,
    rating: p.rating,
    status: p.status,
    likesCount: p.likes_count ?? p.likesCount ?? 0,
    isLikedByMe: p.is_liked_by_me ?? p.isLikedByMe ?? false,
    images: p.images || [],
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
    description: c.description || '',
    isPublic: c.isPublic ?? c.is_public ?? false,
    createdAt: c.createdAt || c.created_at || new Date().toISOString(),
    updatedAt: c.updatedAt || c.updated_at || new Date().toISOString(),
});

// Request Interceptors (Adding the token)
const tokenInterceptor = (config: any) => {
    const token = localStorage.getItem('token');
    if (token && config.headers) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
};

[apiRoot, authApi, guestApi, businessApi].forEach(api => {
    api.interceptors.request.use(tokenInterceptor);
});

// Response Interceptors (Handling 401 & Token Refresh)
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

        const res = await guestApi.get<any>('/posts/', { params });

        const rawPosts = (res.data.posts || []) as RawPost[];
        const mappedPosts = rawPosts.map(p => {
            const mapped = mapRawPostToPost(p);
            if (!mapped.venueId) {
                console.warn(`[API getPosts] Post ${p.id} is missing venueId:`, p);
            }
            return mapped;
        });

        return {
            Posts: mappedPosts,
            Total: res.data.total || 0
        };
    },
    getExploreNearby: async (lat: number, lon: number, limit = 10) => {
        const res = await guestApi.get<any>('/posts/explore', {
            params: { lat, lon, limit }
        });

        const mappedItems = (res.data || []).map((item: any) => ({
            ...item,
            venue_id: item.venue_id || item.venueId,
            posts: (item.posts || []).map((p: any) => mapRawPostToPost({
                ...p,
                customer_id: p.customer_id || p.customerId,
                venueId: item.venue_id || item.venueId,
                text: p.content || p.text,
                likesCount: p.likes_count || p.likesCount,
                isLikedByMe: p.is_liked_by_me || p.isLikedByMe,
                createdAt: p.created_at || p.createdAt,
            }))
        }));

        return mappedItems;
    },
    createPost: async (data: CreatePostInput) => {
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
        if (data.venueId !== undefined) formData.append('venueId', data.venueId.toString());
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

    // Notifications
    getNotifications: async (limit = 20, offset = 0) => {
        try {
            const res = await guestApi.get<{ items: NotificationItem[], total: number }>('/notifications', {
                params: { limit, offset }
            });
            return res.data;
        } catch {
            return { items: [], total: 0 };
        }
    },
    markNotificationAsRead: async (id: number) => {
        try {
            await guestApi.patch(`/notifications/${id}/read`);
        } catch { }
    },

    // Collections
    getCollections: async (customerId?: string) => {
        try {
            const path = customerId ? `/customers/${customerId}/collections` : '/collections/me';
            const res = await guestApi.get<any>(path);

            const collectionsArray = res.data.collections || res.data || [];
            return (Array.isArray(collectionsArray) ? collectionsArray : []).map(mapRawCollection);
        } catch (error) {
            console.error("Failed to fetch collections:", error);
            return [];
        }
    },
    createCollection: async (data: { name: string, description?: string, isPublic: boolean }) => {
        const res = await guestApi.post<{ collection: RawCollection }>('/collections/', data);
        return mapRawCollection(res.data.collection);
    },
    getCollection: async (id: string | number) => {
        const res = await guestApi.get<{ collection: RawCollection }>(`/collections/${id}`);
        return {
            ...res.data,
            collection: mapRawCollection(res.data.collection)
        };
    },
    savePostToCollection: async (collectionId: number | string, venueId: number | string) => {
        await guestApi.post(`/collections/${collectionId}/venues/${venueId}`);
    },
    inviteCollaborator: async (collectionId: number | string, email: string) => {
        await guestApi.post(`/collections/${collectionId}/invitations`, { email });
    },
};