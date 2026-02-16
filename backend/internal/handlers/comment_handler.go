package handlers

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/yourusername/sns-backend/internal/services"
	"github.com/yourusername/sns-backend/internal/utils"
)

// CreateCommentRequest - コメント作成リクエスト
type CreateCommentRequest struct {
	Content string `json:"content" validate:"required,max=280"`
}

// GetComments - コメント一覧取得ハンドラー
// @Summary コメント一覧取得
// @Description 指定された投稿のコメント一覧を取得します
// @Tags コメント
// @Accept json
// @Produce json
// @Param id path int true "投稿ID"
// @Param limit query int false "取得件数（最大100）" default(20)
// @Param cursor query string false "ページネーションカーソル"
// @Success 200 {object} map[string]interface{} "data: []Comment, pagination: {has_more, next_cursor, limit}"
// @Failure 400 {object} map[string]interface{} "無効な投稿ID"
// @Failure 404 {object} map[string]interface{} "投稿が見つかりません"
// @Failure 500 {object} map[string]interface{} "サーバーエラー"
// @Router /posts/{id}/comments [get]
func GetComments(c echo.Context) error {
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

	comments, hasMore, nextCursor, err := services.GetCommentsByPostID(uint(postID), limit, cursorPtr)
	if err != nil {
		if err.Error() == "post not found" {
			return utils.ErrorResponse(c, 404, err.Error())
		}
		return utils.ErrorResponse(c, 500, "Failed to get comments")
	}

	return utils.PaginationResponse(c, comments, hasMore, nextCursor, limit)
}

// CreateComment - コメント作成ハンドラー
// @Summary コメント作成
// @Description 投稿にコメントを追加します
// @Tags コメント
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "投稿ID"
// @Param request body CreateCommentRequest true "コメント内容"
// @Success 201 {object} map[string]interface{} "data: Comment"
// @Failure 400 {object} map[string]interface{} "バリデーションエラー"
// @Failure 401 {object} map[string]interface{} "認証エラー"
// @Failure 404 {object} map[string]interface{} "投稿が見つかりません"
// @Failure 500 {object} map[string]interface{} "サーバーエラー"
// @Router /posts/{id}/comments [post]
func CreateComment(c echo.Context) error {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, 400, "Invalid post ID")
	}

	var req CreateCommentRequest
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

	comment, err := services.CreateComment(userID, uint(postID), sanitizedContent)
	if err != nil {
		if err.Error() == "post not found" {
			return utils.ErrorResponse(c, 404, err.Error())
		}
		return utils.ErrorResponse(c, 500, "Failed to create comment")
	}

	return utils.SuccessResponse(c, 201, comment)
}

// DeleteComment - コメント削除ハンドラー
// @Summary コメント削除
// @Description 自分のコメントを削除します（論理削除）
// @Tags コメント
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "コメントID"
// @Success 204 "削除成功"
// @Failure 400 {object} map[string]interface{} "無効なコメントID"
// @Failure 401 {object} map[string]interface{} "認証エラー"
// @Failure 403 {object} map[string]interface{} "権限エラー（他人のコメントは削除できません）"
// @Failure 404 {object} map[string]interface{} "コメントが見つかりません"
// @Failure 500 {object} map[string]interface{} "サーバーエラー"
// @Router /comments/{id} [delete]
func DeleteComment(c echo.Context) error {
	commentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, 400, "Invalid comment ID")
	}

	// 安全な型アサーション
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return utils.ErrorResponse(c, 401, "Unauthorized")
	}

	if err := services.DeleteComment(uint(commentID), userID); err != nil {
		if err.Error() == "comment not found" {
			return utils.ErrorResponse(c, 404, err.Error())
		}
		if err.Error() == "unauthorized" {
			return utils.ErrorResponse(c, 403, err.Error())
		}
		return utils.ErrorResponse(c, 500, "Failed to delete comment")
	}

	return c.NoContent(204)
}
