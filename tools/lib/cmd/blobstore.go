package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// errBlobStoreNotReady is returned when Minio does not become available after retries.
var errBlobStoreNotReady = errors.New("blob storage not ready after retries")

const (
	blobBucketPrefix = "pubgolf-dev"
	minioDevEndpoint = "localhost:9000"
)

// minioClient constructs a minio-go client pointing at the shared Minio instance
// with credentials from the env provider.
func minioClient(ctx context.Context, ep EnvProvider) (*minio.Client, error) {
	vars := readEnvVars(ctx, ep, config.ServerBinName, config.DopplerEnvName, config.EnvVarPrefix, []string{
		"BLOB_STORE_ACCESS_KEY",
		"BLOB_STORE_SECRET_KEY",
	})

	accessKey := getStr(vars, "BLOB_STORE_ACCESS_KEY", "pubgolf_dev")
	secretKey := getStr(vars, "BLOB_STORE_SECRET_KEY", "pubgolf_dev")

	mc, err := minio.New(minioDevEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, fmt.Errorf("create minio client: %w", err)
	}

	return mc, nil
}

// ensureBucket creates the per-worktree bucket if it does not already exist.
// Retries up to 10 times with 500ms backoff to handle the Minio startup delay
// after docker-compose brings the container up.
func ensureBucket(ctx context.Context, ep EnvProvider, slug string) error {
	mc, err := minioClient(ctx, ep)
	if err != nil {
		return err
	}

	bucket := blobBucketForSlug(slug)

	for attempt := range 10 {
		exists, checkErr := mc.BucketExists(ctx, bucket)
		if checkErr == nil {
			if exists {
				return nil
			}

			mkErr := mc.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
			if mkErr != nil {
				return fmt.Errorf("create bucket %q: %w", bucket, mkErr)
			}

			log.Printf("Created blob storage bucket %q", bucket)

			return nil
		}

		if attempt < 9 {
			time.Sleep(500 * time.Millisecond)
		}
	}

	return fmt.Errorf("%w for bucket %q", errBlobStoreNotReady, bucket)
}

// deleteBucket removes all objects in a bucket and then deletes the bucket.
// Returns nil if the bucket does not exist.
func deleteBucket(ctx context.Context, ep EnvProvider, bucket string) error {
	mc, err := minioClient(ctx, ep)
	if err != nil {
		return err
	}

	exists, err := mc.BucketExists(ctx, bucket)
	if err != nil {
		return fmt.Errorf("check bucket %q: %w", bucket, err)
	}

	if !exists {
		return nil
	}

	// Remove all objects before deleting the bucket (batch API).
	objectsCh := make(chan minio.ObjectInfo)

	go func() {
		defer close(objectsCh)

		for obj := range mc.ListObjects(ctx, bucket, minio.ListObjectsOptions{Recursive: true}) {
			if obj.Err != nil {
				log.Printf("Warning: error listing objects in %q: %v", bucket, obj.Err)

				return
			}

			objectsCh <- obj
		}
	}()

	for err := range mc.RemoveObjects(ctx, bucket, objectsCh, minio.RemoveObjectsOptions{}) {
		return fmt.Errorf("remove objects from %q: %w", bucket, err.Err)
	}

	err = mc.RemoveBucket(ctx, bucket)
	if err != nil {
		return fmt.Errorf("remove bucket %q: %w", bucket, err)
	}

	log.Printf("Deleted blob storage bucket %q", bucket)

	return nil
}

// listDevBuckets returns all buckets matching the pubgolf-dev* naming convention.
func listDevBuckets(ctx context.Context, ep EnvProvider) ([]string, error) {
	mc, err := minioClient(ctx, ep)
	if err != nil {
		return nil, err
	}

	buckets, err := mc.ListBuckets(ctx)
	if err != nil {
		return nil, fmt.Errorf("list buckets: %w", err)
	}

	var devBuckets []string

	for _, b := range buckets {
		if strings.HasPrefix(b.Name, blobBucketPrefix) {
			devBuckets = append(devBuckets, b.Name)
		}
	}

	return devBuckets, nil
}

// slugFromBucket extracts the worktree slug from a bucket name.
// Returns "" for the main bucket ("pubgolf-dev"), and the slug portion
// for worktree buckets ("pubgolf-dev-foo" -> "foo").
func slugFromBucket(bucket string) string {
	after, found := strings.CutPrefix(bucket, blobBucketPrefix+"-")
	if found {
		return after
	}

	return ""
}
