package models

import (
	"time"
)

type PostLike struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	PostID    uint      `gorm:"not null;index;uniqueIndex:idx_post_user" json:"post_id"`
	UserID    uint      `gorm:"not null;index;uniqueIndex:idx_post_user" json:"user_id"`
	CreatedAt time.Time `json:"created_at"`

	// リレーション
	Post Post `gorm:"foreignKey:PostID" json:"post,omitempty"`
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
