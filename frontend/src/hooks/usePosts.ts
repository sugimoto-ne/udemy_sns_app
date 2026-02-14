import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import * as postsApi from '../api/posts';
import type { CreatePostRequest, UpdatePostRequest } from '../types/post';

// タイムライン取得
export const useTimeline = (type: 'all' | 'following' = 'all', cursor?: string, limit: number = 20) => {
  return useQuery({
    queryKey: ['timeline', type, cursor, limit],
    queryFn: () => postsApi.getTimeline(type, cursor, limit),
  });
};

// ユーザーの投稿一覧取得
export const useUserPosts = (username: string, cursor?: string, limit: number = 20) => {
  return useQuery({
    queryKey: ['userPosts', username, cursor, limit],
    queryFn: () => postsApi.getUserPosts(username, cursor, limit),
    enabled: !!username,
  });
};

// 投稿詳細取得
export const usePost = (postId: number) => {
  return useQuery({
    queryKey: ['post', postId],
    queryFn: () => postsApi.getPost(postId),
    enabled: !!postId,
  });
};

// 投稿作成
export const useCreatePost = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreatePostRequest) => postsApi.createPost(data),
    onSuccess: () => {
      // タイムラインを再取得
      queryClient.invalidateQueries({ queryKey: ['timeline'] });
    },
  });
};

// 投稿更新
export const useUpdatePost = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ postId, data }: { postId: number; data: UpdatePostRequest }) =>
      postsApi.updatePost(postId, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['post', variables.postId] });
      queryClient.invalidateQueries({ queryKey: ['timeline'] });
    },
  });
};

// 投稿削除
export const useDeletePost = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (postId: number) => postsApi.deletePost(postId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['timeline'] });
      queryClient.invalidateQueries({ queryKey: ['userPosts'] });
    },
  });
};

// いいね追加
export const useLikePost = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (postId: number) => postsApi.likePost(postId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['timeline'] });
      queryClient.invalidateQueries({ queryKey: ['post'] });
    },
  });
};

// いいね削除
export const useUnlikePost = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (postId: number) => postsApi.unlikePost(postId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['timeline'] });
      queryClient.invalidateQueries({ queryKey: ['post'] });
    },
  });
};
