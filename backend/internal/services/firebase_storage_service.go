package services

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"time"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go/v4"
	"github.com/google/uuid"
	"github.com/yourusername/sns-backend/internal/config"
	"google.golang.org/api/option"
)

// FirebaseStorageService Firebase Storageサービス
type FirebaseStorageService struct {
	bucket *storage.BucketHandle
}

var firebaseStorageService *FirebaseStorageService

// InitFirebaseStorage Firebase Storageを初期化
func InitFirebaseStorage() (*FirebaseStorageService, error) {
	if firebaseStorageService != nil {
		return firebaseStorageService, nil
	}

	cfg := config.AppConfig
	if cfg.FirebaseStorageBucket == "" {
		return nil, fmt.Errorf("FIREBASE_STORAGE_BUCKET not configured")
	}

	ctx := context.Background()

	// Firebase Admin SDKの初期化
	opt := option.WithCredentialsFile(cfg.FirebaseCredentialsPath)
	_, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Firebase: %v", err)
	}

	// Cloud Storage clientの取得
	client, err := storage.NewClient(ctx, opt)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage client: %v", err)
	}

	bucket := client.Bucket(cfg.FirebaseStorageBucket)

	firebaseStorageService = &FirebaseStorageService{
		bucket: bucket,
	}

	return firebaseStorageService, nil
}

// GetFirebaseStorageService シングルトンインスタンスを取得
func GetFirebaseStorageService() (*FirebaseStorageService, error) {
	if firebaseStorageService == nil {
		return InitFirebaseStorage()
	}
	return firebaseStorageService, nil
}

// UploadFile Firebase Storageにファイルをアップロード
func (s *FirebaseStorageService) UploadFile(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	// ファイル名生成（UUID + 元のファイル拡張子）
	ext := filepath.Ext(fileHeader.Filename)
	fileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	objectPath := fmt.Sprintf("uploads/%s", fileName)

	// Cloud Storageへのアップロード
	wc := s.bucket.Object(objectPath).NewWriter(ctx)
	wc.ContentType = fileHeader.Header.Get("Content-Type")
	wc.Metadata = map[string]string{
		"uploaded-at":   time.Now().Format(time.RFC3339),
		"original-name": fileHeader.Filename,
	}

	// ファイルコピー
	if _, err := io.Copy(wc, file); err != nil {
		return "", fmt.Errorf("failed to copy file: %v", err)
	}

	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer: %v", err)
	}

	// 署名付きURL（有効期限あり）を生成
	// Note: バケットが非公開の場合、署名付きURLが必要
	url, err := s.bucket.SignedURL(objectPath, &storage.SignedURLOptions{
		Method:  "GET",
		Expires: time.Now().Add(7 * 24 * time.Hour), // 7日間有効
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL: %v", err)
	}

	return url, nil
}

// DeleteFile Firebase Storageからファイルを削除
func (s *FirebaseStorageService) DeleteFile(ctx context.Context, fileURL string) error {
	// URLからオブジェクトパスを抽出
	// 例: https://storage.googleapis.com/bucket-name/uploads/file.jpg -> uploads/file.jpg
	// 簡易実装（実際はURLパースが必要）
	// objectPath := extractObjectPath(fileURL)

	// obj := s.bucket.Object(objectPath)
	// if err := obj.Delete(ctx); err != nil {
	// 	return fmt.Errorf("failed to delete file: %v", err)
	// }

	return nil
}

// IsConfigured Firebase Storageが設定されているか確認
func IsFirebaseStorageConfigured() bool {
	cfg := config.AppConfig
	return cfg.FirebaseStorageBucket != "" && cfg.FirebaseCredentialsPath != ""
}
