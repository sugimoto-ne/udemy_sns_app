package models

import "time"

// AdminLog - 管理者操作ログ
type AdminLog struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	AdminID      uint      `gorm:"not null;index" json:"admin_id"`
	Admin        User      `gorm:"foreignKey:AdminID" json:"admin"`
	Action       string    `gorm:"type:varchar(50);not null;index" json:"action"` // approve_user, reject_user, password_reset_approve, user_status_change
	TargetUserID *uint     `json:"target_user_id,omitempty"`
	TargetUser   *User     `gorm:"foreignKey:TargetUserID" json:"target_user,omitempty"`
	Details      string    `gorm:"type:text" json:"details"`
	IP           string    `gorm:"type:varchar(50)" json:"ip"`
	CreatedAt    time.Time `gorm:"index" json:"created_at"`
}
