# 管理画面 要件定義書

## 1. プロジェクト概要

### 1.1 目的
SNSアプリケーションの運用管理を効率化するための管理画面を構築する。

### 1.2 対象ユーザー
- システム管理者（adminロールを持つユーザー）

### 1.3 技術スタック
- **バックエンド**: Go + Echo + html/template
- **フロントエンド**: Bulma CSS + Chart.js
- **認証**: Basic認証 + JWT（adminロール）
- **データベース**: PostgreSQL

---

## 2. 機能要件

### 2.1 認証・認可

#### 2.1.1 管理者アカウント
- **デフォルト管理者アカウント**
  - ID: `udemy_sns_admin`
  - PASSWORD: `password`
  - Role: `admin`
  - Status: `approved`

#### 2.1.2 認証フロー
1. Basic認証（`/admin/*` パス全体に適用）
   - ID: `admin`
   - PASSWORD: 環境変数 `ADMIN_BASIC_PASSWORD` から取得
2. JWT認証（adminロールチェック）
   - ログイン画面でJWT発行
   - すべての管理APIでadminロールを検証

#### 2.1.3 ミドルウェア
- **Basic認証ミドルウェア**: `/admin/*` に適用
- **Adminロールチェックミドルウェア**: すべての管理API（`/admin/api/*`）に適用
- **レートリミット**: 管理画面APIにも適用（60req/min）

---

### 2.2 画面構成

#### 2.2.1 管理者ログイン画面
- **パス**: `/admin/login`
- **機能**:
  - ID/パスワード入力フォーム
  - JWT発行（HttpOnly Cookie）
  - ログイン失敗時のエラー表示
- **アクセス制限**: Basic認証のみ（JWT不要）

#### 2.2.2 ダッシュボード
- **パス**: `/admin/dashboard`
- **表示内容**:
  - **ユーザー統計**（数値カード）
    - 総ユーザー数
    - 承認待ちユーザー数（pending）
    - 承認済みユーザー数（approved）
    - 拒否済みユーザー数（rejected）
  - **投稿統計**（数値カード）
    - 総投稿数
    - 本日の投稿数
    - メディア付き投稿率
  - **グラフ**
    - 投稿数推移（直近30日、折れ線グラフ）
    - 新規登録ユーザー推移（直近30日、折れ線グラフ）
  - **アクティブ統計**
    - アクティブユーザー数（直近7日で投稿したユーザー）
    - パスワードリセット申請件数（pending）
  - **アラート**（該当する場合のみ表示）
    - ⚠️ 承認待ちユーザーが10人以上
    - ⚠️ パスワードリセット申請が5件以上
    - ⚠️ 直近1時間で異常なレートリミット超過（50回以上）

#### 2.2.3 ユーザー一覧
- **パス**: `/admin/users`
- **機能**:
  - **フィルタ**
    - status: all / pending / approved / rejected
    - role: all / user / admin
  - **ソート**
    - 登録日時（降順/昇順）
    - 最終ログイン日時（降順/昇順）
  - **検索**
    - username, email で部分一致検索
  - **ページネーション**
    - 20件/ページ
  - **表示項目**
    - ID, username, email, role, status, 登録日時, 投稿数
  - **操作**
    - 詳細表示ボタン
    - 一括操作チェックボックス（承認/拒否）

#### 2.2.4 ユーザー詳細
- **パス**: `/admin/users/:id`
- **表示内容**:
  - **基本情報**
    - ID, username, email, bio
    - role, status
    - 登録日時, 最終ログイン日時
  - **統計**
    - 投稿数, いいね数（受信）
    - フォロワー数, フォロー数
  - **最近の投稿**（5件）
    - 投稿内容, 投稿日時, いいね数, コメント数
  - **操作ボタン**
    - ステータス変更（承認/拒否/保留に戻す）
    - ユーザーの投稿一覧を見る（別タブで `/admin/users/:id/posts`）

#### 2.2.5 パスワードリセット申請一覧
- **パス**: `/admin/password-resets`
- **機能**:
  - **フィルタ**
    - status: all / pending / approved / expired / used
  - **表示項目**
    - 申請ID, ユーザー名, メールアドレス, 申請日時, ステータス
  - **操作**
    - リセットリンク発行（pendingの場合のみ）
      - トークン生成（UUID）
      - 有効期限: 24時間
      - リセットURL: `https://yourdomain.com/reset-password?token=xxx`
      - クリップボードにコピーボタン
    - メールテンプレートコピーボタン
      - ユーザー名、リセットリンクが埋め込まれたテンプレート

