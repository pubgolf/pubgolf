package dao

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao/internal/dbc"
	"go.opentelemetry.io/otel"
)

// Queries holds references to all data stores and provides query methods.
type Queries struct {
	dbc dbc.Querier
}

// New returns a concrete implementation of `QueryProvider`.
func New(ctx context.Context, db *sql.DB) (*Queries, error) {
	q, err := dbc.Prepare(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("prepare dbc queries: %w", err)
	}

	return &Queries{
		dbc: q,
	}, nil
}

// daoSpan annotates a DAO method with a span for tracing.
func daoSpan(ctx *context.Context) func() {
	name := "dao.AnonymousQuery"
	if pc, _, _, ok := runtime.Caller(1); ok {
		name = "dao." + strings.Split(filepath.Base(runtime.FuncForPC(pc).Name()), ".")[2]
	}
	newCtx, span := otel.Tracer("").Start(*ctx, name)
	*ctx = newCtx
	return func() { span.End() }
}
