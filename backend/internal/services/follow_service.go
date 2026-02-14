package services

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/yourusername/sns-backend/internal/database"
	"github.com/yourusername/sns-backend/internal/models"
	"gorm.io/gorm"
)

// FollowUser - ユーザーをフォロー
func FollowUser(followerID uint, followingUsername string) error {
	db := database.GetDB()

	// フォロー対象のユーザーを取得
	var followingUser models.User
	if err := db.Where("username = ?", followingUsername).First(&followingUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	// 自分自身をフォローしようとしていないかチェック
	if followerID == followingUser.ID {
		return errors.New("cannot follow yourself")
	}

	// 既にフォロー済みかチェック
	var existingFollow models.Follow
	if err := db.Where("follower_id = ? AND following_id = ?", followerID, followingUser.ID).First(&existingFollow).Error; err == nil {
		return errors.New("already following")
	}

	// フォロー関係を作成
	follow := &models.Follow{
		FollowerID:  followerID,
		FollowingID: followingUser.ID,
	}

	if err := db.Create(follow).Error; err != nil {
		return err
	}

	return nil
}

// UnfollowUser - ユーザーのフォローを解除
func UnfollowUser(followerID uint, followingUsername string) error {
	db := database.GetDB()

	// フォロー解除対象のユーザーを取得
	var followingUser models.User
	if err := db.Where("username = ?", followingUsername).First(&followingUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	// フォロー関係を取得
	var follow models.Follow
	if err := db.Where("follower_id = ? AND following_id = ?", followerID, followingUser.ID).First(&follow).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("not following")
		}
		return err
	}

	// フォロー関係を削除
	if err := db.Delete(&follow).Error; err != nil {
		return err
	}

	return nil
}

// CheckIfFollowing - ユーザーがフォローしているかチェック
func CheckIfFollowing(followerID, followingID uint) (bool, error) {
	db := database.GetDB()

	var count int64
	if err := db.Model(&models.Follow{}).
		Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetFollowers - フォロワー一覧を取得
func GetFollowers(username string, limit int, cursor *string) ([]models.User, bool, string, error) {
	db := database.GetDB()

	// ユーザーを取得
	var user models.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, "", errors.New("user not found")
		}
		return nil, false, "", err
	}

	query := db.Model(&models.User{}).
		Joins("INNER JOIN follows ON follows.follower_id = users.id").
		Where("follows.following_id = ?", user.ID)

	// カーソルベースページネーション
	if cursor != nil && *cursor != "" {
		cursorID, err := strconv.ParseUint(*cursor, 10, 64)
		if err == nil {
			query = query.Where("users.id < ?", cursorID)
		}
	}

	var users []models.User
	if err := query.Order("users.id DESC").Limit(limit + 1).Find(&users).Error; err != nil {
		return nil, false, "", err
	}

	hasMore := len(users) > limit
	if hasMore {
		users = users[:limit]
	}

	nextCursor := ""
	if hasMore && len(users) > 0 {
		nextCursor = fmt.Sprintf("%d", users[len(users)-1].ID)
	}

	return users, hasMore, nextCursor, nil
}

// GetFollowing - フォロー中ユーザー一覧を取得
func GetFollowing(username string, limit int, cursor *string) ([]models.User, bool, string, error) {
	db := database.GetDB()

	// ユーザーを取得
	var user models.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, "", errors.New("user not found")
		}
		return nil, false, "", err
	}

	query := db.Model(&models.User{}).
		Joins("INNER JOIN follows ON follows.following_id = users.id").
		Where("follows.follower_id = ?", user.ID)

	// カーソルベースページネーション
	if cursor != nil && *cursor != "" {
		cursorID, err := strconv.ParseUint(*cursor, 10, 64)
		if err == nil {
			query = query.Where("users.id < ?", cursorID)
		}
	}

	var users []models.User
	if err := query.Order("users.id DESC").Limit(limit + 1).Find(&users).Error; err != nil {
		return nil, false, "", err
	}

	hasMore := len(users) > limit
	if hasMore {
		users = users[:limit]
	}

	nextCursor := ""
	if hasMore && len(users) > 0 {
		nextCursor = fmt.Sprintf("%d", users[len(users)-1].ID)
	}

	return users, hasMore, nextCursor, nil
}
