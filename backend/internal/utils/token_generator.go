package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateRandomToken ランダムなトークンを生成（パスワードリセット、メール認証用）
func GenerateRandomToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateResetToken パスワードリセット用トークン生成（64文字）
func GenerateResetToken() (string, error) {
	return GenerateRandomToken(32) // 32バイト = 64文字（hex）
}

// GenerateVerificationToken メール認証用トークン生成（64文字）
func GenerateVerificationToken() (string, error) {
	return GenerateRandomToken(32) // 32バイト = 64文字（hex）
}
