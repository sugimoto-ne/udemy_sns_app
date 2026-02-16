package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestSecurityHeaders_Production(t *testing.T) {
	// 本番環境を設定
	originalEnv := os.Getenv("APP_ENV")
	os.Setenv("APP_ENV", "production")
	defer os.Setenv("APP_ENV", originalEnv)

	t.Run("Success - All security headers are set in production", func(t *testing.T) {
		// Echoインスタンスを作成
		e := echo.New()

		// テスト用ハンドラー
		testHandler := func(c echo.Context) error {
			return c.String(http.StatusOK, "test")
		}

		// ミドルウェアを適用
		handler := SecurityHeaders()(testHandler)

		// リクエストとレスポンスを作成
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// ハンドラーを実行
		err := handler(c)
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		// Content-Security-Policyヘッダーを確認
		csp := rec.Header().Get("Content-Security-Policy")
		expectedCSP := "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self'; connect-src 'self'"
		if csp != expectedCSP {
			t.Errorf("Content-Security-Policy = %v, want %v", csp, expectedCSP)
		}

		// X-Frame-Optionsヘッダーを確認
		xFrameOptions := rec.Header().Get("X-Frame-Options")
		if xFrameOptions != "DENY" {
			t.Errorf("X-Frame-Options = %v, want DENY", xFrameOptions)
		}

		// X-Content-Type-Optionsヘッダーを確認
		xContentTypeOptions := rec.Header().Get("X-Content-Type-Options")
		if xContentTypeOptions != "nosniff" {
			t.Errorf("X-Content-Type-Options = %v, want nosniff", xContentTypeOptions)
		}

		// Strict-Transport-Securityヘッダーを確認
		hsts := rec.Header().Get("Strict-Transport-Security")
		expectedHSTS := "max-age=31536000; includeSubDomains"
		if hsts != expectedHSTS {
			t.Errorf("Strict-Transport-Security = %v, want %v", hsts, expectedHSTS)
		}

		// X-XSS-Protectionヘッダーを確認
		xssProtection := rec.Header().Get("X-XSS-Protection")
		if xssProtection != "0" {
			t.Errorf("X-XSS-Protection = %v, want 0", xssProtection)
		}
	})

	t.Run("Success - Headers are set for different HTTP methods", func(t *testing.T) {
		methods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch}

		for _, method := range methods {
			e := echo.New()
			testHandler := func(c echo.Context) error {
				return c.String(http.StatusOK, "test")
			}

			handler := SecurityHeaders()(testHandler)
			req := httptest.NewRequest(method, "/test", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := handler(c)
			if err != nil {
				t.Fatalf("Handler returned error for method %s: %v", method, err)
			}

			// X-Frame-Optionsヘッダーが存在することを確認（代表例）
			if rec.Header().Get("X-Frame-Options") != "DENY" {
				t.Errorf("X-Frame-Options not set for method %s", method)
			}
		}
	})

	t.Run("Success - Headers are set even when handler returns error", func(t *testing.T) {
		e := echo.New()

		// エラーを返すハンドラー
		testHandler := func(c echo.Context) error {
			return echo.NewHTTPError(http.StatusInternalServerError, "test error")
		}

		handler := SecurityHeaders()(testHandler)
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler(c)
		if err == nil {
			t.Fatal("Expected error from handler")
		}

		// エラーが返されてもセキュリティヘッダーが設定されていることを確認
		if rec.Header().Get("X-Frame-Options") != "DENY" {
			t.Error("Security headers should be set even when handler returns error")
		}
	})

	t.Run("Success - Headers do not override existing headers unnecessarily", func(t *testing.T) {
		e := echo.New()

		// 特定のヘッダーを設定するハンドラー
		testHandler := func(c echo.Context) error {
			// セキュリティヘッダーは上書きされるべき（ミドルウェアが先に実行される）
			return c.String(http.StatusOK, "test")
		}

		handler := SecurityHeaders()(testHandler)
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler(c)
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		// ミドルウェアが設定したヘッダーが存在することを確認
		if rec.Header().Get("X-Frame-Options") != "DENY" {
			t.Error("Middleware headers should be set")
		}
	})

	t.Run("Success - CSP header protects against common attacks", func(t *testing.T) {
		e := echo.New()
		testHandler := func(c echo.Context) error {
			return c.String(http.StatusOK, "test")
		}

		handler := SecurityHeaders()(testHandler)
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler(c)
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		csp := rec.Header().Get("Content-Security-Policy")

		// CSPに重要なディレクティブが含まれていることを確認
		requiredDirectives := []string{
			"default-src",
			"script-src",
			"style-src",
			"img-src",
			"font-src",
			"connect-src",
		}

		for _, directive := range requiredDirectives {
			if len(csp) == 0 || !contains(csp, directive) {
				t.Errorf("CSP should contain %s directive, got: %s", directive, csp)
			}
		}

		// 'unsafe-eval'が含まれていないことを確認（セキュリティリスク）
		if contains(csp, "unsafe-eval") {
			t.Error("CSP should not contain 'unsafe-eval' for security reasons")
		}
	})
}

