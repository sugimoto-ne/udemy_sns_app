# Phase 1 - バックエンド開発TODO

## 🎯 目標
基本的なSNS機能（認証、投稿、いいね、コメント、フォロー）のAPIを実装

---

## 📁 プロジェクトセットアップ

### 1. プロジェクト初期化
- [ ] Goプロジェクト初期化（`go mod init`）
- [ ] 必要なパッケージのインストール
  - [ ] Echo: `github.com/labstack/echo/v4`
  - [ ] GORM: `gorm.io/gorm`
  - [ ] PostgreSQL Driver: `gorm.io/driver/postgres`
  - [ ] JWT: `github.com/golang-jwt/jwt/v5`
  - [ ] bcrypt: `golang.org/x/crypto/bcrypt`
  - [ ] Validator: `github.com/go-playground/validator/v10`
  - [ ] godotenv: `github.com/joho/godotenv`
  - [ ] CORS: `github.com/labstack/echo/v4/middleware`

### 2. ディレクトリ構成作成
```
backend/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── database/
│   │   └── database.go
│   ├── models/
│   │   ├── user.go
│   │   ├── post.go
│   │   ├── comment.go
│   │   ├── like.go
│   │   ├── follow.go
│   │   └── media.go
│   ├── handlers/
│   │   ├── auth_handler.go
│   │   ├── user_handler.go
│   │   ├── post_handler.go
│   │   ├── comment_handler.go
│   │   ├── like_handler.go
│   │   ├── follow_handler.go
│   │   └── media_handler.go
│   ├── services/
│   │   ├── auth_service.go
│   │   ├── user_service.go
│   │   ├── post_service.go
│   │   ├── comment_service.go
│   │   ├── like_service.go
│   │   ├── follow_service.go
│   │   └── media_service.go
│   ├── middleware/
│   │   ├── jwt_middleware.go
│   │   ├── cors_middleware.go
│   │   └── error_middleware.go
│   ├── utils/
│   │   ├── jwt.go
│   │   ├── password.go
│   │   ├── validator.go
│   │   └── response.go
│   └── routes/
│       └── routes.go
├── migrations/
├── .env.example
├── .env
├── Dockerfile
├── docker-compose.yml
└── go.mod
```

- [ ] 上記ディレクトリを作成

### 3. 環境設定
- [ ] `.env.example` 作成
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=sns_db
JWT_SECRET=your-secret-key-change-this-in-production
PORT=8080
ENV=development
```

- [ ] `.env` 作成（`.env.example`をコピー）
- [ ] `.gitignore` 作成（`.env`, `tmp/`などを追加）

### 4. Docker設定
- [ ] `Dockerfile` 作成
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
```

- [ ] `docker-compose.yml` 作成
```yaml
version: '3.8'

services:
  db:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: sns_db
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  backend:
    build: ./backend
    ports:
      - "8080:8080"
    depends_on:
      - db
    environment:
      DB_HOST: db
      DB_PORT: 5432
      DB_USER: postgres
      DB_PASSWORD: postgres
      DB_NAME: sns_db
      JWT_SECRET: your-secret-key
      PORT: 8080
      ENV: development
    volumes:
      - ./backend:/app
      - media_data:/app/uploads

volumes:
  postgres_data:
  media_data:
```

---

## 🗄️ データベース・モデル実装

### 5. データベース接続
- [ ] `internal/config/config.go` 実装
  - [ ] 環境変数読み込み
  - [ ] Config構造体定義

- [ ] `internal/database/database.go` 実装
  - [ ] PostgreSQL接続
  - [ ] GORM初期化
  - [ ] 接続プール設定

### 6. モデル定義
- [ ] `internal/models/user.go`
  - [ ] User構造体（JSONタグ、GORMタグ）
  - [ ] BeforeCreateフック（パスワードハッシュ化）
  - [ ] メソッド: `CheckPassword()`

- [ ] `internal/models/post.go`
  - [ ] Post構造体
  - [ ] リレーション: User, Media, Comments, Likes

- [ ] `internal/models/media.go`
  - [ ] Media構造体
  - [ ] リレーション: Post

