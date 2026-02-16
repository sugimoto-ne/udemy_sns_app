package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yourusername/sns-backend/internal/logger"
)

// ErrorHandler - カスタムエラーハンドラー
func ErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	message := "Internal Server Error"
	requestID := GetRequestID(c)

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		if msg, ok := he.Message.(string); ok {
			message = msg
		}
	}

	// 構造化ログでエラーを記録（スタックトレース付き）
	log := logger.GetLogger()
	log.Error().
		Err(err).
		Str("request_id", requestID).
		Str("method", c.Request().Method).
		Str("path", c.Request().URL.Path).
		Int("status", code).
		Str("remote_ip", c.RealIP()).
		Msg("error occurred")

	// レスポンス送信（既に送信済みでない場合）
	if !c.Response().Committed {
		if c.Request().Method == echo.HEAD {
			c.NoContent(code)
		} else {
			// エラーレスポンスにリクエストIDを含める
			c.JSON(code, map[string]interface{}{
				"error": map[string]interface{}{
					"message":    message,
					"request_id": requestID,
				},
			})
		}
	}
}
