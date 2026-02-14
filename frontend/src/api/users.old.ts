import { apiClient } from './client';
import type { User, UpdateProfileRequest } from '../types/user';
import type { PaginatedResponse } from '../types/api';

// ユーザープロフィール取得
export const getProfile = async (username: string): Promise<User> => {
  const response = await apiClient.get<{ data: User }>(`/users/${username}`);
  return response.data.data;
};

// プロフィール更新
export const updateProfile = async (data: UpdateProfileRequest): Promise<User> => {
  const response = await apiClient.put<{ data: User }>('/users/me', data);
  return response.data.data;
};

// フォロワー一覧取得
export const getFollowers = async (
  username: string,
  cursor?: string,
  limit: number = 20
): Promise<PaginatedResponse<User>> => {
  const params = new URLSearchParams();
  if (cursor) params.append('cursor', cursor);
  params.append('limit', limit.toString());

  const response = await apiClient.get<PaginatedResponse<User>>(
    `/users/${username}/followers?${params.toString()}`
  );
  return response.data;
};

// フォロー中一覧取得
export const getFollowing = async (
  username: string,
  cursor?: string,
  limit: number = 20
): Promise<PaginatedResponse<User>> => {
  const params = new URLSearchParams();
  if (cursor) params.append('cursor', cursor);
  params.append('limit', limit.toString());

  const response = await apiClient.get<PaginatedResponse<User>>(
    `/users/${username}/following?${params.toString()}`
  );
  return response.data;
};
