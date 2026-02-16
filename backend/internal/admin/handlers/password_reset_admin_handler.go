package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/yourusername/sns-backend/internal/admin/utils"
	"github.com/yourusername/sns-backend/internal/database"
	"github.com/yourusername/sns-backend/internal/models"
)

type PasswordResetAdminHandler struct{}

func NewPasswordResetAdminHandler() *PasswordResetAdminHandler {
	return &PasswordResetAdminHandler{}
}

// ShowPasswordResetList - パスワードリセット申請一覧画面表示
func (h *PasswordResetAdminHandler) ShowPasswordResetList(c echo.Context) error {
	adminUser := c.Get("admin_user").(models.User)

	return c.Render(http.StatusOK, "password_resets/index.html", map[string]interface{}{
		"Title":         "パスワードリセット申請",
		"AdminUsername": adminUser.Username,
		"Active":        "password-resets",
		"Breadcrumbs": []map[string]interface{}{
			{"Name": "ダッシュボード", "URL": "/admin/dashboard", "Active": false},
			{"Name": "パスワードリセット", "URL": "/admin/password-resets", "Active": true},
		},
	})
}

// GetPasswordResets - パスワードリセット申請一覧取得API
func (h *PasswordResetAdminHandler) GetPasswordResets(c echo.Context) error {
	db := database.GetDB()

	// クエリパラメータ
	status := c.QueryParam("status")
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	// クエリ構築
	query := db.Model(&models.PasswordResetRequest{}).Preload("User")

	if status != "" && status != "all" {
		query = query.Where("status = ?", status)
	}

	// 総件数
	var total int64
	query.Count(&total)

	// ページネーション
	offset := (page - 1) * limit
	var requests []models.PasswordResetRequest
	query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&requests)

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": map[string]interface{}{
			"requests": requests,
			"pagination": map[string]interface{}{
				"total":       total,
				"page":        page,
				"limit":       limit,
				"total_pages": totalPages,
			},
		},
	})
}

// ApproveResetRequest - パスワードリセット承認（トークン発行）
func (h *PasswordResetAdminHandler) ApproveResetRequest(c echo.Context) error {
	db := database.GetDB()
	adminUser := c.Get("admin_user").(models.User)
	requestID := c.Param("id")

	var request models.PasswordResetRequest
	if err := db.Preload("User").First(&request, requestID).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Request not found")
	}

	if request.Status != "pending" {
		return echo.NewHTTPError(http.StatusBadRequest, "Request already processed")
	}

	// トークン生成
	token := uuid.New().String()
	expiresAt := time.Now().Add(24 * time.Hour)

	// レコード更新
	request.Token = token
	request.Status = "approved"
	request.AdminApprovedBy = &adminUser.ID
	now := time.Now()
	request.AdminApprovedAt = &now
	request.ExpiresAt = expiresAt
	db.Save(&request)

	// リセットURL生成
	frontendURL := os.Getenv("FRONTEND_URL")
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", frontendURL, token)

	// メールテンプレート生成
	emailTemplate := fmt.Sprintf(`こんにちは、%sさん

パスワードリセットのリクエストを承認しました。
以下のリンクからパスワードを再設定してください（24時間有効）：

%s

このリクエストに心当たりがない場合は、このメールを無視してください。

よろしくお願いいたします。
SNS運営チーム`, request.User.Username, resetURL)

	// 管理操作ログ記録
	utils.LogAdminAction(db, utils.AdminLogParams{
		AdminID:        adminUser.ID,
		AdminUsername:  adminUser.Username,
		Action:         "password_reset_approve",
		TargetUserID:   &request.UserID,
		TargetUsername: &request.User.Username,
		Details:        fmt.Sprintf("Reset request ID: %d, Token: %s (truncated), Expires: %s", request.ID, token[:8]+"****", expiresAt.Format(time.RFC3339)),
		IP:             c.RealIP(),
	})

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": map[string]interface{}{
			"token":          token,
			"reset_url":      resetURL,
			"expires_at":     expiresAt,
			"email_template": emailTemplate,
			"user_email":     request.User.Email,
		},
		"message": "Password reset approved",
	})
}
