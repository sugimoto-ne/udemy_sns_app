# CLAUDE.md

このファイルは、このリポジトリでコードを扱う際にClaude Code (claude.ai/code) にガイダンスを提供します。

---

## プロジェクト概要

これは**TwitterライクなSNSアプリケーション**で、モダンな技術スタックで構築され、3つのフェーズ（MVP → 拡張機能 → 高度な機能）での段階的開発を想定しています。

**技術スタック:**
- **フロントエンド**: React + TypeScript + Material-UI (MUI)
  - 型生成: openapi-typescript (OpenAPI定義から自動生成)
  - APIクライアント: openapi-fetch (型安全なHTTPクライアント)
- **バックエンド**: Go + Echo + GORM
  - API仕様生成: swaggo/echo-swagger (OpenAPI/Swagger定義の自動生成)
  - ドキュメント: Swagger UI
- **データベース**: PostgreSQL
- **認証**: JWT (JSON Web Token)
- **型安全連携**: OpenAPI 3.0を中心としたフロントエンド⇔バックエンド型連携
- **インフラ**: Docker + Docker Compose (ローカル), Render/Cloud Run (バックエンド), Firebase Hosting (フロントエンド)
- **ストレージ**: Docker Volume (ローカル), Firebase Storage (本番)

---

## アーキテクチャ概要

```
Frontend (React/TS) ←→ REST API (/api/v1/*) ←→ Backend (Go/Echo) ←→ PostgreSQL
  ↑ openapi-fetch         ↑ JWT Auth              ↑ swaggo              ↓
  └─ 型自動生成 ←─ swagger.yaml (OpenAPI 3.0) ←─┘           Firebase Storage
     (openapi-typescript)                                       (メディア)
```

### 型安全な開発フロー

1. **バックエンド開発**: Goコードにswaggoアノテーション追加
2. **OpenAPI生成**: `swag init` で `docs/swagger.yaml` を自動生成
3. **型定義生成**: `openapi-typescript` でTypeScript型を自動生成
4. **フロントエンド開発**: 型安全なAPIクライアント (`openapi-fetch`) で開発
5. **実行時**: Swagger UIでAPIドキュメント確認可能 (`/swagger/index.html`)

### バックエンド構成 (Go - レイヤードアーキテクチャ)

```
backend/
├── cmd/server/main.go           # エントリポイント、DBマイグレーション
├── internal/
│   ├── config/                  # 環境設定
│   ├── database/                # DB接続、GORMセットアップ
│   ├── models/                  # GORMモデル (User, Post, Commentなど)
│   ├── handlers/                # HTTPハンドラー (リクエスト受信) + Swaggerアノテーション
│   ├── services/                # ビジネスロジック層
│   ├── middleware/              # JWT認証、CORS、エラーハンドリング
│   ├── utils/                   # JWT、パスワードハッシュ化、バリデータ、レスポンスヘルパー
│   └── routes/                  # ルート定義
├── docs/                        # Swagger生成ファイル (swagger.yaml, swagger.json)
├── migrations/                  # SQLマイグレーション (golang-migrate使用時)
├── .env                        # ローカル環境変数
└── docker-compose.yml          # PostgreSQL + Backendコンテナ
```

**重要なパターン**: Handlers → Services → Models (Repositoryパターン)
- **Handlers** は入力をバリデートしてサービスを呼び出す
- **Services** はビジネスロジックとデータベース操作を含む
- **Models** はリレーションを持つGORM構造体を定義

### フロントエンド構成 (React/TypeScript)

```
frontend/src/
├── api/
│   ├── client.ts               # openapi-fetchベースのAPIクライアント
│   └── schema.ts               # openapi-typescriptで自動生成された型定義
├── components/
│   ├── common/                 # 共通コンポーネント (AppBar, Sidebar, Layout)
│   ├── auth/                   # ログイン/登録フォーム
│   ├── post/                   # PostCard, PostForm, PostList
│   ├── comment/                # CommentList, CommentItem, CommentForm
│   └── user/                   # UserProfile, UserAvatar, FollowButton
├── pages/                      # ルートページ (HomePage, ProfilePageなど)
├── hooks/                      # React Queryフック (usePosts, useAuthなど)
├── context/                    # グローバル認証状態用AuthContext
├── utils/                      # localStorageヘルパー、日付フォーマット
└── theme/                      # MUIテーマ設定
```

