// Package blobstore provides an S3-compatible blob storage client for storing and retrieving objects.
// It wraps minio-go and works against Minio (local dev), AWS S3, and Cloudflare R2 with no code changes.
//
// Object key conventions (enforced by callers):
//   - Score photos: scores/{score_id}/{ulid}.jpg
//   - Venue images: venues/{venue_id}/{ulid}.jpg
package blobstore

import (
	"context"
	"fmt"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/pubgolf/pubgolf/api/internal/lib/config"
	"github.com/pubgolf/pubgolf/api/internal/lib/telemetry"
)

// Client wraps a minio-go client with the configured bucket.
type Client struct {
	mc     *minio.Client
	bucket string
}

// New creates a new blob storage client from the given config.
func New(cfg config.BlobStoreAuth) (*Client, error) {
	mc, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("create minio client: %w", err)
	}

	return &Client{
		mc:     mc,
		bucket: cfg.Bucket,
	}, nil
}

// BucketExists checks whether the configured bucket is reachable.
func (c *Client) BucketExists(ctx context.Context) (bool, error) {
	defer telemetry.FnSpan(&ctx)()

	exists, err := c.mc.BucketExists(ctx, c.bucket)
	if err != nil {
		return false, fmt.Errorf("check bucket %q: %w", c.bucket, err)
	}

	return exists, nil
}

// PresignedPutURL returns a short-lived pre-signed URL that allows the caller to upload
// an object directly to blob storage without server-side proxying.
func (c *Client) PresignedPutURL(ctx context.Context, objectKey string, expiry time.Duration) (string, error) {
	defer telemetry.FnSpan(&ctx)()

	u, err := c.mc.PresignedPutObject(ctx, c.bucket, objectKey, expiry)
	if err != nil {
		return "", fmt.Errorf("generate presigned PUT URL for %q: %w", objectKey, err)
	}

	return u.String(), nil
}

// PresignedGetURL returns a short-lived pre-signed URL that allows the caller to download
// or display an object directly from blob storage without server-side proxying.
func (c *Client) PresignedGetURL(ctx context.Context, objectKey string, expiry time.Duration) (string, error) {
	defer telemetry.FnSpan(&ctx)()

	u, err := c.mc.PresignedGetObject(ctx, c.bucket, objectKey, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("generate presigned GET URL for %q: %w", objectKey, err)
	}

	return u.String(), nil
}

// DeleteObject removes an object from blob storage.
func (c *Client) DeleteObject(ctx context.Context, objectKey string) error {
	defer telemetry.FnSpan(&ctx)()

	err := c.mc.RemoveObject(ctx, c.bucket, objectKey, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("delete object %q: %w", objectKey, err)
	}

	return nil
}
