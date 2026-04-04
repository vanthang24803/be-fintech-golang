package upload

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/rs/xid"
)

type LocalUploader struct {
	basePath string
	baseURL  string
}

// NewLocalUploader initializes a local storage uploader
func NewLocalUploader(basePath, baseURL string) (*LocalUploader, error) {
	// Create directory if not exists
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		err := os.MkdirAll(basePath, 0755)
		if err != nil {
			return nil, fmt.Errorf("failed to create upload directory: %w", err)
		}
	}

	return &LocalUploader{
		basePath: basePath,
		baseURL:  baseURL,
	}, nil
}

func (l *LocalUploader) UploadFile(ctx context.Context, reader io.Reader, size int64, contentType string, folder string) (*FileInfo, error) {
	// Generate unique filename
	ext := filepath.Ext(folder)
	if ext == "" {
		// Try to guess extension from content type if possible, or use default
		ext = ".bin"
	}
	
	filename := xid.New().String() + ext
	relativeDir := folder
	fullDir := filepath.Join(l.basePath, relativeDir)
	
	// Ensure folder exists
	if err := os.MkdirAll(fullDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create sub-directory: %w", err)
	}

	objectKey := filepath.Join(relativeDir, filename)
	fullPath := filepath.Join(l.basePath, objectKey)

	// Create file
	out, err := os.Create(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Copy content
	written, err := io.Copy(out, reader)
	if err != nil {
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	return &FileInfo{
		Key:      objectKey,
		URL:      fmt.Sprintf("%s/%s", l.baseURL, objectKey),
		Size:     written,
		MimeType: contentType,
	}, nil
}

func (l *LocalUploader) GetFileURL(ctx context.Context, key string) (string, error) {
	return fmt.Sprintf("%s/%s", l.baseURL, key), nil
}

func (l *LocalUploader) DeleteFile(ctx context.Context, key string) error {
	fullPath := filepath.Join(l.basePath, key)
	return os.Remove(fullPath)
}
