package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yourusername/sns-backend/internal/config"
)

type JWTClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateAccessToken - アクセストークンを生成（有効期限: 1時間）
func GenerateAccessToken(userID uint) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)), // 1時間有効
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.AppConfig.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GenerateToken - 後方互換性のため残す（内部的にGenerateAccessTokenを呼ぶ）
func GenerateToken(userID uint) (string, error) {
	return GenerateAccessToken(userID)
}

// ValidateToken - JWTトークンを検証
func ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 署名アルゴリズムの確認
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(config.AppConfig.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token, nil
}

// ExtractUserID - トークンからユーザーIDを取得
func ExtractUserID(token *jwt.Token) (uint, error) {
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}
	return claims.UserID, nil
}
