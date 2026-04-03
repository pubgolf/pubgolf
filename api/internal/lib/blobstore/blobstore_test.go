package blobstore

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pubgolf/pubgolf/api/internal/lib/config"
)

// testConfig returns a BlobStoreAuth pointing at a local Minio instance.
// Uses PUBGOLF_BLOB_STORE_ENDPOINT env var if set, otherwise defaults to localhost:9000.
func testConfig() config.BlobStoreAuth {
	endpoint := os.Getenv("PUBGOLF_BLOB_STORE_ENDPOINT")
	if endpoint == "" {
		endpoint = "localhost:9000"
	}

	return config.BlobStoreAuth{ //nolint:gosec // Dev-only test credentials.
		Endpoint:  endpoint,
		AccessKey: "pubgolf_dev",
		SecretKey: "pubgolf_dev",
		Bucket:    "pubgolf-test",
		UseSSL:    false,
	}
}

// requireMinio skips the test if a live Minio instance is not reachable.
func requireMinio(t *testing.T, client *Client) {
	t.Helper()

	_, err := client.BucketExists(context.Background())
	if err != nil {
		t.Skipf("Minio not available: %v", err)
	}
}

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("creates client with valid config", func(t *testing.T) {
		t.Parallel()

		client, err := New(config.BlobStoreAuth{
			Endpoint:  "localhost:9000",
			AccessKey: "test-key",
			SecretKey: "test-secret",
			Bucket:    "test-bucket",
			UseSSL:    false,
		})

		require.NoError(t, err)
		assert.NotNil(t, client)
		assert.Equal(t, "test-bucket", client.bucket)
	})
}

func TestPresignedPutURL(t *testing.T) {
	t.Parallel()

	client, err := New(testConfig())
	require.NoError(t, err)
	requireMinio(t, client)

	t.Run("returns URL containing bucket and object key", func(t *testing.T) {
		t.Parallel()

		url, err := client.PresignedPutURL(context.Background(), "scores/abc/123.jpg", 15*time.Minute)

		require.NoError(t, err)
		assert.Contains(t, url, testConfig().Bucket)
		assert.Contains(t, url, "scores/abc/123.jpg")
	})
}

func TestPresignedGetURL(t *testing.T) {
	t.Parallel()

	client, err := New(testConfig())
	require.NoError(t, err)
	requireMinio(t, client)

	t.Run("returns URL containing bucket and object key", func(t *testing.T) {
		t.Parallel()

		url, err := client.PresignedGetURL(context.Background(), "venues/xyz/456.jpg", 15*time.Minute)

		require.NoError(t, err)
		assert.Contains(t, url, testConfig().Bucket)
		assert.Contains(t, url, "venues/xyz/456.jpg")
	})
}
