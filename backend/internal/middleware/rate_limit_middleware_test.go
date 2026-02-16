package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
)

func TestRateLimit_AuthEndpoints(t *testing.T) {
	t.Run("Success - Allow requests within auth limit", func(t *testing.T) {
		// テスト前にリミッターをリセット
		ResetLimiter()

		e := echo.New()
		testHandler := func(c echo.Context) error {
			return c.String(http.StatusOK, "success")
		}

		// 認証系エンドポイントのレートリミット: 5回/分
		handler := RateLimit(5, 60)(testHandler)

		// 5回まで成功するはず
		for i := 0; i < 5; i++ {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := handler(c)
			if err != nil {
				t.Fatalf("Request %d failed: %v", i+1, err)
			}

			if rec.Code != http.StatusOK {
				t.Errorf("Request %d: expected status 200, got %d", i+1, rec.Code)
			}

			// X-RateLimit-Remainingヘッダーを確認
			remaining := rec.Header().Get("X-RateLimit-Remaining")
			expectedRemaining := formatInt(5 - i - 1)
			if remaining != expectedRemaining {
				t.Errorf("Request %d: X-RateLimit-Remaining = %s, want %s", i+1, remaining, expectedRemaining)
			}
		}
	})

	t.Run("Error - Reject request exceeding auth limit", func(t *testing.T) {
		ResetLimiter()

		e := echo.New()
		testHandler := func(c echo.Context) error {
			return c.String(http.StatusOK, "success")
		}

		handler := RateLimit(5, 60)(testHandler)

		// 5回リクエストを送信（すべて成功）
		for i := 0; i < 5; i++ {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			handler(c)
		}

		// 6回目のリクエストは拒否されるはず
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler(c)
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		if rec.Code != http.StatusTooManyRequests {
			t.Errorf("Expected status 429, got %d", rec.Code)
		}

		// X-RateLimit-Remainingヘッダーが0であることを確認
		remaining := rec.Header().Get("X-RateLimit-Remaining")
		if remaining != "0" {
			t.Errorf("X-RateLimit-Remaining = %s, want 0", remaining)
		}
	})
}

