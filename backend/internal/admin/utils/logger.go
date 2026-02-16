package utils

import (
	"encoding/json"

	"github.com/yourusername/sns-backend/internal/logger"
	"github.com/yourusername/sns-backend/internal/models"
	"gorm.io/gorm"
)

// AdminLogParams - 管理操作ログのパラメータ
type AdminLogParams struct {
	AdminID        uint
	AdminUsername  string
	Action         string
	TargetUserID   *uint
	TargetUsername *string
	Details        string
	IP             string
}

// LogAdminAction - 管理操作をログに記録
func LogAdminAction(db *gorm.DB, params AdminLogParams) error {
	// データベースに記録
	adminLog := models.AdminLog{
		AdminID:      params.AdminID,
		Action:       params.Action,
		TargetUserID: params.TargetUserID,
		Details:      params.Details,
		IP:           params.IP,
	}

	if err := db.Create(&adminLog).Error; err != nil {
		return err
	}

	// 構造化ログ出力
	log := logger.GetLogger()
	logData := map[string]interface{}{
		"level":           "info",
		"action":          params.Action,
		"admin_id":        params.AdminID,
		"admin_username":  params.AdminUsername,
		"target_user_id":  params.TargetUserID,
		"target_username": params.TargetUsername,
		"details":         params.Details,
		"ip":              params.IP,
	}

	jsonLog, _ := json.Marshal(logData)
	log.Info().RawJSON("admin_action", jsonLog).Msg("Admin action performed")

	return nil
}
