package dao

import (
	"context"
	"database/sql"
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
func New(db *sql.DB) *Queries {
	return &Queries{
		dbc: dbc.New(db),
	}
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
