package dao

import (
	"context"
	"sync"

	"github.com/pubgolf/pubgolf/api/internal/lib/telemetry"
)

type asyncResult struct {
	query func(ctx context.Context)
}

// Run invokes the underlying query.
func (ar *asyncResult) Run(ctx context.Context, wg *sync.WaitGroup) {
	runAsync(ctx, wg, ar.query)
}

func runAsync(ctx context.Context, wg *sync.WaitGroup, fn func(ctx context.Context)) {
	defer telemetry.FnSpan(&ctx)()

	wg.Add(1)

	go func() {
		defer wg.Done()
		fn(ctx)
	}()
}
