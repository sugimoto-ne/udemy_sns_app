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

// AuthResponse - 認証レスポンス（Cookie使用時はトークンをレスポンスに含めない）
type AuthResponse struct {
	User interface{} `json:"user"`
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
		return utils.ErrorResponse(c, 400, "リクエストの形式が正しくありません")
	}

	// バリデーション
	if err := utils.ValidateStruct(req); err != nil {
		return utils.ErrorResponse(c, 400, err.Error())
	}

	// ユーザー登録
	user, err := services.Register(req.Email, req.Password, req.Username)
	if err != nil {
		if err.Error() == "email already exists" {
			return utils.ErrorResponse(c, 409, "このメールアドレスは既に登録されています")
		}
		if err.Error() == "username already exists" {
			return utils.ErrorResponse(c, 409, "このユーザー名は既に使用されています")
		}
		return utils.ErrorResponse(c, 500, "ユーザー登録に失敗しました")
	}

	// メール送信は廃止: 管理者承認制に変更
	// 管理者が手動で Approved = true に設定することでアカウントを有効化

	// アクセストークン生成
	accessToken, err := utils.GenerateAccessToken(user.ID)
	if err != nil {
		return utils.ErrorResponse(c, 500, "トークンの生成に失敗しました")
	}

	// リフレッシュトークン生成
	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		return utils.ErrorResponse(c, 500, "リフレッシュトークンの生成に失敗しました")
	}

	// Cookieに設定
	utils.SetAccessTokenCookie(c, accessToken)
	utils.SetRefreshTokenCookie(c, refreshToken)

	// レスポンス（トークンはCookieに含まれるため、レスポンスボディには含めない）
	// 本人なのでメールアドレスを含める
	response := AuthResponse{
		User: user.ToPublicUser(&user.ID),
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
		return utils.ErrorResponse(c, 400, "リクエストの形式が正しくありません")
	}

	// バリデーション
	if err := utils.ValidateStruct(req); err != nil {
		return utils.ErrorResponse(c, 400, err.Error())
	}

	// ログイン
	user, err := services.Login(req.Email, req.Password)
	if err != nil {
		if err.Error() == "invalid email or password" {
			return utils.ErrorResponse(c, 401, "メールアドレスまたはパスワードが正しくありません")
		}
		if err.Error() == "account not approved" {
			return utils.ErrorResponse(c, 403, "アカウントは管理者による承認待ちです。承認され次第、ログイン可能になります。")
		}
		return utils.ErrorResponse(c, 500, "ログインに失敗しました")
	}

	// アクセストークン生成
	accessToken, err := utils.GenerateAccessToken(user.ID)
	if err != nil {
		return utils.ErrorResponse(c, 500, "トークンの生成に失敗しました")
	}

	// リフレッシュトークン生成
	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		return utils.ErrorResponse(c, 500, "リフレッシュトークンの生成に失敗しました")
	}

	// Cookieに設定
	utils.SetAccessTokenCookie(c, accessToken)
	utils.SetRefreshTokenCookie(c, refreshToken)

	// レスポンス（トークンはCookieに含まれるため、レスポンスボディには含めない）
	// 本人なのでメールアドレスを含める
	response := AuthResponse{
		User: user.ToPublicUser(&user.ID),
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
	// 安全な型アサーション
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return utils.ErrorResponse(c, 401, "認証エラー")
	}

	// ユーザー情報取得
	user, err := services.GetCurrentUser(userID)
	if err != nil {
		if err.Error() == "user not found" {
			return utils.ErrorResponse(c, 404, "ユーザーが見つかりません")
		}
		return utils.ErrorResponse(c, 500, "ユーザー情報の取得に失敗しました")
	}

	// 本人なのでメールアドレスを含める
	return utils.SuccessResponse(c, 200, user.ToPublicUser(&userID))
}

// RefreshToken - トークンリフレッシュハンドラー
// @Summary トークンリフレッシュ
// @Description リフレッシュトークンを使用してアクセストークンを再発行します（トークンローテーション）
// @Tags 認証
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "message: トークンをリフレッシュしました"
// @Failure 401 {object} map[string]interface{} "リフレッシュトークンが無効または期限切れ"
// @Failure 500 {object} map[string]interface{} "サーバーエラー"
// @Router /auth/refresh [post]
func RefreshToken(c echo.Context) error {
	// Cookieからリフレッシュトークンを取得
	refreshToken, err := utils.GetRefreshTokenFromCookie(c)
	if err != nil {
		return utils.ErrorResponse(c, 401, "リフレッシュトークンが見つかりません")
	}

	// リフレッシュトークンを検証
	tokenRecord, err := utils.ValidateRefreshToken(refreshToken)
	if err != nil {
		return utils.ErrorResponse(c, 401, "リフレッシュトークンが無効または期限切れです")
	}

	// 古いリフレッシュトークンを無効化
	if err := utils.RevokeRefreshToken(refreshToken); err != nil {
		return utils.ErrorResponse(c, 500, "リフレッシュトークンの無効化に失敗しました")
	}

	// 新しいアクセストークンを生成
	newAccessToken, err := utils.GenerateAccessToken(tokenRecord.UserID)
	if err != nil {
		return utils.ErrorResponse(c, 500, "アクセストークンの生成に失敗しました")
	}

	// 新しいリフレッシュトークンを生成（トークンローテーション）
	newRefreshToken, err := utils.GenerateRefreshToken(tokenRecord.UserID)
	if err != nil {
		return utils.ErrorResponse(c, 500, "リフレッシュトークンの生成に失敗しました")
	}

	// Cookieに設定
	utils.SetAccessTokenCookie(c, newAccessToken)
	utils.SetRefreshTokenCookie(c, newRefreshToken)

	return utils.SuccessResponse(c, 200, map[string]string{
		"message": "トークンをリフレッシュしました",
	})
}

// Logout - ログアウトハンドラー
// @Summary ログアウト
// @Description リフレッシュトークンを無効化し、Cookieをクリアします
// @Tags 認証
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "message: ログアウトしました"
// @Failure 500 {object} map[string]interface{} "サーバーエラー"
// @Router /auth/logout [post]
func Logout(c echo.Context) error {
	// Cookieからリフレッシュトークンを取得
	refreshToken, err := utils.GetRefreshTokenFromCookie(c)
	if err == nil {
		// リフレッシュトークンを無効化（エラーは無視）
		_ = utils.RevokeRefreshToken(refreshToken)
	}

	// Cookieをクリア
	utils.ClearAuthCookies(c)

	return utils.SuccessResponse(c, 200, map[string]string{
		"message": "ログアウトしました",
	})
}

// RevokeAllTokens - 全デバイスログアウトハンドラー
// @Summary 全デバイスログアウト
// @Description ユーザーのすべてのリフレッシュトークンを無効化します
// @Tags 認証
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "message: すべてのデバイスからログアウトしました"
// @Failure 401 {object} map[string]interface{} "認証エラー"
// @Failure 500 {object} map[string]interface{} "サーバーエラー"
// @Router /auth/revoke-all [post]
func RevokeAllTokens(c echo.Context) error {
	// 安全な型アサーション
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return utils.ErrorResponse(c, 401, "認証エラー")
	}

	// ユーザーのすべてのリフレッシュトークンを無効化
	if err := utils.RevokeAllUserTokens(userID); err != nil {
		return utils.ErrorResponse(c, 500, "トークンの無効化に失敗しました")
	}

	// Cookieをクリア
	utils.ClearAuthCookies(c)

	return utils.SuccessResponse(c, 200, map[string]string{
		"message": "すべてのデバイスからログアウトしました",
	})
}
