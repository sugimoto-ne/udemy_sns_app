package handlers

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/yourusername/sns-backend/internal/services"
	"github.com/yourusername/sns-backend/internal/utils"
)

// CreatePostRequest - 投稿作成リクエスト
type CreatePostRequest struct {
	Content string `json:"content" validate:"required,max=280"`
}

// UpdatePostRequest - 投稿更新リクエスト
type UpdatePostRequest struct {
	Content string `json:"content" validate:"required,max=280"`
}

// GetTimeline - タイムライン取得ハンドラー
// @Summary タイムライン取得
// @Description 投稿のタイムラインを取得します（全体またはフォロー中のユーザー）
// @Tags 投稿
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param type query string false "タイムラインタイプ" Enums(all, following) default(all)
// @Param limit query int false "取得件数（最大100）" default(20)
// @Param cursor query string false "ページネーションカーソル"
// @Success 200 {object} map[string]interface{} "data: []Post, pagination: {has_more, next_cursor, limit}"
// @Failure 500 {object} map[string]interface{} "サーバーエラー"
// @Router /posts [get]
func GetTimeline(c echo.Context) error {
	// クエリパラメータ
	timelineType := c.QueryParam("type") // "all" or "following"
	if timelineType == "" {
		timelineType = "all"
	}

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

	// ユーザーID取得（任意）
	var userIDPtr *uint
	if userID, ok := c.Get("user_id").(uint); ok {
		userIDPtr = &userID
	}

	// タイムライン取得
	posts, hasMore, nextCursor, err := services.GetTimeline(userIDPtr, timelineType, limit, cursorPtr)
	if err != nil {
		return utils.ErrorResponse(c, 500, "Failed to get timeline")
	}

	return utils.PaginationResponse(c, posts, hasMore, nextCursor, limit)
}

// GetPostByID - 投稿をIDで取得ハンドラー
// @Summary 投稿詳細取得
// @Description 指定されたIDの投稿を取得します
// @Tags 投稿
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "投稿ID"
// @Success 200 {object} map[string]interface{} "data: Post"
// @Failure 400 {object} map[string]interface{} "無効な投稿ID"
// @Failure 404 {object} map[string]interface{} "投稿が見つかりません"
// @Failure 500 {object} map[string]interface{} "サーバーエラー"
// @Router /posts/{id} [get]
func GetPostByID(c echo.Context) error {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, 400, "Invalid post ID")
	}

	// ユーザーID取得（任意）
	var userIDPtr *uint
	if userID, ok := c.Get("user_id").(uint); ok {
		userIDPtr = &userID
	}

	post, err := services.GetPostByID(uint(postID), userIDPtr)
	if err != nil {
		if err.Error() == "post not found" {
			return utils.ErrorResponse(c, 404, err.Error())
		}
		return utils.ErrorResponse(c, 500, "Failed to get post")
	}

	return utils.SuccessResponse(c, 200, post)
}

// CreatePost - 投稿作成ハンドラー
// @Summary 投稿作成
// @Description 新しい投稿を作成します
// @Tags 投稿
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreatePostRequest true "投稿内容"
// @Success 201 {object} map[string]interface{} "data: Post"
// @Failure 400 {object} map[string]interface{} "バリデーションエラー"
// @Failure 401 {object} map[string]interface{} "認証エラー"
// @Failure 500 {object} map[string]interface{} "サーバーエラー"
// @Router /posts [post]
func CreatePost(c echo.Context) error {
	var req CreatePostRequest
	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	if err := utils.ValidateStruct(req); err != nil {
		return utils.ErrorResponse(c, 400, err.Error())
	}

	// 安全な型アサーション
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return utils.ErrorResponse(c, 401, "Unauthorized")
	}

	// XSS対策: コンテンツをサニタイズ
	sanitizedContent := utils.SanitizeText(req.Content)

	post, err := services.CreatePost(userID, sanitizedContent)
	if err != nil {
		return utils.ErrorResponse(c, 500, "Failed to create post")
	}

	return utils.SuccessResponse(c, 201, post)
}

// UpdatePost - 投稿更新ハンドラー
// @Summary 投稿更新
// @Description 自分の投稿を更新します
// @Tags 投稿
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "投稿ID"
// @Param request body UpdatePostRequest true "更新内容"
// @Success 200 {object} map[string]interface{} "data: Post"
// @Failure 400 {object} map[string]interface{} "バリデーションエラー"
// @Failure 401 {object} map[string]interface{} "認証エラー"
// @Failure 403 {object} map[string]interface{} "権限エラー（他人の投稿は更新できません）"
// @Failure 404 {object} map[string]interface{} "投稿が見つかりません"
// @Failure 500 {object} map[string]interface{} "サーバーエラー"
// @Router /posts/{id} [put]
func UpdatePost(c echo.Context) error {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, 400, "Invalid post ID")
	}

	var req UpdatePostRequest
	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	if err := utils.ValidateStruct(req); err != nil {
		return utils.ErrorResponse(c, 400, err.Error())
	}

	// 安全な型アサーション
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return utils.ErrorResponse(c, 401, "Unauthorized")
	}

	// XSS対策: コンテンツをサニタイズ
	sanitizedContent := utils.SanitizeText(req.Content)

	post, err := services.UpdatePost(uint(postID), userID, sanitizedContent)
	if err != nil {
		if err.Error() == "post not found" {
			return utils.ErrorResponse(c, 404, err.Error())
		}
		if err.Error() == "unauthorized" {
			return utils.ErrorResponse(c, 403, err.Error())
		}
		return utils.ErrorResponse(c, 500, "Failed to update post")
	}

	return utils.SuccessResponse(c, 200, post)
}

// DeletePost - 投稿削除ハンドラー
// @Summary 投稿削除
// @Description 自分の投稿を削除します（論理削除）
// @Tags 投稿
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "投稿ID"
// @Success 204 "削除成功"
// @Failure 400 {object} map[string]interface{} "無効な投稿ID"
// @Failure 401 {object} map[string]interface{} "認証エラー"
// @Failure 403 {object} map[string]interface{} "権限エラー（他人の投稿は削除できません）"
// @Failure 404 {object} map[string]interface{} "投稿が見つかりません"
// @Failure 500 {object} map[string]interface{} "サーバーエラー"
// @Router /posts/{id} [delete]
func DeletePost(c echo.Context) error {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, 400, "Invalid post ID")
	}

	// 安全な型アサーション
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return utils.ErrorResponse(c, 401, "Unauthorized")
	}

	if err := services.DeletePost(uint(postID), userID); err != nil {
		if err.Error() == "post not found" {
			return utils.ErrorResponse(c, 404, err.Error())
		}
		if err.Error() == "unauthorized" {
			return utils.ErrorResponse(c, 403, err.Error())
		}
		return utils.ErrorResponse(c, 500, "Failed to delete post")
	}

	return c.NoContent(204)
}
