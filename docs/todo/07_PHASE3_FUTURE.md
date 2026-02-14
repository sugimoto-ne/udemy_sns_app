# Phase 3 - 将来的な機能（低優先度）

## 🎯 目標
エンゲージメント向上とユーザー体験の更なる改善

---

## 📊 Phase 3 機能概要

- 🔹 ユーザー検索機能
- 🔹 リツイート（再投稿/シェア）機能
- 🔹 通知機能（いいね・コメント通知）
- 🔹 ダイレクトメッセージ（DM）機能
- 🔹 トレンド/人気投稿表示
- 🔹 ソーシャルログイン（Google、Twitter等）

---

## 🔍 ユーザー検索機能

### バックエンド

#### 1. ユーザー検索サービス
- [ ] `internal/services/user_service.go` 更新
  - [ ] `SearchUsers(query string, limit, cursor int) ([]User, error)`
    - [ ] `username` と `display_name` で部分一致検索
    - [ ] PostgreSQL の `ILIKE` または全文検索使用
    - [ ] スコアリング（完全一致 > 前方一致 > 部分一致）

#### 2. ユーザー検索ハンドラー
- [ ] `internal/handlers/user_handler.go` 更新
  - [ ] `SearchUsers(c echo.Context) error`

#### 3. ルート追加
- [ ] `GET /api/v1/users/search?q=keyword&limit=20&cursor=xxx`

### フロントエンド

#### 4. 検索API
- [ ] `src/api/users.ts` 更新
  - [ ] `searchUsers(query, limit, cursor)`

#### 5. 検索カスタムフック
- [ ] `src/hooks/useUsers.ts` 更新
  - [ ] `useSearchUsers(query)`

#### 6. 検索バー
- [ ] `src/components/common/SearchBar.tsx`
  - [ ] 検索入力フィールド
  - [ ] リアルタイム検索（デバウンス）
  - [ ] オートコンプリート（オプション）

#### 7. 検索結果ページ
- [ ] `src/pages/SearchPage.tsx`
  - [ ] ユーザー検索結果表示
  - [ ] UserCard配列

#### 8. AppBarに検索バー統合
- [ ] `src/components/common/AppBar.tsx` 更新
  - [ ] 検索バー追加

---

## 🔄 リツイート（再投稿/シェア）機能

### バックエンド

#### 9. データベース準備
- [ ] retweets テーブル作成
```sql
CREATE TABLE retweets (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL REFERENCES users(id),
  post_id BIGINT NOT NULL REFERENCES posts(id),
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  UNIQUE(user_id, post_id)
);
```

- [ ] GORMモデル定義
  - [ ] `internal/models/retweet.go`

#### 10. リツイートサービス
- [ ] `internal/services/retweet_service.go`
  - [ ] `RetweetPost(userID, postID uint) error`
  - [ ] `UnretweetPost(userID, postID uint) error`
  - [ ] `GetRetweetsByPostID(postID uint) ([]User, error)`
  - [ ] `CheckIfRetweeted(userID, postID uint) bool`

#### 11. タイムライン取得時にリツイートを含める
- [ ] `internal/services/post_service.go` 更新
  - [ ] タイムラインにリツイートした投稿を含める
  - [ ] レスポンスに `retweeted_by` 情報を含める

#### 12. リツイートハンドラー
- [ ] `internal/handlers/retweet_handler.go`
  - [ ] `RetweetPost(c echo.Context) error`
  - [ ] `UnretweetPost(c echo.Context) error`

#### 13. ルート追加
- [ ] `POST /api/v1/posts/:id/retweet`
- [ ] `DELETE /api/v1/posts/:id/retweet`

### フロントエンド

#### 14. リツイートAPI
- [ ] `src/api/retweets.ts`
  - [ ] `retweetPost(postId)`
  - [ ] `unretweetPost(postId)`

#### 15. リツイートカスタムフック
- [ ] `src/hooks/useRetweets.ts`
  - [ ] `useRetweetPost()` - Mutation
  - [ ] `useUnretweetPost()` - Mutation

#### 16. PostCardにリツイートボタン追加
- [ ] `src/components/post/PostCard.tsx` 更新
  - [ ] リツイートアイコンボタン追加
  - [ ] リツイート数表示
  - [ ] リツイート状態に応じてアイコン変更

#### 17. リツイート表示
- [ ] PostCardにリツイート情報表示
  - [ ] 「○○さんがリツイート」のような表示

---

## 🔔 通知機能

### バックエンド

#### 18. データベース準備
- [ ] notificationsテーブル（既にスキーマ定義済み）をマイグレーション

#### 19. 通知サービス
- [ ] `internal/services/notification_service.go`
  - [ ] `CreateNotification(userID, actorID uint, notifType string, postID, commentID *uint) error`
  - [ ] `GetNotifications(userID uint, limit, cursor int) ([]Notification, error)`
  - [ ] `MarkAsRead(notificationID uint) error`
  - [ ] `MarkAllAsRead(userID uint) error`
  - [ ] `GetUnreadCount(userID uint) (int, error)`

