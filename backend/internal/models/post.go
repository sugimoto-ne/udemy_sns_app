package models

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	Content   string         `gorm:"type:text;not null" json:"content" validate:"required,max=280"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// リレーション
	User      User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Media     []Media    `gorm:"foreignKey:PostID" json:"media,omitempty"`
	Comments  []Comment  `gorm:"foreignKey:PostID" json:"comments,omitempty"`
	PostLikes []PostLike `gorm:"foreignKey:PostID" json:"-"`
	Hashtags  []Hashtag  `gorm:"many2many:post_hashtags;" json:"hashtags,omitempty"`

	// 集計フィールド（DBには保存しない）
	LikesCount    int64    `gorm:"-" json:"likes_count"`
	CommentsCount int64    `gorm:"-" json:"comments_count"`
	IsLiked       bool     `gorm:"-" json:"is_liked"` // 現在のユーザーがいいねしているか
	IsBookmarked  bool     `gorm:"-" json:"is_bookmarked"` // 現在のユーザーがブックマークしているか
	HashtagNames  []string `gorm:"-" json:"hashtag_names,omitempty"` // ハッシュタグ名のリスト
}

// PostWithCounts - いいね数・コメント数を含むレスポンス用構造体
type PostWithCounts struct {
	Post
	LikesCount    int64 `json:"likes_count"`
	CommentsCount int64 `json:"comments_count"`
}
