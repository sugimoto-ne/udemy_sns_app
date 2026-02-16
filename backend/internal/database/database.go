package database

import (
	"log"
	"os"
	"time"

	"github.com/yourusername/sns-backend/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.GetDSN()

	// GORM logger設定（開発環境では詳細ログ）
	var gormLogger logger.Interface
	if cfg.Env == "development" {
		// カスタムロガー設定（クエリ数カウント用）
		gormLogger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             200 * time.Millisecond, // 200ms以上のクエリを警告
				LogLevel:                  logger.Info,            // 全てのSQLを表示
				IgnoreRecordNotFoundError: true,
				Colorful:                  true,
			},
		)
	} else {
		// 本番環境では警告以上のみ
		gormLogger = logger.Default.LogMode(logger.Warn)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})

	if err != nil {
		return nil, err
	}

	// 接続プール設定
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	DB = db
	log.Println("✅ Database connection established")

	return db, nil
}

func GetDB() *gorm.DB {
	return DB
}
