import { apiClient } from './client';

export interface Bookmark {
  id: number;
  user_id: number;
  post_id: number;
  created_at: string;
}

export interface BookmarksResponse {
  posts: any[]; // Post型を使用
  pagination: {
    has_more: boolean;
    next_cursor: string;
  };
}

/**
 * ブックマークを追加
 */
export const bookmarkPost = async (postId: number): Promise<{ message: string }> => {
  const response = await apiClient.post(`/posts/${postId}/bookmark`);
  return response.data;
};

/**
 * ブックマークを解除
 */
export const unbookmarkPost = async (postId: number): Promise<{ message: string }> => {
  const response = await apiClient.delete(`/posts/${postId}/bookmark`);
  return response.data;
};

/**
 * ブックマーク一覧を取得
 */
export const getBookmarks = async (params?: {
  limit?: number;
  cursor?: string;
}): Promise<BookmarksResponse> => {
  const response = await apiClient.get('/bookmarks', { params });
  return response.data;
};
