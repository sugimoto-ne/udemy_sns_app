package models

import "time"

// Bookmark ブックマークモデル
type Bookmark struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index:idx_user_post,unique" json:"user_id"`
	PostID    uint      `gorm:"not null;index:idx_user_post,unique;index:idx_bookmarks_post" json:"post_id"`
	CreatedAt time.Time `json:"created_at"`

	// リレーション
	User User `gorm:"foreignKey:UserID" json:"-"`
	Post Post `gorm:"foreignKey:PostID" json:"-"`
}
