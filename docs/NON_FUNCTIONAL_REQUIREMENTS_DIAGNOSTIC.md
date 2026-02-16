# 非機能要件 診断レポート

**作成日**: 2026-02-16
**対象バージョン**: Phase 1 MVP

このドキュメントでは、現在のコードベースが非機能要件をどの程度満たしているかを診断します。

---

## 📊 サマリー

| 分類 | 項目 | 達成状況 | 優先度 |
|------|------|---------|--------|
| 性能 | レスポンスタイム | 🟡 部分的 | 中 |
| 性能 | クエリ効率（N+1対策） | 🔴 未達成 | 高 |
| 運用・保守性 | ログ（構造化） | 🔴 未達成 | 中 |
| 運用・保守性 | エラーハンドリング | 🟢 達成 | - |
| セキュリティ | 認証（Cookie管理） | 🔴 未達成 | 高 |
| セキュリティ | CORS設定 | 🟡 部分的 | 高 |
| セキュリティ | 入力検証 | 🟢 達成 | - |
| セキュリティ | セキュリティヘッダー | 🔴 未達成 | 高 |
| セキュリティ | レートリミット | 🔴 未達成 | 高 |

**凡例**:
- 🟢 達成: 要件を満たしている
- 🟡 部分的: 一部実装されているが改善が必要
- 🔴 未達成: 実装されていない、または要件を満たしていない

---

## 📝 詳細診断

### 1. 性能要件

#### 1.1 レスポンスタイム（一覧取得API：500ms以内）

**達成状況**: 🟡 **部分的**

**現状**:
- レスポンスタイム計測の仕組みが未実装
- 本番環境でのパフォーマンステストが未実施

**不足点**:
1. レスポンスタイム計測ミドルウェアがない
2. パフォーマンスモニタリングツールの未導入
3. スロークエリログの未設定

**推奨対応**:
- [ ] レスポンスタイム計測ミドルウェアの実装
- [ ] ログにレスポンスタイムを記録
- [ ] PostgreSQLのスロークエリログ設定（500ms以上）
- [ ] 本番環境でのパフォーマンステスト実施

**ファイル**: 該当なし（未実装）

---

#### 1.2 クエリ効率（N+1クエリの排除）

**達成状況**: 🔴 **未達成**

**現状**:
- タイムライン取得時にN+1クエリが発生
- いいね数・コメント数の集計でループ内クエリ実行

**問題箇所**:

**ファイル**: `backend/internal/services/post_service.go:59-70`

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

**影響**:
- 20件の投稿取得時、最大60回のクエリが実行される（20件 × 3クエリ）
- ユーザー体験の悪化（レスポンス遅延）

**推奨対応**:
- [ ] **高優先度**: サブクエリまたはJOINを使用した集計クエリへの変更
- [ ] GROUP BYを使用した一括集計
- [ ] データベースインデックスの追加（`post_likes.post_id`, `comments.post_id`）
- [ ] 集計結果のキャッシュ検討（Phase 2+）

**同様の問題**:
- `backend/internal/services/post_service.go:225-228` (GetUserPosts)
- `backend/internal/services/comment_service.go:27-29` (GetCommentsByPostID)

---

### 2. 運用・保守性

#### 2.1 ログ（構造化ログ）

**達成状況**: 🔴 **未達成**

**現状**:
- Echoのデフォルトロガーを使用（非構造化）
- エラーログが標準出力に出力されるのみ

**ファイル**:
- `backend/cmd/server/main.go:69` - デフォルトのLogger使用
- `backend/internal/middleware/error_middleware.go:21` - 非構造化ログ

```go
// 現在の実装（非構造化）
log.Printf("Error: %v", err)
```

**不足点**:
1. JSON形式のログ出力がない
2. リクエストID、ユーザーID、メソッド、パスなどのコンテキスト情報がない
3. ログレベル（DEBUG, INFO, WARN, ERROR）の適切な分類がない
4. アクセスログの詳細情報不足

