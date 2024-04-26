package dao

import (
	"context"
	"sync"
)

type asyncResult struct {
	query func(ctx context.Context)
}

// Run invokes the underlying query.
func (ar *asyncResult) Run(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)

	go func() {
		defer wg.Done()
		ar.query(ctx)
	}()
}
