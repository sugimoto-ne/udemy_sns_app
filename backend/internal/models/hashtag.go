package models

import (
	"time"

	"gorm.io/gorm"
)

// Hashtag ハッシュタグモデル
type Hashtag struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	Posts     []Post         `gorm:"many2many:post_hashtags;" json:"-"`
}

// TableName テーブル名を指定
func (Hashtag) TableName() string {
	return "hashtags"
}

// BeforeCreate 作成前のバリデーション
func (h *Hashtag) BeforeCreate(tx *gorm.DB) error {
	// ハッシュタグ名は必須
	if h.Name == "" {
		return gorm.ErrInvalidValue
	}
	return nil
}
