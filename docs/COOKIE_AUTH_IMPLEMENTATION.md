# Cookie ベース認証実装ドキュメント

## 概要

このドキュメントは、SNS アプリケーションに実装された **Cookie ベース認証 (HttpOnly Cookie + JWT + Refresh Token)** の詳細を説明します。

従来の localStorage ベースのトークン管理から、より安全な Cookie ベースの認証方式に移行しました。

## 実装日

2026-02-16

## 主な変更点

### セキュリティ向上

| 項目 | 以前 (localStorage) | 現在 (Cookie) |
|------|---------------------|--------------|
| **XSS 対策** | ❌ JavaScript からアクセス可能 | ✅ HttpOnly で JavaScript からアクセス不可 |
| **CSRF 対策** | - | ✅ SameSite 属性で保護 |
| **トークン盗難リスク** | 高 | 低 |
| **トークンリフレッシュ** | なし | ✅ 自動リフレッシュ機構 |
| **トークンローテーション** | なし | ✅ リフレッシュ時に古いトークン無効化 |

### トークン仕様

| トークン | 有効期限 | 保存場所 | 用途 |
|---------|---------|---------|------|
| **Access Token** | 1時間 | HttpOnly Cookie | API アクセス認証 |
| **Refresh Token** | 7日間 | HttpOnly Cookie + DB | Access Token の更新 |

## アーキテクチャ

```
┌─────────────┐         ┌─────────────┐        ┌──────────────┐
│  Frontend   │         │   Backend   │        │  PostgreSQL  │
│  (React)    │         │   (Go/Echo) │        │              │
└──────┬──────┘         └──────┬──────┘        └──────┬───────┘
       │                       │                       │
       │  POST /auth/login     │                       │
       │─────────────────────> │                       │
       │                       │ ① JWT生成             │
       │                       │ ② Refresh Token生成   │
       │                       │───────────────────────>│
       │                       │    (SHA256ハッシュ化)  │
       │ Set-Cookie: access_   │                       │
       │ Set-Cookie: refresh_  │                       │
       │ <───────────────────  │                       │
       │                       │                       │
       │  GET /posts (Cookie)  │                       │
       │─────────────────────> │                       │
       │                       │ ③ Cookie検証          │
       │                       │                       │
       │ (1時間後)              │                       │
       │  GET /posts (401)     │                       │
       │ <───────────────────  │                       │
       │                       │                       │
       │ POST /auth/refresh    │                       │
       │─────────────────────> │                       │
       │                       │ ④ Refresh Token検証   │
       │                       │───────────────────────>│
       │                       │ ⑤ 旧トークン無効化     │
       │                       │───────────────────────>│
       │                       │ ⑥ 新トークン発行       │
       │                       │───────────────────────>│
       │ Set-Cookie: new tokens│                       │
       │ <───────────────────  │                       │
       │                       │                       │
       │  GET /posts (retry)   │                       │
       │─────────────────────> │                       │
       │ <───────────────────  │                       │
       │   200 OK              │                       │
```

## バックエンド実装

### 1. データベーススキーマ

#### RefreshTokens テーブル

```sql
CREATE TABLE refresh_tokens (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL UNIQUE,  -- SHA256ハッシュ化
    expires_at TIMESTAMP NOT NULL,
    revoked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_token ON refresh_tokens(token);
CREATE INDEX idx_refresh_tokens_deleted_at ON refresh_tokens(deleted_at);
```

### 2. JWT ペイロード

```go
// 以前: 多くの情報を含む
type JWTClaims struct {
    UserID   uint   `json:"user_id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    RegisteredClaims jwt.RegisteredClaims
}

// 現在: 最小限の情報のみ
type JWTClaims struct {
    UserID uint `json:"user_id"`
    RegisteredClaims jwt.RegisteredClaims
}
```

**理由**: Cookie のサイズ削減、情報漏洩リスク低減

### 3. Cookie 設定

#### backend/internal/utils/cookie.go

```go
// アクセストークンCookie
func SetAccessTokenCookie(c echo.Context, token string) {
    cookie := &http.Cookie{
        Name:     "access_token",
        Value:    token,
        Path:     "/",
        MaxAge:   3600,            // 1時間
        HttpOnly: true,             // XSS対策
        Secure:   isProduction(),   // HTTPS のみ (本番環境)
        SameSite: getSameSite(),    // CSRF対策
    }
    c.SetCookie(cookie)
}

