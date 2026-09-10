package ports

import (
	"context"
	"io"
	"time"
)

type ObjectStorageService interface {
	EnsureBuckets(ctx context.Context, bucketNames ...string) error
	GeneratePresignedGetURL(ctx context.Context, bucket, objectKey string, expires time.Duration) (string, error)
	GeneratePresignedPutURL(ctx context.Context, bucket, objectKey string, expires time.Duration) (string, error)
	PutObject(ctx context.Context, bucket, objectKey string, reader io.Reader, size int64, contentType string) error
	DeleteObject(ctx context.Context, bucket, objectKey string) error
	IsHealthy(ctx context.Context) bool
}
