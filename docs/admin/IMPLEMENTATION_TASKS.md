# 管理画面 実装タスクリスト

このドキュメントは、管理画面の実装を段階的に進めるための詳細なタスクリストです。

---

## Phase 1: 基盤（推定2日）

### 1.1 データベースマイグレーション

#### タスク 1.1.1: usersテーブルにカラム追加
- [ ] マイグレーションファイル作成
  - `role` VARCHAR(20) NOT NULL DEFAULT 'user'
  - `status` VARCHAR(20) NOT NULL DEFAULT 'pending'
  - `last_login_at` TIMESTAMP NULL
- [ ] User構造体にフィールド追加
- [ ] GORMのAutoMigrateで適用

**ファイル:**
- `backend/internal/models/user.go`

**実装例:**
```go
type User struct {
    gorm.Model
    Username      string         `gorm:"uniqueIndex;not null" json:"username"`
    Email         string         `gorm:"uniqueIndex;not null" json:"email"`
    Password      string         `gorm:"not null" json:"-"`
    Bio           string         `gorm:"type:text" json:"bio"`
    Role          string         `gorm:"type:varchar(20);default:'user';not null" json:"role"`
    Status        string         `gorm:"type:varchar(20);default:'pending';not null" json:"status"`
    LastLoginAt   *time.Time     `json:"last_login_at,omitempty"`
    DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
```

---

#### タスク 1.1.2: password_reset_requestsテーブル作成
- [ ] PasswordResetRequest構造体作成
- [ ] GORMのAutoMigrateで適用

**ファイル:**
- `backend/internal/models/password_reset_request.go`

**実装例:**
```go
package models

import "time"

type PasswordResetRequest struct {
    ID                  uint       `gorm:"primaryKey" json:"id"`
    UserID              uint       `gorm:"not null;index" json:"user_id"`
    Token               string     `gorm:"type:varchar(255);uniqueIndex;not null" json:"token"`
    Status              string     `gorm:"type:varchar(20);default:'pending';not null;index" json:"status"`
    AdminApprovedBy     *uint      `json:"admin_approved_by,omitempty"`
    AdminApprovedByUser *User      `gorm:"foreignKey:AdminApprovedBy" json:"admin_approved_by_user,omitempty"`
    AdminApprovedAt     *time.Time `json:"admin_approved_at,omitempty"`
    ExpiresAt           time.Time  `gorm:"not null" json:"expires_at"`
    CreatedAt           time.Time  `json:"created_at"`
    UpdatedAt           time.Time  `json:"updated_at"`
    User                User       `gorm:"foreignKey:UserID" json:"user"`
}
```

---

#### タスク 1.1.3: admin_logsテーブル作成
- [ ] AdminLog構造体作成
- [ ] GORMのAutoMigrateで適用

**ファイル:**
- `backend/internal/models/admin_log.go`

**実装例:**
```go
package models

import "time"

type AdminLog struct {
    ID           uint      `gorm:"primaryKey" json:"id"`
    AdminID      uint      `gorm:"not null;index" json:"admin_id"`
    Admin        User      `gorm:"foreignKey:AdminID" json:"admin"`
    Action       string    `gorm:"type:varchar(50);not null;index" json:"action"`
    TargetUserID *uint     `json:"target_user_id,omitempty"`
    TargetUser   *User     `gorm:"foreignKey:TargetUserID" json:"target_user,omitempty"`
    Details      string    `gorm:"type:text" json:"details"`
    IP           string    `gorm:"type:varchar(50)" json:"ip"`
    CreatedAt    time.Time `gorm:"index" json:"created_at"`
}
```

---

#### タスク 1.1.4: 管理者アカウント作成（seed）
- [ ] シードデータスクリプト作成
- [ ] 管理者アカウント作成（username: udemy_sns_admin, password: password, role: admin, status: approved）

**ファイル:**
- `backend/cmd/server/main.go` または `backend/internal/database/seed.go`

