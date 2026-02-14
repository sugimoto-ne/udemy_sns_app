package middleware

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

// ErrorHandler - カスタムエラーハンドラー
func ErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	message := "Internal Server Error"

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		message = he.Message.(string)
	}

	// ログ出力
	log.Printf("Error: %v", err)

	// レスポンス送信（既に送信済みでない場合）
	if !c.Response().Committed {
		if c.Request().Method == echo.HEAD {
			c.NoContent(code)
		} else {
			c.JSON(code, map[string]interface{}{
				"error": map[string]interface{}{
					"message": message,
				},
			})
		}
	}
}
