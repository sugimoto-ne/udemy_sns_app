import { apiClient } from './client';

export interface Hashtag {
  id: number;
  name: string;
  count: number;
  created_at: string;
}

export interface HashtagPostsResponse {
  posts: any[]; // Post型を使用
  pagination: {
    has_more: boolean;
    next_cursor: string;
  };
}

/**
 * トレンドハッシュタグを取得
 */
export const getTrendingHashtags = async (params?: {
  limit?: number;
}): Promise<{ hashtags: Hashtag[] }> => {
  const response = await apiClient.get('/hashtags/trending', { params });
  return response.data;
};

/**
 * ハッシュタグ別の投稿を取得
 */
export const getPostsByHashtag = async (
  hashtagName: string,
  params?: {
    limit?: number;
    cursor?: string;
  }
): Promise<HashtagPostsResponse> => {
  const response = await apiClient.get(`/hashtags/${encodeURIComponent(hashtagName)}/posts`, {
    params,
  });
  return response.data;
};
