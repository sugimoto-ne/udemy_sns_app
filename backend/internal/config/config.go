package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost                    string
	DBPort                    string
	DBUser                    string
	DBPassword                string
	DBName                    string
	JWTSecret                 string
	Port                      string
	Env                       string
	IsTestMode                bool
	FirebaseCredentialsPath   string
	FirebaseStorageBucket     string
	ResendAPIKey              string
	FrontendURL               string
	FromEmail                 string
}

var AppConfig *Config

func LoadConfig() *Config {
	// .envファイルを読み込み
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	env := getEnv("ENV", "development")
	isTestMode := getEnv("TEST_MODE", "false") == "true"

	// JWT_SECRETのバリデーション（必須）
	jwtSecret := getEnv("JWT_SECRET", "")
	if jwtSecret == "" {
		log.Fatal("❌ JWT_SECRET environment variable is required")
	}
	// 本番環境では強力なシークレットを強制
	if env == "production" && len(jwtSecret) < 32 {
		log.Fatal("❌ JWT_SECRET must be at least 32 characters in production")
	}

	// テストモードの場合は、テスト用のDB設定を使用
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "")
	dbName := getEnv("DB_NAME", "sns_db")

	// 本番環境ではDB_PASSWORDを必須化
	if env == "production" && dbPassword == "" {
		log.Fatal("❌ DB_PASSWORD environment variable is required in production")
	}
	// 開発環境のみデフォルト値を許可
	if dbPassword == "" {
		dbPassword = "postgres"
	}

	if isTestMode {
		dbHost = getEnv("DB_TEST_HOST", "localhost")
		dbPort = getEnv("DB_TEST_PORT", "5433")
		dbUser = getEnv("DB_TEST_USER", "postgres")
		dbPassword = getEnv("DB_TEST_PASSWORD", "postgres")
		dbName = getEnv("DB_TEST_NAME", "sns_db_test")
		log.Println("⚠️  Running in TEST MODE - using test database:", dbName)
	}

	config := &Config{
		DBHost:                    dbHost,
		DBPort:                    dbPort,
		DBUser:                    dbUser,
		DBPassword:                dbPassword,
		DBName:                    dbName,
		JWTSecret:                 jwtSecret,
		Port:                      getEnv("PORT", "8080"),
		Env:                       env,
		IsTestMode:                isTestMode,
		FirebaseCredentialsPath:   getEnv("FIREBASE_CREDENTIALS_PATH", "./service_account_key.json"),
		FirebaseStorageBucket:     getEnv("FIREBASE_STORAGE_BUCKET", ""),
		ResendAPIKey:              getEnv("RESEND_API_KEY", ""),
		FrontendURL:               getEnv("FRONTEND_URL", "http://localhost:5173"),
		FromEmail:                 getEnv("FROM_EMAIL", "noreply@example.com"),
	}

	AppConfig = config
	return config
}

func (c *Config) GetDSN() string {
	// 本番環境ではSSLを有効化（Neon等のクラウドDBで必須）
	sslMode := "disable"
	if c.Env == "production" {
		sslMode = "require"
	}

	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost,
		c.DBPort,
		c.DBUser,
		c.DBPassword,
		c.DBName,
		sslMode,
	)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
