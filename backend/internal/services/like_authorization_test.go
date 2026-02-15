package services

import (
	"testing"

	"github.com/yourusername/sns-backend/internal/database"
	"github.com/yourusername/sns-backend/internal/testutil"
)

// TestLikePost_Authorization - いいねの認可・冪等性テスト
func TestLikePost_Authorization(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, db)
	defer testutil.CleanupTestDB(t, db)
	database.DB = db

	t.Run("Success - Like a post", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "user", "password123")
		post := testutil.CreateTestPost(t, db, user.ID, "Test post")

		err := LikePost(user.ID, post.ID)
		testutil.AssertNoError(t, err, "Should be able to like a post")

		// いいねが作成されたことを確認
		isLiked, err := CheckIfLiked(user.ID, post.ID)
		testutil.AssertNoError(t, err, "CheckIfLiked should not return error")
		testutil.AssertTrue(t, isLiked, "Post should be liked")
	})

	t.Run("Error - Cannot like same post twice (idempotency)", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "user", "password123")
		post := testutil.CreateTestPost(t, db, user.ID, "Test post")

		// 最初のいいね
		err := LikePost(user.ID, post.ID)
		testutil.AssertNoError(t, err, "First like should succeed")

		// 2回目のいいね（冪等性テスト）
		err = LikePost(user.ID, post.ID)
		testutil.AssertError(t, err, "Should not be able to like same post twice")
	})

	t.Run("Success - Multiple users can like same post", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		user1 := testutil.CreateTestUser(t, db, "user1@example.com", "user1", "password123")
		user2 := testutil.CreateTestUser(t, db, "user2@example.com", "user2", "password123")
		post := testutil.CreateTestPost(t, db, user1.ID, "Test post")

		// user1がいいね
		err := LikePost(user1.ID, post.ID)
		testutil.AssertNoError(t, err, "User1 should be able to like")

		// user2がいいね
		err = LikePost(user2.ID, post.ID)
		testutil.AssertNoError(t, err, "User2 should be able to like")

		// 両方のユーザーがいいねしていることを確認
		isLiked1, _ := CheckIfLiked(user1.ID, post.ID)
		isLiked2, _ := CheckIfLiked(user2.ID, post.ID)
		testutil.AssertTrue(t, isLiked1, "User1 should have liked")
		testutil.AssertTrue(t, isLiked2, "User2 should have liked")
	})

	t.Run("Error - Cannot like non-existent post", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "user", "password123")

		err := LikePost(user.ID, 99999)
		testutil.AssertError(t, err, "Should return error for non-existent post")
	})

	t.Run("Error - Cannot like deleted post", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "user", "password123")
		post := testutil.CreateTestPost(t, db, user.ID, "Test post")

		// 投稿を削除
		err := DeletePost(post.ID, user.ID)
		testutil.AssertNoError(t, err, "Post deletion should succeed")

		// 削除された投稿にいいねしようとする
		err = LikePost(user.ID, post.ID)
		testutil.AssertError(t, err, "Should not be able to like deleted post")
	})
}

// TestUnlikePost_Authorization - いいね解除のテスト
func TestUnlikePost_Authorization(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, db)
	defer testutil.CleanupTestDB(t, db)
	database.DB = db

	t.Run("Success - Unlike a liked post", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "user", "password123")
		post := testutil.CreateTestPost(t, db, user.ID, "Test post")

		// いいね
		err := LikePost(user.ID, post.ID)
		testutil.AssertNoError(t, err, "Like should succeed")

		// いいね解除
		err = UnlikePost(user.ID, post.ID)
		testutil.AssertNoError(t, err, "Unlike should succeed")

		// いいねが解除されたことを確認
		isLiked, err := CheckIfLiked(user.ID, post.ID)
		testutil.AssertNoError(t, err, "CheckIfLiked should not return error")
		testutil.AssertFalse(t, isLiked, "Post should not be liked")
	})

	t.Run("Error - Cannot unlike post that was not liked", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "user", "password123")
		post := testutil.CreateTestPost(t, db, user.ID, "Test post")

		// いいねしていない投稿のいいねを解除しようとする
		err := UnlikePost(user.ID, post.ID)
		testutil.AssertError(t, err, "Should return error when trying to unlike a post that was not liked")
	})

	t.Run("Error - Cannot unlike same post twice", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "user", "password123")
		post := testutil.CreateTestPost(t, db, user.ID, "Test post")

		// いいね
		err := LikePost(user.ID, post.ID)
		testutil.AssertNoError(t, err, "Like should succeed")

		// いいね解除
		err = UnlikePost(user.ID, post.ID)
		testutil.AssertNoError(t, err, "Unlike should succeed")

		// 再度いいね解除しようとする
		err = UnlikePost(user.ID, post.ID)
		testutil.AssertError(t, err, "Should not be able to unlike twice")
	})
}

// TestGetLikesByPostID_EdgeCases - いいね一覧取得のエッジケーステスト
func TestGetLikesByPostID_EdgeCases(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, db)
	defer testutil.CleanupTestDB(t, db)
	database.DB = db

	t.Run("Success - Get likes for post with no likes", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "user", "password123")
		post := testutil.CreateTestPost(t, db, user.ID, "Test post")

		users, hasMore, nextCursor, err := GetLikesByPostID(post.ID, 10, nil)
		testutil.AssertNoError(t, err, "Should not return error for post with no likes")
		testutil.AssertEqual(t, 0, len(users), "Should return empty array")
		testutil.AssertFalse(t, hasMore, "Should not have more")
		testutil.AssertEqual(t, "", nextCursor, "Next cursor should be empty")
	})

	t.Run("Error - Get likes for non-existent post", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		_, _, _, err := GetLikesByPostID(99999, 10, nil)
		testutil.AssertError(t, err, "Should return error for non-existent post")
	})

	t.Run("Error - Get likes for deleted post", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "user", "password123")
		post := testutil.CreateTestPost(t, db, user.ID, "Test post")

		// 投稿を削除
		err := DeletePost(post.ID, user.ID)
		testutil.AssertNoError(t, err, "Post deletion should succeed")

		// 削除された投稿のいいね一覧を取得しようとする
		_, _, _, err = GetLikesByPostID(post.ID, 10, nil)
		testutil.AssertError(t, err, "Should return error for deleted post")
	})
}
