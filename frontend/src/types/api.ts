// APIレスポンス共通型
export interface ApiResponse<T> {
  data: T;
  message: string;
}

// エラーレスポンス型
export interface ApiError {
  error: {
    code: string;
    message: string;
  };
}

// ページネーション型
export interface Pagination {
  has_more: boolean;
  next_cursor: string | null;
  limit: number;
}

// ページネーション付きレスポンス型
export interface PaginatedResponse<T> {
  data: T[];
  pagination: Pagination;
}

// 認証レスポンス型
export interface AuthResponse {
  user: any; // User型を使用
  token: string;
}

// 登録リクエスト型
export interface RegisterRequest {
  email: string;
  password: string;
  username: string;
}

// ログインリクエスト型
export interface LoginRequest {
  email: string;
  password: string;
}
