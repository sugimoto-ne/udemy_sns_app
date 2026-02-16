package models

import "time"

// PasswordResetRequest - パスワードリセット申請（管理画面用）
type PasswordResetRequest struct {
	ID                  uint       `gorm:"primaryKey" json:"id"`
	UserID              uint       `gorm:"not null;index" json:"user_id"`
	Token               string     `gorm:"type:varchar(255);uniqueIndex;not null" json:"token"`
	Status              string     `gorm:"type:varchar(20);default:'pending';not null;index" json:"status"` // pending, approved, expired, used
	AdminApprovedBy     *uint      `json:"admin_approved_by,omitempty"`
	AdminApprovedByUser *User      `gorm:"foreignKey:AdminApprovedBy" json:"admin_approved_by_user,omitempty"`
	AdminApprovedAt     *time.Time `json:"admin_approved_at,omitempty"`
	ExpiresAt           time.Time  `gorm:"not null" json:"expires_at"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	User                User       `gorm:"foreignKey:UserID" json:"user"`
}
