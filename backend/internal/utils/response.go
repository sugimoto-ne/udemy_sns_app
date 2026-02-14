package utils

import (
	"github.com/labstack/echo/v4"
)

// SuccessResponse - 成功レスポンス
func SuccessResponse(c echo.Context, statusCode int, data interface{}) error {
	return c.JSON(statusCode, map[string]interface{}{
		"data": data,
	})
}

// ErrorResponse - エラーレスポンス
func ErrorResponse(c echo.Context, statusCode int, message string) error {
	return c.JSON(statusCode, map[string]interface{}{
		"error": map[string]interface{}{
			"message": message,
		},
	})
}

// PaginationResponse - ページネーション付きレスポンス
func PaginationResponse(c echo.Context, data interface{}, hasMore bool, nextCursor string, limit int) error {
	return c.JSON(200, map[string]interface{}{
		"data": data,
		"pagination": map[string]interface{}{
			"has_more":    hasMore,
			"next_cursor": nextCursor,
			"limit":       limit,
		},
	})
}