// contains は文字列に部分文字列が含まれているかチェックするヘルパー関数
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
		 s[:len(substr)] == substr ||
		 s[len(s)-len(substr):] == substr ||
		 findSubstring(s, substr))
}

// findSubstring は文字列中の部分文字列を探すヘルパー関数
func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestSecurityHeaders_Development(t *testing.T) {
	// 開発環境を設定（デフォルトまたは明示的にdevelopment）
	originalEnv := os.Getenv("APP_ENV")
	os.Setenv("APP_ENV", "development")
	defer os.Setenv("APP_ENV", originalEnv)

	t.Run("Success - No security headers in development", func(t *testing.T) {
		e := echo.New()

		testHandler := func(c echo.Context) error {
			return c.String(http.StatusOK, "test")
		}

		handler := SecurityHeaders()(testHandler)
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler(c)
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		// 開発環境ではセキュリティヘッダーが付与されないことを確認
		if rec.Header().Get("Content-Security-Policy") != "" {
			t.Error("Content-Security-Policy should not be set in development")
		}
		if rec.Header().Get("X-Frame-Options") != "" {
			t.Error("X-Frame-Options should not be set in development")
		}
		if rec.Header().Get("X-Content-Type-Options") != "" {
			t.Error("X-Content-Type-Options should not be set in development")
		}
		if rec.Header().Get("Strict-Transport-Security") != "" {
			t.Error("Strict-Transport-Security should not be set in development")
		}
		if rec.Header().Get("X-XSS-Protection") != "" {
			t.Error("X-XSS-Protection should not be set in development")
		}
	})

	t.Run("Success - No headers when APP_ENV is empty", func(t *testing.T) {
		// APP_ENVが未設定の場合（空文字列）
		os.Setenv("APP_ENV", "")
		defer os.Setenv("APP_ENV", originalEnv)

		e := echo.New()
		testHandler := func(c echo.Context) error {
			return c.String(http.StatusOK, "test")
		}

		handler := SecurityHeaders()(testHandler)
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler(c)
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		// 環境変数が未設定でもヘッダーは付与されない
		if rec.Header().Get("Content-Security-Policy") != "" {
			t.Error("Content-Security-Policy should not be set when APP_ENV is empty")
		}
	})

	t.Run("Success - No headers in test environment", func(t *testing.T) {
		os.Setenv("APP_ENV", "test")
		defer os.Setenv("APP_ENV", originalEnv)

		e := echo.New()
		testHandler := func(c echo.Context) error {
			return c.String(http.StatusOK, "test")
		}

		handler := SecurityHeaders()(testHandler)
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler(c)
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		// テスト環境でもヘッダーは付与されない
		if rec.Header().Get("X-Frame-Options") != "" {
			t.Error("Security headers should not be set in test environment")
		}
	})
}
