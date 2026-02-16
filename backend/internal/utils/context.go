package utils

import (
	"errors"

	"github.com/labstack/echo/v4"
)

// GetUserIDFromContext - コンテキストから安全にユーザーIDを取得
func GetUserIDFromContext(c echo.Context) (uint, error) {
	userIDInterface := c.Get("user_id")
	if userIDInterface == nil {
		return 0, errors.New("user context not found")
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		return 0, errors.New("invalid user context type")
	}

	return userID, nil
}
