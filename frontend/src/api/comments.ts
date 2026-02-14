import { apiClient } from './openapi-client';
import type { Comment, CreateCommentRequest } from '../types/comment';
import type { PaginatedResponse } from '../types/api';
import type { components } from '../types/schema';

// 型エイリアス
type CreateCommentRequestSchema = components['schemas']['handlers.CreateCommentRequest'];

// バックエンドのレスポンス形式
interface BackendCommentResponse {
  data: Comment;
}

// コメント一覧取得
export const getComments = async (
  postId: number,
  cursor?: string,
  limit: number = 20
): Promise<PaginatedResponse<Comment>> => {
  const { data: responseData, error } = await apiClient.GET('/posts/{id}/comments', {
    params: {
      path: { id: postId },
      query: {
        cursor,
        limit,
      },
    },
  });

  if (error) {
    throw new Error('Failed to fetch comments');
  }

  return responseData as unknown as PaginatedResponse<Comment>;
};

// コメント作成
export const createComment = async (
  postId: number,
  data: CreateCommentRequest
): Promise<Comment> => {
  // CreateCommentRequestSchema に変換
  const requestBody: CreateCommentRequestSchema = {
    content: data.content,
  };

  const { data: responseData, error } = await apiClient.POST('/posts/{id}/comments', {
    params: {
      path: { id: postId },
    },
    body: requestBody,
  });

  if (error) {
    throw new Error('Failed to create comment');
  }

  return (responseData as unknown as BackendCommentResponse).data;
};

// コメント削除
export const deleteComment = async (commentId: number): Promise<void> => {
  const { error } = await apiClient.DELETE('/comments/{id}', {
    params: {
      path: { id: commentId },
    },
  });

  if (error) {
    throw new Error('Failed to delete comment');
  }
};