**実装例:**
```go
func seedAdminUser(db *gorm.DB) error {
    var count int64
    db.Model(&models.User{}).Where("username = ?", "udemy_sns_admin").Count(&count)
    if count > 0 {
        log.Println("Admin user already exists")
        return nil
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    admin := models.User{
        Username: "udemy_sns_admin",
        Email:    "admin@example.com",
        Password: string(hashedPassword),
        Role:     "admin",
        Status:   "approved",
    }

    return db.Create(&admin).Error
}
```

---

### 1.2 ミドルウェア実装

#### タスク 1.2.1: Basic認証ミドルウェア
- [ ] Basic認証ミドルウェア作成
- [ ] 環境変数からID/パスワードを取得
- [ ] 認証失敗時は401エラーを返す

**ファイル:**
- `backend/internal/admin/middleware/basic_auth.go`

**実装例:**
```go
package middleware

import (
    "encoding/base64"
    "net/http"
    "os"
    "strings"

    "github.com/labstack/echo/v4"
)

func BasicAuth() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            auth := c.Request().Header.Get("Authorization")
            if auth == "" {
                c.Response().Header().Set("WWW-Authenticate", `Basic realm="Admin Area"`)
                return echo.NewHTTPError(http.StatusUnauthorized, "Missing authorization header")
            }

            if !strings.HasPrefix(auth, "Basic ") {
                return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authorization format")
            }

            payload, err := base64.StdEncoding.DecodeString(auth[6:])
            if err != nil {
                return echo.NewHTTPError(http.StatusUnauthorized, "Invalid base64 encoding")
            }

            pair := strings.SplitN(string(payload), ":", 2)
            if len(pair) != 2 {
                return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials format")
            }

            username := pair[0]
            password := pair[1]

            expectedUser := os.Getenv("ADMIN_BASIC_USER")
            expectedPass := os.Getenv("ADMIN_BASIC_PASSWORD")

            if username != expectedUser || password != expectedPass {
                c.Response().Header().Set("WWW-Authenticate", `Basic realm="Admin Area"`)
                return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
            }

            return next(c)
        }
    }
}
```

---

#### タスク 1.2.2: Adminロールチェックミドルウェア
- [ ] JWT認証ミドルウェア（既存のJWT認証を流用）
- [ ] ユーザーのroleをチェック
- [ ] adminでない場合は403エラーを返す

**ファイル:**
- `backend/internal/admin/middleware/admin_auth.go`

**実装例:**
```go
package middleware

import (
    "net/http"

    "github.com/labstack/echo/v4"
    "your-project/internal/models"
)

func AdminRoleCheck() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            userID := c.Get("user_id")
            if userID == nil {
                return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
            }

            var user models.User
            if err := db.First(&user, userID).Error; err != nil {
                return echo.NewHTTPError(http.StatusUnauthorized, "User not found")
            }

            if user.Role != "admin" {
                return echo.NewHTTPError(http.StatusForbidden, "Admin access required")
            }

            // 管理者情報をコンテキストに保存
            c.Set("admin_user", user)
            return next(c)
        }
    }
}
```

---

#### タスク 1.2.3: 構造化ログヘルパー
- [ ] 管理操作ログを記録するヘルパー関数作成
- [ ] AdminLogテーブルにレコード挿入
- [ ] JSON形式でログ出力

**ファイル:**
- `backend/internal/admin/utils/logger.go`

**実装例:**
```go
package utils

import (
    "encoding/json"
    "log"

    "gorm.io/gorm"
    "your-project/internal/models"
)

type AdminLogParams struct {
    AdminID      uint
    AdminUsername string
    Action       string
    TargetUserID *uint
    TargetUsername *string
    Details      string
    IP           string
}

func LogAdminAction(db *gorm.DB, params AdminLogParams) error {
    // データベースに記録
    adminLog := models.AdminLog{
        AdminID:      params.AdminID,
        Action:       params.Action,
        TargetUserID: params.TargetUserID,
        Details:      params.Details,
        IP:           params.IP,
    }

    if err := db.Create(&adminLog).Error; err != nil {
        return err
    }

    // 構造化ログ出力
    logData := map[string]interface{}{
        "level":           "info",
        "action":          params.Action,
        "admin_id":        params.AdminID,
        "admin_username":  params.AdminUsername,
        "target_user_id":  params.TargetUserID,
        "target_username": params.TargetUsername,
        "details":         params.Details,
        "ip":              params.IP,
    }

    jsonLog, _ := json.Marshal(logData)
    log.Println(string(jsonLog))

    return nil
}
```

