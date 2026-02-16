package services

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yourusername/sns-backend/internal/database"
	"github.com/yourusername/sns-backend/internal/models"
	"github.com/yourusername/sns-backend/internal/testutil"
)

func TestHashtagService_ProcessHashtags(t *testing.T) {
	// テストDBのセットアップ
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, db)
	defer testutil.CleanupTestDB(t, db)

	// グローバルDBを設定
	database.DB = db

	service := NewHashtagService()
	ctx := context.Background()

	t.Run("Success - Process hashtags from content", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		// テストユーザーと投稿を作成
		user := testutil.CreateTestUser(t, db, "hashtag@example.com", "hashtaguser", "password123")
		post := testutil.CreateTestPost(t, db, user.ID, "これは #テスト と #GoLang についての投稿です")

		// ハッシュタグを処理
		err := service.ProcessHashtags(ctx, post.ID, post.Content)
		assert.NoError(t, err, "ハッシュタグの処理に失敗しました")

		// ハッシュタグが作成されたか確認
		var hashtags []models.Hashtag
		err = db.Find(&hashtags).Error
		require.NoError(t, err)
		assert.Len(t, hashtags, 2, "2つのハッシュタグが作成されるべき")

		// post_hashtagsの関連が作成されたか確認
		var postHashtags []models.PostHashtag
		err = db.Where("post_id = ?", post.ID).Find(&postHashtags).Error
		require.NoError(t, err)
		assert.Len(t, postHashtags, 2, "2つの関連が作成されるべき")
	})
}

func TestHashtagService_ProcessHashtags_DuplicateHashtag(t *testing.T) {
	// テストDBのセットアップ
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, db)
	defer testutil.CleanupTestDB(t, db)

	// グローバルDBを設定
	database.DB = db

	service := NewHashtagService()
	ctx := context.Background()

	t.Run("Success - Reuse existing hashtag", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		// テストユーザーを作成
		user := testutil.CreateTestUser(t, db, "dup@example.com", "dupuser", "password123")

		// 最初の投稿
		post1 := testutil.CreateTestPost(t, db, user.ID, "First post with #test")
		err := service.ProcessHashtags(ctx, post1.ID, post1.Content)
		require.NoError(t, err)

		// 2つ目の投稿（同じハッシュタグ）
		post2 := testutil.CreateTestPost(t, db, user.ID, "Second post with #test again")
		err = service.ProcessHashtags(ctx, post2.ID, post2.Content)
		assert.NoError(t, err, "既存のハッシュタグを再利用できるべき")

		// ハッシュタグは1つだけ存在するべき
		var count int64
		err = db.Model(&models.Hashtag{}).Where("name = ?", "test").Count(&count).Error
		require.NoError(t, err)
		assert.Equal(t, int64(1), count, "ハッシュタグは重複せず1つだけ存在するべき")
	})
}

func TestHashtagService_RemovePostHashtags(t *testing.T) {
	// テストDBのセットアップ
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, db)
	defer testutil.CleanupTestDB(t, db)

	// グローバルDBを設定
	database.DB = db

	service := NewHashtagService()
	ctx := context.Background()

	t.Run("Success - Remove post hashtag associations", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		// テストデータ作成
		user := testutil.CreateTestUser(t, db, "remove@example.com", "removeuser", "password123")
		post := testutil.CreateTestPost(t, db, user.ID, "Post with #remove #test")

		err := service.ProcessHashtags(ctx, post.ID, post.Content)
		require.NoError(t, err)

		// 関連を削除
		err = service.RemovePostHashtags(ctx, post.ID)
		assert.NoError(t, err, "ハッシュタグ関連の削除に失敗しました")

		// 関連が削除されたか確認
		var count int64
		err = db.Model(&models.PostHashtag{}).Where("post_id = ?", post.ID).Count(&count).Error
		require.NoError(t, err)
		assert.Equal(t, int64(0), count, "すべての関連が削除されるべき")
	})
}

