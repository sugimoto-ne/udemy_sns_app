package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/yourusername/sns-backend/internal/services"
	"github.com/yourusername/sns-backend/internal/utils"
)

// BookmarkHandler ブックマークハンドラー
type BookmarkHandler struct {
	bookmarkService *services.BookmarkService
}

// NewBookmarkHandler BookmarkHandlerのコンストラクタ
func NewBookmarkHandler() *BookmarkHandler {
	return &BookmarkHandler{
		bookmarkService: services.NewBookmarkService(),
	}
}

// BookmarkPost ブックマーク追加
// @Summary ブックマーク追加
// @Description 投稿をブックマークに追加
// @Tags bookmarks
// @Accept json
// @Produce json
// @Param id path int true "投稿ID"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Success"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Not Found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /posts/{id}/bookmark [post]
func (h *BookmarkHandler) BookmarkPost(c echo.Context) error {
	// 認証済みユーザーID取得
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
	}

	// 投稿ID取得
	postIDStr := c.Param("id")
	postIDUint64, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid post ID")
	}
	postID := uint(postIDUint64)

	// ブックマーク追加
	if err := h.bookmarkService.BookmarkPost(c.Request().Context(), userID, postID); err != nil {
		if err.Error() == "post not found" {
			return utils.ErrorResponse(c, http.StatusNotFound, "Post not found")
		}
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to bookmark post")
	}

	return utils.SuccessResponse(c, http.StatusOK, map[string]interface{}{
		"message": "Post bookmarked successfully",
	})
}

// UnbookmarkPost ブックマーク解除
// @Summary ブックマーク解除
// @Description 投稿のブックマークを解除
// @Tags bookmarks
// @Accept json
// @Produce json
// @Param id path int true "投稿ID"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Success"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /posts/{id}/bookmark [delete]
func (h *BookmarkHandler) UnbookmarkPost(c echo.Context) error {
	// 認証済みユーザーID取得
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
	}

	// 投稿ID取得
	postIDStr := c.Param("id")
	postIDUint64, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid post ID")
	}
	postID := uint(postIDUint64)

	// ブックマーク解除
	if err := h.bookmarkService.UnbookmarkPost(c.Request().Context(), userID, postID); err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to unbookmark post")
	}

	return utils.SuccessResponse(c, http.StatusOK, map[string]interface{}{
		"message": "Post unbookmarked successfully",
	})
}

// GetBookmarks ブックマーク一覧取得
// @Summary ブックマーク一覧取得
// @Description ユーザーのブックマーク一覧を取得（ページネーション対応）
// @Tags bookmarks
// @Accept json
// @Produce json
// @Param limit query int false "取得件数（デフォルト: 20）"
// @Param cursor query string false "ページネーション用カーソル"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Success"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /bookmarks [get]
func (h *BookmarkHandler) GetBookmarks(c echo.Context) error {
	// 認証済みユーザーID取得
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
	}

	// クエリパラメータ取得
	limit := 20
	if limitStr := c.QueryParam("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	var cursor *string
	if cursorStr := c.QueryParam("cursor"); cursorStr != "" {
		cursor = &cursorStr
	}

	// ブックマーク一覧取得
	posts, hasMore, nextCursor, err := h.bookmarkService.GetBookmarks(c.Request().Context(), userID, limit, cursor)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get bookmarks")
	}

	return utils.SuccessResponse(c, http.StatusOK, map[string]interface{}{
		"posts": posts,
		"pagination": map[string]interface{}{
			"has_more":    hasMore,
			"next_cursor": nextCursor,
		},
	})
}
