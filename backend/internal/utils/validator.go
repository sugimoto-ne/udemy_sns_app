package utils

import (
	"errors"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateStruct - 構造体をバリデーション
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

// GetValidator - バリデータインスタンスを取得
func GetValidator() *validator.Validate {
	return validate
}

// バリデーションエラー
var (
	ErrEmptyEmail         = errors.New("email cannot be empty")
	ErrInvalidEmailFormat = errors.New("invalid email format")
	ErrEmailTooLong       = errors.New("email is too long (max 255 characters)")
	ErrEmptyPassword      = errors.New("password cannot be empty")
	ErrPasswordTooShort   = errors.New("password must be at least 6 characters")
	ErrPasswordTooLong    = errors.New("password is too long (max 128 characters)")
	ErrEmptyUsername      = errors.New("username cannot be empty")
	ErrUsernameTooShort   = errors.New("username must be at least 3 characters")
	ErrUsernameTooLong    = errors.New("username is too long (max 30 characters)")
	ErrInvalidUsername    = errors.New("username can only contain letters, numbers, and underscores")
	ErrEmptyContent       = errors.New("content cannot be empty")
	ErrContentTooLong     = errors.New("content is too long")
	ErrInvalidCharacters  = errors.New("content contains invalid characters")
)

// 定数
const (
	MaxEmailLength    = 255
	MinPasswordLength = 6
	MaxPasswordLength = 128
	MinUsernameLength = 3
	MaxUsernameLength = 30
	MaxPostContent    = 280  // Twitterライク
	MaxCommentContent = 500
)

// メールアドレスの正規表現（簡易版）
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// ユーザー名の正規表現（英数字とアンダースコアのみ）
var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

// 危険な文字パターン（SQL injection、XSS）
var dangerousPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(<script|</script|javascript:|on\w+\s*=)`), // XSS
	regexp.MustCompile(`(?i)(;\s*DROP\s+TABLE|;\s*DELETE\s+FROM|';\s*--)`), // SQL injection
}

// ValidateEmail - メールアドレスのバリデーション
func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)

	if email == "" {
		return ErrEmptyEmail
	}

	if utf8.RuneCountInString(email) > MaxEmailLength {
		return ErrEmailTooLong
	}

	if !emailRegex.MatchString(email) {
		return ErrInvalidEmailFormat
	}

	// 危険なパターンチェック
	if containsDangerousPattern(email) {
		return ErrInvalidCharacters
	}

	return nil
}

// ValidatePassword - パスワードのバリデーション
func ValidatePassword(password string) error {
	if password == "" {
		return ErrEmptyPassword
	}

	length := utf8.RuneCountInString(password)

	if length < MinPasswordLength {
		return ErrPasswordTooShort
	}

	if length > MaxPasswordLength {
		return ErrPasswordTooLong
	}

	return nil
}

// ValidateUsername - ユーザー名のバリデーション
func ValidateUsername(username string) error {
	username = strings.TrimSpace(username)

	if username == "" {
		return ErrEmptyUsername
	}

	length := utf8.RuneCountInString(username)

	if length < MinUsernameLength {
		return ErrUsernameTooShort
	}

	if length > MaxUsernameLength {
		return ErrUsernameTooLong
	}

	if !usernameRegex.MatchString(username) {
		return ErrInvalidUsername
	}

	// 危険なパターンチェック
	if containsDangerousPattern(username) {
		return ErrInvalidCharacters
	}

	return nil
}

// ValidatePostContent - 投稿内容のバリデーション
func ValidatePostContent(content string) error {
	content = strings.TrimSpace(content)

	if content == "" {
		return ErrEmptyContent
	}

	if utf8.RuneCountInString(content) > MaxPostContent {
		return errors.New("post content is too long (max 280 characters)")
	}

	// XSSチェック（HTMLタグは警告のみ、エラーにはしない）
	// 実際のアプリケーションではフロントエンドでエスケープする

	return nil
}

// ValidateCommentContent - コメント内容のバリデーション
func ValidateCommentContent(content string) error {
	content = strings.TrimSpace(content)

	if content == "" {
		return ErrEmptyContent
	}

	if utf8.RuneCountInString(content) > MaxCommentContent {
		return errors.New("comment content is too long (max 500 characters)")
	}

	return nil
}

// containsDangerousPattern - 危険なパターンが含まれているかチェック
func containsDangerousPattern(s string) bool {
	for _, pattern := range dangerousPatterns {
		if pattern.MatchString(s) {
			return true
		}
	}
	return false
}
