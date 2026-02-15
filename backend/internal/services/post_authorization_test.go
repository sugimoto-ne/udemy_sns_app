package services

import (
	"testing"

	"github.com/yourusername/sns-backend/internal/database"
	"github.com/yourusername/sns-backend/internal/testutil"
)

// TestUpdatePost_Authorization - 投稿更新の認可テスト
func TestUpdatePost_Authorization(t *testing.T) {
	// テストDBのセットアップ
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, db)
	defer testutil.CleanupTestDB(t, db)

	// グローバルDBを設定
	database.DB = db

	t.Run("Success - Owner can update own post", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		// ユーザーと投稿を作成
		user := testutil.CreateTestUser(t, db, "owner@example.com", "owner", "password123")
		post := testutil.CreateTestPost(t, db, user.ID, "Original content")

		// 投稿者が自分の投稿を更新
		updatedContent := "Updated content"
		updatedPost, err := UpdatePost(post.ID, user.ID, updatedContent)

		testutil.AssertNoError(t, err, "Owner should be able to update own post")
		testutil.AssertEqual(t, updatedContent, updatedPost.Content, "Content should be updated")
	})

	t.Run("Error - Non-owner cannot update post", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		// 2人のユーザーを作成
		owner := testutil.CreateTestUser(t, db, "owner@example.com", "owner", "password123")
		otherUser := testutil.CreateTestUser(t, db, "other@example.com", "other", "password123")

		// 投稿を作成
		post := testutil.CreateTestPost(t, db, owner.ID, "Original content")

		// 他のユーザーが投稿を更新しようとする
		_, err := UpdatePost(post.ID, otherUser.ID, "Hacked content")

		testutil.AssertError(t, err, "Non-owner should not be able to update post")
	})

	t.Run("Error - Update non-existent post", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "user", "password123")

		// 存在しない投稿を更新しようとする
		_, err := UpdatePost(99999, user.ID, "Content")

		testutil.AssertError(t, err, "Should return error for non-existent post")
	})
}

// TestDeletePost_Authorization - 投稿削除の認可テスト
func TestDeletePost_Authorization(t *testing.T) {
	// テストDBのセットアップ
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, db)
	defer testutil.CleanupTestDB(t, db)

	// グローバルDBを設定
	database.DB = db

	t.Run("Success - Owner can delete own post", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		// ユーザーと投稿を作成
		user := testutil.CreateTestUser(t, db, "owner@example.com", "owner", "password123")
		post := testutil.CreateTestPost(t, db, user.ID, "Content to delete")

		// 投稿者が自分の投稿を削除
		err := DeletePost(post.ID, user.ID)

		testutil.AssertNoError(t, err, "Owner should be able to delete own post")

		// 削除されたことを確認（論理削除なのでUnscopedで確認）
		var deletedPost testutil.Post
		result := db.Unscoped().First(&deletedPost, post.ID)
		testutil.AssertNoError(t, result.Error, "Post should exist in database (soft deleted)")
		testutil.AssertTrue(t, deletedPost.DeletedAt.Valid, "Post should be soft deleted")
	})

	t.Run("Error - Non-owner cannot delete post", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		// 2人のユーザーを作成
		owner := testutil.CreateTestUser(t, db, "owner@example.com", "owner", "password123")
		otherUser := testutil.CreateTestUser(t, db, "other@example.com", "other", "password123")

		// 投稿を作成
		post := testutil.CreateTestPost(t, db, owner.ID, "Content")

		// 他のユーザーが投稿を削除しようとする
		err := DeletePost(post.ID, otherUser.ID)

		testutil.AssertError(t, err, "Non-owner should not be able to delete post")

		// 投稿が削除されていないことを確認
		var stillExistingPost testutil.Post
		result := db.First(&stillExistingPost, post.ID)
		testutil.AssertNoError(t, result.Error, "Post should still exist")
	})

	t.Run("Error - Delete non-existent post", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		user := testutil.CreateTestUser(t, db, "user@example.com", "user", "password123")

		// 存在しない投稿を削除しようとする
		err := DeletePost(99999, user.ID)

		testutil.AssertError(t, err, "Should return error for non-existent post")
	})
}
