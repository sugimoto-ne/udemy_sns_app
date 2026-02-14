# Phase 1 - フロントエンド開発TODO

## 🎯 目標
React + TypeScript + MUIでレスポンシブなSNSアプリケーションを構築

---

## 📁 プロジェクトセットアップ

### 1. プロジェクト初期化
- [x] Vite + React + TypeScriptプロジェクト作成
```bash
npm create vite@latest frontend -- --template react-ts
cd frontend
npm install
```

### 2. 必要なパッケージのインストール
- [x] Material-UI
```bash
npm install @mui/material @emotion/react @emotion/styled
npm install @mui/icons-material
```

- [x] ルーティング
```bash
npm install react-router-dom
```

- [x] 状態管理・データフェッチ
```bash
npm install @tanstack/react-query axios
```

- [x] フォーム管理
```bash
npm install react-hook-form
```

- [x] その他
```bash
npm install date-fns
```

### 3. ディレクトリ構成作成
```
frontend/
├── public/
├── src/
│   ├── api/
│   │   ├── client.ts
│   │   ├── auth.ts
│   │   ├── users.ts
│   │   ├── posts.ts
│   │   ├── comments.ts
│   │   ├── likes.ts
│   │   ├── follows.ts
│   │   └── media.ts
│   ├── components/
│   │   ├── common/
│   │   │   ├── AppBar.tsx
│   │   │   ├── Sidebar.tsx
│   │   │   ├── Layout.tsx
│   │   │   ├── Loading.tsx
│   │   │   └── ErrorMessage.tsx
│   │   ├── auth/
│   │   │   ├── LoginForm.tsx
│   │   │   └── RegisterForm.tsx
│   │   ├── post/
│   │   │   ├── PostCard.tsx
│   │   │   ├── PostForm.tsx
│   │   │   ├── PostDetail.tsx
│   │   │   ├── PostList.tsx
│   │   │   └── MediaPreview.tsx
│   │   ├── comment/
│   │   │   ├── CommentList.tsx
│   │   │   ├── CommentItem.tsx
│   │   │   └── CommentForm.tsx
│   │   ├── user/
│   │   │   ├── UserProfile.tsx
│   │   │   ├── UserAvatar.tsx
│   │   │   ├── UserCard.tsx
│   │   │   ├── FollowButton.tsx
│   │   │   └── ProfileEditDialog.tsx
│   │   └── timeline/
│   │       ├── Timeline.tsx
│   │       └── TimelineSwitch.tsx
│   ├── pages/
│   │   ├── LoginPage.tsx
│   │   ├── RegisterPage.tsx
│   │   ├── HomePage.tsx
│   │   ├── PostDetailPage.tsx
│   │   ├── ProfilePage.tsx
│   │   └── ProfileEditPage.tsx
│   ├── hooks/
│   │   ├── useAuth.ts
│   │   ├── usePosts.ts
│   │   ├── useComments.ts
│   │   ├── useLikes.ts
│   │   ├── useFollows.ts
│   │   ├── useUsers.ts
│   │   └── useInfiniteScroll.ts
│   ├── context/
│   │   └── AuthContext.tsx
│   ├── types/
│   │   ├── user.ts
│   │   ├── post.ts
│   │   ├── comment.ts
│   │   └── api.ts
│   ├── utils/
│   │   ├── storage.ts
│   │   ├── formatDate.ts
│   │   └── validation.ts
│   ├── theme/
│   │   └── theme.ts
│   ├── App.tsx
│   ├── main.tsx
│   └── vite-env.d.ts
├── .env.example
├── .env
├── index.html
├── package.json
├── tsconfig.json
└── vite.config.ts
```

- [x] 上記ディレクトリを作成