---

#### タスク 1.2.4: 既存承認APIに権限チェック追加
- [ ] `PATCH /api/v1/users/:id/approve` にAdminロールチェック追加
- [ ] 既存のハンドラーを修正

**ファイル:**
- `backend/internal/handlers/user_handler.go`（既存のハンドラー）
- `backend/internal/routes/routes.go`

**修正箇所:**
```go
// routes.go
api.PATCH("/users/:id/approve", handlers.ApproveUser, middleware.JWTAuth(), adminMiddleware.AdminRoleCheck())
```

---

### 1.3 環境変数追加

#### タスク 1.3.1: .envファイルに管理画面用の環境変数追加
- [ ] ADMIN_BASIC_USER
- [ ] ADMIN_BASIC_PASSWORD
- [ ] ADMIN_JWT_SECRET（オプション: 管理画面専用のJWTシークレット）

**ファイル:**
- `backend/.env`

**追加内容:**
```
ADMIN_BASIC_USER=admin
ADMIN_BASIC_PASSWORD=your-basic-auth-password-change-in-production
ADMIN_JWT_SECRET=your-admin-jwt-secret-key-change-in-production
```

---

## Phase 2: 基本機能（推定3日）

### 2.1 管理者ログイン画面

#### タスク 2.1.1: ログイン画面テンプレート作成
- [ ] html/templateでログイン画面作成
- [ ] Bulma CSSフレームワーク適用
- [ ] フォーム（username, password）

**ファイル:**
- `backend/internal/admin/templates/login.html`

**実装例:**
```html
<!DOCTYPE html>
<html lang="ja">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>管理画面 - ログイン</title>
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bulma@0.9.4/css/bulma.min.css">
</head>
<body>
    <section class="hero is-fullheight">
        <div class="hero-body">
            <div class="container">
                <div class="columns is-centered">
                    <div class="column is-5-tablet is-4-desktop is-3-widescreen">
                        <div class="box">
                            <h1 class="title has-text-centered">管理画面</h1>
                            {{if .Error}}
                            <div class="notification is-danger">
                                {{.Error}}
                            </div>
                            {{end}}
                            <form action="/admin/login" method="POST">
                                <div class="field">
                                    <label class="label">ユーザー名</label>
                                    <div class="control">
                                        <input class="input" type="text" name="username" required>
                                    </div>
                                </div>
                                <div class="field">
                                    <label class="label">パスワード</label>
                                    <div class="control">
                                        <input class="input" type="password" name="password" required>
                                    </div>
                                </div>
                                <div class="field">
                                    <button class="button is-primary is-fullwidth" type="submit">ログイン</button>
                                </div>
                            </form>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </section>
</body>
</html>
```

---

#### タスク 2.1.2: ログインハンドラー作成
- [ ] GET /admin/login: ログイン画面表示
- [ ] POST /admin/login: ログイン処理、JWT発行
- [ ] ログイン失敗時のエラー表示
- [ ] last_login_at を更新

**ファイル:**
- `backend/internal/admin/handlers/auth_handler.go`