- [ ] `internal/models/comment.go`
  - [ ] Comment構造体
  - [ ] リレーション: User, Post

- [ ] `internal/models/like.go`
  - [ ] PostLike構造体
  - [ ] リレーション: User, Post

- [ ] `internal/models/follow.go`
  - [ ] Follow構造体
  - [ ] リレーション: Follower(User), Following(User)

### 7. マイグレーション
- [ ] `cmd/server/main.go`でAutoMigrate実行
```go
db.AutoMigrate(
    &models.User{},
    &models.Post{},
    &models.Media{},
    &models.Comment{},
    &models.PostLike{},
    &models.Follow{},
)
```

---

## 🔧 共通機能実装

### 8. ユーティリティ
- [ ] `internal/utils/password.go`
  - [ ] `HashPassword(password string) (string, error)`
  - [ ] `CheckPassword(hashedPassword, password string) bool`

- [ ] `internal/utils/jwt.go`
  - [ ] `GenerateToken(userID uint) (string, error)`
  - [ ] `ValidateToken(tokenString string) (*jwt.Token, error)`
  - [ ] `ExtractUserID(token *jwt.Token) (uint, error)`

- [ ] `internal/utils/validator.go`
  - [ ] カスタムバリデータ設定
  - [ ] エラーメッセージ整形

- [ ] `internal/utils/response.go`
  - [ ] `SuccessResponse(c echo.Context, data interface{}) error`
  - [ ] `ErrorResponse(c echo.Context, code int, message string) error`
  - [ ] `PaginationResponse()`

### 9. ミドルウェア
- [ ] `internal/middleware/jwt_middleware.go`
  - [ ] JWT認証ミドルウェア
  - [ ] トークン検証
  - [ ] ユーザーIDをコンテキストに設定

- [ ] `internal/middleware/cors_middleware.go`
  - [ ] CORS設定

- [ ] `internal/middleware/error_middleware.go`
  - [ ] エラーハンドリング
  - [ ] 統一されたエラーレスポンス

---

## 🔐 認証機能実装

### 10. 認証サービス
- [ ] `internal/services/auth_service.go`
  - [ ] `Register(email, password, username string) (*User, error)`
  - [ ] `Login(email, password string) (*User, string, error)`
  - [ ] `GetCurrentUser(userID uint) (*User, error)`

### 11. 認証ハンドラー
- [ ] `internal/handlers/auth_handler.go`
  - [ ] `Register(c echo.Context) error`
    - [ ] リクエストバリデーション
    - [ ] 重複チェック（email, username）
    - [ ] ユーザー作成
    - [ ] JWT発行
  - [ ] `Login(c echo.Context) error`
    - [ ] 認証情報検証
    - [ ] JWT発行
  - [ ] `GetMe(c echo.Context) error`
    - [ ] 現在のユーザー情報取得

### 12. 認証ルート
- [ ] `internal/routes/routes.go`
  - [ ] `POST /api/v1/auth/register`
  - [ ] `POST /api/v1/auth/login`
  - [ ] `GET /api/v1/auth/me` (JWT必須)

---

## 👤 ユーザー機能実装

### 13. ユーザーサービス
- [ ] `internal/services/user_service.go`
  - [ ] `GetUserByUsername(username string) (*User, error)`
  - [ ] `UpdateProfile(userID uint, data map[string]interface{}) (*User, error)`
  - [ ] `GetUserPosts(username string, limit, cursor int) ([]Post, error)`
  - [ ] `GetFollowers(username string, limit, cursor int) ([]User, error)`
  - [ ] `GetFollowing(username string, limit, cursor int) ([]User, error)`
  - [ ] フォロー状態チェック機能

### 14. ユーザーハンドラー
- [ ] `internal/handlers/user_handler.go`
  - [ ] `GetUserByUsername(c echo.Context) error`
  - [ ] `UpdateProfile(c echo.Context) error`
  - [ ] `GetUserPosts(c echo.Context) error`
  - [ ] `GetFollowers(c echo.Context) error`
  - [ ] `GetFollowing(c echo.Context) error`

