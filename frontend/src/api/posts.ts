import { apiClient } from './openapi-client';
import type { Post, CreatePostRequest, UpdatePostRequest } from '../types/post';
import type { PaginatedResponse } from '../types/api';
import type { components } from '../types/schema';

// 型エイリアス
type CreatePostRequestSchema = components['schemas']['handlers.CreatePostRequest'];
type UpdatePostRequestSchema = components['schemas']['handlers.UpdatePostRequest'];

// バックエンドのレスポンス形式
interface BackendPostResponse {
  data: Post;
}

// タイムライン取得
export const getTimeline = async (
  type: 'all' | 'following' = 'all',
  cursor?: string,
  limit: number = 20
): Promise<PaginatedResponse<Post>> => {
  const { data: responseData, error } = await apiClient.GET('/posts', {
    params: {
      query: {
        type,
        cursor,
        limit,
      },
    },
  });

  if (error) {
    throw new Error('Failed to fetch timeline');
  }

  return responseData as unknown as PaginatedResponse<Post>;
};

// ユーザーの投稿一覧取得
export const getUserPosts = async (
  username: string,
  cursor?: string,
  limit: number = 20
): Promise<PaginatedResponse<Post>> => {
  const { data: responseData, error } = await apiClient.GET('/users/{username}/posts', {
    params: {
      path: { username },
      query: {
        cursor,
        limit,
      },
    },
  });

  if (error) {
    throw new Error('Failed to fetch user posts');
  }

  return responseData as unknown as PaginatedResponse<Post>;
};

// 投稿詳細取得
export const getPost = async (postId: number): Promise<Post> => {
  const { data: responseData, error } = await apiClient.GET('/posts/{id}', {
    params: {
      path: { id: postId },
    },
  });

  if (error) {
    throw new Error('Failed to fetch post');
  }

  return (responseData as unknown as BackendPostResponse).data;
};

// 投稿作成
export const createPost = async (data: CreatePostRequest): Promise<Post> => {
  // CreatePostRequestSchema に変換
  const requestBody: CreatePostRequestSchema = {
    content: data.content,
  };

  const { data: responseData, error } = await apiClient.POST('/posts', {
    body: requestBody,
  });

  if (error) {
    throw new Error('Failed to create post');
  }

  return (responseData as unknown as BackendPostResponse).data;
};

// 投稿更新
export const updatePost = async (
  postId: number,
  data: UpdatePostRequest
): Promise<Post> => {
  // UpdatePostRequestSchema に変換
  const requestBody: UpdatePostRequestSchema = {
    content: data.content,
  };

  const { data: responseData, error } = await apiClient.PUT('/posts/{id}', {
    params: {
      path: { id: postId },
    },
    body: requestBody,
  });

  if (error) {
    throw new Error('Failed to update post');
  }

  return (responseData as unknown as BackendPostResponse).data;
};

// 投稿削除
export const deletePost = async (postId: number): Promise<void> => {
  const { error } = await apiClient.DELETE('/posts/{id}', {
    params: {
      path: { id: postId },
    },
  });

  if (error) {
    throw new Error('Failed to delete post');
  }
};

// 投稿にいいね
export const likePost = async (postId: number): Promise<void> => {
  const { error } = await apiClient.POST('/posts/{id}/like', {
    params: {
      path: { id: postId },
    },
  });

  if (error) {
    throw new Error('Failed to like post');
  }
};

// いいね解除
export const unlikePost = async (postId: number): Promise<void> => {
  const { error } = await apiClient.DELETE('/posts/{id}/like', {
    params: {
      path: { id: postId },
    },
  });

  if (error) {
    throw new Error('Failed to unlike post');
  }
};