// リフレッシュトークンCookie
func SetRefreshTokenCookie(c echo.Context, token string) {
    cookie := &http.Cookie{
        Name:     "refresh_token",
        Value:    token,
        Path:     "/",
        MaxAge:   7 * 24 * 3600,    // 7日間
        HttpOnly: true,
        Secure:   isProduction(),
        SameSite: getSameSite(),
    }
    c.SetCookie(cookie)
}

// 環境別のSameSite設定
func getSameSite() http.SameSite {
    if isProduction() {
        return http.SameSiteNoneMode  // クロスオリジン対応
    }
    return http.SameSiteLaxMode       // 開発環境
}
```

### 4. リフレッシュトークン管理

#### backend/internal/utils/refresh_token.go

```go
// トークン生成（SHA256でハッシュ化してDB保存）
func GenerateRefreshToken(userID uint) (string, error) {
    tokenBytes := make([]byte, 32)
    rand.Read(tokenBytes)
    tokenString := base64.URLEncoding.EncodeToString(tokenBytes)

    hashedToken := hashToken(tokenString)  // SHA256

    refreshToken := models.RefreshToken{
        UserID:    userID,
        Token:     hashedToken,
        ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
        Revoked:   false,
    }

    database.DB.Create(&refreshToken)
    return tokenString, nil
}

// トークン検証
func ValidateRefreshToken(tokenString string) (*models.RefreshToken, error) {
    hashedToken := hashToken(tokenString)

    var refreshToken models.RefreshToken
    err := database.DB.Where("token = ?", hashedToken).First(&refreshToken).Error
    if err != nil {
        return nil, errors.New("invalid refresh token")
    }

    if !refreshToken.IsValid() {  // 期限切れまたは失効済み
        return nil, errors.New("refresh token is expired or revoked")
    }

    return &refreshToken, nil
}

// トークン失効（単一）
func RevokeRefreshToken(tokenString string) error {
    hashedToken := hashToken(tokenString)
    return database.DB.Model(&models.RefreshToken{}).
        Where("token = ?", hashedToken).
        Update("revoked", true).Error
}