#### 2.2.6 操作ログ一覧
- **パス**: `/admin/logs`
- **機能**:
  - **フィルタ**
    - 操作種別: all / approve_user / reject_user / password_reset_approve
    - 管理者: 管理者username でフィルタ
    - 日付範囲: 開始日〜終了日
  - **表示項目**
    - 日時, 管理者, 操作種別, 対象ユーザー, 詳細
  - **ページネーション**: 50件/ページ

---

### 2.3 API仕様

#### 2.3.1 認証API

##### POST /admin/api/login
管理者ログイン

**Request:**
```json
{
  "username": "udemy_sns_admin",
  "password": "password"
}
```

**Response (200):**
```json
{
  "data": {
    "user": {
      "id": 1,
      "username": "udemy_sns_admin",
      "role": "admin"
    }
  },
  "message": "Login successful"
}
```
- Set-Cookie: `admin_token=xxx; HttpOnly; Secure; SameSite=None; Path=/admin`

**Error (401):**
```json
{
  "error": {
    "code": "INVALID_CREDENTIALS",
    "message": "Invalid username or password"
  }
}
```

##### POST /admin/api/logout
管理者ログアウト

**Response (200):**
```json
{
  "message": "Logout successful"
}
```

---

#### 2.3.2 ダッシュボードAPI

##### GET /admin/api/dashboard/stats
ダッシュボード統計データ取得

**Response (200):**
```json
{
  "data": {
    "users": {
      "total": 1523,
      "pending": 12,
      "approved": 1500,
      "rejected": 11
    },
    "posts": {
      "total": 45230,
      "today": 342,
      "with_media_rate": 68.5
    },
    "active_users_7d": 856,
    "password_reset_pending": 3,
    "alerts": [
      {
        "type": "pending_users",
        "message": "承認待ちユーザーが10人以上います",
        "count": 12
      }
    ]
  }
}
```

##### GET /admin/api/dashboard/charts/posts
投稿数推移データ（直近30日）

**Response (200):**
```json
{
  "data": {
    "labels": ["2025-01-20", "2025-01-21", ..., "2025-02-18"],
    "values": [234, 256, 289, ..., 342]
  }
}
```

##### GET /admin/api/dashboard/charts/users
新規登録ユーザー推移データ（直近30日）

**Response (200):**
```json
{
  "data": {
    "labels": ["2025-01-20", "2025-01-21", ..., "2025-02-18"],
    "values": [12, 15, 8, ..., 23]
  }
}
```

---

#### 2.3.3 ユーザー管理API

##### GET /admin/api/users
ユーザー一覧取得

**Query Parameters:**
- `status`: all | pending | approved | rejected (default: all)
- `role`: all | user | admin (default: all)
- `sort`: created_at | last_login (default: created_at)
- `order`: asc | desc (default: desc)
- `search`: 検索キーワード（username, email）
- `page`: ページ番号（default: 1）
- `limit`: 1ページあたりの件数（default: 20）

**Response (200):**
```json
{
  "data": {
    "users": [
      {
        "id": 123,
        "username": "john_doe",
        "email": "john@example.com",
        "role": "user",
        "status": "pending",
        "created_at": "2025-02-15T10:30:00Z",
        "post_count": 5
      }
    ],
    "pagination": {
      "total": 1523,
      "page": 1,
      "limit": 20,
      "total_pages": 77
    }
  }
}
```

##### GET /admin/api/users/:id
ユーザー詳細取得

**Response (200):**
```json
{
  "data": {
    "user": {
      "id": 123,
      "username": "john_doe",
      "email": "john@example.com",
      "bio": "Software Engineer",
      "role": "user",
      "status": "pending",
      "created_at": "2025-02-15T10:30:00Z",
      "updated_at": "2025-02-15T10:30:00Z"
    },
    "stats": {
      "post_count": 5,
      "like_count": 23,
      "follower_count": 12,
      "following_count": 34
    },
    "recent_posts": [
      {
        "id": 456,
        "content": "Hello world!",
        "created_at": "2025-02-16T14:20:00Z",
        "like_count": 5,
        "comment_count": 2
      }
    ]
  }
}
```

