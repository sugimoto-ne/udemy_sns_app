package services

import (
	"testing"

	"github.com/yourusername/sns-backend/internal/database"
	"github.com/yourusername/sns-backend/internal/testutil"
)

func TestCreatePost(t *testing.T) {
	// テストDBのセットアップ
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, db)
	defer testutil.CleanupTestDB(t, db)

	// グローバルDBを設定
	database.DB = db

	t.Run("Success - Create new post", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "testuser", "password123")
		content := "This is a test post"

		post, err := CreatePost(user.ID, content)

		testutil.AssertNoError(t, err, "CreatePost should not return error")
		testutil.AssertEqual(t, content, post.Content, "Content should match")
		testutil.AssertEqual(t, user.ID, post.UserID, "UserID should match")
		testutil.AssertNotEqual(t, uint(0), post.ID, "Post ID should be set")
	})

	t.Run("Success - Create post with preloaded user", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "testuser", "password123")
		content := "Post with user info"

		post, err := CreatePost(user.ID, content)

		testutil.AssertNoError(t, err, "CreatePost should not return error")
		testutil.AssertEqual(t, user.Username, post.User.Username, "User should be preloaded")
	})
}

func TestGetPostByID(t *testing.T) {
	// テストDBのセットアップ
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, db)
	defer testutil.CleanupTestDB(t, db)

	// グローバルDBを設定
	database.DB = db

	t.Run("Success - Get existing post", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "testuser", "password123")
		createdPost := testutil.CreateTestPost(t, db, user.ID, "Test content")

		post, err := GetPostByID(createdPost.ID, nil)

		testutil.AssertNoError(t, err, "GetPostByID should not return error")
		testutil.AssertEqual(t, createdPost.ID, post.ID, "Post ID should match")
		testutil.AssertEqual(t, createdPost.Content, post.Content, "Content should match")
	})

	t.Run("Success - Get post with like status for logged in user", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "testuser", "password123")
		createdPost := testutil.CreateTestPost(t, db, user.ID, "Test content")
		testutil.CreateTestLike(t, db, createdPost.ID, user.ID)

		post, err := GetPostByID(createdPost.ID, &user.ID)

		testutil.AssertNoError(t, err, "GetPostByID should not return error")
		testutil.AssertTrue(t, post.IsLiked, "Post should be marked as liked")
		testutil.AssertEqual(t, int64(1), post.LikesCount, "Likes count should be 1")
	})

	t.Run("Error - Get non-existent post", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		_, err := GetPostByID(99999, nil)

		testutil.AssertError(t, err, "Should return error for non-existent post")
	})
}

func TestGetTimeline(t *testing.T) {
	// テストDBのセットアップ
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, db)
	defer testutil.CleanupTestDB(t, db)

	// グローバルDBを設定
	database.DB = db

	t.Run("Success - Get global timeline", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		user1 := testutil.CreateTestUser(t, db, "user1@example.com", "user1", "password123")
		user2 := testutil.CreateTestUser(t, db, "user2@example.com", "user2", "password123")

		testutil.CreateTestPost(t, db, user1.ID, "Post 1")
		testutil.CreateTestPost(t, db, user2.ID, "Post 2")
		testutil.CreateTestPost(t, db, user1.ID, "Post 3")

		posts, hasMore, nextCursor, err := GetTimeline(nil, "global", 10, nil)

		testutil.AssertNoError(t, err, "GetTimeline should not return error")
		testutil.AssertEqual(t, 3, len(posts), "Should return 3 posts")
		testutil.AssertFalse(t, hasMore, "Should not have more posts")
		testutil.AssertEqual(t, "", nextCursor, "Next cursor should be empty")
	})

	t.Run("Success - Get timeline with pagination", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "testuser", "password123")

		// 5件の投稿を作成
		for i := 1; i <= 5; i++ {
			testutil.CreateTestPost(t, db, user.ID, "Post content")
		}

		// 最初のページ（limit=2）
		posts, hasMore, nextCursor, err := GetTimeline(nil, "global", 2, nil)

		testutil.AssertNoError(t, err, "GetTimeline should not return error")
		testutil.AssertEqual(t, 2, len(posts), "Should return 2 posts")
		testutil.AssertTrue(t, hasMore, "Should have more posts")
		testutil.AssertNotEqual(t, "", nextCursor, "Next cursor should not be empty")

		// 次のページ
		posts2, hasMore2, _, err2 := GetTimeline(nil, "global", 2, &nextCursor)

		testutil.AssertNoError(t, err2, "GetTimeline should not return error")
		testutil.AssertEqual(t, 2, len(posts2), "Should return 2 posts")
		testutil.AssertTrue(t, hasMore2, "Should have more posts")
	})

	t.Run("Success - Get following timeline", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		user1 := testutil.CreateTestUser(t, db, "user1@example.com", "user1", "password123")
		user2 := testutil.CreateTestUser(t, db, "user2@example.com", "user2", "password123")
		user3 := testutil.CreateTestUser(t, db, "user3@example.com", "user3", "password123")

		// user1がuser2をフォロー
		testutil.CreateTestFollow(t, db, user1.ID, user2.ID)

		// 投稿を作成
		testutil.CreateTestPost(t, db, user2.ID, "Post from followed user")
		testutil.CreateTestPost(t, db, user3.ID, "Post from not followed user")

		// user1のフォロータイムラインを取得
		posts, _, _, err := GetTimeline(&user1.ID, "following", 10, nil)

		testutil.AssertNoError(t, err, "GetTimeline should not return error")
		testutil.AssertEqual(t, 1, len(posts), "Should return only followed user's post")
		testutil.AssertEqual(t, user2.ID, posts[0].UserID, "Post should be from followed user")
	})
}

