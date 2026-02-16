package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/yourusername/sns-backend/internal/admin/utils"
	"github.com/yourusername/sns-backend/internal/database"
	"github.com/yourusername/sns-backend/internal/models"
)

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// ShowUserList - ユーザー一覧画面表示
func (h *UserHandler) ShowUserList(c echo.Context) error {
	adminUser := c.Get("admin_user").(models.User)

	return c.Render(http.StatusOK, "users/index.html", map[string]interface{}{
		"Title":         "ユーザー一覧",
		"AdminUsername": adminUser.Username,
		"Active":        "users",
		"Breadcrumbs": []map[string]interface{}{
			{"Name": "ダッシュボード", "URL": "/admin/dashboard", "Active": false},
			{"Name": "ユーザー一覧", "URL": "/admin/users", "Active": true},
		},
	})
}

// ShowUserDetail - ユーザー詳細画面表示
func (h *UserHandler) ShowUserDetail(c echo.Context) error {
	adminUser := c.Get("admin_user").(models.User)
	userID := c.Param("id")

	return c.Render(http.StatusOK, "users/detail.html", map[string]interface{}{
		"Title":         "ユーザー詳細",
		"AdminUsername": adminUser.Username,
		"Active":        "users",
		"UserID":        userID,
		"Breadcrumbs": []map[string]interface{}{
			{"Name": "ダッシュボード", "URL": "/admin/dashboard", "Active": false},
			{"Name": "ユーザー一覧", "URL": "/admin/users", "Active": false},
			{"Name": "ユーザー詳細", "URL": "/admin/users/" + userID, "Active": true},
		},
	})
}

// GetUsers - ユーザー一覧取得API
func (h *UserHandler) GetUsers(c echo.Context) error {
	db := database.GetDB()

	// クエリパラメータ
	status := c.QueryParam("status")
	role := c.QueryParam("role")
	search := c.QueryParam("search")
	sort := c.QueryParam("sort")
	order := c.QueryParam("order")
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if sort == "" {
		sort = "created_at"
	}
	if order == "" {
		order = "desc"
	}

	// クエリ構築
	query := db.Model(&models.User{})

	if status != "" && status != "all" {
		query = query.Where("status = ?", status)
	}
	if role != "" && role != "all" {
		query = query.Where("role = ?", role)
	}
	if search != "" {
		query = query.Where("username LIKE ? OR email LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// 総件数
	var total int64
	query.Count(&total)

	// ソート
	query = query.Order(sort + " " + order)

	// ページネーション
	offset := (page - 1) * limit
	var users []models.User
	query.Limit(limit).Offset(offset).Find(&users)

	// 投稿数を取得
	type UserWithPostCount struct {
		models.User
		PostCount int `json:"post_count"`
	}

	var usersWithPostCount []UserWithPostCount
	for _, user := range users {
		var postCount int64
		db.Model(&models.Post{}).Where("user_id = ?", user.ID).Count(&postCount)
		usersWithPostCount = append(usersWithPostCount, UserWithPostCount{
			User:      user,
			PostCount: int(postCount),
		})
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": map[string]interface{}{
			"users": usersWithPostCount,
			"pagination": map[string]interface{}{
				"total":       total,
				"page":        page,
				"limit":       limit,
				"total_pages": totalPages,
			},
		},
	})
}

// GetUser - ユーザー詳細取得API
func (h *UserHandler) GetUser(c echo.Context) error {
	db := database.GetDB()
	userID := c.Param("id")

	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	// 統計
	var postCount, likeCount, followerCount, followingCount int64
	db.Model(&models.Post{}).Where("user_id = ?", user.ID).Count(&postCount)
	db.Model(&models.PostLike{}).
		Joins("JOIN posts ON posts.id = post_likes.post_id").
		Where("posts.user_id = ?", user.ID).
		Count(&likeCount)
	db.Model(&models.Follow{}).Where("following_id = ?", user.ID).Count(&followerCount)
	db.Model(&models.Follow{}).Where("follower_id = ?", user.ID).Count(&followingCount)

	// 最近の投稿（5件）
	recentPosts := []struct {
		models.Post
		LikeCount    int64 `json:"like_count"`
		CommentCount int64 `json:"comment_count"`
	}{}

	var posts []models.Post
	db.Where("user_id = ?", user.ID).Order("created_at DESC").Limit(5).Find(&posts)

	for _, post := range posts {
		var lc, cc int64
		db.Model(&models.PostLike{}).Where("post_id = ?", post.ID).Count(&lc)
		db.Model(&models.Comment{}).Where("post_id = ?", post.ID).Count(&cc)
		recentPosts = append(recentPosts, struct {
			models.Post
			LikeCount    int64 `json:"like_count"`
			CommentCount int64 `json:"comment_count"`
		}{
			Post:         post,
			LikeCount:    lc,
			CommentCount: cc,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": map[string]interface{}{
			"user": user,
			"stats": map[string]interface{}{
				"post_count":      postCount,
				"like_count":      likeCount,
				"follower_count":  followerCount,
				"following_count": followingCount,
			},
			"recent_posts": recentPosts,
		},
	})
}

// UpdateUserStatus - ユーザーステータス変更API
func (h *UserHandler) UpdateUserStatus(c echo.Context) error {
	db := database.GetDB()
	adminUser := c.Get("admin_user").(models.User)
	userID := c.Param("id")

	var req struct {
		Status string `json:"status" validate:"required,oneof=pending approved rejected"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}

	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	oldStatus := user.Status
	user.Status = req.Status
	db.Save(&user)

	// 管理操作ログ記録
	action := "user_status_change"
	if req.Status == "approved" {
		action = "approve_user"
	} else if req.Status == "rejected" {
		action = "reject_user"
	}

	utils.LogAdminAction(db, utils.AdminLogParams{
		AdminID:        adminUser.ID,
		AdminUsername:  adminUser.Username,
		Action:         action,
		TargetUserID:   &user.ID,
		TargetUsername: &user.Username,
		Details:        fmt.Sprintf("Status changed from %s to %s", oldStatus, req.Status),
		IP:             c.RealIP(),
	})

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": map[string]interface{}{
			"user": user,
		},
		"message": "User status updated successfully",
	})
}

// BatchUpdateUserStatus - ユーザー一括ステータス変更API
func (h *UserHandler) BatchUpdateUserStatus(c echo.Context) error {
	db := database.GetDB()
	adminUser := c.Get("admin_user").(models.User)

	var req struct {
		UserIDs []uint   `json:"user_ids" validate:"required"`
		Status  string `json:"status" validate:"required,oneof=pending approved rejected"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}

	var users []models.User
	db.Find(&users, req.UserIDs)

	for _, user := range users {
		oldStatus := user.Status
		user.Status = req.Status
		db.Save(&user)

		// 管理操作ログ記録
		action := "user_status_change"
		if req.Status == "approved" {
			action = "approve_user"
		} else if req.Status == "rejected" {
			action = "reject_user"
		}

		utils.LogAdminAction(db, utils.AdminLogParams{
			AdminID:        adminUser.ID,
			AdminUsername:  adminUser.Username,
			Action:         action,
			TargetUserID:   &user.ID,
			TargetUsername: &user.Username,
			Details:        fmt.Sprintf("Batch update: Status changed from %s to %s", oldStatus, req.Status),
			IP:             c.RealIP(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": map[string]interface{}{
			"updated_count": len(users),
		},
		"message": fmt.Sprintf("%d users updated successfully", len(users)),
	})
}
