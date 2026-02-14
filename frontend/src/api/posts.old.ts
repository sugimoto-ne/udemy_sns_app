import { apiClient } from './client';
import type { Post, CreatePostRequest, UpdatePostRequest } from '../types/post';
import type { PaginatedResponse } from '../types/api';

// タイムライン取得
export const getTimeline = async (
  cursor?: string,
  limit: number = 20
): Promise<PaginatedResponse<Post>> => {
  const params = new URLSearchParams();
  if (cursor) params.append('cursor', cursor);
  params.append('limit', limit.toString());

  const response = await apiClient.get<PaginatedResponse<Post>>(
    `/posts/timeline?${params.toString()}`
  );
  return response.data;
};

// ユーザーの投稿一覧取得
export const getUserPosts = async (
  username: string,
  cursor?: string,
  limit: number = 20
): Promise<PaginatedResponse<Post>> => {
  const params = new URLSearchParams();
  if (cursor) params.append('cursor', cursor);
  params.append('limit', limit.toString());

  const response = await apiClient.get<PaginatedResponse<Post>>(
    `/users/${username}/posts?${params.toString()}`
  );
  return response.data;
};

// 投稿詳細取得
export const getPost = async (postId: number): Promise<Post> => {
  const response = await apiClient.get<{ data: Post }>(`/posts/${postId}`);
  return response.data.data;
};

// 投稿作成
export const createPost = async (data: CreatePostRequest): Promise<Post> => {
  const response = await apiClient.post<{ data: Post }>('/posts', data);
  return response.data.data;
};

// 投稿更新
export const updatePost = async (
  postId: number,
  data: UpdatePostRequest
): Promise<Post> => {
  const response = await apiClient.put<{ data: Post }>(`/posts/${postId}`, data);
  return response.data.data;
};

// 投稿削除
export const deletePost = async (postId: number): Promise<void> => {
  await apiClient.delete(`/posts/${postId}`);
};
