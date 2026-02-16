package services

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/yourusername/sns-backend/internal/database"
	"github.com/yourusername/sns-backend/internal/models"
	"github.com/yourusername/sns-backend/internal/utils"
	"gorm.io/gorm"
)

// GetTimeline - タイムライン取得
func GetTimeline(userID *uint, timelineType string, limit int, cursor *string) ([]models.Post, bool, string, error) {
	db := database.GetDB()

	// limitのバリデーション
	if limit <= 0 {
		return nil, false, "", errors.New("limit must be greater than 0")
	}

	// サブクエリを使用した集計で N+1 問題を解消
	type PostWithCounts struct {
		models.Post
		LikesCount    int64 `gorm:"column:likes_count"`
		CommentsCount int64 `gorm:"column:comments_count"`
	}

	query := db.Model(&models.Post{}).
		Select(`posts.*,
			(SELECT COUNT(*) FROM post_likes WHERE post_likes.post_id = posts.id) as likes_count,
			(SELECT COUNT(*) FROM comments WHERE comments.post_id = posts.id AND comments.deleted_at IS NULL) as comments_count`).
		Preload("User").
		Preload("Media")

	// タイムラインタイプによるフィルタリング
	if timelineType == "following" && userID != nil {
		// フォロー中のユーザーの投稿のみ
		query = query.Joins("INNER JOIN follows ON follows.following_id = posts.user_id").
			Where("follows.follower_id = ?", *userID)
	}

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

	// PostWithCounts から models.Post に変換し、集計結果を設定
	posts := make([]models.Post, len(results))
	for i := range results {
		// 埋め込みフィールドを含めて全てコピー
		posts[i] = results[i].Post
		// 集計結果を明示的に設定
		posts[i].LikesCount = results[i].LikesCount
		posts[i].CommentsCount = results[i].CommentsCount
		// Preloadされたリレーションもコピーされている
		posts[i].User = results[i].Post.User
		posts[i].Media = results[i].Post.Media
	}

	// 次のカーソル
	nextCursor := ""
	if hasMore && len(posts) > 0 {
		nextCursor = fmt.Sprintf("%d", posts[len(posts)-1].ID)
	}

	// ログインユーザーのいいね状態を一括取得（N+1解消）
	if userID != nil && len(posts) > 0 {
		postIDs := make([]uint, len(posts))
		for i, post := range posts {
			postIDs[i] = post.ID
		}

		// IN句で一括取得
		var likedPosts []models.PostLike
		db.Where("post_id IN ? AND user_id = ?", postIDs, *userID).Find(&likedPosts)

		// マップ化して高速検索
		likedMap := make(map[uint]bool)
		for _, like := range likedPosts {
			likedMap[like.PostID] = true
		}

		// 投稿にいいね状態を設定
		for i := range posts {
			posts[i].IsLiked = likedMap[posts[i].ID]
		}
	}

	return posts, hasMore, nextCursor, nil
}

// GetPostByID - 投稿をIDで取得
func GetPostByID(postID uint, userID *uint) (*models.Post, error) {
	db := database.GetDB()

	var post models.Post
	if err := db.Preload("User").Preload("Media").First(&post, postID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("post not found")
		}
		return nil, err
	}

	// いいね数・コメント数を集計
	db.Model(&models.PostLike{}).Where("post_id = ?", post.ID).Count(&post.LikesCount)
	db.Model(&models.Comment{}).Where("post_id = ?", post.ID).Count(&post.CommentsCount)

	// ログインユーザーのいいね状態をチェック
	if userID != nil {
		var count int64
		db.Model(&models.PostLike{}).Where("post_id = ? AND user_id = ?", post.ID, *userID).Count(&count)
		post.IsLiked = count > 0
	}

	return &post, nil
}

// CreatePost - 投稿を作成
func CreatePost(userID uint, content string) (*models.Post, error) {
	db := database.GetDB()

	// バリデーション
	if err := utils.ValidatePostContent(content); err != nil {
		return nil, err
	}

	post := &models.Post{
		UserID:  userID,
		Content: content,
	}

	if err := db.Create(post).Error; err != nil {
		return nil, err
	}

	// ユーザー情報をプリロード
	db.Preload("User").First(post, post.ID)

	return post, nil
}

// UpdatePost - 投稿を更新
func UpdatePost(postID, userID uint, content string) (*models.Post, error) {
	db := database.GetDB()

	// バリデーション
	if err := utils.ValidatePostContent(content); err != nil {
		return nil, err
	}

	var post models.Post
	if err := db.First(&post, postID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("post not found")
		}
		return nil, err
	}

	// 投稿者チェック
	if post.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	post.Content = content
	if err := db.Save(&post).Error; err != nil {
		return nil, err
	}

	// ユーザー情報をプリロード
	db.Preload("User").First(&post, post.ID)

	return &post, nil
}

// DeletePost - 投稿を削除（論理削除）
func DeletePost(postID, userID uint) error {
	db := database.GetDB()

	var post models.Post
	if err := db.First(&post, postID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("post not found")
		}
		return err
	}

	// 投稿者チェック
	if post.UserID != userID {
		return errors.New("unauthorized")
	}

	// 論理削除
	if err := db.Delete(&post).Error; err != nil {
		return err
	}

	return nil
}

// GetUserPosts - ユーザーの投稿一覧を取得
func GetUserPosts(username string, limit int, cursor *string) ([]models.Post, bool, string, error) {
	db := database.GetDB()

	// ユーザーを取得
	var user models.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, "", errors.New("user not found")
		}
		return nil, false, "", err
	}

	// サブクエリを使用した集計で N+1 問題を解消
	type PostWithCounts struct {
		models.Post
		LikesCount    int64 `gorm:"column:likes_count"`
		CommentsCount int64 `gorm:"column:comments_count"`
	}

	query := db.Model(&models.Post{}).
		Select(`posts.*,
			(SELECT COUNT(*) FROM post_likes WHERE post_likes.post_id = posts.id) as likes_count,
			(SELECT COUNT(*) FROM comments WHERE comments.post_id = posts.id AND comments.deleted_at IS NULL) as comments_count`).
		Where("user_id = ?", user.ID).
		Preload("User").
		Preload("Media")

	// カーソルベースページネーション
	if cursor != nil && *cursor != "" {
		cursorID, err := strconv.ParseUint(*cursor, 10, 64)
		if err == nil {
			query = query.Where("id < ?", cursorID)
		}
	}

	var results []PostWithCounts
	if err := query.Order("created_at DESC").Limit(limit + 1).Find(&results).Error; err != nil {
		return nil, false, "", err
	}

	hasMore := len(results) > limit
	if hasMore {
		results = results[:limit]
	}

	// PostWithCounts から models.Post に変換し、集計結果を設定
	posts := make([]models.Post, len(results))
	for i := range results {
		// 埋め込みフィールドを含めて全てコピー
		posts[i] = results[i].Post
		// 集計結果を明示的に設定
		posts[i].LikesCount = results[i].LikesCount
		posts[i].CommentsCount = results[i].CommentsCount
		// Preloadされたリレーションもコピーされている
		posts[i].User = results[i].Post.User
		posts[i].Media = results[i].Post.Media
	}

	nextCursor := ""
	if hasMore && len(posts) > 0 {
		nextCursor = fmt.Sprintf("%d", posts[len(posts)-1].ID)
	}

	return posts, hasMore, nextCursor, nil
}
