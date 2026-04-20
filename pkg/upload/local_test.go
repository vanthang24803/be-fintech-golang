package upload

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewLocalUploaderAndFileLifecycle(t *testing.T) {
	t.Parallel()

	basePath := filepath.Join(t.TempDir(), "uploads")
	uploader, err := NewLocalUploader(basePath, "http://localhost/static")
	if err != nil {
		t.Fatalf("NewLocalUploader() error = %v", err)
	}
	if _, err := os.Stat(basePath); err != nil {
		t.Fatalf("expected base path to exist: %v", err)
	}

	info, err := uploader.UploadFile(context.Background(), strings.NewReader("hello world"), int64(len("hello world")), "text/plain", "avatars")
	if err != nil {
		t.Fatalf("UploadFile() error = %v", err)
	}
	if info.Size != int64(len("hello world")) || info.MimeType != "text/plain" {
		t.Fatalf("unexpected file info: %+v", info)
	}
	if !strings.HasPrefix(info.Key, "avatars/") || !strings.HasSuffix(info.Key, ".bin") {
		t.Fatalf("unexpected key format: %q", info.Key)
	}
	if !strings.Contains(info.URL, info.Key) {
		t.Fatalf("expected URL to contain key, got %q", info.URL)
	}
	if _, err := os.Stat(filepath.Join(basePath, info.Key)); err != nil {
		t.Fatalf("expected uploaded file to exist: %v", err)
	}

	url, err := uploader.GetFileURL(context.Background(), info.Key)
	if err != nil || url != "http://localhost/static/"+info.Key {
		t.Fatalf("unexpected GetFileURL() result = %q, err=%v", url, err)
	}

	if err := uploader.DeleteFile(context.Background(), info.Key); err != nil {
		t.Fatalf("DeleteFile() error = %v", err)
	}
	if _, err := os.Stat(filepath.Join(basePath, info.Key)); !os.IsNotExist(err) {
		t.Fatalf("expected file to be removed, stat err=%v", err)
	}
}

type errReader struct{}

func (errReader) Read(p []byte) (int, error) {
	return 0, errors.New("boom")
}

func TestLocalUploader_ErrorBranches(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	uploader, err := NewLocalUploader(filepath.Join(tmpDir, "uploads"), "http://localhost/static")
	if err != nil {
		t.Fatalf("NewLocalUploader(): %v", err)
	}

	blockingFile := filepath.Join(tmpDir, "uploads", "blocker")
	if err := os.WriteFile(blockingFile, []byte("x"), 0644); err != nil {
		t.Fatalf("WriteFile(): %v", err)
	}

	if _, err := uploader.UploadFile(context.Background(), strings.NewReader("ok"), 2, "text/plain", "blocker/child"); err == nil {
		t.Fatal("expected mkdir error when upload folder path is blocked by a file")
	}

	if _, err := uploader.UploadFile(context.Background(), errReader{}, 5, "application/octet-stream", "avatar.png"); err == nil {
		t.Fatal("expected write error from reader")
	}

	if err := uploader.DeleteFile(context.Background(), "missing-file"); err == nil {
		t.Fatal("expected delete error for missing file")
	}
}
