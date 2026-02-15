package utils

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yourusername/sns-backend/internal/config"
)

func TestGenerateToken(t *testing.T) {
	// テスト用のconfig設定
	config.AppConfig = &config.Config{
		JWTSecret: "test-secret-key",
	}

	t.Run("Success - Generate valid token", func(t *testing.T) {
		userID := uint(123)

		token, err := GenerateToken(userID)

		if err != nil {
			t.Fatalf("GenerateToken should not return error: %v", err)
		}

		if token == "" {
			t.Fatal("Token should not be empty")
		}
	})

	t.Run("Success - Token contains correct user ID", func(t *testing.T) {
		userID := uint(456)

		tokenString, err := GenerateToken(userID)
		if err != nil {
			t.Fatalf("GenerateToken should not return error: %v", err)
		}

		// トークンを検証してユーザーIDを取得
		token, err := ValidateToken(tokenString)
		if err != nil {
			t.Fatalf("ValidateToken should not return error: %v", err)
		}

		extractedUserID, err := ExtractUserID(token)
		if err != nil {
			t.Fatalf("ExtractUserID should not return error: %v", err)
		}

		if extractedUserID != userID {
			t.Fatalf("Expected user ID %d, got %d", userID, extractedUserID)
		}
	})
}

func TestValidateToken(t *testing.T) {
	// テスト用のconfig設定
	config.AppConfig = &config.Config{
		JWTSecret: "test-secret-key",
	}

	t.Run("Success - Validate valid token", func(t *testing.T) {
		userID := uint(789)

		// トークンを生成
		tokenString, err := GenerateToken(userID)
		if err != nil {
			t.Fatalf("GenerateToken should not return error: %v", err)
		}

		// トークンを検証
		token, err := ValidateToken(tokenString)
		if err != nil {
			t.Fatalf("ValidateToken should not return error: %v", err)
		}

		if !token.Valid {
			t.Fatal("Token should be valid")
		}
	})

	t.Run("Error - Invalid token string", func(t *testing.T) {
		invalidToken := "invalid.token.string"

		_, err := ValidateToken(invalidToken)
		if err == nil {
			t.Fatal("ValidateToken should return error for invalid token")
		}
	})

	t.Run("Error - Wrong secret", func(t *testing.T) {
		userID := uint(111)

		// 正しいシークレットでトークンを生成
		tokenString, err := GenerateToken(userID)
		if err != nil {
			t.Fatalf("GenerateToken should not return error: %v", err)
		}

		// 異なるシークレットで検証
		config.AppConfig.JWTSecret = "different-secret"
		_, err = ValidateToken(tokenString)
		if err == nil {
			t.Fatal("ValidateToken should return error for wrong secret")
		}

		// シークレットを元に戻す
		config.AppConfig.JWTSecret = "test-secret-key"
	})

	t.Run("Error - Expired token", func(t *testing.T) {
		// 期限切れトークンを作成
		claims := JWTClaims{
			UserID: 999,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // 1時間前に期限切れ
				IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(config.AppConfig.JWTSecret))
		if err != nil {
			t.Fatalf("Failed to create expired token: %v", err)
		}

		// 期限切れトークンを検証
		_, err = ValidateToken(tokenString)
		if err == nil {
			t.Fatal("ValidateToken should return error for expired token")
		}
	})
}

func TestExtractUserID(t *testing.T) {
	// テスト用のconfig設定
	config.AppConfig = &config.Config{
		JWTSecret: "test-secret-key",
	}

	t.Run("Success - Extract user ID from valid token", func(t *testing.T) {
		userID := uint(555)

		// トークンを生成して検証
		tokenString, err := GenerateToken(userID)
		if err != nil {
			t.Fatalf("GenerateToken should not return error: %v", err)
		}

		token, err := ValidateToken(tokenString)
		if err != nil {
			t.Fatalf("ValidateToken should not return error: %v", err)
		}

		// ユーザーIDを取得
		extractedUserID, err := ExtractUserID(token)
		if err != nil {
			t.Fatalf("ExtractUserID should not return error: %v", err)
		}

		if extractedUserID != userID {
			t.Fatalf("Expected user ID %d, got %d", userID, extractedUserID)
		}
	})

	t.Run("Error - Invalid token claims", func(t *testing.T) {
		// 不正なクレームを持つトークンを作成
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"invalid": "claims",
		})

		_, err := ExtractUserID(token)
		if err == nil {
			t.Fatal("ExtractUserID should return error for invalid claims")
		}
	})
}
