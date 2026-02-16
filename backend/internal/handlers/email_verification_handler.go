package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yourusername/sns-backend/internal/services"
	"github.com/yourusername/sns-backend/internal/utils"
)

// EmailVerificationHandler メール認証ハンドラー
type EmailVerificationHandler struct {
	emailVerificationService *services.EmailVerificationService
}

// NewEmailVerificationHandler EmailVerificationHandlerのコンストラクタ
func NewEmailVerificationHandler() *EmailVerificationHandler {
	return &EmailVerificationHandler{
		emailVerificationService: services.NewEmailVerificationService(),
	}
}

// EmailVerifyRequest メール認証リクエスト
type EmailVerifyRequest struct {
	Token string `json:"token" validate:"required"`
}

// VerifyEmail メールアドレスを認証
// @Summary メールアドレス認証
// @Description トークンを検証してメールアドレスを認証済みにする
// @Tags auth
// @Accept json
// @Produce json
// @Param request body EmailVerifyRequest true "認証トークン"
// @Success 200 {object} map[string]interface{} "Success"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /auth/email/verify [post]
func (h *EmailVerificationHandler) VerifyEmail(c echo.Context) error {
	var req EmailVerifyRequest
	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	// バリデーション
	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	// メール認証処理
	if err := h.emailVerificationService.VerifyEmail(c.Request().Context(), req.Token); err != nil {
		if err.Error() == "invalid or expired token" || err.Error() == "token has expired" {
			return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid or expired token")
		}
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to verify email")
	}

	return utils.SuccessResponse(c, http.StatusOK, map[string]interface{}{
		"message": "Email verified successfully",
	})
}

// ResendVerificationEmail 認証メールを再送信
// @Summary 認証メール再送信
// @Description 認証メールを再送信（既存トークンは削除される）
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Success"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /auth/email/resend [post]
func (h *EmailVerificationHandler) ResendVerificationEmail(c echo.Context) error {
	// 認証済みユーザーID取得
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
	}

	// 認証メール再送信
	if err := h.emailVerificationService.ResendVerificationEmail(c.Request().Context(), userID); err != nil {
		if err.Error() == "email already verified" {
			return utils.ErrorResponse(c, http.StatusBadRequest, "Email is already verified")
		}
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to resend verification email")
	}

	return utils.SuccessResponse(c, http.StatusOK, map[string]interface{}{
		"message": "Verification email sent successfully",
	})
}
