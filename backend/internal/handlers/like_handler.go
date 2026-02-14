package handlers

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/yourusername/sns-backend/internal/services"
	"github.com/yourusername/sns-backend/internal/utils"
)

// LikePost - いいねハンドラー
// @Summary 投稿にいいね
// @Description 投稿にいいねを追加します
// @Tags いいね
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "投稿ID"
// @Success 204 "いいね成功"
// @Failure 400 {object} map[string]interface{} "無効な投稿ID"
// @Failure 401 {object} map[string]interface{} "認証エラー"
// @Failure 404 {object} map[string]interface{} "投稿が見つかりません"
// @Failure 409 {object} map[string]interface{} "既にいいね済み"
// @Failure 500 {object} map[string]interface{} "サーバーエラー"
// @Router /posts/{id}/like [post]
func LikePost(c echo.Context) error {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, 400, "Invalid post ID")
	}

	userID := c.Get("user_id").(uint)

	if err := services.LikePost(userID, uint(postID)); err != nil {
		if err.Error() == "post not found" {
			return utils.ErrorResponse(c, 404, err.Error())
		}
		if err.Error() == "already liked" {
			return utils.ErrorResponse(c, 409, err.Error())
		}
		return utils.ErrorResponse(c, 500, "Failed to like post")
	}

	return c.NoContent(204)
}

// UnlikePost - いいね解除ハンドラー
// @Summary いいね解除
// @Description 投稿のいいねを解除します
// @Tags いいね
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "投稿ID"
// @Success 204 "いいね解除成功"
// @Failure 400 {object} map[string]interface{} "無効な投稿ID"
// @Failure 401 {object} map[string]interface{} "認証エラー"
// @Failure 404 {object} map[string]interface{} "いいねが見つかりません"
// @Failure 500 {object} map[string]interface{} "サーバーエラー"
// @Router /posts/{id}/like [delete]
func UnlikePost(c echo.Context) error {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, 400, "Invalid post ID")
	}

	userID := c.Get("user_id").(uint)

	if err := services.UnlikePost(userID, uint(postID)); err != nil {
		if err.Error() == "like not found" {
			return utils.ErrorResponse(c, 404, err.Error())
		}
		return utils.ErrorResponse(c, 500, "Failed to unlike post")
	}

	return c.NoContent(204)
}

// GetLikes - いいね一覧取得ハンドラー
// @Summary いいね一覧取得
// @Description 投稿にいいねしたユーザー一覧を取得します
// @Tags いいね
// @Accept json
// @Produce json
// @Param id path int true "投稿ID"
// @Param limit query int false "取得件数（最大100）" default(20)
// @Param cursor query string false "ページネーションカーソル"
// @Success 200 {object} map[string]interface{} "data: []User, pagination: {has_more, next_cursor, limit}"
// @Failure 400 {object} map[string]interface{} "無効な投稿ID"
// @Failure 404 {object} map[string]interface{} "投稿が見つかりません"
// @Failure 500 {object} map[string]interface{} "サーバーエラー"
// @Router /posts/{id}/likes [get]
func GetLikes(c echo.Context) error {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, 400, "Invalid post ID")
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

	users, hasMore, nextCursor, err := services.GetLikesByPostID(uint(postID), limit, cursorPtr)
	if err != nil {
		if err.Error() == "post not found" {
			return utils.ErrorResponse(c, 404, err.Error())
		}
		return utils.ErrorResponse(c, 500, "Failed to get likes")
	}

	return utils.PaginationResponse(c, users, hasMore, nextCursor, limit)
}
