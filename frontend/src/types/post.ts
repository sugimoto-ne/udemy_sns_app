import type { User } from './user';

// メディア型定義
export interface Media {
  id: number;
  media_type: 'image' | 'video' | 'audio';
  media_url: string;
  file_size: number;
  duration?: number;
  order_index: number;
}

// 投稿型定義
export interface Post {
  id: number;
  user: User;
  content: string;
  media: Media[];
  likes_count: number;
  comments_count: number;
  is_liked: boolean;
  is_bookmarked?: boolean;
  created_at: string;
  updated_at: string;
}

// 投稿作成リクエスト型
export interface CreatePostRequest {
  content: string;
  media_urls?: string[];
}

// 投稿更新リクエスト型
export interface UpdatePostRequest {
  content: string;
}
