import { apiClient } from './client';

// フォロー
export const followUser = async (username: string): Promise<void> => {
  await apiClient.post(`/users/${username}/follow`);
};

// アンフォロー
export const unfollowUser = async (username: string): Promise<void> => {
  await apiClient.delete(`/users/${username}/follow`);
};
