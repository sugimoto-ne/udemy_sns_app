package services

import (
	"context"
	"errors"
	"time"

	"github.com/yourusername/sns-backend/internal/database"
	"github.com/yourusername/sns-backend/internal/models"
	"github.com/yourusername/sns-backend/internal/utils"
	"gorm.io/gorm"
)

// HashtagService ハッシュタグサービス
type HashtagService struct {
	db *gorm.DB
}

// NewHashtagService ハッシュタグサービスのコンストラクタ
func NewHashtagService() *HashtagService {
	return &HashtagService{
		db: database.GetDB(),
	}
}

// ProcessHashtags 投稿内容からハッシュタグを抽出して保存・関連付け
// @param ctx コンテキスト
// @param postID 投稿ID
// @param content 投稿内容
// @return error
func (s *HashtagService) ProcessHashtags(ctx context.Context, postID uint, content string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// コンテンツからハッシュタグを抽出
	hashtagNames := utils.ExtractHashtags(content)
	if len(hashtagNames) == 0 {
		return nil // ハッシュタグがない場合は何もしない
	}

	// トランザクション開始
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, name := range hashtagNames {
			// ハッシュタグの存在確認・作成
			var hashtag models.Hashtag
			if err := tx.Where("name = ?", name).First(&hashtag).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					// 存在しない場合は新規作成
					hashtag = models.Hashtag{Name: name}
					if err := tx.Create(&hashtag).Error; err != nil {
						return err
					}
				} else {
					return err
				}
			}

			// post_hashtagsテーブルに関連付け（重複チェック）
			var count int64
			if err := tx.Model(&models.PostHashtag{}).
				Where("post_id = ? AND hashtag_id = ?", postID, hashtag.ID).
				Count(&count).Error; err != nil {
				return err
			}

			// 既に関連付けされていない場合のみ作成
			if count == 0 {
				postHashtag := models.PostHashtag{
					PostID:    postID,
					HashtagID: hashtag.ID,
				}
				if err := tx.Create(&postHashtag).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// RemovePostHashtags 投稿に関連付けられたハッシュタグを削除
// @param ctx コンテキスト
// @param postID 投稿ID
// @return error
func (s *HashtagService) RemovePostHashtags(ctx context.Context, postID uint) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return s.db.WithContext(ctx).
		Where("post_id = ?", postID).
		Delete(&models.PostHashtag{}).Error
}

// GetPostsByHashtag ハッシュタグで投稿を検索（ページネーション対応）
// @param ctx コンテキスト
// @param hashtagName ハッシュタグ名
// @param currentUserID 現在のユーザーID（いいね・ブックマーク状態取得用、0の場合は未認証）
// @param limit 取得件数
// @param cursor カーソル（最後の投稿ID）
// @return 投稿リスト, 次のカーソル, さらにデータがあるか, error
func (s *HashtagService) GetPostsByHashtag(ctx context.Context, hashtagName string, currentUserID uint, limit int, cursor uint) ([]models.Post, uint, bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// ハッシュタグを検索
	var hashtag models.Hashtag
	if err := s.db.WithContext(ctx).Where("name = ?", hashtagName).First(&hashtag).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []models.Post{}, 0, false, nil // ハッシュタグが存在しない場合は空リストを返す
		}
		return nil, 0, false, err
	}

	// 投稿を検索
	query := s.db.WithContext(ctx).
		Table("posts").
		Select("posts.*").
		Joins("INNER JOIN post_hashtags ON post_hashtags.post_id = posts.id").
		Where("post_hashtags.hashtag_id = ?", hashtag.ID).
		Where("posts.deleted_at IS NULL")

	if cursor > 0 {
		query = query.Where("posts.id < ?", cursor)
	}

	var posts []models.Post
	if err := query.
		Order("posts.created_at DESC").
		Limit(limit + 1). // hasMoreを判定するために1件多く取得
		Preload("User").
		Preload("Media", func(db *gorm.DB) *gorm.DB {
			return db.Order("order_index ASC")
		}).
		Preload("Hashtags").
		Find(&posts).Error; err != nil {
		return nil, 0, false, err
	}

	// hasMoreの判定
	hasMore := len(posts) > limit
	if hasMore {
		posts = posts[:limit] // 余分な1件を削除
	}

	// 各投稿のカウントといいね状態を取得
	postService := NewPostService()
	for i := range posts {
		posts[i].LikesCount, _ = postService.GetLikesCount(ctx, posts[i].ID)
		posts[i].CommentsCount, _ = postService.GetCommentsCount(ctx, posts[i].ID)

		if currentUserID > 0 {
			posts[i].IsLiked, _ = postService.CheckIfLiked(ctx, currentUserID, posts[i].ID)
			posts[i].IsBookmarked, _ = postService.CheckIfBookmarked(ctx, currentUserID, posts[i].ID)
		}

		// ハッシュタグ名のリストを作成
		posts[i].HashtagNames = make([]string, len(posts[i].Hashtags))
		for j, h := range posts[i].Hashtags {
			posts[i].HashtagNames[j] = h.Name
		}
	}

	// 次のカーソル
	var nextCursor uint
	if hasMore && len(posts) > 0 {
		nextCursor = posts[len(posts)-1].ID
	}

	return posts, nextCursor, hasMore, nil
}

// GetTrendingHashtags トレンドハッシュタグを取得（過去7日間で最も使用されたもの）
// @param ctx コンテキスト
// @param limit 取得件数
// @return ハッシュタグリスト（投稿数付き）, error
func (s *HashtagService) GetTrendingHashtags(ctx context.Context, limit int) ([]map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	sevenDaysAgo := time.Now().AddDate(0, 0, -7)

	type Result struct {
		HashtagID  uint   `json:"hashtag_id"`
		Name       string `json:"name"`
		PostsCount int64  `json:"posts_count"`
	}

	var results []Result
	err := s.db.WithContext(ctx).
		Table("hashtags").
		Select("hashtags.id as hashtag_id, hashtags.name, COUNT(post_hashtags.post_id) as posts_count").
		Joins("INNER JOIN post_hashtags ON post_hashtags.hashtag_id = hashtags.id").
		Joins("INNER JOIN posts ON posts.id = post_hashtags.post_id").
		Where("posts.deleted_at IS NULL").
		Where("post_hashtags.created_at >= ?", sevenDaysAgo).
		Group("hashtags.id, hashtags.name").
		Order("posts_count DESC").
		Limit(limit).
		Find(&results).Error

	if err != nil {
		return nil, err
	}

	// map形式に変換
	trending := make([]map[string]interface{}, len(results))
	for i, r := range results {
		trending[i] = map[string]interface{}{
			"id":          r.HashtagID,
			"name":        r.Name,
			"posts_count": r.PostsCount,
		}
	}

	return trending, nil
}
