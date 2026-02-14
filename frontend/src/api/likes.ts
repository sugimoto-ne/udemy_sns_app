import { apiClient } from './client';

// いいね追加
export const likePost = async (postId: number): Promise<void> => {
  await apiClient.post(`/posts/${postId}/like`);
};

// いいね削除
export const unlikePost = async (postId: number): Promise<void> => {
  await apiClient.delete(`/posts/${postId}/like`);
};
