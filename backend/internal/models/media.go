package models

import (
	"time"
)

type Media struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	PostID     uint      `gorm:"not null;index" json:"post_id"`
	MediaType  string    `gorm:"type:varchar(20);not null" json:"media_type"` // image, video, audio
	MediaURL   string    `gorm:"type:varchar(500);not null" json:"media_url"`
	FileSize   int64     `gorm:"not null" json:"file_size"`
	Duration   *int      `json:"duration"` // 動画・音声の長さ（秒）
	OrderIndex int       `gorm:"default:0" json:"order_index"`
	CreatedAt  time.Time `json:"created_at"`

	// リレーション
	Post Post `gorm:"foreignKey:PostID" json:"-"`
}

// ValidMediaTypes - 許可されているメディアタイプ
var ValidMediaTypes = []string{"image", "video", "audio"}

// IsValidMediaType - メディアタイプが有効かチェック
func IsValidMediaType(mediaType string) bool {
	for _, validType := range ValidMediaTypes {
		if mediaType == validType {
			return true
		}
	}
	return false
}
