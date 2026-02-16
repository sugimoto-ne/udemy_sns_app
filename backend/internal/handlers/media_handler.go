package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/yourusername/sns-backend/internal/services"
	"github.com/yourusername/sns-backend/internal/utils"
)

// MediaHandler メディアハンドラー
type MediaHandler struct {
	mediaService *services.MediaService
}

// NewMediaHandler MediaHandlerのコンストラクタ
func NewMediaHandler() *MediaHandler {
	return &MediaHandler{
		mediaService: services.NewMediaService(),
	}
}

// UploadMedia メディアをアップロード
// @Summary メディアアップロード
// @Description 投稿に画像/動画/音声をアップロード（最大4ファイル）
// @Tags media
// @Accept multipart/form-data
// @Produce json
// @Param post_id formData int true "投稿ID"
// @Param files formData file true "メディアファイル（最大4つ）"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Success"
// @Failure 400 {object} utils.ErrorResponse "Bad Request"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /media/upload [post]
func (h *MediaHandler) UploadMedia(c echo.Context) error {
	// 認証済みユーザーID取得
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
	}

	// 投稿ID取得
	postIDStr := c.FormValue("post_id")
	postIDUint64, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid post ID")
	}
	postID := uint(postIDUint64)

	// 投稿の所有者確認
	postService := services.NewPostService()
	post, err := postService.GetPostByID(c.Request().Context(), postID, &userID)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusNotFound, "Post not found")
	}

	if post.UserID != userID {
		return utils.ErrorResponse(c, http.StatusForbidden, "You can only upload media to your own posts")
	}

	// ファイル取得（複数）
	form, err := c.MultipartForm()
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Failed to parse multipart form")
	}

	files := form.File["files"]
	if len(files) == 0 {
		return utils.ErrorResponse(c, http.StatusBadRequest, "No files provided")
	}

	if len(files) > 4 {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Maximum 4 files allowed")
	}

	// メディアアップロード
	mediaList, err := h.mediaService.UploadMultipleMedia(c.Request().Context(), files, postID)
	if err != nil {
		// エラーをログに記録（内部詳細は含まない）
		c.Logger().Error(err)
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to upload media")
	}

	return utils.SuccessResponse(c, http.StatusOK, map[string]interface{}{
		"message": "Media uploaded successfully",
		"media":   mediaList,
	})
}

// DeleteMedia メディアを削除
// @Summary メディア削除
// @Description メディアを削除
// @Tags media
// @Accept json
// @Produce json
// @Param id path int true "メディアID"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Success"
// @Failure 400 {object} utils.ErrorResponse "Bad Request"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 404 {object} utils.ErrorResponse "Not Found"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /media/{id} [delete]
func (h *MediaHandler) DeleteMedia(c echo.Context) error {
	// 認証済みユーザーID取得
	_, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
	}

	// メディアID取得
	mediaIDStr := c.Param("id")
	mediaIDUint64, err := strconv.ParseUint(mediaIDStr, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid media ID")
	}
	mediaID := uint(mediaIDUint64)

	// メディア削除
	if err := h.mediaService.DeleteMedia(c.Request().Context(), mediaID); err != nil {
		if err.Error() == "media not found" {
			return utils.ErrorResponse(c, http.StatusNotFound, "Media not found")
		}
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete media")
	}

	return utils.SuccessResponse(c, http.StatusOK, map[string]interface{}{
		"message": "Media deleted successfully",
	})
}
