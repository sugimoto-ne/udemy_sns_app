package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yourusername/sns-backend/internal/database"
	"github.com/yourusername/sns-backend/internal/models"
)

type DashboardHandler struct{}

func NewDashboardHandler() *DashboardHandler {
	return &DashboardHandler{}
}

// ShowDashboard - ダッシュボード画面表示
func (h *DashboardHandler) ShowDashboard(c echo.Context) error {
	adminUser := c.Get("admin_user").(models.User)

	return c.Render(http.StatusOK, "dashboard.html", map[string]interface{}{
		"Title":         "ダッシュボード",
		"AdminUsername": adminUser.Username,
		"Active":        "dashboard",
		"Breadcrumbs": []map[string]interface{}{
			{"Name": "ダッシュボード", "URL": "/admin/dashboard", "Active": true},
		},
	})
}

// GetDashboardStats - ダッシュボード統計データAPI
func (h *DashboardHandler) GetDashboardStats(c echo.Context) error {
	db := database.GetDB()

	// ユーザー統計
	var totalUsers, pendingUsers, approvedUsers, rejectedUsers int64
	db.Model(&models.User{}).Count(&totalUsers)
	db.Model(&models.User{}).Where("status = ?", "pending").Count(&pendingUsers)
	db.Model(&models.User{}).Where("status = ?", "approved").Count(&approvedUsers)
	db.Model(&models.User{}).Where("status = ?", "rejected").Count(&rejectedUsers)

	// 投稿統計
	var totalPosts, todayPosts, postsWithMedia int64
	db.Model(&models.Post{}).Count(&totalPosts)
	today := time.Now().Truncate(24 * time.Hour)
	db.Model(&models.Post{}).Where("created_at >= ?", today).Count(&todayPosts)
	db.Model(&models.Post{}).
		Joins("JOIN media ON media.post_id = posts.id").
		Distinct("posts.id").
		Count(&postsWithMedia)

	var withMediaRate float64
	if totalPosts > 0 {
		withMediaRate = float64(postsWithMedia) / float64(totalPosts) * 100
	}

	// アクティブユーザー（直近7日で投稿したユーザー）
	var activeUsers7d int64
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)
	db.Model(&models.Post{}).
		Where("created_at >= ?", sevenDaysAgo).
		Distinct("user_id").
		Count(&activeUsers7d)

	// パスワードリセット申請件数（pending）
	var passwordResetPending int64
	db.Model(&models.PasswordResetRequest{}).Where("status = ?", "pending").Count(&passwordResetPending)

	// アラート判定
	alerts := []map[string]interface{}{}
	if pendingUsers >= 10 {
		alerts = append(alerts, map[string]interface{}{
			"type":    "pending_users",
			"message": "承認待ちユーザーが10人以上います",
			"count":   pendingUsers,
		})
	}
	if passwordResetPending >= 5 {
		alerts = append(alerts, map[string]interface{}{
			"type":    "password_reset",
			"message": "パスワードリセット申請が5件以上あります",
			"count":   passwordResetPending,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": map[string]interface{}{
			"users": map[string]interface{}{
				"total":    totalUsers,
				"pending":  pendingUsers,
				"approved": approvedUsers,
				"rejected": rejectedUsers,
			},
			"posts": map[string]interface{}{
				"total":           totalPosts,
				"today":           todayPosts,
				"with_media_rate": withMediaRate,
			},
			"active_users_7d":         activeUsers7d,
			"password_reset_pending": passwordResetPending,
			"alerts":                  alerts,
		},
	})
}

// GetPostsChartData - 投稿数推移データ（直近30日）
func (h *DashboardHandler) GetPostsChartData(c echo.Context) error {
	db := database.GetDB()

	labels := []string{}
	values := []int64{}

	for i := 29; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i).Truncate(24 * time.Hour)
		nextDate := date.Add(24 * time.Hour)

		var count int64
		db.Model(&models.Post{}).
			Where("created_at >= ? AND created_at < ?", date, nextDate).
			Count(&count)

		labels = append(labels, date.Format("2006-01-02"))
		values = append(values, count)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": map[string]interface{}{
			"labels": labels,
			"values": values,
		},
	})
}

// GetUsersChartData - 新規登録ユーザー推移データ（直近30日）
func (h *DashboardHandler) GetUsersChartData(c echo.Context) error {
	db := database.GetDB()

	labels := []string{}
	values := []int64{}

	for i := 29; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i).Truncate(24 * time.Hour)
		nextDate := date.Add(24 * time.Hour)

		var count int64
		db.Model(&models.User{}).
			Where("created_at >= ? AND created_at < ?", date, nextDate).
			Count(&count)

		labels = append(labels, date.Format("2006-01-02"))
		values = append(values, count)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": map[string]interface{}{
			"labels": labels,
			"values": values,
		},
	})
}
