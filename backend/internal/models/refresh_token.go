package models

import (
	"time"

	"gorm.io/gorm"
)

// RefreshToken - リフレッシュトークンモデル
type RefreshToken struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	Token     string         `gorm:"uniqueIndex;not null" json:"token"` // ハッシュ化されたトークン
	ExpiresAt time.Time      `gorm:"not null" json:"expires_at"`
	Revoked   bool           `gorm:"default:false" json:"revoked"` // 失効フラグ
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// リレーション
	User User `gorm:"foreignKey:UserID" json:"-"`
}

// IsValid - トークンが有効かチェック
func (rt *RefreshToken) IsValid() bool {
	return !rt.Revoked && time.Now().Before(rt.ExpiresAt)
}
