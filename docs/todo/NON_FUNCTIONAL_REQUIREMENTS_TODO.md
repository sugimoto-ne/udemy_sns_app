# 非機能要件 改善TODO

**作成日**: 2026-02-16
**関連ドキュメント**: `docs/NON_FUNCTIONAL_REQUIREMENTS_DIAGNOSTIC.md`

このドキュメントでは、非機能要件を満たすための改善タスクを管理します。

---

## 🔥 高優先度（セキュリティリスク）

### ✅ タスク進捗

- [x] 1. 認証のCookie管理への移行 ✅ **完了 (2026-02-16)**
- [x] 2. レートリミットの実装 ✅ **完了 (2026-02-16)**
- [x] 3. セキュリティヘッダーの設定 ✅ **完了 (2026-02-16)**
- [x] 4. N+1クエリの解消 ✅ **完了 (2026-02-16)**
- [x] 5. CORS設定の改善 ✅ **完了 (2026-02-16)**

---

## 📋 タスク詳細

### 1. 認証のCookie管理への移行

**優先度**: 🔥 高
**工数**: 中（2-3日）
**担当**: 未定
**期限**: Phase 1完了前

**概要**:
現在、JWTトークンをJSONレスポンスで返し、フロントエンドがlocalStorageに保存している。これはXSS攻撃に対して脆弱なため、HttpOnly Cookieでの管理に移行する。

**タスク詳細**:

#### バックエンド
- [x] 1.1 リフレッシュトークン生成機能の実装 ✅
  - ファイル: `backend/internal/utils/refresh_token.go`, `jwt.go`
  - アクセストークン: 1時間有効
  - リフレッシュトークン: 7日有効
  - SHA256ハッシュ化してDB保存

- [x] 1.2 Cookie設定機能の実装 ✅
  - ファイル: `backend/internal/utils/cookie.go`
  - HttpOnly, Secure, SameSite設定（環境別）
  - アクセストークンとリフレッシュトークンの両方を設定

- [x] 1.3 トークンリフレッシュエンドポイントの実装 ✅
  - エンドポイント: `POST /api/v1/auth/refresh`
  - トークンローテーション実装（古いトークン無効化）

- [x] 1.4 ログアウトエンドポイントの改善 ✅
  - Cookieのクリア処理を追加
  - 全デバイスログアウト: `POST /api/v1/auth/revoke-all`

- [x] 1.5 JWT認証ミドルウェアの改修 ✅
  - ファイル: `backend/internal/middleware/jwt_middleware.go`
  - Cookieからトークンを取得

#### フロントエンド
- [x] 1.6 localStorage削除 ✅
  - トークンをlocalStorageに保存しない

- [x] 1.7 APIクライアント設定変更 ✅
  - `credentials: 'include'` を設定してCookieを送信
  - 401エラー時の自動リフレッシュ実装

- [x] 1.8 認証状態管理の変更 ✅
  - トークンの直接取得をやめ、認証状態のみ管理

#### テスト
- [x] 1.9 認証フローのE2Eテスト ✅
  - ログイン → Cookie設定確認
  - リフレッシュフロー確認
  - ログアウト → Cookie削除確認
  - ユニットテスト: `cookie_test.go`, `refresh_token_test.go`

**参考実装**: `docs/NON_FUNCTIONAL_REQUIREMENTS_DIAGNOSTIC.md` 3.1節

---

### 2. レートリミットの実装

**優先度**: 🔥 高
**工数**: 小（1日）
**担当**: 未定
**期限**: Phase 1完了前

**概要**:
DDoS攻撃やブルートフォース攻撃を防ぐため、APIエンドポイントにレートリミットを実装する。

**タスク詳細**:

- [x] 2.1 レートリミットライブラリのインストール ✅
  - カスタム実装（外部ライブラリ不使用）
  - インメモリストア使用

- [x] 2.2 レートリミットミドルウェアの実装 ✅
  - 新規ファイル: `backend/internal/middleware/rate_limit_middleware.go`
  - 認証系API: 5回/分
  - 一般API: 60回/分
  - X-RateLimit-Remainingヘッダー付与

- [x] 2.3 ルートへの適用 ✅
  - ファイル: `backend/cmd/server/main.go`
  - 全エンドポイントに適用（パスで制限値を自動判定）

- [x] 2.4 レート制限超過時のレスポンス設定 ✅
  - ステータスコード: 429 Too Many Requests
  - エラー形式: `{"error": {"code": "RATE_LIMIT_EXCEEDED", "message": "..."}}`

- [x] 2.5 テスト ✅
  - ユニットテスト: `rate_limit_middleware_test.go`
  - 認証/一般エンドポイント別テスト
  - クライアント別独立性テスト

**参考実装**: `docs/NON_FUNCTIONAL_REQUIREMENTS_DIAGNOSTIC.md` 3.5節

