package utils

import (
	"html"
	"strings"
)

// SanitizeText - テキストをサニタイズ（HTMLエスケープ）
// XSS攻撃を防ぐため、HTMLタグをエスケープする
func SanitizeText(text string) string {
	// 前後の空白を除去
	text = strings.TrimSpace(text)

	// HTMLエスケープ
	return html.EscapeString(text)
}

// SanitizeMultiline - 複数行テキストをサニタイズ
func SanitizeMultiline(text string) string {
	// 改行を保持しながらエスケープ
	text = strings.TrimSpace(text)
	return html.EscapeString(text)
}
