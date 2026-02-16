import { apiClient } from './client';

export interface EmailVerifyData {
  token: string;
}

/**
 * メールアドレスを認証
 */
export const verifyEmail = async (data: EmailVerifyData): Promise<{ message: string }> => {
  const response = await apiClient.post('/auth/email/verify', data);
  return response.data;
};

/**
 * 認証メールを再送信
 */
export const resendVerificationEmail = async (): Promise<{ message: string }> => {
  const response = await apiClient.post('/auth/email/resend');
  return response.data;
};