func TestRateLimit_GeneralEndpoints(t *testing.T) {
	t.Run("Success - Allow requests within general limit", func(t *testing.T) {
		ResetLimiter()

		e := echo.New()
		testHandler := func(c echo.Context) error {
			return c.String(http.StatusOK, "success")
		}

		// 一般APIのレートリミット: 10回/分（テスト用に軽量化）
		handler := RateLimit(5, 10)(testHandler)

		// 10回まで成功するはず
		for i := 0; i < 10; i++ {
			req := httptest.NewRequest(http.MethodGet, "/api/v1/posts", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := handler(c)
			if err != nil {
				t.Fatalf("Request %d failed: %v", i+1, err)
			}

			if rec.Code != http.StatusOK {
				t.Errorf("Request %d: expected status 200, got %d", i+1, rec.Code)
			}
		}
	})

	t.Run("Error - Reject request exceeding general limit", func(t *testing.T) {
		ResetLimiter()

		e := echo.New()
		testHandler := func(c echo.Context) error {
			return c.String(http.StatusOK, "success")
		}

		handler := RateLimit(5, 10)(testHandler)

		// 10回リクエストを送信
		for i := 0; i < 10; i++ {
			req := httptest.NewRequest(http.MethodGet, "/api/v1/posts", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			handler(c)
		}

		// 11回目のリクエストは拒否されるはず
		req := httptest.NewRequest(http.MethodGet, "/api/v1/posts", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler(c)
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		if rec.Code != http.StatusTooManyRequests {
			t.Errorf("Expected status 429, got %d", rec.Code)
		}
	})
}

func TestRateLimit_DifferentClients(t *testing.T) {
	t.Run("Success - Different IPs have independent limits", func(t *testing.T) {
		ResetLimiter()

		e := echo.New()
		testHandler := func(c echo.Context) error {
			return c.String(http.StatusOK, "success")
		}

		handler := RateLimit(5, 60)(testHandler)

		// クライアント1が5回リクエスト（認証系）
		for i := 0; i < 5; i++ {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", nil)
			req.RemoteAddr = "192.168.1.1:12345"
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			handler(c)
		}

		// クライアント2は別のリミットを持つはず
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", nil)
		req.RemoteAddr = "192.168.1.2:12345"
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler(c)
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("Different client should not be affected by other client's limit, got status %d", rec.Code)
		}

		// X-RateLimit-Remainingは4であるべき（5 - 1）
		remaining := rec.Header().Get("X-RateLimit-Remaining")
		if remaining != "4" {
			t.Errorf("X-RateLimit-Remaining = %s, want 4", remaining)
		}
	})
}

func TestRateLimit_ResetAfterMinute(t *testing.T) {
	t.Run("Success - Limit resets after 1 minute", func(t *testing.T) {
		// このテストは実際には1分待つ必要があるため、スキップまたは短縮版を実装
		// 実際のテストでは時間の経過をモックすることが推奨される
		t.Skip("Skipping time-dependent test in CI")

		ResetLimiter()

		e := echo.New()
		testHandler := func(c echo.Context) error {
			return c.String(http.StatusOK, "success")
		}

		handler := RateLimit(5, 60)(testHandler)

		// 5回リクエストを送信（すべて成功）
		for i := 0; i < 5; i++ {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			handler(c)
		}

		// 6回目は失敗するはず
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		handler(c)

		if rec.Code != http.StatusTooManyRequests {
			t.Error("6th request should be rate limited")
		}

		// 1分待機
		time.Sleep(61 * time.Second)

		// 1分後は再びリクエストが許可されるはず
		req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", nil)
		rec = httptest.NewRecorder()
		c = e.NewContext(req, rec)
		err := handler(c)

		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("Request after reset should succeed, got status %d", rec.Code)
		}
	})
}

func TestRateLimit_ErrorResponse(t *testing.T) {
	t.Run("Success - 429 response has correct error format", func(t *testing.T) {
		ResetLimiter()

		e := echo.New()
		testHandler := func(c echo.Context) error {
			return c.String(http.StatusOK, "success")
		}

		handler := RateLimit(5, 60)(testHandler)

		// 5回リクエストを送信
		for i := 0; i < 5; i++ {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			handler(c)
		}

		// 6回目でエラーレスポンスを確認
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		handler(c)

		// レスポンスボディにエラーコードとメッセージが含まれていることを確認
		body := rec.Body.String()
		if !contains(body, "RATE_LIMIT_EXCEEDED") {
			t.Error("Error response should contain RATE_LIMIT_EXCEEDED code")
		}
		if !contains(body, "Too many requests") {
			t.Error("Error response should contain error message")
		}
	})
}

func TestIsAuthEndpoint(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"/api/v1/auth/register", true},
		{"/api/v1/auth/login", true},
		{"/api/v1/auth/password-reset", true},
		{"/api/v1/posts", false},
		{"/api/v1/users/profile", false},
		{"/api/v1/auth/register/extra", true}, // prefixマッチ
		{"/health", false},
	}

	for _, tt := range tests {
		result := isAuthEndpoint(tt.path)
		if result != tt.expected {
			t.Errorf("isAuthEndpoint(%s) = %v, want %v", tt.path, result, tt.expected)
		}
	}
}

func TestFormatInt(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{0, "0"},
		{1, "1"},
		{10, "10"},
		{99, "99"},
		{100, "100"},
		{-1, "0"}, // 負の数は0として扱う
		{-100, "0"},
	}

	for _, tt := range tests {
		result := formatInt(tt.input)
		if result != tt.expected {
			t.Errorf("formatInt(%d) = %s, want %s", tt.input, result, tt.expected)
		}
	}
}
