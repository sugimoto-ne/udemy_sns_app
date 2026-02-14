import { apiClient } from './openapi-client';
import type { User, UpdateProfileRequest } from '../types/user';
import type { PaginatedResponse } from '../types/api';
import type { components } from '../types/schema';

// 型エイリアス
type UpdateProfileRequestSchema = components['schemas']['handlers.UpdateProfileRequest'];

// バックエンドのレスポンス形式
interface BackendUserResponse {
  data: User;
}

// ユーザープロフィール取得
export const getProfile = async (username: string): Promise<User> => {
  const { data: responseData, error } = await apiClient.GET('/users/{username}', {
    params: {
      path: { username },
    },
  });

  if (error) {
    throw new Error('Failed to fetch profile');
  }

  return (responseData as unknown as BackendUserResponse).data;
};

// プロフィール更新
export const updateProfile = async (data: UpdateProfileRequest): Promise<User> => {
  // UpdateProfileRequestSchema に変換
  const requestBody: UpdateProfileRequestSchema = {
    display_name: data.display_name,
    bio: data.bio,
    avatar_url: data.avatar_url,
    header_url: data.header_url,
    website: data.website,
    birth_date: data.birth_date,
    occupation: data.occupation,
  };

  const { data: responseData, error } = await apiClient.PUT('/users/me', {
    body: requestBody,
  });

  if (error) {
    throw new Error('Failed to update profile');
  }

  return (responseData as unknown as BackendUserResponse).data;
};

// フォロワー一覧取得
export const getFollowers = async (
  username: string,
  cursor?: string,
  limit: number = 20
): Promise<PaginatedResponse<User>> => {
  const { data: responseData, error } = await apiClient.GET('/users/{username}/followers', {
    params: {
      path: { username },
      query: {
        cursor,
        limit,
      },
    },
  });

  if (error) {
    throw new Error('Failed to fetch followers');
  }

  return responseData as unknown as PaginatedResponse<User>;
};

// フォロー中一覧取得
export const getFollowing = async (
  username: string,
  cursor?: string,
  limit: number = 20
): Promise<PaginatedResponse<User>> => {
  const { data: responseData, error } = await apiClient.GET('/users/{username}/following', {
    params: {
      path: { username },
      query: {
        cursor,
        limit,
      },
    },
  });

  if (error) {
    throw new Error('Failed to fetch following');
  }

  return responseData as unknown as PaginatedResponse<User>;
};

// ユーザーをフォロー
export const followUser = async (username: string): Promise<void> => {
  const { error } = await apiClient.POST('/users/{username}/follow', {
    params: {
      path: { username },
    },
  });

  if (error) {
    throw new Error('Failed to follow user');
  }
};

// フォロー解除
export const unfollowUser = async (username: string): Promise<void> => {
  const { error } = await apiClient.DELETE('/users/{username}/follow', {
    params: {
      path: { username },
    },
  });

  if (error) {
    throw new Error('Failed to unfollow user');
  }
};
