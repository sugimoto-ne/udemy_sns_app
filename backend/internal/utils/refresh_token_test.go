package utils

import (
	"os"
	"testing"
	"time"

	"github.com/yourusername/sns-backend/internal/config"
	"github.com/yourusername/sns-backend/internal/database"
	"github.com/yourusername/sns-backend/internal/models"
)

// テスト前の環境チェック
func skipIfNoTestDB(t *testing.T) {
	if database.DB == nil {
		t.Skip("Skipping test: database connection not available")
	}
}

// テスト用の設定初期化
func setupRefreshTokenTestConfig(t *testing.T) {
	// 環境変数から取得、なければデフォルト値
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "test-secret-key-for-testing-only"
	}

	config.AppConfig = &config.Config{
		JWTSecret: jwtSecret,
	}
}

// テスト用のユーザーをクリーンアップ
func cleanupTestUser(t *testing.T, userID uint) {
	// リフレッシュトークンを削除
	database.DB.Unscoped().Where("user_id = ?", userID).Delete(&models.RefreshToken{})
	// ユーザーを削除
	database.DB.Unscoped().Delete(&models.User{}, userID)
}

// TestGenerateRefreshToken - リフレッシュトークン生成テスト
func TestGenerateRefreshToken(t *testing.T) {
	skipIfNoTestDB(t)
	setupRefreshTokenTestConfig(t)

	// テストユーザー作成
	user := models.User{
		Username: "testuser_gen",
		Email:    "testgen@example.com",
		Password: "hashedpassword",
	}
	database.DB.Create(&user)
	defer cleanupTestUser(t, user.ID)

	// リフレッシュトークン生成
	token, err := GenerateRefreshToken(user.ID)
	if err != nil {
		t.Fatalf("GenerateRefreshToken failed: %v", err)
	}

	// トークンが空でないことを確認
	if token == "" {
		t.Error("Generated token is empty")
	}

	// DBにトークンが保存されていることを確認
	var refreshToken models.RefreshToken
	result := database.DB.Where("user_id = ?", user.ID).First(&refreshToken)
	if result.Error != nil {
		t.Fatalf("Failed to find refresh token in DB: %v", result.Error)
	}

	// トークンが有効期限内であることを確認
	if refreshToken.ExpiresAt.Before(time.Now()) {
		t.Error("Refresh token is already expired")
	}

	// トークンが失効していないことを確認
	if refreshToken.Revoked {
		t.Error("Newly generated token should not be revoked")
	}
}

// TestValidateRefreshToken - リフレッシュトークン検証テスト
func TestValidateRefreshToken(t *testing.T) {
	skipIfNoTestDB(t)
	setupRefreshTokenTestConfig(t)

	// テストユーザー作成
	user := models.User{
		Username: "testuser_val",
		Email:    "testval@example.com",
		Password: "hashedpassword",
	}
	database.DB.Create(&user)
	defer cleanupTestUser(t, user.ID)

	// リフレッシュトークン生成
	token, err := GenerateRefreshToken(user.ID)
	if err != nil {
		t.Fatalf("GenerateRefreshToken failed: %v", err)
	}

	// トークン検証（正常系）
	tokenRecord, err := ValidateRefreshToken(token)
	if err != nil {
		t.Fatalf("ValidateRefreshToken failed: %v", err)
	}

	// トークンレコードが正しいことを確認
	if tokenRecord.UserID != user.ID {
		t.Errorf("UserID mismatch: got %d, want %d", tokenRecord.UserID, user.ID)
	}

	if tokenRecord.Revoked {
		t.Error("Token should not be revoked")
	}

	if !tokenRecord.IsValid() {
		t.Error("Token should be valid")
	}
}

// TestValidateRefreshToken_Invalid - 無効なトークンの検証テスト
func TestValidateRefreshToken_Invalid(t *testing.T) {
	skipIfNoTestDB(t)
	setupRefreshTokenTestConfig(t)

	// 存在しないトークンで検証
	_, err := ValidateRefreshToken("invalid-token-string-12345")
	if err == nil {
		t.Error("Expected error for invalid token, got nil")
	}
}