### 4. 環境設定
- [x] `.env.example` 作成
```env
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

- [x] `.env.local` 作成（`.env.example`をコピー）
- [x] `.gitignore` に `.env.local` 追加

---

## 🎨 UI基盤構築

### 5. MUIテーマ設定
- [x] `src/theme/theme.ts` 実装
  - [x] カラーパレット設定（プライマリ、セカンダリ）
  - [x] タイポグラフィ設定
  - [x] レスポンシブブレークポイント設定
  - [ ] ダークモード対応（オプション）

### 6. 型定義
- [x] `src/types/user.ts`
```typescript
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
```

- [x] `src/types/post.ts`
```typescript
export interface Post {
  id: number;
  user: User;
  content: string;
  media: Media[];
  likes_count: number;
  comments_count: number;
  is_liked: boolean;
  is_bookmarked?: boolean;
  created_at: string;
  updated_at: string;
}

export interface Media {
  id: number;
  media_type: 'image' | 'video' | 'audio';
  media_url: string;
  file_size: number;
  duration?: number;
  order_index: number;
}
```

- [x] `src/types/comment.ts`
- [x] `src/types/api.ts` (APIレスポンス型)

---

## 🔐 認証機能実装

### 7. API Client設定
- [x] `src/api/client.ts`
  - [x] axios instance作成
  - [x] ベースURL設定
  - [x] インターセプター（リクエスト: JWT付与、レスポンス: エラーハンドリング）
  - [ ] トークンリフレッシュ（オプション）

### 8. 認証API
- [x] `src/api/auth.ts`
  - [x] `register(email, password, username)`
  - [x] `login(email, password)`
  - [x] `logout()`
  - [x] `getCurrentUser()`

### 9. ローカルストレージ管理
- [x] `src/utils/storage.ts`
  - [x] `setToken(token: string)`
  - [x] `getToken(): string | null`
  - [x] `removeToken()`
  - [x] `setUser(user: User)`
  - [x] `getUser(): User | null`
  - [x] `removeUser()`

### 10. 認証コンテキスト
- [x] `src/context/AuthContext.tsx`
  - [x] AuthProvider実装
  - [x] 状態: `user`, `isAuthenticated`, `isLoading`
  - [x] 関数: `login`, `register`, `logout`
  - [x] 初期化時にトークンからユーザー情報取得

- [x] `src/hooks/useAuth.ts`
  - [x] AuthContextを使用するカスタムフック

### 11. 認証フォーム
- [x] `src/components/auth/LoginForm.tsx`
  - [x] react-hook-form使用
  - [x] メールアドレス、パスワード入力
  - [x] バリデーション
  - [x] エラーメッセージ表示
  - [x] ログイン処理

- [x] `src/components/auth/RegisterForm.tsx`
  - [x] メールアドレス、パスワード、ユーザー名入力
  - [x] バリデーション
  - [x] 登録処理

### 12. 認証ページ
- [x] `src/pages/LoginPage.tsx`
  - [x] LoginForm表示
  - [x] 登録ページへのリンク
  - [x] レスポンシブレイアウト

- [x] `src/pages/RegisterPage.tsx`
  - [x] RegisterForm表示
  - [x] ログインページへのリンク

### 13. ルーティング設定
- [x] `src/App.tsx`
  - [x] React Router設定
  - [x] 認証ルート（ログイン、登録）
  - [x] プライベートルート（ホーム、プロフィールなど）
  - [x] 未認証時のリダイレクト

---

## 🏠 メインレイアウト

### 14. 共通コンポーネント
- [x] `src/components/common/AppBar.tsx` (Header.tsx として実装)
  - [x] ロゴ
  - [x] ユーザーアバター（認証済み）
  - [x] メニュー（プロフィール、ログアウト）
  - [ ] レスポンシブ対応（ハンバーガーメニュー）

- [ ] `src/components/common/Sidebar.tsx`
  - [ ] ナビゲーションリンク（ホーム、プロフィール、ブックマーク等）
  - [ ] 投稿ボタン
  - [ ] デスクトップ: サイドバー、モバイル: ボトムナビゲーション

- [x] `src/components/common/Layout.tsx` (MainLayout.tsx として実装)
  - [x] AppBar + メインコンテンツ
  - [x] レスポンシブグリッドレイアウト

- [x] `src/components/common/Loading.tsx` (ProtectedRoute内で実装)
  - [x] CircularProgress

- [ ] `src/components/common/ErrorMessage.tsx`
  - [ ] エラー表示コンポーネント

---

## 📝 投稿機能実装

### 15. 投稿API
- [x] `src/api/posts.ts`
  - [x] `getTimeline(type, limit, cursor)`
  - [x] `getPostById(id)`
  - [x] `createPost(content, media_urls)`
  - [x] `updatePost(id, content)`
  - [x] `deletePost(id)`

### 16. 投稿カスタムフック
- [x] `src/hooks/usePosts.ts`
  - [x] `useTimeline(type)` - React Query
  - [x] `usePost(id)` - React Query
  - [x] `useCreatePost()` - Mutation
  - [x] `useUpdatePost()` - Mutation
  - [x] `useDeletePost()` - Mutation

- [ ] `src/hooks/useInfiniteScroll.ts`
  - [ ] Intersection Observer使用
  - [ ] 無限スクロール実装

### 17. 投稿コンポーネント
- [x] `src/components/post/PostCard.tsx`
  - [x] ユーザー情報表示（アバター、名前、ユーザー名）
  - [x] 投稿内容表示
  - [x] メディア表示（画像/動画/音声）
  - [x] いいねボタン、いいね数
  - [x] コメントボタン、コメント数
  - [x] 投稿時刻表示
  - [x] 投稿者の場合: 削除ボタン
  - [x] レスポンシブデザイン

- [x] `src/components/post/MediaPreview.tsx` (PostCard内に実装)
  - [x] 画像プレビュー（Grid表示）
  - [x] 動画プレビュー（再生可能）
  - [x] 音声プレビュー（再生可能）

- [x] `src/components/post/PostForm.tsx`
  - [x] テキストエリア（280文字制限）
  - [ ] 文字数カウンター
  - [ ] メディアアップロードボタン
  - [ ] メディアプレビュー
  - [x] 投稿ボタン
  - [x] Card形式

- [x] `src/components/post/PostList.tsx` (HomePage内に実装)
  - [x] PostCardを配列で表示
  - [ ] 無限スクロール対応
  - [x] ローディング表示

- [x] `src/components/post/PostDetail.tsx` (PostDetailPage内に実装)
  - [x] 投稿詳細表示
  - [x] コメント一覧表示

---

## 🏠 ホームページ（タイムライン）

### 18. タイムライン
- [ ] `src/components/timeline/TimelineSwitch.tsx`
  - [ ] タブ切り替え（全体 / フォロー中）
  - [ ] MUI Tabs使用

- [x] `src/components/timeline/Timeline.tsx` (HomePage内に実装)
  - [x] PostList表示
  - [ ] 無限スクロール
  - [ ] Pull to Refresh（オプション）

- [x] `src/pages/HomePage.tsx`
  - [x] Layout適用
  - [x] PostForm（上部固定）
  - [x] Timeline表示

---

## 💬 コメント機能実装

### 19. コメントAPI
- [x] `src/api/comments.ts`
  - [x] `getComments(postId, limit, cursor)`
  - [x] `createComment(postId, content)`
  - [x] `deleteComment(commentId)`

### 20. コメントカスタムフック
- [x] `src/hooks/useComments.ts`
  - [x] `useComments(postId)` - React Query
  - [x] `useCreateComment()` - Mutation
  - [x] `useDeleteComment()` - Mutation

### 21. コメントコンポーネント
- [x] `src/components/comment/CommentItem.tsx` (CommentList内に実装)
  - [x] ユーザー情報（アバター、名前）
  - [x] コメント内容
  - [x] 投稿時刻
  - [x] 削除ボタン（自分のコメント）

- [x] `src/components/comment/CommentList.tsx`
  - [x] CommentItem配列表示
  - [ ] 無限スクロール対応

- [x] `src/components/comment/CommentForm.tsx`
  - [x] テキスト入力
  - [x] 投稿ボタン

### 22. 投稿詳細ページ
- [x] `src/pages/PostDetailPage.tsx`
  - [x] PostDetail表示
  - [x] CommentList表示
  - [x] CommentForm表示

---

## ❤️ いいね機能実装

### 23. いいねAPI
- [x] `src/api/likes.ts`
  - [x] `likePost(postId)`
  - [x] `unlikePost(postId)`
  - [ ] `getLikes(postId, limit, cursor)`

### 24. いいねカスタムフック
- [x] `src/hooks/useLikes.ts` (usePosts.ts内に実装)
  - [x] `useLikePost()` - Mutation
  - [x] `useUnlikePost()` - Mutation
  - [ ] `useLikes(postId)` - いいね一覧取得

### 25. いいね機能統合
- [x] PostCardにいいねボタン統合
  - [x] いいね状態に応じてアイコン変更（FavoriteBorder / Favorite）
  - [x] いいね数表示
  - [x] クリックでいいね/いいね解除

- [ ] いいね一覧ダイアログ（オプション）
  - [ ] PostCardからいいね数クリックで表示
  - [ ] ユーザーリスト表示

---

## 👥 フォロー機能実装

### 26. フォローAPI
- [x] `src/api/follows.ts`
  - [x] `followUser(username)`
  - [x] `unfollowUser(username)`

### 27. フォローカスタムフック
- [x] `src/hooks/useFollows.ts` (useUsers.ts内に実装)
  - [x] `useFollowUser()` - Mutation
  - [x] `useUnfollowUser()` - Mutation

### 28. フォローボタン
- [x] `src/components/user/FollowButton.tsx`
  - [x] フォロー状態に応じて表示切り替え
  - [x] フォロー/フォロー解除処理

---

## 👤 ユーザープロフィール機能

### 29. ユーザーAPI
- [x] `src/api/users.ts`
  - [x] `getUserByUsername(username)` (getProfile)
  - [x] `updateProfile(data)`
  - [x] `getUserPosts(username, limit, cursor)`
  - [x] `getFollowers(username, limit, cursor)`
  - [x] `getFollowing(username, limit, cursor)`

### 30. ユーザーカスタムフック
- [x] `src/hooks/useUsers.ts`
  - [x] `useUser(username)` - React Query (useUserProfile)
  - [x] `useUpdateProfile()` - Mutation
  - [x] `useUserPosts(username)` - Query
  - [x] `useFollowers(username)` - Query
  - [x] `useFollowing(username)` - Query

### 31. ユーザーコンポーネント
- [x] `src/components/user/UserAvatar.tsx` (MUI Avatar を各コンポーネント内で使用)
  - [x] MUI Avatar
  - [x] デフォルトアバター対応

- [ ] `src/components/user/UserCard.tsx`
  - [ ] ユーザー情報簡易表示
  - [ ] フォローボタン
  - [ ] プロフィールページへのリンク

- [x] `src/components/user/UserProfile.tsx` (UserProfilePage内に実装)
  - [x] ヘッダー画像
  - [x] アバター画像
  - [x] 表示名、ユーザー名
  - [x] 自己紹介
  - [x] ウェブサイト、職業、誕生日
  - [x] フォロー数、フォロワー数
  - [x] フォローボタン（他人のプロフィール）
  - [ ] 編集ボタン（自分のプロフィール）
  - [ ] タブ（投稿 / フォロー / フォロワー）

- [ ] `src/components/user/ProfileEditDialog.tsx`
  - [ ] プロフィール編集フォーム
  - [ ] 画像アップロード（アバター、ヘッダー）
  - [ ] 各種項目編集
  - [ ] 保存ボタン

### 32. プロフィールページ
- [x] `src/pages/ProfilePage.tsx` (UserProfilePage として実装)
  - [x] UserProfile表示
  - [x] 投稿一覧表示
  - [ ] タブ切り替え
    - [x] 投稿タブ: PostList（ユーザーの投稿）
    - [ ] フォロワータブ: UserCard配列
    - [ ] フォロー中タブ: UserCard配列

- [ ] `src/pages/ProfileEditPage.tsx` (オプション)
  - [ ] または ProfileEditDialog で対応

---

## 📷 メディアアップロード実装

### 33. メディアAPI
- [ ] `src/api/media.ts`
  - [ ] `uploadMedia(file: File)`

### 34. メディアアップロード
- [ ] PostForm にメディアアップロード機能統合
  - [ ] ファイル選択ボタン
  - [ ] プレビュー表示
  - [ ] アップロード処理
  - [ ] 進捗表示（オプション）
  - [ ] バリデーション（サイズ、形式）

- [ ] ProfileEditDialog にアバター・ヘッダー画像アップロード機能統合

---

## 🎨 レスポンシブデザイン

### 35. レスポンシブ対応
- [ ] MUIのBreakpoints使用
  - [ ] xs (0px-600px): モバイル
  - [ ] sm (600px-960px): タブレット
  - [ ] md (960px-1280px): デスクトップ小
  - [ ] lg (1280px+): デスクトップ大

- [ ] レイアウト調整
  - [ ] モバイル: 1カラム、ボトムナビゲーション
  - [ ] タブレット: 2カラム
  - [ ] デスクトップ: 3カラム（サイドバー + メイン + サイドウィジェット）

- [ ] コンポーネントのレスポンシブ調整
  - [ ] AppBar: モバイルでハンバーガーメニュー
  - [ ] PostCard: 画像サイズ調整
  - [ ] UserProfile: レイアウト変更

---

## ✅ テスト・最適化

### 36. 基本テスト
- [ ] すべてのページの動作確認
- [ ] 認証フローのテスト
- [ ] 投稿作成・編集・削除のテスト
- [ ] いいね・コメント・フォロー機能のテスト
- [ ] レスポンシブデザインの確認（Chrome DevTools）

### 37. パフォーマンス最適化
- [ ] React.memo 使用（不要な再レンダリング防止）
- [ ] 画像の遅延読み込み（Lazy Loading）
- [ ] React Query のキャッシュ設定
- [ ] コード分割（React.lazy, Suspense）

### 38. UX改善
- [ ] ローディング状態の表示
- [ ] エラーハンドリング
- [ ] 楽観的UI更新（いいね、フォロー）
- [ ] スケルトンローディング（オプション）
- [ ] トースト通知（オプション: notistack）

---

## 📚 ドキュメント

### 39. README作成
- [ ] プロジェクト概要
- [ ] セットアップ手順
- [ ] 開発サーバー起動方法
- [ ] ビルド方法

---

## 🚀 デプロイ準備（Phase 1完了後）

### 40. ビルド設定
- [ ] 環境変数の本番設定（`.env.production`）
- [ ] ビルド実行 (`npm run build`)
- [ ] ビルド結果確認

### 41. Firebase Hostingデプロイ
- [ ] Firebase CLIインストール
```bash
npm install -g firebase-tools
```

- [ ] Firebaseプロジェクト作成
- [ ] Firebase初期化
```bash
firebase init hosting
```

- [ ] デプロイ
```bash
npm run build
firebase deploy
```

---

## ✅ Phase 1 完了チェックリスト

- [x] ユーザー登録・ログインができる
- [x] プロフィールを表示できる
- [ ] プロフィールを編集できる
- [x] 投稿を作成・削除できる
- [ ] 投稿を編集できる
- [x] 投稿にコメントできる
- [x] 投稿にいいねできる
- [x] ユーザーをフォロー/フォロー解除できる
- [x] タイムラインを表示できる
- [ ] タイムラインを切り替えできる（全体 / フォロー中）
- [ ] 無限スクロールが機能する
- [ ] メディアをアップロードして投稿できる
- [ ] レスポンシブデザインが実装されている（基本的なレスポンシブは実装済み）
- [x] 基本機能が正常に動作する

---

## 📝 開発の進め方

1. **プロジェクトセットアップ** (項目1-4)
2. **UI基盤構築** (項目5-6)
3. **認証機能** (項目7-13)
4. **メインレイアウト** (項目14)
5. **投稿機能** (項目15-17)
6. **ホームページ** (項目18)
7. **コメント機能** (項目19-22)
8. **いいね機能** (項目23-25)
9. **フォロー機能** (項目26-28)
10. **ユーザープロフィール** (項目29-32)
11. **メディアアップロード** (項目33-34)
12. **レスポンシブデザイン** (項目35)
13. **テスト・最適化** (項目36-38)
14. **ドキュメント** (項目39)
15. **デプロイ** (項目40-41)

---

**作成日**: 2026-02-14
**最終更新**: 2026-02-14
