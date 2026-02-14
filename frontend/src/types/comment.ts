import type { User } from './user';

// コメント型定義
export interface Comment {
  id: number;
  user: User;
  post_id: number;
  content: string;
  created_at: string;
  updated_at?: string;
}

// コメント作成リクエスト型
export interface CreateCommentRequest {
  content: string;
}