**実装例:**
```go
package handlers

import (
    "net/http"
    "time"

    "github.com/labstack/echo/v4"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
    "your-project/internal/models"
    "your-project/internal/utils"
)

type AuthHandler struct {
    DB *gorm.DB
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
    return &AuthHandler{DB: db}
}

// GET /admin/login
func (h *AuthHandler) ShowLoginPage(c echo.Context) error {
    return c.Render(http.StatusOK, "login.html", nil)
}

// POST /admin/login
func (h *AuthHandler) Login(c echo.Context) error {
    username := c.FormValue("username")
    password := c.FormValue("password")

    var user models.User
    if err := h.DB.Where("username = ? AND role = ?", username, "admin").First(&user).Error; err != nil {
        return c.Render(http.StatusOK, "login.html", map[string]interface{}{
            "Error": "ユーザー名またはパスワードが正しくありません",
        })
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        return c.Render(http.StatusOK, "login.html", map[string]interface{}{
            "Error": "ユーザー名またはパスワードが正しくありません",
        })
    }

    // JWT生成
    token, err := utils.GenerateJWT(user.ID)
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate token")
    }

    // last_login_at を更新
    now := time.Now()
    user.LastLoginAt = &now
    h.DB.Save(&user)

    // HttpOnly Cookieに保存
    cookie := &http.Cookie{
        Name:     "admin_token",
        Value:    token,
        Path:     "/admin",
        HttpOnly: true,
        Secure:   true,
        SameSite: http.SameSiteNoneMode,
        MaxAge:   86400, // 24時間
    }
    c.SetCookie(cookie)

    return c.Redirect(http.StatusSeeOther, "/admin/dashboard")
}

// POST /admin/logout
func (h *AuthHandler) Logout(c echo.Context) error {
    cookie := &http.Cookie{
        Name:     "admin_token",
        Value:    "",
        Path:     "/admin",
        HttpOnly: true,
        MaxAge:   -1,
    }
    c.SetCookie(cookie)

    return c.Redirect(http.StatusSeeOther, "/admin/login")
}
```

---

### 2.2 レイアウトテンプレート

#### タスク 2.2.1: 共通レイアウト作成
- [ ] ヘッダー（ロゴ、ログアウトボタン）
- [ ] サイドバー（ナビゲーション）
- [ ] フッター
- [ ] パンくずリスト

**ファイル:**
- `backend/internal/admin/templates/layout.html`

**実装例:**
```html
<!DOCTYPE html>
<html lang="ja">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}} - 管理画面</title>
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bulma@0.9.4/css/bulma.min.css">
    <script src="https://cdn.jsdelivr.net/npm/chart.js@4.4.0/dist/chart.umd.min.js"></script>
    <link rel="stylesheet" href="/static/css/admin.css">
</head>
<body>
    <!-- Header -->
    <nav class="navbar is-dark" role="navigation">
        <div class="navbar-brand">
            <a class="navbar-item" href="/admin/dashboard">
                <strong>SNS管理画面</strong>
            </a>
        </div>
        <div class="navbar-menu">
            <div class="navbar-end">
                <div class="navbar-item">
                    <span class="tag is-light">{{.AdminUsername}}</span>
                </div>
                <div class="navbar-item">
                    <form action="/admin/logout" method="POST">
                        <button class="button is-light is-small" type="submit">ログアウト</button>
                    </form>
                </div>
            </div>
        </div>
    </nav>

    <div class="columns is-gapless">
        <!-- Sidebar -->
        <aside class="column is-2 menu has-background-light" style="min-height: calc(100vh - 52px);">
            <p class="menu-label">メニュー</p>
            <ul class="menu-list">
                <li><a href="/admin/dashboard" class="{{if eq .Active "dashboard"}}is-active{{end}}">ダッシュボード</a></li>
                <li><a href="/admin/users" class="{{if eq .Active "users"}}is-active{{end}}">ユーザー一覧</a></li>
                <li><a href="/admin/password-resets" class="{{if eq .Active "password-resets"}}is-active{{end}}">パスワードリセット</a></li>
                <li><a href="/admin/logs" class="{{if eq .Active "logs"}}is-active{{end}}">操作ログ</a></li>
            </ul>
        </aside>

        <!-- Main Content -->
        <div class="column">
            <section class="section">
                <!-- Breadcrumb -->
                <nav class="breadcrumb" aria-label="breadcrumbs">
                    <ul>
                        {{range .Breadcrumbs}}
                        <li class="{{if .Active}}is-active{{end}}">
                            <a href="{{.URL}}">{{.Name}}</a>
                        </li>
                        {{end}}
                    </ul>
                </nav>

                <!-- Page Content -->
                {{template "content" .}}
            </section>
        </div>
    </div>

    <script src="/static/js/admin.js"></script>
</body>
</html>
```

---

