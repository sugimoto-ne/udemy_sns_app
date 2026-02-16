package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/yourusername/sns-backend/internal/database"
	"github.com/yourusername/sns-backend/internal/models"
	"github.com/yourusername/sns-backend/internal/utils"
)

// JWTAuth - JWT認証ミドルウェア（Cookie対応）
func JWTAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Cookieからアクセストークンを取得
			tokenString, err := utils.GetAccessTokenFromCookie(c)
			if err != nil {
				return utils.ErrorResponse(c, 401, "認証が必要です")
			}

			// トークンを検証
			token, err := utils.ValidateToken(tokenString)
			if err != nil {
				return utils.ErrorResponse(c, 401, "トークンが無効または期限切れです")
			}

			// トークンからユーザーIDを取得
			userID, err := utils.ExtractUserID(token)
			if err != nil {
				return utils.ErrorResponse(c, 401, "トークンのクレームが不正です")
			}

			// データベースからユーザーステータスを確認
			db := database.GetDB()
			var user models.User
			if err := db.Select("status").First(&user, userID).Error; err != nil {
				return utils.ErrorResponse(c, 401, "ユーザーが見つかりません")
			}

			// ステータスチェック（承認済みユーザーのみアクセス可能）
			if user.Status != "approved" {
				// Cookieをクリアしてログアウトさせる
				utils.ClearAuthCookies(c)
				return utils.ErrorResponse(c, 403, "アカウントが承認されていないため、アクセスできません")
			}

			// コンテキストにユーザーIDを設定
			c.Set("user_id", userID)

			return next(c)
		}
	}
}

// OptionalJWTAuth - 任意のJWT認証ミドルウェア（Cookie対応、トークンがあれば検証、なければスキップ）
func OptionalJWTAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Cookieからアクセストークンを取得
			tokenString, err := utils.GetAccessTokenFromCookie(c)
			if err != nil {
				// トークンがない場合はそのまま次へ
				return next(c)
			}

			// トークンを検証
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