**重要なパターン**: データ取得にReact Query + 認証にAuthContext + 型安全なAPIクライアント
- タイムライン/リストには **Infinite Queries** を使用（カーソルベースページネーション）
- いいね/フォローには楽観的更新を伴う **Mutations** を使用
- **型安全性**: openapi-fetchでコンパイル時に型チェック

---

## 開発コマンド

### バックエンド (Go)

**前提条件:**
- Go 1.21+
- Docker & Docker Compose

**初期セットアップ:**
```bash
cd backend
go mod init github.com/yourusername/sns-backend
go mod tidy
```

**Docker Composeで実行:**
```bash
# PostgreSQL + Backend を起動
docker-compose up -d

# ログを表示
docker-compose logs -f backend

# サービスを停止
docker-compose down  # -vフラグは使用しない (データ損失!)
```

**ローカルで実行 (Dockerなし):**
```bash
# PostgreSQLのみ起動
docker-compose up -d db

# バックエンドを実行
cd backend
go run cmd/server/main.go
```

**テスト:**
```bash
# すべてのテストを実行
go test ./...

# 特定のパッケージをテスト
go test ./internal/services

# カバレッジ付きで実行
go test -cover ./...
```

**データベースマイグレーション:**
```bash
# 自動マイグレーションはmain.goのGORM AutoMigrateで起動時に実行
# 手動マイグレーション (golang-migrate使用時):
migrate -path migrations -database "postgres://..." up
```

### フロントエンド (React)

**前提条件:**
- Node.js 18+
- npm または yarn

**初期セットアップ:**
```bash
npm create vite@latest frontend -- --template react-ts
cd frontend
npm install
```

**開発:**
```bash
cd frontend
npm run dev          # 開発サーバー起動 (http://localhost:5173)
npm run build        # 本番ビルド
npm run preview      # 本番ビルドをプレビュー
```

**テスト:**
```bash
npm run test         # テスト実行 (設定済みの場合)
```

**デプロイ:**
```bash
# Firebase Hosting
npm run build
firebase deploy
```

---

## 重要なアーキテクチャの決定事項

### 1. **論理削除（Soft Delete）**
- Users, Posts, Commentsは `deleted_at` カラムを使用
- **常に** `WHERE deleted_at IS NULL` でクエリ
- GORMは `gorm.DeletedAt` で自動的に処理

### 2. **JWT認証フロー**
1. ユーザーが登録/ログイン → バックエンドがJWTを返す
2. フロントエンドがJWTをlocalStorageに保存
3. すべての保護されたAPI呼び出しに `Authorization: Bearer <token>` を含める
4. バックエンドミドルウェアがJWTを検証してuser_idを抽出

### 3. **ページネーション戦略**
- 無限スクロール用の**カーソルベースページネーション**
- クエリパラメータ: `?limit=20&cursor=<last_post_id>`
- レスポンスに `has_more` と `next_cursor` を含む

### 4. **メディアアップロードフロー**
- Phase 1: ローカルの `/uploads` ディレクトリにアップロード（Dockerボリューム）
- Phase 2+: Firebase Storageにアップロード
- バックエンドが公開URLを返し、`media` テーブルに保存

### 5. **データベースリレーション**
```
Users ←→ Posts (1:N)
Posts ←→ Media (1:N, order_indexで順序付け)
Posts ←→ Comments (1:N)
Posts ←→ PostLikes (1:N)
Users ←→ Follows (follower_id/following_idによる自己参照M:N)
```

---

## 🐳 Dockerベース開発の厳格なルール

### **絶対ルール: ホストOSでgoコマンドを直接実行しない**

このプロジェクトは**Dockerベース**で開発します。ホストOSの環境に依存しないように、**すべてのgoコマンドはDockerコンテナ内で実行**してください。

### ✅ 正しいコマンド実行方法

```bash
# Goコマンドは必ずコンテナ内で実行
docker compose exec api go mod init github.com/yourusername/sns-backend
docker compose exec api go mod tidy
docker compose exec api go get <package>
docker compose exec api go run cmd/server/main.go
docker compose exec api go test ./...

# コンテナの起動・停止
docker compose up -d          # コンテナ起動（デタッチモード）
docker compose down           # コンテナ停止（データは保持）
docker compose restart api    # apiサービスのみ再起動
docker compose logs -f api    # apiサービスのログを表示
```

