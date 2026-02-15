package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	Email         string         `gorm:"uniqueIndex;not null" json:"email" validate:"required,email"`
	Password      string         `gorm:"not null" json:"-"`
	Username      string         `gorm:"uniqueIndex;not null" json:"username" validate:"required,min=3,max=50"`
	DisplayName   *string        `json:"display_name"`
	Bio           *string        `json:"bio"`
	AvatarURL     *string        `json:"avatar_url"`
	HeaderURL     *string        `json:"header_url"`
	Website       *string        `json:"website"`
	BirthDate     *time.Time     `json:"birth_date"`
	Occupation    *string        `json:"occupation"`
	EmailVerified bool           `gorm:"default:false" json:"email_verified"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// リレーション
	Posts     []Post     `gorm:"foreignKey:UserID" json:"posts,omitempty"`
	Comments  []Comment  `gorm:"foreignKey:UserID" json:"comments,omitempty"`
	PostLikes []PostLike `gorm:"foreignKey:UserID" json:"-"`

	// フォロー関係
	Followers []Follow `gorm:"foreignKey:FollowingID" json:"-"` // このユーザーをフォローしているユーザー
	Following []Follow `gorm:"foreignKey:FollowerID" json:"-"`  // このユーザーがフォローしているユーザー
}

// BeforeCreate - パスワードをハッシュ化するフック
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return nil
}

// CheckPassword - パスワードを検証
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// PublicUser - パスワードを含まない公開用のユーザー情報
type PublicUser struct {
	ID             uint       `json:"id"`
	Email          string     `json:"email"`
	Username       string     `json:"username"`
	DisplayName    *string    `json:"display_name"`
	Bio            *string    `json:"bio"`
	AvatarURL      *string    `json:"avatar_url"`
	HeaderURL      *string    `json:"header_url"`
	Website        *string    `json:"website"`
	BirthDate      *time.Time `json:"birth_date"`
	Occupation     *string    `json:"occupation"`
	EmailVerified  bool       `json:"email_verified"`
	FollowersCount int        `json:"followers_count"`
	FollowingCount int        `json:"following_count"`
	IsFollowing    *bool      `json:"is_following,omitempty"`
	IsFollowedBy   *bool      `json:"is_followed_by,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// ToPublicUser - Userを PublicUserに変換
func (u *User) ToPublicUser() *PublicUser {
	return &PublicUser{
		ID:            u.ID,
		Email:         u.Email,
		Username:      u.Username,
		DisplayName:   u.DisplayName,
		Bio:           u.Bio,
		AvatarURL:     u.AvatarURL,
		HeaderURL:     u.HeaderURL,
		Website:       u.Website,
		BirthDate:     u.BirthDate,
		Occupation:    u.Occupation,
		EmailVerified: u.EmailVerified,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
	}
}