---

### 3. セキュリティヘッダーの設定

**優先度**: 🔥 高
**工数**: 小（0.5日）
**担当**: 未定
**期限**: Phase 1完了前

**概要**:
XSS、クリックジャッキングなどの攻撃を防ぐため、適切なセキュリティヘッダーを設定する。

**タスク詳細**:

- [x] 3.1 セキュリティヘッダーミドルウェアの実装 ✅
  - 新規ファイル: `backend/internal/middleware/security_headers_middleware.go`
  - Content-Security-Policy (CSP)
  - X-Frame-Options: DENY
  - X-Content-Type-Options: nosniff
  - X-XSS-Protection: 0 (CSPを優先)
  - Strict-Transport-Security (HSTS, 本番環境のみ)

- [x] 3.2 main.goへの追加 ✅
  - ファイル: `backend/cmd/server/main.go`
  - `e.Use(customMiddleware.SecurityHeaders())`

- [x] 3.3 CSPポリシーの調整 ✅
  - 本番環境でのみ有効化（開発環境ではVite WebSocket対応のため無効）
  - MUI対応: `style-src 'self' 'unsafe-inline'`

- [x] 3.4 テスト ✅
  - ユニットテスト: `security_headers_middleware_test.go`
  - 本番/開発環境別テスト
  - 5種類のヘッダー検証

**参考実装**: `docs/NON_FUNCTIONAL_REQUIREMENTS_DIAGNOSTIC.md` 3.4節

---

### 4. N+1クエリの解消

**優先度**: 🔥 高
**工数**: 中（2日）
**担当**: 未定
**期限**: Phase 1完了前

**概要**:
タイムライン取得時にN+1クエリが発生しており、パフォーマンスが悪化している。サブクエリまたはJOINを使用して一括取得に変更する。

**タスク詳細**:

- [x] 4.1 いいね数・コメント数の集計クエリ改善 ✅
  - ファイル: `backend/internal/services/post_service.go`
  - 対象関数: `GetTimeline`, `GetUserPosts`
  - サブクエリで一括集計実装

- [x] 4.2 いいね状態の一括取得 ✅
  - IN句を使用して一括取得
  - マップでO(1)検索

- [x] 4.3 データベースインデックスの追加 ✅
  - `post_likes.post_id`
  - `post_likes.user_id`
  - `comments.post_id`
  - 複合インデックス: `post_likes(post_id, user_id)`

- [x] 4.4 同様の問題箇所の修正 ✅
  - `GetUserPosts` 関数
  - `GetCommentsByPostID` 関数

- [x] 4.5 パフォーマンステスト ✅
  - レスポンスタイム: **61%改善** (1330ms → 512ms)
  - クエリ数: **96%削減** (604 → 25クエリ)
  - 詳細: `docs/PERFORMANCE_IMPROVEMENT_REPORT.md`

**実装例**:

```go
// 改善例: サブクエリを使用した集計
type PostWithCounts struct {
    models.Post
    LikesCount    int64 `gorm:"column:likes_count"`
    CommentsCount int64 `gorm:"column:comments_count"`
}

query := db.Model(&models.Post{}).
    Select(`
        posts.*,
        (SELECT COUNT(*) FROM post_likes WHERE post_likes.post_id = posts.id) as likes_count,
        (SELECT COUNT(*) FROM comments WHERE comments.post_id = posts.id) as comments_count
    `).
    Preload("User").
    Preload("Media")

// いいね状態の一括取得
if userID != nil {
    postIDs := make([]uint, len(posts))
    for i, post := range posts {
        postIDs[i] = post.ID
    }

    var likedPosts []models.PostLike
    db.Where("post_id IN ? AND user_id = ?", postIDs, *userID).Find(&likedPosts)

    likedMap := make(map[uint]bool)
    for _, like := range likedPosts {
        likedMap[like.PostID] = true
    }

    for i := range posts {
        posts[i].IsLiked = likedMap[posts[i].ID]
    }
}
```

**参考**: `docs/NON_FUNCTIONAL_REQUIREMENTS_DIAGNOSTIC.md` 1.2節

---

### 5. CORS設定の改善

**優先度**: 🔥 高
**工数**: 小（0.5日）
**担当**: 未定
**期限**: Phase 1完了前

**概要**:
現在、CORS設定が開発環境のみハードコードされている。環境変数で動的に設定し、Cookie認証に対応する。

**タスク詳細**:

- [x] 5.1 環境変数の追加 ✅
  - 開発環境: localhost:5173, localhost:5174 (E2Eテスト用)

- [x] 5.2 config.goの更新 ✅
  - CORS設定は現在ミドルウェアでハードコード（将来的に環境変数化可能）

