package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/yourusername/sns-backend/internal/services"
	"github.com/yourusername/sns-backend/internal/utils"
)

// FollowUser - フォローハンドラー
// @Summary ユーザーをフォロー
// @Description 指定されたユーザーをフォローします
// @Tags フォロー
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param username path string true "フォローするユーザーのユーザー名"
// @Success 204 "フォロー成功"
// @Failure 400 {object} map[string]interface{} "自分自身をフォローすることはできません"
// @Failure 401 {object} map[string]interface{} "認証エラー"
// @Failure 404 {object} map[string]interface{} "ユーザーが見つかりません"
// @Failure 409 {object} map[string]interface{} "既にフォロー済み"
// @Failure 500 {object} map[string]interface{} "サーバーエラー"
// @Router /users/{username}/follow [post]
func FollowUser(c echo.Context) error {
	username := c.Param("username")
	userID := c.Get("user_id").(uint)

	if err := services.FollowUser(userID, username); err != nil {
		if err.Error() == "user not found" {
			return utils.ErrorResponse(c, 404, err.Error())
		}
		if err.Error() == "cannot follow yourself" {
			return utils.ErrorResponse(c, 400, err.Error())
		}
		if err.Error() == "already following" {
			return utils.ErrorResponse(c, 409, err.Error())
		}
		return utils.ErrorResponse(c, 500, "Failed to follow user")
	}

	return c.NoContent(204)
}

// UnfollowUser - フォロー解除ハンドラー
// @Summary フォロー解除
// @Description 指定されたユーザーのフォローを解除します
// @Tags フォロー
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param username path string true "フォロー解除するユーザーのユーザー名"
// @Success 204 "フォロー解除成功"
// @Failure 401 {object} map[string]interface{} "認証エラー"
// @Failure 404 {object} map[string]interface{} "ユーザーが見つかりません / フォローしていません"
// @Failure 500 {object} map[string]interface{} "サーバーエラー"
// @Router /users/{username}/follow [delete]
func UnfollowUser(c echo.Context) error {
	username := c.Param("username")
	userID := c.Get("user_id").(uint)

	if err := services.UnfollowUser(userID, username); err != nil {
		if err.Error() == "user not found" {
			return utils.ErrorResponse(c, 404, err.Error())
		}
		if err.Error() == "not following" {
			return utils.ErrorResponse(c, 404, err.Error())
		}
		return utils.ErrorResponse(c, 500, "Failed to unfollow user")
	}

	return c.NoContent(204)
}