### 2.3 ユーザー一覧

#### タスク 2.3.1: ユーザー一覧API作成
- [ ] GET /admin/api/users
- [ ] フィルタ（status, role）
- [ ] 検索（username, email）
- [ ] ソート（created_at, last_login_at）
- [ ] ページネーション

**ファイル:**
- `backend/internal/admin/handlers/user_handler.go`

**実装例:**
```go
package handlers

import (
    "net/http"
    "strconv"

    "github.com/labstack/echo/v4"
    "gorm.io/gorm"
    "your-project/internal/models"
)

type UserHandler struct {
    DB *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
    return &UserHandler{DB: db}
}

func (h *UserHandler) GetUsers(c echo.Context) error {
    // クエリパラメータ
    status := c.QueryParam("status")
    role := c.QueryParam("role")
    search := c.QueryParam("search")
    sort := c.QueryParam("sort")
    order := c.QueryParam("order")
    page, _ := strconv.Atoi(c.QueryParam("page"))
    limit, _ := strconv.Atoi(c.QueryParam("limit"))

    if page < 1 {
        page = 1
    }
    if limit < 1 {
        limit = 20
    }
    if sort == "" {
        sort = "created_at"
    }
    if order == "" {
        order = "desc"
    }

    // クエリ構築
    query := h.DB.Model(&models.User{})

    if status != "" && status != "all" {
        query = query.Where("status = ?", status)
    }
    if role != "" && role != "all" {
        query = query.Where("role = ?", role)
    }
    if search != "" {
        query = query.Where("username LIKE ? OR email LIKE ?", "%"+search+"%", "%"+search+"%")
    }

    // 総件数
    var total int64
    query.Count(&total)

    // ソート
    query = query.Order(sort + " " + order)

    // ページネーション
    offset := (page - 1) * limit
    var users []models.User
    query.Limit(limit).Offset(offset).Find(&users)

    // 投稿数を取得（サブクエリ）
    type UserWithPostCount struct {
        models.User
        PostCount int `json:"post_count"`
    }

    var usersWithPostCount []UserWithPostCount
    for _, user := range users {
        var postCount int64
        h.DB.Model(&models.Post{}).Where("user_id = ?", user.ID).Count(&postCount)
        usersWithPostCount = append(usersWithPostCount, UserWithPostCount{
            User:      user,
            PostCount: int(postCount),
        })
    }

    totalPages := int(total) / limit
    if int(total)%limit > 0 {
        totalPages++
    }

    return c.JSON(http.StatusOK, map[string]interface{}{
        "data": map[string]interface{}{
            "users": usersWithPostCount,
            "pagination": map[string]interface{}{
                "total":       total,
                "page":        page,
                "limit":       limit,
                "total_pages": totalPages,
            },
        },
    })
}
```

---

#### タスク 2.3.2: ユーザー一覧画面テンプレート作成
- [ ] フィルタフォーム
- [ ] 検索フォーム
- [ ] ユーザーテーブル
- [ ] ページネーション
- [ ] 一括操作チェックボックス

**ファイル:**
- `backend/internal/admin/templates/users/index.html`

---

### 2.4 ユーザー詳細

#### タスク 2.4.1: ユーザー詳細API作成
- [ ] GET /admin/api/users/:id
- [ ] ユーザー基本情報
- [ ] 統計（投稿数、いいね数、フォロワー数、フォロー数）
- [ ] 最近の投稿（5件）

**ファイル:**
- `backend/internal/admin/handlers/user_handler.go`

---

#### タスク 2.4.2: ユーザー詳細画面テンプレート作成
- [ ] 基本情報カード
- [ ] 統計カード
- [ ] 最近の投稿一覧
- [ ] 操作ボタン（承認/拒否/ステータス変更）

**ファイル:**
- `backend/internal/admin/templates/users/detail.html`

---

### 2.5 ユーザー承認/拒否

#### タスク 2.5.1: ユーザーステータス変更API
- [ ] PATCH /admin/api/users/:id/status
- [ ] リクエスト: {"status": "approved"}
- [ ] 管理操作ログ記録
- [ ] メールテンプレート生成

