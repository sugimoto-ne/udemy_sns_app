package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yourusername/sns-backend/internal/database"
	"github.com/yourusername/sns-backend/internal/models"
)

type LogHandler struct{}

func NewLogHandler() *LogHandler {
	return &LogHandler{}
}

// ShowLogList - 操作ログ一覧画面表示
func (h *LogHandler) ShowLogList(c echo.Context) error {
	adminUser := c.Get("admin_user").(models.User)

	return c.Render(http.StatusOK, "logs/index.html", map[string]interface{}{
		"Title":         "操作ログ",
		"AdminUsername": adminUser.Username,
		"Active":        "logs",
		"Breadcrumbs": []map[string]interface{}{
			{"Name": "ダッシュボード", "URL": "/admin/dashboard", "Active": false},
			{"Name": "操作ログ", "URL": "/admin/logs", "Active": true},
		},
	})
}

// GetLogs - 操作ログ一覧取得API
func (h *LogHandler) GetLogs(c echo.Context) error {
	db := database.GetDB()

	// クエリパラメータ
	action := c.QueryParam("action")
	adminUsername := c.QueryParam("admin_username")
	startDate := c.QueryParam("start_date")
	endDate := c.QueryParam("end_date")
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 50
	}

	// クエリ構築
	query := db.Model(&models.AdminLog{}).
		Preload("Admin").
		Preload("TargetUser")

	if action != "" && action != "all" {
		query = query.Where("action = ?", action)
	}

	if adminUsername != "" {
		query = query.Joins("JOIN users ON users.id = admin_logs.admin_id").
			Where("users.username = ?", adminUsername)
	}

	if startDate != "" {
		start, err := time.Parse("2006-01-02", startDate)
		if err == nil {
			query = query.Where("admin_logs.created_at >= ?", start)
		}
	}

	if endDate != "" {
		end, err := time.Parse("2006-01-02", endDate)
		if err == nil {
			// 終了日の23:59:59まで含める
			end = end.Add(24*time.Hour - time.Second)
			query = query.Where("admin_logs.created_at <= ?", end)
		}
	}

	// 総件数
	var total int64
	query.Count(&total)

	// ページネーション
	offset := (page - 1) * limit
	var logs []models.AdminLog
	query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&logs)

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": map[string]interface{}{
			"logs": logs,
			"pagination": map[string]interface{}{
				"total":       total,
				"page":        page,
				"limit":       limit,
				"total_pages": totalPages,
			},
		},
	})
}
