import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import * as usersApi from '../api/users';
import * as followsApi from '../api/follows';
import type { UpdateProfileRequest } from '../types/user';

// ユーザープロフィール取得
export const useUserProfile = (username: string) => {
  return useQuery({
    queryKey: ['user', username],
    queryFn: () => usersApi.getProfile(username),
    enabled: !!username,
  });
};

// プロフィール更新
export const useUpdateProfile = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: UpdateProfileRequest) => usersApi.updateProfile(data),
    onSuccess: (updatedUser) => {
      queryClient.invalidateQueries({ queryKey: ['user', updatedUser.username] });
    },
  });
};

// フォロワー一覧取得
export const useFollowers = (username: string, cursor?: string, limit: number = 20) => {
  return useQuery({
    queryKey: ['followers', username, cursor, limit],
    queryFn: () => usersApi.getFollowers(username, cursor, limit),
    enabled: !!username,
  });
};

// フォロー中一覧取得
export const useFollowing = (username: string, cursor?: string, limit: number = 20) => {
  return useQuery({
    queryKey: ['following', username, cursor, limit],
    queryFn: () => usersApi.getFollowing(username, cursor, limit),
    enabled: !!username,
  });
};

// フォロー
export const useFollowUser = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (username: string) => followsApi.followUser(username),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['user'] });
      queryClient.invalidateQueries({ queryKey: ['followers'] });
      queryClient.invalidateQueries({ queryKey: ['following'] });
    },
  });
};

// アンフォロー
export const useUnfollowUser = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (username: string) => followsApi.unfollowUser(username),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['user'] });
      queryClient.invalidateQueries({ queryKey: ['followers'] });
      queryClient.invalidateQueries({ queryKey: ['following'] });
    },
  });
};
