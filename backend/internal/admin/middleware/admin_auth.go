package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yourusername/sns-backend/internal/database"
	"github.com/yourusername/sns-backend/internal/models"
	"github.com/yourusername/sns-backend/internal/utils"
)

// AdminJWTAuth - 管理者専用JWT認証ミドルウェア（admin_tokenを使用）
func AdminJWTAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Cookieから管理者トークンを取得
			tokenString, err := utils.GetAdminTokenFromCookie(c)
			if err != nil {
				// 管理者トークンがない場合、ログインページへリダイレクト
				return c.Redirect(http.StatusSeeOther, "/admin/login")
			}

			// トークンを検証
			token, err := utils.ValidateToken(tokenString)
			if err != nil {
				// トークンが無効な場合、ログインページへリダイレクト
				return c.Redirect(http.StatusSeeOther, "/admin/login")
			}

			// トークンからユーザーIDを取得
			userID, err := utils.ExtractUserID(token)
			if err != nil {
				return c.Redirect(http.StatusSeeOther, "/admin/login")
			}

			// データベースからユーザーを取得
			db := database.GetDB()
			var user models.User
			if err := db.First(&user, userID).Error; err != nil {
				return c.Redirect(http.StatusSeeOther, "/admin/login")
			}

			// adminロールかチェック
			if user.Role != "admin" {
				return echo.NewHTTPError(http.StatusForbidden, "Admin access required")
			}

			// コンテキストにユーザーIDと管理者情報を設定
			c.Set("user_id", userID)
			c.Set("admin_user", user)
			return next(c)
		}
	}
}

// AdminRoleCheck - 管理者ロールチェックミドルウェア（非推奨: AdminJWTAuthを使用してください）
// JWT認証ミドルウェアの後に適用すること
func AdminRoleCheck() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// JWT認証ミドルウェアで設定されたuser_idを取得
			userID := c.Get("user_id")
			if userID == nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
			}

			// データベースからユーザーを取得
			db := database.GetDB()
			var user models.User
			if err := db.First(&user, userID).Error; err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "User not found")
			}

			// adminロールかチェック
			if user.Role != "admin" {
				return echo.NewHTTPError(http.StatusForbidden, "Admin access required")
			}

			// 管理者情報をコンテキストに保存
			c.Set("admin_user", user)
			return next(c)
		}
	}
}