**推奨対応**:
- [ ] 構造化ログライブラリの導入（`zap`または`logrus`推奨）
- [ ] リクエストIDミドルウェアの実装
- [ ] アクセスログミドルウェアの強化（レスポンスタイム、ステータスコード、ユーザーID含む）
- [ ] ログレベルの適切な使い分け

**推奨実装例**:
```go
// zap使用例
logger.Info("User logged in",
    zap.String("request_id", requestID),
    zap.Uint("user_id", userID),
    zap.String("method", "POST"),
    zap.String("path", "/api/v1/auth/login"),
    zap.Int("status", 200),
    zap.Duration("response_time", duration),
)
```

---

#### 2.2 エラーハンドリング

**達成状況**: 🟢 **達成**

**現状**:
- カスタムエラーハンドラーが実装されている
- ユーザー向けの適切なエラーメッセージを返している
- 内部エラーを露出させていない

**ファイル**:
- `backend/internal/middleware/error_middleware.go` - カスタムエラーハンドラー
- `backend/internal/utils/response.go` - エラーレスポンスヘルパー
- `backend/internal/handlers/auth_handler.go:54-60` - 適切なエラーメッセージ

**良い点**:
```go
if err.Error() == "email already exists" {
    return utils.ErrorResponse(c, 409, "このメールアドレスは既に登録されています")
}
// 内部エラーは隠蔽
return utils.ErrorResponse(c, 500, "ユーザー登録に失敗しました")
```

**改善提案**（任意）:
- [ ] エラーコード体系の導入（例: `E001_EMAIL_ALREADY_EXISTS`）
- [ ] エラーログの強化（スタックトレースの記録）

---

### 3. セキュリティ

#### 3.1 認証（HttpOnly Cookie管理）

**達成状況**: 🔴 **未達成**

**現状**:
- JWTトークンをJSONレスポンスで返している
- フロントエンドがlocalStorageに保存する設計（XSS脆弱性）
- Cookie管理が未実装

**ファイル**:
- `backend/internal/handlers/auth_handler.go:64-75` - トークンをJSON返却
- `backend/internal/utils/jwt.go:16-33` - トークン生成（有効期限24時間）

**現在の実装**:
```go
// レスポンス
response := AuthResponse{
    User:  user.ToPublicUser(),
    Token: token,  // ❌ JSONで返している
}
return utils.SuccessResponse(c, 201, response)
```

**セキュリティリスク**:
1. **XSS攻撃**: JavaScriptからlocalStorageのトークンにアクセス可能
2. **トークン盗難**: XSSでトークンが盗まれる可能性

**推奨対応** (**高優先度**):
- [ ] **必須**: HttpOnly Cookieでのトークン管理への移行
- [ ] アクセストークン有効期限を1時間に短縮
- [ ] リフレッシュトークンの実装（有効期限7日）
- [ ] Cookie属性の設定: `HttpOnly`, `Secure`, `SameSite=None`（本番環境）

**推奨実装例**:
```go
// アクセストークンをHttpOnly Cookieに設定
c.SetCookie(&http.Cookie{
    Name:     "access_token",
    Value:    accessToken,
    Path:     "/",
    HttpOnly: true,
    Secure:   true,           // HTTPS環境でのみ送信
    SameSite: http.SameSiteNoneMode,  // クロスサイトリクエスト許可
    MaxAge:   3600,           // 1時間
})

// リフレッシュトークンもHttpOnly Cookieに設定
c.SetCookie(&http.Cookie{
    Name:     "refresh_token",
    Value:    refreshToken,
    Path:     "/api/v1/auth/refresh",
    HttpOnly: true,
    Secure:   true,
    SameSite: http.SameSiteNoneMode,
    MaxAge:   604800,         // 7日間
})
```

---

#### 3.2 CORS設定

**達成状況**: 🟡 **部分的**

**現状**:
- CORS設定は実装されている
- 開発環境のオリジンのみ許可
- 本番環境の設定が未実装

**ファイル**: `backend/internal/middleware/cors_middleware.go:10-14`

