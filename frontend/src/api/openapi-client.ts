import createClient from 'openapi-fetch';
import type { paths } from '../types/schema';

const BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1';

// デバッグログ
if (typeof window !== 'undefined') {
  console.log('🔍 [API Client] BASE_URL:', BASE_URL);
  console.log('🔍 [API Client] VITE_API_BASE_URL:', import.meta.env.VITE_API_BASE_URL);
}

// リフレッシュ中フラグ
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

// OpenAPI Fetch クライアントを作成
export const apiClient = createClient<paths>({
  baseUrl: BASE_URL,
  credentials: 'include', // Cookie送信を有効化
});

// レスポンスインターセプター: 401エラー時の自動リフレッシュ
apiClient.use({
  async onResponse({ response, request }) {
    // 401エラーの場合、リフレッシュトークンで再試行
    if (response.status === 401) {
      const requestUrl = new URL(request.url);

      // リフレッシュAPIへのリクエストは再試行しない
      if (requestUrl.pathname.includes('/auth/refresh')) {
        // リフレッシュも失敗したらログインページへ
        const currentPath = window.location.pathname;
        if (currentPath !== '/login' && currentPath !== '/register') {
          window.location.href = '/login';
        }
        return response;
      }

      // ログイン/登録ページでは何もしない
      const currentPath = window.location.pathname;
      if (currentPath === '/login' || currentPath === '/register') {
        return response;
      }

      if (isRefreshing) {
        // 既にリフレッシュ中の場合は、キューに追加して待機
        await new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject });
        });

        // リフレッシュ成功後、元のリクエストを再実行
        const retryResponse = await fetch(request);
        return retryResponse;
      }

      isRefreshing = true;

      try {
        // リフレッシュAPIを呼び出し
        const refreshResponse = await fetch(`${BASE_URL}/auth/refresh`, {
          method: 'POST',
          credentials: 'include',
        });

        if (!refreshResponse.ok) {
          throw new Error('Refresh failed');
        }

        // リフレッシュ成功、キュー内のリクエストを再実行
        processQueue(null);
        isRefreshing = false;

        // 元のリクエストを再実行
        const retryResponse = await fetch(request);
        return retryResponse;
      } catch (refreshError) {
        // リフレッシュ失敗、キュー内のリクエストをすべて拒否
        processQueue(refreshError as Error);
        isRefreshing = false;

        // ログインページへリダイレクト
        window.location.href = '/login';
        return response;
      }
    }

    return response;
  },
});
