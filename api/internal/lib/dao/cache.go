package dao

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var (
	// cacheSizeEvent is the default item size for caches that hold event-scoped data.
	cacheSizeEvent cacheSize = 8
	// cacheSizeVenue is the default item size for caches that hold venue-scoped data.
	cacheSizeVenue cacheSize = 16
	// cacheSizePlayer is the default item size for caches that hold player-scoped data.
	cacheSizePlayer cacheSize = 128
)

var (
	// cacheExpirationDurable is the default expiration duration for caches that hold long-lasting data.
	cacheExpirationDurable = cacheExpiration(24 * time.Hour)
	// cacheExpirationModerate is the default expiration duration for caches that hold medium-lasting data.
	cacheExpirationModerate = cacheExpiration(10 * time.Minute)
	// cacheExpirationVolatile is the default expiration duration for caches that hold short-lived data.
	cacheExpirationVolatile = cacheExpiration(1 * time.Minute)
)

// errDoNotCacheResult is returned from a cache query function to prevent caching the result.
var errDoNotCacheResult = errors.New("uncachable result")

// Purger is an interface for in-memory caches that can be purged.
type Purger interface {
	Purge()
}

// allCaches is a list of all in-memory caches to allow purging.
var allCaches []Purger

// PurgeAllCaches purges all in-memory caches.
func PurgeAllCaches() {
	for _, c := range allCaches {
		c.Purge()
	}
}

type (
	cache[K comparable, V any] *expirable.LRU[K, V]
	cacheSize                  int
	cacheExpiration            time.Duration
)

func emptyEvictionCallback[K comparable, V any](_ K, _ V) {}

// makeCache creates a new in-memory cache with the provided size and expiration, registering it to allow purging.
func makeCache[K comparable, V any](size cacheSize, exp cacheExpiration) cache[K, V] {
	c := expirable.NewLRU(int(size), emptyEvictionCallback[K, V], time.Duration(exp))
	allCaches = append(allCaches, c)

	return c
}

// wrapWithCache handles access and instrumentation of the provided cache, falling back to access via the provided query function.
func wrapWithCache[K comparable, V any](ctx context.Context, cache *expirable.LRU[K, V], query func(context.Context, K) (V, error), params K) (V, //nolint:ireturn // Needs to be *inside* the function signature to disable ireturn but not trigger nolintlint.
	error,
) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.Bool("dao.cache.available", true))

	if val, ok := cache.Get(params); ok {
		span.SetAttributes(attribute.Bool("dao.cache.hit", true))

		return val, nil
	}

	span.SetAttributes(attribute.Bool("dao.cache.hit", false))

	val, err := query(ctx, params)
	if err != nil {
		if errors.Is(err, errDoNotCacheResult) {
			span.SetAttributes(attribute.Bool("dao.cache.placed", false))

			return val, nil
		}

		return val, fmt.Errorf("fetch from database on cache miss: %w", err)
	}

	evicted := cache.Add(params, val)
	span.SetAttributes(attribute.Bool("dao.cache.placed", true), attribute.Bool("dao.cache.evicted", evicted))

	return val, nil
}
