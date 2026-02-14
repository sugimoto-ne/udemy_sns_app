package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/yourusername/sns-backend/internal/services"
	"github.com/yourusername/sns-backend/internal/utils"
)

// RegisterRequest - 登録リクエスト
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Username string `json:"username" validate:"required,min=3,max=50"`
}

// LoginRequest - ログインリクエスト
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// AuthResponse - 認証レスポンス
type AuthResponse struct {
	User  interface{} `json:"user"`
	Token string      `json:"token"`
}

// Register - ユーザー登録ハンドラー
// @Summary ユーザー登録
// @Description 新しいユーザーを登録し、JWTトークンを発行します
// @Tags 認証
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "登録情報"
// @Success 201 {object} map[string]interface{} "data: AuthResponse"
// @Failure 400 {object} map[string]interface{} "バリデーションエラー"
// @Failure 409 {object} map[string]interface{} "メールアドレスまたはユーザー名が既に存在"
// @Failure 500 {object} map[string]interface{} "サーバーエラー"
// @Router /auth/register [post]
func Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	// バリデーション
	if err := utils.ValidateStruct(req); err != nil {
		return utils.ErrorResponse(c, 400, err.Error())
	}

	// ユーザー登録
	user, err := services.Register(req.Email, req.Password, req.Username)
	if err != nil {
		if err.Error() == "email already exists" || err.Error() == "username already exists" {
			return utils.ErrorResponse(c, 409, err.Error())
		}
		return utils.ErrorResponse(c, 500, "Failed to register user")
	}

	// JWTトークン生成
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return utils.ErrorResponse(c, 500, "Failed to generate token")
	}

	// レスポンス
	response := AuthResponse{
		User:  user.ToPublicUser(),
		Token: token,
	}

	return utils.SuccessResponse(c, 201, response)
}

// Login - ログインハンドラー
// @Summary ユーザーログイン
// @Description メールアドレスとパスワードでログインし、JWTトークンを発行します
// @Tags 認証
// @Accept json
// @Produce json
// @Param request body LoginRequest true "ログイン情報"
// @Success 200 {object} map[string]interface{} "data: AuthResponse"
// @Failure 400 {object} map[string]interface{} "バリデーションエラー"
// @Failure 401 {object} map[string]interface{} "メールアドレスまたはパスワードが正しくない"
// @Failure 500 {object} map[string]interface{} "サーバーエラー"
// @Router /auth/login [post]
func Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	// バリデーション
	if err := utils.ValidateStruct(req); err != nil {
		return utils.ErrorResponse(c, 400, err.Error())
	}

	// ログイン
	user, err := services.Login(req.Email, req.Password)
	if err != nil {
		if err.Error() == "invalid email or password" {
			return utils.ErrorResponse(c, 401, err.Error())
		}
		return utils.ErrorResponse(c, 500, "Failed to login")
	}

	// JWTトークン生成
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return utils.ErrorResponse(c, 500, "Failed to generate token")
	}

	// レスポンス
	response := AuthResponse{
		User:  user.ToPublicUser(),
		Token: token,
	}

	return utils.SuccessResponse(c, 200, response)
}

// GetMe - 現在のユーザー情報取得ハンドラー
// @Summary 現在のユーザー情報取得
// @Description JWTトークンから現在ログイン中のユーザー情報を取得します
// @Tags 認証
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "data: PublicUser"
// @Failure 401 {object} map[string]interface{} "認証エラー"
// @Failure 404 {object} map[string]interface{} "ユーザーが見つかりません"
// @Failure 500 {object} map[string]interface{} "サーバーエラー"
// @Router /auth/me [get]
func GetMe(c echo.Context) error {
	// ミドルウェアで設定されたユーザーIDを取得
	userID := c.Get("user_id").(uint)

	// ユーザー情報取得
	user, err := services.GetCurrentUser(userID)
	if err != nil {
		if err.Error() == "user not found" {
			return utils.ErrorResponse(c, 404, err.Error())
		}
		return utils.ErrorResponse(c, 500, "Failed to get user")
	}

	return utils.SuccessResponse(c, 200, user.ToPublicUser())
}
