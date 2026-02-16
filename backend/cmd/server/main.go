package main

import (
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/yourusername/sns-backend/internal/config"
	"github.com/yourusername/sns-backend/internal/database"
	"github.com/yourusername/sns-backend/internal/logger"
	customMiddleware "github.com/yourusername/sns-backend/internal/middleware"
	"github.com/yourusername/sns-backend/internal/models"
	"github.com/yourusername/sns-backend/internal/routes"

	_ "github.com/yourusername/sns-backend/docs" // Swagger生成ファイルをインポート
)

// @title SNS API
// @version 1.0
// @description TwitterライクなSNSアプリケーションのREST API
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT認証トークン。形式: "Bearer {token}"

func main() {
	// ロガーを初期化
	logger.InitLogger()
	log := logger.GetLogger()

	// 設定を読み込み
	cfg := config.LoadConfig()
	log.Info().Msg("Configuration loaded")

	// データベース接続
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	// AutoMigrate - テーブルを自動作成
	log.Info().Msg("Running database migrations...")
	if err := db.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Media{},
		&models.Comment{},
		&models.PostLike{},
		&models.Follow{},
		&models.RefreshToken{},
		// Phase 2
		&models.Hashtag{},
		&models.PostHashtag{},
		&models.Bookmark{},
		&models.PasswordResetToken{},
		&models.EmailVerificationToken{},
	); err != nil {
		log.Fatal().Err(err).Msg("Failed to migrate database")
	}
	log.Info().Msg("Database migrations completed")

	// Echoインスタンスを作成
	e := echo.New()

	// カスタムエラーハンドラー設定
	e.HTTPErrorHandler = customMiddleware.ErrorHandler

	// ミドルウェア（順序重要）
	e.Use(customMiddleware.RequestID())      // リクエストID生成（最初）
	e.Use(customMiddleware.AccessLog())      // アクセスログ
	e.Use(echoMiddleware.Recover())          // パニック回復
	e.Use(customMiddleware.CORS())           // CORS
	e.Use(customMiddleware.SecurityHeaders()) // セキュリティヘッダー

	// レート制限（テスト・開発環境では緩和）
	authLimit := 5   // 認証系: 5回/分（本番）
	generalLimit := 60 // 一般: 60回/分（本番）
	if cfg.Env == "test" || cfg.Env == "development" {
		authLimit = 1000    // テスト・開発環境: 1000回/分
		generalLimit = 1000 // テスト・開発環境: 1000回/分
		log.Warn().Msg("Rate limit relaxed for development/test environment")
	}
	e.Use(customMiddleware.RateLimit(authLimit, generalLimit))

	// ヘルスチェックエンドポイント
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status": "ok",
			"env":    cfg.Env,
		})
	})

	// ルート設定
	routes.SetupRoutes(e)

	// Swagger UIエンドポイント
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	log.Info().Str("url", "http://localhost:"+cfg.Port+"/swagger/index.html").Msg("Swagger UI available")

	// サーバー起動
	log.Info().Str("port", cfg.Port).Msg("Server starting")
	if err := e.Start(":" + cfg.Port); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
