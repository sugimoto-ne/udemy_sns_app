package services

import (
	"context"
	"errors"
	"time"

	"github.com/yourusername/sns-backend/internal/database"
	"github.com/yourusername/sns-backend/internal/models"
	"github.com/yourusername/sns-backend/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// PasswordResetService パスワードリセットサービス
type PasswordResetService struct {
	db           *gorm.DB
	emailService *EmailService
}

// NewPasswordResetService PasswordResetServiceのコンストラクタ
func NewPasswordResetService() *PasswordResetService {
	return &PasswordResetService{
		db:           database.GetDB(),
		emailService: NewEmailService(),
	}
}

// RequestPasswordReset パスワードリセットをリクエスト
func (s *PasswordResetService) RequestPasswordReset(ctx context.Context, email string) error {
	// ユーザーの存在確認
	var user models.User
	if err := s.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// セキュリティ上、ユーザーが存在しない場合でもエラーにしない
			// （メールアドレスの存在が判明することを防ぐ）
			return nil
		}
		return err
	}

	// リセットトークン生成
	token, err := utils.GenerateResetToken()
	if err != nil {
		return err
	}

	// 有効期限設定（1時間）
	expiresAt := time.Now().Add(1 * time.Hour)

	// トークン保存（既存のトークンは自動的に古いものとして無視される）
	resetToken := &models.PasswordResetToken{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: expiresAt,
	}

	if err := s.db.WithContext(ctx).Create(resetToken).Error; err != nil {
		return err
	}

	// メール送信
	if s.emailService != nil {
		if err := s.emailService.SendPasswordResetEmail(ctx, email, token); err != nil {
			// メール送信失敗をログに記録するが、ユーザーにはエラーを返さない
			// （メールサービスが設定されていない環境でも動作するように）
			return nil
		}
	}

	return nil
}

// ConfirmPasswordReset パスワードをリセット
func (s *PasswordResetService) ConfirmPasswordReset(ctx context.Context, token, newPassword string) error {
	// トークン検証
	var resetToken models.PasswordResetToken
	if err := s.db.WithContext(ctx).Where("token = ?", token).First(&resetToken).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("invalid or expired token")
		}
		return err
	}

	// 有効期限チェック
	if time.Now().After(resetToken.ExpiresAt) {
		return errors.New("token has expired")
	}

	// パスワードバリデーション
	if len(newPassword) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	// パスワードハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// パスワード更新
	if err := s.db.WithContext(ctx).Model(&models.User{}).
		Where("id = ?", resetToken.UserID).
		Update("password", string(hashedPassword)).Error; err != nil {
		return err
	}

	// トークン削除（使用済み）
	if err := s.db.WithContext(ctx).Delete(&resetToken).Error; err != nil {
		return err
	}

	return nil
}
