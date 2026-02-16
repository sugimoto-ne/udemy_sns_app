package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yourusername/sns-backend/internal/services"
	"github.com/yourusername/sns-backend/internal/utils"
)

// PasswordResetHandler パスワードリセットハンドラー
type PasswordResetHandler struct {
	passwordResetService *services.PasswordResetService
}

// NewPasswordResetHandler PasswordResetHandlerのコンストラクタ
func NewPasswordResetHandler() *PasswordResetHandler {
	return &PasswordResetHandler{
		passwordResetService: services.NewPasswordResetService(),
	}
}

// PasswordResetRequest パスワードリセットリクエスト
type PasswordResetRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// PasswordResetConfirm パスワードリセット確認リクエスト
type PasswordResetConfirm struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

// RequestPasswordReset パスワードリセットをリクエスト
// @Summary パスワードリセットリクエスト
// @Description パスワードリセット用のメールを送信
// @Tags auth
// @Accept json
// @Produce json
// @Param request body PasswordResetRequest true "メールアドレス"
// @Success 200 {object} map[string]interface{} "Success"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /auth/password-reset/request [post]
func (h *PasswordResetHandler) RequestPasswordReset(c echo.Context) error {
	var req PasswordResetRequest
	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	// バリデーション
	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	// パスワードリセットリクエスト処理
	if err := h.passwordResetService.RequestPasswordReset(c.Request().Context(), req.Email); err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to process password reset request")
	}

	// セキュリティ上、常に成功レスポンスを返す（メールアドレスの存在を露出しない）
	return utils.SuccessResponse(c, http.StatusOK, map[string]interface{}{
		"message": "If an account with that email exists, a password reset email has been sent",
	})
}

// ConfirmPasswordReset パスワードリセット確認
// @Summary パスワードリセット確認
// @Description トークンを検証して新しいパスワードに更新
// @Tags auth
// @Accept json
// @Produce json
// @Param request body PasswordResetConfirm true "トークンと新しいパスワード"
// @Success 200 {object} map[string]interface{} "Success"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /auth/password-reset/confirm [post]
func (h *PasswordResetHandler) ConfirmPasswordReset(c echo.Context) error {
	var req PasswordResetConfirm
	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	// バリデーション
	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	// パスワードリセット確認処理
	if err := h.passwordResetService.ConfirmPasswordReset(c.Request().Context(), req.Token, req.NewPassword); err != nil {
		if err.Error() == "invalid or expired token" || err.Error() == "token has expired" {
			return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid or expired token")
		}
		if err.Error() == "password must be at least 8 characters" {
			return utils.ErrorResponse(c, http.StatusBadRequest, "Password must be at least 8 characters")
		}
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to reset password")
	}

	return utils.SuccessResponse(c, http.StatusOK, map[string]interface{}{
		"message": "Password reset successfully",
	})
}
