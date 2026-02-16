import axios, { AxiosError } from 'axios';

// Axiosインスタンス作成
export const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
  withCredentials: true, // Cookie送信を有効化
});

// リフレッシュ中フラグ（複数の401エラーが同時に発生した場合の重複防止）
let isRefreshing = false;
let failedQueue: Array<{
  resolve: (value?: unknown) => void;
  reject: (reason?: unknown) => void;
}> = [];

// 失敗したリクエストをキューに追加
const processQueue = (error: Error | null) => {
  failedQueue.forEach((prom) => {
    if (error) {
      prom.reject(error);
    } else {
      prom.resolve();
    }
  });

  failedQueue = [];
};

// レスポンスインターセプター（エラーハンドリング + 自動リフレッシュ）
apiClient.interceptors.response.use(
  (response) => {
    return response;
  },
  async (error: AxiosError) => {
    const originalRequest = error.config as typeof error.config & { _retry?: boolean };

    // 401エラー（認証エラー）の場合、リフレッシュトークンで再試行
    if (error.response?.status === 401 && originalRequest && !originalRequest._retry) {
      // リフレッシュAPIへのリクエストは再試行しない
      if (originalRequest.url?.includes('/auth/refresh')) {
        // リフレッシュも失敗したらログインページへ
        window.location.href = '/login';
        return Promise.reject(error);
      }

      if (isRefreshing) {
        // 既にリフレッシュ中の場合は、キューに追加して待機
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject });
        })
          .then(() => {
            return apiClient(originalRequest);
          })
          .catch((err) => {
            return Promise.reject(err);
          });
      }

      originalRequest._retry = true;
      isRefreshing = true;

      try {
        // リフレッシュAPIを呼び出し
        await apiClient.post('/auth/refresh');

        // リフレッシュ成功、キュー内のリクエストを再実行
        processQueue(null);
        isRefreshing = false;

        // 元のリクエストを再実行
        return apiClient(originalRequest);
      } catch (refreshError) {
        // リフレッシュ失敗、キュー内のリクエストをすべて拒否
        processQueue(refreshError as Error);
        isRefreshing = false;

        // ログインページへリダイレクト
        window.location.href = '/login';
        return Promise.reject(refreshError);
      }
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