func TestHashtagService_GetPostsByHashtag(t *testing.T) {
	// テストDBのセットアップ
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, db)
	defer testutil.CleanupTestDB(t, db)

	// グローバルDBを設定
	database.DB = db

	service := NewHashtagService()
	ctx := context.Background()

	t.Run("Success - Get posts by hashtag", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		// テストユーザーを作成
		user := testutil.CreateTestUser(t, db, "search@example.com", "searchuser", "password123")

		// 複数の投稿を作成
		for i := 0; i < 3; i++ {
			post := testutil.CreateTestPost(t, db, user.ID, "Post with #golang")
			err := service.ProcessHashtags(ctx, post.ID, post.Content)
			require.NoError(t, err)
		}

		// ハッシュタグで検索
		posts, nextCursor, hasMore, err := service.GetPostsByHashtag(ctx, "golang", 0, 20, 0)
		assert.NoError(t, err, "ハッシュタグ検索に失敗しました")
		assert.Len(t, posts, 3, "3つの投稿が見つかるべき")
		assert.False(t, hasMore, "次のページはないはず")
		assert.Equal(t, uint(0), nextCursor, "次のカーソルは0のはず")
	})
}

func TestHashtagService_GetPostsByHashtag_WithPagination(t *testing.T) {
	// テストDBのセットアップ
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, db)
	defer testutil.CleanupTestDB(t, db)

	// グローバルDBを設定
	database.DB = db

	service := NewHashtagService()
	ctx := context.Background()

	t.Run("Success - Paginate posts by hashtag", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		// テストユーザーを作成
		user := testutil.CreateTestUser(t, db, "page@example.com", "pageuser", "password123")

		// 5つの投稿を作成
		var postIDs []uint
		for i := 0; i < 5; i++ {
			post := testutil.CreateTestPost(t, db, user.ID, "Post with #pagination")
			err := service.ProcessHashtags(ctx, post.ID, post.Content)
			require.NoError(t, err)
			postIDs = append(postIDs, post.ID)
			time.Sleep(10 * time.Millisecond) // 順序を保証
		}

		// 最初のページ（2件取得）
		posts, nextCursor, hasMore, err := service.GetPostsByHashtag(ctx, "pagination", 0, 2, 0)
		assert.NoError(t, err)
		assert.Len(t, posts, 2, "2件取得されるべき")
		assert.True(t, hasMore, "次のページがあるはず")
		assert.NotEqual(t, uint(0), nextCursor, "次のカーソルが設定されるべき")

		// 2ページ目
		posts2, _, hasMore2, err := service.GetPostsByHashtag(ctx, "pagination", 0, 2, nextCursor)
		assert.NoError(t, err)
		assert.Len(t, posts2, 2, "2件取得されるべき")
		assert.True(t, hasMore2, "まだ次のページがあるはず")
	})
}

func TestHashtagService_GetPostsByHashtag_NonExistent(t *testing.T) {
	// テストDBのセットアップ
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, db)
	defer testutil.CleanupTestDB(t, db)

	// グローバルDBを設定
	database.DB = db

	service := NewHashtagService()
	ctx := context.Background()

	t.Run("Success - Non-existent hashtag returns empty", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		// 存在しないハッシュタグで検索
		posts, _, hasMore, err := service.GetPostsByHashtag(ctx, "nonexistent", 0, 20, 0)
		assert.NoError(t, err, "エラーは発生しないべき")
		assert.Len(t, posts, 0, "投稿は見つからないべき")
		assert.False(t, hasMore, "次のページはないはず")
	})
}

