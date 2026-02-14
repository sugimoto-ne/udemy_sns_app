// ユーザー型定義
export interface User {
  id: number;
  username: string;
  email?: string;
  display_name: string | null;
  bio: string | null;
  avatar_url: string | null;
  header_url: string | null;
  website: string | null;
  birth_date: string | null;
  occupation: string | null;
  followers_count: number;
  following_count: number;
  is_following?: boolean;
  is_followed_by?: boolean;
  created_at: string;
}

// プロフィール更新用リクエスト型
export interface UpdateProfileRequest {
  display_name?: string;
  bio?: string;
  avatar_url?: string;
  header_url?: string;
  website?: string;
  birth_date?: string;
  occupation?: string;
}