```go
return middleware.CORSWithConfig(middleware.CORSConfig{
    AllowOrigins: []string{"http://localhost:3000", "http://localhost:5173"}, // ✅ 開発環境
    AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH, echo.OPTIONS},
    AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
})
```

**不足点**:
1. 本番環境のオリジン設定がハードコード
2. Cookie使用時の `AllowCredentials: true` 設定がない
3. 環境変数による動的設定がない

**推奨対応**:
- [ ] 環境変数 `CORS_ALLOWED_ORIGINS` を追加
- [ ] `AllowCredentials: true` を追加（Cookie認証実装時）
- [ ] `AllowHeaders` に `Cookie` を追加

**推奨実装例**:
```go
allowedOrigins := strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ",")
return middleware.CORSWithConfig(middleware.CORSConfig{
    AllowOrigins:     allowedOrigins,
    AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH, echo.OPTIONS},
    AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, "Cookie"},
    AllowCredentials: true,  // Cookie使用時に必須
})
```

`.env`:
```bash
CORS_ALLOWED_ORIGINS=http://localhost:5173,https://your-production-domain.com
```

---

#### 3.3 入力検証

**達成状況**: 🟢 **達成**

**現状**:
- `validator/v10` を使用したバリデーションが実装されている
- リクエスト構造体にバリデーションタグが設定されている
- カスタムバリデーション関数も実装

**ファイル**:
- `backend/internal/handlers/auth_handler.go:11-13` - バリデーションタグ
- `backend/internal/utils/validator.go` - バリデーション実装

**良い点**:
```go
type RegisterRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
    Username string `json:"username" validate:"required,min=3,max=50"`
}
```

**継続推奨**:
- すべての入力に対してバリデーションを継続実施
- カスタムバリデーション関数の拡充

---

#### 3.4 セキュリティヘッダー

**達成状況**: 🔴 **未達成**

**現状**:
- セキュリティヘッダーが未設定
- XSS、クリックジャッキング対策が未実装

**不足しているヘッダー**:
1. `Content-Security-Policy` (CSP) - XSS対策
2. `X-Frame-Options` - クリックジャッキング対策
3. `X-Content-Type-Options` - MIMEタイプスニッフィング対策
4. `Strict-Transport-Security` (HSTS) - HTTPS強制（本番環境）
5. `X-XSS-Protection` - XSS対策（旧ブラウザ向け）

**推奨対応** (**高優先度**):
- [ ] セキュリティヘッダーミドルウェアの実装
- [ ] CSPポリシーの定義
- [ ] 本番環境でのHSTS有効化

**推奨実装**:

新規ファイル: `backend/internal/middleware/security_headers_middleware.go`

```go
package middleware

import (
    "github.com/labstack/echo/v4"
)

// SecurityHeaders - セキュリティヘッダーを設定するミドルウェア
func SecurityHeaders() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            // Content-Security-Policy (CSP)
            c.Response().Header().Set("Content-Security-Policy",
                "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self'; connect-src 'self'")

            // X-Frame-Options
            c.Response().Header().Set("X-Frame-Options", "DENY")

            // X-Content-Type-Options
            c.Response().Header().Set("X-Content-Type-Options", "nosniff")

            // X-XSS-Protection (旧ブラウザ向け)
            c.Response().Header().Set("X-XSS-Protection", "1; mode=block")

            // HSTS (本番環境でHTTPS使用時)
            if c.Request().Header.Get("X-Forwarded-Proto") == "https" {
                c.Response().Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
            }

            return next(c)
        }
    }
}
```

`backend/cmd/server/main.go` に追加:
```go
e.Use(customMiddleware.SecurityHeaders())
```

---

#### 3.5 レートリミット

**達成状況**: 🔴 **未達成**

**現状**:
- レートリミット機能が未実装
- DDoS攻撃、ブルートフォース攻撃に対して脆弱

**不足点**:
1. レートリミットミドルウェアがない
2. IP/ユーザーベースの制限がない
3. 認証APIの保護がない

