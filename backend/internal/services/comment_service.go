package services

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/yourusername/sns-backend/internal/database"
	"github.com/yourusername/sns-backend/internal/models"
	"gorm.io/gorm"
)

// GetCommentsByPostID - 投稿のコメント一覧を取得
func GetCommentsByPostID(postID uint, limit int, cursor *string) ([]models.Comment, bool, string, error) {
	db := database.GetDB()

	// 投稿が存在するかチェック
	var post models.Post
	if err := db.First(&post, postID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, "", errors.New("post not found")
		}
		return nil, false, "", err
	}

	query := db.Model(&models.Comment{}).
		Where("post_id = ?", postID).
		Preload("User")

	// カーソルベースページネーション
	if cursor != nil && *cursor != "" {
		cursorID, err := strconv.ParseUint(*cursor, 10, 64)
		if err == nil {
			query = query.Where("id < ?", cursorID)
		}
	}

	var comments []models.Comment
	if err := query.Order("created_at DESC").Limit(limit + 1).Find(&comments).Error; err != nil {
		return nil, false, "", err
	}

	hasMore := len(comments) > limit
	if hasMore {
		comments = comments[:limit]
	}

	nextCursor := ""
	if hasMore && len(comments) > 0 {
		nextCursor = fmt.Sprintf("%d", comments[len(comments)-1].ID)
	}

	return comments, hasMore, nextCursor, nil
}

// CreateComment - コメントを作成
func CreateComment(userID, postID uint, content string) (*models.Comment, error) {
	db := database.GetDB()

	// 投稿が存在するかチェック
	var post models.Post
	if err := db.First(&post, postID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("post not found")
		}
		return nil, err
	}

	comment := &models.Comment{
		PostID:  postID,
		UserID:  userID,
		Content: content,
	}

	if err := db.Create(comment).Error; err != nil {
		return nil, err
	}

	// ユーザー情報をプリロード
	db.Preload("User").First(comment, comment.ID)

	return comment, nil
}

// DeleteComment - コメントを削除（論理削除）
func DeleteComment(commentID, userID uint) error {
	db := database.GetDB()

	var comment models.Comment
	if err := db.First(&comment, commentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("comment not found")
		}
		return err
	}

	// コメント投稿者チェック
	if comment.UserID != userID {
		return errors.New("unauthorized")
	}

	// 論理削除
	if err := db.Delete(&comment).Error; err != nil {
		return err
	}

	return nil
}