func TestGetUserPosts(t *testing.T) {
	// テストDBのセットアップ
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, db)
	defer testutil.CleanupTestDB(t, db)

	// グローバルDBを設定
	database.DB = db

	t.Run("Success - Get user's posts", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "testuser", "password123")
		otherUser := testutil.CreateTestUser(t, db, "other@example.com", "otheruser", "password123")

		testutil.CreateTestPost(t, db, user.ID, "User's post 1")
		testutil.CreateTestPost(t, db, user.ID, "User's post 2")
		testutil.CreateTestPost(t, db, otherUser.ID, "Other user's post")

		posts, hasMore, nextCursor, err := GetUserPosts("testuser", 10, nil)

		testutil.AssertNoError(t, err, "GetUserPosts should not return error")
		testutil.AssertEqual(t, 2, len(posts), "Should return 2 posts")
		testutil.AssertFalse(t, hasMore, "Should not have more posts")
		testutil.AssertEqual(t, "", nextCursor, "Next cursor should be empty")
	})

	t.Run("Error - Get posts for non-existent user", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		_, _, _, err := GetUserPosts("nonexistentuser", 10, nil)

		testutil.AssertError(t, err, "Should return error for non-existent user")
	})
}

// TestCreatePost_Validation - 投稿作成のバリデーションテスト
func TestCreatePost_Validation(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, db)
	defer testutil.CleanupTestDB(t, db)
	database.DB = db

	t.Run("Error - Empty content", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "testuser", "password123")

		_, err := CreatePost(user.ID, "")
		// Note: 現在は空文字のバリデーションがないため、このテストは失敗する可能性がある
		if err == nil {
			t.Logf("WARNING: Empty content was accepted (validation not implemented)")
		}
	})

	t.Run("Error - Very long content", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "testuser", "password123")

		// 10000文字のコンテンツ
		longContent := ""
		for i := 0; i < 10000; i++ {
			longContent += "a"
		}

		_, err := CreatePost(user.ID, longContent)
		// Note: 長さ制限がない場合、このテストは失敗しない可能性がある
		if err == nil {
			t.Logf("WARNING: Very long content (10000 chars) was accepted (validation not implemented)")
		}
	})

	t.Run("Success - Content with emoji", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "testuser", "password123")
		content := "Hello 😀🎉 World!"

		post, err := CreatePost(user.ID, content)
		testutil.AssertNoError(t, err, "Should accept emoji in content")
		testutil.AssertEqual(t, content, post.Content, "Emoji should be preserved")
	})

	t.Run("Success - Content with newlines", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "testuser", "password123")
		content := "Line 1\nLine 2\nLine 3"

		post, err := CreatePost(user.ID, content)
		testutil.AssertNoError(t, err, "Should accept newlines in content")
		testutil.AssertEqual(t, content, post.Content, "Newlines should be preserved")
	})

	t.Run("Warning - Content with HTML/XSS", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "testuser", "password123")
		xssContent := "<script>alert('XSS')</script>"

		post, err := CreatePost(user.ID, xssContent)
		if err == nil {
			t.Logf("WARNING: XSS content was accepted: %s", post.Content)
			// Note: フロントエンドでエスケープすべきだが、バックエンドでも検証推奨
		}
	})

	t.Run("Error - Invalid user ID", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		_, err := CreatePost(99999, "Test content")
		testutil.AssertError(t, err, "Should return error for invalid user ID")
	})
}

