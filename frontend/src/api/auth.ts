import { apiClient, setToken, removeToken } from './openapi-client';
import type { AuthResponse } from '../types/api';
import type { User } from '../types/user';
import type { components } from '../types/schema';

// 型エイリアス
type LoginRequest = components['schemas']['handlers.LoginRequest'];
type RegisterRequest = components['schemas']['handlers.RegisterRequest'];

// バックエンドのレスポンス形式: { data: { user, token } }
interface BackendAuthResponse {
  data: AuthResponse;
}

interface BackendUserResponse {
  data: User;
}

// 新規登録
export const register = async (data: RegisterRequest): Promise<AuthResponse> => {
  const { data: responseData, error } = await apiClient.POST('/auth/register', {
    body: data,
  });

  if (error) {
    throw new Error('Registration failed');
  }

  // レスポンスから data プロパティを取り出す
  const authResponse = (responseData as unknown as BackendAuthResponse).data;

  // トークンを保存
  setToken(authResponse.token);

  return authResponse;
};

// ログイン
export const login = async (data: LoginRequest): Promise<AuthResponse> => {
  const { data: responseData, error } = await apiClient.POST('/auth/login', {
    body: data,
  });

  if (error) {
    throw new Error('Login failed');
  }

  // レスポンスから data プロパティを取り出す
  const authResponse = (responseData as unknown as BackendAuthResponse).data;

  // トークンを保存
  setToken(authResponse.token);

  return authResponse;
};

// ログアウト
export const logout = async (): Promise<void> => {
  // トークンを削除
  removeToken();
};

// 現在のユーザー情報取得
export const getCurrentUser = async (): Promise<User> => {
  const { data: responseData, error } = await apiClient.GET('/auth/me');

  if (error) {
    throw new Error('Failed to get current user');
  }

  // レスポンスから data プロパティを取り出す
  return (responseData as unknown as BackendUserResponse).data;
};
