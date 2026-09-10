package s3

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Config struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	PublicEndpoint  string
}

func DefaultConfig() Config {
	return Config{
		Endpoint:        "127.0.0.1:9000",
		AccessKeyID:     "minioadmin",
		SecretAccessKey: "minioadmin",
		UseSSL:          false,
		PublicEndpoint:  "127.0.0.1:9000",
	}
}

type MinIOClient struct {
	client *minio.Client
	cfg    Config
}

func NewMinIOClient(cfg Config) (*MinIOClient, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize minio client: %w", err)
	}

	return &MinIOClient{
		client: client,
		cfg:    cfg,
	}, nil
}

func (m *MinIOClient) EnsureBuckets(ctx context.Context, bucketNames ...string) error {
	for _, bucket := range bucketNames {
		exists, err := m.client.BucketExists(ctx, bucket)
		if err != nil {
			return fmt.Errorf("failed to check bucket existence '%s': %w", bucket, err)
		}
		if !exists {
			err = m.client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
			if err != nil {
				return fmt.Errorf("failed to create bucket '%s': %w", bucket, err)
			}
			log.Printf("✅ [MinIO S3] Created bucket: %s\n", bucket)
		}
	}
	return nil
}

func (m *MinIOClient) GeneratePresignedGetURL(ctx context.Context, bucket, objectKey string, expires time.Duration) (string, error) {
	reqParams := make(url.Values)
	presignedURL, err := m.client.PresignedGetObject(ctx, bucket, objectKey, expires, reqParams)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned GET url for %s/%s: %w", bucket, objectKey, err)
	}
	return presignedURL.String(), nil
}

func (m *MinIOClient) GeneratePresignedPutURL(ctx context.Context, bucket, objectKey string, expires time.Duration) (string, error) {
	presignedURL, err := m.client.PresignedPutObject(ctx, bucket, objectKey, expires)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned PUT url for %s/%s: %w", bucket, objectKey, err)
	}
	return presignedURL.String(), nil
}

func (m *MinIOClient) PutObject(ctx context.Context, bucket, objectKey string, reader io.Reader, size int64, contentType string) error {
	_, err := m.client.PutObject(ctx, bucket, objectKey, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("failed to put object %s/%s: %w", bucket, objectKey, err)
	}
	return nil
}

func (m *MinIOClient) DeleteObject(ctx context.Context, bucket, objectKey string) error {
	err := m.client.RemoveObject(ctx, bucket, objectKey, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete object %s/%s: %w", bucket, objectKey, err)
	}
	return nil
}

func (m *MinIOClient) IsHealthy(ctx context.Context) bool {
	ctxTimeout, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	_, err := m.client.ListBuckets(ctxTimeout)
	return err == nil
}
