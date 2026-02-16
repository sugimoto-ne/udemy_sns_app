package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yourusername/sns-backend/internal/logger"
)

// AccessLog - アクセスログを出力するミドルウェア
func AccessLog() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// リクエスト開始時刻
			start := time.Now()

			// リクエストIDを取得
			requestID := GetRequestID(c)

			// ユーザーIDを取得（ログイン済みの場合）
			var userID uint
			if id, ok := c.Get("user_id").(uint); ok {
				userID = id
			}

			// 次のハンドラーを実行
			err := next(c)

			// レスポンスタイム計算
			latency := time.Since(start)

			// アクセスログを出力
			log := logger.GetLogger()
			logEvent := log.Info().
				Str("request_id", requestID).
				Str("method", c.Request().Method).
				Str("path", c.Request().URL.Path).
				Int("status", c.Response().Status).
				Dur("latency_ms", latency).
				Str("latency", latency.String()).
				Str("remote_ip", c.RealIP())

			// ユーザーIDがあれば追加
			if userID > 0 {
				logEvent = logEvent.Uint("user_id", userID)
			}

			// クエリパラメータ
			if c.Request().URL.RawQuery != "" {
				logEvent = logEvent.Str("query", c.Request().URL.RawQuery)
			}

			// User-Agent
			if ua := c.Request().UserAgent(); ua != "" {
				logEvent = logEvent.Str("user_agent", ua)
			}

			logEvent.Msg("access")

			return err
		}
	}
}
