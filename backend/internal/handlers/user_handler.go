package handlers

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/yourusername/sns-backend/internal/services"
	"github.com/yourusername/sns-backend/internal/utils"
)

// UpdateProfileRequest - プロフィール更新リクエスト
type UpdateProfileRequest struct {
	DisplayName *string `json:"display_name"`
	Bio         *string `json:"bio"`
	AvatarURL   *string `json:"avatar_url"`
	HeaderURL   *string `json:"header_url"`
	Website     *string `json:"website"`
	BirthDate   *string `json:"birth_date"`
	Occupation  *string `json:"occupation"`
}

// GetUserByUsername - ユーザー名でユーザーを取得ハンドラー
// @Summary ユーザー情報取得
// @Description ユーザー名でユーザー情報を取得します
// @Tags ユーザー
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param username path string true "ユーザー名"
// @Success 200 {object} map[string]interface{} "data: PublicUser"
// @Failure 404 {object} map[string]interface{} "ユーザーが見つかりません"
// @Failure 500 {object} map[string]interface{} "サーバーエラー"
// @Router /users/{username} [get]
func GetUserByUsername(c echo.Context) error {
	username := c.Param("username")

	// 現在のユーザーID取得（任意）
	var currentUserIDPtr *uint
	if userID, ok := c.Get("user_id").(uint); ok {
		currentUserIDPtr = &userID
	}

	publicUser, err := services.GetUserByUsername(username, currentUserIDPtr)
	if err != nil {
		if err.Error() == "user not found" {
			return utils.ErrorResponse(c, 404, err.Error())
		}
		return utils.ErrorResponse(c, 500, "Failed to get user")
	}

	return utils.SuccessResponse(c, 200, publicUser)
}

// UpdateProfile - プロフィール更新ハンドラー
// @Summary プロフィール更新
// @Description 自分のプロフィール情報を更新します
// @Tags ユーザー
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UpdateProfileRequest true "更新内容"
// @Success 200 {object} map[string]interface{} "data: PublicUser"
// @Failure 400 {object} map[string]interface{} "バリデーションエラー"
// @Failure 401 {object} map[string]interface{} "認証エラー"
// @Failure 404 {object} map[string]interface{} "ユーザーが見つかりません"
// @Failure 500 {object} map[string]interface{} "サーバーエラー"
// @Router /users/me [put]
func UpdateProfile(c echo.Context) error {
	var req UpdateProfileRequest
	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	userID := c.Get("user_id").(uint)

	// 更新データをマップに変換
	updates := make(map[string]interface{})
	if req.DisplayName != nil {
		updates["display_name"] = req.DisplayName
	}
	if req.Bio != nil {
		updates["bio"] = req.Bio
	}
	if req.AvatarURL != nil {
		updates["avatar_url"] = req.AvatarURL
	}
	if req.HeaderURL != nil {
		updates["header_url"] = req.HeaderURL
	}
	if req.Website != nil {
		updates["website"] = req.Website
	}
	if req.BirthDate != nil {
		updates["birth_date"] = req.BirthDate
	}
	if req.Occupation != nil {
		updates["occupation"] = req.Occupation
	}

	user, err := services.UpdateProfile(userID, updates)
	if err != nil {
		if err.Error() == "user not found" {
			return utils.ErrorResponse(c, 404, err.Error())
		}
		return utils.ErrorResponse(c, 500, "Failed to update profile")
	}

	return utils.SuccessResponse(c, 200, user.ToPublicUser())
}

// GetUserPosts - ユーザーの投稿一覧を取得ハンドラー
// @Summary ユーザーの投稿一覧取得
// @Description 指定されたユーザーの投稿一覧を取得します
// @Tags ユーザー
// @Accept json
// @Produce json
// @Param username path string true "ユーザー名"
// @Param limit query int false "取得件数（最大100）" default(20)
// @Param cursor query string false "ページネーションカーソル"
// @Success 200 {object} map[string]interface{} "data: []Post, pagination: {has_more, next_cursor, limit}"
// @Failure 404 {object} map[string]interface{} "ユーザーが見つかりません"
// @Failure 500 {object} map[string]interface{} "サーバーエラー"
// @Router /users/{username}/posts [get]
func GetUserPosts(c echo.Context) error {
	username := c.Param("username")

	limitStr := c.QueryParam("limit")
	limit := 20
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	cursor := c.QueryParam("cursor")
	var cursorPtr *string
	if cursor != "" {
		cursorPtr = &cursor
	}

	posts, hasMore, nextCursor, err := services.GetUserPosts(username, limit, cursorPtr)
	if err != nil {
		if err.Error() == "user not found" {
			return utils.ErrorResponse(c, 404, err.Error())
		}
		return utils.ErrorResponse(c, 500, "Failed to get user posts")
	}

	return utils.PaginationResponse(c, posts, hasMore, nextCursor, limit)
}

// GetFollowers - フォロワー一覧を取得ハンドラー
// @Summary フォロワー一覧取得
// @Description 指定されたユーザーのフォロワー一覧を取得します
// @Tags ユーザー
// @Accept json
// @Produce json
// @Param username path string true "ユーザー名"
// @Param limit query int false "取得件数（最大100）" default(20)
// @Param cursor query string false "ページネーションカーソル"
// @Success 200 {object} map[string]interface{} "data: []User, pagination: {has_more, next_cursor, limit}"
// @Failure 404 {object} map[string]interface{} "ユーザーが見つかりません"
// @Failure 500 {object} map[string]interface{} "サーバーエラー"
// @Router /users/{username}/followers [get]
func GetFollowers(c echo.Context) error {
	username := c.Param("username")

	limitStr := c.QueryParam("limit")
	limit := 20
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	cursor := c.QueryParam("cursor")
	var cursorPtr *string
	if cursor != "" {
		cursorPtr = &cursor
	}

	users, hasMore, nextCursor, err := services.GetFollowers(username, limit, cursorPtr)
	if err != nil {
		if err.Error() == "user not found" {
			return utils.ErrorResponse(c, 404, err.Error())
		}
		return utils.ErrorResponse(c, 500, "Failed to get followers")
	}

	return utils.PaginationResponse(c, users, hasMore, nextCursor, limit)
}

// GetFollowing - フォロー中ユーザー一覧を取得ハンドラー
// @Summary フォロー中ユーザー一覧取得
// @Description 指定されたユーザーがフォロー中のユーザー一覧を取得します
// @Tags ユーザー
// @Accept json
// @Produce json
// @Param username path string true "ユーザー名"
// @Param limit query int false "取得件数（最大100）" default(20)
// @Param cursor query string false "ページネーションカーソル"
// @Success 200 {object} map[string]interface{} "data: []User, pagination: {has_more, next_cursor, limit}"
// @Failure 404 {object} map[string]interface{} "ユーザーが見つかりません"
// @Failure 500 {object} map[string]interface{} "サーバーエラー"
// @Router /users/{username}/following [get]
func GetFollowing(c echo.Context) error {
	username := c.Param("username")

	limitStr := c.QueryParam("limit")
	limit := 20
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	cursor := c.QueryParam("cursor")
	var cursorPtr *string
	if cursor != "" {
		cursorPtr = &cursor
	}

	users, hasMore, nextCursor, err := services.GetFollowing(username, limit, cursorPtr)
	if err != nil {
		if err.Error() == "user not found" {
			return utils.ErrorResponse(c, 404, err.Error())
		}
		return utils.ErrorResponse(c, 500, "Failed to get following")
	}

	return utils.PaginationResponse(c, users, hasMore, nextCursor, limit)
}