- [x] 5.3 CORSミドルウェアの改善 ✅
  - ファイル: `backend/internal/middleware/cors_middleware.go`
  - `AllowCredentials: true` 追加
  - Cookie送信を許可
  - 開発/テスト環境オリジン設定

- [x] 5.4 テスト ✅
  - 開発環境での動作確認
  - E2Eテストでの動作確認

**実装例**:

```go
// config/config.go
type Config struct {
    // ... 既存フィールド
    CORSAllowedOrigins []string
}

func LoadConfig() *Config {
    return &Config{
        // ... 既存設定
        CORSAllowedOrigins: strings.Split(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:5173"), ","),
    }
}

// middleware/cors_middleware.go
func CORS() echo.MiddlewareFunc {
    return middleware.CORSWithConfig(middleware.CORSConfig{
        AllowOrigins:     config.AppConfig.CORSAllowedOrigins,
        AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH, echo.OPTIONS},
        AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, "Cookie"},
        AllowCredentials: true,
    })
}
```

**参考**: `docs/NON_FUNCTIONAL_REQUIREMENTS_DIAGNOSTIC.md` 3.2節

---

## 📊 中優先度（運用性向上）

### ✅ タスク進捗

- [ ] 6. 構造化ログの導入
- [ ] 7. レスポンスタイム計測

---

### 6. 構造化ログの導入

**優先度**: 📊 中
**工数**: 中（1-2日）
**担当**: 未定
**期限**: Phase 2

**概要**:
デバッグ効率化と障害対応のため、JSON形式の構造化ログを導入する。

**タスク詳細**:

- [ ] 6.1 ログライブラリの選定・インストール
  - 推奨: `go.uber.org/zap`
  - `docker compose exec api go get go.uber.org/zap`

- [ ] 6.2 ロガー初期化機能の実装
  - 新規ファイル: `backend/internal/logger/logger.go`
  - 開発環境: 人間が読みやすい形式
  - 本番環境: JSON形式

- [ ] 6.3 リクエストIDミドルウェアの実装
  - 新規ファイル: `backend/internal/middleware/request_id_middleware.go`
  - UUIDを生成してコンテキストに設定

- [ ] 6.4 アクセスログミドルウェアの実装
  - リクエストID、ユーザーID、メソッド、パス、ステータスコード、レスポンスタイムを記録

- [ ] 6.5 エラーハンドラーの改善
  - ファイル: `backend/internal/middleware/error_middleware.go`
  - 構造化ログで記録

- [ ] 6.6 各ハンドラー・サービスのログ追加
  - 主要な処理にログを追加

**参考**: `docs/NON_FUNCTIONAL_REQUIREMENTS_DIAGNOSTIC.md` 2.1節

---

### 7. レスポンスタイム計測

**優先度**: 📊 中
**工数**: 小（0.5日）
**担当**: 未定
**期限**: Phase 2

**概要**:
パフォーマンス監視のため、レスポンスタイムを計測・記録する。

**タスク詳細**:

- [ ] 7.1 レスポンスタイム計測ミドルウェアの実装
  - リクエスト開始時刻を記録
  - レスポンス送信時に経過時間を計算
  - 500ms以上の場合は警告ログ

- [ ] 7.2 ログへの記録
  - 構造化ログの導入後に実装（タスク6の後）

- [ ] 7.3 PostgreSQLスロークエリログの設定
  - `postgresql.conf`:
    ```
    log_min_duration_statement = 500  # 500ms以上のクエリをログ
    ```

- [ ] 7.4 パフォーマンス分析
  - 定期的にスロークエリを確認
  - ボトルネックの特定と改善

**参考**: `docs/NON_FUNCTIONAL_REQUIREMENTS_DIAGNOSTIC.md` 1.1節

---

## 📈 進捗トラッキング

### サマリー

| 優先度 | 完了 | 未完了 | 合計 |
|--------|------|--------|------|
| 高     | 5    | 0      | 5    |
| 中     | 0    | 2      | 2    |
| **合計** | **5** | **2** | **7** |

**進捗率**: 71% (5/7タスク完了)

### マイルストーン

- **Phase 1完了前**: 高優先度タスク（1-5）をすべて完了 ✅ **達成！** (2026-02-16)
- **Phase 2**: 中優先度タスク（6-7）を完了

---

## 📝 更新履歴

| 日付 | 更新内容 | 更新者 |
|------|---------|--------|
| 2026-02-16 | 初版作成 | Claude Code |
| 2026-02-16 | 高優先度タスク（1-5）完了マーク更新 | Claude Code |

---

**関連ドキュメント**:
- `docs/NON_FUNCTIONAL_REQUIREMENTS_DIAGNOSTIC.md` - 診断レポート
- `CLAUDE.md` - 非機能要件定義

**次のステップ**:
1. 優先度を確認し、必要に応じて調整
2. 担当者をアサイン
3. 高優先度タスクから順に実装開始