### ❌ 禁止事項

```bash
# ホストOSで直接goコマンドを実行してはいけない
go mod init                   # ❌ 禁止
go get ...                    # ❌ 禁止
go run ...                    # ❌ 禁止
```

### 📁 サービス構成

- **db**: PostgreSQLデータベース
- **api**: Goバックエンド（旧名: backend）
  - コンテナ名: `sns_api`
  - ポート: `8080`
  - ボリュームマウント: `./backend:/app`

---

## データ保護に関する重要なルール

### ⚠️ 絶対に使用してはいけないコマンド:
```bash
docker compose down -v        # ボリュームを削除 = データ損失
docker volume rm <volume>     # データベースデータを削除
```

### ✅ 安全なコマンド:
```bash
docker compose restart        # コンテナを再起動
docker compose down           # コンテナを停止（データは保持）
docker compose up -d --build  # 再ビルドして再起動
```

### データベース操作:
- **常に** 論理削除を使用 (`deleted_at`)
- ユーザーの明示的な確認なしに `DROP TABLE` や `TRUNCATE` を**実行しない**
- スキーマ変更前にバックアップを推奨

---

## API設計パターン

### エンドポイント構造
```
/api/v1/auth/*           # 認証 (register, login, me)
/api/v1/users/*          # ユーザープロフィール、フォロワー
/api/v1/posts/*          # 投稿CRUD、タイムライン
/api/v1/posts/:id/comments  # ネストされたリソース
/api/v1/posts/:id/like      # アクション (POST/DELETE)
```

### 標準レスポンス形式
```json
{
  "data": { /* レスポンスデータ */ },
  "message": "Success"
}
```

### エラーレスポンス
```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Error message"
  }
}
```

### ページネーションレスポンス
```json
{
  "data": [ /* 配列 */ ],
  "pagination": {
    "has_more": true,
    "next_cursor": "1234567890",
    "limit": 20
  }
}
```

---

## 命名規則

### データベース (PostgreSQL)
- テーブル: `snake_case` 複数形 (`users`, `posts`, `post_likes`)
- カラム: `snake_case` (`user_id`, `created_at`)

### バックエンド (Go)
- 構造体: `PascalCase` (`User`, `PostLike`)
- 関数: `camelCase` (`GetUserByID`, `CreatePost`)
- パッケージ: `lowercase` (`auth`, `models`)
- ファイル: `snake_case` (`auth_service.go`)

### フロントエンド (TypeScript)
- コンポーネント: `PascalCase` (`PostCard`, `UserProfile`)
- 関数/フック: `camelCase` (`fetchPosts`, `useAuth`)
- 定数: `UPPER_SNAKE_CASE` (`API_BASE_URL`)
- ファイル: コンポーネントは `PascalCase` (`PostCard.tsx`)、ユーティリティは `camelCase`

---

## 開発フェーズ

### **Phase 1 (MVP)** - 現在の焦点
コアSNS機能: 認証、投稿、コメント、いいね、フォロー、タイムライン
- ドキュメント: `docs/todo/03_PHASE1_BACKEND.md`, `04_PHASE1_FRONTEND.md`
- 目標: 機能的なTwitterライクSNS

### **Phase 2 (拡張機能)**
ハッシュタグ、複数画像、ブックマーク、パスワードリセット、メール認証
- ドキュメント: `docs/todo/05_PHASE2_BACKEND.md`, `06_PHASE2_FRONTEND.md`

### **Phase 3 (高度な機能)**
ユーザー検索、リツイート、通知、DM、トレンド投稿、ソーシャルログイン
- ドキュメント: `docs/todo/07_PHASE3_FUTURE.md`

---

## ドキュメント参照

- **プロジェクト概要**: `docs/todo/00_OVERVIEW.md`
- **データベーススキーマ**: `docs/todo/01_DATABASE_SCHEMA.md` (ER図、テーブル定義)
- **API仕様**: `docs/todo/02_API_SPECIFICATION.md` (全エンドポイント、リクエスト/レスポンス形式)
- **Phase別TODO**: `docs/todo/03-07_PHASE*.md` (ステップバイステップ実装ガイド)

---

## よく使うパターン & ヘルパー

### バックエンド: 新しいエンドポイントの作成

