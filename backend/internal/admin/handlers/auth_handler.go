package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yourusername/sns-backend/internal/database"
	"github.com/yourusername/sns-backend/internal/models"
	"github.com/yourusername/sns-backend/internal/utils"
)

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

// ShowLoginPage - ログイン画面表示
func (h *AuthHandler) ShowLoginPage(c echo.Context) error {
	return c.Render(http.StatusOK, "login.html", nil)
}

// Login - 管理者ログイン処理
func (h *AuthHandler) Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	db := database.GetDB()
	var user models.User
	if err := db.Where("username = ? AND role = ?", username, "admin").First(&user).Error; err != nil {
		return c.Render(http.StatusOK, "login.html", map[string]interface{}{
			"Error": "ユーザー名またはパスワードが正しくありません",
		})
	}

	if !user.CheckPassword(password) {
		return c.Render(http.StatusOK, "login.html", map[string]interface{}{
			"Error": "ユーザー名またはパスワードが正しくありません",
		})
	}

	// JWT生成
	token, err := utils.GenerateAccessToken(user.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate token")
	}

	// last_login_at を更新
	now := time.Now()
	user.LastLoginAt = &now
	db.Save(&user)

	// HttpOnly Cookieに保存（管理者専用Cookie）
	utils.SetAdminTokenCookie(c, token)

	return c.Redirect(http.StatusSeeOther, "/admin/dashboard")
}

// Logout - 管理者ログアウト
func (h *AuthHandler) Logout(c echo.Context) error {
	utils.ClearAdminCookie(c)
	return c.Redirect(http.StatusSeeOther, "/admin/login")
}