**ファイル:**
- `backend/internal/admin/handlers/user_handler.go`

---

#### タスク 2.5.2: 一括ステータス変更API
- [ ] POST /admin/api/users/batch-update-status
- [ ] リクエスト: {"user_ids": [1,2,3], "status": "approved"}
- [ ] 管理操作ログ記録

**ファイル:**
- `backend/internal/admin/handlers/user_handler.go`

---

### 2.6 ルート設定

#### タスク 2.6.1: 管理画面ルート定義
- [ ] /admin/login (GET, POST)
- [ ] /admin/logout (POST)
- [ ] /admin/dashboard (GET)
- [ ] /admin/users (GET)
- [ ] /admin/users/:id (GET)
- [ ] /admin/api/* (管理API)

**ファイル:**
- `backend/internal/routes/admin_routes.go`

**実装例:**
```go
package routes

import (
    "github.com/labstack/echo/v4"
    "gorm.io/gorm"
    "your-project/internal/admin/handlers"
    adminMiddleware "your-project/internal/admin/middleware"
)

func SetupAdminRoutes(e *echo.Echo, db *gorm.DB) {
    // Basic認証を全体に適用
    admin := e.Group("/admin", adminMiddleware.BasicAuth())

    // ハンドラー初期化
    authHandler := handlers.NewAuthHandler(db)
    userHandler := handlers.NewUserHandler(db)

    // 認証不要のルート
    admin.GET("/login", authHandler.ShowLoginPage)
    admin.POST("/login", authHandler.Login)

    // JWT認証 + Adminロール必須のルート
    adminAuth := admin.Group("", middleware.JWTAuth(), adminMiddleware.AdminRoleCheck())

    // 画面
    adminAuth.GET("/dashboard", dashboardHandler.ShowDashboard)
    adminAuth.GET("/users", userHandler.ShowUserList)
    adminAuth.GET("/users/:id", userHandler.ShowUserDetail)

    // API
    api := adminAuth.Group("/api")
    api.POST("/logout", authHandler.Logout)
    api.GET("/users", userHandler.GetUsers)
    api.GET("/users/:id", userHandler.GetUser)
    api.PATCH("/users/:id/status", userHandler.UpdateUserStatus)
    api.POST("/users/batch-update-status", userHandler.BatchUpdateUserStatus)
}
```

---

## Phase 3: パスワードリセット（推定2日）

### 3.1 パスワードリセット申請API（フロントエンド側）

#### タスク 3.1.1: パスワードリセット申請API作成
- [ ] POST /api/v1/password-reset/request
- [ ] リクエスト: {"email": "user@example.com"}
- [ ] password_reset_requestsテーブルにレコード作成（status: pending）

**ファイル:**
- `backend/internal/handlers/password_reset_handler.go`（フロントエンド用API）

---

### 3.2 管理画面: パスワードリセット申請一覧

#### タスク 3.2.1: パスワードリセット申請一覧API
- [ ] GET /admin/api/password-resets
- [ ] フィルタ（status）
- [ ] ページネーション

**ファイル:**
- `backend/internal/admin/handlers/password_reset_handler.go`

---

#### タスク 3.2.2: パスワードリセット申請一覧画面テンプレート
- [ ] フィルタフォーム
- [ ] 申請テーブル
- [ ] リセットリンク発行ボタン

**ファイル:**
- `backend/internal/admin/templates/password_resets/index.html`

---

### 3.3 トークン発行機能

#### タスク 3.3.1: トークン発行API
- [ ] POST /admin/api/password-resets/:id/approve
- [ ] UUID生成
- [ ] 有効期限: 24時間
- [ ] status: approved
- [ ] リセットURL生成
- [ ] メールテンプレート生成
- [ ] 管理操作ログ記録

**ファイル:**
- `backend/internal/admin/handlers/password_reset_handler.go`

**実装例:**
```go
func (h *PasswordResetHandler) ApproveResetRequest(c echo.Context) error {
    id := c.Param("id")
    adminUser := c.Get("admin_user").(models.User)

    var request models.PasswordResetRequest
    if err := h.DB.Preload("User").First(&request, id).Error; err != nil {
        return echo.NewHTTPError(http.StatusNotFound, "Request not found")
    }

    if request.Status != "pending" {
        return echo.NewHTTPError(http.StatusBadRequest, "Request already processed")
    }

    // トークン生成
    token := uuid.New().String()
    expiresAt := time.Now().Add(24 * time.Hour)

    // レコード更新
    request.Token = token
    request.Status = "approved"
    request.AdminApprovedBy = &adminUser.ID
    now := time.Now()
    request.AdminApprovedAt = &now
    request.ExpiresAt = expiresAt
    h.DB.Save(&request)

    // リセットURL生成
    frontendURL := os.Getenv("FRONTEND_URL")
    resetURL := fmt.Sprintf("%s/reset-password?token=%s", frontendURL, token)

    // メールテンプレート生成
    emailTemplate := fmt.Sprintf(`こんにちは、%sさん

パスワードリセットのリクエストを承認しました。
以下のリンクからパスワードを再設定してください（24時間有効）：

%s

このリクエストに心当たりがない場合は、このメールを無視してください。

よろしくお願いいたします。
SNS運営チーム`, request.User.Username, resetURL)

    // 管理操作ログ記録
    utils.LogAdminAction(h.DB, utils.AdminLogParams{
        AdminID:        adminUser.ID,
        AdminUsername:  adminUser.Username,
        Action:         "password_reset_approve",
        TargetUserID:   &request.UserID,
        TargetUsername: &request.User.Username,
        Details:        fmt.Sprintf("Reset request ID: %d, Token: %s (truncated), Expires: %s", request.ID, token[:8]+"****", expiresAt.Format(time.RFC3339)),
        IP:             c.RealIP(),
    })

    return c.JSON(http.StatusOK, map[string]interface{}{
        "data": map[string]interface{}{
            "token":          token,
            "reset_url":      resetURL,
            "expires_at":     expiresAt,
            "email_template": emailTemplate,
        },
        "message": "Password reset approved",
    })
}
```

---

### 3.4 クリップボードコピー機能（JavaScript）

#### タスク 3.4.1: クリップボードコピーボタン
- [ ] リセットURLコピーボタン
- [ ] メールテンプレートコピーボタン
- [ ] コピー成功時の通知

**ファイル:**
- `backend/static/js/admin.js`

**実装例:**
```javascript
function copyToClipboard(text, buttonElement) {
    navigator.clipboard.writeText(text).then(() => {
        // 成功通知
        const originalText = buttonElement.textContent;
        buttonElement.textContent = 'コピーしました！';
        buttonElement.classList.add('is-success');

        setTimeout(() => {
            buttonElement.textContent = originalText;
            buttonElement.classList.remove('is-success');
        }, 2000);
    }).catch(err => {
        alert('コピーに失敗しました');
    });
}
```

---

## Phase 4: ダッシュボード（推定3日）

### 4.1 統計データAPI

#### タスク 4.1.1: ダッシュボード統計API
- [ ] GET /admin/api/dashboard/stats
- [ ] ユーザー統計（total, pending, approved, rejected）
- [ ] 投稿統計（total, today, with_media_rate）
- [ ] アクティブユーザー数（直近7日）
- [ ] パスワードリセット申請件数（pending）
- [ ] アラート判定

**ファイル:**
- `backend/internal/admin/handlers/dashboard_handler.go`

---

#### タスク 4.1.2: 投稿数推移API
- [ ] GET /admin/api/dashboard/charts/posts
- [ ] 直近30日の投稿数
- [ ] Chart.js用のJSON形式

**ファイル:**
- `backend/internal/admin/handlers/dashboard_handler.go`

---

#### タスク 4.1.3: 新規登録ユーザー推移API
- [ ] GET /admin/api/dashboard/charts/users
- [ ] 直近30日の新規登録ユーザー数
- [ ] Chart.js用のJSON形式

**ファイル:**
- `backend/internal/admin/handlers/dashboard_handler.go`

---

### 4.2 ダッシュボード画面

#### タスク 4.2.1: ダッシュボードテンプレート作成
- [ ] 統計カード（ユーザー、投稿）
- [ ] グラフエリア（Chart.js）
- [ ] アラート表示

**ファイル:**
- `backend/internal/admin/templates/dashboard.html`

---

#### タスク 4.2.2: Chart.js グラフ描画
- [ ] 投稿数推移グラフ
- [ ] 新規登録ユーザー推移グラフ

**ファイル:**
- `backend/static/js/admin.js`

**実装例:**
```javascript
async function loadPostsChart() {
    const response = await fetch('/admin/api/dashboard/charts/posts');
    const data = await response.json();

    const ctx = document.getElementById('postsChart').getContext('2d');
    new Chart(ctx, {
        type: 'line',
        data: {
            labels: data.data.labels,
            datasets: [{
                label: '投稿数',
                data: data.data.values,
                borderColor: '#00d1b2',
                backgroundColor: 'rgba(0, 209, 178, 0.1)',
                tension: 0.4
            }]
        },
        options: {
            responsive: true,
            plugins: {
                legend: {
                    display: true
                }
            }
        }
    });
}
```

---

### 4.3 アラート機能

#### タスク 4.3.1: アラート判定ロジック
- [ ] 承認待ちユーザーが10人以上
- [ ] パスワードリセット申請が5件以上
- [ ] レートリミット超過が直近1時間で50回以上

**ファイル:**
- `backend/internal/admin/services/dashboard_service.go`

---

## Phase 5: 操作ログ（推定1日）

### 5.1 操作ログAPI

#### タスク 5.1.1: 操作ログ一覧API
- [ ] GET /admin/api/logs
- [ ] フィルタ（action, admin_username, 日付範囲）
- [ ] ページネーション（50件/ページ）

**ファイル:**
- `backend/internal/admin/handlers/log_handler.go`

---

### 5.2 操作ログ画面

#### タスク 5.2.1: 操作ログ一覧テンプレート
- [ ] フィルタフォーム
- [ ] ログテーブル
- [ ] ページネーション

**ファイル:**
- `backend/internal/admin/templates/logs/index.html`

---

## Phase 6: テスト・デプロイ（推定2日）

### 6.1 手動テスト

#### タスク 6.1.1: 機能テスト
- [ ] 管理者ログイン
- [ ] ユーザー一覧（フィルタ、検索、ページネーション）
- [ ] ユーザー詳細
- [ ] ユーザー承認/拒否
- [ ] パスワードリセット申請承認
- [ ] ダッシュボード（統計、グラフ）
- [ ] 操作ログ

---

#### タスク 6.1.2: セキュリティテスト
- [ ] Basic認証が正しく動作するか
- [ ] Adminロールでないユーザーがアクセスできないか
- [ ] JWT認証が正しく動作するか
- [ ] レートリミットが動作するか
- [ ] 管理操作ログが記録されているか

---

### 6.2 本番環境設定

#### タスク 6.2.1: 環境変数設定
- [ ] ADMIN_BASIC_USER
- [ ] ADMIN_BASIC_PASSWORD（強力なパスワードに変更）
- [ ] ADMIN_JWT_SECRET（強力なシークレットに変更）
- [ ] FRONTEND_URL

---

#### タスク 6.2.2: デプロイ
- [ ] バックエンドデプロイ（Render/Cloud Run）
- [ ] 静的ファイル配信（CSS, JS）
- [ ] データベースマイグレーション実行

---

## 追加タスク（オプション）

### メール送信自動化（Phase 3以降）
- [ ] SendGridやGmail SMTP設定
- [ ] メール送信機能実装
- [ ] 管理画面から直接メール送信ボタン

### CSVエクスポート
- [ ] ユーザー一覧をCSVでエクスポート
- [ ] 操作ログをCSVでエクスポート

### 高度な分析（Phase 3以降）
- [ ] いいね数推移グラフ
- [ ] コメント数推移グラフ
- [ ] エンゲージメント率

---

以上が管理画面の実装タスクリストです。
各タスクを段階的に実装していくことで、効率的に管理画面を構築できます。