#### 20. 通知トリガー実装
- [ ] いいねサービス更新
  - [ ] `LikePost` 実行時に通知作成
- [ ] コメントサービス更新
  - [ ] `CreateComment` 実行時に通知作成
- [ ] フォローサービス更新
  - [ ] `FollowUser` 実行時に通知作成

#### 21. 通知ハンドラー
- [ ] `internal/handlers/notification_handler.go`
  - [ ] `GetNotifications(c echo.Context) error`
  - [ ] `MarkAsRead(c echo.Context) error`
  - [ ] `MarkAllAsRead(c echo.Context) error`
  - [ ] `GetUnreadCount(c echo.Context) error`

#### 22. ルート追加
- [ ] `GET /api/v1/notifications`
- [ ] `PUT /api/v1/notifications/:id/read`
- [ ] `PUT /api/v1/notifications/read-all`
- [ ] `GET /api/v1/notifications/unread-count`

#### 23. WebSocket/SSE実装（リアルタイム通知）（オプション）
- [ ] WebSocketまたはServer-Sent Events導入
- [ ] リアルタイムで通知をプッシュ

### フロントエンド

#### 24. 通知API
- [ ] `src/api/notifications.ts`
  - [ ] `getNotifications(limit, cursor)`
  - [ ] `markAsRead(notificationId)`
  - [ ] `markAllAsRead()`
  - [ ] `getUnreadCount()`

#### 25. 通知カスタムフック
- [ ] `src/hooks/useNotifications.ts`
  - [ ] `useNotifications()` - Infinite Query
  - [ ] `useUnreadCount()` - Query（ポーリング）
  - [ ] `useMarkAsRead()` - Mutation
  - [ ] `useMarkAllAsRead()` - Mutation

#### 26. 通知アイコン
- [ ] `src/components/common/NotificationIcon.tsx`
  - [ ] AppBarに配置
  - [ ] 未読通知数のバッジ表示
  - [ ] クリックで通知ドロップダウンまたはページへ

#### 27. 通知ページ
- [ ] `src/pages/NotificationsPage.tsx`
  - [ ] 通知一覧表示
  - [ ] 通知タイプ別にアイコン・メッセージ表示
  - [ ] クリックで関連投稿へ遷移

#### 28. リアルタイム通知（オプション）
- [ ] WebSocket接続
- [ ] 新しい通知をリアルタイムで受信・表示

---

## 💬 ダイレクトメッセージ（DM）機能

### バックエンド

