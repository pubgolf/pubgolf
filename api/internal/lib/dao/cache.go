package dao

import (
	"context"
	"fmt"

	"github.com/hashicorp/golang-lru/v2/expirable"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var (
	// defaultEventCacheSize is the default item size for caches that hold event-scoped data.
	defaultEventCacheSize = 8
	// defaultVenueCacheSize is the default item size for caches that hold venue-scoped data.
	defaultVenueCacheSize = 16
	// defaultPlayerCacheSize is the default item size for caches that hold player-scoped data.
	defaultPlayerCacheSize = 128
)

// wrapWithCache handles access and instrumentation of the provided cache, falling back to access via the provided query function.
func wrapWithCache[K comparable, V any](ctx context.Context, cache *expirable.LRU[K, V], query func(context.Context, K) (V, error), params K) (V, error) { //nolint:ireturn
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.Bool("dao.cache.available", true))

	if val, ok := cache.Get(params); ok {
		span.SetAttributes(attribute.Bool("dao.cache.hit", true))

		return val, nil
	}

	val, err := query(ctx, params)
	if err != nil {
		return val, fmt.Errorf("fetch from database on cache miss: %w", err)
	}

	span.SetAttributes(attribute.Bool("dao.cache.hit", false))

	evicted := cache.Add(params, val)
	span.SetAttributes(attribute.Bool("dao.cache.placed", true), attribute.Bool("dao.cache.evicted", evicted))

	return val, nil
}
