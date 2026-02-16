import { apiClient } from './openapi-client';
import type { User } from '../types/user';
import type { components } from '../types/schema';

// 型エイリアス
type LoginRequest = components['schemas']['handlers.LoginRequest'];
type RegisterRequest = components['schemas']['handlers.RegisterRequest'];

// バックエンドのレスポンス形式: { data: { user } } (トークンはCookieに含まれる)
interface BackendAuthResponse {
  data: {
    user: User;
  };
}

interface BackendUserResponse {
  data: User;
}

// 新規登録レスポンス（管理者承認制）
interface BackendRegisterResponse {
  data: {
    message: string;
    user: {
      username: string;
      email: string;
      status: string;
    };
  };
}

// 新規登録（管理者承認制: トークン発行なし、null返却）
export const register = async (data: RegisterRequest): Promise<null> => {
  const { data: responseData, error } = await apiClient.POST('/auth/register', {
    body: data,
  });

  if (error) {
    // バックエンドのエラーメッセージを含むエラーオブジェクトを投げる
    const apiError: any = new Error('Registration failed');
    apiError.response = { data: error };
    throw apiError;
  }

  // 管理者承認制: ユーザー情報は返さず、承認待ちメッセージのみ
  // AuthContextではnullとして扱う
  return null;
};

// ログイン
export const login = async (data: LoginRequest): Promise<User> => {
  const { data: responseData, error } = await apiClient.POST('/auth/login', {
    body: data,
  });

  if (error) {
    // バックエンドのエラーメッセージを含むエラーオブジェクトを投げる
    const apiError: any = new Error('Login failed');
    apiError.response = { data: error };
    throw apiError;
  }

  // レスポンスから data.user プロパティを取り出す
  const authResponse = (responseData as unknown as BackendAuthResponse).data;

  // トークンはCookieに保存されるため、ここでは何もしない
  return authResponse.user;
};

// ログアウト
export const logout = async (): Promise<void> => {
  // 型定義に含まれていないため、anyにキャストして呼び出し
  const { error } = await (apiClient.POST as any)('/auth/logout', {});

  if (error) {
    console.error('Logout failed:', error);
    // ログアウトは失敗してもCookieをクリアする（サーバー側で削除される）
  }
};

// 全デバイスログアウト
export const revokeAllTokens = async (): Promise<void> => {
  // 型定義に含まれていないため、anyにキャストして呼び出し
  const { error } = await (apiClient.POST as any)('/auth/revoke-all', {});

  if (error) {
    throw new Error('Failed to revoke all tokens');
  }
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
