import { apiClient } from './client';

export interface PasswordResetRequestData {
  email: string;
}

export interface PasswordResetConfirmData {
  token: string;
  new_password: string;
}

/**
 * パスワードリセットをリクエスト
 */
export const requestPasswordReset = async (
  data: PasswordResetRequestData
): Promise<{ message: string }> => {
  const response = await apiClient.post('/auth/password-reset/request', data);
  return response.data;
};

/**
 * パスワードリセットを確認・実行
 */
export const confirmPasswordReset = async (
  data: PasswordResetConfirmData
): Promise<{ message: string }> => {
  const response = await apiClient.post('/auth/password-reset/confirm', data);
  return response.data;
};