// 全トークン失効（全デバイスログアウト）
func RevokeAllUserTokens(userID uint) error {
    return database.DB.Model(&models.RefreshToken{}).
        Where("user_id = ? AND revoked = ?", userID, false).
        Update("revoked", true).Error
}
```

### 5. API エンドポイント

#### POST /api/v1/auth/login

```go
func Login(c echo.Context) error {
    // ユーザー認証...

    // アクセストークン生成
    accessToken, _ := utils.GenerateAccessToken(user.ID)

    // リフレッシュトークン生成
    refreshToken, _ := utils.GenerateRefreshToken(user.ID)

    // Cookieに設定
    utils.SetAccessTokenCookie(c, accessToken)
    utils.SetRefreshTokenCookie(c, refreshToken)

    // レスポンス（トークンは含めない）
    return utils.SuccessResponse(c, 200, AuthResponse{
        User: user.ToPublicUser(),
    })
}
```

#### POST /api/v1/auth/refresh

```go
func RefreshToken(c echo.Context) error {
    // Cookieからリフレッシュトークンを取得
    refreshToken, err := utils.GetRefreshTokenFromCookie(c)
    if err != nil {
        return utils.ErrorResponse(c, 401, "リフレッシュトークンが見つかりません")
    }

    // トークン検証
    tokenRecord, err := utils.ValidateRefreshToken(refreshToken)
    if err != nil {
        return utils.ErrorResponse(c, 401, "リフレッシュトークンが無効または期限切れです")
    }

    // 古いトークンを無効化（トークンローテーション）
    utils.RevokeRefreshToken(refreshToken)

    // 新しいトークンを生成
    newAccessToken, _ := utils.GenerateAccessToken(tokenRecord.UserID)
    newRefreshToken, _ := utils.GenerateRefreshToken(tokenRecord.UserID)

    // Cookieに設定
    utils.SetAccessTokenCookie(c, newAccessToken)
    utils.SetRefreshTokenCookie(c, newRefreshToken)

    return utils.SuccessResponse(c, 200, map[string]string{
        "message": "トークンをリフレッシュしました",
    })
}
```

#### POST /api/v1/auth/logout

```go
func Logout(c echo.Context) error {
    // Cookieからリフレッシュトークンを取得
    refreshToken, err := utils.GetRefreshTokenFromCookie(c)
    if err == nil {
        // DBのトークンを無効化
        utils.RevokeRefreshToken(refreshToken)
    }

    // Cookieをクリア
    utils.ClearAuthCookies(c)

    return utils.SuccessResponse(c, 200, map[string]string{
        "message": "ログアウトしました",
    })
}
```

#### POST /api/v1/auth/revoke-all

```go
func RevokeAllTokens(c echo.Context) error {
    userID := c.Get("user_id").(uint)

    // ユーザーの全リフレッシュトークンを無効化
    err := utils.RevokeAllUserTokens(userID)
    if err != nil {
        return utils.ErrorResponse(c, 500, "トークンの無効化に失敗しました")
    }

    // Cookieをクリア
    utils.ClearAuthCookies(c)

    return utils.SuccessResponse(c, 200, map[string]string{
        "message": "全デバイスでログアウトしました",
    })
}
```

### 6. JWT ミドルウェア

#### backend/internal/middleware/jwt_middleware.go

```go
func JWTAuth() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            // Cookieからアクセストークンを取得
            tokenString, err := utils.GetAccessTokenFromCookie(c)
            if err != nil {
                return utils.ErrorResponse(c, 401, "認証が必要です")
            }

            // トークンを検証
            token, err := utils.ValidateToken(tokenString)
            if err != nil {
                return utils.ErrorResponse(c, 401, "トークンが無効または期限切れです")
            }

            // user_id を抽出してコンテキストに設定
            userID, _ := utils.ExtractUserID(token)
            c.Set("user_id", userID)

            return next(c)
        }
    }
}
```

### 7. CORS 設定

#### backend/internal/middleware/cors_middleware.go

```go
func CORS() echo.MiddlewareFunc {
    return middleware.CORSWithConfig(middleware.CORSConfig{
        AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173"},
        AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH, echo.OPTIONS},
        AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
        AllowCredentials: true,  // Cookie送信を許可
    })
}
```

## フロントエンド実装

### 1. API クライアント設定

#### frontend/src/api/openapi-client.ts

```typescript
export const apiClient = createClient<paths>({
  baseUrl: BASE_URL,
  credentials: 'include',  // Cookie送信を有効化
});
```

### 2. 401 エラー時の自動リフレッシュ

#### frontend/src/api/openapi-client.ts

```typescript
let isRefreshing = false;
let failedQueue: Array<{
  resolve: (value?: unknown) => void;
  reject: (reason?: unknown) => void;
}> = [];

const processQueue = (error: Error | null) => {
  failedQueue.forEach((prom) => {
    if (error) {
      prom.reject(error);
    } else {
      prom.resolve();
    }
  });
  failedQueue = [];
};

