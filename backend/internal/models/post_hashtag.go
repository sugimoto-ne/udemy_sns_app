package models

import "time"

// PostHashtag 投稿とハッシュタグの中間テーブル
type PostHashtag struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	PostID    uint      `gorm:"not null;index:idx_hashtag_created" json:"post_id"`
	HashtagID uint      `gorm:"not null;index:idx_hashtag_created" json:"hashtag_id"`
	CreatedAt time.Time `gorm:"index:idx_hashtag_created" json:"created_at"`

	// リレーション
	Post    Post    `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE" json:"-"`
	Hashtag Hashtag `gorm:"foreignKey:HashtagID;constraint:OnDelete:CASCADE" json:"-"`
}

// TableName テーブル名を指定
func (PostHashtag) TableName() string {
	return "post_hashtags"
}
