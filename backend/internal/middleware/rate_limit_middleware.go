package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

// RateLimiter - レートリミット管理構造体
type RateLimiter struct {
	clients map[string]*ClientLimitInfo
	mu      sync.RWMutex
}

// ClientLimitInfo - クライアントごとのリミット情報
type ClientLimitInfo struct {
	count       int
	resetTime   time.Time
	lastCleanup time.Time
}

var (
	limiter = &RateLimiter{
		clients: make(map[string]*ClientLimitInfo),
	}
)

// RateLimit - レートリミットミドルウェア
// authLimit: 認証系APIのリクエスト上限（回/分）
// generalLimit: 一般APIのリクエスト上限（回/分）
func RateLimit(authLimit, generalLimit int) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// クライアントIPを取得
			clientIP := c.RealIP()
			if clientIP == "" {
				clientIP = c.Request().RemoteAddr
			}

			// パスに基づいてリミットを決定
			path := c.Request().URL.Path
			var limit int
			if isAuthEndpoint(path) {
				limit = authLimit
			} else {
				limit = generalLimit
			}

			// レートリミットをチェック
			allowed, remaining := limiter.checkLimit(clientIP, limit)

			// X-RateLimit-Remainingヘッダーを設定
			c.Response().Header().Set("X-RateLimit-Remaining", formatInt(remaining))

			if !allowed {
				return c.JSON(http.StatusTooManyRequests, map[string]interface{}{
					"error": map[string]interface{}{
						"code":    "RATE_LIMIT_EXCEEDED",
						"message": "Too many requests. Please try again later.",
					},
				})
			}

			return next(c)
		}
	}
}

// checkLimit - レートリミットをチェック
func (rl *RateLimiter) checkLimit(clientIP string, limit int) (bool, int) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	// クライアント情報を取得または作成
	info, exists := rl.clients[clientIP]
	if !exists || now.After(info.resetTime) {
		// 新規クライアントまたはリセット時刻を過ぎた場合
		rl.clients[clientIP] = &ClientLimitInfo{
			count:       1,
			resetTime:   now.Add(1 * time.Minute),
			lastCleanup: now,
		}
		return true, limit - 1
	}

	// リミットをチェック
	if info.count >= limit {
		return false, 0
	}

	// カウントを増加
	info.count++

	// 定期的に古いエントリをクリーンアップ（1分ごと）
	if now.Sub(info.lastCleanup) > 1*time.Minute {
		rl.cleanup(now)
		info.lastCleanup = now
	}

	return true, limit - info.count
}

// cleanup - 古いクライアント情報を削除
func (rl *RateLimiter) cleanup(now time.Time) {
	for ip, info := range rl.clients {
		if now.After(info.resetTime) {
			delete(rl.clients, ip)
		}
	}
}

// ResetLimiter - レートリミッターをリセット（テスト用）
func ResetLimiter() {
	limiter.mu.Lock()
	defer limiter.mu.Unlock()
	limiter.clients = make(map[string]*ClientLimitInfo)
}

// isAuthEndpoint - 認証系エンドポイントかどうかを判定
func isAuthEndpoint(path string) bool {
	authPaths := []string{
		"/api/v1/auth/register",
		"/api/v1/auth/login",
		"/api/v1/auth/password-reset", // パスワードリセット
	}

	for _, authPath := range authPaths {
		if strings.HasPrefix(path, authPath) {
			return true
		}
	}
	return false
}

// formatInt - int を文字列に変換
func formatInt(n int) string {
	if n < 0 {
		return "0"
	}
	// 簡易的な変換（標準ライブラリを使わない）
	if n == 0 {
		return "0"
	}

	digits := []byte{}
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}
