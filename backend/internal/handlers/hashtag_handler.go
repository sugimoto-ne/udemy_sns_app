package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/yourusername/sns-backend/internal/services"
	"github.com/yourusername/sns-backend/internal/utils"
)

// HashtagHandler ハッシュタグハンドラー
type HashtagHandler struct {
	hashtagService *services.HashtagService
}

// NewHashtagHandler ハッシュタグハンドラーのコンストラクタ
func NewHashtagHandler() *HashtagHandler {
	return &HashtagHandler{
		hashtagService: services.NewHashtagService(),
	}
}

// GetPostsByHashtag ハッシュタグで投稿を検索
// @Summary ハッシュタグで投稿を検索
// @Description 指定したハッシュタグに関連する投稿を取得します
// @Tags hashtags
// @Accept json
// @Produce json
// @Param name path string true "ハッシュタグ名（#を除く）"
// @Param limit query int false "取得件数" default(20)
// @Param cursor query int false "カーソル（最後の投稿ID）"
// @Success 200 {object} map[string]interface{} "投稿リストとページネーション情報"
// @Failure 400 {object} map[string]interface{} "バリデーションエラー"
// @Failure 500 {object} map[string]interface{} "サーバーエラー"
// @Security BearerAuth
// @Router /api/v1/hashtags/{name}/posts [get]
func (h *HashtagHandler) GetPostsByHashtag(c echo.Context) error {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	// パラメータ取得
	hashtagName := c.Param("name")
	if hashtagName == "" {
		return utils.ErrorResponse(c, http.StatusBadRequest, "ハッシュタグ名は必須です")
	}

	// limit取得（デフォルト20）
	limitStr := c.QueryParam("limit")
	limit := 20
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 && parsedLimit <= 100 {
			limit = parsedLimit
		}
	}

	// cursor取得
	cursorStr := c.QueryParam("cursor")
	var cursor uint
	if cursorStr != "" {
		parsedCursor, err := strconv.ParseUint(cursorStr, 10, 32)
		if err == nil {
			cursor = uint(parsedCursor)
		}
	}

	// 現在のユーザーID取得（オプション）
	var currentUserID uint
	userIDInterface := c.Get("user_id")
	if userIDInterface != nil {
		currentUserID = userIDInterface.(uint)
	}

	// 投稿を検索
	posts, nextCursor, hasMore, err := h.hashtagService.GetPostsByHashtag(ctx, hashtagName, currentUserID, limit, cursor)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "投稿の取得に失敗しました")
	}

	return utils.SuccessResponse(c, http.StatusOK, map[string]interface{}{
		"posts": posts,
		"pagination": map[string]interface{}{
			"has_more":    hasMore,
			"next_cursor": nextCursor,
			"limit":       limit,
		},
	})
}

// GetTrendingHashtags トレンドハッシュタグを取得
// @Summary トレンドハッシュタグを取得
// @Description 過去7日間で最も使用されたハッシュタグを取得します
// @Tags hashtags
// @Accept json
// @Produce json
// @Param limit query int false "取得件数" default(10)
// @Success 200 {object} map[string]interface{} "トレンドハッシュタグリスト"
// @Failure 500 {object} map[string]interface{} "サーバーエラー"
// @Router /api/v1/hashtags/trending [get]
func (h *HashtagHandler) GetTrendingHashtags(c echo.Context) error {
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	// limit取得（デフォルト10）
	limitStr := c.QueryParam("limit")
	limit := 10
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 && parsedLimit <= 50 {
			limit = parsedLimit
		}
	}

	// トレンドハッシュタグを取得
	trending, err := h.hashtagService.GetTrendingHashtags(ctx, limit)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "トレンドハッシュタグの取得に失敗しました")
	}

	return utils.SuccessResponse(c, http.StatusOK, map[string]interface{}{
		"hashtags": trending,
	})
}