##### PATCH /admin/api/users/:id/status
ユーザーステータス変更

**Request:**
```json
{
  "status": "approved"
}
```

**Response (200):**
```json
{
  "data": {
    "user": {
      "id": 123,
      "username": "john_doe",
      "status": "approved"
    }
  },
  "message": "User status updated successfully"
}
```

**構造化ログ出力:**
```json
{
  "level": "info",
  "time": "2025-02-18T12:00:00Z",
  "action": "approve_user",
  "admin_id": 1,
  "admin_username": "udemy_sns_admin",
  "target_user_id": 123,
  "target_username": "john_doe",
  "old_status": "pending",
  "new_status": "approved",
  "ip": "192.168.1.1"
}
```

##### POST /admin/api/users/batch-update-status
ユーザー一括ステータス変更

**Request:**
```json
{
  "user_ids": [123, 456, 789],
  "status": "approved"
}
```

**Response (200):**
```json
{
  "data": {
    "updated_count": 3
  },
  "message": "3 users updated successfully"
}
```

---

#### 2.3.4 パスワードリセットAPI

##### GET /admin/api/password-resets
パスワードリセット申請一覧取得

**Query Parameters:**
- `status`: all | pending | approved | expired | used (default: all)
- `page`: ページ番号（default: 1）
- `limit`: 1ページあたりの件数（default: 20）

**Response (200):**
```json
{
  "data": {
    "requests": [
      {
        "id": 1,
        "user": {
          "id": 123,
          "username": "john_doe",
          "email": "john@example.com"
        },
        "status": "pending",
        "created_at": "2025-02-17T09:00:00Z",
        "expires_at": null
      }
    ],
    "pagination": {
      "total": 15,
      "page": 1,
      "limit": 20,
      "total_pages": 1
    }
  }
}
```

##### POST /admin/api/password-resets/:id/approve
パスワードリセット承認（トークン発行）

**Response (200):**
```json
{
  "data": {
    "token": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "reset_url": "https://yourdomain.com/reset-password?token=a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "expires_at": "2025-02-19T12:00:00Z",
    "email_template": "こんにちは、john_doeさん\n\nパスワードリセットのリクエストを承認しました。\n以下のリンクからパスワードを再設定してください（24時間有効）：\n\nhttps://yourdomain.com/reset-password?token=a1b2c3d4-e5f6-7890-abcd-ef1234567890\n\nこのリクエストに心当たりがない場合は、このメールを無視してください。\n\nよろしくお願いいたします。"
  },
  "message": "Password reset approved"
}
```

**構造化ログ出力:**
```json
{
  "level": "info",
  "time": "2025-02-18T12:00:00Z",
  "action": "password_reset_approve",
  "admin_id": 1,
  "admin_username": "udemy_sns_admin",
  "target_user_id": 123,
  "target_username": "john_doe",
  "reset_request_id": 1,
  "token": "a1b2c3d4-****-****",
  "expires_at": "2025-02-19T12:00:00Z",
  "ip": "192.168.1.1"
}
```

---

#### 2.3.5 操作ログAPI

##### GET /admin/api/logs
操作ログ一覧取得

**Query Parameters:**
- `action`: all | approve_user | reject_user | password_reset_approve (default: all)
- `admin_username`: 管理者username でフィルタ
- `start_date`: 開始日（YYYY-MM-DD）
- `end_date`: 終了日（YYYY-MM-DD）
- `page`: ページ番号（default: 1）
- `limit`: 1ページあたりの件数（default: 50）

**Response (200):**
```json
{
  "data": {
    "logs": [
      {
        "id": 1,
        "time": "2025-02-18T12:00:00Z",
        "action": "approve_user",
        "admin": {
          "id": 1,
          "username": "udemy_sns_admin"
        },
        "target_user": {
          "id": 123,
          "username": "john_doe"
        },
        "details": "Status changed from pending to approved",
        "ip": "192.168.1.1"
      }
    ],
    "pagination": {
      "total": 234,
      "page": 1,
      "limit": 50,
      "total_pages": 5
    }
  }
}
```

---

## 3. データベース設計

### 3.1 既存テーブルの変更

