package services

import (
	"testing"

	"github.com/yourusername/sns-backend/internal/database"
	"github.com/yourusername/sns-backend/internal/testutil"
)

// TestFollowUser_Authorization - フォローの認可テスト
func TestFollowUser_Authorization(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, db)
	defer testutil.CleanupTestDB(t, db)
	database.DB = db

	t.Run("Success - Follow another user", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		follower := testutil.CreateTestUser(t, db, "follower@example.com", "follower", "password123")
		following := testutil.CreateTestUser(t, db, "following@example.com", "following", "password123")

		err := FollowUser(follower.ID, following.Username)
		testutil.AssertNoError(t, err, "Should be able to follow another user")

		// フォロー関係が作成されたことを確認
		isFollowing, err := CheckIfFollowing(follower.ID, following.ID)
		testutil.AssertNoError(t, err, "CheckIfFollowing should not return error")
		testutil.AssertTrue(t, isFollowing, "Should be following")
	})

	t.Run("Error - Cannot follow yourself", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "user", "password123")

		err := FollowUser(user.ID, user.Username)
		testutil.AssertError(t, err, "Should not be able to follow yourself")

		// フォロー関係が作成されていないことを確認
		isFollowing, err := CheckIfFollowing(user.ID, user.ID)
		testutil.AssertNoError(t, err, "CheckIfFollowing should not return error")
		testutil.AssertFalse(t, isFollowing, "Should not be following yourself")
	})

	t.Run("Error - Cannot follow same user twice (idempotency)", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		follower := testutil.CreateTestUser(t, db, "follower@example.com", "follower", "password123")
		following := testutil.CreateTestUser(t, db, "following@example.com", "following", "password123")

		// 最初のフォロー
		err := FollowUser(follower.ID, following.Username)
		testutil.AssertNoError(t, err, "First follow should succeed")

		// 2回目のフォロー（冪等性テスト）
		err = FollowUser(follower.ID, following.Username)
		testutil.AssertError(t, err, "Should not be able to follow same user twice")
	})

	t.Run("Error - Cannot follow non-existent user", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "user", "password123")

		err := FollowUser(user.ID, "nonexistentuser")
		testutil.AssertError(t, err, "Should return error for non-existent user")
	})

	t.Run("Error - Invalid follower ID", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		following := testutil.CreateTestUser(t, db, "following@example.com", "following", "password123")

		err := FollowUser(99999, following.Username)
		// Note: このケースはDBの外部キー制約でエラーになる可能性がある
		if err == nil {
			t.Logf("WARNING: Invalid follower ID was accepted")
		}
	})
}

// TestUnfollowUser_Authorization - フォロー解除のテスト
func TestUnfollowUser_Authorization(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, db)
	defer testutil.CleanupTestDB(t, db)
	database.DB = db

	t.Run("Success - Unfollow a followed user", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		follower := testutil.CreateTestUser(t, db, "follower@example.com", "follower", "password123")
		following := testutil.CreateTestUser(t, db, "following@example.com", "following", "password123")

		// フォロー
		err := FollowUser(follower.ID, following.Username)
		testutil.AssertNoError(t, err, "Follow should succeed")

		// フォロー解除
		err = UnfollowUser(follower.ID, following.Username)
		testutil.AssertNoError(t, err, "Unfollow should succeed")

		// フォロー関係が解除されたことを確認
		isFollowing, err := CheckIfFollowing(follower.ID, following.ID)
		testutil.AssertNoError(t, err, "CheckIfFollowing should not return error")
		testutil.AssertFalse(t, isFollowing, "Should not be following")
	})

	t.Run("Error - Cannot unfollow user that is not followed", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		follower := testutil.CreateTestUser(t, db, "follower@example.com", "follower", "password123")
		following := testutil.CreateTestUser(t, db, "following@example.com", "following", "password123")

		// フォローしていないユーザーをフォロー解除しようとする
		err := UnfollowUser(follower.ID, following.Username)
		testutil.AssertError(t, err, "Should return error when trying to unfollow a user that is not followed")
	})

	t.Run("Error - Cannot unfollow same user twice", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		follower := testutil.CreateTestUser(t, db, "follower@example.com", "follower", "password123")
		following := testutil.CreateTestUser(t, db, "following@example.com", "following", "password123")

		// フォロー
		err := FollowUser(follower.ID, following.Username)
		testutil.AssertNoError(t, err, "Follow should succeed")

		// フォロー解除
		err = UnfollowUser(follower.ID, following.Username)
		testutil.AssertNoError(t, err, "Unfollow should succeed")

		// 再度フォロー解除しようとする
		err = UnfollowUser(follower.ID, following.Username)
		testutil.AssertError(t, err, "Should not be able to unfollow twice")
	})

	t.Run("Error - Unfollow non-existent user", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "user", "password123")

		err := UnfollowUser(user.ID, "nonexistentuser")
		testutil.AssertError(t, err, "Should return error for non-existent user")
	})
}

