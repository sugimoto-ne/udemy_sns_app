import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import * as commentsApi from '../api/comments';
import type { CreateCommentRequest } from '../types/comment';

// コメント一覧取得
export const useComments = (postId: number, cursor?: string, limit: number = 20) => {
  return useQuery({
    queryKey: ['comments', postId, cursor, limit],
    queryFn: () => commentsApi.getComments(postId, cursor, limit),
    enabled: !!postId,
  });
};

// コメント作成
export const useCreateComment = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ postId, data }: { postId: number; data: CreateCommentRequest }) =>
      commentsApi.createComment(postId, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['comments', variables.postId] });
      queryClient.invalidateQueries({ queryKey: ['post', variables.postId] });
      queryClient.invalidateQueries({ queryKey: ['timeline'] });
    },
  });
};

// コメント削除
export const useDeleteComment = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ commentId }: { commentId: number; postId?: number }) =>
      commentsApi.deleteComment(commentId),
    onSuccess: (_, variables) => {
      if (variables.postId) {
        queryClient.invalidateQueries({ queryKey: ['comments', variables.postId] });
        queryClient.invalidateQueries({ queryKey: ['post', variables.postId] });
      }
      queryClient.invalidateQueries({ queryKey: ['timeline'] });
    },
  });
};
