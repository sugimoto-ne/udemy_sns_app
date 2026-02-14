# Phase 2 - バックエンド開発TODO（中優先度機能）

## 🎯 目標
ユーザー体験を向上させる追加機能の実装

---

## 📊 Phase 2 機能概要

- ✅ ハッシュタグ機能
- ✅ 複数画像添付機能
- ✅ ブックマーク機能
- ✅ パスワードリセット機能
- ✅ メールアドレス認証機能

---

## 🏷️ ハッシュタグ機能

### 1. データベース準備
- [ ] マイグレーション実行（hashtagsテーブル、post_hashtagsテーブル）
- [ ] GORMモデル定義
  - [ ] `internal/models/hashtag.go`
  - [ ] `internal/models/post_hashtag.go`

### 2. ハッシュタグ抽出ユーティリティ
- [ ] `internal/utils/hashtag.go`
  - [ ] `ExtractHashtags(content string) []string`
    - [ ] 正規表現で `#[a-zA-Z0-9_\p{L}]+` を抽出
    - [ ] 重複削除
    - [ ] 最大10個に制限

### 3. 投稿作成・更新時のハッシュタグ処理
- [ ] `internal/services/post_service.go` 更新
  - [ ] `CreatePost` 修正
    - [ ] ハッシュタグ抽出
    - [ ] ハッシュタグテーブルに存在確認・作成
    - [ ] post_hashtagsテーブルにリレーション作成
  - [ ] `UpdatePost` 修正
    - [ ] 既存のハッシュタグ関連削除
    - [ ] 新しいハッシュタグで再作成

### 4. ハッシュタグサービス
- [ ] `internal/services/hashtag_service.go`
  - [ ] `GetPostsByHashtag(hashtagName string, limit, cursor int) ([]Post, error)`
  - [ ] `GetTrendingHashtags(limit int) ([]Hashtag, error)`
    - [ ] 過去7日間で最も使用されたハッシュタグ
    - [ ] post_hashtagsテーブルをCOUNTでグループ化

### 5. ハッシュタグハンドラー
- [ ] `internal/handlers/hashtag_handler.go`
  - [ ] `GetPostsByHashtag(c echo.Context) error`
  - [ ] `GetTrendingHashtags(c echo.Context) error`

### 6. ハッシュタグルート
- [ ] `GET /api/v1/hashtags/:name/posts`
- [ ] `GET /api/v1/hashtags/trending`

### 7. 投稿レスポンスにハッシュタグ追加
- [ ] 投稿取得時にハッシュタグ情報をプリロード
- [ ] レスポンスに `hashtags: []string` を含める

---

## 🖼️ 複数画像添付機能

### 8. Mediaモデル更新確認
- [ ] `order_index` カラムが存在することを確認
- [ ] 既存のメディアサービスが複数対応していることを確認

### 9. 投稿作成時の複数メディア対応
- [ ] `internal/services/post_service.go` 更新
  - [ ] `CreatePost` 修正
    - [ ] `media_urls` を配列で受け取る
    - [ ] 最大4件に制限
    - [ ] 各メディアに `order_index` を設定（0, 1, 2, 3）

### 10. メディアアップロードAPI改善
- [ ] `internal/handlers/media_handler.go` 更新
  - [ ] 複数ファイル同時アップロード対応（オプション）
  - [ ] レスポンスに `order_index` 含める

### 11. 投稿取得時のメディアソート
- [ ] メディア取得時に `order_index` でソート

---

## 🔖 ブックマーク機能

### 12. データベース準備
- [ ] マイグレーション実行（bookmarksテーブル）
- [ ] GORMモデル定義
  - [ ] `internal/models/bookmark.go`

### 13. ブックマークサービス
- [ ] `internal/services/bookmark_service.go`
  - [ ] `BookmarkPost(userID, postID uint) error`
    - [ ] 重複チェック
  - [ ] `UnbookmarkPost(userID, postID uint) error`
  - [ ] `GetBookmarks(userID uint, limit, cursor int) ([]Post, error)`
  - [ ] `CheckIfBookmarked(userID, postID uint) bool`

