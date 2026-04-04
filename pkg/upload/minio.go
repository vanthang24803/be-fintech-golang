package upload

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"path"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/xid"
)

type MinioUploader struct {
	client     *minio.Client
	bucketName string
	endpoint   string
	useSSL     bool
}

// NewMinioUploader initializes a new MinIO client and ensures the bucket exists
func NewMinioUploader(endpoint, accessKey, secretKey, bucketName string, useSSL bool) (*MinioUploader, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize minio client: %w", err)
	}

	// Ensure bucket exists
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if !exists {
		err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}

		// Set public policy for the bucket (Read-Only for anonymous)
		policy := fmt.Sprintf(`{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::%s/*"]}]}`, bucketName)
		err = client.SetBucketPolicy(ctx, bucketName, policy)
		if err != nil {
			return nil, fmt.Errorf("failed to set bucket policy: %w", err)
		}
	}

	return &MinioUploader{
		client:     client,
		bucketName: bucketName,
		endpoint:   endpoint,
		useSSL:     useSSL,
	}, nil
}

// UploadFile uploads an io.Reader to MinIO and returns the FileInfo
func (m *MinioUploader) UploadFile(ctx context.Context, reader io.Reader, size int64, contentType string, folder string) (*FileInfo, error) {
	// Generate unique filename
	filename := xid.New().String() + path.Ext(folder) // Simplified: folder here acts as a prefix or original filename hint
	objectKey := path.Join(folder, filename)

	_, err := m.client.PutObject(ctx, m.bucketName, objectKey, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload object: %w", err)
	}

	return &FileInfo{
		Key:      objectKey,
		URL:      m.buildPublicURL(objectKey),
		Size:     size,
		MimeType: contentType,
	}, nil
}

// GetFileURL returns a pre-signed URL (or public URL depending on bucket policy)
func (m *MinioUploader) GetFileURL(ctx context.Context, key string) (string, error) {
	// For this implementation, we assume the bucket is public-read as set in NewMinioUploader
	return m.buildPublicURL(key), nil
}

// DeleteFile removes an object from MinIO
func (m *MinioUploader) DeleteFile(ctx context.Context, key string) error {
	err := m.client.RemoveObject(ctx, m.bucketName, key, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to remove object: %w", err)
	}
	return nil
}

func (m *MinioUploader) buildPublicURL(key string) string {
	scheme := "http"
	if m.useSSL {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s/%s/%s", scheme, m.endpoint, m.bucketName, key)
}

// GetPresignedURL returns a temporary URL valid for the specified duration
func (m *MinioUploader) GetPresignedURL(ctx context.Context, key string, expires time.Duration) (string, error) {
	reqParams := make(url.Values)
	presignedURL, err := m.client.PresignedGetObject(ctx, m.bucketName, key, expires, reqParams)
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}
