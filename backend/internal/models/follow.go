package models

import (
	"time"
)

type Follow struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	FollowerID  uint      `gorm:"not null;index;uniqueIndex:idx_follower_following" json:"follower_id"`   // フォローする側
	FollowingID uint      `gorm:"not null;index;uniqueIndex:idx_follower_following" json:"following_id"` // フォローされる側
	CreatedAt   time.Time `json:"created_at"`

	// リレーション
	Follower  User `gorm:"foreignKey:FollowerID" json:"follower,omitempty"`
	Following User `gorm:"foreignKey:FollowingID" json:"following,omitempty"`
}