func TestHashtagService_GetTrendingHashtags(t *testing.T) {
	// テストDBのセットアップ
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, db)
	defer testutil.CleanupTestDB(t, db)

	// グローバルDBを設定
	database.DB = db

	service := NewHashtagService()
	ctx := context.Background()

	t.Run("Success - Get trending hashtags", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		// テストユーザーを作成
		user := testutil.CreateTestUser(t, db, "trend@example.com", "trenduser", "password123")

		// 異なる頻度でハッシュタグを含む投稿を作成
		hashtags := map[string]int{
			"#popular":  5,
			"#trending": 3,
			"#normal":   1,
		}

		for tag, count := range hashtags {
			for i := 0; i < count; i++ {
				post := testutil.CreateTestPost(t, db, user.ID, "Post with "+tag)
				err := service.ProcessHashtags(ctx, post.ID, post.Content)
				require.NoError(t, err)
			}
		}

		// トレンドハッシュタグを取得
		trending, err := service.GetTrendingHashtags(ctx, 10)
		assert.NoError(t, err, "トレンドハッシュタグの取得に失敗しました")
		assert.Len(t, trending, 3, "3つのハッシュタグがあるはず")

		// 最も人気のあるハッシュタグが最初に来るはず
		if len(trending) > 0 {
			firstTag := trending[0]
			assert.Equal(t, "popular", firstTag["name"], "最初のハッシュタグは'popular'であるべき")
			assert.Equal(t, int64(5), firstTag["posts_count"], "投稿数は5であるべき")
		}
	})
}

func TestHashtagService_GetTrendingHashtags_Within7Days(t *testing.T) {
	// テストDBのセットアップ
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, db)
	defer testutil.CleanupTestDB(t, db)

	// グローバルDBを設定
	database.DB = db

	service := NewHashtagService()
	ctx := context.Background()

	t.Run("Success - Only recent hashtags in trending", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		// テストユーザーを作成
		user := testutil.CreateTestUser(t, db, "timetest@example.com", "timetestuser", "password123")

		// 古い投稿（8日前）
		oldPost := testutil.CreateTestPost(t, db, user.ID, "Old post with #oldtag")

		// post_hashtagsの created_at を8日前に設定
		hashtag := &models.Hashtag{Name: "oldtag"}
		err := db.Create(hashtag).Error
		require.NoError(t, err)

		postHashtag := &models.PostHashtag{
			PostID:    oldPost.ID,
			HashtagID: hashtag.ID,
			CreatedAt: time.Now().AddDate(0, 0, -8),
		}
		err = db.Create(postHashtag).Error
		require.NoError(t, err)

		// 新しい投稿（今日）
		newPost := testutil.CreateTestPost(t, db, user.ID, "New post with #newtag")
		err = service.ProcessHashtags(ctx, newPost.ID, newPost.Content)
		require.NoError(t, err)

		// トレンドハッシュタグを取得（過去7日間のみ）
		trending, err := service.GetTrendingHashtags(ctx, 10)
		assert.NoError(t, err)

		// 新しいハッシュタグのみが含まれるはず
		foundNew := false
		foundOld := false
		for _, tag := range trending {
			if tag["name"] == "newtag" {
				foundNew = true
			}
			if tag["name"] == "oldtag" {
				foundOld = true
			}
		}

		assert.True(t, foundNew, "新しいハッシュタグが含まれるべき")
		assert.False(t, foundOld, "8日前のハッシュタグは含まれないべき")
	})
}

func TestHashtagService_GetTrendingHashtags_Limit(t *testing.T) {
	// テストDBのセットアップ
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, db)
	defer testutil.CleanupTestDB(t, db)

	// グローバルDBを設定
	database.DB = db

	service := NewHashtagService()
	ctx := context.Background()

	t.Run("Success - Limit trending hashtags", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		// テストユーザーを作成
		user := testutil.CreateTestUser(t, db, "limit@example.com", "limituser", "password123")

		// 10個のハッシュタグを作成
		for i := 0; i < 10; i++ {
			post := testutil.CreateTestPost(t, db, user.ID, "Post with #tag"+string(rune('A'+i)))
			err := service.ProcessHashtags(ctx, post.ID, post.Content)
			require.NoError(t, err)
		}

		// 上位5件のみ取得
		trending, err := service.GetTrendingHashtags(ctx, 5)
		assert.NoError(t, err)
		assert.LessOrEqual(t, len(trending), 5, "最大5件までに制限されるべき")
	})
}
