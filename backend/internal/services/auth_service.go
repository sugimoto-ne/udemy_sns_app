package services

import (
	"errors"

	"github.com/yourusername/sns-backend/internal/database"
	"github.com/yourusername/sns-backend/internal/models"
	"gorm.io/gorm"
)

// Register - ユーザー登録
func Register(email, password, username string) (*models.User, error) {
	db := database.GetDB()

	// メールアドレスの重複チェック
	var existingUser models.User
	if err := db.Where("email = ?", email).First(&existingUser).Error; err == nil {
		return nil, errors.New("email already exists")
	}

	// ユーザー名の重複チェック
	if err := db.Where("username = ?", username).First(&existingUser).Error; err == nil {
		return nil, errors.New("username already exists")
	}

	// ユーザー作成（パスワードはBeforeCreateフックで自動ハッシュ化）
	user := &models.User{
		Email:    email,
		Password: password,
		Username: username,
	}

	if err := db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// Login - ログイン
func Login(email, password string) (*models.User, error) {
	db := database.GetDB()

	var user models.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	// パスワード検証
	if !user.CheckPassword(password) {
		return nil, errors.New("invalid email or password")
	}

	return &user, nil
}

// GetCurrentUser - 現在のユーザー情報取得
func GetCurrentUser(userID uint) (*models.User, error) {
	db := database.GetDB()

	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}
