package services

import (
	"testing"

	"github.com/yourusername/sns-backend/internal/database"
	"github.com/yourusername/sns-backend/internal/testutil"
)

func TestRegister(t *testing.T) {
	// テストDBのセットアップ
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, db)
	defer testutil.CleanupTestDB(t, db)

	// グローバルDBを設定（サービスがdatabase.GetDB()を使用するため）
	database.DB = db

	t.Run("Success - Register new user", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		email := "test@example.com"
		password := "password123"
		username := "testuser"

		user, err := Register(email, password, username)

		testutil.AssertNoError(t, err, "Register should not return error")
		testutil.AssertEqual(t, email, user.Email, "Email should match")
		testutil.AssertEqual(t, username, user.Username, "Username should match")
		testutil.AssertNotEqual(t, password, user.Password, "Password should be hashed")
	})

	t.Run("Error - Duplicate email", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		email := "duplicate@example.com"
		password := "password123"
		username := "user1"

		// 最初のユーザーを作成
		_, err := Register(email, password, username)
		testutil.AssertNoError(t, err, "First registration should succeed")

		// 同じメールで再登録
		_, err = Register(email, password, "user2")
		testutil.AssertError(t, err, "Should return error for duplicate email")
	})

	t.Run("Error - Duplicate username", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		password := "password123"
		username := "duplicateuser"

		// 最初のユーザーを作成
		_, err := Register("user1@example.com", password, username)
		testutil.AssertNoError(t, err, "First registration should succeed")

		// 同じユーザー名で再登録
		_, err = Register("user2@example.com", password, username)
		testutil.AssertError(t, err, "Should return error for duplicate username")
	})
}

func TestLogin(t *testing.T) {
	// テストDBのセットアップ
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, db)
	defer testutil.CleanupTestDB(t, db)

	// グローバルDBを設定
	database.DB = db

	t.Run("Success - Login with valid credentials", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		email := "login@example.com"
		password := "password123"
		username := "loginuser"

		// ユーザーを作成
		_, err := Register(email, password, username)
		testutil.AssertNoError(t, err, "User registration should succeed")

		// ログイン
		user, err := Login(email, password)
		testutil.AssertNoError(t, err, "Login should not return error")
		testutil.AssertEqual(t, email, user.Email, "Email should match")
		testutil.AssertEqual(t, username, user.Username, "Username should match")
	})

	t.Run("Error - Invalid email", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		_, err := Login("nonexistent@example.com", "password123")
		testutil.AssertError(t, err, "Should return error for invalid email")
	})

	t.Run("Error - Invalid password", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		email := "wrongpass@example.com"
		password := "correctpassword"
		username := "wrongpassuser"

		// ユーザーを作成
		_, err := Register(email, password, username)
		testutil.AssertNoError(t, err, "User registration should succeed")

		// 間違ったパスワードでログイン
		_, err = Login(email, "wrongpassword")
		testutil.AssertError(t, err, "Should return error for invalid password")
	})
}

func TestGetCurrentUser(t *testing.T) {
	// テストDBのセットアップ
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, db)
	defer testutil.CleanupTestDB(t, db)

	// グローバルDBを設定
	database.DB = db

	t.Run("Success - Get existing user", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		email := "getuser@example.com"
		password := "password123"
		username := "getusertest"

		// ユーザーを作成
		createdUser, err := Register(email, password, username)
		testutil.AssertNoError(t, err, "User registration should succeed")

		// ユーザー情報を取得
		user, err := GetCurrentUser(createdUser.ID)
		testutil.AssertNoError(t, err, "GetCurrentUser should not return error")
		testutil.AssertEqual(t, email, user.Email, "Email should match")
		testutil.AssertEqual(t, username, user.Username, "Username should match")
	})

	t.Run("Error - User not found", func(t *testing.T) {
		testutil.CleanupTestDB(t, db)

		_, err := GetCurrentUser(99999)
		testutil.AssertError(t, err, "Should return error for non-existent user")
	})
}