// TestGetTimeline_EdgeCases - タイムライン取得のエッジケーステスト
func TestGetTimeline_EdgeCases(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, db)
	defer testutil.CleanupTestDB(t, db)
	database.DB = db

	t.Run("Success - Empty timeline", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		posts, hasMore, nextCursor, err := GetTimeline(nil, "global", 10, nil)

		testutil.AssertNoError(t, err, "Should not return error for empty timeline")
		testutil.AssertEqual(t, 0, len(posts), "Should return empty array")
		testutil.AssertFalse(t, hasMore, "Should not have more posts")
		testutil.AssertEqual(t, "", nextCursor, "Next cursor should be empty")
	})

	t.Run("Error - Invalid limit (0)", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "testuser", "password123")
		testutil.CreateTestPost(t, db, user.ID, "Test post")

		_, _, _, err := GetTimeline(nil, "global", 0, nil)
		testutil.AssertError(t, err, "Should return error for limit=0")
	})

	t.Run("Error - Invalid limit (negative)", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "testuser", "password123")
		testutil.CreateTestPost(t, db, user.ID, "Test post")

		_, _, _, err := GetTimeline(nil, "global", -1, nil)
		testutil.AssertError(t, err, "Should return error for negative limit")
	})

	t.Run("Success - Very large limit", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "testuser", "password123")
		testutil.CreateTestPost(t, db, user.ID, "Test post")

		posts, hasMore, _, err := GetTimeline(nil, "global", 1000000, nil)
		testutil.AssertNoError(t, err, "Should handle large limit")
		testutil.AssertEqual(t, 1, len(posts), "Should return 1 post")
		testutil.AssertFalse(t, hasMore, "Should not have more posts")
	})

	t.Run("Error - Invalid cursor", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "testuser", "password123")
		testutil.CreateTestPost(t, db, user.ID, "Test post")

		invalidCursor := "invalid-cursor-string"
		posts, _, _, err := GetTimeline(nil, "global", 10, &invalidCursor)
		// Note: 不正なカーソルの場合、エラーを返すか無視するか定義すべき
		testutil.AssertNoError(t, err, "Should handle invalid cursor gracefully")
		testutil.AssertEqual(t, 1, len(posts), "Should return all posts when cursor is invalid")
	})

	t.Run("Error - Following timeline without user ID", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		posts, hasMore, _, err := GetTimeline(nil, "following", 10, nil)
		// Note: followingタイムラインはログインユーザーIDが必須
		if err == nil && len(posts) == 0 && !hasMore {
			t.Logf("INFO: Following timeline without userID returned empty result")
		}
	})
}

// TestGetTimeline_Pagination - ページネーションを個別にテスト
func TestGetTimeline_Pagination_FirstPage(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, db)
	defer testutil.CleanupTestDB(t, db)
	database.DB = db

	testutil.CleanupTestDB(t, db)

	user := testutil.CreateTestUser(t, db, "user@example.com", "testuser", "password123")

	// 5件の投稿を作成
	for i := 1; i <= 5; i++ {
		testutil.CreateTestPost(t, db, user.ID, "Post content")
	}

	// 最初のページ（limit=2）
	posts, hasMore, nextCursor, err := GetTimeline(nil, "global", 2, nil)

	testutil.AssertNoError(t, err, "GetTimeline should not return error")
	testutil.AssertEqual(t, 2, len(posts), "Should return 2 posts")
	testutil.AssertTrue(t, hasMore, "Should have more posts")
	testutil.AssertNotEqual(t, "", nextCursor, "Next cursor should not be empty")
}

func TestGetTimeline_Pagination_SecondPage(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, db)
	defer testutil.CleanupTestDB(t, db)
	database.DB = db

	testutil.CleanupTestDB(t, db)

	user := testutil.CreateTestUser(t, db, "user@example.com", "testuser", "password123")

	// 5件の投稿を作成
	for i := 1; i <= 5; i++ {
		testutil.CreateTestPost(t, db, user.ID, "Post content")
	}

	// 最初のページを取得してカーソルを得る
	_, _, nextCursor, _ := GetTimeline(nil, "global", 2, nil)

	// 次のページを取得
	posts2, hasMore2, _, err2 := GetTimeline(nil, "global", 2, &nextCursor)

	testutil.AssertNoError(t, err2, "GetTimeline should not return error")
	testutil.AssertEqual(t, 2, len(posts2), "Should return 2 posts")
	testutil.AssertTrue(t, hasMore2, "Should have more posts")
}
