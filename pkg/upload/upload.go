package upload

import (
	"context"
	"io"
)

// FileInfo contains basic information about the uploaded file
type FileInfo struct {
	Key      string `json:"key"`
	URL      string `json:"url"`
	Size     int64  `json:"size"`
	MimeType string `json:"mime_type"`
}

// Uploader defines the interface for file storage operations
type Uploader interface {
	UploadFile(ctx context.Context, reader io.Reader, size int64, contentType string, folder string) (*FileInfo, error)
	GetFileURL(ctx context.Context, key string) (string, error)
	DeleteFile(ctx context.Context, key string) error
}