### 15. ユーザールート
- [ ] `GET /api/v1/users/:username`
- [ ] `PUT /api/v1/users/me` (JWT必須)
- [ ] `GET /api/v1/users/:username/posts`
- [ ] `GET /api/v1/users/:username/followers`
- [ ] `GET /api/v1/users/:username/following`

---

## 📝 投稿機能実装

### 16. 投稿サービス
- [ ] `internal/services/post_service.go`
  - [ ] `GetTimeline(userID uint, timelineType string, limit, cursor int) ([]Post, error)`
    - [ ] `all`: 全体タイムライン
    - [ ] `following`: フォロー中タイムライン
  - [ ] `GetPostByID(postID uint) (*Post, error)`
  - [ ] `CreatePost(userID uint, content string, mediaURLs []string) (*Post, error)`
  - [ ] `UpdatePost(postID, userID uint, content string) (*Post, error)`
  - [ ] `DeletePost(postID, userID uint) error` (論理削除)
  - [ ] いいね数・コメント数の集計

### 17. 投稿ハンドラー
- [ ] `internal/handlers/post_handler.go`
  - [ ] `GetTimeline(c echo.Context) error`
    - [ ] クエリパラメータ: `type`, `limit`, `cursor`
  - [ ] `GetPostByID(c echo.Context) error`
  - [ ] `CreatePost(c echo.Context) error`
  - [ ] `UpdatePost(c echo.Context) error`
    - [ ] 投稿者チェック
  - [ ] `DeletePost(c echo.Context) error`
    - [ ] 投稿者チェック

### 18. 投稿ルート
- [ ] `GET /api/v1/posts` (JWT任意)
- [ ] `GET /api/v1/posts/:id`
- [ ] `POST /api/v1/posts` (JWT必須)
- [ ] `PUT /api/v1/posts/:id` (JWT必須)
- [ ] `DELETE /api/v1/posts/:id` (JWT必須)

---

## 💬 コメント機能実装

### 19. コメントサービス
- [ ] `internal/services/comment_service.go`
  - [ ] `GetCommentsByPostID(postID uint, limit, cursor int) ([]Comment, error)`
  - [ ] `CreateComment(userID, postID uint, content string) (*Comment, error)`
  - [ ] `DeleteComment(commentID, userID uint) error` (論理削除)

### 20. コメントハンドラー
- [ ] `internal/handlers/comment_handler.go`
  - [ ] `GetComments(c echo.Context) error`
  - [ ] `CreateComment(c echo.Context) error`
  - [ ] `DeleteComment(c echo.Context) error`
    - [ ] コメント投稿者チェック

### 21. コメントルート
- [ ] `GET /api/v1/posts/:id/comments`
- [ ] `POST /api/v1/posts/:id/comments` (JWT必須)
- [ ] `DELETE /api/v1/comments/:id` (JWT必須)

---

## ❤️ いいね機能実装

### 22. いいねサービス
- [ ] `internal/services/like_service.go`
  - [ ] `LikePost(userID, postID uint) error`
    - [ ] 重複チェック
  - [ ] `UnlikePost(userID, postID uint) error`
  - [ ] `GetLikesByPostID(postID uint, limit, cursor int) ([]User, error)`
  - [ ] `CheckIfLiked(userID, postID uint) bool`

### 23. いいねハンドラー
- [ ] `internal/handlers/like_handler.go`
  - [ ] `LikePost(c echo.Context) error`
  - [ ] `UnlikePost(c echo.Context) error`
  - [ ] `GetLikes(c echo.Context) error`

### 24. いいねルート
- [ ] `POST /api/v1/posts/:id/like` (JWT必須)
- [ ] `DELETE /api/v1/posts/:id/like` (JWT必須)
- [ ] `GET /api/v1/posts/:id/likes`

---

## 👥 フォロー機能実装

### 25. フォローサービス
- [ ] `internal/services/follow_service.go`
  - [ ] `FollowUser(followerID, followingID uint) error`
    - [ ] 自分自身のフォロー防止
    - [ ] 重複チェック
  - [ ] `UnfollowUser(followerID, followingID uint) error`
  - [ ] `CheckIfFollowing(followerID, followingID uint) bool`

