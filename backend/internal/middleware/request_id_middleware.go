package middleware

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const RequestIDHeader = "X-Request-ID"
const RequestIDKey = "request_id"

// RequestID - リクエストIDを生成・設定するミドルウェア
func RequestID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// リクエストIDを生成
			requestID := uuid.New().String()

			// コンテキストに設定
			c.Set(RequestIDKey, requestID)

			// レスポンスヘッダーに設定
			c.Response().Header().Set(RequestIDHeader, requestID)

			return next(c)
		}
	}
}

// GetRequestID - コンテキストからリクエストIDを取得
func GetRequestID(c echo.Context) string {
	if requestID, ok := c.Get(RequestIDKey).(string); ok {
		return requestID
	}
	return ""
}
