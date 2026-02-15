package testutil

import (
	"testing"

	"github.com/yourusername/sns-backend/internal/models"
	"gorm.io/gorm"
)

// Post は論理削除の確認用
type Post struct {
	gorm.Model
	DeletedAt gorm.DeletedAt
}

// CreateTestUser creates a test user in the database
func CreateTestUser(t *testing.T, db *gorm.DB, email, username, password string) *models.User {
	t.Helper()

	user := &models.User{
		Email:    email,
		Username: username,
		Password: password, // BeforeCreateフックでハッシュ化される
	}

	if err := db.Create(user).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	return user
}

// CreateTestPost creates a test post in the database
func CreateTestPost(t *testing.T, db *gorm.DB, userID uint, content string) *models.Post {
	t.Helper()

	post := &models.Post{
		UserID:  userID,
		Content: content,
	}

	if err := db.Create(post).Error; err != nil {
		t.Fatalf("Failed to create test post: %v", err)
	}

	return post
}

// CreateTestComment creates a test comment in the database
func CreateTestComment(t *testing.T, db *gorm.DB, postID, userID uint, content string) *models.Comment {
	t.Helper()

	comment := &models.Comment{
		PostID:  postID,
		UserID:  userID,
		Content: content,
	}

	if err := db.Create(comment).Error; err != nil {
		t.Fatalf("Failed to create test comment: %v", err)
	}

	return comment
}

// CreateTestLike creates a test like in the database
func CreateTestLike(t *testing.T, db *gorm.DB, postID, userID uint) *models.PostLike {
	t.Helper()

	like := &models.PostLike{
		PostID: postID,
		UserID: userID,
	}

	if err := db.Create(like).Error; err != nil {
		t.Fatalf("Failed to create test like: %v", err)
	}

	return like
}

// CreateTestFollow creates a test follow relationship in the database
func CreateTestFollow(t *testing.T, db *gorm.DB, followerID, followingID uint) *models.Follow {
	t.Helper()

	follow := &models.Follow{
		FollowerID:  followerID,
		FollowingID: followingID,
	}

	if err := db.Create(follow).Error; err != nil {
		t.Fatalf("Failed to create test follow: %v", err)
	}

	return follow
}

// AssertNoError asserts that there is no error
func AssertNoError(t *testing.T, err error, msg string) {
	t.Helper()
	if err != nil {
		t.Fatalf("%s: %v", msg, err)
	}
}

// AssertError asserts that there is an error
func AssertError(t *testing.T, err error, msg string) {
	t.Helper()
	if err == nil {
		t.Fatalf("%s: expected error but got nil", msg)
	}
}

// AssertEqual asserts that two values are equal
func AssertEqual(t *testing.T, expected, actual interface{}, msg string) {
	t.Helper()
	if expected != actual {
		t.Fatalf("%s: expected %v, got %v", msg, expected, actual)
	}
}

// AssertNotEqual asserts that two values are not equal
func AssertNotEqual(t *testing.T, expected, actual interface{}, msg string) {
	t.Helper()
	if expected == actual {
		t.Fatalf("%s: expected %v to not equal %v", msg, expected, actual)
	}
}

// AssertTrue asserts that a condition is true
func AssertTrue(t *testing.T, condition bool, msg string) {
	t.Helper()
	if !condition {
		t.Fatalf("%s: expected true but got false", msg)
	}
}

// AssertFalse asserts that a condition is false
func AssertFalse(t *testing.T, condition bool, msg string) {
	t.Helper()
	if condition {
		t.Fatalf("%s: expected false but got true", msg)
	}
}
