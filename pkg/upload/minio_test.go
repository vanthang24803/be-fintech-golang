package upload

import (
	"context"
	"testing"
	"time"
)

func TestMinioUploaderHelpers(t *testing.T) {
	t.Parallel()

	uploader := &MinioUploader{
		bucketName: "bucket",
		endpoint:   "minio.local",
		useSSL:     true,
	}

	if got := uploader.buildPublicURL("dir/file.txt"); got != "https://minio.local/bucket/dir/file.txt" {
		t.Fatalf("unexpected public url: %q", got)
	}
	uploader.useSSL = false
	if got := uploader.buildPublicURL("dir/file.txt"); got != "http://minio.local/bucket/dir/file.txt" {
		t.Fatalf("unexpected public url without ssl: %q", got)
	}
	url, err := uploader.GetFileURL(context.Background(), "dir/file.txt")
	if err != nil || url != "http://minio.local/bucket/dir/file.txt" {
		t.Fatalf("unexpected GetFileURL result: %q err=%v", url, err)
	}
}

func TestNewMinioUploaderAndPresignedURLErrors(t *testing.T) {
	t.Parallel()

	if _, err := NewMinioUploader("", "a", "b", "bucket", false); err == nil {
		t.Fatal("expected invalid endpoint error")
	}

	uploader := &MinioUploader{}
	if _, err := uploader.GetPresignedURL(context.Background(), "key", time.Minute); err == nil {
		t.Fatal("expected nil client presign error")
	}
}

func TestMinioDeleteFileNilClientReturnsError(t *testing.T) {
	t.Parallel()

	if err := (&MinioUploader{}).DeleteFile(context.Background(), "key"); err == nil {
		t.Fatal("expected delete error when client is nil")
	}
}