// レスポンスインターセプター
apiClient.use({
  async onResponse({ response, request }) {
    if (response.status === 401) {
      const requestUrl = new URL(request.url);

      // リフレッシュAPI自体のエラーはリダイレクト
      if (requestUrl.pathname.includes('/auth/refresh')) {
        window.location.href = '/login';
        return response;
      }

      // ログイン/登録ページでは何もしない
      const currentPath = window.location.pathname;
      if (currentPath === '/login' || currentPath === '/register') {
        return response;
      }

      // 既にリフレッシュ中の場合はキューに追加
      if (isRefreshing) {
        await new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject });
        });
        return await fetch(request);
      }

      isRefreshing = true;

      try {
        // リフレッシュAPIを呼び出し
        const refreshResponse = await fetch(`${BASE_URL}/auth/refresh`, {
          method: 'POST',
          credentials: 'include',
        });

        if (!refreshResponse.ok) {
          throw new Error('Refresh failed');
        }

        // リフレッシュ成功、キュー内のリクエストを再実行
        processQueue(null);
        isRefreshing = false;

        // 元のリクエストを再実行
        return await fetch(request);
      } catch (refreshError) {
        // リフレッシュ失敗、ログインページへ
        processQueue(refreshError as Error);
        isRefreshing = false;
        window.location.href = '/login';
        return response;
      }
    }

    return response;
  },
});
```

### 3. AuthContext の更新

#### frontend/src/contexts/AuthContext.tsx

```typescript
export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [user, setUserState] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  // 初期化：Cookieからユーザー情報を取得
  useEffect(() => {
    const initAuth = async () => {
      try {
        // Cookieにトークンがあれば、現在のユーザー情報を取得
        const currentUser = await authApi.getCurrentUser();
        setUserState(currentUser);
      } catch (error) {
        // Cookieにトークンがない、または無効な場合は何もしない
        setUserState(null);
      } finally {
        setIsLoading(false);
      }
    };

    initAuth();
  }, []);

  const login = async (data: LoginRequest): Promise<void> => {
    const user = await authApi.login(data);
    setUserState(user);
    // localStorage への保存は不要（Cookie に保存される）
  };

  const logout = async (): Promise<void> => {
    try {
      await authApi.logout();
    } catch (error) {
      console.error('Logout API error:', error);
    } finally {
      setUserState(null);
      // localStorage.removeItem は不要（Cookie はサーバーで削除）
    }
  };

  // ...
};
```

### 4. 認証 API の更新

#### frontend/src/api/auth.ts

```typescript
// 新規登録
export const register = async (data: RegisterRequest): Promise<User> => {
  const { data: responseData, error } = await apiClient.POST('/auth/register', {
    body: data,
  });

  if (error) {
    const apiError: any = new Error('Registration failed');
    apiError.response = { data: error };
    throw apiError;
  }

  const authResponse = (responseData as unknown as BackendAuthResponse).data;

  // トークンはCookieに保存されるため、ここでは何もしない
  return authResponse.user;
};

// ログアウト
export const logout = async (): Promise<void> => {
  const { error } = await apiClient.POST('/auth/logout', {});

  if (error) {
    console.error('Logout failed:', error);
    // ログアウトは失敗してもCookieをクリアする（サーバー側で削除される）
  }
};

// 全デバイスログアウト
export const revokeAllTokens = async (): Promise<void> => {
  const { error } = await apiClient.POST('/auth/revoke-all', {});

  if (error) {
    throw new Error('Failed to revoke all tokens');
  }
};
```

## テスト

### バックエンドユニットテスト

#### backend/internal/utils/refresh_token_test.go

実装済みのテスト:

- ✅ リフレッシュトークン生成
- ✅ リフレッシュトークン検証
- ✅ 無効なトークンの検証エラー
- ✅ リフレッシュトークン失効
- ✅ 全トークン失効（全デバイスログアウト）
- ✅ トークンローテーション

#### backend/internal/utils/cookie_test.go

実装済みのテスト:

- ✅ アクセストークンCookie設定
- ✅ リフレッシュトークンCookie設定
- ✅ 認証Cookie削除
- ✅ CookieからアクセストークンX取得
- ✅ Cookieからリフレッシュトークンを取得
- ✅ Cookie不在時のエラー処理
- ✅ 開発環境でのSecure属性（false）
- ✅ 本番環境でのSecure属性（true）

### テスト実行方法

#### バックエンド単体テスト

```bash
# make コマンド（推奨）
make test-backend

# または直接実行
docker compose --profile test up -d db_test api_test
docker compose exec -T api_test go test -v ./...
docker compose stop api_test db_test
```

#### E2Eテスト

**重要**: E2Eテストは**テスト用データベース** (`sns_db_test`) を使用します。

```bash
# make コマンド（推奨）
make test-e2e

# または直接実行
docker compose --profile test up -d db_test api_test
cd frontend && npm run test:e2e
cd .. && docker compose stop api_test db_test
```

**E2Eテスト環境の仕組み:**

1. **テスト用Dockerサービス** (`--profile test`)
   - `db_test`: PostgreSQL (ポート5433, データベース: `sns_db_test`)
   - `api_test`: Go APIサーバー (ポート8081, ENV=test)

2. **Playwrightの設定** (`frontend/playwright.config.ts`)
   ```typescript
   webServer: {
     command: process.env.CI ? 'npm run dev' : 'npm run dev:test',
     // 'npm run dev:test' は VITE_API_BASE_URL=http://localhost:8081/api/v1 を設定
   }
   ```

3. **環境変数** (`frontend/package.json`)
   ```json
   "dev:test": "VITE_API_BASE_URL=http://localhost:8081/api/v1 vite"
   ```

**データベース分離の確認:**

```bash
# テストDBのユーザー確認
docker compose exec -T db_test psql -U postgres -d sns_db_test -c "SELECT username FROM users ORDER BY created_at DESC LIMIT 5;"