#### usersテーブル
```sql
ALTER TABLE users ADD COLUMN role VARCHAR(20) NOT NULL DEFAULT 'user';
ALTER TABLE users ADD COLUMN status VARCHAR(20) NOT NULL DEFAULT 'pending';
ALTER TABLE users ADD COLUMN last_login_at TIMESTAMP NULL;

-- role: 'user' | 'admin'
-- status: 'pending' | 'approved' | 'rejected'
```

#### GORM Model
```go
type User struct {
    gorm.Model
    Username      string         `gorm:"uniqueIndex;not null"`
    Email         string         `gorm:"uniqueIndex;not null"`
    Password      string         `gorm:"not null"`
    Bio           string         `gorm:"type:text"`
    Role          string         `gorm:"type:varchar(20);default:'user';not null"`
    Status        string         `gorm:"type:varchar(20);default:'pending';not null"`
    LastLoginAt   *time.Time
    DeletedAt     gorm.DeletedAt `gorm:"index"`
}
```

---

### 3.2 新規テーブル

#### password_reset_requestsテーブル
```sql
CREATE TABLE password_reset_requests (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) UNIQUE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    admin_approved_by INTEGER NULL REFERENCES users(id) ON DELETE SET NULL,
    admin_approved_at TIMESTAMP NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- status: 'pending' | 'approved' | 'expired' | 'used'

CREATE INDEX idx_password_reset_requests_user_id ON password_reset_requests(user_id);
CREATE INDEX idx_password_reset_requests_token ON password_reset_requests(token);
CREATE INDEX idx_password_reset_requests_status ON password_reset_requests(status);
```

#### GORM Model
```go
type PasswordResetRequest struct {
    ID                uint      `gorm:"primaryKey"`
    UserID            uint      `gorm:"not null;index"`
    Token             string    `gorm:"type:varchar(255);uniqueIndex;not null"`
    Status            string    `gorm:"type:varchar(20);default:'pending';not null;index"`
    AdminApprovedBy   *uint
    AdminApprovedByUser *User  `gorm:"foreignKey:AdminApprovedBy"`
    AdminApprovedAt   *time.Time
    ExpiresAt         time.Time `gorm:"not null"`
    CreatedAt         time.Time
    UpdatedAt         time.Time
    User              User      `gorm:"foreignKey:UserID"`
}
```

---

#### admin_logsテーブル
```sql
CREATE TABLE admin_logs (
    id SERIAL PRIMARY KEY,
    admin_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    action VARCHAR(50) NOT NULL,
    target_user_id INTEGER NULL REFERENCES users(id) ON DELETE SET NULL,
    details TEXT,
    ip VARCHAR(50),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- action: 'approve_user' | 'reject_user' | 'password_reset_approve' | 'user_status_change'

CREATE INDEX idx_admin_logs_admin_id ON admin_logs(admin_id);
CREATE INDEX idx_admin_logs_action ON admin_logs(action);
CREATE INDEX idx_admin_logs_created_at ON admin_logs(created_at);
```

#### GORM Model
```go
type AdminLog struct {
    ID           uint      `gorm:"primaryKey"`
    AdminID      uint      `gorm:"not null;index"`
    Admin        User      `gorm:"foreignKey:AdminID"`
    Action       string    `gorm:"type:varchar(50);not null;index"`
    TargetUserID *uint
    TargetUser   *User     `gorm:"foreignKey:TargetUserID"`
    Details      string    `gorm:"type:text"`
    IP           string    `gorm:"type:varchar(50)"`
    CreatedAt    time.Time `gorm:"index"`
}
```

---

## 4. セキュリティ要件

### 4.1 認証・認可
- ✅ Basic認証（`/admin/*` 全体）
- ✅ JWT認証（adminロールチェック）
- ✅ 既存の承認APIに権限チェック追加（`PATCH /api/v1/users/:id/approve`）
  - 現在は認証済みユーザー誰でも承認可能 → adminロールのみに制限

### 4.2 ロギング
- ✅ すべての管理操作を構造化ログに記録
  - 管理者ID, 管理者username
  - 操作種別（action）
  - 対象ユーザーID, username
  - 変更内容（old_status → new_status）
  - IPアドレス
  - タイムスタンプ

