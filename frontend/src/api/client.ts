import axios, { AxiosError } from 'axios';
import { getToken, clearAuth } from '../utils/storage';

// Axiosインスタンス作成
export const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// リクエストインターセプター（JWT付与）
apiClient.interceptors.request.use(
  (config) => {
    const token = getToken();
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error: AxiosError) => {
    return Promise.reject(error);
  }
);

// レスポンスインターセプター（エラーハンドリング）
apiClient.interceptors.response.use(
  (response) => {
    return response;
  },
  (error: AxiosError) => {
    // 401エラー（認証エラー）の場合、ローカルストレージをクリアしてログインページへ
    if (error.response?.status === 401) {
      clearAuth();
      // ログインページへリダイレクト（React Routerを使用する場合は別途実装）
      window.location.href = '/login';
    }

    // 403エラー（権限エラー）
    if (error.response?.status === 403) {
      console.error('Access forbidden');
    }

    // 404エラー（Not Found）
    if (error.response?.status === 404) {
      console.error('Resource not found');
    }

    // 500エラー（サーバーエラー）
    if (error.response?.status === 500) {
      console.error('Internal server error');
    }

    return Promise.reject(error);
  }
);
