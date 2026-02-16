package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// CORS - CORS設定ミドルウェア
func CORS() echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			// 開発環境
			"http://localhost:3000", // React開発サーバー（旧）
			"http://localhost:5173", // Vite開発サーバー
			"http://localhost:5174", // Vite E2Eテストサーバー
			// 本番環境（Firebase Hosting）
			"https://udemy-sns-b9e40.web.app",
			"https://udemy-sns-b9e40.firebaseapp.com",
		},
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH, echo.OPTIONS},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true, // Cookie送信を許可
	})
}
