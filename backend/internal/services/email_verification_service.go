package services

import (
	"context"
	"errors"
	"time"

	"github.com/yourusername/sns-backend/internal/database"
	"github.com/yourusername/sns-backend/internal/models"
	"github.com/yourusername/sns-backend/internal/utils"
	"gorm.io/gorm"
)

// EmailVerificationService メール認証サービス
type EmailVerificationService struct {
	db           *gorm.DB
	emailService *EmailService
}

// NewEmailVerificationService EmailVerificationServiceのコンストラクタ
func NewEmailVerificationService() *EmailVerificationService {
	return &EmailVerificationService{
		db:           database.GetDB(),
		emailService: NewEmailService(),
	}
}

// SendVerificationEmail 認証メールを送信
func (s *EmailVerificationService) SendVerificationEmail(ctx context.Context, userID uint) error {
	// ユーザー情報取得
	var user models.User
	if err := s.db.WithContext(ctx).First(&user, userID).Error; err != nil {
		return err
	}

	// 既に認証済みの場合はスキップ
	if user.EmailVerified {
		return errors.New("email already verified")
	}

	// トークン生成
	token, err := utils.GenerateVerificationToken()
	if err != nil {
		return err
	}

	// 有効期限設定（24時間）
	expiresAt := time.Now().Add(24 * time.Hour)

	// トークン保存
	verificationToken := &models.EmailVerificationToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
	}

	if err := s.db.WithContext(ctx).Create(verificationToken).Error; err != nil {
		return err
	}

	// メール送信
	if s.emailService != nil {
		if err := s.emailService.SendVerificationEmail(ctx, user.Email, token); err != nil {
			// メール送信失敗時はエラーを返す
			return err
		}
	}

	return nil
}

// VerifyEmail メールアドレスを認証
func (s *EmailVerificationService) VerifyEmail(ctx context.Context, token string) error {
	// トークン検証
	var verificationToken models.EmailVerificationToken
	if err := s.db.WithContext(ctx).Where("token = ?", token).First(&verificationToken).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("invalid or expired token")
		}
		return err
	}

	// 有効期限チェック
	if time.Now().After(verificationToken.ExpiresAt) {
		return errors.New("token has expired")
	}

	// ユーザーのemail_verifiedをtrueに更新
	if err := s.db.WithContext(ctx).Model(&models.User{}).
		Where("id = ?", verificationToken.UserID).
		Update("email_verified", true).Error; err != nil {
		return err
	}

	// トークン削除（使用済み）
	if err := s.db.WithContext(ctx).Delete(&verificationToken).Error; err != nil {
		return err
	}

	return nil
}

// ResendVerificationEmail 認証メールを再送信
func (s *EmailVerificationService) ResendVerificationEmail(ctx context.Context, userID uint) error {
	// 既存トークン削除
	s.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.EmailVerificationToken{})

	// 新規トークン送信
	return s.SendVerificationEmail(ctx, userID)
}