### 4.3 レートリミット
- ✅ 管理画面APIにもレートリミット適用（60req/min）
- ログインAPI: 5req/min（ブルートフォース対策）

### 4.4 パスワードリセットトークン
- トークン: UUID v4
- 有効期限: 24時間
- 1回限り使用可能（使用後は status='used' に変更）
- 有効期限切れチェック（cronジョブで定期的に status='expired' に更新）

---

## 5. 非機能要件

### 5.1 パフォーマンス
- ダッシュボード統計データ: 1秒以内
- ユーザー一覧: 1秒以内（ページネーション適用）
- N+1クエリの排除（Preload/Joins使用）

### 5.2 スケーラビリティ
- ログテーブルのパーティショニング（将来的に月別）
- 統計データのキャッシュ（Redis、Phase 3で実装）

### 5.3 保守性
- すべてのAPIに構造化ログ
- エラーハンドリングの統一
- コードコメントの充実

---

## 6. UI/UX要件

### 6.1 デザイン
- **CSSフレームワーク**: Bulma
- **カラースキーム**: Bulmaのデフォルトテーマ（プライマリカラー: #00d1b2）
- **レスポンシブ**: タブレット以上を想定（モバイル対応は不要）

### 6.2 ユーザビリティ
- **パンくずリスト**: すべてのページに表示
- **フィードバック**: 操作成功時は緑の通知、エラー時は赤の通知（Bulma notification）
- **確認ダイアログ**: 削除・拒否操作時に確認ダイアログを表示
- **ローディング**: データ取得中はスピナーを表示

### 6.3 Chart.js設定
- **折れ線グラフ**: 投稿数推移、ユーザー登録推移
- **オプション**:
  - レスポンシブ: true
  - アニメーション: 有効
  - ツールチップ: 有効
  - 凡例: 表示

---

## 7. 実装フェーズ

### Phase 1: 基盤（2日）
- ✅ DB migration（role, status, last_login_at カラム追加）
- ✅ 新規テーブル作成（password_reset_requests, admin_logs）
- ✅ 管理者アカウント作成（seed）
- ✅ Basic認証ミドルウェア
- ✅ Adminロールチェックミドルウェア
- ✅ 構造化ログ（admin操作用）
- ✅ 既存承認APIに権限チェック追加

### Phase 2: 基本機能（3日）
- ✅ 管理者ログイン画面（html/template）
- ✅ ユーザー一覧（フィルタ、検索、ページネーション）
- ✅ ユーザー詳細
- ✅ ユーザー承認/拒否API + UI
- ✅ レイアウトテンプレート（ヘッダー、サイドバー）

### Phase 3: パスワードリセット（2日）
- ✅ パスワードリセット申請API（フロントエンド側で実装）
- ✅ 管理画面: リセット申請一覧
- ✅ トークン発行機能
- ✅ メールテンプレート生成
- ✅ クリップボードコピー機能

### Phase 4: ダッシュボード（3日）
- ✅ 統計データAPI（ユーザー、投稿、アクティブユーザー）
- ✅ グラフデータAPI（投稿数推移、ユーザー登録推移）
- ✅ Chart.jsでグラフ描画
- ✅ アラート機能

### Phase 5: 操作ログ（1日）
- ✅ 操作ログAPI
- ✅ ログ閲覧UI（フィルタ、ページネーション）

### Phase 6: テスト・デプロイ（2日）
- ✅ 手動テスト
- ✅ セキュリティチェック
- ✅ 本番環境変数設定
- ✅ デプロイ

**合計: 13日**

---

## 8. 環境変数

### backend/.env
```
# 既存の環境変数
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=sns_db
JWT_SECRET=your-secret-key
PORT=8080
ENV=development

# 新規追加: 管理画面用
ADMIN_BASIC_USER=admin
ADMIN_BASIC_PASSWORD=your-basic-auth-password-change-in-production
FRONTEND_URL=http://localhost:5173
ADMIN_JWT_SECRET=your-admin-jwt-secret-key-change-in-production
```

---

## 9. ディレクトリ構成