#### 29. データベース準備
- [ ] conversationsテーブル作成
```sql
CREATE TABLE conversations (
  id BIGSERIAL PRIMARY KEY,
  user1_id BIGINT NOT NULL REFERENCES users(id),
  user2_id BIGINT NOT NULL REFERENCES users(id),
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

- [ ] messagesテーブル作成
```sql
CREATE TABLE messages (
  id BIGSERIAL PRIMARY KEY,
  conversation_id BIGINT NOT NULL REFERENCES conversations(id),
  sender_id BIGINT NOT NULL REFERENCES users(id),
  content TEXT NOT NULL,
  is_read BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

#### 30. DMサービス
- [ ] `internal/services/conversation_service.go`
  - [ ] `GetOrCreateConversation(user1ID, user2ID uint) (*Conversation, error)`
  - [ ] `GetConversations(userID uint) ([]Conversation, error)`
  - [ ] `GetMessages(conversationID uint, limit, cursor int) ([]Message, error)`
  - [ ] `SendMessage(conversationID, senderID uint, content string) (*Message, error)`
  - [ ] `MarkAsRead(messageID uint) error`

#### 31. DMハンドラー
- [ ] `internal/handlers/dm_handler.go`
  - [ ] `GetConversations(c echo.Context) error`
  - [ ] `GetMessages(c echo.Context) error`
  - [ ] `SendMessage(c echo.Context) error`

#### 32. ルート追加
- [ ] `GET /api/v1/conversations`
- [ ] `GET /api/v1/conversations/:id/messages`
- [ ] `POST /api/v1/conversations/:id/messages`

#### 33. WebSocket実装（リアルタイムメッセージ）
- [ ] WebSocket導入
- [ ] リアルタイムでメッセージ送受信

### フロントエンド

#### 34. DM API
- [ ] `src/api/messages.ts`
  - [ ] `getConversations()`
  - [ ] `getMessages(conversationId, limit, cursor)`
  - [ ] `sendMessage(conversationId, content)`

#### 35. DMカスタムフック
- [ ] `src/hooks/useMessages.ts`
  - [ ] `useConversations()`
  - [ ] `useMessages(conversationId)`
  - [ ] `useSendMessage()`

#### 36. DMコンポーネント
- [ ] `src/components/dm/ConversationList.tsx`
  - [ ] 会話一覧表示

- [ ] `src/components/dm/MessageList.tsx`
  - [ ] メッセージ一覧表示
  - [ ] 自分/相手でレイアウト切り替え

- [ ] `src/components/dm/MessageInput.tsx`
  - [ ] メッセージ入力フォーム

#### 37. DMページ
- [ ] `src/pages/MessagesPage.tsx`
  - [ ] 2カラムレイアウト（会話一覧 + メッセージ）
  - [ ] モバイルでは切り替え

#### 38. ナビゲーションにDMリンク追加
- [ ] Sidebarに追加

---

## 📈 トレンド/人気投稿表示

### バックエンド

#### 39. 人気投稿サービス
- [ ] `internal/services/post_service.go` 更新
  - [ ] `GetTrendingPosts(limit int) ([]Post, error)`
    - [ ] 過去24時間でいいね数・コメント数が多い投稿
    - [ ] スコアリングアルゴリズム実装

#### 40. 人気投稿ハンドラー
- [ ] `internal/handlers/post_handler.go` 更新
  - [ ] `GetTrendingPosts(c echo.Context) error`

#### 41. ルート追加
- [ ] `GET /api/v1/posts/trending`

### フロントエンド

#### 42. 人気投稿API
- [ ] `src/api/posts.ts` 更新
  - [ ] `getTrendingPosts(limit)`

#### 43. 人気投稿カスタムフック
- [ ] `src/hooks/usePosts.ts` 更新
  - [ ] `useTrendingPosts()`

#### 44. 人気投稿ウィジェット
- [ ] `src/components/timeline/TrendingPosts.tsx`
  - [ ] サイドバーに表示
  - [ ] 簡易的なPostCard

---

## 🔐 ソーシャルログイン

### バックエンド

#### 45. OAuth 2.0実装
- [ ] Google OAuth 2.0設定
- [ ] Twitter OAuth設定（オプション）

#### 46. OAuth認証サービス
- [ ] `internal/services/oauth_service.go`
  - [ ] `GoogleLogin(code string) (*User, string, error)`
  - [ ] `TwitterLogin(code string) (*User, string, error)`

#### 47. OAuthハンドラー
- [ ] `internal/handlers/oauth_handler.go`
  - [ ] `GoogleCallback(c echo.Context) error`
  - [ ] `TwitterCallback(c echo.Context) error`

#### 48. ルート追加
- [ ] `GET /api/v1/auth/google`
- [ ] `GET /api/v1/auth/google/callback`
- [ ] `GET /api/v1/auth/twitter`
- [ ] `GET /api/v1/auth/twitter/callback`

### フロントエンド

#### 49. ソーシャルログインボタン
- [ ] `src/components/auth/SocialLoginButtons.tsx`
  - [ ] Google ログインボタン
  - [ ] Twitter ログインボタン

#### 50. ログインページに統合
- [ ] `src/pages/LoginPage.tsx` 更新
  - [ ] SocialLoginButtons追加

---

## 🎨 その他の改善案

### パフォーマンス最適化
- [ ] CDN導入（画像配信）
- [ ] Redis導入（キャッシュ）
- [ ] データベースインデックス最適化
- [ ] N+1問題の解消

### セキュリティ強化
- [ ] レート制限の厳格化
- [ ] CSRF対策
- [ ] XSS対策
- [ ] スパム投稿検出

### アクセシビリティ
- [ ] ARIA属性追加
- [ ] キーボードナビゲーション対応
- [ ] スクリーンリーダー対応

### 国際化（i18n）
- [ ] 多言語対応（日本語・英語）
- [ ] react-i18next導入

### アナリティクス
- [ ] Google Analytics導入
- [ ] ユーザー行動分析

### モバイルアプリ
- [ ] React Native での iOS/Androidアプリ開発

---

## 📝 開発の優先順位

Phase 3は長期的な改善項目です。ユーザーのフィードバックや使用状況に応じて、必要な機能から順次実装していくことを推奨します。

### 推奨実装順序
1. **ユーザー検索** - 基本的な機能として早めに実装
2. **通知機能** - ユーザーエンゲージメント向上
3. **リツイート機能** - SNSとしての拡散機能
4. **トレンド/人気投稿** - コンテンツ発見性向上
5. **ダイレクトメッセージ** - 大規模な機能、慎重に実装
6. **ソーシャルログイン** - UX改善

---

## ✅ Phase 3 完了チェックリスト

- [ ] ユーザーを検索できる
- [ ] 投稿をリツイートできる
- [ ] 通知が届く（いいね・コメント・フォロー）
- [ ] ダイレクトメッセージを送受信できる
- [ ] 人気投稿が表示される
- [ ] Googleアカウントでログインできる
- [ ] すべての機能が安定動作する

---

**作成日**: 2026-02-14
**最終更新**: 2026-02-14
