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
	EmailVerified bool           `gorm:"default:false" json:"email_verified"` // 廃止予定: 現在は未使用
	Approved      bool           `gorm:"default:false" json:"approved"`       // 廃止予定: statusカラムを使用
	Role          string         `gorm:"type:varchar(20);default:'user';not null" json:"role"`
	Status        string         `gorm:"type:varchar(20);default:'pending';not null" json:"status"`
	LastLoginAt   *time.Time     `json:"last_login_at,omitempty"`
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
	Email          *string    `json:"email,omitempty"` // 本人のみ表示（omitempty）
	Username       string     `json:"username"`
	DisplayName    *string    `json:"display_name"`
	Bio            *string    `json:"bio"`
	AvatarURL      *string    `json:"avatar_url"`
	HeaderURL      *string    `json:"header_url"`
	Website        *string    `json:"website"`
	BirthDate      *time.Time `json:"birth_date"`
	Occupation     *string    `json:"occupation"`
	EmailVerified  bool       `json:"email_verified"` // 廃止予定: 現在は未使用
	Approved       bool       `json:"approved"`       // 廃止予定: statusカラムを使用
	Role           string     `json:"role"`
	Status         string     `json:"status"`
	LastLoginAt    *time.Time `json:"last_login_at,omitempty"`
	FollowersCount int        `json:"followers_count"`
	FollowingCount int        `json:"following_count"`
	IsFollowing    *bool      `json:"is_following,omitempty"`
	IsFollowedBy   *bool      `json:"is_followed_by,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// ToPublicUser - Userを PublicUserに変換（閲覧者を考慮）
// viewerID: 現在閲覧しているユーザーのID（本人の場合のみEmailを含める）
func (u *User) ToPublicUser(viewerID *uint) *PublicUser {
	publicUser := &PublicUser{
		ID:            u.ID,
		Username:      u.Username,
		DisplayName:   u.DisplayName,
		Bio:           u.Bio,
		AvatarURL:     u.AvatarURL,
		HeaderURL:     u.HeaderURL,
		Website:       u.Website,
		BirthDate:     u.BirthDate,
		Occupation:    u.Occupation,
		EmailVerified: u.EmailVerified,
		Approved:      u.Approved,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
	}

	// 本人のみメールアドレスを含める
	if viewerID != nil && *viewerID == u.ID {
		publicUser.Email = &u.Email
	}

	// 管理画面用のフィールド
	publicUser.Role = u.Role
	publicUser.Status = u.Status
	publicUser.LastLoginAt = u.LastLoginAt

	return publicUser
}
