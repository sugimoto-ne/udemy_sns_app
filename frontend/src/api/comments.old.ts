import { apiClient } from './client';
import type { Comment, CreateCommentRequest } from '../types/comment';
import type { PaginatedResponse } from '../types/api';

// コメント一覧取得
export const getComments = async (
  postId: number,
  cursor?: string,
  limit: number = 20
): Promise<PaginatedResponse<Comment>> => {
  const params = new URLSearchParams();
  if (cursor) params.append('cursor', cursor);
  params.append('limit', limit.toString());

  const response = await apiClient.get<PaginatedResponse<Comment>>(
    `/posts/${postId}/comments?${params.toString()}`
  );
  return response.data;
};

// コメント作成
export const createComment = async (
  postId: number,
  data: CreateCommentRequest
): Promise<Comment> => {
  const response = await apiClient.post<{ data: Comment }>(
    `/posts/${postId}/comments`,
    data
  );
  return response.data.data;
};

// コメント削除
export const deleteComment = async (postId: number, commentId: number): Promise<void> => {
  await apiClient.delete(`/posts/${postId}/comments/${commentId}`);
};