### 26. フォローハンドラー
- [ ] `internal/handlers/follow_handler.go`
  - [ ] `FollowUser(c echo.Context) error`
  - [ ] `UnfollowUser(c echo.Context) error`

### 27. フォロールート
- [ ] `POST /api/v1/users/:username/follow` (JWT必須)
- [ ] `DELETE /api/v1/users/:username/follow` (JWT必須)

---

## 📷 メディアアップロード実装

### 28. メディアサービス
- [ ] `internal/services/media_service.go`
  - [ ] `UploadMedia(file multipart.File, fileHeader *multipart.FileHeader) (string, error)`
    - [ ] ファイルタイプ検証（画像/動画/音声）
    - [ ] ファイルサイズ検証
    - [ ] ローカルストレージ保存（開発環境）
    - [ ] Firebase Storage保存（本番環境）※Phase 2で実装
    - [ ] URLを返す
  - [ ] `SaveMediaRecord(postID uint, mediaType, mediaURL string, fileSize int64) error`

### 29. メディアハンドラー
- [ ] `internal/handlers/media_handler.go`
  - [ ] `UploadMedia(c echo.Context) error`
    - [ ] `multipart/form-data` 受け取り
    - [ ] バリデーション

### 30. メディアルート
- [ ] `POST /api/v1/media/upload` (JWT必須)
- [ ] 静的ファイル配信: `/uploads/*`

---

## ✅ テスト

### 31. 基本テスト
- [ ] 各エンドポイントの動作確認（Postman/Thunder Client）
- [ ] 認証フローのテスト
- [ ] エラーハンドリングのテスト
- [ ] バリデーションのテスト

### 32. 統合テスト（オプション）
- [ ] ユーザー登録 → ログイン → 投稿作成フロー
- [ ] フォロー → タイムライン取得フロー
- [ ] いいね・コメント機能のテスト

---

## 📚 ドキュメント

### 33. README作成
- [ ] プロジェクト概要
- [ ] セットアップ手順
- [ ] Docker起動方法
- [ ] API使用例

### 34. API仕様書更新
- [ ] 実装したエンドポイントの動作確認
- [ ] レスポンス例の更新

---

## 🚀 デプロイ準備（Phase 1完了後）

### 35. 本番環境対応
- [ ] 環境変数の本番設定
- [ ] CORS設定の調整
- [ ] ログ出力の設定
- [ ] ヘルスチェックエンドポイント追加 (`GET /health`)

### 36. Renderデプロイ
- [ ] Renderアカウント作成
- [ ] PostgreSQLインスタンス作成
- [ ] Webサービス作成（Dockerビルド）
- [ ] 環境変数設定

---

## ✅ Phase 1 完了チェックリスト

- [ ] ユーザー登録・ログインができる
- [ ] プロフィールを編集できる
- [ ] 投稿を作成・編集・削除できる
- [ ] 投稿にコメントできる
- [ ] 投稿にいいねできる
- [ ] ユーザーをフォロー/フォロー解除できる
- [ ] タイムラインを取得できる（全体 / フォロー中）
- [ ] メディアをアップロードできる
- [ ] すべてのエンドポイントが正常に動作する
- [ ] Docker環境で動作する

---

## 📝 開発の進め方

1. **プロジェクトセットアップ** (項目1-4)
2. **データベース・モデル** (項目5-7)
3. **共通機能** (項目8-9)
4. **認証機能** (項目10-12)
5. **ユーザー機能** (項目13-15)
6. **投稿機能** (項目16-18)
7. **コメント機能** (項目19-21)
8. **いいね機能** (項目22-24)
9. **フォロー機能** (項目25-27)
10. **メディアアップロード** (項目28-30)
11. **テスト** (項目31-32)
12. **ドキュメント** (項目33-34)
13. **デプロイ** (項目35-36)

---

**作成日**: 2026-02-14
**最終更新**: 2026-02-14
