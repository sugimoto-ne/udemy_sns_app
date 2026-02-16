package utils

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

const (
	AccessTokenCookieName  = "access_token"
	RefreshTokenCookieName = "refresh_token"
	AdminTokenCookieName   = "admin_token" // 管理者専用
)

// SetAccessTokenCookie - アクセストークンをCookieに設定
func SetAccessTokenCookie(c echo.Context, token string) {
	cookie := &http.Cookie{
		Name:     AccessTokenCookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   3600, // 1時間（秒単位）
		HttpOnly: true,
		Secure:   isProduction(),
		SameSite: getSameSite(),
	}
	c.SetCookie(cookie)
}

// SetRefreshTokenCookie - リフレッシュトークンをCookieに設定
func SetRefreshTokenCookie(c echo.Context, token string) {
	cookie := &http.Cookie{
		Name:     RefreshTokenCookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   7 * 24 * 3600, // 7日間（秒単位）
		HttpOnly: true,
		Secure:   isProduction(),
		SameSite: getSameSite(),
	}
	c.SetCookie(cookie)
}

// ClearAuthCookies - 認証関連のCookieをクリア
func ClearAuthCookies(c echo.Context) {
	// アクセストークンを削除
	accessCookie := &http.Cookie{
		Name:     AccessTokenCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // 削除
		HttpOnly: true,
		Secure:   isProduction(),
		SameSite: getSameSite(),
	}
	c.SetCookie(accessCookie)

	// リフレッシュトークンを削除
	refreshCookie := &http.Cookie{
		Name:     RefreshTokenCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // 削除
		HttpOnly: true,
		Secure:   isProduction(),
		SameSite: getSameSite(),
	}
	c.SetCookie(refreshCookie)
}

// GetAccessTokenFromCookie - CookieからアクセストークンX取得
func GetAccessTokenFromCookie(c echo.Context) (string, error) {
	cookie, err := c.Cookie(AccessTokenCookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// GetRefreshTokenFromCookie - Cookieからリフレッシュトークンを取得
func GetRefreshTokenFromCookie(c echo.Context) (string, error) {
	cookie, err := c.Cookie(RefreshTokenCookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// isProduction - 本番環境かどうかを判定
func isProduction() bool {
	env := os.Getenv("APP_ENV")
	return env == "production"
}

// getSameSite - SameSite属性を取得
func getSameSite() http.SameSite {
	// 本番環境ではSameSite=None（異なるオリジン間でCookieを送信）
	// 開発環境ではSameSite=Lax（同一サイト内で送信）
	if isProduction() {
		return http.SameSiteNoneMode
	}
	return http.SameSiteLaxMode
}

// === 管理者専用のCookie関数 ===

// SetAdminTokenCookie - 管理者トークンをCookieに設定
func SetAdminTokenCookie(c echo.Context, token string) {
	cookie := &http.Cookie{
		Name:     AdminTokenCookieName,
		Value:    token,
		Path:     "/admin", // 管理画面専用のパス
		MaxAge:   3600,     // 1時間（秒単位）
		HttpOnly: true,
		Secure:   isProduction(),
		SameSite: getSameSite(),
	}
	c.SetCookie(cookie)
}

// GetAdminTokenFromCookie - Cookieから管理者トークンを取得
func GetAdminTokenFromCookie(c echo.Context) (string, error) {
	cookie, err := c.Cookie(AdminTokenCookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// ClearAdminCookie - 管理者Cookieをクリア
func ClearAdminCookie(c echo.Context) {
	cookie := &http.Cookie{
		Name:     AdminTokenCookieName,
		Value:    "",
		Path:     "/admin",
		MaxAge:   -1, // 削除
		HttpOnly: true,
		Secure:   isProduction(),
		SameSite: getSameSite(),
	}
	c.SetCookie(cookie)
}
