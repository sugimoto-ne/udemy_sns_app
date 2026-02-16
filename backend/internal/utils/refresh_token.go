package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"time"

	"github.com/yourusername/sns-backend/internal/database"
	"github.com/yourusername/sns-backend/internal/models"
)

const (
	RefreshTokenLength     = 32                // トークンの長さ（バイト）
	RefreshTokenExpiration = 7 * 24 * time.Hour // 7日間
)

// GenerateRefreshToken - リフレッシュトークンを生成してDBに保存
func GenerateRefreshToken(userID uint) (string, error) {
	// ランダムなトークンを生成
	tokenBytes := make([]byte, RefreshTokenLength)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}

	// Base64エンコード
	tokenString := base64.URLEncoding.EncodeToString(tokenBytes)

	// トークンのハッシュを計算（DBに保存）
	hashedToken := hashToken(tokenString)

	// DBに保存
	refreshToken := models.RefreshToken{
		UserID:    userID,
		Token:     hashedToken,
		ExpiresAt: time.Now().Add(RefreshTokenExpiration),
		Revoked:   false,
	}

	if err := database.DB.Create(&refreshToken).Error; err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateRefreshToken - リフレッシュトークンを検証
func ValidateRefreshToken(tokenString string) (*models.RefreshToken, error) {
	// トークンのハッシュを計算
	hashedToken := hashToken(tokenString)

	// DBからトークンを検索
	var refreshToken models.RefreshToken
	err := database.DB.Where("token = ?", hashedToken).First(&refreshToken).Error
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// トークンが有効かチェック
	if !refreshToken.IsValid() {
		return nil, errors.New("refresh token is expired or revoked")
	}

	return &refreshToken, nil
}

// RevokeRefreshToken - リフレッシュトークンを無効化
func RevokeRefreshToken(tokenString string) error {
	hashedToken := hashToken(tokenString)

	result := database.DB.Model(&models.RefreshToken{}).
		Where("token = ?", hashedToken).
		Update("revoked", true)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("refresh token not found")
	}

	return nil
}

// RevokeAllUserTokens - ユーザーのすべてのリフレッシュトークンを無効化
func RevokeAllUserTokens(userID uint) error {
	result := database.DB.Model(&models.RefreshToken{}).
		Where("user_id = ? AND revoked = ?", userID, false).
		Update("revoked", true)

	return result.Error
}

// CleanupExpiredTokens - 期限切れトークンをDBから削除（定期実行用）
func CleanupExpiredTokens() error {
	result := database.DB.Where("expires_at < ?", time.Now()).
		Delete(&models.RefreshToken{})

	return result.Error
}

// hashToken - トークンをSHA256でハッシュ化
func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(hash[:])
}