// TestRevokeRefreshToken - リフレッシュトークン失効テスト
func TestRevokeRefreshToken(t *testing.T) {
	skipIfNoTestDB(t)
	setupRefreshTokenTestConfig(t)

	// テストユーザー作成
	user := models.User{
		Username: "testuser_rev",
		Email:    "testrev@example.com",
		Password: "hashedpassword",
	}
	database.DB.Create(&user)
	defer cleanupTestUser(t, user.ID)

	// リフレッシュトークン生成
	token, err := GenerateRefreshToken(user.ID)
	if err != nil {
		t.Fatalf("GenerateRefreshToken failed: %v", err)
	}

	// トークン失効
	err = RevokeRefreshToken(token)
	if err != nil {
		t.Fatalf("RevokeRefreshToken failed: %v", err)
	}

	// DBでトークンが失効済みになっていることを確認
	var refreshToken models.RefreshToken
	hashedToken := hashToken(token)
	result := database.DB.Where("token = ?", hashedToken).First(&refreshToken)
	if result.Error != nil {
		t.Fatalf("Failed to find refresh token in DB: %v", result.Error)
	}

	if !refreshToken.Revoked {
		t.Error("Token should be revoked")
	}

	// 失効済みトークンの検証が失敗することを確認
	_, err = ValidateRefreshToken(token)
	if err == nil {
		t.Error("Expected error for revoked token, got nil")
	}
}

// TestRevokeAllUserTokens - 全トークン失効テスト
func TestRevokeAllUserTokens(t *testing.T) {
	skipIfNoTestDB(t)
	setupRefreshTokenTestConfig(t)

	// テストユーザー作成
	user := models.User{
		Username: "testuser_revall",
		Email:    "testrevall@example.com",
		Password: "hashedpassword",
	}
	database.DB.Create(&user)
	defer cleanupTestUser(t, user.ID)

	// 複数のリフレッシュトークンを生成
	token1, _ := GenerateRefreshToken(user.ID)
	token2, _ := GenerateRefreshToken(user.ID)
	token3, _ := GenerateRefreshToken(user.ID)

	// 全トークン失効
	err := RevokeAllUserTokens(user.ID)
	if err != nil {
		t.Fatalf("RevokeAllUserTokens failed: %v", err)
	}

	// 全トークンが失効済みになっていることを確認
	tokens := []string{token1, token2, token3}
	for _, token := range tokens {
		_, err := ValidateRefreshToken(token)
		if err == nil {
			t.Errorf("Expected error for revoked token, got nil")
		}
	}

	// DBで全トークンが失効済みになっていることを確認
	var count int64
	database.DB.Model(&models.RefreshToken{}).Where("user_id = ? AND revoked = ?", user.ID, false).Count(&count)
	if count != 0 {
		t.Errorf("Expected 0 non-revoked tokens, got %d", count)
	}
}

// TestTokenRotation - トークンローテーションテスト
func TestTokenRotation(t *testing.T) {
	skipIfNoTestDB(t)
	setupRefreshTokenTestConfig(t)

	// テストユーザー作成
	user := models.User{
		Username: "testuser_rot",
		Email:    "testrot@example.com",
		Password: "hashedpassword",
	}
	database.DB.Create(&user)
	defer cleanupTestUser(t, user.ID)

	// 最初のリフレッシュトークン生成
	oldToken, err := GenerateRefreshToken(user.ID)
	if err != nil {
		t.Fatalf("GenerateRefreshToken failed: %v", err)
	}

	// トークンローテーション: 古いトークンを失効し、新しいトークンを生成
	err = RevokeRefreshToken(oldToken)
	if err != nil {
		t.Fatalf("RevokeRefreshToken failed: %v", err)
	}

	newToken, err := GenerateRefreshToken(user.ID)
	if err != nil {
		t.Fatalf("GenerateRefreshToken failed: %v", err)
	}

	// 古いトークンが無効になっていることを確認
	_, err = ValidateRefreshToken(oldToken)
	if err == nil {
		t.Error("Old token should be invalid after rotation")
	}

	// 新しいトークンが有効であることを確認
	_, err = ValidateRefreshToken(newToken)
	if err != nil {
		t.Errorf("New token should be valid: %v", err)
	}

	// トークンが異なることを確認
	if oldToken == newToken {
		t.Error("Old token and new token should be different")
	}
}
