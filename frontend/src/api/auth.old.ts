import { apiClient } from './client';
import type { AuthResponse, RegisterRequest, LoginRequest } from '../types/api';
import type { User } from '../types/user';

// 新規登録
export const register = async (data: RegisterRequest): Promise<AuthResponse> => {
  const response = await apiClient.post<{ data: AuthResponse }>('/auth/register', data);
  return response.data.data;
};

// ログイン
export const login = async (data: LoginRequest): Promise<AuthResponse> => {
  const response = await apiClient.post<{ data: AuthResponse }>('/auth/login', data);
  return response.data.data;
};

// ログアウト
export const logout = async (): Promise<void> => {
  await apiClient.post('/auth/logout');
};

// 現在のユーザー情報取得
export const getCurrentUser = async (): Promise<User> => {
  const response = await apiClient.get<{ data: User }>('/auth/me');
  return response.data.data;
};