```
backend/
├── cmd/server/main.go
├── internal/
│   ├── admin/                          # 新規: 管理画面専用
│   │   ├── handlers/
│   │   │   ├── auth_handler.go        # 管理者ログイン/ログアウト
│   │   │   ├── dashboard_handler.go   # ダッシュボード統計
│   │   │   ├── user_handler.go        # ユーザー管理
│   │   │   ├── password_reset_handler.go  # パスワードリセット
│   │   │   └── log_handler.go         # 操作ログ
│   │   ├── services/
│   │   │   ├── auth_service.go
│   │   │   ├── dashboard_service.go
│   │   │   ├── user_service.go
│   │   │   ├── password_reset_service.go
│   │   │   └── log_service.go
│   │   ├── middleware/
│   │   │   ├── basic_auth.go          # Basic認証
│   │   │   └── admin_auth.go          # Adminロールチェック
│   │   └── templates/                 # html/template
│   │       ├── layout.html            # 共通レイアウト
│   │       ├── login.html             # ログイン画面
│   │       ├── dashboard.html         # ダッシュボード
│   │       ├── users/
│   │       │   ├── index.html         # ユーザー一覧
│   │       │   └── detail.html        # ユーザー詳細
│   │       ├── password_resets/
│   │       │   └── index.html         # パスワードリセット一覧
│   │       └── logs/
│   │           └── index.html         # 操作ログ一覧
│   ├── models/
│   │   ├── user.go                    # 既存: role, status カラム追加
│   │   ├── password_reset_request.go  # 新規
│   │   └── admin_log.go               # 新規
│   ├── routes/
│   │   └── admin_routes.go            # 新規: 管理画面ルート
│   └── ...
├── static/                             # 新規: 静的ファイル
│   ├── css/
│   │   └── admin.css                  # 管理画面専用CSS
│   └── js/
│       ├── chart.min.js               # Chart.js
│       └── admin.js                   # 管理画面専用JS
└── migrations/
    └── XXXXXX_add_admin_features.sql  # 新規マイグレーション
```

---

## 10. メールテンプレート

### 10.1 ユーザー承認通知
```
件名: アカウントが承認されました

こんにちは、{{.Username}}さん

あなたのアカウントが承認されました。
以下のURLからログインしてご利用ください：

{{.FrontendURL}}/login

ご不明な点がございましたら、お気軽にお問い合わせください。

よろしくお願いいたします。
SNS運営チーム
```

### 10.2 ユーザー拒否通知
```
件名: アカウント申請について

こんにちは、{{.Username}}さん

誠に申し訳ございませんが、あなたのアカウント申請は承認されませんでした。

理由: {{.Reason}}

ご不明な点がございましたら、お問い合わせください。

SNS運営チーム
```

### 10.3 パスワードリセット承認通知
```
件名: パスワードリセットのご案内

こんにちは、{{.Username}}さん

パスワードリセットのリクエストを承認しました。
以下のリンクからパスワードを再設定してください（24時間有効）：

{{.ResetURL}}

このリクエストに心当たりがない場合は、このメールを無視してください。

よろしくお願いいたします。
SNS運営チーム
```

---

## 11. API Rate Limit

| エンドポイント | リミット | 期間 |
|--------------|---------|------|
| POST /admin/api/login | 5回 | 1分 |
| その他の管理API | 60回 | 1分 |

---

## 12. エラーコード

| コード | 説明 |
|--------|------|
| INVALID_CREDENTIALS | 認証情報が無効 |
| UNAUTHORIZED | 認証されていない |
| FORBIDDEN | 権限がない（adminロールでない） |
| NOT_FOUND | リソースが見つからない |
| VALIDATION_ERROR | バリデーションエラー |
| INTERNAL_ERROR | サーバー内部エラー |

---

## 13. 今後の拡張（Phase 3以降）

### 13.1 高度な分析
- いいね数推移グラフ
- コメント数推移グラフ
- フォロー関係増加数グラフ
- エンゲージメント率（いいね数/投稿数）
- ユーザー投稿率（登録後1回以上投稿したユーザーの割合）

### 13.2 通知機能
- 承認待ちユーザーが10人以上になったらメール通知
- パスワードリセット申請が5件以上になったらメール通知

### 13.3 一括操作
- CSVエクスポート（ユーザー一覧、投稿一覧）
- CSVインポート（ユーザー一括登録）

### 13.4 監視ダッシュボード
- Redis統計（キャッシュヒット率）
- API呼び出し統計（エンドポイント別）
- エラーログ一覧

---

以上が管理画面の要件定義です。
