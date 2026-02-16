package models

import "time"

// PasswordResetToken パスワードリセットトークンモデル
type PasswordResetToken struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	Token     string    `gorm:"type:varchar(255);not null;uniqueIndex" json:"token"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`

	// リレーション
	User User `gorm:"foreignKey:UserID" json:"-"`
}
