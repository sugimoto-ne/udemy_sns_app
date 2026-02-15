package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
	Port       string
	Env        string
	IsTestMode bool
}

var AppConfig *Config

func LoadConfig() *Config {
	// .envファイルを読み込み
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	isTestMode := getEnv("TEST_MODE", "false") == "true"

	// テストモードの場合は、テスト用のDB設定を使用
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "sns_db")

	if isTestMode {
		dbHost = getEnv("DB_TEST_HOST", "localhost")
		dbPort = getEnv("DB_TEST_PORT", "5433")
		dbUser = getEnv("DB_TEST_USER", "postgres")
		dbPassword = getEnv("DB_TEST_PASSWORD", "postgres")
		dbName = getEnv("DB_TEST_NAME", "sns_db_test")
		log.Println("⚠️  Running in TEST MODE - using test database:", dbName)
	}

	config := &Config{
		DBHost:     dbHost,
		DBPort:     dbPort,
		DBUser:     dbUser,
		DBPassword: dbPassword,
		DBName:     dbName,
		JWTSecret:  getEnv("JWT_SECRET", "secret"),
		Port:       getEnv("PORT", "8080"),
		Env:        getEnv("ENV", "development"),
		IsTestMode: isTestMode,
	}

	AppConfig = config
	return config
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost,
		c.DBPort,
		c.DBUser,
		c.DBPassword,
		c.DBName,
	)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