**推奨対応** (**高優先度**):
- [ ] レートリミットミドルウェアの実装
- [ ] 認証系API: 5回/分の制限
- [ ] 一般API: 60回/分の制限
- [ ] IPベースの制限とユーザーベースの制限の両方実装

**推奨ライブラリ**:
- `golang.org/x/time/rate` (標準ライブラリ)
- `github.com/ulule/limiter/v3` (Redis対応、推奨)

**推奨実装**:

```bash
# パッケージインストール
docker compose exec api go get github.com/ulule/limiter/v3
docker compose exec api go get github.com/ulule/limiter/v3/drivers/middleware/echo
docker compose exec api go get github.com/ulule/limiter/v3/drivers/store/memory
```

新規ファイル: `backend/internal/middleware/rate_limit_middleware.go`

```go
package middleware

import (
    "github.com/labstack/echo/v4"
    "github.com/ulule/limiter/v3"
    echomiddleware "github.com/ulule/limiter/v3/drivers/middleware/echo"
    "github.com/ulule/limiter/v3/drivers/store/memory"
)

// RateLimitAuth - 認証系APIのレートリミット (5回/分)
func RateLimitAuth() echo.MiddlewareFunc {
    rate := limiter.Rate{
        Period: 1 * time.Minute,
        Limit:  5,
    }
    store := memory.NewStore()
    middleware := echomiddleware.NewMiddleware(limiter.New(store, rate))
    return middleware
}

// RateLimitGeneral - 一般APIのレートリミット (60回/分)
func RateLimitGeneral() echo.MiddlewareFunc {
    rate := limiter.Rate{
        Period: 1 * time.Minute,
        Limit:  60,
    }
    store := memory.NewStore()
    middleware := echomiddleware.NewMiddleware(limiter.New(store, rate))
    return middleware
}
```

`backend/internal/routes/routes.go` に適用:
```go
// 認証系エンドポイント（5回/分制限）
authGroup := api.Group("/auth")
authGroup.Use(customMiddleware.RateLimitAuth())
authGroup.POST("/register", handlers.Register)
authGroup.POST("/login", handlers.Login)

// 一般エンドポイント（60回/分制限）
api.Use(customMiddleware.RateLimitGeneral())
```

---

## 🎯 優先度別改善ロードマップ

### 高優先度（セキュリティリスク）

1. **認証のCookie管理への移行** (3.1)
   - リスク: XSS攻撃によるトークン盗難
   - 影響: 全ユーザー
   - 工数: 中（2-3日）

2. **レートリミットの実装** (3.5)
   - リスク: DDoS攻撃、ブルートフォース攻撃
   - 影響: サービス可用性
   - 工数: 小（1日）

3. **セキュリティヘッダーの設定** (3.4)
   - リスク: XSS、クリックジャッキング
   - 影響: 全ユーザー
   - 工数: 小（0.5日）

4. **N+1クエリの解消** (1.2)
   - リスク: パフォーマンス劣化
   - 影響: ユーザー体験
   - 工数: 中（2日）

5. **CORS設定の改善** (3.2)
   - リスク: 本番環境での動作不良
   - 影響: 本番環境のみ
   - 工数: 小（0.5日）

### 中優先度（運用性向上）

6. **構造化ログの導入** (2.1)
   - 目的: デバッグ効率化、障害対応
   - 影響: 運用チーム
   - 工数: 中（1-2日）

7. **レスポンスタイム計測** (1.1)
   - 目的: パフォーマンス監視
   - 影響: 開発・運用
   - 工数: 小（0.5日）

---

## 📌 次のステップ

1. 本レポートを確認し、優先度を再調整
2. `docs/todo/NON_FUNCTIONAL_REQUIREMENTS_TODO.md` で各項目の実装タスクを追跡
3. Phase 1完了前に高優先度項目をすべて対応
4. Phase 2以降で中優先度項目を実装

---

**レポート作成者**: Claude Code
**最終更新**: 2026-02-16
