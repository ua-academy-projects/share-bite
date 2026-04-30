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

export interface User {
  id: string;
  name: string;
  handle?: string;
  avatar?: string | null;
}

// Guest DTOs
export interface PostResponse {
  id: string;
  customerId: string;
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

export interface PostItem {
  id: string;
  content: string;
  createdAt: string;
  images: string[];
}

export interface ExploreVenueItem {
  venue_id: number;
  posts: PostItem[];
}

export interface CreatePostInput {
  venueId: number;
  text: string;
  rating: number;
  images?: File[];
}

export interface ReviewResponse {
  id: string;
  customerId: string;
  userName: string;
  avatarURL?: string | null;
  venueId: number;
  rating: number;
  text: string;
  createdAt: string;
}

export type Review = ReviewResponse;

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

// Customer DTOs
export interface UpdateCustomerRequest {
  userName?: string;
  firstName?: string;
  lastName?: string;
  bio?: string;
}

export interface Customer {
  id: string;
  userId: string;
  userName: string;
  firstName: string;
  lastName: string;
  avatarUrl?: string | null;
  bio?: string | null;
  createdAt: string;
}

// Collection DTOs
export interface Collection {
  id: string;
  name: string;
  description?: string | null;
  isPublic: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface CollectionVenue {
  id: number;
  name: string;
  description?: string | null;
  avatarUrl?: string | null;
  bannerUrl?: string | null;
  sortOrder: number;
  addedAt: string;
}

export interface CreateCollectionRequest {
  name: string;
  description?: string;
  isPublic: boolean;
}

export interface UpdateCollectionRequest {
  name?: string;
  description?: string;
  isPublic?: boolean;
}

export interface ListCollectionsResponse {
  collections: Collection[];
  next_page_token?: string;
}

export interface ListVenuesResponse {
  venues: CollectionVenue[];
}

export interface ReorderVenueRequest {
  prevVenueId?: number;
  nextVenueId?: number;
}

// Business DTOs
export interface RestaurantResponse {
  id: number;
  name: string;
  category: string;
  rating: number;
  image: string;
  description: string;
  location: string;
}

export type Restaurant = RestaurantResponse;