1. **モデル定義** (`internal/models/*.go`)
```go
type Post struct {
    gorm.Model
    UserID  uint   `gorm:"not null"`
    Content string `gorm:"type:text;not null"`
    User    User   `gorm:"foreignKey:UserID"`
}
```

2. **サービス作成** (`internal/services/*.go`)
```go
func CreatePost(userID uint, content string) (*Post, error) {
    post := &Post{UserID: userID, Content: content}
    result := db.Create(post)
    return post, result.Error
}
```

3. **ハンドラー作成** (`internal/handlers/*.go`)
```go
func CreatePost(c echo.Context) error {
    var req PostCreateRequest
    if err := c.Bind(&req); err != nil {
        return utils.ErrorResponse(c, 400, "Invalid request")
    }
    userID := c.Get("user_id").(uint)
    post, err := services.CreatePost(userID, req.Content)
    return utils.SuccessResponse(c, post)
}
```

4. **ルート登録** (`internal/routes/routes.go`)
```go
api.POST("/posts", handlers.CreatePost, middleware.JWTAuth())
```

### フロントエンド: React Queryでデータ取得

```typescript
// hooks/usePosts.ts 内
export const usePosts = () => {
  return useInfiniteQuery({
    queryKey: ['posts', 'timeline'],
    queryFn: ({ pageParam = null }) => getPosts(pageParam),
    getNextPageParam: (lastPage) => lastPage.pagination.next_cursor,
  });
};

// コンポーネント内
const { data, fetchNextPage, hasNextPage } = usePosts();
```

---

## 環境変数

### バックエンド (.env)
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=sns_db
JWT_SECRET=your-secret-key-change-in-production
PORT=8080
ENV=development
FIREBASE_CREDENTIALS_PATH=/path/to/serviceAccountKey.json  # Phase 2+
```

### フロントエンド (.env.local)
```
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

---

## トラブルシューティング

### バックエンドの問題
- **DB接続失敗**: `docker-compose ps` を確認、PostgreSQLが起動しているか確認
- **JWTエラー**: `.env` とトークン生成間で `JWT_SECRET` が一致しているか確認
- **CORSエラー**: `internal/middleware/cors_middleware.go` の設定を確認

### フロントエンドの問題
- **API呼び出し失敗**: `VITE_API_BASE_URL` が起動中のバックエンドを指しているか確認
- **認証リダイレクト**: `AuthContext` とlocalStorageの有効なトークンを確認

### Dockerの問題
- **ポート競合**: ポート5432 (PostgreSQL) または8080 (バックエンド) が既に使用中か確認
- **データ損失**: `docker-compose down -v` を**絶対に使用しない** (データ保護ルール参照)

---

## テスト戦略

### バックエンド
- Goの組み込み `testing` パッケージを使用
- モックDBでサービスを独立してテスト
- `httptest` を使用したハンドラーのHTTPテスト

### フロントエンド
- React Testing Libraryでコンポーネントテスト
- Playwrightを使用したE2Eテスト (Phase 2+)

---

## セキュリティ考慮事項

- **パスワード保存**: 常にbcryptを使用 (`User.BeforeCreate` フックで処理)
- **JWTシークレット**: 本番環境では強力でランダムなシークレットを使用
- **入力バリデーション**: リクエスト構造体に `validator/v10` タグを使用
- **SQLインジェクション**: GORMがプリペアドステートメントで自動的に防止
- **CORS**: 本番環境では許可されたオリジンを適切に設定
- **レート制限**: Phase 2+で実装 (`02_API_SPECIFICATION.md` 参照)

---

## メディアファイル制限

| 種類  | 最大サイズ | 最大長さ | 対応フォーマット        |
|-------|----------|---------|----------------------|
| 画像  | 5 MB     | -       | jpg, png, gif, heic  |
| 動画  | 50 MB    | 30秒    | mp4, mov             |
| 音声  | -        | -       | mp3                  |

---

## 新規開発者向けクイックスタート

1. **読む**: プロジェクトビジョンについて `docs/todo/00_OVERVIEW.md` を読む
2. **セットアップ**: `03_PHASE1_BACKEND.md` の項目1-4（プロジェクトセットアップ）に従う
3. **実行**: プロジェクトルートで `docker-compose up -d`
4. **開発**: 1つの機能モジュール（例：認証）から始める
5. **参照**: エンドポイント詳細は `02_API_SPECIFICATION.md` を使用

---

**最終更新**: 2026-02-14