### 14. ブックマークハンドラー
- [ ] `internal/handlers/bookmark_handler.go`
  - [ ] `BookmarkPost(c echo.Context) error`
  - [ ] `UnbookmarkPost(c echo.Context) error`
  - [ ] `GetBookmarks(c echo.Context) error`

### 15. ブックマークルート
- [ ] `POST /api/v1/posts/:id/bookmark` (JWT必須)
- [ ] `DELETE /api/v1/posts/:id/bookmark` (JWT必須)
- [ ] `GET /api/v1/bookmarks` (JWT必須)

### 16. 投稿レスポンスにブックマーク状態追加
- [ ] 投稿取得時に `is_bookmarked` フィールドを含める

---

## 🔑 パスワードリセット機能

### 17. データベース準備
- [ ] password_reset_tokensテーブル作成（マイグレーション）
```sql
CREATE TABLE password_reset_tokens (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL REFERENCES users(id),
  token VARCHAR(255) NOT NULL UNIQUE,
  expires_at TIMESTAMP NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

- [ ] GORMモデル定義
  - [ ] `internal/models/password_reset_token.go`

### 18. パスワードリセットトークン生成
- [ ] `internal/utils/token.go`
  - [ ] `GenerateResetToken() string`
    - [ ] ランダムな64文字の文字列生成

### 19. メール送信サービス
- [ ] メール送信パッケージインストール
```bash
go get github.com/sendgrid/sendgrid-go
# または
go get gopkg.in/gomail.v2
```

- [ ] `internal/services/email_service.go`
  - [ ] `SendPasswordResetEmail(email, token string) error`
    - [ ] リセットリンク生成（フロントエンドURL + token）
    - [ ] メール送信

### 20. パスワードリセットサービス
- [ ] `internal/services/password_reset_service.go`
  - [ ] `RequestPasswordReset(email string) error`
    - [ ] ユーザー存在確認
    - [ ] リセットトークン生成
    - [ ] 有効期限設定（1時間）
    - [ ] トークン保存
    - [ ] メール送信
  - [ ] `ConfirmPasswordReset(token, newPassword string) error`
    - [ ] トークン検証
    - [ ] 有効期限チェック
    - [ ] パスワード更新
    - [ ] トークン削除

### 21. パスワードリセットハンドラー
- [ ] `internal/handlers/password_reset_handler.go`
  - [ ] `RequestPasswordReset(c echo.Context) error`
  - [ ] `ConfirmPasswordReset(c echo.Context) error`

### 22. パスワードリセットルート
- [ ] `POST /api/v1/auth/password-reset/request`
- [ ] `POST /api/v1/auth/password-reset/confirm`

---

## ✉️ メールアドレス認証機能

### 23. データベース準備
- [ ] email_verification_tokensテーブル作成（マイグレーション）
```sql
CREATE TABLE email_verification_tokens (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL REFERENCES users(id),
  token VARCHAR(255) NOT NULL UNIQUE,
  expires_at TIMESTAMP NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

- [ ] GORMモデル定義
  - [ ] `internal/models/email_verification_token.go`

### 24. メール認証トークン生成
- [ ] `internal/utils/token.go`（既存に追加）
  - [ ] `GenerateVerificationToken() string`

### 25. メール認証サービス
- [ ] `internal/services/email_verification_service.go`
  - [ ] `SendVerificationEmail(userID uint, email string) error`
    - [ ] トークン生成
    - [ ] 有効期限設定（24時間）
    - [ ] トークン保存
    - [ ] 確認メール送信
  - [ ] `VerifyEmail(token string) error`
    - [ ] トークン検証
    - [ ] 有効期限チェック
    - [ ] users.email_verified を true に更新
    - [ ] トークン削除
  - [ ] `ResendVerificationEmail(userID uint) error`
    - [ ] 既存トークン削除
    - [ ] 新規トークン送信

### 26. ユーザー登録時のメール送信
- [ ] `internal/services/auth_service.go` 更新
  - [ ] `Register` 修正
    - [ ] ユーザー作成後に確認メール送信

### 27. メール認証ハンドラー
- [ ] `internal/handlers/email_verification_handler.go`
  - [ ] `VerifyEmail(c echo.Context) error`
  - [ ] `ResendVerificationEmail(c echo.Context) error`

### 28. メール認証ルート
- [ ] `POST /api/v1/auth/email/verify`
- [ ] `POST /api/v1/auth/email/resend` (JWT必須)

### 29. メール未認証時の制限（オプション）
- [ ] ミドルウェア作成（オプション）
  - [ ] 特定の機能（投稿作成など）でメール認証を必須にする

---

## 🚀 Firebase Storage統合

### 30. Firebase Admin SDK導入
- [ ] Firebase Admin SDKインストール
```bash
go get firebase.google.com/go/v4
go get google.golang.org/api/option
```

### 31. Firebase設定
- [ ] Firebaseプロジェクト作成
- [ ] サービスアカウントキー取得（JSON）
- [ ] 環境変数に設定
```env
FIREBASE_CREDENTIALS_PATH=/path/to/serviceAccountKey.json
FIREBASE_STORAGE_BUCKET=your-project.appspot.com
```

### 32. Firebase Storage サービス
- [ ] `internal/services/firebase_storage_service.go`
  - [ ] Firebase初期化
  - [ ] `UploadToFirebase(file multipart.File, fileName string) (string, error)`
    - [ ] ファイルアップロード
    - [ ] 公開URL取得
  - [ ] `DeleteFromFirebase(fileName string) error`（オプション）

### 33. メディアサービス更新
- [ ] `internal/services/media_service.go` 更新
  - [ ] 環境変数で保存先切り替え（ローカル / Firebase）
  - [ ] 本番環境ではFirebase Storage使用

---

## ✅ テスト

### 34. Phase 2機能のテスト
- [ ] ハッシュタグ機能
  - [ ] 投稿作成時にハッシュタグが抽出・保存される
  - [ ] ハッシュタグ別投稿一覧が取得できる
  - [ ] トレンドハッシュタグが取得できる

- [ ] 複数画像添付
  - [ ] 最大4枚の画像を投稿できる
  - [ ] order_index順に表示される

- [ ] ブックマーク機能
  - [ ] ブックマーク追加・削除ができる
  - [ ] ブックマーク一覧が取得できる

- [ ] パスワードリセット
  - [ ] リセットメールが送信される
  - [ ] トークンでパスワードをリセットできる
  - [ ] 期限切れトークンは無効

- [ ] メール認証
  - [ ] 登録時に確認メールが送信される
  - [ ] トークンでメール認証できる
  - [ ] 確認メール再送信ができる

- [ ] Firebase Storage
  - [ ] 画像がFirebaseにアップロードされる
  - [ ] 公開URLが取得できる

---

## 📚 ドキュメント更新

### 35. API仕様書更新
- [ ] ハッシュタグエンドポイント追加
- [ ] ブックマークエンドポイント追加
- [ ] パスワードリセットエンドポイント追加
- [ ] メール認証エンドポイント追加

### 36. README更新
- [ ] Phase 2機能の説明追加
- [ ] Firebase設定手順追加
- [ ] メール送信設定手順追加

---

## 🚀 デプロイ

### 37. 環境変数追加
- [ ] Firebase認証情報
- [ ] メール送信設定（SendGrid API Key等）

### 38. 本番デプロイ
- [ ] Renderに環境変数追加
- [ ] 再デプロイ

---

## ✅ Phase 2 完了チェックリスト

- [ ] 投稿にハッシュタグを付けられる
- [ ] ハッシュタグ別に投稿を検索できる
- [ ] トレンドハッシュタグが表示される
- [ ] 1つの投稿に最大4枚の画像を添付できる
- [ ] 投稿をブックマークできる
- [ ] ブックマーク一覧を表示できる
- [ ] パスワードリセット機能が動作する
- [ ] メール認証機能が動作する
- [ ] 画像がFirebase Storageにアップロードされる
- [ ] すべてのPhase 2機能が正常に動作する

---

## 📝 開発の進め方

1. **ハッシュタグ機能** (項目1-7)
2. **複数画像添付** (項目8-11)
3. **ブックマーク機能** (項目12-16)
4. **パスワードリセット** (項目17-22)
5. **メール認証** (項目23-29)
6. **Firebase Storage** (項目30-33)
7. **テスト** (項目34)
8. **ドキュメント** (項目35-36)
9. **デプロイ** (項目37-38)

---

**作成日**: 2026-02-14
**最終更新**: 2026-02-14
