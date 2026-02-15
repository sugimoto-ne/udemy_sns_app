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
