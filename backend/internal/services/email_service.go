package services

import (
	"context"
	"fmt"

	"github.com/resend/resend-go/v2"
	"github.com/yourusername/sns-backend/internal/config"
)

// EmailService メール送信サービス
type EmailService struct {
	client *resend.Client
}

// NewEmailService EmailServiceのコンストラクタ
func NewEmailService() *EmailService {
	cfg := config.AppConfig
	if cfg.ResendAPIKey == "" {
		return nil
	}

	client := resend.NewClient(cfg.ResendAPIKey)
	return &EmailService{
		client: client,
	}
}

// SendPasswordResetEmail パスワードリセットメールを送信
func (s *EmailService) SendPasswordResetEmail(ctx context.Context, toEmail, token string) error {
	cfg := config.AppConfig
	resetLink := fmt.Sprintf("%s/auth/password-reset/confirm?token=%s", cfg.FrontendURL, token)

	// 開発・テスト環境ではログ出力のみ（本番環境では実際に送信）
	if cfg.Env != "production" {
		fmt.Printf("\n🔑 [DEV] Password Reset Link for %s:\n%s\n\n", toEmail, resetLink)
		return nil
	}

	if s.client == nil {
		return fmt.Errorf("email service not configured")
	}

	htmlBody := fmt.Sprintf(`
		<html>
		<head>
			<style>
				body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
				.container { max-width: 600px; margin: 0 auto; padding: 20px; }
				.header { background-color: #1976d2; color: white; padding: 20px; text-align: center; }
				.content { background-color: #f9f9f9; padding: 30px; }
				.button { background-color: #1976d2; color: white; padding: 12px 24px; text-decoration: none; border-radius: 4px; display: inline-block; margin-top: 20px; }
				.footer { text-align: center; padding: 20px; color: #777; font-size: 12px; }
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header">
					<h1>パスワードリセット</h1>
				</div>
				<div class="content">
					<p>パスワードリセットのリクエストを受け付けました。</p>
					<p>以下のリンクをクリックして、新しいパスワードを設定してください：</p>
					<a href="%s" class="button">パスワードをリセット</a>
					<p style="margin-top: 20px; color: #777; font-size: 14px;">
						このリンクは1時間後に無効になります。<br>
						もしこのメールに心当たりがない場合は、無視してください。
					</p>
				</div>
				<div class="footer">
					<p>&copy; 2026 SNS App. All rights reserved.</p>
				</div>
			</div>
		</body>
		</html>
	`, resetLink)

	params := &resend.SendEmailRequest{
		From:    cfg.FromEmail,
		To:      []string{toEmail},
		Subject: "パスワードリセットのお知らせ",
		Html:    htmlBody,
		Text:    fmt.Sprintf("パスワードリセットリンク: %s\n\nこのリンクは1時間後に無効になります。", resetLink),
	}

	_, err := s.client.Emails.SendWithContext(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}

// SendVerificationEmail メール認証メールを送信
func (s *EmailService) SendVerificationEmail(ctx context.Context, toEmail, token string) error {
	cfg := config.AppConfig
	verificationLink := fmt.Sprintf("%s/auth/email/verify?token=%s", cfg.FrontendURL, token)

	// 開発・テスト環境ではログ出力のみ（本番環境では実際に送信）
	if cfg.Env != "production" {
		fmt.Printf("\n🔗 [DEV] Email Verification Link for %s:\n%s\n\n", toEmail, verificationLink)
		return nil
	}

	if s.client == nil {
		return fmt.Errorf("email service not configured")
	}

	htmlBody := fmt.Sprintf(`
		<html>
		<head>
			<style>
				body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
				.container { max-width: 600px; margin: 0 auto; padding: 20px; }
				.header { background-color: #1976d2; color: white; padding: 20px; text-align: center; }
				.content { background-color: #f9f9f9; padding: 30px; }
				.button { background-color: #1976d2; color: white; padding: 12px 24px; text-decoration: none; border-radius: 4px; display: inline-block; margin-top: 20px; }
				.footer { text-align: center; padding: 20px; color: #777; font-size: 12px; }
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header">
					<h1>メールアドレス認証</h1>
				</div>
				<div class="content">
					<p>SNS Appへようこそ！</p>
					<p>アカウント登録ありがとうございます。以下のリンクをクリックしてメールアドレスを認証してください：</p>
					<a href="%s" class="button">メールアドレスを認証</a>
					<p style="margin-top: 20px; color: #777; font-size: 14px;">
						このリンクは24時間後に無効になります。<br>
						もしこのメールに心当たりがない場合は、無視してください。
					</p>
				</div>
				<div class="footer">
					<p>&copy; 2026 SNS App. All rights reserved.</p>
				</div>
			</div>
		</body>
		</html>
	`, verificationLink)

	params := &resend.SendEmailRequest{
		From:    cfg.FromEmail,
		To:      []string{toEmail},
		Subject: "メールアドレス認証のお知らせ",
		Html:    htmlBody,
		Text:    fmt.Sprintf("メールアドレス認証リンク: %s\n\nこのリンクは24時間後に無効になります。", verificationLink),
	}

	_, err := s.client.Emails.SendWithContext(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}

// IsEmailServiceConfigured メール送信サービスが設定されているか確認
func IsEmailServiceConfigured() bool {
	cfg := config.AppConfig
	return cfg.ResendAPIKey != ""
}
