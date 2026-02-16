package testutil

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/yourusername/sns-backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SetupTestDB sets up a test database connection
func SetupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	// テスト用のDB接続情報
	host := getEnv("DB_TEST_HOST", "db_test")
	port := getEnv("DB_TEST_PORT", "5432")
	user := getEnv("DB_TEST_USER", "postgres")
	password := getEnv("DB_TEST_PASSWORD", "postgres")
	dbname := getEnv("DB_TEST_NAME", "sns_db_test")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	// テスト用DBに接続
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// マイグレーション実行
	err = db.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Media{},
		&models.Comment{},
		&models.PostLike{},
		&models.Follow{},
		&models.Hashtag{},
		&models.PostHashtag{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	log.Println("✅ Test database setup completed")

	return db
}

// CleanupTestDB cleans up all data from test database
func CleanupTestDB(t *testing.T, db *gorm.DB) {
	t.Helper()

	// テーブルの順序に注意（外部キー制約のため）
	tables := []interface{}{
		&models.PostLike{},
		&models.Comment{},
		&models.Media{},
		&models.PostHashtag{},
		&models.Hashtag{},
		&models.Post{},
		&models.Follow{},
		&models.User{},
	}

	for _, table := range tables {
		if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(table).Error; err != nil {
			t.Logf("Warning: Failed to cleanup table: %v", err)
		}
	}

	log.Println("🧹 Test database cleaned up")
}

// TeardownTestDB closes the database connection
func TeardownTestDB(t *testing.T, db *gorm.DB) {
	t.Helper()

	sqlDB, err := db.DB()
	if err != nil {
		t.Logf("Warning: Failed to get database instance: %v", err)
		return
	}

	if err := sqlDB.Close(); err != nil {
		t.Logf("Warning: Failed to close database connection: %v", err)
	}

	log.Println("👋 Test database connection closed")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
