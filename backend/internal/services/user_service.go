package services

import (
	"errors"

	"github.com/yourusername/sns-backend/internal/database"
	"github.com/yourusername/sns-backend/internal/models"
	"gorm.io/gorm"
)

// GetUserByUsername - ユーザー名でユーザーを取得
func GetUserByUsername(username string, currentUserID *uint) (*models.User, error) {
	db := database.GetDB()

	var user models.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// UpdateProfile - プロフィールを更新
func UpdateProfile(userID uint, updates map[string]interface{}) (*models.User, error) {
	db := database.GetDB()

	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// 許可されたフィールドのみ更新
	allowedFields := map[string]bool{
		"display_name": true,
		"bio":          true,
		"avatar_url":   true,
		"header_url":   true,
		"website":      true,
		"birth_date":   true,
		"occupation":   true,
	}

	filteredUpdates := make(map[string]interface{})
	for key, value := range updates {
		if allowedFields[key] {
			filteredUpdates[key] = value
		}
	}

	if err := db.Model(&user).Updates(filteredUpdates).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
