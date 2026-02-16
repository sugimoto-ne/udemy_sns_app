package services

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/yourusername/sns-backend/internal/database"
	"github.com/yourusername/sns-backend/internal/models"
	"gorm.io/gorm"
)

// BookmarkService ブックマークサービス
type BookmarkService struct {
	db *gorm.DB
}

// NewBookmarkService BookmarkServiceのコンストラクタ
func NewBookmarkService() *BookmarkService {
	return &BookmarkService{
		db: database.GetDB(),
	}
}

// BookmarkPost 投稿をブックマーク
func (s *BookmarkService) BookmarkPost(ctx context.Context, userID, postID uint) error {
	// 投稿の存在確認
	var post models.Post
	if err := s.db.WithContext(ctx).First(&post, postID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("post not found")
		}
		return err
	}

	// 重複チェック
	var count int64
	s.db.WithContext(ctx).Model(&models.Bookmark{}).
		Where("user_id = ? AND post_id = ?", userID, postID).
		Count(&count)

	if count > 0 {
		// 既にブックマーク済み（冪等性: エラーにしない）
		return nil
	}

	// ブックマーク作成
	bookmark := &models.Bookmark{
		UserID: userID,
		PostID: postID,
	}

	if err := s.db.WithContext(ctx).Create(bookmark).Error; err != nil {
		return err
	}

	return nil
}

// UnbookmarkPost ブックマークを解除
func (s *BookmarkService) UnbookmarkPost(ctx context.Context, userID, postID uint) error {
	result := s.db.WithContext(ctx).
		Where("user_id = ? AND post_id = ?", userID, postID).
		Delete(&models.Bookmark{})

	if result.Error != nil {
		return result.Error
	}

	// 削除件数が0でもエラーにしない（冪等性）
	return nil
}

// GetBookmarks ブックマーク一覧を取得（ページネーション対応）
func (s *BookmarkService) GetBookmarks(ctx context.Context, userID uint, limit int, cursor *string) ([]models.Post, bool, string, error) {
	if limit <= 0 {
		return nil, false, "", errors.New("limit must be greater than 0")
	}

	// サブクエリで集計（N+1問題解消）
	type PostWithCounts struct {
		models.Post
		LikesCount    int64 `gorm:"column:likes_count"`
		CommentsCount int64 `gorm:"column:comments_count"`
	}

	query := s.db.WithContext(ctx).Model(&models.Post{}).
		Select(`posts.*,
			(SELECT COUNT(*) FROM post_likes WHERE post_likes.post_id = posts.id) as likes_count,
			(SELECT COUNT(*) FROM comments WHERE comments.post_id = posts.id AND comments.deleted_at IS NULL) as comments_count`).
		Joins("INNER JOIN bookmarks ON bookmarks.post_id = posts.id").
		Where("bookmarks.user_id = ?", userID).
		Preload("User").
		Preload("Media")

	// カーソルベースページネーション
	if cursor != nil && *cursor != "" {
		cursorID, err := strconv.ParseUint(*cursor, 10, 64)
		if err == nil {
			query = query.Where("posts.id < ?", cursorID)
		}
	}

	// 取得件数+1を取得して、次のページがあるか判定
	var results []PostWithCounts
	if err := query.Order("posts.created_at DESC").Limit(limit + 1).Find(&results).Error; err != nil {
		return nil, false, "", err
	}

	hasMore := len(results) > limit
	if hasMore {
		results = results[:limit]
	}

	// PostWithCounts から models.Post に変換
	posts := make([]models.Post, len(results))
	for i := range results {
		posts[i] = results[i].Post
		posts[i].LikesCount = results[i].LikesCount
		posts[i].CommentsCount = results[i].CommentsCount
		posts[i].User = results[i].Post.User
		posts[i].Media = results[i].Post.Media
	}

	// 次のカーソル
	nextCursor := ""
	if hasMore && len(posts) > 0 {
		nextCursor = fmt.Sprintf("%d", posts[len(posts)-1].ID)
	}

	// いいね状態を一括取得
	if len(posts) > 0 {
		postIDs := make([]uint, len(posts))
		for i, post := range posts {
			postIDs[i] = post.ID
		}

		var likedPosts []models.PostLike
		s.db.WithContext(ctx).Where("post_id IN ? AND user_id = ?", postIDs, userID).Find(&likedPosts)

		likedMap := make(map[uint]bool)
		for _, like := range likedPosts {
			likedMap[like.PostID] = true
		}

		for i := range posts {
			posts[i].IsLiked = likedMap[posts[i].ID]
			posts[i].IsBookmarked = true // ブックマーク一覧なので常にtrue
		}
	}

	return posts, hasMore, nextCursor, nil
}

// CheckIfBookmarked ブックマーク済みかチェック
func (s *BookmarkService) CheckIfBookmarked(ctx context.Context, userID, postID uint) (bool, error) {
	var count int64
	err := s.db.WithContext(ctx).Model(&models.Bookmark{}).
		Where("user_id = ? AND post_id = ?", userID, postID).
		Count(&count).Error
	return count > 0, err
}
