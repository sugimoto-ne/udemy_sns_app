package services

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/yourusername/sns-backend/internal/config"
	"github.com/yourusername/sns-backend/internal/database"
	"github.com/yourusername/sns-backend/internal/models"
	"gorm.io/gorm"
)

// MediaService メディアサービス
type MediaService struct {
	db                     *gorm.DB
	firebaseStorageService *FirebaseStorageService
	useFirebase            bool
}

// NewMediaService MediaServiceのコンストラクタ
func NewMediaService() *MediaService {
	cfg := config.AppConfig

	// Firebase Storageが設定されているか確認
	useFirebase := cfg.FirebaseCredentialsPath != "" && cfg.FirebaseStorageBucket != ""

	var firebaseService *FirebaseStorageService
	if useFirebase {
		var err error
		firebaseService, err = InitFirebaseStorage()
		if err != nil {
			// Firebase初期化失敗時はローカルモードにフォールバック
			fmt.Printf("Warning: Failed to initialize Firebase Storage: %v\n", err)
			useFirebase = false
		}
	}

	return &MediaService{
		db:                     database.GetDB(),
		firebaseStorageService: firebaseService,
		useFirebase:            useFirebase,
	}
}

// UploadMedia メディアをアップロード
func (s *MediaService) UploadMedia(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader, postID uint, orderIndex int) (*models.Media, error) {
	// ファイルサイズ取得
	fileSize := fileHeader.Size

	// メディアタイプを判定
	mediaType, err := s.getMediaType(fileHeader.Filename)
	if err != nil {
		return nil, err
	}

	// メディアタイプごとのサイズ制限チェック
	if err := s.validateFileSize(mediaType, fileSize); err != nil {
		return nil, err
	}

	// アップロード処理
	var mediaURL string
	if s.useFirebase {
		// Firebase Storageにアップロード
		mediaURL, err = s.firebaseStorageService.UploadFile(ctx, file, fileHeader)
		if err != nil {
			return nil, fmt.Errorf("failed to upload to Firebase Storage: %w", err)
		}
	} else {
		// ローカルストレージにアップロード（Phase 1の動作）
		return nil, errors.New("local storage upload not implemented in Phase 2")
	}

	// メディアレコード作成
	media := &models.Media{
		PostID:     postID,
		MediaType:  mediaType,
		MediaURL:   mediaURL,
		FileSize:   fileSize,
		OrderIndex: orderIndex,
	}

	if err := s.db.WithContext(ctx).Create(media).Error; err != nil {
		return nil, err
	}

	return media, nil
}

// UploadMultipleMedia 複数メディアをアップロード（最大4枚）
func (s *MediaService) UploadMultipleMedia(ctx context.Context, files []*multipart.FileHeader, postID uint) ([]models.Media, error) {
	if len(files) > 4 {
		return nil, errors.New("maximum 4 media files allowed")
	}

	if len(files) == 0 {
		return []models.Media{}, nil
	}

	mediaList := make([]models.Media, 0, len(files))

	for i, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open file %d: %w", i, err)
		}
		defer file.Close()

		media, err := s.UploadMedia(ctx, file, fileHeader, postID, i)
		if err != nil {
			return nil, fmt.Errorf("failed to upload media %d: %w", i, err)
		}

		mediaList = append(mediaList, *media)
	}

	return mediaList, nil
}

// getMediaType ファイル名から媒体タイプを判定
func (s *MediaService) getMediaType(filename string) (string, error) {
	ext := strings.ToLower(filepath.Ext(filename))

	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".heic":
		return "image", nil
	case ".mp4", ".mov":
		return "video", nil
	case ".mp3":
		return "audio", nil
	default:
		return "", fmt.Errorf("unsupported file format: %s", ext)
	}
}

// validateFileSize ファイルサイズを検証
func (s *MediaService) validateFileSize(mediaType string, fileSize int64) error {
	const (
		MaxImageSize = 5 * 1024 * 1024   // 5 MB
		MaxVideoSize = 50 * 1024 * 1024  // 50 MB
		MaxAudioSize = 10 * 1024 * 1024  // 10 MB
	)

	switch mediaType {
	case "image":
		if fileSize > MaxImageSize {
			return fmt.Errorf("image file size exceeds limit (max 5MB)")
		}
	case "video":
		if fileSize > MaxVideoSize {
			return fmt.Errorf("video file size exceeds limit (max 50MB)")
		}
	case "audio":
		if fileSize > MaxAudioSize {
			return fmt.Errorf("audio file size exceeds limit (max 10MB)")
		}
	}

	return nil
}

// DeleteMedia メディアを削除
func (s *MediaService) DeleteMedia(ctx context.Context, mediaID uint) error {
	// メディア情報取得
	var media models.Media
	if err := s.db.WithContext(ctx).First(&media, mediaID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("media not found")
		}
		return err
	}

	// DBから削除
	if err := s.db.WithContext(ctx).Delete(&media).Error; err != nil {
		return err
	}

	// Firebase Storageからの削除は実装しない（将来的に実装可能）
	// 理由: URLからファイルパスを抽出して削除する必要があるが、複雑性を避けるため省略

	return nil
}

// GetMediaByPostID 投稿に紐づくメディア一覧を取得
func (s *MediaService) GetMediaByPostID(ctx context.Context, postID uint) ([]models.Media, error) {
	var mediaList []models.Media
	if err := s.db.WithContext(ctx).
		Where("post_id = ?", postID).
		Order("order_index ASC").
		Find(&mediaList).Error; err != nil {
		return nil, err
	}

	return mediaList, nil
}
