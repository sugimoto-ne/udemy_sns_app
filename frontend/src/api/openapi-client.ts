import createClient from 'openapi-fetch';
import type { paths } from '../types/schema';

const BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';

// トークン管理関数
export const getToken = (): string | null => {
  return localStorage.getItem('token');
};

export const setToken = (token: string): void => {
  localStorage.setItem('token', token);
};

export const removeToken = (): void => {
  localStorage.removeItem('token');
};

// OpenAPI Fetch クライアントを作成
export const apiClient = createClient<paths>({
  baseUrl: BASE_URL,
});

// リクエストインターセプター: JWTトークンをヘッダーに追加
apiClient.use({
  async onRequest({ request }) {
    const token = getToken();
    if (token) {
      request.headers.set('Authorization', `Bearer ${token}`);
    }
    return request;
  },
});

// レスポンスインターセプター: エラーハンドリング
apiClient.use({
  async onResponse({ response }) {
    // 401エラーの場合はトークンを削除してログインページにリダイレクト
    if (response.status === 401) {
      removeToken();
      window.location.href = '/login';
    }
    return response;
  },
});
