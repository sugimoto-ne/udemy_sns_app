package services

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/yourusername/sns-backend/internal/database"
	"github.com/yourusername/sns-backend/internal/models"
	"gorm.io/gorm"
)

// LikePost - 投稿にいいね
func LikePost(userID, postID uint) error {
	db := database.GetDB()

	// 投稿が存在するかチェック
	var post models.Post
	if err := db.First(&post, postID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("post not found")
		}
		return err
	}

	// 既にいいね済みかチェック
	var existingLike models.PostLike
	if err := db.Where("post_id = ? AND user_id = ?", postID, userID).First(&existingLike).Error; err == nil {
		return errors.New("already liked")
	}

	// いいねを作成
	like := &models.PostLike{
		PostID: postID,
		UserID: userID,
	}

	if err := db.Create(like).Error; err != nil {
		return err
	}

	return nil
}

// UnlikePost - 投稿のいいねを解除
func UnlikePost(userID, postID uint) error {
	db := database.GetDB()

	// いいねを取得
	var like models.PostLike
	if err := db.Where("post_id = ? AND user_id = ?", postID, userID).First(&like).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("like not found")
		}
		return err
	}

	// いいねを削除
	if err := db.Delete(&like).Error; err != nil {
		return err
	}

	return nil
}

// GetLikesByPostID - 投稿のいいね一覧を取得
func GetLikesByPostID(postID uint, limit int, cursor *string) ([]models.User, bool, string, error) {
	db := database.GetDB()

	// 投稿が存在するかチェック
	var post models.Post
	if err := db.First(&post, postID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, "", errors.New("post not found")
		}
		return nil, false, "", err
	}

	query := db.Model(&models.User{}).
		Joins("INNER JOIN post_likes ON post_likes.user_id = users.id").
		Where("post_likes.post_id = ?", postID)

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

// CheckIfLiked - ユーザーが投稿にいいねしているかチェック
func CheckIfLiked(userID, postID uint) (bool, error) {
	db := database.GetDB()

	var count int64
	if err := db.Model(&models.PostLike{}).
		Where("post_id = ? AND user_id = ?", postID, userID).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}
