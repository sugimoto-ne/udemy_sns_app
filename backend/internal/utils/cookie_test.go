package utils

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
)

// TestSetAccessTokenCookie - アクセストークンCookie設定テスト
func TestSetAccessTokenCookie(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// テストトークン
	testToken := "test-access-token"

	// Cookieを設定
	SetAccessTokenCookie(c, testToken)

	// レスポンスからCookieを取得
	cookies := rec.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("No cookies set")
	}

	// Cookieの検証
	var accessCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == AccessTokenCookieName {
			accessCookie = cookie
			break
		}
	}

	if accessCookie == nil {
		t.Fatal("Access token cookie not found")
	}

	// Cookie属性の検証
	if accessCookie.Value != testToken {
		t.Errorf("Cookie value mismatch: got %s, want %s", accessCookie.Value, testToken)
	}

	if accessCookie.Path != "/" {
		t.Errorf("Cookie path mismatch: got %s, want /", accessCookie.Path)
	}

	if accessCookie.MaxAge != 3600 {
		t.Errorf("Cookie MaxAge mismatch: got %d, want 3600", accessCookie.MaxAge)
	}

	if !accessCookie.HttpOnly {
		t.Error("Cookie should be HttpOnly")
	}
}

// TestSetRefreshTokenCookie - リフレッシュトークンCookie設定テスト
func TestSetRefreshTokenCookie(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// テストトークン
	testToken := "test-refresh-token"

	// Cookieを設定
	SetRefreshTokenCookie(c, testToken)

	// レスポンスからCookieを取得
	cookies := rec.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("No cookies set")
	}

	// Cookieの検証
	var refreshCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == RefreshTokenCookieName {
			refreshCookie = cookie
			break
		}
	}

	if refreshCookie == nil {
		t.Fatal("Refresh token cookie not found")
	}

	// Cookie属性の検証
	if refreshCookie.Value != testToken {
		t.Errorf("Cookie value mismatch: got %s, want %s", refreshCookie.Value, testToken)
	}

	if refreshCookie.Path != "/" {
		t.Errorf("Cookie path mismatch: got %s, want /", refreshCookie.Path)
	}

	expectedMaxAge := 7 * 24 * 3600 // 7日間
	if refreshCookie.MaxAge != expectedMaxAge {
		t.Errorf("Cookie MaxAge mismatch: got %d, want %d", refreshCookie.MaxAge, expectedMaxAge)
	}

	if !refreshCookie.HttpOnly {
		t.Error("Cookie should be HttpOnly")
	}
}

// TestClearAuthCookies - 認証Cookie削除テスト
func TestClearAuthCookies(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Cookieをクリア
	ClearAuthCookies(c)

	// レスポンスからCookieを取得
	cookies := rec.Result().Cookies()
	if len(cookies) != 2 {
		t.Errorf("Expected 2 cookies, got %d", len(cookies))
	}

	// 両方のCookieが削除用（MaxAge = -1）であることを確認
	for _, cookie := range cookies {
		if cookie.MaxAge != -1 {
			t.Errorf("Cookie %s should have MaxAge = -1 for deletion, got %d", cookie.Name, cookie.MaxAge)
		}

		if cookie.Value != "" {
			t.Errorf("Cookie %s should have empty value, got %s", cookie.Name, cookie.Value)
		}

		if !cookie.HttpOnly {
			t.Errorf("Cookie %s should be HttpOnly", cookie.Name)
		}
	}
}

// TestGetAccessTokenFromCookie - CookieからアクセストークンX取得テスト
func TestGetAccessTokenFromCookie(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	testToken := "test-access-token"

	// Cookieを含むリクエストを作成
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{
		Name:  AccessTokenCookieName,
		Value: testToken,
	})
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Cookieからトークンを取得
	token, err := GetAccessTokenFromCookie(c)
	if err != nil {
		t.Fatalf("GetAccessTokenFromCookie failed: %v", err)
	}

	if token != testToken {
		t.Errorf("Token mismatch: got %s, want %s", token, testToken)
	}
}

// TestGetAccessTokenFromCookie_NotFound - Cookie不在時のエラーテスト
func TestGetAccessTokenFromCookie_NotFound(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Cookieからトークンを取得（存在しない）
	_, err := GetAccessTokenFromCookie(c)
	if err == nil {
		t.Error("Expected error when cookie not found, got nil")
	}
}

// TestGetRefreshTokenFromCookie - Cookieからリフレッシュトークンを取得テスト
func TestGetRefreshTokenFromCookie(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	testToken := "test-refresh-token"

	// Cookieを含むリクエストを作成
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{
		Name:  RefreshTokenCookieName,
		Value: testToken,
	})
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Cookieからトークンを取得
	token, err := GetRefreshTokenFromCookie(c)
	if err != nil {
		t.Fatalf("GetRefreshTokenFromCookie failed: %v", err)
	}

	if token != testToken {
		t.Errorf("Token mismatch: got %s, want %s", token, testToken)
	}
}

// TestGetRefreshTokenFromCookie_NotFound - Cookie不在時のエラーテスト
func TestGetRefreshTokenFromCookie_NotFound(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Cookieからトークンを取得（存在しない）
	_, err := GetRefreshTokenFromCookie(c)
	if err == nil {
		t.Error("Expected error when cookie not found, got nil")
	}
}

// TestCookieSecureAttribute_Development - 開発環境でSecure属性がfalseであることを確認
func TestCookieSecureAttribute_Development(t *testing.T) {
	// 開発環境を設定
	originalEnv := os.Getenv("APP_ENV")
	os.Setenv("APP_ENV", "development")
	defer os.Setenv("APP_ENV", originalEnv)

	// Echoのセットアップ
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Cookieを設定
	SetAccessTokenCookie(c, "test-token")

	// レスポンスからCookieを取得
	cookies := rec.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("No cookies set")
	}

	// Secure属性がfalseであることを確認
	if cookies[0].Secure {
		t.Error("Cookie should not be Secure in development environment")
	}

	// SameSiteがLaxであることを確認
	if cookies[0].SameSite != http.SameSiteLaxMode {
		t.Errorf("Cookie SameSite should be Lax in development, got %v", cookies[0].SameSite)
	}
}

// TestCookieSecureAttribute_Production - 本番環境でSecure属性がtrueであることを確認
func TestCookieSecureAttribute_Production(t *testing.T) {
	// 本番環境を設定
	originalEnv := os.Getenv("APP_ENV")
	os.Setenv("APP_ENV", "production")
	defer os.Setenv("APP_ENV", originalEnv)

	// Echoのセットアップ
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Cookieを設定
	SetAccessTokenCookie(c, "test-token")

	// レスポンスからCookieを取得
	cookies := rec.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("No cookies set")
	}

	// Secure属性がtrueであることを確認
	if !cookies[0].Secure {
		t.Error("Cookie should be Secure in production environment")
	}

	// SameSiteがNoneであることを確認
	if cookies[0].SameSite != http.SameSiteNoneMode {
		t.Errorf("Cookie SameSite should be None in production, got %v", cookies[0].SameSite)
	}
}
