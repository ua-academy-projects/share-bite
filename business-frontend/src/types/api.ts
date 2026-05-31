// Auth DTOs
export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
  slug: string;
}

export interface AuthResponse {
  access_token: string;
  refresh_token: string;
}

export interface CollectionItem {
  id: string;
  name: string;
  description: string;
  isPublic: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface PostResponse {
  id: string;
  customerId: string;
  customerUsername: string;
  userName: string;
  avatarURL?: string | null;
  venueId: number;
  text: string;
  rating: number;
  status: string;
  likesCount: number;
  isLikedByMe: boolean;
  images: string[];
  createdAt: string;
  updatedAt: string;
  publishedAt?: string;
}

export type Post = PostResponse;

export interface CreatePostInput {
  venueId: number;
  text: string;
  rating: number;
  images?: File[];
}

export interface CustomerResponse {
  id: string;
  userName: string;
  firstName: string;
  lastName: string;
  avatarURL?: string | null;
  bio?: string;
  createdAt: string;
}

export interface CommentCustomer {
  id: string;
  userName: string;
  firstName: string;
  lastName: string;
  avatarURL?: string | null;
}

export interface CommentResponse {
  id: number;
  postId: number;
  text: string;
  createdAt: string;
  updatedAt: string;
  customer: CommentCustomer;
}

export interface PaginatedComments {
  total: number;
  entities: CommentResponse[];
}

export interface AdminUserListItem {
  id: string;
  email: string;
  role_slug: string;
  status: string;
  created_at: string;
}

export interface PaginatedAdminUsers {
  items: AdminUserListItem[];
  total_count: number;
}

export interface CustomerProfileData {
  username: string;
  first_name: string;
  last_name: string;
  avatar_object_key: string;
  bio: string;
}

export interface BusinessProfileData {
  profile_type: string;
  name: string;
  avatar: string;
  banner: string;
  description: string;
  latitude?: number | null;
  longitude?: number | null;
}

export interface FullUserDetails {
  id: string;
  email: string;
  role_slug: string;
  status: string;
  created_at: string;
  customer_profile?: CustomerProfileData | null;
  business_profile?: BusinessProfileData | null;
}

export interface AdminUsersParams {
  limit?: number;
  offset?: number;
  search_email?: string;
  role?: string;
  status?: string;
  sort_order?: string;
}
