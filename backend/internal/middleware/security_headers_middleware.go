package middleware

import (
	"os"

	"github.com/labstack/echo/v4"
)

// SecurityHeaders - セキュリティヘッダーを設定するミドルウェア
// 本番環境（APP_ENV=production）でのみセキュリティヘッダーを付与
// 開発環境ではCSPがViteのWebSocketをブロックし、HSTSがlocalhostを拒否するため無効化
func SecurityHeaders() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 本番環境でのみセキュリティヘッダーを付与
			env := os.Getenv("APP_ENV")
			if env != "production" {
				// 開発環境ではヘッダーを付与せずに次のハンドラーへ
				return next(c)
			}

			// Content-Security-Policy (CSP)
			// default-src 'self': 自サイトのリソースのみ許可
			// script-src 'self': 自サイトのスクリプトのみ許可
			// style-src 'self' 'unsafe-inline': 自サイトのスタイル + インラインスタイル許可（MUI対応）
			// img-src 'self' data: https:: 自サイトの画像 + data URI + HTTPS画像許可
			// font-src 'self': 自サイトのフォントのみ許可
			// connect-src 'self': 自サイトへの接続のみ許可
			c.Response().Header().Set("Content-Security-Policy",
				"default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self'; connect-src 'self'")

			// X-Frame-Options: DENY
			// クリックジャッキング攻撃を防ぐため、iframe内での表示を禁止
			c.Response().Header().Set("X-Frame-Options", "DENY")

			// X-Content-Type-Options: nosniff
			// ブラウザによるMIMEタイプの推測を防ぎ、宣言されたContent-Typeを強制
			c.Response().Header().Set("X-Content-Type-Options", "nosniff")

			// Strict-Transport-Security (HSTS)
			// HTTPS接続を強制（本番環境でのみ有効）
			// max-age=31536000: 1年間HSTSを有効化
			// includeSubDomains: サブドメインにも適用
			c.Response().Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

			// X-XSS-Protection: 0
			// CSPがあるため、レガシーなXSS Protection機能は無効化
			// 最新のブラウザではCSPを使用することが推奨される
			c.Response().Header().Set("X-XSS-Protection", "0")

			return next(c)
		}
	}
}