# 開発DBのユーザー確認（E2Eテスト実行後も変化なし）
docker compose exec -T db psql -U postgres -d sns_db -c "SELECT username FROM users ORDER BY created_at DESC LIMIT 5;"
```

### 動作確認済み機能

- ✅ ユーザー登録
- ✅ ログイン
- ✅ ログアウト
- ✅ Cookie ベース認証
- ✅ 401 エラー時の自動リフレッシュ
- ✅ E2Eテスト (23/25 成功)
  - 認証フロー
  - 投稿操作
  - プロフィール管理
  - ソーシャルインタラクション (いいね、フォロー)

## セキュリティ上の利点

### 1. XSS (Cross-Site Scripting) 対策

- **HttpOnly Cookie**: JavaScript からアクセス不可
  - `document.cookie` で読み取れない
  - XSS 攻撃でトークンを盗まれるリスクが大幅に低減

### 2. CSRF (Cross-Site Request Forgery) 対策

- **SameSite 属性**: クロスサイトからのリクエストを制限
  - 開発環境: `SameSite=Lax`
  - 本番環境: `SameSite=None; Secure`（異なるオリジン対応）

### 3. トークン盗難対策

- **Refresh Token のハッシュ化**: DB に保存されるのは SHA256 ハッシュ
  - DB が漏洩しても元のトークンは復元不可
- **トークンローテーション**: リフレッシュ時に古いトークンを無効化
  - トークンの使い回しを防止

### 4. トークン有効期限の最小化

- **Access Token: 1時間**: 短い有効期限で漏洩リスク低減
- **Refresh Token: 7日間**: 長期間ログイン状態を維持

## 環境別設定

| 環境 | Secure | SameSite | HTTPS |
|------|--------|----------|-------|
| **開発** | false | Lax | 不要 |
| **本番** | true | None | 必須 |

## トラブルシューティング

### 問題: Cookie が設定されない

**原因**: CORS設定で `AllowCredentials: true` が設定されていない

**解決策**:
```go
// backend/internal/middleware/cors_middleware.go
AllowCredentials: true,
```

### 問題: リフレッシュが無限ループする

**原因**: リフレッシュAPIへのリクエストが再試行されている

**解決策**: リフレッシュAPIを再試行対象から除外
```typescript
if (requestUrl.pathname.includes('/auth/refresh')) {
    window.location.href = '/login';
    return response;
}
```

### 問題: 本番環境で Cookie が送信されない

**原因**: Secure 属性が true だが HTTPS を使用していない

**解決策**: 本番環境では必ず HTTPS を使用する

## 今後の拡張案

### 1. リフレッシュトークンの期限切れクリーンアップ

定期的に期限切れトークンを DB から削除

```go
func CleanupExpiredTokens() error {
    return database.DB.Where("expires_at < ?", time.Now()).
        Delete(&models.RefreshToken{}).Error
}
```

### 2. デバイス管理機能

ユーザーがログイン中のデバイス一覧を表示・管理

```go
type RefreshToken struct {
    // ...
    DeviceName string `json:"device_name"`
    IPAddress  string `json:"ip_address"`
    UserAgent  string `json:"user_agent"`
}
```

### 3. レート制限

リフレッシュAPIへの過剰なリクエストを制限

```go
// ミドルウェアで実装
middleware.RateLimiter(middleware.RateLimiterConfig{
    Max: 10,               // 10回まで
    Duration: time.Minute, // 1分間
})
```

## 参考資料

- [OWASP - HttpOnly Cookie](https://owasp.org/www-community/HttpOnly)
- [MDN - SameSite cookies](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Set-Cookie/SameSite)
- [RFC 6749 - OAuth 2.0 (Refresh Token)](https://datatracker.ietf.org/doc/html/rfc6749#section-1.5)
- [JWT Best Practices](https://tools.ietf.org/html/rfc8725)

## 変更履歴

| 日付 | バージョン | 変更内容 |
|------|-----------|---------|
| 2026-02-16 | 1.0.0 | 初版作成 - Cookie ベース認証実装完了 |
