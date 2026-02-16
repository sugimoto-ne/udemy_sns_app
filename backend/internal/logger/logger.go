package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Logger zerolog.Logger

// InitLogger - ロガーを初期化
func InitLogger() {
	env := os.Getenv("APP_ENV")

	// 環境に応じてログ出力形式を切り替え
	if env == "development" || env == "test" {
		// 開発/テスト環境: 人間が読みやすい形式
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		})
	} else {
		// 本番環境: JSON形式
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}

	// ログレベルを設定
	logLevel := os.Getenv("LOG_LEVEL")
	switch logLevel {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		// デフォルトはINFO
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	Logger = log.Logger
}

// GetLogger - グローバルロガーを取得
func GetLogger() zerolog.Logger {
	return Logger
}