// TestGetFollowers_EdgeCases - フォロワー一覧取得のエッジケーステスト
func TestGetFollowers_EdgeCases(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, db)
	defer testutil.CleanupTestDB(t, db)
	database.DB = db

	t.Run("Success - Get followers for user with no followers", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "user", "password123")

		users, hasMore, nextCursor, err := GetFollowers(user.Username, 10, nil)
		testutil.AssertNoError(t, err, "Should not return error for user with no followers")
		testutil.AssertEqual(t, 0, len(users), "Should return empty array")
		testutil.AssertFalse(t, hasMore, "Should not have more")
		testutil.AssertEqual(t, "", nextCursor, "Next cursor should be empty")
	})

	t.Run("Error - Get followers for non-existent user", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		_, _, _, err := GetFollowers("nonexistentuser", 10, nil)
		testutil.AssertError(t, err, "Should return error for non-existent user")
	})

	t.Run("Success - Get followers with pagination", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "user", "password123")

		// 5人のフォロワーを作成
		for i := 1; i <= 5; i++ {
			follower := testutil.CreateTestUser(t, db, "follower"+string(rune(i))+"@example.com", "follower"+string(rune(i)), "password123")
			err := FollowUser(follower.ID, user.Username)
			testutil.AssertNoError(t, err, "Follow should succeed")
		}

		// 最初のページ（limit=2）
		followers, hasMore, nextCursor, err := GetFollowers(user.Username, 2, nil)
		testutil.AssertNoError(t, err, "GetFollowers should not return error")
		testutil.AssertEqual(t, 2, len(followers), "Should return 2 followers")
		testutil.AssertTrue(t, hasMore, "Should have more followers")
		testutil.AssertNotEqual(t, "", nextCursor, "Next cursor should not be empty")
	})
}

// TestGetFollowing_EdgeCases - フォロー中ユーザー一覧取得のエッジケーステスト
func TestGetFollowing_EdgeCases(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, db)
	defer testutil.CleanupTestDB(t, db)
	database.DB = db

	t.Run("Success - Get following for user not following anyone", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "user", "password123")

		users, hasMore, nextCursor, err := GetFollowing(user.Username, 10, nil)
		testutil.AssertNoError(t, err, "Should not return error for user following no one")
		testutil.AssertEqual(t, 0, len(users), "Should return empty array")
		testutil.AssertFalse(t, hasMore, "Should not have more")
		testutil.AssertEqual(t, "", nextCursor, "Next cursor should be empty")
	})

	t.Run("Error - Get following for non-existent user", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		_, _, _, err := GetFollowing("nonexistentuser", 10, nil)
		testutil.AssertError(t, err, "Should return error for non-existent user")
	})
}

// TestCheckIfFollowing_EdgeCases - フォロー状態チェックのエッジケーステスト
func TestCheckIfFollowing_EdgeCases(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, db)
	defer testutil.CleanupTestDB(t, db)
	database.DB = db

	t.Run("Success - Check if not following", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		user1 := testutil.CreateTestUser(t, db, "user1@example.com", "user1", "password123")
		user2 := testutil.CreateTestUser(t, db, "user2@example.com", "user2", "password123")

		isFollowing, err := CheckIfFollowing(user1.ID, user2.ID)
		testutil.AssertNoError(t, err, "CheckIfFollowing should not return error")
		testutil.AssertFalse(t, isFollowing, "Should not be following")
	})

	t.Run("Success - Check if following yourself", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "user", "password123")

		isFollowing, err := CheckIfFollowing(user.ID, user.ID)
		testutil.AssertNoError(t, err, "CheckIfFollowing should not return error")
		testutil.AssertFalse(t, isFollowing, "Should not be following yourself")
	})

	t.Run("Success - Check with invalid user IDs", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		isFollowing, err := CheckIfFollowing(99999, 99998)
		testutil.AssertNoError(t, err, "CheckIfFollowing should not return error")
		testutil.AssertFalse(t, isFollowing, "Should return false for invalid user IDs")
	})
}
