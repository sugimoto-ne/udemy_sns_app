package middleware

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/yourusername/sns-backend/internal/utils"
)

// JWTAuth - JWT認証ミドルウェア
func JWTAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Authorizationヘッダーからトークンを取得
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return utils.ErrorResponse(c, 401, "Authorization header is required")
			}

			// "Bearer "プレフィックスを削除
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				return utils.ErrorResponse(c, 401, "Invalid authorization header format")
			}

			// トークンを検証
			token, err := utils.ValidateToken(tokenString)
			if err != nil {
				return utils.ErrorResponse(c, 401, "Invalid or expired token")
			}

			// トークンからユーザーIDを取得
			userID, err := utils.ExtractUserID(token)
			if err != nil {
				return utils.ErrorResponse(c, 401, "Invalid token claims")
			}

			// コンテキストにユーザーIDを設定
			c.Set("user_id", userID)

			return next(c)
		}
	}
}

// OptionalJWTAuth - 任意のJWT認証ミドルウェア（トークンがあれば検証、なければスキップ）
func OptionalJWTAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				// トークンがない場合はそのまま次へ
				return next(c)
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				// フォーマットが不正でもエラーにせずスキップ
				return next(c)
			}

			token, err := utils.ValidateToken(tokenString)
			if err == nil {
				// トークンが有効な場合のみユーザーIDを設定
				userID, err := utils.ExtractUserID(token)
				if err == nil {
					c.Set("user_id", userID)
				}
			}

			return next(c)
		}
	}
}
